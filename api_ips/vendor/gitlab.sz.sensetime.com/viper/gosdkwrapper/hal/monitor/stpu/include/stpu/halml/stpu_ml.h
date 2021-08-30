/**
 * @file stpu_ml.h
 */

#ifndef __HAL_ML_H__
#define __HAL_ML_H__

#include <unistd.h>
#include <stdint.h>
#include "version.h"

#ifdef __cplusplus
#define STPU_HAL_API extern "C" __attribute__((visibility("default")))
#else
#define STPU_HAL_API __attribute__((visibility("default")))
#endif

typedef int32_t HalError;

#define MAX_DRV_UUID_LEN 64
#define STPU_ML_DEVICE_TREE_ID_BUFFER_SIZE 128
#define STPU_ML_MAX_CHIP_PER_BOARD 4
#define STPU_ML_MAX_PU_PER_CHIP 4
#define STPU_ML_MAX_SUBSYS_PER_CHIP 2
#define STPU_ML_MAX_PU_PER_SUBSYS 2
#define STPU_ML_MAX_PU_PER_BOARD (STPU_ML_MAX_CHIP_PER_BOARD * STPU_ML_MAX_PU_PER_CHIP)
#define STPU_ML_MAX_BOARD 16
#define STPU_ML_MAX_PU_COUNT                                                                     \
    (STPU_ML_MAX_BOARD * STPU_ML_MAX_CHIP_PER_BOARD * STPU_ML_MAX_SUBSYS_PER_CHIP *              \
     STPU_ML_MAX_PU_PER_SUBSYS)

typedef enum {
    STPU_ML_SUCCESS = 0, // The operation was successful
    STPU_ML_ERROR_UNINITIALIZED = 1, // ML was not first initialized with stpuMlInit()
    STPU_ML_ERROR_INVALID_ARGUMENT = 2, // A supplied argument is invalid
    STPU_ML_ERROR_NOT_SUPPORTED = 3, // The requested operation is not available on target device
    STPU_ML_ERROR_NO_PERMISSION = 4, // The current user does not have permission for operation
    STPU_ML_ERROR_ALREADY_INITIALIZED =
        5, // Deprecated: Multiple initializations are now allowed through ref counting
    STPU_ML_ERROR_NOT_FOUND = 6, // A query to find an object was unsuccessful
    STPU_ML_ERROR_INSUFFICIENT_SIZE = 7, // An input argument is not large enough
    STPU_ML_ERROR_INSUFFICIENT_POWER =
        8, // A device's external power cables are not properly attached
    STPU_ML_ERROR_DRIVER_NOT_LOADED = 9, // STPU driver is not loaded
    STPU_ML_ERROR_TIMEOUT = 10, // User provided timeout passed
    STPU_ML_ERROR_LIBRARY_NOT_FOUND = 11, // ML Shared Library couldn't be found or loaded
    STPU_ML_ERROR_FUNCTION_NOT_FOUND = 12, // Local version of ML doesn't implement this function
    STPU_ML_ERROR_STPU_IS_LOST =
        13, // The STPU has fallen off the bus or has otherwise become inaccessible
    STPU_ML_ERROR_RESET_REQUIRED = 14, // The STPU requires a reset before it can be used again
    STPU_ML_ERROR_UNKNOWN = 999 // An internal driver error occurred
} StpuHalMlReturn;

typedef enum StpuHalMlPuResources {
    STPU_PU_RES_RESIZE = 0, // resize
    STPU_PU_RES_VPU, // video decoder
    STPU_PU_RES_JPEG, // jpeg decoder & encode
    STPU_PU_RES_HIST,
    STPU_PU_RES_CROP,
    STPU_PU_RES_TPU,
    STPU_PU_RES_ARM,
    STPU_PU_RES_DDR,
    STPU_PU_RES_MAX
} StpuHalMlPuResources;

typedef enum StpuHalMlDeviceType {
    STPU_DEVICE_TYPE_BOARD = 0,
    STPU_DEVICE_TYPE_CHIP,
    STPU_DEVICE_TYPE_PU,
    STPU_DEVICE_TYPE_MAX
} StpuHalMlDeviceType;

typedef struct StpuHalMlDevice {
    StpuHalMlDeviceType device_type; // device type
} StpuHalMlDevice;

typedef struct StpuHalMlPuArmInfo {
    uint32_t arm_mem_size; // arm top memsize limit */
    uint32_t arm_core_freq; // arm core freq max
    float arm_core_limit; // arm core limit
} StpuHalMlPuArmInfo;

typedef struct StpuHalMlPuVpuInfo {
    uint32_t enc_ability; // fps, 1080p
    uint32_t dec_ability; // fps, 1080p
    uint32_t freq_limit; // vpu freq limit
} StpuHalMlPuVpuInfo;

typedef struct StpuHalMlPuJpegInfo {
    uint32_t enc_ability; // fps, 1080p
    uint32_t dec_ability; // fps, 1080p
    uint32_t freq_limit; // jpeg freq limit
} StpuHalMlPuJpegInfo;

typedef struct StpuHalMlPuHistInfo {
    uint32_t hist_ability; // fps, 1080p
    uint32_t freq_limit; // hist freq limit
} StpuHalMlPuHistInfo;

typedef struct StpuHalMlPuCropInfo {
    uint32_t crop_ability; // fps, 1080p
    uint32_t freq_limit; // crop freq limit
} StpuHalMlPuCropInfo;

typedef struct StpuHalMlPuTpuInfo {
    uint32_t tpu_mem_size; // tpu memsize limit
    uint32_t tpu_freq; // tpu freq limit
} StpuHalMlPlTpuInfo;

typedef struct StpuHalMlPuInfo_st {
    StpuHalMlPuArmInfo arm;
    StpuHalMlPuVpuInfo vpu;
    StpuHalMlPuJpegInfo jpeg;
    StpuHalMlPuHistInfo hist;
    StpuHalMlPuCropInfo crop;
    StpuHalMlPlTpuInfo tpu;
    uint32_t top_mem_size; // pu location soc memory
    uint32_t local_mem_szie; // pu local memory size
    //
    uint32_t mem_size; // MB
    uint32_t pu_index;
    uint32_t unitResources[STPU_PU_RES_MAX];
    // uint8_t sub_sys_id;
    uint32_t minor_number;
    char uuid[STPU_ML_DEVICE_TREE_ID_BUFFER_SIZE]; // pu uuid
} StpuHalMlPuInfo;

typedef struct StpuHalMlSubsysInfo_st {
    char id[STPU_ML_DEVICE_TREE_ID_BUFFER_SIZE];
    uint32_t pu_count;
    uint32_t pu_index[STPU_ML_MAX_PU_PER_SUBSYS];
    uint32_t subsys_index;
    uint8_t sub_sys_id;
    uint8_t res[3];
} StpuHalMlSubsysInfo;

typedef struct StpuHalMlChipInfo_st {
    char cuuid[MAX_DRV_UUID_LEN + 4];
    uint32_t subsys_count; // chip's subsys count
    uint32_t subsys_index[STPU_ML_MAX_SUBSYS_PER_CHIP]; // chip's subsys index
    uint32_t pu_count; // chip's pu count
    uint32_t pu_index[STPU_ML_MAX_PU_PER_CHIP]; // chip's pu index
    uint32_t chip_index;
    uint32_t top_mem; // soc total memory
} StpuHalMlChipInfo;

typedef struct StpuHalMlBoardInfo_st {
    char buuid[MAX_DRV_UUID_LEN + 4]; // board uuid
    uint32_t chip_count;
    uint32_t chip_index[STPU_ML_MAX_CHIP_PER_BOARD];
    uint32_t board_index;
    char *board_serial;
} StpuHalMlBoardInfo;

typedef void *StpuHalMlPuHandler;
typedef void *StpuHalMlSubsysHandler;
typedef void *StpuHalMlChipHandler;
typedef void *StpuHalMlBoardHandler;

#define STPU_ML_DEVICE_PCI_BUS_ID_LEGACY_BUFFER_SIZE 20
#define STPU_ML_DEVICE_PCI_BUS_ID_BUFFER_SIZE 32

typedef struct StpuHalMlPciInfo_st {
    char bus_id_legacy[STPU_ML_DEVICE_PCI_BUS_ID_LEGACY_BUFFER_SIZE];
    uint32_t domain;
    uint32_t bus;
    uint32_t device;
    uint32_t pci_device_id;
    uint32_t pci_subsys_id;
    char bus_id[STPU_ML_DEVICE_PCI_BUS_ID_BUFFER_SIZE];
} StpuHalMlPciInfo;

typedef struct StpuHalMlDevicePu_st {
    uint64_t events_subscribed;
    StpuHalMlPuInfo *pu;
    char uuid[STPU_ML_DEVICE_TREE_ID_BUFFER_SIZE];
    // uint32_t minor_number;
    StpuHalMlPciInfo pci_info;
} StpuHalMlDevicePu;
typedef void *StpuHalMlDevicePuHandle;

typedef enum StpuHalMlParamType {
    STPU_HAL_ML_PARAM_TEMPERATEURE,
    STPU_HAL_ML_PARAM_POWER_STATUS,
    STPU_HAL_ML_PARAM_UTILIZATION,
    STPU_HAL_ML_PARAM_TYPE_MAX
} StpuHalMlParamType;

/* chip temperateure */
#define CHIP_SENSOR_COUNT 10
typedef struct StpuHalMlChipTemperateure {
    float current_max; // current max temperateure
    uint32_t max_index; // the index of senseor
    float sensor[CHIP_SENSOR_COUNT]; // temperateure value of every sensor
} StpuHalMlChipTemperateure;

#define BOARD_SENSOR_COUNT_MAX 4
typedef struct StpuHalMlBoardTemperature {
    float current_max; // current max temperateure
    float board_sensor[BOARD_SENSOR_COUNT_MAX]; // temperateure value of every sensor
    StpuHalMlChipTemperateure chip_temp[STPU_ML_MAX_CHIP_PER_BOARD];
} StpuHalMlBoardTemperature;

typedef enum StpuHalMlLowPowerStatus {
    STPU_HAL_ML_ALL_ON,
    STPU_HAL_ML_CLOCK_OFF,
    STPU_HAL_ML_POWER_OFF,
    STPU_HAL_ML_AUTO_SET,
    STPU_HAL_ML_POWER_STATUS_MAX,

    STPU_HAL_ML_ALWAYS_ON = STPU_HAL_ML_ALL_ON,
    STPU_HAL_ML_SLEEP = STPU_HAL_ML_CLOCK_OFF,
    STPU_HAL_ML_DEEP_SLEEP = STPU_HAL_ML_POWER_OFF
} StpuHalMlLowPowerStatus;

typedef struct StpuHalMlPuPower {
    StpuHalMlLowPowerStatus lp_status[STPU_PU_RES_MAX]; // low power status
} StpuHalMlPuPower;

typedef struct StpuHalMlChipPower {
    float chip_power; // W
    StpuHalMlLowPowerStatus lp_status; // all chip low power status
    StpuHalMlPuPower pu_power[STPU_ML_MAX_PU_PER_CHIP]; // pu power status
} StpuHalMlChipPower;

typedef struct StpuHalMlBoardPower {
    float board_power; // total power
    float board_current; // current
} StpuHalMlBoardPower;

typedef struct StpuHalMlPuUtilization {
    // StpuHalMlPuResources pu_res; // pu resources
    float util[STPU_PU_RES_MAX]; // Utilization
} StpuHalMlPuUtilization;

// status

typedef struct StpuHalMlBoardStatus {
    StpuHalMlBoardTemperature board_tem; // board temperature
    StpuHalMlBoardPower board_power; // board power
} StpuHalMlBoardStatus;

typedef struct StpuHalMlChipStatus {
    StpuHalMlChipTemperateure chip_tem; // chip temperature
    StpuHalMlChipPower chip_power; // chip power
    float chip_util; // chip cpu util 0-100%
    float running_time; // chip boot time, second
    uint32_t top_mem; // MB
} StpuHalMlChipStatus;

typedef struct StpuHalMlSubsysStatus {

} StpuHalMlSubsysStatus;

typedef struct StpuHalMlPuStatus {
    StpuHalMlPuPower pu_power; // pu power
    StpuHalMlPuUtilization pu_util; // pu util
    uint32_t local_mem; // pu local MB
    uint32_t top_mem; // pu location's chip memory MB
} StpuHalMlPuStatus;

typedef struct StpuHalMlPuHealth {
    uint32_t online; // pu service is online
    double running_time; // pu service running time
    int32_t restart_count; // pu service restart count
} StpuHalMlPuHealth;

// set param
typedef struct StpuHalMlPowerStatus {
    StpuHalMlLowPowerStatus lp_status;
} StpuHalMlPowerStatus;

typedef struct StpuHalMlPowerLimit {
    float power_limit;
} StpuHalMlPowerLimit;

/* no events */
#define STPU_ML_EVENT_TYPE_NONE 0x00000000LL

/* power state change, used for pu device */
#define STPU_ML_EVENT_TYPE_PSTATE 0x00000001LL

/* XID event, critical error occurred */
#define STPU_ML_EVENT_TYPE_XID_ERROR 0x00000002LL

/* clock change */
#define STPU_ML_EVENT_TYPE_CLOCK 0x00000004LL

#define STPU_ML_EVENT_PU_ONLINE 0x00000008LL

#define STPU_ML_EVENT_PU_OFFLINE 0x00000010LL

/* all event mask */
#define STPU_ML_TYPE_ALL                                                                         \
    (STPU_ML_EVENT_TYPE_PSTATE | STPU_ML_EVENT_TYPE_XID_ERROR | STPU_ML_EVENT_TYPE_CLOCK |       \
     STPU_ML_EVENT_PU_ONLINE | STPU_ML_EVENT_PU_OFFLINE)

typedef void *StpuHalMlEventSet;

/**/
typedef struct StpuHalMlEventData {
    StpuHalMlDevicePuHandle device; // Specific device where the event occurred
    uint64_t event_data; // Stores last XID error for the device,
    // event_data is 0 for any other event.
    uint64_t event_type; // Information about what specific event occurred
    // event_type is set as 999 for unknown xid error.
} StpuHalMlEventData;

/**
 * @brief  ML init
 * @return STPU_HAL_API stpuHalMlInit
 */
STPU_HAL_API HalError stpuHalMlInit(void);

/**
 * @brief  ML exit, must do stpuHalMlInit first
 * @return STPU_HAL_API stpuHalMlExit
 */
STPU_HAL_API HalError stpuHalMlExit(void);

/**
 * @brief  init ml lib, thread local
 * @return STPU_HAL_API stpuHalLibMLInit
 */
STPU_HAL_API HalError stpuHalLibMLInit(void);

/**
 * @brief  get pci driver version
 * @param  version          [out]version string addr
 * @param  length           [out]version string length
 * @return STPU_HAL_API stpuHalMlSystemGetDriverVersion
 */
STPU_HAL_API HalError stpuHalMlSystemGetDriverVersion(char *version, uint32_t length);

/**
 * @brief  get ml lib version
 * @param  version          [out]version string addr
 * @param  length           [out]version string length
 * @return STPU_HAL_API stpuHalMlSystemGetStmlVersion
 */
STPU_HAL_API HalError stpuHalMlSystemGetStmlVersion(char *version, uint32_t length);

/**
 * @brief err code to string
 *
 * @param errorCode         [in]input errcode
 * @return STPU_HAL_API const*
 */
STPU_HAL_API const char *stpuHalMlErrorString(int32_t errorCode);

//********************************* device inf (for k8s)start
//*******************************************//
/**
 * @brief get device count, pu
 *
 * @param device_count      [out]return device count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetCount(uint32_t *device_count);
/**
 * @brief get device handle by index, index(0 - device_count-1)
 *
 * @param index             [in]input index
 * @param device_pu         [out]return handle
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetHandleByIndex(uint32_t index,
                                                      StpuHalMlDevicePuHandle *device_pu);

/**
 * @brief get device uuid
 *
 * @param device_pu         [in]device pu handle
 * @param uuid              [out]pu uuid
 * @param length            [out]length of uuid
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetUuid(StpuHalMlDevicePuHandle device_pu, char *uuid,
                                             uint32_t length);
/**
 * @brief get device is online
 *
 * @param index             [in]device index
 * @param pu_online         [out]pu is online
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetHealthByIndex(uint32_t index, uint32_t *pu_online);

/**
 * @brief get device minor number
 *
 * @param device_pu         [in]device pu handle
 * @param minor_num         [out]minor number
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetMinorNumber(StpuHalMlDevicePuHandle device_pu,
                                                    uint32_t *minor_num);
/**
 * @brief get device pci info
 *
 * @param device_pu         [in]device pu handle
 * @param pci_info          [out]device pci info
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetPciInfo(StpuHalMlDevicePuHandle device_pu,
                                                StpuHalMlPciInfo *pci_info);
/**
 * @brief  get pu name from pu index
 * @param  pu_index         [in]input pu index
 * @param  pu_name          [out]out pu name, like "/dev/stpu/pu/pcie_ipc_rc.0.pu1"
 * @param  length           [out]input pu_name string length
 * @return STPU_HAL_API stpuHalMlGetDevNameByPuIndex
 */
STPU_HAL_API HalError stpuHalMlGetDevNameByPuIndex(int32_t pu_index, char *pu_name,
                                                   uint32_t length);
//*********************************** device inf end *******************************************//

//*************************************** host inf start ***************************//
/**
 * @brief get board count
 *
 * @param board_count       [out]all board count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetCount(uint32_t *board_count);

/**
 * @brief get board handle by index
 *
 * @param index             [in]board index
 * @param board             [out]board handle
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetHandleByIndex(uint32_t index,
                                                     StpuHalMlBoardHandler *board);
/**
 * @brief get board uuid, equal to board product s/n
 *
 * @param board             [in]board handle
 * @param buuid             [out]board uuid
 * @param length            [out]length of buuid
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetBuuid(StpuHalMlBoardHandler board, char *buuid,
                                             uint32_t length);
/**
 * @brief get one board's chip count
 *
 * @param board             [in]board handle
 * @param chip_count        [out]board's chip count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetChipCount(StpuHalMlBoardHandler board,
                                                 uint32_t *chip_count);
/**
 * @brief get board status, include power and Temperature
 *
 * @param board             [in]board handle
 * @param board_status      [out]board status
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetStatus(StpuHalMlBoardHandler board,
                                              StpuHalMlBoardStatus *board_status);
/**
 * @brief get board info
 *
 * @param board             [in]board handle
 * @param board_info        [out]board info
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetInfo(StpuHalMlBoardHandler board,
                                            StpuHalMlBoardInfo *board_info);
/**
 * @brief get one board's pu count and pu index
 *
 * @param board             [in]board handle
 * @param pu_index          [out]board's pu index
 * @param pu_count          [out]board's pu count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetPuIndex(StpuHalMlBoardHandler board, uint32_t *pu_index,
                                               uint32_t *pu_count);
/**
 * @brief get one board's chip count and chip index
 *
 * @param board             [in]board handle
 * @param chip_index        [out]board's chip index
 * @param chip_count        [out]board's chip count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlBoardGetChipIndex(StpuHalMlBoardHandler board,
                                                 uint32_t *chip_index, uint32_t *chip_count);
// chip
/**
 * @brief get chip handle
 *
 * @param board             [in]board handle
 * @param index             [in]chip index
 * @param chip              [out]chip handle
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetHandleByIndex(StpuHalMlBoardHandler board, uint32_t index,
                                                    StpuHalMlChipHandler *chip);
/**
 * @brief get chip cuuid
 *
 * @param chip              [in]chip handle
 * @param cuuid             [out]chip uuid
 * @param length            [out]length of cuuid
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetCuuid(StpuHalMlChipHandler chip, char *cuuid,
                                            uint32_t length);
/**
 * @brief get chip's subsys count
 *
 * @param chip              [in]chip handle
 * @param subsys_count      [out]chip's subsys count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetSubsysCount(StpuHalMlChipHandler chip,
                                                  uint32_t *subsys_count);
/**
 * @brief get chip's subsys index
 *
 * @param chip              [in]chip handle
 * @param subsys_count      [out]chip's subsys count
 * @param subsys_index      [out]chip's subsys index
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetSubsysIndex(StpuHalMlChipHandler chip,
                                                  uint32_t *subsys_count, uint32_t *subsys_index);
/**
 * @brief get chip status
 *
 * @param chip              [in]chip handle
 * @param chip_status       [out]chip status
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetStatus(StpuHalMlChipHandler chip,
                                             StpuHalMlChipStatus *chip_status);
/**
 * @brief get chip info
 *
 * @param chip              [in]chip handle
 * @param chip_info         [out]chip info
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetInfo(StpuHalMlChipHandler chip,
                                           StpuHalMlChipInfo *chip_info);
/**
 * @brief get chip's pu index
 *
 * @param chip              [in]chip handle
 * @param pu_index          [out]chip's pu index
 * @param pu_count          [out]chip's pu count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetPuIndex(StpuHalMlChipHandler chip, uint32_t *pu_index,
                                              uint32_t *pu_count);

/**
 * @brief get chip manage service is online
 *
 * @param chip_index        [in]chip index
 * @param online            [out]chip is online
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetManageProcStatus(uint32_t chip_index, uint32_t *online);

/**
 * @brief get chip ip
 *
 * @param chip              [in]chip handle
 * @param ip                [out]chip ip, like "192.168.1.10"
 * @param len               [out]length of ip
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetIp(StpuHalMlChipHandler chip, char *ip, uint32_t len);

/**
 * @brief get pu count and pu index
 *
 * @param board_id          [in]board id
 * @param chip_id           [in]chip id
 * @param pu_index          [out]pu index
 * @param pu_count          [out]pu count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetPuIndex(uint32_t board_id, uint32_t chip_id, uint32_t *pu_index,
                                          uint32_t *pu_count);
// subsys
/**
 * @brief get subsys handle
 *
 * @param chip              [in]chip handle
 * @param index             [in]subsys index
 * @param subsys            [out]subsys handle
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlSubsysGetHandleByIndex(StpuHalMlChipHandler chip, uint32_t index,
                                                      StpuHalMlSubsysHandler *subsys);
/**
 * @brief get subsys's pu count
 *
 * @param subsys            [in]subsys handle
 * @param pu_count          [out]subsys's pu count
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlSubsysGetPuCount(StpuHalMlSubsysHandler subsys,
                                                uint32_t *pu_count);
/**
 * @brief get subsys's pu count and pu index
 *
 * @param subsys            [in]subsys handle
 * @param pu_count          [out]subsys's pu count
 * @param pu_index          [out]subsys's pu index
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlSubsysGetPuIndex(StpuHalMlSubsysHandler subsys, uint32_t *pu_count,
                                                uint32_t *pu_index);
/**
 * @brief get subsys status
 *
 * @param subsys            [in]subsys handle
 * @param subsys_status     [out]subsys status
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlSubsysGetStatus(StpuHalMlSubsysHandler subsys,
                                               StpuHalMlSubsysStatus *subsys_status);
/**
 * @brief get subsys info
 *
 * @param subsys            [in]subsys handle
 * @param subsys_info       [out]subsys info
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlSubsysGetInfo(StpuHalMlSubsysHandler subsys,
                                             StpuHalMlSubsysInfo *subsys_info);

// pu
/**
 * @brief get pu handle
 *
 * @param subsys            [in]subsys handle
 * @param index             [in]pu index
 * @param pu                [out]pu handle
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetHandleByIndex(StpuHalMlSubsysHandler subsys, uint32_t index,
                                                  StpuHalMlPuHandler *pu);
/**
 * @brief get pu uuid
 *
 * @param pu                [in]pu handle
 * @param uuid              [out]pu uuid
 * @param length            [out]length of uuid
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetUuid(StpuHalMlPuHandler pu, char *uuid, uint32_t length);

/**
 * @brief get pu status
 *
 * @param pu                [in]pu handle
 * @param pu_status         [out]pu status
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetStatus(StpuHalMlPuHandler pu, StpuHalMlPuStatus *pu_status);

/**
 * @brief get pu info
 *
 * @param pu                [in]pu handle
 * @param pu_info           [out]pu info
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetInfo(StpuHalMlPuHandler pu, StpuHalMlPuInfo *pu_info);

/**
 * @brief get pu is online
 *
 * @param pu_index          [in]pu index
 * @param pu_online         [out]pu is online
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetHealth(uint32_t pu_index, uint32_t *pu_online);

/**
 * @brief get pu health sync, make sure bind pu service is alive
 *
 * @param pu_index          [in]pu index
 * @param pu_health         [out]pu's health status
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetHealthSync(uint32_t pu_index, StpuHalMlPuHealth *pu_health);

/**
 * @brief get all pu is online
 *
 * @param pu_count          [out]all pu count
 * @param pu_health         [out]all pu health
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetHealth(uint32_t *pu_count, StpuHalMlPuHealth *pu_health);

// event
/**
 * @brief get support event type
 *
 * @param device            [in]device handle
 * @param event_type        [out]event type
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceGetSupportedEventTypes(StpuHalMlDevicePuHandle device,
                                                            uint64_t *event_type);
/**
 * @brief create a eventset
 *
 * @param event_set         [out]event set
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlEventSetCreate(StpuHalMlEventSet *event_set);

/**
 * @brief destort a eventset
 *
 * @param event_set         [in]event set
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlEventSetDestory(StpuHalMlEventSet event_set);

/**
 * @brief register event
 *
 * @param device            [in]device handle
 * @param event_type        [in]event type
 * @param event_set         [in]event set
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlDeviceRegisterEvent(StpuHalMlDevicePuHandle device,
                                                   uint64_t event_type,
                                                   StpuHalMlEventSet event_set);

/**
 * @brief wait timeout_ms for registered event
 *
 * @param event_set         [in]event set
 * @param data              [out]event data
 * @param timeout_ms        [in]wait time
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlEventSetWait(StpuHalMlEventSet event_set, StpuHalMlEventData *data,
                                            int32_t timeout_ms);

// bmc
/**
 * @brief upgrade bmc fw
 *
 * @param board             [in]board handle
 * @param hex_path          [in]fw file
 * @param dl_path           [in]exec shell
 * @param stm_path          [in]stm32flash
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlUpgradeFw(StpuHalMlBoardHandler board, char *hex_path,
                                         char *dl_path, char *stm_path);

/**
 * @brief get bmc fw version
 *
 * @param board             [in]board handle
 * @param fw_version        [out]board bmc fw version
 * @param len               [out]lenth of fw_version
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetFwVersion(StpuHalMlBoardHandler board, char *fw_version,
                                            uint32_t len);

/**
 * @brief get board product version
 *
 * @param board             [in]board handle
 * @param product_version   [out]board product version
 * @param len               [out]length of product_version
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetProductVersion(StpuHalMlBoardHandler board,
                                                 char *product_version, uint32_t len);

/**
 * @brief get board product name
 *
 * @param board             [in]board handle
 * @param product_name      [out]product name
 * @param len               [out]length of product_name
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetProductName(StpuHalMlBoardHandler board, char *product_name,
                                              uint32_t len);

/**
 * @brief get board product model number
 *
 * @param board             [in]board handle
 * @param product_model_number  [out]board product model number
 * @param len               [out]length of product_model_number
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlGetProductModelNumber(StpuHalMlBoardHandler board,
                                                     char *product_model_number, uint32_t len);

//*************************************** host inf end ***************************//

//*************************************** docker/app inf start ***************************//
/**
 * @brief get chip top mem
 *
 * @param chip_mem_use      [out]chip/soc top mem in use(MB)
 * @param chip_mem_total    [out]chip/soc top total mem(MB)
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetMem(uint32_t *chip_mem_use, uint32_t *chip_mem_total);

/**
 * @brief chip/soc cpu util
 *
 * @param chip_cpu_util     [out]range: 0.0-100.0
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlChipGetCpuUtil(float *chip_cpu_util);

/**
 * @brief get bind pu local mem
 *
 * @param pu_mem_use        [out]pu local mem in use(MB)
 * @param pu_mem_total      [out]pu total local mem(MB)
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetMem(uint32_t *pu_mem_use, uint32_t *pu_mem_total);

/**
 * @brief get bind pu use util
 *
 * @param pu_util           [out]range: 0.0-100.0
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetUtil(float *pu_util);

/**
 * @brief get bind pu service process's cpu util
 *
 * @param pu_util           [out]range: 0.0-800.0(one chip has 8 cpu)
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlPuGetArmCpuUtil(float *pu_util);

/**
 * @brief get bind pu's board product version, like "A101"
 *
 * @param product_version   [out]product version
 * @param len               [out]length of product_version
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlAppGetProductVersion(char *product_version, uint32_t len);

/**
 * @brief get bind pu's board product name, like "SenseVenus PCIe"
 *
 * @param product_name      [out]product name
 * @param len               [out]length of product_name
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlAppGetProductName(char *product_name, uint32_t len);

/**
 * @brief get bind pu's board product part/model number, like "SV-116E"
 *
 * @param product_model_number  [out]product model number
 * @param len               [out]length of product_model_number
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlAppGetProductModelNumber(char *product_model_number, uint32_t len);

/**
 * @brief get bind pu's board product serial number, like "LKD0335848"
 *
 * @param product_serial_number [out]product serial number
 * @param len               [out]length of product_serial_number
 * @return STPU_HAL_API
 */
STPU_HAL_API HalError stpuHalMlAppGetProductSerialNumber(char *product_serial_number,
                                                         uint32_t len);

//*************************************** docker/app inf end ***************************//

#endif