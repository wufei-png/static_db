package cudaruntime

/*
#cgo CFLAGS: -I/usr/local/cuda/include
#include <cuda_runtime.h>
*/
import "C"
import (
	"fmt"
)

type Error struct {
	Code int
}

func newError(r C.cudaError_t) error {
	if r == 0 {
		return nil
	}
	return Error{
		Code: int(r),
	}
}

func NewRuntimeErrorFromCode(r int) error {
	return newError(C.cudaError_t(r))
}

func (e Error) Error() string {
	return fmt.Sprintf("CUDA Runtime Error: %d, %s", e.Code, GetErrorString(e.Code))
}
