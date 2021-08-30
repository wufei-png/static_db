#ifndef __dynlink_nvjpeg_h
#define __dynlink_nvjpeg_h

#include "library_types.h"
#include "dynlink_cuda.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct CUstream_st *cudaStream_t;
#define NVJPEG_MAX_COMPONENT 4
/* nvJPEG status enums, returned by nvJPEG API */
typedef enum
{
    NVJPEG_STATUS_SUCCESS            = 0,
    NVJPEG_STATUS_NOT_INITIALIZED    = 1,
    NVJPEG_STATUS_INVALID_PARAMETER  = 2,
    NVJPEG_STATUS_BAD_JPEG           = 3,
    NVJPEG_STATUS_JPEG_NOT_SUPPORTED = 4,
    NVJPEG_STATUS_ALLOCATOR_FAILURE  = 5,
    NVJPEG_STATUS_EXECUTION_FAILED   = 6,
    NVJPEG_STATUS_ARCH_MISMATCH      = 7,
    NVJPEG_STATUS_INTERNAL_ERROR     = 8,
    NVJPEG_STATUS_IMPLEMENTATION_NOT_SUPPORTED = 9,
} nvjpegStatus_t;


// Enumeration returned by getImageInfo identifies image chroma subsampling stored inside JPEG input stream
// In the case of NVJPEG_CSS_GRAY only 1 luminance channel is encoded in JPEG input stream
// Otherwise both chroma planes are present
// Initial release support: 4:4:4, 4:2:0, 4:2:2, Grayscale
typedef enum
{
    NVJPEG_CSS_444 = 0,
    NVJPEG_CSS_422 = 1,
    NVJPEG_CSS_420 = 2,
    NVJPEG_CSS_440 = 3,
    NVJPEG_CSS_411 = 4,
    NVJPEG_CSS_410 = 5,
    NVJPEG_CSS_GRAY = 6,
    NVJPEG_CSS_UNKNOWN = -1
} nvjpegChromaSubsampling_t;

// Parameter of this type specifies what type of output user wants for image decoding
typedef enum
{
    NVJPEG_OUTPUT_UNCHANGED   = 0, // return decoded image as it is - write planar output
    NVJPEG_OUTPUT_YUV         = 1, // return planar luma and chroma
    NVJPEG_OUTPUT_Y           = 2, // return luma component only, write to 1-st channel of nvjpegImage_t
    NVJPEG_OUTPUT_RGB         = 3, // convert to planar RGB
    NVJPEG_OUTPUT_BGR         = 4, // convert to planar BGR
    NVJPEG_OUTPUT_RGBI        = 5, // convert to interleaved RGB and write to 1-st channel of nvjpegImage_t
    NVJPEG_OUTPUT_BGRI        = 6  // convert to interleaved BGR and write to 1-st channel of nvjpegImage_t
} nvjpegOutputFormat_t;

// Implementation
// Initial release support: NVJPEG_BACKEND_DEFAULT, NVJPEG_BACKEND_HYBRID
typedef enum 
{
    NVJPEG_BACKEND_DEFAULT = 0,
    NVJPEG_BACKEND_HYBRID  = 1,
    NVJPEG_BACKEND_GPU     = 2,
} nvjpegBackend_t;

// Output descriptor.
// Data that is written to planes depends on output forman
typedef struct
{
    unsigned char * channel[NVJPEG_MAX_COMPONENT];
    unsigned int    pitch[NVJPEG_MAX_COMPONENT];
} nvjpegImage_t;


typedef int (*tDevMalloc)(void**, size_t); 
typedef int (*tDevFree)(void*); 
typedef struct 
{ 
	tDevMalloc dev_malloc; 
	tDevFree dev_free; 
} nvjpegDevAllocator_t;

typedef int (*tPinnedMalloc)(void**, size_t, unsigned int flags);
typedef int (*tPinnedFree)(void*);
typedef struct
{
    tPinnedMalloc pinned_malloc;
    tPinnedFree pinned_free;
} nvjpegPinnedAllocator_t;

struct nvjpegHandle;
typedef struct nvjpegHandle* nvjpegHandle_t;

struct nvjpegJpegState;
typedef struct nvjpegJpegState* nvjpegJpegState_t;

///////////////////////////////////////////////////////////////////////////////////
// Decode parameters //
///////////////////////////////////////////////////////////////////////////////////
// decode parameters structure. Used to set decode-related tweaks
struct nvjpegDecodeParams;
typedef struct nvjpegDecodeParams* nvjpegDecodeParams_t;

//creates decoder implementation
struct nvjpegJpegDecoder;
typedef struct nvjpegJpegDecoder* nvjpegJpegDecoder_t;

///////////////////////////////////////////////////////////////////////////////////
// JPEG stream parameters //
///////////////////////////////////////////////////////////////////////////////////

struct nvjpegJpegStream;
typedef struct nvjpegJpegStream* nvjpegJpegStream_t;

///////////////////////////////////////////////////////////////////////////////////
// NVJPEG buffers //
///////////////////////////////////////////////////////////////////////////////////
struct nvjpegBufferDevice;
typedef struct nvjpegBufferDevice* nvjpegBufferDevice_t;

#define NVJPEGAPI

// returns library's property values, such as MAJOR_VERSION, MINOR_VERSION or PATCH_LEVEL
typedef nvjpegStatus_t NVJPEGAPI tnvjpegGetProperty(libraryPropertyType type, int *value);

// Initalization of nvjpeg handle with additional parameters. This handle is used for all consecutive nvjpeg calls
// IN         backend       : Backend to use. Currently Default or Hybrid (which is the same at the moment) is supported.
// IN         dev_allocator : Pointer to nvjpegDevAllocator. If NULL - use default cuda calls (cudaMalloc/cudaFree)
// IN         pinned_allocator : Pointer to nvjpegPinnedAllocator. If NULL - use default cuda calls (cudaHostAlloc/cudaFreeHost)
// IN         flags         : Parameters for the operation. Must be 0.
// INT/OUT    handle        : Codec instance, use for other calls
typedef nvjpegStatus_t NVJPEGAPI tnvjpegCreateEx(nvjpegBackend_t backend,
        nvjpegDevAllocator_t *dev_allocator,
        nvjpegPinnedAllocator_t *pinned_allocator,
        unsigned int flags,
        nvjpegHandle_t *handle);

// Release the handle and resources.
// IN/OUT     handle: instance handle to release 
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDestroy(nvjpegHandle_t handle);


// Initalization of decoding state
// IN         handle        : Library handle
// INT/OUT    jpeg_handle   : Decoded jpeg image state handle
typedef nvjpegStatus_t NVJPEGAPI tnvjpegJpegStateCreate(nvjpegHandle_t handle, nvjpegJpegState_t *jpeg_handle);

// Release the jpeg image handle.
// INT/OUT    jpeg_handle   : Decoded jpeg image state handle
typedef nvjpegStatus_t NVJPEGAPI tnvjpegJpegStateDestroy(nvjpegJpegState_t jpeg_handle);
// 
// Retrieve the image info, including channel, width and height of each component, and chroma subsampling.
// If less than NVJPEG_MAX_COMPONENT channels are encoded, then zeros would be set to absent channels information
// If the image is 3-channel, all three groups are valid.
// This function is thread safe.
// IN         handle      : Library handle
// IN         data        : Pointer to the buffer containing the jpeg stream data to be decoded. 
// IN         length      : Length of the jpeg image buffer.
// OUT        nComponent  : Number of componenets of the image, currently only supports 1-channel (grayscale) or 3-channel.
// OUT        subsampling : Chroma subsampling used in this JPEG, see nvjpegChromaSubsampling_t
// OUT        widths      : pointer to NVJPEG_MAX_COMPONENT of ints, returns width of each channel. 0 if channel is not encoded  
// OUT        heights     : pointer to NVJPEG_MAX_COMPONENT of ints, returns height of each channel. 0 if channel is not encoded 
typedef nvjpegStatus_t NVJPEGAPI tnvjpegGetImageInfo(
          nvjpegHandle_t handle,
          const unsigned char *data, 
          size_t length,
          int *nComponents, 
          nvjpegChromaSubsampling_t *subsampling,
          int *widths,
          int *heights);
                   

// Decodes single image. Destination buffers should be large enough to be able to store 
// output of specified format. For each color plane sizes could be retrieved for image using nvjpegGetImageInfo()
// and minimum required memory buffer for each plane is nPlaneHeight*nPlanePitch where nPlanePitch >= nPlaneWidth for
// planar output formats and nPlanePitch >= nPlaneWidth*nOutputComponents for interleaved output format.
// 
// IN/OUT     handle        : Library handle
// INT/OUT    jpeg_handle   : Decoded jpeg image state handle
// IN         data          : Pointer to the buffer containing the jpeg image to be decoded. 
// IN         length        : Length of the jpeg image buffer.
// IN         output_format : Output data format. See nvjpegOutputFormat_t for description
// IN/OUT     destination   : Pointer to structure with information about output buffers. See nvjpegImage_t description.
// IN/OUT     stream        : CUDA stream where to submit all GPU work
// 
// \return NVJPEG_STATUS_SUCCESS if successful
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecode(
          nvjpegHandle_t handle,
          nvjpegJpegState_t jpeg_handle,
          const unsigned char *data,
          size_t length, 
          nvjpegOutputFormat_t output_format,
          nvjpegImage_t *destination,
          cudaStream_t stream);

//////////////////////////////////////////////
/////////////// Batch decoding ///////////////
//////////////////////////////////////////////

// Resets and initizlizes batch decoder for working on the batches of specified size
// Should be called once for decoding bathes of this specific size, also use to reset failed batches
// IN/OUT     handle          : Library handle
// INT/OUT    jpeg_handle     : Decoded jpeg image state handle
// IN         batch_size      : Size of the batch
// IN         max_cpu_threads : Maximum number of CPU threads that will be processing this batch
// IN         output_format   : Output data format. Will be the same for every image in batch
//
// \return NVJPEG_STATUS_SUCCESS if successful
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeBatchedInitialize(
          nvjpegHandle_t handle,
          nvjpegJpegState_t jpeg_handle,
          int batch_size,
          int max_cpu_threads,
          nvjpegOutputFormat_t output_format);

// Decodes batch of images. Output buffers should be large enough to be able to store 
// outputs of specified format, see single image decoding description for details. Call to 
// nvjpegDecodeBatchedInitialize() is required prior to this call, batch size is expected to be the same as 
// parameter to this batch initialization function.
// 
// IN/OUT     handle        : Library handle
// INT/OUT    jpeg_handle   : Decoded jpeg image state handle
// IN         data          : Array of size batch_size of pointers to the input buffers containing the jpeg images to be decoded. 
// IN         lengths       : Array of size batch_size with lengths of the jpeg images' buffers in the batch.
// IN/OUT     destinations  : Array of size batch_size with pointers to structure with information about output buffers, 
// IN/OUT     stream        : CUDA stream where to submit all GPU work
// 
// \return NVJPEG_STATUS_SUCCESS if successful
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeBatched(
          nvjpegHandle_t handle,
          nvjpegJpegState_t jpeg_handle,
          const unsigned char *const *data,
          const size_t *lengths, 
          nvjpegImage_t *destinations,
          cudaStream_t stream);

// starts decoding on host and save decode parameters to the state
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeJpegHost(
          nvjpegHandle_t handle,
          nvjpegJpegDecoder_t decoder,
          nvjpegJpegState_t decoder_state,
          nvjpegDecodeParams_t decode_params,
          nvjpegJpegStream_t jpeg_stream);

// hybrid stage of decoding image,  involves device async calls
// note that jpeg stream is a parameter here - because we still might need copy
// parts of bytestream to device
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeJpegTransferToDevice(
          nvjpegHandle_t handle,
          nvjpegJpegDecoder_t decoder,
          nvjpegJpegState_t decoder_state,
          nvjpegJpegStream_t jpeg_stream,
          cudaStream_t stream);

// finishing async operations on the device
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeJpegDevice(
          nvjpegHandle_t handle,
          nvjpegJpegDecoder_t decoder,
          nvjpegJpegState_t decoder_state,
          nvjpegImage_t *destination,
          cudaStream_t stream);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegJpegStreamParse(
          nvjpegHandle_t handle,
          const unsigned char *data,
          size_t length,
          int save_metadata,
          int save_stream,
          nvjpegJpegStream_t jpeg_stream);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegStateAttachDeviceBuffer(
          nvjpegJpegState_t decoder_state,
          nvjpegBufferDevice_t device_buffer);

///////////////////////////////////////////////////////////////////////////////////
// Decoder helper functions //
///////////////////////////////////////////////////////////////////////////////////

// creates decoder state
typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecoderStateCreate(nvjpegHandle_t nvjpeg_handle,
    nvjpegJpegDecoder_t decoder_handle,
    nvjpegJpegState_t* decoder_state);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeParamsCreate(
    nvjpegHandle_t handle,
    nvjpegDecodeParams_t *decode_params);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeParamsDestroy(nvjpegDecodeParams_t decode_params);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecoderCreate(nvjpegHandle_t nvjpeg_handle,
    nvjpegBackend_t implementation,
    nvjpegJpegDecoder_t* decoder_handle);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecoderDestroy(nvjpegJpegDecoder_t decoder_handle);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegBufferDeviceCreate(nvjpegHandle_t handle,
    nvjpegDevAllocator_t* device_allocator,
    nvjpegBufferDevice_t* buffer);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegBufferDeviceDestroy(nvjpegBufferDevice_t buffer);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegJpegStreamCreate(
    nvjpegHandle_t handle,
    nvjpegJpegStream_t *jpeg_stream);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegJpegStreamDestroy(nvjpegJpegStream_t jpeg_stream);

typedef nvjpegStatus_t NVJPEGAPI tnvjpegDecodeParamsSetOutputFormat(
     	nvjpegDecodeParams_t decode_params,
     	nvjpegOutputFormat_t output_format);

typedef struct {
    tnvjpegGetProperty *nvjpegGetProperty;
    tnvjpegCreateEx *nvjpegCreateEx;
    tnvjpegJpegStateCreate *nvjpegJpegStateCreate;
    tnvjpegGetImageInfo *nvjpegGetImageInfo;
    tnvjpegJpegStateDestroy *nvjpegJpegStateDestroy;
    tnvjpegDestroy *nvjpegDestroy;
    tnvjpegDecodeBatchedInitialize *nvjpegDecodeBatchedInitialize;
    tnvjpegDecodeBatched *nvjpegDecodeBatched;
    tnvjpegDecodeJpegHost *nvjpegDecodeJpegHost;
    tnvjpegDecodeJpegTransferToDevice *nvjpegDecodeJpegTransferToDevice;
    tnvjpegDecodeJpegDevice *nvjpegDecodeJpegDevice;
    tnvjpegJpegStreamParse *nvjpegJpegStreamParse;
    tnvjpegStateAttachDeviceBuffer *nvjpegStateAttachDeviceBuffer;
    tnvjpegDecoderCreate *nvjpegDecoderCreate;
    tnvjpegDecoderDestroy *nvjpegDecoderDestroy;
    tnvjpegJpegStreamCreate *nvjpegJpegStreamCreate;
    tnvjpegJpegStreamDestroy *nvjpegJpegStreamDestroy;
    tnvjpegDecoderStateCreate *nvjpegDecoderStateCreate;
    tnvjpegDecodeParamsCreate *nvjpegDecodeParamsCreate;
    tnvjpegDecodeParamsDestroy *nvjpegDecodeParamsDestroy;
    tnvjpegBufferDeviceCreate *nvjpegBufferDeviceCreate;
    tnvjpegBufferDeviceDestroy *nvjpegBufferDeviceDestroy;
    tnvjpegDecodeParamsSetOutputFormat *nvjpegDecodeParamsSetOutputFormat;
} nvjpegFunctions_t;

extern nvjpegFunctions_t *nvjpegFunctions;
extern nvjpegStatus_t NVJPEGAPI nvjpegInit(unsigned int flags);

// implement a simple, global memory pool
extern nvjpegStatus_t NVJPEGAPI _nvjpegInitMemoryPool(unsigned int blobsize, int count);
extern nvjpegStatus_t NVJPEGAPI _nvjpegDestroyMemoryPool();
extern nvjpegDevAllocator_t* NVJPEGAPI _nvjpegGetMemoryPoolAllocator();

#ifdef __cplusplus
}
#endif

#endif
