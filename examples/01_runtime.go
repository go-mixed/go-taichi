package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: Create and manage Taichi runtime
// Features: Demonstrates how to create runtime, query architecture info, and properly release resources

func main() {
	fmt.Println("=== Runtime Management Example ===\n")

	// Method 1: Auto-select best architecture
	fmt.Println("--- Method 1: Auto-select Architecture ---")
	runtime1, err := taichi.NewRuntimeAuto()
	if err != nil {
		panic(err)
	}
	defer runtime1.Release()

	fmt.Printf("✅ Runtime created successfully\n")
	fmt.Printf("📌 Architecture: %s\n", runtime1.ArchName())
	fmt.Printf("📌 Architecture code: %d\n\n", runtime1.Arch())

	// Method 2: Manually specify architecture
	fmt.Println("--- Method 2: Manually Specify Architecture ---")

	// Get all available architectures
	archs := taichi.GetAvailableArchs()
	fmt.Printf("Available architectures: %d\n", len(archs))
	for i, arch := range archs {
		// Create temporary runtime to get architecture name
		tmpRuntime, _ := taichi.NewRuntime(arch)
		if tmpRuntime != nil {
			fmt.Printf("  [%d] %s\n", i, tmpRuntime.ArchName())
			tmpRuntime.Release()
		}
	}

	// Use the first available architecture
	if len(archs) > 0 {
		runtime2, err := taichi.NewRuntime(archs[0])
		if err != nil {
			panic(err)
		}
		defer runtime2.Release()

		fmt.Printf("\n✅ Created with specified architecture: %s\n", runtime2.ArchName())
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • Use defer runtime.Release() to ensure resource cleanup")
	fmt.Println("   • NewRuntimeAuto() auto-selects best architecture")
	fmt.Println("   • NewRuntime(arch) manually specifies architecture")
	fmt.Println("   • Priority: Vulkan > CUDA > CPU")
}
