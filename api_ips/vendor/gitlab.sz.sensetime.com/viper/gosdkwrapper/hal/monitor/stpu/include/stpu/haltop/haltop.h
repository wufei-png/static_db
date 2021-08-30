/**
 * @file haltop.h
 */

#ifndef __stpuHal_H__
#define __stpuHal_H__

#include "version.h"
#include <stpu/stpu_errno.h>

#ifdef __cplusplus
#define STPU_HAL_API extern "C" __attribute__((visibility("default")))
#else
#define STPU_HAL_API __attribute__((visibility("default")))
#endif

/**
 * @brief stpuHalTopInit Remote procedure call initialization
 * @param[in] config (SA: Unix:/var/run/stpuhal ACC: PCIe:dev_id)
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalTopInit(const char *config);

/**
 * @brief stpuHalTopDeinit Remote procedure call deinitialize
 *
 */
STPU_HAL_API void stpuHalTopDeinit(void);

/**
 * @brief Rpc connect (sync stream)
 *
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalConnection(void);

/**
 * @brief Rpc disconnect (sync stream)
 *
 */
STPU_HAL_API void stpuHalDeConnection(void);

/**
 * @brief Rpc get current connection (sync stream)
 * @param[out] context connection context
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalGetCurConnection(void **context);

/**
 * @brief Rpc put current connection (sync stream)
 * @param[in] context connection context
 *
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalPutConnection(void *context);

/**
 * @brief Rpc create async stream
 *
 * @return
 *   @retval stream
 */
STPU_HAL_API struct stream *stpuHalCreateStream(void);

/**
 * @brief Rpc release async stream
 * @param[in] st stream
 *
 */
STPU_HAL_API void stpuHalReleaseStream(struct stream *st);

/**
 * @brief Rpc wait async method complete
 * @param[in] stream stream
 * @param[in] async_hdl asynchronous handle
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuhalMethodSynchronize(struct stream *stream, void *async_hdl);

/**
 * @brief Rpc wait async stream complete
 * @param[in] stream stream
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuhalStreamSynchronize(struct stream *stream);

#endif