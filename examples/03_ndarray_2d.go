package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：2D 矩阵操作
// 功能：创建、初始化、操作 2D NdArray

func main() {
	fmt.Println("=== 2D 矩阵操作示例 ===\n")

	// 创建运行时
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 创建 4x4 矩阵
	matrix, err := taichi.NewNdArray2D(runtime, 4, 4, taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}
	defer matrix.Release()

	shape := matrix.Shape()
	rows, cols := shape[0], shape[1]

	fmt.Printf("✅ 创建 2D 矩阵\n")
	fmt.Printf("   形状: %dx%d\n", rows, cols)
	fmt.Printf("   元素数: %d\n", matrix.TotalElements())
	fmt.Printf("   数据类型: F32\n\n")

	// 初始化为单位矩阵
	fmt.Println("--- 初始化为单位矩阵 ---")
	data, err := matrix.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := uint32(0); i < rows; i++ {
		for j := uint32(0); j < cols; j++ {
			if i == j {
				data[i*cols+j] = 1.0
			} else {
				data[i*cols+j] = 0.0
			}
		}
	}
	matrix.Unmap()

	fmt.Println("✅ 初始化完成\n")

	// 读取并打印矩阵
	fmt.Println("--- 矩阵内容 ---")
	data, err = matrix.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := uint32(0); i < rows; i++ {
		for j := uint32(0); j < cols; j++ {
			fmt.Printf("%.0f ", data[i*cols+j])
		}
		fmt.Println()
	}

	matrix.Unmap()

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • 2D 数组在内存中是行优先存储")
	fmt.Println("   • 索引计算: data[row*cols + col]")
	fmt.Println("   • Shape() 返回 [rows, cols]")
}
