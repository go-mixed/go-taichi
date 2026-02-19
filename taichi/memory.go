package taichi

import (
	"fmt"
	"go-taichi/taichi/c_api"
	"unsafe"
)

// Memory 内存基类
type Memory struct {
	runtime *Runtime
	handle  c_api.TiMemory
	size    uint64
	mapped  bool
	ptr     unsafe.Pointer
}

// NewMemory 创建新的内存对象
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
		return nil, fmt.Errorf("内存分配失败 [%d]: %s", errCode, errMsg)
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

// Release 释放内存
func (m *Memory) Release() {
	if m.mapped {
		m.Unmap()
	}
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.FreeMemory(m.runtime.handle, m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
}

// Map 映射内存到主机
func (m *Memory) Map() (unsafe.Pointer, error) {
	if m.mapped {
		return m.ptr, nil
	}

	ptr := c_api.MapMemory(m.runtime.handle, m.handle)
	if ptr == nil {
		return nil, fmt.Errorf("内存映射失败")
	}

	m.ptr = ptr
	m.mapped = true
	return ptr, nil
}

// Unmap 取消内存映射
func (m *Memory) Unmap() {
	if m.mapped {
		c_api.UnmapMemory(m.runtime.handle, m.handle)
		m.mapped = false
		m.ptr = nil
	}
}

// Size 获取内存大小
func (m *Memory) Size() uint64 {
	return m.size
}

// IsMapped 检查是否已映射
func (m *Memory) IsMapped() bool {
	return m.mapped
}

// Handle 获取底层句柄（用于内部或测试）
func (m *Memory) Handle() c_api.TiMemory {
	return m.handle
}

// CopyTo 复制到另一个内存（设备端）
func (m *Memory) CopyTo(dst *Memory) error {
	if m.size != dst.size {
		return fmt.Errorf("内存大小不匹配: %d vs %d", m.size, dst.size)
	}

	srcSlice := c_api.NewMemorySlice(m.handle, 0, m.size)
	dstSlice := c_api.NewMemorySlice(dst.handle, 0, dst.size)
	c_api.CopyMemoryDeviceToDevice(m.runtime.handle, &dstSlice, &srcSlice)

	return nil
}

// CopyFrom 从另一个内存复制（设备端）
func (m *Memory) CopyFrom(src *Memory) error {
	return src.CopyTo(m)
}
