package main

import (
	"fmt"
	"unsafe"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：CUDA 内存导入和流管理
// 功能：演示 CUDA 内存导入 API 和流管理（概念演示）

func main() {
	fmt.Println("=== CUDA 内存导入和流管理 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// === CUDA 内存导入示例 ===
	fmt.Println("--- CUDA 内存导入 API ---")
	fmt.Println("⚠️  注意：CUDA 内存导入需要 CGO 和 CUDA SDK")
	fmt.Println()
	fmt.Println("示例代码（需要 CGO）：")
	fmt.Println("```go")
	fmt.Println("// #include <cuda_runtime.h>")
	fmt.Println("// import \"C\"")
	fmt.Println()
	fmt.Println("var cudaPtr unsafe.Pointer")
	fmt.Println("C.cudaMalloc(&cudaPtr, 4000)")
	fmt.Println()
	fmt.Println("memory, err := taichi.ImportCUDAMemory(runtime, cudaPtr, 4000)")
	fmt.Println("if err != nil {")
	fmt.Println("    panic(err)")
	fmt.Println("}")
	fmt.Println("defer memory.Release()")
	fmt.Println("defer C.cudaFree(cudaPtr)")
	fmt.Println("```")

	// 尝试导入无效指针（仅用于演示错误处理）
	fmt.Println("\n--- 错误处理演示 ---")
	_, err = taichi.ImportCUDAMemory(runtime, unsafe.Pointer(uintptr(0x12345678)), 1000)
	if err != nil {
		fmt.Printf("❌ CUDA 内存导入失败: %v\n", err)
		fmt.Println("   (这是预期的，因为传入了无效指针)")
	}

	// === CUDA 流管理 ===
	fmt.Println("\n--- CUDA 流管理 ---")

	// 获取当前 CUDA 流
	stream := taichi.GetCUDAStream()
	if stream != nil {
		fmt.Printf("✅ 当前 CUDA 流: %p\n", stream)
	} else {
		fmt.Println("⚠️  CUDA 流不可用")
		fmt.Println("   原因：非 CUDA 后端或不支持")
	}

	// 设置 CUDA 流示例
	fmt.Println("\n--- 设置自定义 CUDA 流 ---")
	fmt.Println("示例代码（需要 CGO）：")
	fmt.Println("```go")
	fmt.Println("// 创建 CUDA 流")
	fmt.Println("var myStream unsafe.Pointer")
	fmt.Println("C.cudaStreamCreate(&myStream)")
	fmt.Println()
	fmt.Println("// 设置为 Taichi 使用的流")
	fmt.Println("taichi.SetCUDAStream(myStream)")
	fmt.Println()
	fmt.Println("// 现在 Taichi 的所有操作都会使用这个流")
	fmt.Println("kernel.Launch().ArgNdArray(arr).Run()")
	fmt.Println()
	fmt.Println("// 清理")
	fmt.Println("C.cudaStreamDestroy(myStream)")
	fmt.Println("```")

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 CUDA 内存导入要点：")
	fmt.Println("   • 需要 CGO 和 CUDA SDK")
	fmt.Println("   • 与现有 CUDA 代码集成")
	fmt.Println("   • 避免 CPU-GPU 数据传输")
	fmt.Println("   • 导入的内存生命周期必须超过 Memory 对象")
	fmt.Println("\n💡 CUDA 流管理要点：")
	fmt.Println("   • 控制 Taichi 操作的执行流")
	fmt.Println("   • 与其他 CUDA 操作同步")
	fmt.Println("   • 实现更精细的并发控制")
	fmt.Println("\n⚡ 使用场景：")
	fmt.Println("   • 与现有 CUDA 应用集成")
	fmt.Println("   • 多流并发执行")
	fmt.Println("   • 精确控制执行顺序")
}
