// +build gsc

package gsc

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L../../libs -Llibs -lgpu_search_core -lopenblas
#include <stdio.h>
#include <stdlib.h>
#include "gpu_search_core.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/dsc"
)

// nolint: golint
type CpuIndex struct {
	handle C.gsc_cpu_index_t
}

type GpuIndex struct {
	handle C.gsc_gpu_index_t
}

type GpuResource struct {
	handle C.gsc_resource_t
	Device *GpuDevice
}

type GpuDevice struct {
	handle C.int
}

type GpuInfo struct {
	handle                 C.gsc_gpu_info_t
	AvailableGpuMemorySize uint64
	TotalGpuMemorySize     uint64
}

// make sure all methods implement
var _ dsc.SearchIndex = &CpuIndex{}
var _ dsc.SearchIndex = &GpuIndex{}

func convertFilterByTimespaceFilter(filter *dsc.TimespaceFilter) (C.gsc_timespace_filter_t, bool) {
	var f C.gsc_timespace_filter_t
	var ignoreFilter bool = true
	if filter.TimeRange[0] != 0 || filter.TimeRange[1] != 0 {
		ignoreFilter = false
	}
	f.time_range[0] = C.uint(filter.TimeRange[0])
	f.time_range[1] = C.uint(filter.TimeRange[1])
	for i := 0; i < len(filter.CameraMask); i++ {
		if filter.CameraMask[i] != 0xff {
			ignoreFilter = false
		}
		f.camera_mask[i] = C.uchar(filter.CameraMask[i])
	}
	return f, ignoreFilter
}

func (index *GpuIndex) WriteToFile(filepath string) error {
	cpuIndex, err := index.WriteToCpuIndex(true)
	if err != nil {
		return err
	}
	return cpuIndex.WriteToFile(filepath)
}

func (index *CpuIndex) WriteToFile(filepath string) error {
	cfilepath := C.CString(filepath)
	// #nosec
	defer C.free(unsafe.Pointer(cfilepath))
	if err := int(C.gsc_write_cpu_index(index.handle, cfilepath)); err != 0 {
		return fmt.Errorf("gsc_write_cpu_index error code: %d", err)
	}
	return nil
}

func (index *GpuIndex) AddBatch(n int64, x []float32, ids []int64) error {
	if err := int(C.gsc_add_index_batch(index.handle, C.long(n), (*C.float)(&x[0]), (*C.long)(&ids[0]))); err != 0 {
		return fmt.Errorf("gsc_add_index_batch error code: %d", err)
	}
	return nil
}

func (index *CpuIndex) AddBatch(n int64, x []float32, ids []int64) error {
	if err := int(C.gsc_add_cpu_index_batch(index.handle, C.long(n), (*C.float)(&x[0]), (*C.long)(&ids[0]))); err != 0 {
		return fmt.Errorf("gsc_add_cpu_index_batch error code: %d", err)
	}
	return nil
}

func (index *GpuIndex) Search(n int64, x []float32, k int32) ([]float32, []int64, error) {
	if k > 1024 {
		k = 1024
	}
	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))
	if err :=
		int(C.gsc_search_index(index.handle, C.long(n), (*C.float)(&x[0]), C.int(k), (*C.float)(&distances[0]), (*C.long)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("gsc_search_index error code: %d", err)
	}
	return distances, ids, nil
}

// nolint: dupl
func (index *GpuIndex) SearchWithTimespace(n int64, x []float32, k int32, filter []*dsc.TimespaceFilter) ([]float32, []int64, error) {
	if len(filter) > 1 {
		return nil, nil, fmt.Errorf("gsc_search_index_with_timespace_filter error, filter more than 1")
	} else if len(filter) < 1 {
		return nil, nil, fmt.Errorf("gsc_search_index_with_timespace_filter error, filter is nil")
	}
	f, _ := convertFilterByTimespaceFilter(filter[0])

	if k > 1024 {
		k = 1024
	}
	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))
	if err :=
		int(C.gsc_search_index_with_timespace_filter(index.handle, C.long(n), (*C.float)(&x[0]), C.int(k), &f, (*C.float)(&distances[0]), (*C.long)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("gsc_search_index_with_timespace_filter error code: %d", err)
	}
	return distances, ids, nil
}

func (index *CpuIndex) Search(n int64, x []float32, k int32) ([]float32, []int64, error) {
	if k > 1024 {
		k = 1024
	}
	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))
	if err :=
		int(C.gsc_search_cpu_index(index.handle, C.long(n), (*C.float)(&x[0]), C.int(k), (*C.float)(&distances[0]), (*C.long)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("gsc_search_cpu_index error code: %d", err)
	}
	return distances, ids, nil
}

// nolint: dupl
func (index *CpuIndex) SearchWithTimespace(n int64, x []float32, k int32, filter []*dsc.TimespaceFilter) ([]float32, []int64, error) {
	if len(filter) > 1 {
		return nil, nil, fmt.Errorf("gsc_search_cpu_index_with_timespace_filter error, filter more than 1")
	} else if len(filter) < 1 {
		return nil, nil, fmt.Errorf("gsc_search_cpu_index_with_timespace_filter error, filter is nil")
	}
	f, _ := convertFilterByTimespaceFilter(filter[0])

	if k > 1024 {
		k = 1024
	}
	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))
	if err :=
		int(C.gsc_search_cpu_index_with_timespace_filter(index.handle, C.long(n), (*C.float)(&x[0]), C.int(k), &f, (*C.float)(&distances[0]), (*C.long)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("gsc_search_cpu_index_with_timespace_filter error code: %d", err)
	}
	return distances, ids, nil
}

// nolint: golint
func (index *GpuIndex) WriteToCpuIndex(reclaimMemory bool) (dsc.SearchIndex, error) {
	cpuIndex := &CpuIndex{}
	r := 0
	if reclaimMemory {
		r = 1
	}
	if err := int(C.gsc_index_gpu_to_cpu(index.handle, &cpuIndex.handle, C.int(r))); err != 0 {
		return nil, fmt.Errorf("gsc_index_gpu_to_cpu error code: %d", err)
	}
	return cpuIndex, nil
}

// nolint: golint
// WriteToCpuIndex :Cpu mode just clone the index
func (index *CpuIndex) WriteToCpuIndex(reclaimMemory bool) (dsc.SearchIndex, error) {
	if index == nil {
		return nil, fmt.Errorf("writeToCpuIndex Fail: nil CpuIndex")
	}
	cpuIndex := &CpuIndex{}
	if err := int(C.gsc_clone_cpu_index(index.handle, &cpuIndex.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_clone_cpu_index error code: %d", err)
	}
	return cpuIndex, nil
}

func (index *CpuIndex) WriteToDeviceIndex(resource interface{}) (dsc.SearchIndex, error) {
	r, ok := resource.(*GpuResource)
	if !ok {
		return nil, fmt.Errorf("gsc_index_cpu_to_gpu error, invalid gpu resource")
	}
	gpuIndex := &GpuIndex{}
	if err := int(C.gsc_index_cpu_to_gpu(index.handle, r.handle, r.Device.handle, &gpuIndex.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_index_cpu_to_gpu error code: %d", err)
	}

	return gpuIndex, nil
}

// WriteToDeviceIndex Gpu mode do not need to write index to gpu, just for abstract
func (index *GpuIndex) WriteToDeviceIndex(resource interface{}) (dsc.SearchIndex, error) {
	if index == nil {
		return nil, fmt.Errorf("WriteToDeviceIndex Fail: nil GpuIndex")
	}
	return index, nil
}

func (index *GpuIndex) FreeIndex() error {
	if index.handle != nil {
		if err := C.gsc_free_gpu_index(index.handle); err != 0 {
			return fmt.Errorf("gsc_free_gpu_index error code: %d", err)
		}
		index.handle = nil
	}
	return nil
}

func (index *CpuIndex) FreeIndex() error {
	if index.handle != nil {
		if err := C.gsc_free_cpu_index(index.handle); err != 0 {
			return fmt.Errorf("gsc_free_cpu_index error code: %d", err)
		}
		index.handle = nil
	}
	return nil
}

func (index *GpuIndex) GetIndexInfo() (*dsc.IndexInfo, error) {
	status := C.gsc_index_status_t{}
	if err := C.gsc_get_gpu_index_status(index.handle, &status); err != 0 {
		return nil, fmt.Errorf("gsc_get_gpu_index_status error code: %d", err)
	}
	gpuIndexInfo := &dsc.IndexInfo{}
	gpuIndexInfo.Size = uint64(status.index_size)
	gpuIndexInfo.Dimension = int32(status.dimension)
	gpuIndexInfo.BitsPerCode = int32(status.bitsPerCode)
	gpuIndexInfo.Nprobe = int32(status.nprobe)
	gpuIndexInfo.SubQuantizers = int32(status.subQuantizers)
	gpuIndexInfo.Nlist = int32(status.nlist)
	gpuIndexInfo.MaxListSize = int32(status.max_list_size)
	gpuIndexInfo.SDKType = dsc.GSC
	gpuIndexInfo.IndexType = dsc.IndexIVFPQ
	gpuIndexInfo.IsCPUModel = false
	if int(status.is_trained) > 0 {
		gpuIndexInfo.IsTrained = true
	} else {
		gpuIndexInfo.IsTrained = false
	}
	return gpuIndexInfo, nil
}

func (index *CpuIndex) GetIndexInfo() (*dsc.IndexInfo, error) {
	status := C.gsc_index_status_t{}
	if err := C.gsc_get_cpu_index_status(index.handle, &status); err != 0 {
		return nil, fmt.Errorf("gsc_get_cpu_index_status error code: %d", err)
	}
	cpuIndexInfo := &dsc.IndexInfo{}
	cpuIndexInfo.Size = uint64(status.index_size)
	cpuIndexInfo.Dimension = int32(status.dimension)
	cpuIndexInfo.BitsPerCode = int32(status.bitsPerCode)
	cpuIndexInfo.Nprobe = int32(status.nprobe)
	cpuIndexInfo.SubQuantizers = int32(status.subQuantizers)
	cpuIndexInfo.Nlist = int32(status.nlist)
	cpuIndexInfo.MaxListSize = int32(status.max_list_size)
	cpuIndexInfo.SDKType = dsc.GSC
	cpuIndexInfo.IndexType = dsc.IndexIVFPQ
	cpuIndexInfo.IsCPUModel = true
	if int(status.is_trained) > 0 {
		cpuIndexInfo.IsTrained = true
	} else {
		cpuIndexInfo.IsTrained = false
	}
	return cpuIndexInfo, nil
}

func (index *GpuIndex) GetIndexIds() ([]int64, error) {
	info, err := index.GetIndexInfo()
	if err != nil {
		return nil, err
	}
	if info.Size == 0 {
		return nil, nil
	}
	ids := make([]int64, info.Size)
	if err := int(C.gsc_get_gpu_index_ids(index.handle, (*C.long)(&ids[0]))); err != 0 {
		return nil, fmt.Errorf("gsc_get_gpu_index_ids error code: %d", err)
	}
	return ids, nil
}

func (index *CpuIndex) GetIndexIds() ([]int64, error) {
	info, err := index.GetIndexInfo()
	if err != nil {
		return nil, err
	}
	if info.Size == 0 {
		return nil, nil
	}
	ids := make([]int64, info.Size)
	if err := int(C.gsc_get_cpu_index_ids(index.handle, (*C.long)(&ids[0]))); err != 0 {
		return nil, fmt.Errorf("gsc_get_cpu_index_ids error code: %d", err)
	}
	return ids, nil
}

//RemoveIndexIds GpuIndex can not remove feature by ids, just for abstract
func (index *GpuIndex) RemoveIndexIds([]int64) error {
	return nil
}

func (index *CpuIndex) RemoveIndexIds(ids []int64) error {
	n := len(ids)
	if n == 0 {
		return nil
	}

	if err := int(C.gsc_remove_cpu_index_ids(index.handle, C.long(n), (*C.long)(&ids[0]))); err != 0 {
		return fmt.Errorf("gsc_remove_cpu_index_ids error code: %d", err)
	}
	return nil
}

func (index *GpuIndex) GetMaxMayReserveMemory(n int64) int64 {
	return int64(C.gsc_get_gpu_max_may_reserve_memory(index.handle, C.long(n)))
}

func (resource *GpuResource) LoadGpuIndex(filepath string) (*GpuIndex, error) {
	index := &GpuIndex{}
	cfilepath := C.CString(filepath)
	// #nosec
	defer C.free(unsafe.Pointer(cfilepath))
	if err := int(C.gsc_load_gpu_index(resource.handle, cfilepath, resource.Device.handle, &index.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_load_gpu_index error code: %d", err)
	}
	return index, nil
}

// nolint: golint
func LoadCpuIndex(filepath string) (*CpuIndex, error) {
	index := &CpuIndex{}
	cfilepath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cfilepath))
	if err := int(C.gsc_load_cpu_index(cfilepath, &index.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_load_cpu_index error code: %d", err)
	}
	return index, nil
}

func (resource *GpuResource) InitGpuIndex(indexConfig *dsc.GeneralIndexConfig) (*GpuIndex, error) {
	index := &GpuIndex{}
	config := C.gsc_index_config_t{}
	config.dimension = C.int(indexConfig.Dimension)
	config.nlist = C.int(indexConfig.Nlist)
	config.subQuantizers = C.int(indexConfig.SubQuantizers)
	config.bitsPerCode = C.int(indexConfig.BitsPerCode)
	config.nprobe = C.int(indexConfig.Nprobe)
	if err := int(C.gsc_init_gpu_index(resource.handle, &config, &index.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_init_gpu_index error code: %d", err)
	}
	return index, nil

}

// nolint: golint
func InitCpuIndex(indexConfig *dsc.GeneralIndexConfig) (*CpuIndex, error) {
	index := &CpuIndex{}
	config := C.gsc_index_config_t{}
	config.dimension = C.int(indexConfig.Dimension)
	config.nlist = C.int(indexConfig.Nlist)
	config.subQuantizers = C.int(indexConfig.SubQuantizers)
	config.bitsPerCode = C.int(indexConfig.BitsPerCode)
	config.nprobe = C.int(indexConfig.Nprobe)
	if err := int(C.gsc_init_cpu_index(&config, &index.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_init_cpu_index error code: %d", err)
	}
	return index, nil

}

func (index *GpuIndex) TrainIndex(n int64, fs []float32) error {
	if err := int(C.gsc_train_gpu_index(index.handle, C.long(n), (*C.float)(&fs[0]))); err != 0 {
		return fmt.Errorf("gsc_train_gpu_index error code: %d", err)
	}
	return nil
}

func (index *CpuIndex) TrainIndex(n int64, fs []float32) error {
	if err := int(C.gsc_train_cpu_index(index.handle, C.long(n), (*C.float)(&fs[0]))); err != 0 {
		return fmt.Errorf("gsc_train_cpu_index error code: %d", err)
	}
	return nil
}

func (device *GpuDevice) GetGpuInfo() (*GpuInfo, error) {
	gpuInfo := &GpuInfo{}
	if err := int(C.gsc_get_gpu_info(device.handle, &gpuInfo.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_get_gpu_info error code: %d", err)
	}
	gpuInfo.AvailableGpuMemorySize = uint64(gpuInfo.handle.free_memory_size)
	gpuInfo.TotalGpuMemorySize = uint64(gpuInfo.handle.total_memory_size)
	return gpuInfo, nil
}

func (device *GpuDevice) CreateGpuResource(memory uint64) (*GpuResource, error) {
	resource := &GpuResource{
		Device: device,
	}
	// 300MB Temp Memory
	if err := int(C.gsc_create_gpu_resource(device.handle, C.ulong(memory), 0, &resource.handle)); err != 0 {
		return nil, fmt.Errorf("gsc_create_gpu_resource error code: %d", err)
	}
	return resource, nil
}

func (resource *GpuResource) DestroyGpuResource() error {
	if resource.handle != nil {
		if err := C.gsc_destroy_gpu_resource(resource.handle); err != 0 {
			return fmt.Errorf("gsc_destroy_gpu_resource error code: %d", err)
		}
		resource.handle = nil
	}
	return nil
}

func GetGpuDevice() (*GpuDevice, error) {
	return &GpuDevice{
		handle: C.gsc_get_current_device_id(),
	}, nil
}

type DeviceProperties struct {
	Major int
	Minor int
	Name  string
}

func GetDeviceProperties(deviceID int) (*DeviceProperties, error) {
	var deviceProperties DeviceProperties
	var cp C.struct_gsc_device_properties

	if err := C.gsc_get_device_properties(&cp, C.int(deviceID)); err != 0 {
		return &deviceProperties, fmt.Errorf("gsc_get_device_properties error code: %d", err)
	}
	deviceProperties.Major = int(cp.major)
	deviceProperties.Minor = int(cp.minor)
	deviceProperties.Name = C.GoString(&(cp.name[0]))
	return &deviceProperties, nil
}
