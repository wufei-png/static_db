#include <stdio.h>
#include <string.h>
#include <pthread.h>
#include "decoder.h"
#include "ringbuffer.h"

// typedef struct frame_queue_t {
// } frame_queue_t;

typedef struct decoder_output_config_t {
    int nb_surfaces;
    // Deinterlacing mode
    int deint_mode_current;
    int drop_second_field;
} decoder_output_config_t;

ringBuffer_typedef(CuvidParsedFrame, frame_buffer_t);

typedef struct decoder_context_t {
    CUvideoparser parser;
    CUvideodecoder decoder;

    CUVIDPARSERPARAMS params;
    CUVIDEOFORMATEX extra;

    decoder_output_config_t config;

    int inited;
    // received EOF Packet
    int decoder_flushing;

    CUresult internal_error;

    // protect frame_buffer
    pthread_mutex_t mu;
    // internal decoded frame ring buffer
    frame_buffer_t frame_buffer;

    int coded_width;
    int coded_height;
    int width;
    int height;
    cudaVideoCodec codec_type;
    cudaVideoChromaFormat chroma_format;
} decoder_context_t;

static int sequenceCallback(void *data, CUVIDEOFORMAT *format) {
    decoder_context_t *ctx = data;
    CUVIDDECODECREATEINFO cuinfo;

    memset(&cuinfo, 0, sizeof(cuinfo));

    // apply cropping
    cuinfo.display_area.left = format->display_area.left; // + ctx->crop.left;
    cuinfo.display_area.top = format->display_area.top; // + ctx->crop.top;
    cuinfo.display_area.right = format->display_area.right; // - ctx->crop.right;
    cuinfo.display_area.bottom = format->display_area.bottom; // - ctx->crop.bottom;

    int width = cuinfo.display_area.right - cuinfo.display_area.left;
    int height = cuinfo.display_area.bottom - cuinfo.display_area.top;
    if (ctx->inited) {
        // XXX not necessary
        if (ctx->width != width || ctx->height != height) {
            fprintf(stderr, "WARNING: re-initialize decoder with different dimension not supported: %dx%d to %dx%d\n",
                    ctx->width, ctx->height, width, height);
            ctx->internal_error = CUDA_ERROR_INVALID_VALUE;
            return 0;
        }
        if (ctx->width == width
            && ctx->height == height
            && ctx->codec_type == format->codec
            && ctx->chroma_format == format->chroma_format
            && ctx->coded_width == format->coded_width
            && ctx->coded_height == format->coded_height
           ) {
            return 1;
        }
        fprintf(stderr, "INFO: re-initializing cuvid decoder\n");
        ctx->internal_error = cuvidDestroyDecoder(ctx->decoder);
        if (ctx->internal_error != 0) {
            fprintf(stderr, "WARNING: failed to destory previous decoder: %d\n", ctx->internal_error);
            return 0;
        }
        ctx->decoder = NULL;
    }

    ctx->internal_error = 0;

    ctx->coded_width = cuinfo.ulWidth = format->coded_width;
    ctx->coded_height = cuinfo.ulHeight = format->coded_height;

    ctx->width = width;
    ctx->height = height;

    // target width/height need to be multiples of two
    cuinfo.ulTargetWidth = ctx->width = (ctx->width + 1) & ~1;
    cuinfo.ulTargetHeight = ctx->height = (ctx->height + 1) & ~1;

    // aspect ratio conversion, 1:1, depends on scaled resolution
    cuinfo.target_rect.left = 0;
    cuinfo.target_rect.top = 0;
    cuinfo.target_rect.right = cuinfo.ulTargetWidth;
    cuinfo.target_rect.bottom = cuinfo.ulTargetHeight;

    fprintf(stderr, "cuvid: initialize new decoder: size=%dx%d, format=%d\n", width, height, format->chroma_format);
    if (format->chroma_format != cudaVideoChromaFormat_420) {
        fprintf(stderr, "WARNING: Chroma formats other than 420 are not supported\n");
        ctx->internal_error = CUDA_ERROR_INVALID_VALUE;
        return 0;
    }
    ctx->chroma_format = format->chroma_format;

    cuinfo.CodecType = ctx->codec_type = format->codec;
    cuinfo.ChromaFormat = format->chroma_format;

    // nvidia decoder current only support NV12
    cuinfo.OutputFormat = cudaVideoSurfaceFormat_NV12;

    // Maximum number of internal decode surfaces
    // In order to minimize decode latencies, there should be always at least 2
    // pictures in the decode queue at any time, in order to make sure that all
    // decode engines are always busy.
    // https://www.ffmpeg.org/doxygen/3.2/group__VIDEO__DECODER.html#ga118a3e1fc92b013f4b564c53d24f0023
    cuinfo.ulNumDecodeSurfaces = ctx->config.nb_surfaces;
    // There is a limit to how many pictures can be mapped simultaneously (ulNumOutputSurfaces)
    cuinfo.ulNumOutputSurfaces = 1;

    // Use dedicated video engines directly
    cuinfo.ulCreationFlags = cudaVideoCreate_PreferCUVID;
    // Must be 0 (only 8-bit supported)
    cuinfo.bitDepthMinus8 = format->bit_depth_luma_minus8;
    // Deinterlacing mode
    cuinfo.DeinterlaceMode = ctx->config.deint_mode_current;

    CUresult res = cuvidCreateDecoder(&ctx->decoder, &cuinfo);
    if (res != 0) {
        ctx->internal_error = res;
        fprintf(stderr, "WARNING: cuvidCreateDecoder failed: %d\n", res);
        return 0;
    }
    ctx->inited = 1;
    ctx->decoder_flushing = 0;
    return 1;
}

static int decodeCallback(void *data, CUVIDPICPARAMS *params) {
    decoder_context_t *ctx = data;
    if (!ctx->inited) {
        fprintf(stderr, "WARNING: ctx is not inited in decodeCallback\n");
        ctx->internal_error = CUDA_ERROR_INVALID_VALUE;
        return 0;
    }
    // fprintf(stderr, "HERE decodePicture %d\n", params->CurrPicIdx);
    CUresult res = cuvidDecodePicture(ctx->decoder, params);
    if (res != 0) {
        ctx->internal_error = res;
        fprintf(stderr, "WARNING: cuvidDecodePicture failed: %d\n", res);
        return 0;
    }
    return 1;
}

static int displayCallback(void *data, CUVIDPARSERDISPINFO *dispinfo) {
    decoder_context_t *ctx = data;
    if (!ctx->inited) {
        fprintf(stderr, "WARNING: ctx is not inited in displayCallback\n");
        ctx->internal_error = CUDA_ERROR_INVALID_VALUE;
        return 0;
    }
    // fprintf(stderr, "HERE displayPicture\n");

    CuvidParsedFrame parsed_frame;
    memset(&parsed_frame, 0, sizeof(parsed_frame));

    parsed_frame.dispinfo = *dispinfo;
    parsed_frame.desc.width = ctx->width;
    parsed_frame.desc.height = ctx->height;
    parsed_frame.desc.format = cudaVideoSurfaceFormat_NV12;

    pthread_mutex_lock(&ctx->mu);
    ctx->internal_error = 0;

    /*
    if (isBufferFull((&ctx->frame_buffer))) {
        fprintf(stderr, "QUEUE FULL\n");
        return 0;
    }
    */

    // always true
    if (ctx->config.deint_mode_current == cudaVideoDeinterlaceMode_Weave) {
        bufferWrite(&ctx->frame_buffer, parsed_frame);
    } else {
        parsed_frame.is_deinterlacing = 1;
        bufferWrite((&ctx->frame_buffer), parsed_frame);
        if (!ctx->config.drop_second_field) {
            parsed_frame.second_field = 1;
            bufferWrite((&ctx->frame_buffer), parsed_frame);
        }
    }
    pthread_mutex_unlock(&ctx->mu);

    return 1;
}

CUresult init_decoder_context(decoder_context_t **out, CUVIDPARSERPARAMS* params, CUVIDEOFORMATEX *extra_params) {
    decoder_context_t *ctx = calloc(1, sizeof(decoder_context_t));
    ctx->params = *params;
    ctx->extra = *extra_params;
    ctx->params.pExtVideoInfo = &ctx->extra;

    ctx->params.pUserData = ctx;
    // cuvidCreateDecoder Called before decoding frames and/or whenever there is a format change
    // prepare decoder creation param preparation and call cuvidCreateDecoder
    ctx->params.pfnSequenceCallback = sequenceCallback;
    // pfnDecodePicture Called when a picture is ready to be decoded (decode order)
    // just do decode param preparation, and call cuvidDecodePicture. it can not get decoded frame
    ctx->params.pfnDecodePicture = decodeCallback;
    // pfnDisplayPicture Called whenever a picture is ready to be displayed (display order)
    // will get the decoded frame and push it to `frame_buffer`
    ctx->params.pfnDisplayPicture = displayCallback;

    ctx->config.nb_surfaces = ctx->params.ulMaxNumDecodeSurfaces;
    // cudaVideoDeinterlaceMode_Weave is Weave both fields (no deinterlacing)
    ctx->config.deint_mode_current = cudaVideoDeinterlaceMode_Weave;

    CUresult res = CUDA_SUCCESS;
    bufferInit(ctx->frame_buffer,ctx->config.nb_surfaces,CuvidParsedFrame);
    if (pthread_mutex_init(&ctx->mu, NULL) != 0) {
        bufferDestroy(&ctx->frame_buffer);
        free(ctx);
        return CUDA_ERROR_UNKNOWN;
    }

    res = cuvidCreateVideoParser(&ctx->parser, &ctx->params);
    if (res != 0) {
        fprintf(stderr, "WARNING: failed to call cuvidCreateVideoParser in init_decoder_context: %d\n", res);
        goto err;
    }

    CUVIDSOURCEDATAPACKET pkt = {0};
    pkt.payload = (const unsigned char*)ctx->extra.raw_seqhdr_data;
    pkt.payload_size = ctx->extra.format.seqhdr_data_length;
    if (pkt.payload && pkt.payload_size) {
        res = cuvidParseVideoData(ctx->parser, &pkt);
        if (res != 0) {
            fprintf(stderr, "WARNING: failed to call cuvidParseVideoData in init_decoder_context: %d\n", res);
            goto err;
        }
    }

    *out = ctx;
    return CUDA_SUCCESS;
err:
    close_decoder_context(ctx);
    return res;
}

CUresult close_decoder_context(decoder_context_t *ctx) {
    if (ctx->inited) {
        CUresult res = cuvidDestroyDecoder(ctx->decoder);
        if (res != 0) {
            fprintf(stderr, "WARNING: failed to destory decoder: %d\n", res);
        }
    }
    CUresult res = CUDA_SUCCESS;
    if (ctx->parser) {
        res = cuvidDestroyVideoParser(ctx->parser);
    }
    if (res !=0) {
        fprintf(stderr, "WARNING: failed to destory parser: %d\n", res);
    }

    pthread_mutex_destroy(&ctx->mu);
    bufferDestroy(&ctx->frame_buffer);
    free(ctx);
    return res;
}

CUresult decoder_push_packet(decoder_context_t *ctx, const void *payload, unsigned long payload_size, long long timestamp, unsigned int flags) {
    pthread_mutex_lock(&ctx->mu);
    int is_flusing = ctx->decoder_flushing;
    if (is_flusing && payload && payload_size) {
        return (CUresult)DECODER_CUSTOM_ERROR_EOF;
    }
    if ((bufferSize(ctx->frame_buffer) + 2 > ctx->config.nb_surfaces) && payload && payload_size) {
        // fprintf(stderr, "XXX1 %d\n", bufferSize(ctx->frame_buffer));
        pthread_mutex_unlock(&ctx->mu);
        return (CUresult)DECODER_CUSTOM_ERROR_AGAIN;
    }

    CUVIDSOURCEDATAPACKET pkt = {0};
    if (payload && payload_size) {
        pkt.payload = (const unsigned char*)payload;
        pkt.payload_size = payload_size;
        pkt.flags = flags;
        pkt.timestamp = timestamp;
    } else {
        pkt.flags = CUVID_PKT_ENDOFSTREAM;
        ctx->decoder_flushing = 1;
    }

    pthread_mutex_unlock(&ctx->mu);
    CUresult res = cuvidParseVideoData(ctx->parser, &pkt);
    if (res != 0){
        fprintf(stderr, "WARNING: failed to call cuvidParseVideoData in decoder_push_packet: %d\n", res);
    }
    return res;
}

CUresult decoder_pull_frame(decoder_context_t *ctx, CuvidParsedFrame *frame, int *more) {
    pthread_mutex_lock(&ctx->mu);
    CUresult err = ctx->internal_error;
    if (err) {
        goto done;
    }
    if (!isBufferEmpty(&ctx->frame_buffer)) {
        bufferRead(&ctx->frame_buffer, *frame);
        *more = 1;
    } else if (ctx->decoder_flushing) {
        err = (CUresult)DECODER_CUSTOM_ERROR_EOF;
        goto done;
    } else {
        *more = 0;
        goto done;
    }
done:
    pthread_mutex_unlock(&ctx->mu);
    return err;
}

CUvideodecoder decoder_get_internal(decoder_context_t *ctx) {
    return ctx->decoder;
}

