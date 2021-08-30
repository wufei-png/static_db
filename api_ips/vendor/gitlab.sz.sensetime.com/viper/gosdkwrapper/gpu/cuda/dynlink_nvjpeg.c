#include <stdlib.h>
#include <dlfcn.h>
#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <pthread.h>
#include <cuda.h>
#include "dynlink_nvjpeg.h"

static char __NVJPEGLibName[64] = "libnvjpeg.so";

typedef void *DLLDRIVER;

static nvjpegStatus_t LOAD_LIBRARY(DLLDRIVER *pInstance)
{
    *pInstance = dlopen(__NVJPEGLibName, RTLD_NOW);
    // symlink to libnvcuvid.so is not available in nvidia-docker
    if (*pInstance == NULL) {
        strncat(__NVJPEGLibName, ".1", 3);
        *pInstance = dlopen(__NVJPEGLibName, RTLD_NOW);
    }

    if (*pInstance == NULL)
    {
        fprintf(stderr, "dlopen \"%s\" failed: %s\n", __NVJPEGLibName, dlerror());
        return NVJPEG_STATUS_NOT_INITIALIZED;
    }

    return NVJPEG_STATUS_SUCCESS;
}

#define STRINGIFY(X) #X
#define GET_PROC_EX(name, alias, required)                              \
    defaultFuncs.alias = (t##name *)dlsym(DriverLib, #name);                        \
    if (defaultFuncs.alias == NULL && required) {                                    \
        printf("Failed to find required function \"%s\" in %s\n",       \
               #name, __NVJPEGLibName);                                  \
        return NVJPEG_STATUS_NOT_INITIALIZED;                                      \
    }

#define GET_PROC_EX_V2(name, alias, required)                           \
    alias = (t##name *)dlsym(DriverLib, STRINGIFY(name##_v2));         \
    if (alias == NULL && required) {                                    \
        printf("Failed to find required function \"%s\" in %s\n",       \
               STRINGIFY(name##_v2), __DriverLibName);                    \
        return NVJPEG_STATUS_NOT_INITIALIZED;                                      \
    }

#define CHECKED_CALL(call)              \
    do {                                \
        nvjpegStatus_t result = (call);       \
        if (NVJPEG_STATUS_SUCCESS != result) {   \
            return result;              \
        }                               \
    } while(0)

#define GET_PROC_REQUIRED(name) GET_PROC_EX(name,name,1)
#define GET_PROC_OPTIONAL(name) GET_PROC_EX(name,name,0)
#define GET_PROC(name)          GET_PROC_REQUIRED(name)
#define GET_PROC_V2(name)       GET_PROC_EX_V2(name,name,1)


nvjpegFunctions_t *nvjpegFunctions;
static nvjpegFunctions_t defaultFuncs;

nvjpegStatus_t NVJPEGAPI nvjpegInit(unsigned int flags) {
    DLLDRIVER DriverLib;

    CHECKED_CALL(LOAD_LIBRARY(&DriverLib));
    GET_PROC(nvjpegGetProperty);
    GET_PROC(nvjpegGetImageInfo);

    GET_PROC(nvjpegCreateEx);
    GET_PROC(nvjpegDestroy);

    GET_PROC(nvjpegJpegStateCreate);
    GET_PROC(nvjpegJpegStateDestroy);

    GET_PROC(nvjpegDecodeBatchedInitialize);
    GET_PROC(nvjpegDecodeBatched);
    GET_PROC(nvjpegDecodeJpegHost);
    GET_PROC(nvjpegDecodeJpegTransferToDevice);
    GET_PROC(nvjpegDecodeJpegDevice);

    GET_PROC(nvjpegJpegStreamParse);
    GET_PROC(nvjpegStateAttachDeviceBuffer);
    GET_PROC(nvjpegDecoderCreate);
    GET_PROC(nvjpegDecoderDestroy);
    GET_PROC(nvjpegJpegStreamCreate);
    GET_PROC(nvjpegJpegStreamDestroy);
    GET_PROC(nvjpegDecoderStateCreate);
    GET_PROC(nvjpegDecodeParamsCreate);
    GET_PROC(nvjpegDecodeParamsDestroy);
    GET_PROC(nvjpegBufferDeviceCreate);
    GET_PROC(nvjpegBufferDeviceDestroy);
    GET_PROC(nvjpegDecodeParamsSetOutputFormat);

    nvjpegFunctions = &defaultFuncs;
    return NVJPEG_STATUS_SUCCESS;
}

typedef struct nvjpegMemoryPool_t {
    nvjpegDevAllocator_t alloc;
    pthread_mutex_t mu;
    CUdeviceptr *ptrs;
    int *refs;
    int count;
    unsigned int blobsize;
} nvjpegMemoryPool_t;

static nvjpegMemoryPool_t gMemPool;

nvjpegDevAllocator_t* NVJPEGAPI _nvjpegGetMemoryPoolAllocator() {
    if (gMemPool.count <= 0) return NULL;
    return &gMemPool.alloc;
}

static inline int _fallback_malloc(void **ptr, size_t size) {
    CUdeviceptr p = 0;
    fprintf(stderr, "nvjpeg mempool: fallback malloc: %ld\n", size);
    CUresult r = cuMemAlloc(&p, size);
    if (r!=0) return r;
    *ptr = (void*)p;
    return 0;
}

static int _pool_malloc(void** ptr, size_t size) {
    int i;
    if (size > gMemPool.blobsize) {
        return _fallback_malloc(ptr, size);
    }
    pthread_mutex_lock(&gMemPool.mu);
    for (i = 0; i < gMemPool.count; i++) {
        if (gMemPool.refs[i] <= 0) {
            fprintf(stderr, "nvjpeg mempool: checkout size %ld\n", size);
            gMemPool.refs[i]++;
            *ptr = (void*)gMemPool.ptrs[i];
            pthread_mutex_unlock(&gMemPool.mu);
            return 0;
        }
    }
    pthread_mutex_unlock(&gMemPool.mu);
    // no empty slots
    return _fallback_malloc(ptr, size);
}

static int _pool_free(void *ptr) {
    int i;
    pthread_mutex_lock(&gMemPool.mu);
    for (i = 0; i < gMemPool.count; i++) {
        if ((void*)(gMemPool.ptrs[i]) == ptr) {
            assert(gMemPool.refs[i] > 0);
            gMemPool.refs[i]--;
            pthread_mutex_unlock(&gMemPool.mu);
            return 0;
        }
    }
    pthread_mutex_unlock(&gMemPool.mu);
    return cuMemFree((CUdeviceptr)ptr);
}

nvjpegStatus_t NVJPEGAPI _nvjpegInitMemoryPool(unsigned int blobsize, int count) {
    int i, j;
    if (count <= 0) return NVJPEG_STATUS_INVALID_PARAMETER;
    if (gMemPool.count > 0) return NVJPEG_STATUS_NOT_INITIALIZED;
    CUdeviceptr *ptrs = calloc(sizeof(CUdeviceptr*), count);
    if (!ptrs) return NVJPEG_STATUS_ALLOCATOR_FAILURE;
    int *refs = calloc(sizeof(int), count);
    if (!refs) {
        free(ptrs);
        return NVJPEG_STATUS_ALLOCATOR_FAILURE;
    }
    for (i = 0; i < count; i++) {
        CUresult r = cuMemAlloc(&ptrs[i], blobsize);
        if (r != 0) {
            for (j = 0; j < count; j++) {
                if (ptrs[j]) cuMemFree(ptrs[j]);
            }
            free(ptrs);
            free(refs);
            return NVJPEG_STATUS_ALLOCATOR_FAILURE;
        }
    }
    pthread_mutex_init(&gMemPool.mu, NULL);
    gMemPool.alloc.dev_malloc = _pool_malloc;
    gMemPool.alloc.dev_free = _pool_free;
    gMemPool.ptrs = ptrs;
    gMemPool.refs = refs;
    gMemPool.count = count;
    gMemPool.blobsize = blobsize;
    return NVJPEG_STATUS_SUCCESS;
}

nvjpegStatus_t NVJPEGAPI _nvjpegDestroyMemoryPool() {
    int i;
    for (i = 0; i < gMemPool.count; i++) {
        if (gMemPool.ptrs[i]) cuMemFree(gMemPool.ptrs[i]);
    }
    free(gMemPool.ptrs);
    free(gMemPool.refs);
    pthread_mutex_destroy(&gMemPool.mu);
    gMemPool.count = 0;
    gMemPool.blobsize = 0;
    return NVJPEG_STATUS_SUCCESS;
}


