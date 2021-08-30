// +build hal_stpu

package stpu

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L/usr/local/stpu/lib -lstpuml -lhaltop -lhalcomm -lpcie_commlib_host
#include <stdlib.h>
#include "stpu/haltop/haltop.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"strconv"
	"unsafe"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
)

func (s *stpuMonitor) initHalTop(id int) error {

	var option string

	switch runtime.GOARCH {
	case "arm64":
		option = "Unix:/var/run/stpuhal"
	case "amd64":
		option = "PCIe:" + strconv.Itoa(id)
	default:
		return errors.New("device type invalid")
	}

	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))

	rc := C.stpuHalTopInit(cOption)
	if rc != 0 {
		return common.NewSDKError("stpuHalTopInit", int(rc))
	}
	return nil
}

func (s *stpuMonitor) deinitHalTop() {
	C.stpuHalTopDeinit()
}
