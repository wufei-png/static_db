package cudaruntime

/*
#cgo CFLAGS: -I/usr/local/cuda/include
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart -lcuda
#include <cuda_runtime.h>
#include <cuda.h>
*/
import "C"

import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/utils/device"
)

const (
	// TODO: DEPRECATED due to redefine in hal/monitor/define.go
	P4 = "Tesla P4"
	// 部分驱动下(如367.48), P4的name为"Graphics Device"
	GraphicsDevice = device.NvidiaGraphicsDevice
	T4             = device.NvidiaT4
	RTX4000        = device.NvidiaRTX4000
	GTX1060        = device.NvidiaGTX1060
	GTX1660        = device.NvidiaGTX1660
	GTX1070        = device.NvidiaGTX1070
	GTX2080        = device.NvidiaGTX2080
)

func DeviceGetCount() (int, error) {
	var count C.int
	ce := C.cuDeviceGetCount(&count)
	if ce != 0 {
		return 0, common.NewSDKError("cuda", int(ce))
	}
	return int(count), nil
}

func SetDevice(id int) error {
	return newError(C.cudaSetDevice(C.int(id)))
}

func RuntimeContext() error {
	return newError(C.cudaFree(nil))
}

func ResetDevice() error {
	return newError(C.cudaDeviceReset())
}

func DeviceSynchronize() error {
	return newError(C.cudaDeviceSynchronize())
}

func GetErrorString(code int) string {
	return C.GoString(C.cudaGetErrorString(C.cudaError_t(code)))
}

func GetLastError() error {
	return newError(C.cudaGetLastError())
}

func GetMemInfo() (int, int, error) {
	var freeByte, totalByte C.size_t
	r := C.cudaMemGetInfo(&freeByte, &totalByte)
	if r == 0 {
		return int(freeByte), int(totalByte), nil
	}
	return 0, 0, newError(r)
}

type CUDADeviceProperties struct {
	Major int
	Minor int
	Name  string
}

func GetCUDADeviceProperties(deviceID int) (*CUDADeviceProperties, error) {
	var deviceProperties CUDADeviceProperties
	var cp C.struct_cudaDeviceProp

	// https://docs.nvidia.com/cuda/cuda-runtime-api/group__CUDART__DEVICE.html
	err := newError(C.cudaGetDeviceProperties(&cp, C.int(deviceID)))
	deviceProperties.Major = int(cp.major)
	deviceProperties.Minor = int(cp.minor)
	deviceProperties.Name = C.GoString(&(cp.name[0]))
	return &deviceProperties, err
}

func GetCUDAVersion() (int, error) {
	var version C.int
	err := newError(C.cudaRuntimeGetVersion(&version))
	return int(version), err
}
