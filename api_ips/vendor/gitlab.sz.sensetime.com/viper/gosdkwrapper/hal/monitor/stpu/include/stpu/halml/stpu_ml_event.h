#ifndef _STPU_ML_EVENT_H_
#define _STPU_ML_EVENT_H_

#include <stpu/halml/stpu_ml.h>

typedef enum {
    STPU_HAL_ML_EVENT_SOCK_ROLE_UNINITED = 0,
    STPU_HAL_ML_EVENT_SOCK_ROLE_READER,
    STPU_HAL_ML_EVENT_SOCK_ROLE_WRITER
} StpuHalMlEventSockRole;

typedef enum {
    STPU_HAL_ML_EVENT_SOCK_CMD_REGISTER_READER = 1, // writer -> server
    STPU_HAL_ML_EVENT_SOCK_CMD_REGISTER_WRITER, // reader -> server
    STPU_HAL_ML_EVENT_SOCK_CMD_GENERATE_NEW_EVENT, // writer -> server
    STPU_HAL_ML_EVENT_SOCK_CMD_FORWARD_EVENT // server -> reader
} StpuHalMlEventSockCmdType;

typedef struct StpuHalMlEventSet_st {
    uint64_t event_mask;
    char device_uuid[STPU_ML_MAX_PU_COUNT][STPU_ML_DEVICE_TREE_ID_BUFFER_SIZE];
} StpuHalMlEventSet_st;

typedef struct StpuHalMlEventSockInfo {
    int32_t sock;
    int32_t type; // 0 for just connected; 1 for  event reader ; 2 for event writer
    StpuHalMlEventSet_st event_set; // only available for reader
} StpuHalMlEventSockInfo;

typedef struct StpuHalMlEventRawData {
    uint32_t pu_index;
    uint32_t reserved;
    uint64_t event_type;
    uint64_t event_data;
} StpuHalMlEventRawData;

#define STPU_HAL_ML_EVENT_SOCK_CMD_MAGIC 0x11223344
typedef struct StpuHalMlEventSockCmd {
    uint32_t magic;
    uint32_t cmd;
    uint32_t length;
    uint32_t mode;
} StpuHalMlEventSockCmd;

#endif
