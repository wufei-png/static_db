// +build hal_atlas300 hal_atlas500

package acldevice

/*
#cgo CXXFLAGS: -std=c++11
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/usr/local/Ascend/driver/lib64 -ldrvdsmi_host -lc_sec -lmmpa -lascend_hal
#include "dsmi_common_interface.h"
*/
import "C"
import (
	"unsafe"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
)

func getDeviceName(deviceID int) (deviceName string, err error) {
	info := C.struct_dsmi_chip_info_stru{}
	rc := C.dsmi_get_chip_info(C.int(deviceID), &info)
	if rc != 0 {
		return "", common.NewSDKError("dsmi_get_chip_info", int(rc))
	}
	chipType := C.GoString((*C.char)(unsafe.Pointer(&(info.chip_type[0]))))
	chipName := C.GoString((*C.char)(unsafe.Pointer(&(info.chip_name[0]))))

	return chipType + chipName, nil
}

// getMemoryInfo return mem info in MB
func getMemoryInfo(deviceID int) (use int, total int, err error) {
	info := C.struct_dsmi_memory_info_stru{}
	rc := C.dsmi_get_memory_info(C.int(deviceID), &info)
	if rc != 0 {
		return 0, 0, common.NewSDKError("dsmi_get_memory_info", int(rc))
	}
	mem := uint32(info.memory_size)
	uti := uint32(info.utiliza)
	return int(float32(mem) * float32(uti) / 100), int(mem), nil
}
