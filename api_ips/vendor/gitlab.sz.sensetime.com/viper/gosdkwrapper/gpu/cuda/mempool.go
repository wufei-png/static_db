package cuda

import (
	"sync"
	"unsafe"
)

type MemPool struct {
	mu       sync.Mutex
	closed   bool
	blobsize int
	free     []unsafe.Pointer
	alloced  map[unsafe.Pointer]bool
}

func NewMemPool(blobsize int, count int) (*MemPool, error) {
	m := make([]unsafe.Pointer, count)
	for i := 0; i < count; i++ {
		dm, err := AllocDeviceMemory(blobsize)
		if err != nil {
			return nil, err
		}
		m[i] = dm.UnsafePtr()
	}
	return &MemPool{
		blobsize: blobsize,
		free:     m,
		alloced:  make(map[unsafe.Pointer]bool, count),
	}, nil
}

func (m *MemPool) Alloc(size int) (DeviceMemory, error) {
	if size > m.blobsize {
		return AllocDeviceMemory(size)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || len(m.free) == 0 {
		return AllocDeviceMemory(size)
	}
	free := m.free[len(m.free)-1]
	m.free = m.free[:len(m.free)-1]
	m.alloced[free] = true
	return DeviceMemoryFromUnsafePointer(free), nil
}

func (m *MemPool) Free(dm DeviceMemory) error {
	if dm.IsNil() {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		if !dm.owned {
			dm = DeviceMemory{ptr: dm.ptr, owned: true}
		}
		return dm.Free()
	}
	if m.alloced[dm.UnsafePtr()] {
		m.free = append(m.free, dm.UnsafePtr())
		delete(m.alloced, dm.UnsafePtr())
	} else {
		dm.Free() // nolint
	}
	return nil
}

func (m *MemPool) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.free {
		dm := DeviceMemoryFromUnsafePointer(m.free[i])
		dm.owned = true
		dm.Free() // nolint
	}
	m.free = nil
	m.alloced = nil
	m.closed = true
}
