package cuda

/*
#cgo CFLAGS: -I/usr/local/cuda/include
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcuda
#include <string.h>
#include <cuda.h>
#include <cudaProfiler.h>

CUresult _wrapCuMemcpy2D(CUDA_MEMCPY2D m2d) {
	return cuMemcpy2D(&m2d);
}

CUresult _wrapCuMemcpy2DAsync(CUDA_MEMCPY2D m2d, CUstream s) {
	return cuMemcpy2DAsync(&m2d, s);
}
*/
import "C"
import (
	"unsafe"
)

const CudaAPIVersion = 4000

// nolint: golint
const (
	CU_CTX_SCHED_AUTO          = C.CU_CTX_SCHED_AUTO
	CU_CTX_SCHED_SPIN          = C.CU_CTX_SCHED_SPIN
	CU_CTX_SCHED_YIELD         = C.CU_CTX_SCHED_YIELD
	CU_CTX_SCHED_BLOCKING_SYNC = C.CU_CTX_SCHED_BLOCKING_SYNC

	CU_CTX_BLOCKING_SYNC = C.CU_CTX_BLOCKING_SYNC

	CU_CTX_MAP_HOST           = C.CU_CTX_MAP_HOST
	CU_CTX_LMEM_RESIZE_TO_MAX = C.CU_CTX_LMEM_RESIZE_TO_MAX
)

func InitCuda(flags uint, version int) error {
	// ce := C.cuInit(C.uint(flags), C.int(version), nil)
	ce := C.cuInit(C.uint(flags))
	return newCudaError(ce)
}

func DeviceGetCount() (int, error) {
	var count C.int
	ce := C.cuDeviceGetCount(&count)
	if ce != 0 {
		return 0, newCudaError(ce)
	}
	return int(count), nil
}

type Device struct {
	h C.CUdevice
}

type DevicePointer C.CUdeviceptr

type Context struct {
	h C.CUcontext
}

func DeviceGet(ord int) (*Device, error) {
	var dev C.CUdevice
	ce := C.cuDeviceGet(&dev, C.int(ord))
	if ce != 0 {
		return nil, newCudaError(ce)
	}
	return &Device{h: dev}, nil
}

func (d *Device) CtxCreate(flags uint) (Context, error) {
	var ctx C.CUcontext
	ce := C.cuCtxCreate(&ctx, C.uint(flags), d.h)
	if ce != 0 {
		return Context{}, newCudaError(ce)
	}
	return Context{h: ctx}, nil
}

func CtxGetCurrent() (Context, error) {
	var ctx C.CUcontext
	ce := C.cuCtxGetCurrent(&ctx)
	if ce != 0 {
		return Context{}, newCudaError(ce)
	}
	return Context{h: ctx}, nil
}

func (c Context) Destroy() error {
	ce := C.cuCtxDestroy(c.h)
	return newCudaError(ce)
}

func CtxSynchronize() error {
	ce := C.cuCtxSynchronize()
	return newCudaError(ce)
}

func CtxPopCurrent() (Context, error) {
	var ctx C.CUcontext
	ce := C.cuCtxPopCurrent(&ctx)
	if ce != 0 {
		return Context{}, newCudaError(ce)
	}
	return Context{h: ctx}, nil
}

func CtxPushCurrent(ctx Context) error {
	ce := C.cuCtxPushCurrent(ctx.h)
	return newCudaError(ce)
}

type Stream struct {
	s C.CUstream
}

var DefaultStream = Stream{s: nil}

func CreateStream(flags uint) (Stream, error) {
	var s C.CUstream
	ce := C.cuStreamCreate(&s, C.uint(flags))
	if ce != 0 {
		return Stream{}, newCudaError(ce)
	}
	return Stream{s: s}, nil
}

func (s Stream) Synchronize() error {
	ce := C.cuStreamSynchronize(s.s)
	return newCudaError(ce)
}

func (s Stream) Destroy() error {
	ce := C.cuStreamDestroy(s.s)
	return newCudaError(ce)
}

type DeviceMemory struct {
	owned bool
	ptr   C.CUdeviceptr
}

func DeviceMemoryFromPointer(dptr DevicePointer) DeviceMemory {
	return DeviceMemory{
		ptr: C.CUdeviceptr(dptr),
	}
}

func DeviceMemoryFromUnsafePointer(dptr unsafe.Pointer) DeviceMemory {
	return DeviceMemory{
		ptr: C.CUdeviceptr(uintptr(dptr)),
	}
}

func AllocDeviceMemory(size int) (DeviceMemory, error) {
	var ptr C.CUdeviceptr
	ce := C.cuMemAlloc(&ptr, C.size_t(size))
	if ce != 0 {
		return DeviceMemory{}, newCudaError(ce)
	}
	return DeviceMemory{owned: true, ptr: ptr}, nil
}

func (m DeviceMemory) IsNil() bool {
	return m.ptr == 0
}

func (m DeviceMemory) Free() error {
	if !m.owned {
		panic("device memory not owned")
	}
	if m.ptr == 0 {
		return nil
	}
	ce := C.cuMemFree(m.ptr)
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToDevice(dst DeviceMemory, size int) error {
	ce := C.cuMemcpyDtoD(dst.ptr, m.ptr, C.size_t(size))
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToHost(dst HostMemory, size int) error {
	ce := C.cuMemcpyDtoH(dst.unsafePtr(), m.ptr, C.size_t(size))
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToDeviceAsync(dst DeviceMemory, size int, stream Stream) error {
	ce := C.cuMemcpyDtoDAsync(dst.ptr, m.ptr, C.size_t(size), stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToHostAsync(dst HostMemory, size int, stream Stream) error {
	ce := C.cuMemcpyDtoHAsync(dst.unsafePtr(), m.ptr, C.size_t(size), stream.s)
	return newCudaError(ce)
}

func paramsToM2D(
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
) C.CUDA_MEMCPY2D {
	var m2d C.CUDA_MEMCPY2D
	m2d.srcXInBytes = C.size_t(srcXInBytes)
	m2d.srcY = C.size_t(srcY)
	m2d.srcPitch = C.size_t(srcPitch)

	m2d.dstXInBytes = C.size_t(dstXInBytes)
	m2d.dstY = C.size_t(dstY)
	m2d.dstPitch = C.size_t(dstPitch)

	m2d.WidthInBytes = C.size_t(widthInBytes)
	m2d.Height = C.size_t(height)

	return m2d
}

func (m DeviceMemory) CopyToHost2D(dst HostMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.srcDevice = m.ptr

	m2d.dstMemoryType = C.CU_MEMORYTYPE_HOST
	m2d.dstHost = dst.unsafePtr()

	ce := C._wrapCuMemcpy2D(m2d)
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToHost2DAsync(dst HostMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
	stream Stream,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.srcDevice = m.ptr

	m2d.dstMemoryType = C.CU_MEMORYTYPE_HOST
	m2d.dstHost = dst.unsafePtr()

	ce := C._wrapCuMemcpy2DAsync(m2d, stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToDevice2D(dst DeviceMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.srcDevice = m.ptr

	m2d.dstMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.dstDevice = dst.ptr

	ce := C._wrapCuMemcpy2D(m2d)
	return newCudaError(ce)
}

func (m DeviceMemory) CopyToDevice2DAsync(dst DeviceMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
	stream Stream,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.srcDevice = m.ptr

	m2d.dstMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.dstDevice = dst.ptr

	ce := C._wrapCuMemcpy2DAsync(m2d, stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) Memset8(c byte, n int) error {
	ce := C.cuMemsetD8(m.ptr, C.uchar(c), C.size_t(n))
	return newCudaError(ce)
}

func (m DeviceMemory) Memset8Async(c byte, n int, stream Stream) error {
	ce := C.cuMemsetD8Async(m.ptr, C.uchar(c), C.size_t(n), stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) Memset16(c uint16, n int) error {
	ce := C.cuMemsetD16(m.ptr, C.ushort(c), C.size_t(n))
	return newCudaError(ce)
}

func (m DeviceMemory) Memset16Async(c uint16, n int, stream Stream) error {
	ce := C.cuMemsetD16Async(m.ptr, C.ushort(c), C.size_t(n), stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) Memset32(c uint32, n int) error {
	ce := C.cuMemsetD32(m.ptr, C.uint(c), C.size_t(n))
	return newCudaError(ce)
}

func (m DeviceMemory) Memset32Async(c uint32, n int, stream Stream) error {
	ce := C.cuMemsetD32Async(m.ptr, C.uint(c), C.size_t(n), stream.s)
	return newCudaError(ce)
}

func (m DeviceMemory) UnsafePtr() unsafe.Pointer {
	// nolint
	return unsafe.Pointer(uintptr(m.ptr))
}

type HostMemory struct {
	owned bool
	// managed bool

	ptr   unsafe.Pointer
	goPtr []byte
}

func AllocLockedHostMemory(size int) (HostMemory, error) {
	var p unsafe.Pointer
	ce := C.cuMemAllocHost(&p, C.size_t(size))
	if ce != 0 {
		return HostMemory{}, newCudaError(ce)
	}
	return HostMemory{
		owned: true,

		ptr: p,
	}, nil
}

func HostMemoryFromBytes(data []byte) HostMemory {
	return HostMemory{
		goPtr: data,
	}
}

func HostMemoryFromUnsafePointer(ptr unsafe.Pointer) HostMemory {
	return HostMemory{
		ptr: ptr,
	}
}

func AllocHostMemory(size int) (HostMemory, error) {
	return HostMemory{
		owned: true,

		goPtr: make([]byte, size),
	}, nil
}

func (m HostMemory) Locked() bool {
	return m.ptr != nil
}

func (m HostMemory) IsNil() bool {
	return m.unsafePtr() == nil
}

func (m HostMemory) unsafePtr() unsafe.Pointer {
	if m.Locked() {
		return m.ptr
	}
	if len(m.goPtr) == 0 {
		return nil
	}
	// #nosec
	return unsafe.Pointer(&m.goPtr[0])
}

func (m HostMemory) UnsafePtr() unsafe.Pointer {
	return m.unsafePtr()
}

func (m HostMemory) GoBytes() []byte {
	if m.Locked() {
		panic("locked memory")
	}
	return m.goPtr
}

func (m HostMemory) Free() error {
	if !m.owned {
		panic("host memory not owned")
	}
	if m.Locked() {
		ce := C.cuMemFreeHost(m.ptr)
		return newCudaError(ce)
	}
	return nil
}

func (m HostMemory) CopyToDevice(dst DeviceMemory, size int) error {
	ce := C.cuMemcpyHtoD(dst.ptr, m.unsafePtr(), C.size_t(size))
	return newCudaError(ce)
}

func (m HostMemory) CopyToHost(dst HostMemory, size int) error {
	// ce := C.cuMemcpyHtoH(dst.unsafePtr(), m.unsafePtr(), C.size_t(size))
	// return newCudaError(ce)
	C.memcpy(dst.unsafePtr(), m.unsafePtr(), C.size_t(size))
	return nil
}

func (m HostMemory) CopyToDeviceAsync(dst DeviceMemory, size int, stream Stream) error {
	ce := C.cuMemcpyHtoDAsync(dst.ptr, m.unsafePtr(), C.size_t(size), stream.s)
	return newCudaError(ce)
}

func (m HostMemory) CopyToDevice2D(dst DeviceMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_HOST
	m2d.srcHost = m.unsafePtr()

	m2d.dstMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.dstDevice = dst.ptr

	ce := C._wrapCuMemcpy2D(m2d)
	return newCudaError(ce)
}

func (m HostMemory) CopyToDevice2DAsync(dst DeviceMemory,
	srcXInBytes, srcY, srcPitch int,
	dstXInBytes, dstY, dstPitch int,
	widthInBytes, height int,
	stream Stream,
) error {
	m2d := paramsToM2D(srcXInBytes, srcY, srcPitch, dstXInBytes, dstY, dstPitch, widthInBytes, height)

	m2d.srcMemoryType = C.CU_MEMORYTYPE_HOST
	m2d.srcHost = m.unsafePtr()

	m2d.dstMemoryType = C.CU_MEMORYTYPE_DEVICE
	m2d.dstDevice = dst.ptr

	ce := C._wrapCuMemcpy2DAsync(m2d, stream.s)
	return newCudaError(ce)
}

func ProfilerStart() error {
	ce := C.cuProfilerStart()
	return newCudaError(ce)
}

func ProfilerStop() error {
	ce := C.cuProfilerStop()
	return newCudaError(ce)
}
