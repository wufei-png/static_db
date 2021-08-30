/**
 * @file stpu_device.h
 */

#ifndef __STPUHAL_DEVICE_H__
#define __STPUHAL_DEVICE_H__

#include <unistd.h>
#include <stdint.h>

#ifdef __cplusplus
#define STPU_HAL_API extern "C" __attribute__((visibility("default")))
#else
#define STPU_HAL_API __attribute__((visibility("default")))
#endif

typedef int32_t halError;

/**
 * @brief memory category enum definition.
 */
typedef enum stpuHalMemType {
    stpuHalMemStpu = 0, /**< Device local memory */
    stpuHalMemStupCacheable, /**< Device local memory cacheable */
    stpuHalMemSoc, /**< Host top memory */
    stpuHalMemSocCacheable, /**< Host top memory cacheable */
    stpuHalMemHost, /**< Host memory (don't support)*/
    stpuHalMemHostCacheable, /**< Host memory cacheable (don't support)*/
} StpuHalMemType;

/**
 * @brief memcpy category enum definition.
 */
typedef enum stpuHalMemcpyType {
    stpuHalMemcpyHostToStpu = 1, /**< Memory transimssion form host to device local */
    stpuHalMemcpyHostToSoc, /**< Memory transimssion form host to device top */
    stpuHalMemcpyStpuToHost, /**< Memory transimssion form device local to host */
    stpuHalMemcpyStpuToStpu, /**< Memory transimssion form device local to device local */
    stpuHalMemcpyStpuToSoc, /**< Memory transimssion form device local to device top */
    stpuHalMemcpySocToHost, /**< Memory transimssion form device top to host */
    stpuHalMemcpySocToStpu, /**< Memory transimssion form device top to device local */
    stpuHalMemcpySocToSoc, /**< Memory transimssion form device top to device top */
    stpuHalMemcpyUnkown,
} StpuHalMemcpyType;

/**
 * @brief Get Current Thread Device ID
 *
 * @param[out] id device id
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalCurrentDeviceID(uint32_t *id);

/**
 * @brief Set Current Thread Device ID
 *
 * @param[in] id device id
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalSetDeviceID(uint32_t *id);

/**
 * @brief Stpu Device Initialization
 *
 * @param[in] dev_id device id
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalDeviceOpen(uint32_t dev_id);

/**
 * @brief Stpu Device Close
 *
 * @param[in] dev_id device id
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalDeviceClose(uint32_t dev_id);

/**
 * @brief Memory malloc
 *
 * @param[in] dev_id device id
 * @param[out] memptr assigned address
 * @param[in] size expectation size
 * @param[in] type memory kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalMalloc(uint32_t dev_id, void **memptr, size_t size,
                                    StpuHalMemType type);

/**
 * @brief Memory calloc
 *
 * @param[in] dev_id device id
 * @param[out] memptr assigned address
 * @param[in] num nmemb elements
 * @param[in] size elements size
 * @param[in] type memory kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalCalloc(uint32_t dev_id, void **memptr, size_t num, size_t size,
                                    StpuHalMemType type);

/**
 * @brief Memory realloc
 *
 * @param[in] dev_id device id
 * @param[out] memptr vaddr
 * @param[in] size expectation size
 * @param[in] type memory kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalRealloc(uint32_t dev_id, void **memptr, size_t size,
                                     StpuHalMemType type);

/**
 * @brief Memory free
 *
 * @param[in] dev_id device id
 * @param[in] type memory kind
 * @param[out] memptr destination address
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalFree(uint32_t dev_id, void *memptr, StpuHalMemType type);

/**
 * @brief Memset
 *
 * @param[in] dev_id device id
 * @param[in] dst destination address
 * @param[in] c set value
 * @param[in] size memory size
 * @param[in] type memory kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalMemset(uint32_t dev_id, void *dst, int32_t c, size_t size,
                                    StpuHalMemType type);

/**
 * @brief Memcory copy
 *
 * @param[in] dev_id device id
 * @param[in] dst destination address
 * @param[in] src source address
 * @param[in] size size
 * @param[in] type memcory copy kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalMemcopy(uint32_t dev_id, void *dst, void *src, size_t size,
                                     StpuHalMemcpyType type);

/**
 * @brief Memcory copy2D
 *
 * @param[in] dev_id device id
 * @param[in] dst destination address
 * @param[in] dpitch destination stride
 * @param[in] src source address
 * @param[in] spitch source stride
 * @param[in] width width
 * @param[in] height height
 * @param[in] type memcory copy kind
 * @return
 *   @retval halError halSuccess for succeed, otherwise return error code
 */
STPU_HAL_API halError stpuHalMemcopy2D(uint32_t dev_id, void *dst, size_t dpitch, void *src,
                                       size_t spitch, size_t width, size_t height,
                                       StpuHalMemcpyType type);

#endif