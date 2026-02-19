package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
	"unsafe"
)

// ImportCPUMemory 将现有的 CPU 内存包装为 Memory 对象
//
// 这允许直接使用现有的 Go 切片或其他 CPU 内存，避免数据复制。
// 注意：原始内存的生命周期必须超过返回的 Memory 对象。
//
// 参数：
//   - runtime: Taichi 运行时
//   - ptr: CPU 内存指针
//   - size: 内存大小（字节）
//
// 返回：
//   - *Memory: 包装后的内存对象
//   - error: 如果导入失败
//
// 示例：
//
//	data := make([]float32, 1000)
//	memory, err := taichi.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
//	if err != nil {
//	    panic(err)
//	}
//	defer memory.Release() // 注意：这不会释放原始 Go 切片
//
//	// 现在可以直接使用这个内存创建 NdArray
//	arr, _ := taichi.NewNdArray1DFromMemory(memory, 1000, taichi.DATA_TYPE_F32)
func ImportCPUMemory(runtime *Runtime, ptr unsafe.Pointer, size uint64) (*Memory, error) {
	handle := c_api.ImportCPUMemory(runtime.handle, ptr, size)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("导入CPU内存失败 [%d]: %s", errCode, errMsg)
	}

	return &Memory{
		runtime: runtime,
		handle:  handle,
		size:    size,
		mapped:  false,
	}, nil
}

// ImportCUDAMemory 将现有的 CUDA 内存包装为 Memory 对象
//
// 这允许直接使用现有的 CUDA 内存，避免数据复制。
// 注意：原始内存的生命周期必须超过返回的 Memory 对象。
//
// 参数：
//   - runtime: Taichi 运行时（必须是 CUDA 后端）
//   - ptr: CUDA 内存指针（通过 cudaMalloc 等分配）
//   - size: 内存大小（字节）
//
// 返回：
//   - *Memory: 包装后的内存对象
//   - error: 如果导入失败
//
// 示例：
//
//	// 假设已有 CUDA 内存指针 cudaPtr (通过 CGO 或其他方式获得)
//	memory, err := taichi.ImportCUDAMemory(runtime, cudaPtr, 4000)
//	if err != nil {
//	    panic(err)
//	}
//	defer memory.Release() // 注意：这不会释放原始 CUDA 内存
func ImportCUDAMemory(runtime *Runtime, ptr unsafe.Pointer, size uint64) (*Memory, error) {
	handle := c_api.ImportCUDAMemory(runtime.handle, ptr, size)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("导入CUDA内存失败 [%d]: %s", errCode, errMsg)
	}

	return &Memory{
		runtime: runtime,
		handle:  handle,
		size:    size,
		mapped:  false,
	}, nil
}

// GetCUDAStream 获取当前的 CUDA 流
//
// 返回当前 Taichi 使用的 CUDA 流指针。
// 只在 CUDA 后端可用。
//
// 返回：
//   - unsafe.Pointer: CUDA 流指针，如果不支持返回 nil
//
// 示例：
//
//	stream := taichi.GetCUDAStream()
//	if stream != nil {
//	    // 可以将此流传递给其他 CUDA 代码使用
//	    fmt.Printf("当前 CUDA 流: %p\n", stream)
//	}
func GetCUDAStream() unsafe.Pointer {
	var stream unsafe.Pointer
	c_api.GetCUDAStream(&stream)
	return stream
}

// SetCUDAStream 设置 Taichi 使用的 CUDA 流
//
// 允许 Taichi 与现有的 CUDA 代码共享同一个流，
// 实现更好的同步和性能。
// 只在 CUDA 后端可用。
//
// 参数：
//   - stream: CUDA 流指针
//
// 示例：
//
//	// 假设已有 CUDA 流 myStream (通过 CGO 或其他方式获得)
//	taichi.SetCUDAStream(myStream)
//
//	// 现在 Taichi 的所有操作都会使用这个流
//	kernel.Launch().ArgNdArray(arr).Run()
func SetCUDAStream(stream unsafe.Pointer) {
	c_api.SetCUDAStream(stream)
}

// NewNdArray1DFromMemory 从现有 Memory 创建 1D NdArray
//
// 这是一个便捷函数，用于从导入的内存创建 NdArray。
//
// 参数：
//   - memory: 内存对象（可以是导入的或分配的）
//   - length: 数组长度
//   - elemType: 元素类型
//
// 返回：
//   - *NdArray: 创建的数组对象
//   - error: 如果创建失败
func NewNdArray1DFromMemory(memory *Memory, length uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{length},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}

// NewNdArray2DFromMemory 从现有 Memory 创建 2D NdArray
func NewNdArray2DFromMemory(memory *Memory, rows, cols uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{rows, cols},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}

// NewNdArray3DFromMemory 从现有 Memory 创建 3D NdArray
func NewNdArray3DFromMemory(memory *Memory, dim0, dim1, dim2 uint32, elemType c_api.TiDataType) (*NdArray, error) {
	return &NdArray{
		Memory:   memory,
		shape:    []uint32{dim0, dim1, dim2},
		elemType: elemType,
		elemSize: getElemSize(elemType),
	}, nil
}
