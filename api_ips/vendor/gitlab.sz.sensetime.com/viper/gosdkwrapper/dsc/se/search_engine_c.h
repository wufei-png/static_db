#ifndef SEARCH_ENGINE_INCLUDE_SEARCH_ENGINE_C_H_
#define SEARCH_ENGINE_INCLUDE_SEARCH_ENGINE_C_H_

#ifdef __cplusplus
#define SE_API extern "C"
#else
#define SE_API
#endif

#include <stddef.h>
#include <stdint.h>

#define SE_OK 0
#define SE_E_INIT -1
#define SE_E_HANDLE -2
#define SE_E_ALLOC_FAILED -3
#define SE_E_INVALID_VALUE -4
#define SE_E_EXECUTION_FAILED -5
#define SE_E_NOT_SUPPORTED -6

/// @defgroup search engine process context api
/// @{

/// @brief 初始化search engine的整个进程环境
/// @note 整个进程中只需要调用一次
SE_API int se_init(const char *product_name, const char *licence_file);

/// @brief 释放search engine的进程环境资源
SE_API void se_deinit();

/// @brief 日志级别
typedef enum {
  SE_LL_TRACE = 0,
  SE_LL_DEBUG = 1,
  SE_LL_INFO = 2,
  SE_LL_WARNING = 3,
  SE_LL_ERROR = 4
} se_log_level_e;

/// @brief 设置search engine内部的日志打印级别
/// @note log级别设置与获取的接口是进程级别生效的,
/// 需要先调用init接口初骀化searchengine的进程环境
SE_API int se_set_log_level(se_log_level_e level);

/// @brief 获取search engine内部的日志打印级别
SE_API int se_get_log_level(se_log_level_e *level);

/// @}

/// @defgroup Search Engine Context API
/// @{

/// @brief 检索上下文句柄
typedef void *se_context_t;

/// @brief 创建检索上下文句柄
/// @param[in] device_id 为当前检索上下文所绑定的设备号
/// @note 同一个进程或线程环境下多个context可以绑定不同的设备号
/// context可以跨线程创建、使用与销毁，但不支持多线程并发调用，需要调用者来加锁
SE_API int se_context_create(int32_t device_id, se_context_t *context);

/// @brief 销毁检索上下文，销毁后，其创建出来的index都将属于不可用状态
/// 所以在销毁检索上下文前，要提前释放所有index
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API void se_context_destroy(se_context_t context);

/// @brief 设置检索内部计算时，所使用的线程数，默认为1
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_context_set_thread_num(se_context_t context, int32_t num_thread);

/// @brief 设置当前检索上下文中预分配的存储空间的大小
/// @param[in] memory_size 预分配的存储空间大小，单位是byte
/// @note 设置的越大，内部有足够的空间缓存中间结果，则吞吐量就会越大
/// 如果context绑定在设备上，则预分配的为设备存储，否则预分配的为host侧的存储
SE_API int se_context_set_reserved_memory_size(se_context_t context, int32_t memory_size);

/// @brief 获取当前context所绑定的设备内存类型
/// @param[out] mem_type 0代表Host侧内存，1代表Device侧内存
SE_API int se_context_get_mem_type(se_context_t context, int32_t *mem_type);

/// @}

/// @defgroup Search Engine Index API
/// @{

/// @brief 用于进行底库管理的特征索引句柄
/// @note 同一个context下创建的多个index，将会共享一部分内部资源，所以只能串行使用
/// index可以跨线程创建、使用与销毁，但相关的操作接口，不支持并发调用，需要调用者加锁
typedef void *se_index_t;

/// @brief 创建Deepcode索引
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_dc_create(se_context_t context, const char *model_path, int32_t feature_dim,
                              int32_t num_ivf_lists, int32_t num_ivf_probes, se_index_t *index);

/// @brief 创建PQ索引
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_pq_create(se_context_t context, int32_t feature_dim, int32_t num_ivf_lists,
                              int32_t num_ivf_probes, se_index_t *index);

/// @brief 创建Flat索引
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_flat_create(se_context_t context, int32_t feature_dim, int32_t use_int8,
                                se_index_t *index);

/// @brief 创建HNSW索引
/// @param[in] context SearchEngine的线程上下文
/// @param[in] feature_dim 特征的维度
/// @param[in] max_db_size　当前索引最大支持的底库大小
/// @param[in] M 每个点的边数
/// @param[in] ef_construction 插入数据时，搜索的邻近点的个数
/// @note HNSW的算法参数意义与设置指导可以参考： doc/HNSW_ALGO_PARAMS.md
SE_API int se_index_hnsw_create(se_context_t context, int32_t feature_dim, int32_t max_db_size,
                                int32_t M, int32_t ef_construction, int32_t ef, int32_t use_int8,
                                se_index_t *index);

/// @brief 创建精确搜索的Index
/// @param[in] context SearchEngine的线程上下文
/// @param[in] rough_index 精搜内部的粗搜器, 粗搜器内不可以添加特征，创建完成后rough index可以释放
/// @param[in] feature_dim 特征的维度
SE_API int se_index_exact_create(se_context_t context, se_index_t rough_index,
                                 int32_t feature_dim, int32_t use_int8, se_index_t *index);

/// @brief 删除底库索引
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_destroy(se_index_t index);

/// @brief 获取index所绑定的对应的context
/// @param[out] context 批向index所绑定的context的指针
/// @note 由该接口返回的context指针与创建index时的context指向同一个地址，不需要重复释放
SE_API int se_index_get_context(se_index_t index, se_context_t *context);

/// @brief 开启提前对feature id进行时间信息提取
/// @return 调用成功，则返回SE_OK，否则返回错误码
/// @note　开启该功能后，每个特征将额外占用5bytes存储空间，但会加速时空过滤检索时计算的效率
SE_API int se_index_pre_compute_id_filters(se_index_t index);

/// @brief 根据一批采样到的特征来训练index的内部参数
/// @param[in] sampling_features 采样的特征数据，多条特征的数据直接平铺排列
/// @param[in] num 用于训练的特征的数量
/// @note 特征的长度由当前检索index的Context来决定
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_train(se_index_t index, const float *sampling_features, int64_t num);

/// @brief 将特征加入index中,可以分多次插入
/// @param[in] features 加入到index的特征数据，多条特征的数据直接平铺排列
/// @param[in] ids 特征数据对应的特征id，添加到index中的所有id不能有重复
/// @param[in] n　特征的数量
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_add_features(se_index_t index, const float *features, const int64_t *ids,
                                 int64_t n);

/// @brief 根据id来将index中已加入的特征删除
/// @return 调用成功，则返回SE_OK，否则返回错误码
SE_API int se_index_remove_features(se_index_t index, const int64_t *ids, int64_t num);

/// @brief 获取内部所有底库特征的id
/// @param[in] ids 已经分配好存储空间的数组，数组的长度等于底库中特征的数量
/// @note 底库中特征的数量，可以调用状态接口获取
/// @see se_index_status
SE_API int se_index_get_added_ids(se_index_t index, int64_t *ids);

/// @brief 获取index的平衡系数,仅对带ivf的Index适用
/// @param[out] factor index的平衡系数
SE_API int se_index_imbalance_factor(se_index_t index, float *factor);

/// @brief 清空index中的所有特征
SE_API int se_index_reset(se_index_t index);

/// @brief SearchEngine当前支持的Index类型
typedef enum {
  SE_INDEX_UNKNOW = -1,
  SE_INDEX_DC = 0,
  SE_INDEX_PQ = 1,
  SE_INDEX_FLAT = 2,
  SE_INDEX_HNSW = 3,
  SE_INDEX_EXACT = 4
} se_index_type_e;

typedef struct {
  int32_t index_type;  ///< index类型
  int32_t size;        ///< 已加入的特征的数量
  int32_t is_trained;  ///< index是否被训练
} se_index_status_t;

/// @brief 返回index的状态
SE_API int se_index_status(se_index_t index, se_index_status_t *status);

#define CAMERA_MASK_BYTES 16
/// @brief 时空过滤条件，如果time_range[0]与time_range[1]相同，则表示不进行时空过滤
typedef struct {
  uint8_t camera_id_mask[CAMERA_MASK_BYTES];  ///< 相机点位过滤掩码，最大支持128个相机点位过滤
  uint32_t time_range[2];                     ///< 时间范围
} se_id_filter_t;

typedef enum { DISTANCE_L2 = 0, DISTANCE_INNER_PRODUCT = 1 } se_distance_type_t;

/// @brief 检索接口的配置
typedef struct {
  int64_t k;                     ///< Search最终返回的最相似的特征的个数
  float threshold;               ///< 限定返回的结果中最小的比分
  int32_t batch_size;            ///< 检索时的特征并发数
  se_distance_type_t dist_type;  ///< 检索返回的距离类型
} se_search_config_t;

/// @brief 在检索库中进行特征相度性查询
/// @param[in] features 要查询的特征的数量，多条特征按序排列
/// @param[in] num 查询特征的数量
/// @param[in] filters　每条检索特征对应的时空过滤的条件，如果设置为空，则内部不进行时空过滤
/// @param[in] config 检索的控制参数
/// @param[out] scores 返回的每条检索特征和底库特征最相似的top k条结果的距离分数
/// @param[out] indexes 返回的每条检索特征和底库特征最相似的top k条结果在底库中的id
SE_API int se_index_search(se_index_t index, const float *features, int64_t num,
                           const se_id_filter_t *filters, const se_search_config_t *config,
                           float *scores, int64_t *indexes);

/// @brief 从序列化文件中创建index
SE_API int se_index_create_from_file(se_context_t context, const char *pathname,
                                     se_index_t *index);
/// @brief 将index序列化到文件
SE_API int se_index_serialize_to_file(se_index_t index, const char *pathname);

/// @brief search engine序列化内存数据的结构表示
typedef struct {
  char *data;
  size_t size;
} se_mem_data_t;

/// @brief 从内存数据创建search engine
/// @note 如果创建的是device上的search engine，需要当前线程先绑定好相应的device
SE_API int se_index_create_from_memory(se_context_t context, const char *data, size_t size,
                                       se_index_t *index);

/// @brief 将search engine序列化到内存数据中
/// @note memory_data由接口内部分配，需要用户在外部调用se_index_memory_data_free进行释放
SE_API int se_index_serialize_to_memory(se_index_t index, se_mem_data_t *mem_data);

/// @brief 释放search engine内存数据
SE_API void se_index_memory_data_free(se_mem_data_t *mem_data);

/// @brief 将search engine内存数据写入文件
SE_API int se_index_memory_data_to_file(const char *pathname, const se_mem_data_t *mem_data);

/// @brief 从文件中读取数据到search engine内存数据
/// @note memory_data由接口内部分配，需要用户在外部调用se_index_memory_data_free进行释放
SE_API int se_index_file_to_memory_data(const char *pathname, se_mem_data_t *mem_data);

/// @}

#endif  // SEARCH_ENGINE_INCLUDE_SEARCH_ENGINE_C_H_
