package main

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi"
)

// 示例：图像基础操作
// 功能：创建图像、布局转换、查询属性

func main() {
	fmt.Println("=== 图像基础操作示例 ===\n")

	// 初始化
	taichi.Init()

	// 创建运行时
	runtime, err := taichi.NewRuntimeAuto()
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ 运行时: %s\n\n", runtime.ArchName())

	// 创建 512x512 RGBA8 图像
	img, err := taichi.NewImage2D(runtime, 512, 512, taichi.FormatRgba8)
	if err != nil {
		panic(err)
	}
	defer img.Release()

	fmt.Printf("✅ 创建图像\n")
	fmt.Printf("   尺寸: %dx%d\n", img.Width(), img.Height())
	fmt.Printf("   格式: RGBA8\n")
	fmt.Printf("   通道数: 4\n\n")

	// 布局转换示例
	fmt.Println("--- 布局转换 ---")

	// 转换为 Shader 写入布局
	img.TransitionLayout(taichi.ImageLayoutShaderWrite)
	fmt.Println("✅ 转换为 ShaderWrite 布局")

	// 转换为 Shader 读取布局
	img.TransitionLayout(taichi.ImageLayoutShaderRead)
	fmt.Println("✅ 转换为 ShaderRead 布局")

	// 转换为传输目标布局
	img.TransitionLayout(taichi.ImageLayoutTransferDst)
	fmt.Println("✅ 转换为 TransferDst 布局")

	// 转换为传输源布局
	img.TransitionLayout(taichi.ImageLayoutTransferSrc)
	fmt.Println("✅ 转换为 TransferSrc 布局")

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 要点：")
	fmt.Println("   • 图像布局必须与操作匹配")
	fmt.Println("   • ShaderWrite - 用于 Shader 写入")
	fmt.Println("   • ShaderRead - 用于 Shader 读取")
	fmt.Println("   • TransferDst - 用于接收数据")
	fmt.Println("   • TransferSrc - 用于发送数据")
	fmt.Println("\n📚 常用格式：")
	fmt.Println("   • RGBA8 - 8位RGBA (常用)")
	fmt.Println("   • RGBA16F - 16位浮点RGBA")
	fmt.Println("   • RGBA32F - 32位浮点RGBA")
	fmt.Println("   • R32F - 32位浮点单通道")
}
