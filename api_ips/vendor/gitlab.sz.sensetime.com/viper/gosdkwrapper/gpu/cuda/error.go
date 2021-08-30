package cuda

/*
#include "dynlink_cuda.h"
*/
import "C"
import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/video/decode"
)

type Error = decode.Error

func newCudaError(r C.CUresult) error {
	return decode.NewDecodeErrorFromCode(int(r), "cuda")
}

var (
	ErrUnknownFormat = decode.ErrUnknownFormat
	ErrInvalidPacket = decode.ErrInvalidPacket
)
