/**
 * @file stpu_errno.h
 */

#ifndef __STPUHAL_ERRNO_H__
#define __STPUHAL_ERRNO_H__

#include <stdint.h>

#ifdef __cplusplus
#define STPU_HAL_API extern "C" __attribute__((visibility("default")))
#else
#define STPU_HAL_API __attribute__((visibility("default")))
#endif

typedef int32_t halError;

/**
 * @brief errno enum definition.
 */
typedef enum halErrorCode {
    halSuccess,
    halErrorInitialization,
    halErrorInvalidArg,
    halErrorInvalidMsg,
    halErrorInvalidValue,
    halErrorInvalidStatus,
    halErrorInvalidDevice,
    halErrorInvalidModel,
    halErrorInvalidFileFormat,
    halErrorInvalidConnectType,
    halErrorMemoryAllocation,
    halErrorSocketFail,
    halErrorNotExist,
    halErrorNullPtr,
    halErrorTimeout,
    halErrorOutOfMemory,
    halErrorFileNotFound,
    halErrorCreateNNFail,
    halErrorNNForward,
    halErrorNNRelease,
    halErrorLoadVpuFail,
    halErrorCreatCodecFail,
    halErrorDecodePacketFail,
    halErrorDestoryCodecFail,
    halErrorEAGAIN,
    halErrorEOF,
    halErrorNotSupport,
    halErrorResizeFail,
    halErrorUnknown
} HalErrorCode;

/**
 * @brief returns a pointer to a string that describes the error code
 *
 * @param[in] halerrno HalErrorCode
 * @return
 *   @retval string describing error number
 */
STPU_HAL_API const char *stpuHalStrError(halError halerrno);

#endif