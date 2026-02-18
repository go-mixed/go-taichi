package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// 内存导入相关函数
// 这些函数允许将外部内存（CPU、CUDA等）导入为 TiMemory，避免数据复制

var (
	// CPU 内存导入
	tiImportCpuMemory func(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory

	// CUDA 内存导入
	tiImportCudaMemory func(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory

	// CUDA 流管理
	tiGetCudaStream func(stream *unsafe.Pointer)
	tiSetCudaStream func(stream unsafe.Pointer)
)

// registerMemoryImportFunctions 注册内存导入相关函数
func registerMemoryImportFunctions(libHandle uintptr) error {
	// 导入 purego 包
	// 注意：这些函数可能不被所有后端支持，失败时设为 nil

	// CPU 内存导入 - 可能不被所有后端支持
	if err := tryRegisterFunction(&tiImportCpuMemory, libHandle, "ti_import_cpu_memory"); err != nil {
		tiImportCpuMemory = nil
	}

	// CUDA 内存导入 - 只在 CUDA 后端可用
	if err := tryRegisterFunction(&tiImportCudaMemory, libHandle, "ti_import_cuda_memory"); err != nil {
		tiImportCudaMemory = nil
	}

	// CUDA 流管理 - 只在 CUDA 后端可用
	if err := tryRegisterFunction(&tiGetCudaStream, libHandle, "ti_get_cuda_stream"); err != nil {
		tiGetCudaStream = nil
	}
	if err := tryRegisterFunction(&tiSetCudaStream, libHandle, "ti_set_cuda_stream"); err != nil {
		tiSetCudaStream = nil
	}

	return nil
}

// tryRegisterFunction 尝试注册函数，失败时不返回错误
func tryRegisterFunction(fn interface{}, libHandle uintptr, name string) error {
	// 使用 purego 注册函数，如果函数不存在会 panic
	// 我们捕获 panic 并转换为错误
	defer func() {
		if r := recover(); r != nil {
			// 函数不存在，这是正常的（某些后端不支持某些功能）
		}
	}()

	purego.RegisterLibFunc(fn, libHandle, name)
	return nil
}

// ImportCPUMemory 将现有的 CPU 内存指针包装为 TiMemory
//
// 这允许直接使用现有的 CPU 内存，避免数据复制。
// 注意：原始内存的生命周期必须超过返回的 TiMemory。
//
// 参数：
//   - runtime: Taichi 运行时
//   - ptr: CPU 内存指针
//   - size: 内存大小（字节）
//
// 返回：
//   - TiMemory: 包装后的内存句柄，如果失败返回 TI_NULL_HANDLE
//
// 示例：
//
//	data := make([]float32, 1000)
//	memory := c_api.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
//	if memory == c_api.TI_NULL_HANDLE {
//	    // 处理错误
//	}
//	defer c_api.FreeMemory(runtime, memory) // 注意：这不会释放原始内存
func ImportCPUMemory(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory {
	if tiImportCpuMemory == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CPU memory import not supported by current backend")
		return TI_NULL_HANDLE
	}
	return tiImportCpuMemory(runtime, ptr, size)
}

// ImportCUDAMemory 将现有的 CUDA 内存指针包装为 TiMemory
//
// 这允许直接使用现有的 CUDA 内存，避免数据复制。
// 注意：原始内存的生命周期必须超过返回的 TiMemory。
//
// 参数：
//   - runtime: Taichi 运行时（必须是 CUDA 后端）
//   - ptr: CUDA 内存指针（通过 cudaMalloc 等分配）
//   - size: 内存大小（字节）
//
// 返回：
//   - TiMemory: 包装后的内存句柄，如果失败返回 TI_NULL_HANDLE
//
// 示例：
//
//	// 假设已有 CUDA 内存指针 cudaPtr
//	memory := c_api.ImportCUDAMemory(runtime, cudaPtr, 4000)
//	if memory == c_api.TI_NULL_HANDLE {
//	    // 处理错误
//	}
//	defer c_api.FreeMemory(runtime, memory) // 注意：这不会释放原始 CUDA 内存
func ImportCUDAMemory(runtime TiRuntime, ptr unsafe.Pointer, size uint64) TiMemory {
	if tiImportCudaMemory == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA memory import not supported by current backend")
		return TI_NULL_HANDLE
	}
	return tiImportCudaMemory(runtime, ptr, size)
}

// GetCUDAStream 获取当前的 CUDA 流
//
// 返回当前 Taichi 使用的 CUDA 流指针。
// 只在 CUDA 后端可用。
//
// 参数：
//   - stream: 用于接收流指针的指针
//
// 示例：
//
//	var stream unsafe.Pointer
//	c_api.GetCUDAStream(&stream)
//	if stream != nil {
//	    // 使用 CUDA 流
//	}
func GetCUDAStream(stream *unsafe.Pointer) {
	if tiGetCudaStream == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA stream management not supported by current backend")
		*stream = nil
		return
	}
	tiGetCudaStream(stream)
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
//	// 假设已有 CUDA 流 myStream
//	c_api.SetCUDAStream(myStream)
func SetCUDAStream(stream unsafe.Pointer) {
	if tiSetCudaStream == nil {
		SetLastError(TI_ERROR_NOT_SUPPORTED, "CUDA stream management not supported by current backend")
		return
	}
	tiSetCudaStream(stream)
}
