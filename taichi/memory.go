package taichi

import (
	"fmt"
	"unsafe"

	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Memory base class for memory
type Memory struct {
	runtime *Runtime
	handle  c_api.TiMemory
	size    uint64
}

// NewMemory creates a new memory object
func NewMemory(tiRuntime *Runtime, size uint64) (*Memory, error) {
	allocInfo := c_api.TiMemoryAllocateInfo{
		Size:      size,
		HostWrite: c_api.TI_TRUE,
		HostRead:  c_api.TI_TRUE,
		Usage:     c_api.TiMemoryUsageFlags(c_api.TI_MEMORY_USAGE_STORAGE_BIT),
	}

	handle := c_api.AllocateMemory(tiRuntime.handle, &allocInfo)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("memory allocation failed [%d]: %s", errCode, errMsg)
	}

	return &Memory{
		runtime: tiRuntime,
		handle:  handle,
		size:    size,
	}, nil
}

// Release releases memory
func (m *Memory) Release() {
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.FreeMemory(m.runtime.handle, m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
}

// mapMemory maps memory to host (internal use only)
func (m *Memory) mapMemory() (unsafe.Pointer, error) {
	ptr := c_api.MapMemory(m.runtime.handle, m.handle)
	if ptr == nil {
		return nil, fmt.Errorf("memory mapping failed")
	}
	return ptr, nil
}

// unmapMemory unmaps memory (internal use only)
func (m *Memory) unmapMemory() {
	c_api.UnmapMemory(m.runtime.handle, m.handle)
}

// Size gets memory size
func (m *Memory) Size() uint64 {
	return m.size
}

// Handle gets the underlying handle (for internal use or testing)
func (m *Memory) Handle() c_api.TiMemory {
	return m.handle
}

// CopyTo copies to another memory (device-side)
func (m *Memory) CopyTo(dst *Memory) error {
	if m.size != dst.size {
		return fmt.Errorf("memory size mismatch: %d vs %d", m.size, dst.size)
	}

	srcSlice := c_api.NewMemorySlice(m.handle, 0, m.size)
	dstSlice := c_api.NewMemorySlice(dst.handle, 0, dst.size)
	c_api.CopyMemoryDeviceToDevice(m.runtime.handle, &dstSlice, &srcSlice)

	return nil
}

// CopyFrom copies from another memory (device-side)
func (m *Memory) CopyFrom(src *Memory) error {
	return src.CopyTo(m)
}

// MapMemory maps memory, executes the write/read function, and unmaps.
// This ensures thread safety by using syncCall which serializes all C-API calls.
func MapMemory(fn func(ptrs ...unsafe.Pointer) error, memories ...*Memory) error {
	var mappedMemories []*Memory
	defer func() {
		for _, ptr := range mappedMemories {
			ptr.unmapMemory()
		}
	}()

	var ptrs []unsafe.Pointer
	for _, memory := range memories {
		ptr := c_api.MapMemory(memory.runtime.handle, memory.handle)
		if ptr == nil {
			return fmt.Errorf("memory mapping failed")
		}
		mappedMemories = append(mappedMemories, memory)
		ptrs = append(ptrs, ptr)
	}

	return c_api.SyncCall(func() error {
		// Execute user's write function
		return fn(ptrs...)
	})
}

// Read reads data from memory into the provided slice
func (m *Memory) Read(data []byte) error {
	ptr := c_api.MapMemory(m.runtime.handle, m.handle)
	if ptr == nil {
		return fmt.Errorf("memory mapping failed")
	}
	defer c_api.UnmapMemory(m.runtime.handle, m.handle)

	c_api.SyncCallVoid(func() {
		copy(data, unsafe.Slice((*byte)(ptr), m.size))
	})
	return nil
}
