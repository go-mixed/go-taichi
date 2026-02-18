package main

import (
	"fmt"
	"go-taichi/taichi"
)

// 示例：创建和管理 Taichi 运行时
// 功能：演示如何创建运行时、查询架构信息、正确释放资源

func main() {
	fmt.Println("=== 运行时管理示例 ===\n")

	// 初始化 Taichi
	if err := taichi.Init(); err != nil {
		panic(err)
	}

	// 方式1: 自动选择最佳架构
	fmt.Println("--- 方式1: 自动选择架构 ---")
	runtime1, err := taichi.NewRuntimeAuto()
	if err != nil {
		panic(err)
	}
	defer runtime1.Release()

	fmt.Printf("✅ 运行时创建成功\n")
	fmt.Printf("📌 架构: %s\n", runtime1.ArchName())
	fmt.Printf("📌 架构代码: %d\n\n", runtime1.Arch())

	// 方式2: 手动指定架构
	fmt.Println("--- 方式2: 手动指定架构 ---")

	// 获取所有可用架构
	archs := taichi.GetAvailableArchs()
	fmt.Printf("可用架构数量: %d\n", len(archs))
	for i, arch := range archs {
		// 创建临时运行时来获取架构名称
		tmpRuntime, _ := taichi.NewRuntime(arch)
		if tmpRuntime != nil {
			fmt.Printf("  [%d] %s\n", i, tmpRuntime.ArchName())
			tmpRuntime.Release()
		}
	}

	// 使用第一个可用架构
	if len(archs) > 0 {
		runtime2, err := taichi.NewRuntime(archs[0])
		if err != nil {
			panic(err)
		}
		defer runtime2.Release()

		fmt.Printf("\n✅ 使用指定架构创建成功: %s\n", runtime2.ArchName())
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • 使用 defer runtime.Release() 确保资源释放")
	fmt.Println("   • NewRuntimeAuto() 自动选择最佳架构")
	fmt.Println("   • NewRuntime(arch) 手动指定架构")
	fmt.Println("   • 优先级: Vulkan > CUDA > CPU")
}
