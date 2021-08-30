// +build hal_cambricon

package cambricon

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L/usr/local/cambricon/lib -lcndev -lcnrt
#include "cndev.h"
#include "cnrt.h"
*/
import "C"

import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
	_ "gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor/cambricon/include"
)

func cnInit() error {
	cndevRet := C.cndevInit(0)
	if C.CNDEV_SUCCESS != cndevRet {
		return common.NewSDKError("cndevInit "+C.GoString(C.cndevGetErrorString(cndevRet)), int(cndevRet))
	}

	cnrtRet := C.cnrtInit(0)
	if C.CNRT_RET_SUCCESS != cnrtRet {
		return common.NewSDKError("cnrtInit "+C.GoString(C.cnrtGetErrorStr(cnrtRet)), int(cnrtRet))
	}

	return nil
}

func cndevCardName() string {
	return C.GoString(C.getCardNameStringByDevId(0))
}

func cnrtGetLibVersion() (int, int) {
	var major, minor, patch C.int
	C.cnrtGetLibVersion(&major, &minor, &patch)
	return int(major), int(minor)
}

func cnrtGetDeviceCount() (int, error) {
	var devNum C.uint
	cnrtRet := C.cnrtGetDeviceCount(&devNum)
	if C.CNRT_RET_SUCCESS != cnrtRet {
		return 0, common.NewSDKError("cnrtGetDeviceCount "+C.GoString(C.cnrtGetErrorStr(cnrtRet)), int(cnrtRet))
	}
	return int(devNum), nil
}

// cndevGetMemoryUsage return mem info in MB
func cndevGetMemoryUsage(deviceId int) (use int, total int, err error) {
	var memInfo C.cndevMemoryInfo_t
	cDevRet := C.cndevGetMemoryUsage(&memInfo, C.int(deviceId))
	if C.CNDEV_SUCCESS != cDevRet {
		return 0, 0, common.NewSDKError("cndevGetMemoryUsage", int(cDevRet))
	}
	return int(memInfo.PhysicalMemoryUsed), int(memInfo.PhysicalMemoryTotal), nil
}

func cnUnInit() error {
	cndevRet := C.cndevRelease()
	if C.CNDEV_SUCCESS != cndevRet {
		return common.NewSDKError("cndevRelease "+C.GoString(C.cndevGetErrorString(cndevRet)), int(cndevRet))
	}

	C.cnrtDestroy()

	return nil
}
