package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：AOT Kernel 基础执行
// 功能：加载 AOT 模块，执行预编译的 Kernel

func main() {
	fmt.Println("=== AOT Kernel 基础执行 ===\n")

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

	fmt.Println("✅ AOT 模块加载成功")

	// 获取 kernel
	kernel, err := module.GetKernel("add_kernel")
	if err != nil {
		fmt.Printf("❌ 获取 kernel 失败: %v\n", err)
		return
	}

	fmt.Println("✅ 获取 Kernel: add_kernel\n")

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
		dataA[i] = float32(i)
		dataB[i] = float32(i) * 2
	}
	a.Unmap()
	b.Unmap()

	fmt.Println("✅ 测试数据准备完成")

	// 执行 kernel (Builder 模式)
	fmt.Println("\n--- 执行 Kernel ---")
	kernel.Launch().
		ArgNdArray(a).
		ArgNdArray(b).
		ArgNdArray(c).
		Run()

	fmt.Println("✅ Kernel 执行完成")

	// 检查结果
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\n前5个结果: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	fmt.Printf("预期结果:  [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataA[0]+dataB[0], dataA[1]+dataB[1], dataA[2]+dataB[2], dataA[3]+dataB[3], dataA[4]+dataB[4])
	c.Unmap()
}
