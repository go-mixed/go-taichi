package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：AOT Kernel 批量执行
// 功能：批量提交多个异步任务，充分利用 GPU 并行能力

func main() {
	fmt.Println("=== AOT Kernel 批量执行 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan, "")
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

	// 创建输入数据
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()

	// 初始化输入数据
	dataA, _ := a.AsSliceFloat32()
	dataB, _ := b.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i) * 0.5
		dataB[i] = float32(i) * 1.5
	}
	a.Unmap()
	b.Unmap()

	// 创建多个输出数组
	batchSize := 5
	results := make([]*taichi.NdArray, batchSize)
	for i := 0; i < batchSize; i++ {
		results[i], _ = taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
		defer results[i].Release()
	}

	fmt.Printf("✅ 准备批量执行 %d 个任务\n\n", batchSize)

	// 批量提交任务
	fmt.Println("--- 批量提交任务 ---")
	for i := 0; i < batchSize; i++ {
		kernel.Launch().
			ArgNdArray(a).
			ArgNdArray(b).
			ArgNdArray(results[i]).
			RunAsync()
		fmt.Printf("✅ 任务 %d 已提交\n", i+1)
	}

	fmt.Printf("\n✅ 已提交 %d 个异步任务\n", batchSize)

	// 等待所有任务完成
	fmt.Println("\n--- 等待所有任务完成 ---")
	runtime.Wait()
	fmt.Println("✅ 所有任务执行完成")

	// 验证结果
	fmt.Println("\n--- 验证结果 ---")
	for i := 0; i < batchSize; i++ {
		data, _ := results[i].AsSliceFloat32()
		fmt.Printf("结果 %d 前3个: [%.1f, %.1f, %.1f]\n", i+1, data[0], data[1], data[2])
		results[i].Unmap()
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • 批量提交多个异步任务")
	fmt.Println("   • GPU 可以并行执行多个任务")
	fmt.Println("   • 一次 Wait() 等待所有任务")
	fmt.Println("   • 适合大规模并行计算")
	fmt.Println("\n⚡ 性能优势：")
	fmt.Println("   • 减少 CPU-GPU 同步开销")
	fmt.Println("   • 充分利用 GPU 并行能力")
	fmt.Println("   • 提高整体吞吐量")
}
