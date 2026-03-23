package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: AOT Kernel Batch Execution
// Features: Submit multiple asynchronous tasks in batch, fully utilize GPU parallel capabilities

func main() {
	fmt.Println("=== AOT Kernel Batch Execution ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan)
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Load AOT module (directory containing metadata.json)
	module, err := taichi.LoadAotModuleFile(runtime, "./examples")
	if err != nil {
		fmt.Printf("❌ Failed to load AOT module: %v\n", err)
		fmt.Println("\nPlease run the following command to generate AOT module: uv run ./examples/10_aot_kernel.py")
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

	// Create input data
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()

	// Initialize input data
	err = taichi.MapNdArray(func(arrays ...taichi.NdArrayPtr) error {
		_a := arrays[0].AsFloat32()
		_b := arrays[1].AsFloat32()
		for i := range _a {
			_a[i] = float32(i) * 0.5
			_b[i] = float32(i) * 1.5
		}
		return nil
	}, a, b)
	if err != nil {
		panic(err)
	}

	// Create multiple output arrays
	batchSize := 5
	results := make([]*taichi.NdArray, batchSize)
	for i := 0; i < batchSize; i++ {
		results[i], _ = taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
		defer results[i].Release()
	}

	fmt.Printf("✅ Prepared for batch execution of %d tasks\n\n", batchSize)

	// Submit tasks in batch
	fmt.Println("--- Submitting Tasks in Batch ---")
	for i := 0; i < batchSize; i++ {
		kernel.Launch().
			ArgNdArray(a).
			ArgNdArray(b).
			ArgNdArray(results[i]).
			RunAsync()
		fmt.Printf("✅ Task %d submitted\n", i+1)
	}

	fmt.Printf("\n✅ Submitted %d asynchronous tasks\n", batchSize)

	// Wait for all tasks to complete
	fmt.Println("\n--- Waiting for All Tasks to Complete ---")
	runtime.Wait()
	fmt.Println("✅ All tasks completed")

	// Verify results
	fmt.Println("\n--- Verifying Results ---")
	for i := 0; i < batchSize; i++ {
		results[i].MapFloat32(func(data []float32) error {
			fmt.Printf("Result %d first 3: [%.1f, %.1f, %.1f]\n", i+1, data[0], data[1], data[2])
			return nil
		})
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • Submit multiple asynchronous tasks in batch")
	fmt.Println("   • GPU can execute multiple tasks in parallel")
	fmt.Println("   • Single Wait() waits for all tasks")
	fmt.Println("   • Suitable for large-scale parallel computing")
	fmt.Println("\n⚡ Performance Benefits:")
	fmt.Println("   • Reduce CPU-GPU synchronization overhead")
	fmt.Println("   • Fully utilize GPU parallel capabilities")
	fmt.Println("   • Improve overall throughput")
}
