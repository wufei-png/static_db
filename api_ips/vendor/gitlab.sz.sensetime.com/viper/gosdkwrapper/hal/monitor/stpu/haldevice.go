// +build hal_stpu

package stpu

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L/usr/local/stpu/lib -lstpuml -lhaltop -lhalcomm -lpcie_commlib_host
#include <stdlib.h>
#include "stpu/stpu_device.h"
*/
import "C"
import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
)

func (s *stpuMonitor) openHalDevice(id int) error {
	rc := C.stpuHalDeviceOpen(C.uint32_t(id))
	if rc != 0 {
		return common.NewSDKError("stpuHalDeviceOpen", int(rc))
	}
	return nil
}

func (s *stpuMonitor) closeHalDevice(id int) {
	C.stpuHalDeviceClose(C.uint32_t(id))
}
