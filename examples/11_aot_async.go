package main

import (
	"fmt"
	"go-taichi/taichi"
)

// 示例：AOT Kernel 异步执行
// 功能：使用 RunAsync() 异步执行 Kernel，使用 Wait() 等待完成

func main() {
	fmt.Println("=== AOT Kernel 异步执行 ===\n")

	// 初始化
	taichi.Init()

	// 创建运行时
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan)
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 加载 AOT 模块
	module, err := taichi.LoadAotModule(runtime, "./examples/10_aot_module.tcm")
	if err != nil {
		fmt.Printf("❌ 加载 AOT 模块失败: %v\n", err)
		fmt.Println("\n请先运行以下命令生成 AOT 模块： uv run ./examples/10_aot_kenerl.py")
		return
	}
	defer module.Release()

	// 获取 kernel
	kernel, err := module.GetKernel("add_kernel")
	if err != nil {
		fmt.Printf("❌ 获取 kernel 失败: %v\n", err)
		return
	}

	fmt.Println("✅ AOT 模块和 Kernel 加载成功\n")

	// 创建测试数据
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	c, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()
	defer c.Release()

	// 初始化数据
	dataA, _ := a.AsSliceFloat32()
	dataB, _ := b.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i) * 0.5
		dataB[i] = float32(i) * 1.5
	}
	a.Unmap()
	b.Unmap()

	fmt.Println("✅ 测试数据准备完成")

	// 异步执行 kernel
	fmt.Println("\n--- 异步执行 ---")
	kernel.Launch().
		ArgNdArray(a).
		ArgNdArray(b).
		ArgNdArray(c).
		RunAsync()

	fmt.Println("✅ 异步任务已提交")
	fmt.Println("   (可以继续执行其他操作...)")

	// 等待完成
	fmt.Println("\n--- 等待任务完成 ---")
	runtime.Wait()
	fmt.Println("✅ 异步任务执行完成")

	// 检查结果
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\n前5个结果: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	c.Unmap()

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • RunAsync() 异步提交任务")
	fmt.Println("   • 提交后立即返回，不阻塞")
	fmt.Println("   • runtime.Wait() 等待所有任务完成")
	fmt.Println("   • 适合 CPU 和 GPU 并行工作")
}
