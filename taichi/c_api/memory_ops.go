package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// ===== Memory Management Function Pointers =====

var (
	tiAllocateMemory           func(runtime TiRuntime, allocateInfo *TiMemoryAllocateInfo) TiMemory
	tiFreeMemory               func(runtime TiRuntime, memory TiMemory)
	tiMapMemory                func(runtime TiRuntime, memory TiMemory) unsafe.Pointer
	tiUnmapMemory              func(runtime TiRuntime, memory TiMemory)
	tiCopyMemoryDeviceToDevice func(runtime TiRuntime, dstMemory *TiMemorySlice, srcMemory *TiMemorySlice)
)

// registerMemoryFunctions registers memory management functions
func registerMemoryFunctions() error {
	purego.RegisterLibFunc(&tiAllocateMemory, libHandle, "ti_allocate_memory")
	purego.RegisterLibFunc(&tiFreeMemory, libHandle, "ti_free_memory")
	purego.RegisterLibFunc(&tiMapMemory, libHandle, "ti_map_memory")
	purego.RegisterLibFunc(&tiUnmapMemory, libHandle, "ti_unmap_memory")
	purego.RegisterLibFunc(&tiCopyMemoryDeviceToDevice, libHandle, "ti_copy_memory_device_to_device")
	return nil
}

// ===== Exported Memory Management Functions =====

// AllocateMemory allocates contiguous device memory with the provided parameters
//
// Parameters:
//   - runtime: Runtime handle
//   - allocateInfo: Memory allocation information
//
// Returns:
//   - Memory handle, or TI_NULL_HANDLE if allocation fails
//
// Note:
//   - Allocated memory is automatically freed when the associated runtime is destroyed
//   - Can also be manually freed by calling FreeMemory
//
// Example:
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
	return SyncCall(func() TiMemory {
		return tiAllocateMemory(runtime, allocateInfo)
	})
}

// FreeMemory frees a memory allocation
//
// Parameters:
//   - runtime: Runtime handle
//   - memory: Memory handle to free
//
// Example:
//
//	taichi.FreeMemory(runtime, memory)
func FreeMemory(runtime TiRuntime, memory TiMemory) {
	SyncCallVoid(func() {
		tiFreeMemory(runtime, memory)
		asyncTasks.Add(1)
	})
}

// MapMemory maps device memory to host-addressable space
//
// Before mapping, ensure the device is not being used by any device commands.
//
// Parameters:
//   - runtime: Runtime handle
//   - memory: Memory handle to map
//
// Returns:
//   - Mapped host address, or nil if mapping fails
//
// Example:
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
	return SyncCall(func() unsafe.Pointer {
		return tiMapMemory(runtime, memory)
	})
}

// UnmapMemory unmaps device memory and makes any host-side changes visible to the device
//
// Must ensure the previously mapped host-addressable space is no longer accessed.
//
// Parameters:
//   - runtime: Runtime handle
//   - memory: Memory handle to unmap
//
// Example:
//
//	ptr := taichi.MapMemory(runtime, memory)
//	// ... operate on ptr ...
//	taichi.UnmapMemory(runtime, memory)
func UnmapMemory(runtime TiRuntime, memory TiMemory) {
	SyncCallVoid(func() {
		tiUnmapMemory(runtime, memory)
		asyncTasks.Add(1)
	})
}

// CopyMemoryDeviceToDevice copies a contiguous subsection of memory within the device
//
// The two subsections must not overlap. This is a device command.
//
// Parameters:
//   - runtime: Runtime handle
//   - dst: Destination memory slice
//   - src: Source memory slice
//
// Example:
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
	SyncCallVoid(func() {
		tiCopyMemoryDeviceToDevice(runtime, dst, src)
		asyncTasks.Add(1)
	})
}
