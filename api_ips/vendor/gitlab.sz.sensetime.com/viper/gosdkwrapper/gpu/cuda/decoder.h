#ifndef __GO_CUVID_DECODER_H
#define __GO_CUVID_DECODER_H

#include "dynlink_nvcuvid.h"
#include "dynlink_cuviddec.h"

#define DECODER_CUSTOM_ERROR_EOF   (10000 + 1)
#define DECODER_CUSTOM_ERROR_AGAIN (10000 + 2)

#define DECODER_MAX_OUTPUT_PANNEL 4

typedef struct decoded_frame_desc_t {
    int format;
    int width;
    int height;

    unsigned char *data[DECODER_MAX_OUTPUT_PANNEL];
    unsigned int linesize[DECODER_MAX_OUTPUT_PANNEL];
} decoded_frame_desc_t;

typedef struct CuvidParsedFrame
{
    CUVIDPARSERDISPINFO dispinfo;
    int second_field;
    int is_deinterlacing;

    decoded_frame_desc_t desc;
} CuvidParsedFrame;

typedef struct frame_queue_t frame_queue_t;
typedef struct decoder_context_t decoder_context_t;

CUresult init_decoder_context(decoder_context_t **ctx, CUVIDPARSERPARAMS* params, CUVIDEOFORMATEX *extra_params);

CUresult close_decoder_context(decoder_context_t *ctx);

CUresult decoder_push_packet(decoder_context_t *ctx, const void *payload, unsigned long payload_size, long long timestamp, unsigned int flags);

CUresult decoder_pull_frame(decoder_context_t *ctx, CuvidParsedFrame *frame, int *more);

CUvideodecoder decoder_get_internal(decoder_context_t *ctx);

#endif
