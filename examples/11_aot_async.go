package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: AOT Kernel Asynchronous Execution
// Features: Use RunAsync() to execute kernel asynchronously, use Wait() to wait for completion

func main() {
	fmt.Println("=== AOT Kernel Asynchronous Execution ===\n")

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

	// Get kernel
	kernel, err := module.GetKernel("add_kernel")
	if err != nil {
		fmt.Printf("❌ Failed to get kernel: %v\n", err)
		return
	}

	fmt.Println("✅ AOT module and kernel loaded successfully\n")

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
		dataA[i] = float32(i) * 0.5
		dataB[i] = float32(i) * 1.5
	}
	a.Unmap()
	b.Unmap()

	fmt.Println("✅ Test data prepared")

	// Execute kernel asynchronously
	fmt.Println("\n--- Asynchronous Execution ---")
	kernel.Launch().
		ArgNdArray(a).
		ArgNdArray(b).
		ArgNdArray(c).
		RunAsync()

	fmt.Println("✅ Asynchronous task submitted")
	fmt.Println("   (Can continue with other operations...)")

	// Wait for completion
	fmt.Println("\n--- Waiting for Task Completion ---")
	runtime.Wait()
	fmt.Println("✅ Asynchronous task completed")

	// Check results
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\nFirst 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	c.Unmap()

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • RunAsync() submits task asynchronously")
	fmt.Println("   • Returns immediately after submission, non-blocking")
	fmt.Println("   • runtime.Wait() waits for all tasks to complete")
	fmt.Println("   • Suitable for CPU and GPU parallel work")
}
