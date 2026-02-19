package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
	"unsafe"
)

// Memory base class for memory
type Memory struct {
	runtime *Runtime
	handle  c_api.TiMemory
	size    uint64
	mapped  bool
	ptr     unsafe.Pointer
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

	m := &Memory{
		runtime: tiRuntime,
		handle:  handle,
		size:    size,
		mapped:  false,
		ptr:     nil,
	}

	return m, nil
}

// Release releases memory
func (m *Memory) Release() {
	if m.mapped {
		m.Unmap()
	}
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.FreeMemory(m.runtime.handle, m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
}

// Map maps memory to host
func (m *Memory) Map() (unsafe.Pointer, error) {
	if m.mapped {
		return m.ptr, nil
	}

	ptr := c_api.MapMemory(m.runtime.handle, m.handle)
	if ptr == nil {
		return nil, fmt.Errorf("memory mapping failed")
	}

	m.ptr = ptr
	m.mapped = true
	return ptr, nil
}

// Unmap unmaps memory
func (m *Memory) Unmap() {
	if m.mapped {
		c_api.UnmapMemory(m.runtime.handle, m.handle)
		m.mapped = false
		m.ptr = nil
	}
}

// Size gets memory size
func (m *Memory) Size() uint64 {
	return m.size
}

// IsMapped checks if memory is mapped
func (m *Memory) IsMapped() bool {
	return m.mapped
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
