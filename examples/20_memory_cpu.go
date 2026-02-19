package main

import (
	"fmt"
	"unsafe"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：CPU 内存导入
// 功能：将 Go 切片内存导入到 Taichi，避免数据复制

func main() {
	fmt.Println("=== CPU 内存导入示例 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 创建 Go 切片
	data := make([]float32, 1000)
	for i := range data {
		data[i] = float32(i) * 0.5
	}

	fmt.Printf("✅ 创建 Go 切片: %d 个元素\n", len(data))
	fmt.Printf("   内存大小: %d bytes\n", len(data)*4)
	fmt.Printf("   内存地址: %p\n\n", &data[0])

	// 导入 CPU 内存
	fmt.Println("--- 导入内存到 Taichi ---")
	memory, err := taichi.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
	if err != nil {
		fmt.Printf("❌ CPU 内存导入失败: %v\n", err)
		fmt.Println("\n⚠️  某些后端不支持 CPU 内存导入")
		fmt.Println("   支持的后端: CPU (x64/ARM64)")
		fmt.Println("   不支持的后端: Vulkan, CUDA, OpenGL")
		return
	}
	defer memory.Release()

	fmt.Printf("✅ 成功导入 CPU 内存\n")
	fmt.Printf("   大小: %d bytes\n\n", memory.Size())

	// 从导入的内存创建 NdArray
	fmt.Println("--- 创建 NdArray ---")
	arr, err := taichi.NewNdArray1DFromMemory(memory, uint32(len(data)), taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ 创建 NdArray\n")
	fmt.Printf("   形状: %v\n", arr.Shape())
	fmt.Printf("   元素数: %d\n\n", arr.TotalElements())

	// 注意事项
	fmt.Println("--- 注意事项 ---")
	fmt.Println("⚠️  导入的内存可能无法直接映射")
	fmt.Println("   • 导入的内存主要用于与 Taichi kernel 交互")
	fmt.Println("   • 不一定支持 AsSliceXXX() 映射")
	fmt.Println("   • 数据修改会直接影响原始 Go 切片")

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • ImportCPUMemory() 导入现有内存")
	fmt.Println("   • 避免数据复制，提高性能")
	fmt.Println("   • 导入的内存生命周期必须超过 Memory 对象")
	fmt.Println("   • 仅 CPU 后端支持")
	fmt.Println("\n⚡ 使用场景：")
	fmt.Println("   • 与现有 Go 代码集成")
	fmt.Println("   • 大数据量避免复制")
	fmt.Println("   • 零拷贝数据传输")
}
