package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: Basic Image Operations
// Features: Create images, layout transitions, query properties

func main() {
	fmt.Println("=== Basic Image Operations Example ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Create 512x512 RGBA8 image
	img, err := taichi.NewImage2D(runtime, 512, 512, taichi.FormatRgba8)
	if err != nil {
		panic(err)
	}
	defer img.Release()

	fmt.Printf("✅ Image created\n")
	fmt.Printf("   Size: %dx%d\n", img.Width(), img.Height())
	fmt.Printf("   Format: RGBA8\n")
	fmt.Printf("   Channels: 4\n\n")

	// Layout transition examples
	fmt.Println("--- Layout Transitions ---")

	// Transition to shader write layout
	img.TransitionLayout(taichi.ImageLayoutShaderWrite)
	fmt.Println("✅ Transitioned to ShaderWrite layout")

	// Transition to shader read layout
	img.TransitionLayout(taichi.ImageLayoutShaderRead)
	fmt.Println("✅ Transitioned to ShaderRead layout")

	// Transition to transfer destination layout
	img.TransitionLayout(taichi.ImageLayoutTransferDst)
	fmt.Println("✅ Transitioned to TransferDst layout")

	// Transition to transfer source layout
	img.TransitionLayout(taichi.ImageLayoutTransferSrc)
	fmt.Println("✅ Transitioned to TransferSrc layout")

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • Image layout must match the operation")
	fmt.Println("   • ShaderWrite - for shader writes")
	fmt.Println("   • ShaderRead - for shader reads")
	fmt.Println("   • TransferDst - for receiving data")
	fmt.Println("   • TransferSrc - for sending data")
	fmt.Println("\n📚 Common Formats:")
	fmt.Println("   • RGBA8 - 8-bit RGBA (common)")
	fmt.Println("   • RGBA16F - 16-bit float RGBA")
	fmt.Println("   • RGBA32F - 32-bit float RGBA")
	fmt.Println("   • R32F - 32-bit float single channel")
}
