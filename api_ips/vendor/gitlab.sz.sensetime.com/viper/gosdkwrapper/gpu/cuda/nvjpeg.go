package cuda

/*
#include "dynlink_nvjpeg.h"

nvjpegStatus_t _nvjpegCreate(nvjpegBackend_t backend, nvjpegDevAllocator_t *allocator, nvjpegHandle_t *handle) {
	return nvjpegFunctions->nvjpegCreateEx(backend, allocator, NULL, 0, handle);
}

nvjpegStatus_t _nvjpegDestroy(nvjpegHandle_t handle) {
	return nvjpegFunctions->nvjpegDestroy(handle);
}

nvjpegStatus_t _nvjpegJpegStateCreate(nvjpegHandle_t handle, nvjpegJpegState_t *jpeg_handle) {
	return nvjpegFunctions->nvjpegJpegStateCreate(handle, jpeg_handle);
}

nvjpegStatus_t _nvjpegJpegStateDestroy(nvjpegJpegState_t jpeg_handle) {
	return nvjpegFunctions->nvjpegJpegStateDestroy(jpeg_handle);
}

nvjpegStatus_t _nvjpegDecodeBatchedInitialize(
          nvjpegHandle_t handle,
          nvjpegJpegState_t jpeg_handle,
          int batch_size,
          int max_cpu_threads,
          nvjpegOutputFormat_t output_format) {
	return nvjpegFunctions->nvjpegDecodeBatchedInitialize(handle, jpeg_handle, batch_size, max_cpu_threads, output_format);
}

nvjpegStatus_t _nvjpegDecodeJpegHost(
		nvjpegHandle_t handle,
		nvjpegJpegDecoder_t decoder,
		nvjpegJpegState_t decoder_state,
		nvjpegDecodeParams_t decode_params,
		nvjpegJpegStream_t jpeg_stream) {
	return nvjpegFunctions->nvjpegDecodeJpegHost(handle, decoder, decoder_state, decode_params, jpeg_stream);
}

nvjpegStatus_t _nvjpegDecodeJpegTransferToDevice(
		nvjpegHandle_t handle,
		nvjpegJpegDecoder_t decoder,
		nvjpegJpegState_t decoder_state,
		nvjpegJpegStream_t jpeg_stream,
		cudaStream_t stream) {
	return nvjpegFunctions->nvjpegDecodeJpegTransferToDevice(handle, decoder, decoder_state, jpeg_stream, stream);
}

nvjpegStatus_t _nvjpegDecodeJpegDevice(
		nvjpegHandle_t handle,
		nvjpegJpegDecoder_t decoder,
		nvjpegJpegState_t decoder_state,
		nvjpegImage_t *destination,
		cudaStream_t stream) {
	return nvjpegFunctions->nvjpegDecodeJpegDevice(handle, decoder, decoder_state, destination, stream);
}

nvjpegStatus_t _nvjpegJpegStreamCreate(nvjpegHandle_t handle, nvjpegJpegStream_t *jpeg_stream) {
	return nvjpegFunctions->nvjpegJpegStreamCreate(handle, jpeg_stream);
}

nvjpegStatus_t _nvjpegJpegStreamDestroy(nvjpegJpegStream_t jpeg_stream) {
	return nvjpegFunctions->nvjpegJpegStreamDestroy(jpeg_stream);
}

nvjpegStatus_t _nvjpegDecoderCreate(
		nvjpegHandle_t nvjpeg_handle,
		nvjpegBackend_t implementation,
		nvjpegJpegDecoder_t* decoder_handle) {
	return nvjpegFunctions->nvjpegDecoderCreate(nvjpeg_handle, implementation, decoder_handle);
}

nvjpegStatus_t _nvjpegDecoderDestroy(nvjpegJpegDecoder_t decoder_handle) {
	return nvjpegFunctions->nvjpegDecoderDestroy(decoder_handle);
}

nvjpegStatus_t _nvjpegJpegStreamParse(
		nvjpegHandle_t handle,
		const unsigned char *data,
		size_t length,
		int save_metadata,
		int save_stream,
		nvjpegJpegStream_t jpeg_stream) {
	return nvjpegFunctions->nvjpegJpegStreamParse(handle, data, length, save_metadata, save_stream, jpeg_stream);
}

nvjpegStatus_t _nvjpegDecoderStateCreate(
		nvjpegHandle_t nvjpeg_handle,
		nvjpegJpegDecoder_t decoder_handle,
		nvjpegJpegState_t* decoder_state) {
	return nvjpegFunctions->nvjpegDecoderStateCreate(nvjpeg_handle, decoder_handle, decoder_state);
}

nvjpegStatus_t _nvjpegDecodeParamsCreate(
		nvjpegHandle_t handle,
		nvjpegDecodeParams_t *decode_params) {
	return nvjpegFunctions->nvjpegDecodeParamsCreate(handle, decode_params);
}

nvjpegStatus_t _nvjpegDecodeParamsDestroy(nvjpegDecodeParams_t decode_params) {
	return nvjpegFunctions->nvjpegDecodeParamsDestroy(decode_params);
}

nvjpegStatus_t _nvjpegDecodeParamsSetOutputFormat(
		nvjpegDecodeParams_t decode_params,
		nvjpegOutputFormat_t output_format) {
	return nvjpegFunctions->nvjpegDecodeParamsSetOutputFormat(decode_params, output_format);
}

nvjpegStatus_t _nvjpegBufferDeviceCreate(
	nvjpegHandle_t handle,
	nvjpegBufferDevice_t* buffer) {
	return nvjpegFunctions->nvjpegBufferDeviceCreate(handle, NULL, buffer);
}

nvjpegStatus_t _nvjpegBufferDeviceDestroy(nvjpegBufferDevice_t buffer) {
	return nvjpegFunctions->nvjpegBufferDeviceDestroy(buffer);
}

nvjpegStatus_t _nvjpegStateAttachDeviceBuffer(
		nvjpegJpegState_t decoder_state,
		nvjpegBufferDevice_t device_buffer) {
	return nvjpegFunctions->nvjpegStateAttachDeviceBuffer(decoder_state, device_buffer);
}

nvjpegStatus_t _nvjpegGetImageInfo(
          nvjpegHandle_t handle,
          const unsigned char *data,
          size_t length,
          int *nComponents,
          nvjpegChromaSubsampling_t *subsampling,
          int *widths,
          int *heights) {
	return nvjpegFunctions->nvjpegGetImageInfo(handle, data, length, nComponents, subsampling, widths, heights);
}

*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/common"
	kestrel "gitlab.sz.sensetime.com/viper/gosdkwrapper/kestrelV1/KestrelGo"
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/kestrelV1/keson"
)

type NVJPEGError struct {
	Code int
}

func (e NVJPEGError) Error() string {
	return fmt.Sprintf("NVJPEG Error: %d", e.Code)
}

func newNVJPEGError(s C.nvjpegStatus_t) error {
	if s == 0 {
		return nil
	}
	return NVJPEGError{
		Code: int(s),
	}
}

type NVJPEGImage struct {
	dm           DeviceMemory
	w, h, stride int

	pool *MemPool
}

func (i *NVJPEGImage) OriginInfo() *common.OriginInfo {
	return nil
}

func (i *NVJPEGImage) Data() unsafe.Pointer {
	return i.dm.UnsafePtr()
}

func (i *NVJPEGImage) UserData() unsafe.Pointer {
	return nil
}

func (i *NVJPEGImage) ExtraInfoBuffer() unsafe.Pointer {
	return nil
}

func (i *NVJPEGImage) DeviceType() common.DeviceType {
	return common.DeviceGPU
}

func (i *NVJPEGImage) Width() int {
	return i.w
}

func (i *NVJPEGImage) Height() int {
	return i.h
}

func (i *NVJPEGImage) Stride() int {
	return i.stride
}

func (i *NVJPEGImage) CommonPixelFormat() common.PixelFormat {
	return common.PIXEL_FORMAT_BGR
}

func (i *NVJPEGImage) Crop(x int, y int, w int, h int) (common.Image, error) {
	strides := [4]int{i.Stride(), 0, 0, 0}
	tmpFrame := kestrel.NewFrameFromMem(kestrel.DeviceMem, kestrel.PixelFormatBGR, i.Data(), i.Width(), i.Height(), strides, 0)
	defer tmpFrame.Release()

	img, err := tmpFrame.Crop(x, y, w, h)
	if err != nil {
		return nil, err
	}
	return &keson.Frame{Frame: img}, nil
}

func (i *NVJPEGImage) Release() {
	if i.pool == nil {
		i.dm.Free() // nolint
	} else {
		i.pool.Free(i.dm) // nolint: errcheck
	}
}

type NVJPEGDecoder struct {
	h       C.nvjpegHandle_t
	p       C.nvjpegDecodeParams_t
	stream  Stream
	mempool *MemPool
}

func NewNVJPEGDecoder(blobsize int, poolsize int) (*NVJPEGDecoder, error) {
	stream, err := CreateStream(0)
	if err != nil {
		return nil, err
	}

	var pool *MemPool
	if poolsize > 0 {
		pool, err = NewMemPool(blobsize, poolsize)
		if err != nil {
			log.Println("ERROR: failed to create pool: ", err)
			stream.Destroy() // nolint: errcheck
			return nil, err
		}
	}

	var h C.nvjpegHandle_t
	alloc := C._nvjpegGetMemoryPoolAllocator()
	rc := C._nvjpegCreate(C.NVJPEG_BACKEND_DEFAULT, alloc, &h)
	if rc != 0 {
		stream.Destroy() // nolint: errcheck
		return nil, newNVJPEGError(rc)
	}

	var params C.nvjpegDecodeParams_t
	if rt := C._nvjpegDecodeParamsCreate(h, &params); rt != 0 {
		log.Println("create decode params failed", rt)
		return nil, newNVJPEGError(rt)
	}
	if rt := C._nvjpegDecodeParamsSetOutputFormat(params, C.NVJPEG_OUTPUT_BGRI); rt != 0 {
		log.Println("params set output format failed", rt)
		return nil, newNVJPEGError(rt)
	}

	dec := &NVJPEGDecoder{
		h:       h,
		p:       params,
		stream:  stream,
		mempool: pool,
	}
	return dec, nil
}

func (d *NVJPEGDecoder) getImageInfo(img []byte) (int, int, int, error) {
	if len(img) == 0 {
		return 0, 0, 0, newNVJPEGError(C.NVJPEG_STATUS_BAD_JPEG)
	}
	var nc C.int
	var hs, ws [4]C.int
	var subsampling C.nvjpegChromaSubsampling_t
	rc := C._nvjpegGetImageInfo(d.h, (*C.uchar)(&img[0]), C.size_t(len(img)), &nc, &subsampling, &ws[0], &hs[0])
	if rc != 0 {
		return 0, 0, 0, newNVJPEGError(rc)
	}
	return int(ws[0]), int(hs[0]), int(nc), nil
}

func (d *NVJPEGDecoder) allocDeviceMemory(size int) (DeviceMemory, error) {
	if d.mempool == nil {
		return AllocDeviceMemory(size)
	}
	return d.mempool.Alloc(size)
}

func (d *NVJPEGDecoder) Decode(data []byte) (common.Image, error) {
	w, h, _, err := d.getImageInfo(data)
	if err != nil {
		return nil, err
	}
	s := w * 3
	size := s * h

	var stream C.nvjpegJpegStream_t
	if rt := C._nvjpegJpegStreamCreate(d.h, &stream); rt != 0 {
		log.Println("create jpeg stream failed", rt)
		return nil, newNVJPEGError(rt)
	}
	defer func() { C._nvjpegJpegStreamDestroy(stream) }()

	var decoder C.nvjpegJpegDecoder_t
	if rt := C._nvjpegDecoderCreate(d.h, C.NVJPEG_BACKEND_DEFAULT, &decoder); rt != 0 {
		log.Println("create jpeg decoder failed", rt)
		return nil, newNVJPEGError(rt)
	}
	defer func() { C._nvjpegDecoderDestroy(decoder) }()

	var state C.nvjpegJpegState_t
	if rt := C._nvjpegDecoderStateCreate(d.h, decoder, &state); rt != 0 {
		log.Println("create jpeg decoder state failed", rt)
		return nil, newNVJPEGError(rt)
	}
	defer func() { C._nvjpegJpegStateDestroy(state) }()

	var buffer C.nvjpegBufferDevice_t
	if rt := C._nvjpegBufferDeviceCreate(d.h, &buffer); rt != 0 {
		log.Println("create jpeg decoder state failed", rt)
		return nil, newNVJPEGError(rt)
	}
	defer func() { C._nvjpegBufferDeviceDestroy(buffer) }()

	if rt := C._nvjpegJpegStreamParse(d.h, (*C.uchar)(&data[0]), C.size_t(len(data)), 0, 0, stream); rt != 0 {
		log.Println("jpeg stream parse failed", rt)
		return nil, newNVJPEGError(rt)
	}

	if rt := C._nvjpegDecodeJpegHost(d.h, decoder, state, d.p, stream); rt != 0 {
		log.Println("decode jpeg host failed", rt)
		return nil, newNVJPEGError(rt)
	}

	if rt := C._nvjpegStateAttachDeviceBuffer(state, buffer); rt != 0 {
		log.Println("attach device buffer failed", rt)
		return nil, newNVJPEGError(rt)
	}

	if rt := C._nvjpegDecodeJpegTransferToDevice(d.h, decoder, state, stream, d.stream.s); rt != 0 {
		log.Println("transfer to device failed", rt)
		return nil, newNVJPEGError(rt)
	}

	dm, err := d.allocDeviceMemory(size)
	if err == nil {
		hm := HostMemoryFromUnsafePointer(dm.UnsafePtr())
		err = hm.CopyToDevice(dm, size)
	}

	var img C.nvjpegImage_t
	img.pitch[0] = C.uint(s)
	img.channel[0] = (*C.uchar)(dm.UnsafePtr())
	if rt := C._nvjpegDecodeJpegDevice(d.h, decoder, state, &img, d.stream.s); rt != 0 {
		log.Println(" decode jpeg device failed", rt)
		return nil, newNVJPEGError(rt)
	}

	if err != nil {
		if !dm.IsNil() {
			dm.Free() // nolint
		}

		return nil, err
	}

	return &NVJPEGImage{
		dm:     dm,
		w:      w,
		h:      h,
		stride: s,
		pool:   d.mempool,
	}, nil
}

func (d *NVJPEGDecoder) Synchronize() error {
	return d.stream.Synchronize()
}

func (d *NVJPEGDecoder) Close() error {
	if d.p != nil {
		C._nvjpegDecodeParamsDestroy(d.p)
	}

	if d.h != nil {
		C._nvjpegDestroy(d.h)
	}
	d.stream.Destroy() // nolint: errcheck

	if d.mempool != nil {
		d.mempool.Close()
	}
	return nil
}

func NVJPEGInit(flags uint) error {
	rc := C.nvjpegInit(C.uint(flags))
	return newNVJPEGError(rc)
}

func NVJPEGInitMemoryPool(blobsize uint, count int) error {
	rc := C._nvjpegInitMemoryPool(C.uint(blobsize), C.int(count))
	return newNVJPEGError(rc)
}

func NVJPEGDestroyMemoryPool() error {
	rc := C._nvjpegDestroyMemoryPool()
	return newNVJPEGError(rc)
}
