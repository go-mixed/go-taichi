package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// Memory import related functions
// These functions allow importing external memory (CPU, CUDA, etc.) as TiMemory, avoiding data copying

var (
	// CPU memory import
	tiImportCpuMemory func(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory

	// CUDA memory import
	tiImportCudaMemory func(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory

	// CUDA stream management
	tiGetCudaStream func(stream *unsafe.Pointer)
	tiSetCudaStream func(stream unsafe.Pointer)
)

// registerMemoryImportFunctions registers memory import related functions
func registerMemoryImportFunctions(libHandle uintptr) error {
	// Import purego package
	// Note: These functions may not be supported by all backends, set to nil on failure

	// CPU memory import - may not be supported by all backends
	if err := tryRegisterFunction(&tiImportCpuMemory, libHandle, "ti_import_cpu_memory"); err != nil {
		tiImportCpuMemory = nil
	}

	// CUDA memory import - only available on CUDA backend
	if err := tryRegisterFunction(&tiImportCudaMemory, libHandle, "ti_import_cuda_memory"); err != nil {
		tiImportCudaMemory = nil
	}

	// CUDA stream management - only available on CUDA backend
	if err := tryRegisterFunction(&tiGetCudaStream, libHandle, "ti_get_cuda_stream"); err != nil {
		tiGetCudaStream = nil
	}
	if err := tryRegisterFunction(&tiSetCudaStream, libHandle, "ti_set_cuda_stream"); err != nil {
		tiSetCudaStream = nil
	}

	return nil
}

// tryRegisterFunction attempts to register a function, does not return error on failure
func tryRegisterFunction(fn interface{}, libHandle uintptr, name string) error {
	// Use purego to register function, will panic if function doesn't exist
	// We catch the panic and convert it to an error
	defer func() {
		if r := recover(); r != nil {
			// Function doesn't exist, this is normal (some backends don't support certain features)
		}
	}()

	purego.RegisterLibFunc(fn, libHandle, name)
	return nil
}

// ImportCPUMemory wraps an existing CPU memory pointer as TiMemory
//
// This allows direct use of existing CPU memory, avoiding data copying.
// Note: The lifetime of the original memory must exceed that of the returned TiMemory.
//
// Parameters:
//   - runtime: Taichi runtime
//   - ptr: CPU memory pointer
//   - size: Memory size (bytes)
//
// Returns:
//   - TiMemory: Wrapped memory handle, or TI_NULL_HANDLE if failed
//
// Example:
//
//	data := make([]float32, 1000)
//	memory := c_api.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
//	if memory == c_api.TI_NULL_HANDLE {
//	    // Handle error
//	}
//	defer c_api.FreeMemory(runtime, memory) // Note: This does not free the original memory
func ImportCPUMemory(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory {
	if tiImportCpuMemory == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CPU memory import not supported by current backend")
		return TI_NULL_HANDLE
	}
	return SyncCall(func() TiMemory {
		return tiImportCpuMemory(runtime, ptr, size)
	})
}

// ImportCUDAMemory wraps an existing CUDA memory pointer as TiMemory
//
// This allows direct use of existing CUDA memory, avoiding data copying.
// Note: The lifetime of the original memory must exceed that of the returned TiMemory.
//
// Parameters:
//   - runtime: Taichi runtime (must be CUDA backend)
//   - ptr: CUDA memory pointer (allocated via cudaMalloc, etc.)
//   - size: Memory size (bytes)
//
// Returns:
//   - TiMemory: Wrapped memory handle, or TI_NULL_HANDLE if failed
//
// Example:
//
//	// Assuming existing CUDA memory pointer cudaPtr
//	memory := c_api.ImportCUDAMemory(runtime, cudaPtr, 4000)
//	if memory == c_api.TI_NULL_HANDLE {
//	    // Handle error
//	}
//	defer c_api.FreeMemory(runtime, memory) // Note: This does not free the original CUDA memory
func ImportCUDAMemory(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory {
	if tiImportCudaMemory == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA memory import not supported by current backend")
		return TI_NULL_HANDLE
	}
	return SyncCall(func() TiMemory {
		return tiImportCudaMemory(runtime, ptr, size)
	})
}

// GetCUDAStream gets the current CUDA stream
//
// Returns the CUDA stream pointer currently used by Taichi.
// Only available on CUDA backend.
//
// Parameters:
//   - stream: Pointer to receive the stream pointer
//
// Example:
//
//	var stream unsafe.Pointer
//	c_api.GetCUDAStream(&stream)
//	if stream != nil {
//	    // Use CUDA stream
//	}
func GetCUDAStream(stream *unsafe.Pointer) {
	if tiGetCudaStream == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA stream management not supported by current backend")
		*stream = nil
		return
	}
	SyncCallVoid(func() {
		tiGetCudaStream(stream)
	})
}

// SetCUDAStream sets the CUDA stream used by Taichi
//
// Allows Taichi to share the same stream with existing CUDA code,
// enabling better synchronization and performance.
// Only available on CUDA backend.
//
// Parameters:
//   - stream: CUDA stream pointer
//
// Example:
//
//	// Assuming existing CUDA stream myStream
//	c_api.SetCUDAStream(myStream)
func SetCUDAStream(stream unsafe.Pointer) {
	if tiSetCudaStream == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA stream management not supported by current backend")
		return
	}
	SyncCallVoid(func() {
		tiSetCudaStream(stream)
	})
}
