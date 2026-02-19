package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: AOT Kernel Basic Execution
// Features: Load AOT module, execute precompiled kernel

func main() {
	fmt.Println("=== AOT Kernel Basic Execution ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan, "")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Load AOT module
	module, err := taichi.LoadAotModule(runtime, "./examples/10_aot_module.tcm")
	if err != nil {
		fmt.Printf("❌ Failed to load AOT module: %v\n", err)
		fmt.Println("\nPlease run the following command to generate AOT module: uv run ./examples/10_aot_kenerl.py")
		return
	}
	defer module.Release()

	fmt.Println("✅ AOT module loaded successfully")

	// Get kernel
	kernel, err := module.GetKernel("add_kernel")
	if err != nil {
		fmt.Printf("❌ Failed to get kernel: %v\n", err)
		return
	}

	fmt.Println("✅ Got kernel: add_kernel\n")

	// Create test data
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	c, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()
	defer c.Release()

	// Initialize data
	dataA, _ := a.AsSliceFloat32()
	dataB, _ := b.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i)
		dataB[i] = float32(i) * 2
	}
	a.Unmap()
	b.Unmap()

	fmt.Println("✅ Test data prepared")

	// Execute kernel (Builder pattern)
	fmt.Println("\n--- Execute Kernel ---")
	kernel.Launch().
		ArgNdArray(a).
		ArgNdArray(b).
		ArgNdArray(c).
		Run()

	fmt.Println("✅ Kernel execution completed")

	// Check results
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\nFirst 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	fmt.Printf("Expected results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataA[0]+dataB[0], dataA[1]+dataB[1], dataA[2]+dataB[2], dataA[3]+dataB[3], dataA[4]+dataB[4])
	c.Unmap()
}
