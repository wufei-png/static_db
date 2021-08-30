package dsc

import (
	"fmt"
	"strings"
	"time"
)

type TimespaceFilter struct {
	TimeRange  [2]uint32
	CameraMask [128 / 8]uint8
}

type IndexType int

const (
	IndexUnknown IndexType = -1
	IndexFlat    IndexType = 0
	IndexPQ      IndexType = 1
	IndexIVFPQ   IndexType = 2
	IndexDC      IndexType = 3
	IndexHNSW    IndexType = 4
)

type GeneralIndexConfig struct {
	Dimension     int32
	Nlist         int32
	Nprobe        int32
	SubQuantizers int32
	BitsPerCode   int32
	MaxSize       int64
	IndexType     IndexType
	ModelPath     string
	InitFilePath  string
	M             int32
	EfConstruct   int32
	Ef            int32
	UseInt8       int32
}

type IndexInfo struct {
	Size          uint64
	Dimension     int32
	IsTrained     bool
	Nlist         int32
	Nprobe        int32
	MaxListSize   int32
	SubQuantizers int32
	BitsPerCode   int32
	SDKType       string
	IndexType     IndexType
	IsCPUModel    bool
}

type SearchIndex interface {
	AddBatch(n int64, x []float32, ids []int64) error
	Search(n int64, x []float32, k int32) ([]float32, []int64, error)
	SearchWithTimespace(n int64, x []float32, k int32, filter []*TimespaceFilter) ([]float32, []int64, error)
	GetIndexInfo() (*IndexInfo, error)
	GetIndexIds() ([]int64, error)
	TrainIndex(n int64, fs []float32) error
	FreeIndex() error
	WriteToCpuIndex(reclaimMemory bool) (SearchIndex, error)
	WriteToDeviceIndex(resource interface{}) (SearchIndex, error)
	WriteToFile(filepath string) error
	RemoveIndexIds([]int64) error
}

func (f *TimespaceFilter) SetTimeRange(start, end time.Time) {
	f.TimeRange[0] = uint32(start.Unix())
	f.TimeRange[1] = uint32(end.Unix())
}

func (f *TimespaceFilter) SetCameraMask(idx uint) {
	if idx > 127 {
		panic("index out of range")
	}
	f.CameraMask[idx/8] |= 1 << (idx % 8)
}

func ParseIndexType(s string) (IndexType, error) {
	var t IndexType
	if err := t.UnmarshalText([]byte(strings.ToUpper(s))); err != nil {
		return IndexUnknown, err
	}
	return t, nil
}

func (t IndexType) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t *IndexType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "FLAT":
		*t = IndexFlat
	case "PQ":
		*t = IndexPQ
	case "IVFPQ":
		*t = IndexIVFPQ
	case "DC":
		*t = IndexDC
	case "HNSW":
		*t = IndexHNSW
	default:
		return fmt.Errorf("invalid index type %q", text)
	}
	return nil
}

func (t IndexType) String() string {
	switch t {
	case IndexFlat:
		return "FLAT"
	case IndexPQ:
		return "PQ"
	case IndexIVFPQ:
		return "IVFPQ"
	case IndexDC:
		return "DC"
	case IndexHNSW:
		return "HNSW"
	default:
		return fmt.Sprintf("unknown type %v", int(t))
	}
}
