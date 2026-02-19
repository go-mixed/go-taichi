package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：1D 数组操作
// 功能：创建、初始化、读写 1D NdArray

func main() {
	fmt.Println("=== 1D 数组操作示例 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 创建 1D 数组 (1000 个 float32)
	arr, err := taichi.NewNdArray1D(runtime, 1000, taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}
	defer arr.Release()

	fmt.Printf("✅ 创建 1D 数组\n")
	fmt.Printf("   形状: %v\n", arr.Shape())
	fmt.Printf("   元素数: %d\n", arr.TotalElements())
	fmt.Printf("   数据类型: F32\n\n")

	// 写入数据
	fmt.Println("--- 写入数据 ---")
	data, err := arr.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := range data {
		data[i] = float32(i) * 0.5
	}
	arr.Unmap()

	fmt.Printf("✅ 已写入 %d 个元素\n\n", len(data))

	// 读取数据
	fmt.Println("--- 读取数据 ---")
	data, err = arr.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	// 计算总和
	sum := float32(0)
	for _, v := range data {
		sum += v
	}

	fmt.Printf("前10个元素: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%.1f ", data[i])
	}
	fmt.Printf("\n")
	fmt.Printf("数据总和: %.2f\n", sum)

	arr.Unmap()

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • AsSliceFloat32() 映射为 Go 切片")
	fmt.Println("   • 修改切片会直接修改 GPU 内存")
	fmt.Println("   • 使用完毕后调用 Unmap()")
	fmt.Println("   • defer arr.Release() 释放资源")
}
