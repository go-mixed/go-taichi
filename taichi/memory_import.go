package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
	"unsafe"
)

// ImportCPUMemory wraps existing CPU memory as a Memory object
//
// This allows direct use of existing Go slices or other CPU memory, avoiding data copying.
// Note: The lifetime of the original memory must exceed that of the returned Memory object.
//
// Parameters:
//   - runtime: Taichi runtime
//   - ptr: CPU memory pointer
//   - size: Memory size (bytes)
//
// Returns:
//   - *Memory: Wrapped memory object
//   - error: If import fails
//
// Example:
//
//	data := make([]float32, 1000)
//	memory, err := taichi.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
//	if err != nil {
//	    panic(err)
//	}
//	defer memory.Release() // Note: This does not free the original Go slice
//
//	// Now you can directly use this memory to create NdArray
//	arr, _ := taichi.NewNdArray1DFromMemory(memory, 1000, taichi.DATA_TYPE_F32)
func ImportCPUMemory(runtime *Runtime, ptr unsafe.Pointer, size uint64) (*Memory, error) {
	handle := c_api.ImportCPUMemory(runtime.handle, ptr, size)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to import CPU memory [%d]: %s", errCode, errMsg)
	}

	return &Memory{
		runtime: runtime,
		handle:  handle,
		size:    size,
		mapped:  false,
	}, nil
}

// ImportCUDAMemory wraps existing CUDA memory as a Memory object
//
// This allows direct use of existing CUDA memory, avoiding data copying.
// Note: The lifetime of the original memory must exceed that of the returned Memory object.
//
// Parameters:
//   - runtime: Taichi runtime (must be CUDA backend)
//   - ptr: CUDA memory pointer (allocated via cudaMalloc, etc.)
//   - size: Memory size (bytes)
//
// Returns:
//   - *Memory: Wrapped memory object
//   - error: If import fails
//
// Example:
//
//	// Assuming existing CUDA memory pointer cudaPtr (obtained via CGO or other means)
//	memory, err := taichi.ImportCUDAMemory(runtime, cudaPtr, 4000)
//	if err != nil {
//	    panic(err)
//	}
//	defer memory.Release() // Note: This does not free the original CUDA memory
func ImportCUDAMemory(runtime *Runtime, ptr unsafe.Pointer, size uint64) (*Memory, error) {
	handle := c_api.ImportCUDAMemory(runtime.handle, ptr, size)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to import CUDA memory [%d]: %s", errCode, errMsg)
	}

	return &Memory{
		runtime: runtime,
		handle:  handle,
		size:    size,
		mapped:  false,
	}, nil
}

// GetCUDAStream gets the current CUDA stream
//
// Returns the CUDA stream pointer currently used by Taichi.
// Only available on CUDA backend.
//
// Returns:
//   - unsafe.Pointer: CUDA stream pointer, or nil if not supported
//
// Example:
//
//	stream := taichi.GetCUDAStream()
//	if stream != nil {
//	    // Can pass this stream to other CUDA code
//	    fmt.Printf("Current CUDA stream: %p\n", stream)
//	}
func GetCUDAStream() unsafe.Pointer {
	var stream unsafe.Pointer
	c_api.GetCUDAStream(&stream)
	return stream
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
//	// Assuming existing CUDA stream myStream (obtained via CGO or other means)
//	taichi.SetCUDAStream(myStream)
//
//	// Now all Taichi operations will use this stream
//	kernel.Launch().ArgNdArray(arr).Run()
func SetCUDAStream(stream unsafe.Pointer) {
	c_api.SetCUDAStream(stream)
}

// NewNdArray1DFromMemory creates a 1D NdArray from existing Memory
//
// This is a convenience function for creating NdArray from imported memory.
//
// Parameters:
//   - memory: Memory object (can be imported or allocated)
//   - length: Array length
//   - elemType: Element type
//
// Returns:
//   - *NdArray: Created array object
//   - error: If creation fails
func NewNdArray1DFromMemory(memory *Memory, length uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{length},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}

// NewNdArray2DFromMemory creates a 2D NdArray from existing Memory
func NewNdArray2DFromMemory(memory *Memory, rows, cols uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{rows, cols},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}

// NewNdArray3DFromMemory creates a 3D NdArray from existing Memory
func NewNdArray3DFromMemory(memory *Memory, dim0, dim1, dim2 uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{dim0, dim1, dim2},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}
