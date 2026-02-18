package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// ===== 内存管理函数指针 =====

var (
	tiAllocateMemory           func(runtime TiRuntime, allocateInfo *TiMemoryAllocateInfo) TiMemory
	tiFreeMemory               func(runtime TiRuntime, memory TiMemory)
	tiMapMemory                func(runtime TiRuntime, memory TiMemory) unsafe.Pointer
	tiUnmapMemory              func(runtime TiRuntime, memory TiMemory)
	tiCopyMemoryDeviceToDevice func(runtime TiRuntime, dstMemory *TiMemorySlice, srcMemory *TiMemorySlice)
)

// registerMemoryFunctions 注册内存管理函数
func registerMemoryFunctions() error {
	purego.RegisterLibFunc(&tiAllocateMemory, libHandle, "ti_allocate_memory")
	purego.RegisterLibFunc(&tiFreeMemory, libHandle, "ti_free_memory")
	purego.RegisterLibFunc(&tiMapMemory, libHandle, "ti_map_memory")
	purego.RegisterLibFunc(&tiUnmapMemory, libHandle, "ti_unmap_memory")
	purego.RegisterLibFunc(&tiCopyMemoryDeviceToDevice, libHandle, "ti_copy_memory_device_to_device")
	return nil
}

// ===== 导出的内存管理函数 =====

// AllocateMemory 使用提供的参数分配连续的设备内存
//
// 参数:
//   - runtime: 运行时句柄
//   - allocateInfo: 内存分配信息
//
// 返回:
//   - 内存句柄,如果分配失败则返回TI_NULL_HANDLE
//
// 注意:
//   - 分配的内存在相关运行时销毁时自动释放
//   - 也可以手动调用FreeMemory释放内存
//
// 示例:
//
//	allocInfo := taichi.TiMemoryAllocateInfo{
//	    Size:      1024,
//	    HostWrite: taichi.TI_TRUE,
//	    HostRead:  taichi.TI_TRUE,
//	    Usage:     taichi.TiMemoryUsageFlags(taichi.TI_MEMORY_USAGE_STORAGE_BIT),
//	}
//	memory := taichi.AllocateMemory(runtime, &allocInfo)
//	defer taichi.FreeMemory(runtime, memory)
func AllocateMemory(runtime TiRuntime, allocateInfo *TiMemoryAllocateInfo) TiMemory {
	return tiAllocateMemory(runtime, allocateInfo)
}

// FreeMemory 释放内存分配
//
// 参数:
//   - runtime: 运行时句柄
//   - memory: 要释放的内存句柄
//
// 示例:
//
//	taichi.FreeMemory(runtime, memory)
func FreeMemory(runtime TiRuntime, memory TiMemory) {
	tiFreeMemory(runtime, memory)
}

// MapMemory 将设备内存映射到主机可寻址空间
//
// 在映射之前,必须确保设备没有被任何设备命令使用。
//
// 参数:
//   - runtime: 运行时句柄
//   - memory: 要映射的内存句柄
//
// 返回:
//   - 映射的主机地址,如果映射失败则返回nil
//
// 示例:
//
//	ptr := taichi.MapMemory(runtime, memory)
//	if ptr != nil {
//	    data := (*[256]uint32)(ptr)
//	    for i := 0; i < 256; i++ {
//	        data[i] = uint32(i)
//	    }
//	    taichi.UnmapMemory(runtime, memory)
//	}
func MapMemory(runtime TiRuntime, memory TiMemory) unsafe.Pointer {
	return tiMapMemory(runtime, memory)
}

// UnmapMemory 取消设备内存映射并使主机端的任何更改对设备可见
//
// 必须确保不再访问先前映射的主机可寻址空间。
//
// 参数:
//   - runtime: 运行时句柄
//   - memory: 要取消映射的内存句柄
//
// 示例:
//
//	ptr := taichi.MapMemory(runtime, memory)
//	// ... 操作ptr ...
//	taichi.UnmapMemory(runtime, memory)
func UnmapMemory(runtime TiRuntime, memory TiMemory) {
	tiUnmapMemory(runtime, memory)
}

// CopyMemoryDeviceToDevice 在设备内复制内存的连续子部分
//
// 两个子部分不能重叠。这是一个设备命令。
//
// 参数:
//   - runtime: 运行时句柄
//   - dst: 目标内存切片
//   - src: 源内存切片
//
// 示例:
//
//	srcSlice := &taichi.TiMemorySlice{
//	    Memory: srcMemory,
//	    Offset: 0,
//	    Size:   1024,
//	}
//	dstSlice := &taichi.TiMemorySlice{
//	    Memory: dstMemory,
//	    Offset: 0,
//	    Size:   1024,
//	}
//	taichi.CopyMemoryDeviceToDevice(runtime, dstSlice, srcSlice)
func CopyMemoryDeviceToDevice(runtime TiRuntime, dst *TiMemorySlice, src *TiMemorySlice) {
	tiCopyMemoryDeviceToDevice(runtime, dst, src)
}
