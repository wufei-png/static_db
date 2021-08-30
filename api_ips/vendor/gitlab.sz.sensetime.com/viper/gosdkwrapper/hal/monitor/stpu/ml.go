// +build hal_stpu

package stpu

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L/usr/local/stpu/lib -lstpuml -lhaltop -lhalcomm -lpcie_commlib_host
#include <stdlib.h>
#include "stpu/halml/stpu_ml.h"
*/
import "C"
import (
	"unsafe"

	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
)

const (
	PRODUCT_MODEL_NUMBER_LEN_MAX = 128
)

func (s *stpuMonitor) initHalLibML() {
	rc := C.stpuHalLibMLInit()
	if rc != 0 {
		log.Fatal(common.NewSDKError("stpu_hal_libml_init", int(rc)))
	}
}

func (s *stpuMonitor) getProductModelNumber() (productModelNumber string, err error) {
	// init hal libml first, and only init once
	s.once.Do(s.initHalLibML)

	var len C.uint32_t = PRODUCT_MODEL_NUMBER_LEN_MAX
	cProductModelNumber := (*C.char)(C.malloc(C.ulong(len + 1)))
	defer func() {
		C.free(unsafe.Pointer(cProductModelNumber))
	}()
	rc := C.stpuHalMlAppGetProductModelNumber(cProductModelNumber, len)
	if rc != 0 {
		return "", common.NewSDKError("stpu_halml_app_get_product_model_number", int(rc))
	}
	return C.GoString(cProductModelNumber), nil
}

// getPuMem returan mem info in MB
func (s *stpuMonitor) getPuMem() (use int, total int, err error) {
	// init hal libml first, and only init once
	s.once.Do(s.initHalLibML)

	var cMemUse C.uint32_t
	var cMemTotal C.uint32_t
	rc := C.stpuHalMlPuGetMem(&cMemUse, &cMemTotal)
	if rc != 0 {
		return 0, 0, common.NewSDKError("stpuHalMlPuGetMem ", int(rc))
	}
	return int(cMemUse), int(cMemTotal), nil
}
