// +build se

package se

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L../../libs -Llibs -lsearch_engine -lopenblas
#include <stdio.h>
#include <stdlib.h>
#include "search_engine_c.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/dsc"
)

type SearchEngineContext interface{}

type SearchEngineIndex struct {
	// isOnlyMem will set to true only when a SearchEngineIndex invoke WriteToCpuIndex.
	// if isOnlyMem is true:
	//     1. SearchEngineIndex's handle is invalid, and only mem is valid, so ONLY WriteToFile can be invoked later.
	//     2. SearchEngineIndex will be recognized as isCPUModel
	isOnlyMem bool
	handle    C.se_index_t
	mem       *C.se_mem_data_t
}

// make sure all methods implement
var _ dsc.SearchIndex = &SearchEngineIndex{}

func convertFilterByTimespaceFilter(filter *dsc.TimespaceFilter) (C.se_id_filter_t, bool) {
	var f C.se_id_filter_t
	var ignoreFilter bool = true
	if filter.TimeRange[0] != 0 || filter.TimeRange[1] != 0 {
		ignoreFilter = false
	}
	f.time_range[0] = C.uint32_t(filter.TimeRange[0])
	f.time_range[1] = C.uint32_t(filter.TimeRange[1])
	for i := 0; i < len(filter.CameraMask); i++ {
		if filter.CameraMask[i] != 0xff {
			ignoreFilter = false
		}
		f.camera_id_mask[i] = C.uint8_t(filter.CameraMask[i])
	}

	return f, ignoreFilter
}

func InitSearchEngineEnv(productName, licPath string) error {
	cLicPath := C.CString(licPath)
	defer C.free(unsafe.Pointer(cLicPath))
	cProdName := C.CString(productName)
	defer C.free(unsafe.Pointer(cProdName))
	if err := int(C.se_init(cProdName, cLicPath)); err != 0 {
		return fmt.Errorf("se_init error, code: %d", err)
	}
	return nil
}

func DestroySearchEngineEnv() {
	C.se_deinit()
}

func SetSearchEngineLogLevel() {
	C.se_set_log_level(C.SE_LL_TRACE)
}

func BindDevice(deviceID int32) error {
	return nil
}

func UnbindDevice() error {
	return nil
}

func InitSearchEngineContext(deviceID int32) (SearchEngineContext, error) {
	var seContext C.se_context_t
	if err := int(C.se_context_create(C.int32_t(deviceID), &seContext)); err != 0 {
		return nil, fmt.Errorf("se_context_create error, code: %d", err)
	}
	return seContext, nil
}

func SetSearchEngineContextReservedMemory(seContext SearchEngineContext, memReserved int32) error {
	ctx, ok := seContext.(C.se_context_t)
	if !ok {
		// XXX: should panic?
		return fmt.Errorf("SetSearchEngineContextReservedMemory fail: invalid se context")
	}
	if err := int(C.se_context_set_reserved_memory_size(ctx, C.int32_t(memReserved))); err != 0 {
		return fmt.Errorf("se_context_set_reserved_memory_size error, code: %d", err)
	}

	return nil
}

func SetSearchEngineContextThreadNum(seContext SearchEngineContext, threadNum int32) error {
	c, ok := seContext.(C.se_context_t)
	if !ok {
		return fmt.Errorf("SetSearchEngineContextThreadNum fail: invalid se context")
	}

	number := C.int32_t(threadNum)
	if err := C.se_context_set_thread_num(c, number); err != 0 {
		return fmt.Errorf("se_context_set_thread_num error, code:  %d", err)
	}

	return nil
}

func DestroySearchEngineContext(seContext SearchEngineContext) error {
	c, ok := seContext.(C.se_context_t)
	if !ok {
		// XXX: should panic?
		return fmt.Errorf("DestroySearchEngineContext fail: invalid se context")
	}
	C.se_context_destroy(c)
	return nil
}

func InitSearchEngineIndex(seContext SearchEngineContext, indexConfig *dsc.GeneralIndexConfig) (*SearchEngineIndex, error) {
	if len(indexConfig.ModelPath) == 0 && indexConfig.IndexType == dsc.IndexDC {
		return nil, fmt.Errorf("InitSearchEngineIndex fail, modelPath not set")
	}

	c, ok := seContext.(C.se_context_t)
	if !ok {
		return nil, fmt.Errorf("InitSearchEngineIndex fail: invalid se context")
	}

	cModelPath := C.CString(indexConfig.ModelPath)
	defer C.free(unsafe.Pointer(cModelPath))
	dim := C.int32_t(indexConfig.Dimension)
	nList := C.int32_t(indexConfig.Nlist)
	nProbe := C.int32_t(indexConfig.Nprobe)

	maxDBSize := C.int32_t(indexConfig.MaxSize)
	M := C.int(indexConfig.M)
	efConstruction := C.int(indexConfig.EfConstruct)
	ef := C.int(indexConfig.Ef)
	useInt8 := C.int(indexConfig.UseInt8)

	se := SearchEngineIndex{}
	switch indexConfig.IndexType {
	case dsc.IndexDC:
		if err := int(C.se_index_dc_create(c, cModelPath, dim, nList, nProbe, &se.handle)); err != 0 {
			return nil, fmt.Errorf("se_index_dc_create error, code:  %d", err)
		}
	case dsc.IndexPQ:
		if err := int(C.se_index_pq_create(c, dim, nList, nProbe, &se.handle)); err != 0 {
			return nil, fmt.Errorf("se_index_pq_create error, code:  %d", err)
		}
	case dsc.IndexHNSW:
		if err := int(C.se_index_hnsw_create(c, dim, maxDBSize, M, efConstruction, ef, useInt8, &se.handle)); err != 0 {
			return nil, fmt.Errorf("se_index_hnsw_create error, code: %d", err)
		}
	default:
		return nil, fmt.Errorf("InitSearchEngineIndex fail, unknow IndexType")
	}

	return &se, nil
}

func LoadSearchEngineIndex(seContext SearchEngineContext, filepath string) (*SearchEngineIndex, error) {
	c, ok := seContext.(C.se_context_t)
	if !ok {
		return nil, fmt.Errorf("LoadSearchEngineIndex fail: invalid se context")
	}
	se := SearchEngineIndex{}
	cfilepath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cfilepath))
	if err := C.se_index_create_from_file(c, cfilepath, &se.handle); err != 0 {
		return nil, fmt.Errorf("se_index_create_from_file error, code: %d", err)
	}
	return &se, nil
}

func (index *SearchEngineIndex) AddBatch(n int64, x []float32, ids []int64) error {

	if err := int(C.se_index_add_features(index.handle, (*C.float)(&x[0]), (*C.int64_t)(&ids[0]), C.int64_t(n))); err != 0 {
		return fmt.Errorf("se_index_add_features error, code: %d", err)
	}
	return nil
}

func (index *SearchEngineIndex) Search(n int64, x []float32, k int32) ([]float32, []int64, error) {
	config := C.se_search_config_t{}
	config.k = C.int64_t(k)
	config.threshold = C.float(0)
	config.batch_size = C.int32_t(n)
	config.dist_type = C.DISTANCE_L2
	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))

	if err := int(C.se_index_search(index.handle, (*C.float)(&x[0]), C.int64_t(n), nil, &config, (*C.float)(&distances[0]), (*C.int64_t)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("se_index_search error, code: %d", err)
	}
	return distances, ids, nil
}

func (index *SearchEngineIndex) SearchWithTimespace(n int64, x []float32, k int32, filter []*dsc.TimespaceFilter) ([]float32, []int64, error) {
	var filterPtr *C.se_id_filter_t
	filterConverted := make([]C.se_id_filter_t, len(filter))

	var ignoreFilters bool = true
	for i, f := range filter {
		if f == nil {
			return nil, nil, fmt.Errorf("search engine index search with batch filter fail, some filter is nil")
		}
		var ignoreFilter bool
		filterConverted[i], ignoreFilter = convertFilterByTimespaceFilter(f)
		if !ignoreFilter {
			ignoreFilters = false
		}
	}

	if len(filterConverted) != 0 && !ignoreFilters {
		filterPtr = (*C.se_id_filter_t)(unsafe.Pointer((&filterConverted[0])))
	}

	if k > 1024 {
		k = 1024
	}

	config := C.se_search_config_t{}
	config.k = C.int64_t(k)
	config.threshold = C.float(0)
	config.batch_size = C.int32_t(n)
	config.dist_type = C.DISTANCE_L2

	distances := make([]float32, n*int64(k))
	ids := make([]int64, n*int64(k))
	if err := int(C.se_index_search(index.handle, (*C.float)(&x[0]), C.int64_t(n), filterPtr, &config, (*C.float)(&distances[0]), (*C.int64_t)(&ids[0]))); err != 0 {
		return nil, nil, fmt.Errorf("search engine index search fail, error code: %d", err)
	}
	return distances, ids, nil
}

// GetIndexInfo
func (index *SearchEngineIndex) GetIndexInfo() (*dsc.IndexInfo, error) {
	indexInfo := &dsc.IndexInfo{
		SDKType: dsc.SE,
	}

	// if SearchEngineIndex is a onlyMem index, it can not call se_index_status, just return the base info
	if index.mem != nil && index.mem.data != nil && index.isOnlyMem == true {
		indexInfo.IsCPUModel = true
		// a onlyMem index is not a norm index, so set the indextype to Unknown
		indexInfo.IndexType = dsc.IndexUnknown
	} else {
		// if SearchEngineIndex is a normal index, return info get from se_index_status
		s := C.se_index_status_t{}
		if err := C.se_index_status(index.handle, &s); err != 0 {
			return nil, fmt.Errorf("se_index_status error, code: %d", err)
		}
		indexInfo.IndexType = dsc.IndexDC
		indexInfo.Size = uint64(s.size)
		indexInfo.IsTrained = (int(s.is_trained) == 1)
	}

	return indexInfo, nil
}

func (index *SearchEngineIndex) GetIndexIds() ([]int64, error) {
	info, err := index.GetIndexInfo()
	if err != nil {
		return nil, err
	}
	if info.Size == 0 {
		return nil, nil
	}
	ids := make([]int64, info.Size)
	if err := int(C.se_index_get_added_ids(index.handle, (*C.long)(&ids[0]))); err != 0 {
		return nil, fmt.Errorf("se_index_get_added_ids error, code: %d", err)
	}
	return ids, nil
}

func (index *SearchEngineIndex) TrainIndex(n int64, fs []float32) error {
	if err := int(C.se_index_train(index.handle, (*C.float)(&fs[0]), C.long(n))); err != 0 {
		return fmt.Errorf("se_index_train error, code: %d", err)
	}
	return nil
}

func (index *SearchEngineIndex) FreeIndex() error {
	if index.mem != nil && index.mem.data != nil {
		C.se_index_memory_data_free(index.mem)
	}
	if err := int(C.se_index_destroy(index.handle)); err != 0 {
		return fmt.Errorf("se_index_destroy error, code: %d", err)
	}
	return nil
}

func (index *SearchEngineIndex) WriteToCpuIndex(reclaimMemory bool) (dsc.SearchIndex, error) {
	cpuHandler := &SearchEngineIndex{
		mem:       &C.se_mem_data_t{},
		isOnlyMem: true,
	}
	if err := int(C.se_index_serialize_to_memory(index.handle, cpuHandler.mem)); err != 0 {
		return nil, fmt.Errorf("se_index_serialize_to_memory error, code: %d", err)
	}
	return cpuHandler, nil
}

func (index *SearchEngineIndex) WriteToDeviceIndex(resource interface{}) (dsc.SearchIndex, error) {
	return nil, nil
}

func (index *SearchEngineIndex) RemoveIndexIds(ids []int64) error {

	n := len(ids)
	if n == 0 {
		return nil
	}

	if err := int(C.se_index_remove_features(index.handle, (*C.int64_t)(&ids[0]), C.int64_t(n))); err != 0 {
		return fmt.Errorf("se_index_remove_features error, code: %d", err)
	}
	return nil
}

func (index *SearchEngineIndex) WriteToFile(filepath string) error {
	if index.mem == nil || index.mem.data == nil || index.isOnlyMem == false {
		return fmt.Errorf("SearchEngineIndex WriteToFile fail, no data to write, isOnlyMem: %v", index.isOnlyMem)
	}
	cfilepath := C.CString(filepath)
	// #nosec
	defer C.free(unsafe.Pointer(cfilepath))
	if err := int(C.se_index_memory_data_to_file(cfilepath, index.mem)); err != 0 {
		return fmt.Errorf("se_index_memory_data_to_file error, code: %d", err)
	}

	return nil
}

func SEIndexTypeToDSCIndexType(t int) dsc.IndexType {
	switch C.se_index_type_e(t) {
	case C.SE_INDEX_DC:
		return dsc.IndexDC
	case C.SE_INDEX_PQ:
		return dsc.IndexPQ
	case C.SE_INDEX_FLAT:
		return dsc.IndexFlat
	default:
		return dsc.IndexUnknown
	}
}
