package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：Compute Graph 执行
// 功能：执行包含多个 Kernel 的计算图，使用命名参数

func main() {
	fmt.Println("=== Compute Graph 示例 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan, "")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 加载 AOT 模块
	module, err := taichi.LoadAotModule(runtime, "./exmaples/10_aot_module.tcm")
	if err != nil {
		fmt.Printf("❌ 加载 AOT 模块失败: %v\n", err)
		fmt.Println("\n请先运行以下命令生成包含 Compute Graph 的 AOT 模块：")
		fmt.Println("  python generate_compute_graph.py")
		return
	}
	defer module.Release()

	fmt.Println("✅ AOT 模块加载成功")

	// 获取 compute graph
	graph, err := module.GetComputeGraph("my_compute_graph")
	if err != nil {
		fmt.Printf("❌ 获取 Compute Graph 失败: %v\n", err)
		fmt.Println("\n当前 AOT 模块可能只包含 kernel，没有 compute graph。")
		fmt.Println("Compute Graph 是多个 kernel 的组合，需要在 Python 中定义。")
		fmt.Println("\n💡 提示：如果只需要使用 Kernel，请运行 10_aot_kernel.go")
		return
	}

	fmt.Printf("✅ 获取 Compute Graph: %s\n\n", graph.Name())

	// 创建测试数据
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	c, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()
	defer c.Release()

	// 初始化输入数据
	dataA, _ := a.AsSliceFloat32()
	dataB, _ := b.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i)
		dataB[i] = float32(i) * 2
	}
	a.Unmap()
	b.Unmap()

	fmt.Println("✅ 测试数据准备完成")

	// 执行 Compute Graph（使用命名参数）
	fmt.Println("\n--- 执行 Compute Graph ---")
	graph.Launch().
		ArgNdArray("input_a", a).
		ArgNdArray("input_b", b).
		ArgNdArray("output_c", c).
		ArgFloat32("scale_factor", 1.5).
		Run()

	fmt.Println("✅ Compute Graph 执行完成")

	// 检查结果
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\n前10个结果: ")
	for i := 0; i < 10 && i < len(dataC); i++ {
		fmt.Printf("%.1f ", dataC[i])
	}
	fmt.Println()
	c.Unmap()

	// 异步执行示例
	fmt.Println("\n--- 异步执行 ---")

	// 重置数据
	dataA, _ = a.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i) * 0.1
	}
	a.Unmap()

	// 异步执行
	graph.Launch().
		ArgNdArray("input_a", a).
		ArgNdArray("input_b", b).
		ArgNdArray("output_c", c).
		ArgFloat32("scale_factor", 2.0).
		RunAsync()

	fmt.Println("✅ 异步任务已提交")

	// 等待完成
	runtime.Wait()
	fmt.Println("✅ 异步任务执行完成")

	// 检查结果
	dataC, _ = c.AsSliceFloat32()
	fmt.Printf("\n异步执行结果前5个: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	c.Unmap()

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 Compute Graph vs Kernel：")
	fmt.Println("   • Kernel: 单个计算函数")
	fmt.Println("   • Compute Graph: 多个 kernel 的有向无环图")
	fmt.Println("   • Compute Graph 可以优化整个计算流程")
	fmt.Println("   • 使用命名参数，更灵活的参数传递")
}
