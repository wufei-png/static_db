
#ifndef INCLUDE_GPU_SEARCH_CORE_H_
#define INCLUDE_GPU_SEARCH_CORE_H_

#include <stddef.h>

#ifdef _MSC_VER
#ifdef __cplusplus
#ifdef GPU_SEARCH_CORE_EXPORTS
#define GPU_SEARCH_CORE_API extern "C" __declspec(dllexport)
#else
#define GPU_SEARCH_CORE_API extern "C" __declspec(dllimport)
#endif
#else
#ifdef GPU_SEARCH_CORE_EXPORTS
#define GPU_SEARCH_CORE_API __declspec(dllexport)
#else
#define GPU_SEARCH_CORE_API __declspec(dllimport)
#endif
#endif
#else /* _MSC_VER */
#ifdef __cplusplus
#ifdef GPU_SEARCH_CORE_EXPORTS
#define GPU_SEARCH_CORE_API extern "C" __attribute__((visibility("default")))
#else
#define GPU_SEARCH_CORE_API extern "C"
#endif
#else
#ifdef GPU_SEARCH_CORE_EXPORTS
#define GPU_SEARCH_CORE_API __attribute__((visibility("default")))
#else
#define GPU_SEARCH_CORE_API
#endif
#endif
#endif

#define GSC_OK 0
#define GSC_INVALID_HANDLE (-1)
#define GSC_FILE_NOT_FOUND (-2)
#define GSC_NOMEMORY (-3)
#define GSC_INVALID_ARGUMENT (-4)
#define GSC_BAD_INDEX (-5)
#define GSC_SYSCALL_ERROR (-6)

typedef void *gsc_gpu_index_t;
typedef void *gsc_cpu_index_t;
typedef void *gsc_index_shards_t;
typedef void *gsc_resource_t;

typedef struct
{
    // Index
    int dimension;
    long index_size;
    int is_trained;
    // GpuIndexIVF
    int nlist;
    int nprobe;
    int max_list_size;
    // GpuIndexIVFPQ
    int subQuantizers;
    int bitsPerCode;
} gsc_index_status_t;

typedef struct
{
    // Index
    int dimension;
    // GpuIndexIVF
    int nlist;
    int nprobe;
    // GpuIndexIVFPQ
    int subQuantizers;
    int bitsPerCode;
} gsc_index_config_t;

typedef struct
{
    size_t total_memory_size;
    size_t free_memory_size;
} gsc_gpu_info_t;

typedef struct
{
    unsigned int time_range[2];
    // bits from 0 to 127
    unsigned char camera_mask[128 / 8];
} gsc_timespace_filter_t;

typedef struct gsc_device_properties
{
    int major;
    int minor;
    char name[256];
} gsc_device_properties;

//GPU Search
GPU_SEARCH_CORE_API int gsc_get_current_device_id();

GPU_SEARCH_CORE_API int gsc_get_gpu_info(int device_id, gsc_gpu_info_t *info);

GPU_SEARCH_CORE_API int gsc_create_gpu_resource(int device_id, unsigned long temp_memory_size, float temp_memory_fraction, gsc_resource_t *res);
GPU_SEARCH_CORE_API int gsc_destroy_gpu_resource(gsc_resource_t resource);

GPU_SEARCH_CORE_API int gsc_init_gpu_index(gsc_resource_t resource, gsc_index_config_t *config, gsc_gpu_index_t *idx);
GPU_SEARCH_CORE_API int gsc_load_gpu_index(gsc_resource_t resource, const char *path, int device_id, gsc_gpu_index_t *idx);
GPU_SEARCH_CORE_API int gsc_create_index_shards(int dimension, int threaded, gsc_index_shards_t *shards);

GPU_SEARCH_CORE_API int gsc_free_gpu_index(gsc_gpu_index_t index);
GPU_SEARCH_CORE_API int gsc_free_cpu_index(gsc_cpu_index_t index);
GPU_SEARCH_CORE_API int gsc_free_index_shards(gsc_index_shards_t shards);

GPU_SEARCH_CORE_API int gsc_get_gpu_index_status(gsc_gpu_index_t index, gsc_index_status_t *status);
GPU_SEARCH_CORE_API int gsc_get_gpu_index_ids(gsc_gpu_index_t index, long *ids);
GPU_SEARCH_CORE_API long gsc_get_gpu_max_may_reserve_memory(gsc_gpu_index_t index, long n);

GPU_SEARCH_CORE_API int gsc_train_gpu_index(gsc_gpu_index_t index, long n, const float *x);
GPU_SEARCH_CORE_API int gsc_add_index_batch(gsc_gpu_index_t index, long n, const float *x, const long *ids);
GPU_SEARCH_CORE_API int gsc_add_gpu_index_to_shards(gsc_gpu_index_t index, gsc_index_shards_t shards);
GPU_SEARCH_CORE_API int gsc_add_shards_to_shards(gsc_index_shards_t from, gsc_index_shards_t to);

GPU_SEARCH_CORE_API int gsc_index_gpu_to_cpu(gsc_gpu_index_t index, gsc_cpu_index_t *cpu_index, int reclaim_gpu_memory);
GPU_SEARCH_CORE_API int gsc_index_cpu_to_gpu(gsc_cpu_index_t index, gsc_resource_t resource, int device_id, gsc_gpu_index_t *gpu_index);
GPU_SEARCH_CORE_API int gsc_write_cpu_index(gsc_cpu_index_t index, const char *filename);

GPU_SEARCH_CORE_API int gsc_search_index(gsc_gpu_index_t index, long n, const float *x, int k, float *distances, long *ids);
GPU_SEARCH_CORE_API int gsc_search_index_with_timespace_filter(gsc_gpu_index_t index, long n, const float *x, int k, gsc_timespace_filter_t *filter, float *distances, long *ids);

GPU_SEARCH_CORE_API int gsc_search_index_shards(gsc_index_shards_t shards, long n, const float *x, int k, float *distances, long *ids);

//CPU Search
GPU_SEARCH_CORE_API int gsc_init_cpu_index(gsc_index_config_t *config, gsc_cpu_index_t *idx);
GPU_SEARCH_CORE_API int gsc_load_cpu_index(const char *path, gsc_cpu_index_t *idx);

GPU_SEARCH_CORE_API int gsc_get_cpu_index_status(gsc_cpu_index_t index, gsc_index_status_t *status);
GPU_SEARCH_CORE_API int gsc_get_cpu_index_ids(gsc_cpu_index_t index, long *ids);

GPU_SEARCH_CORE_API int gsc_train_cpu_index(gsc_cpu_index_t index, long n, const float *x);
GPU_SEARCH_CORE_API int gsc_add_cpu_index_batch(gsc_cpu_index_t index, long n, const float *x, const long *ids);
GPU_SEARCH_CORE_API int gsc_remove_cpu_index_ids(gsc_cpu_index_t index, long n, const long *ids);

GPU_SEARCH_CORE_API int gsc_search_cpu_index(gsc_cpu_index_t index, long n, const float *x, int k, float *distances, long *ids);
GPU_SEARCH_CORE_API int gsc_search_cpu_index_with_timespace_filter(gsc_cpu_index_t index, long n, const float *x, int k, gsc_timespace_filter_t *filter, float *distances, long *ids);

GPU_SEARCH_CORE_API int gsc_clone_cpu_index(gsc_cpu_index_t input_index, gsc_cpu_index_t *output_index);

GPU_SEARCH_CORE_API int gsc_get_device_properties(gsc_device_properties *prop, int device);
#endif // INCLUDE_GPU_SEARCH_CORE_H_
