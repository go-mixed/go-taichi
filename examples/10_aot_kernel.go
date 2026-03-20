package main

import (
	"fmt"
	"os"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: AOT Kernel Basic Execution
// Features: Load AOT module, execute precompiled kernel

func main() {
	fmt.Println("=== AOT Kernel Basic Execution ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntime(taichi.ArchCuda, "./lib")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	buf, err := os.ReadFile("./examples/10_aot_kernel_cuda.tcm")
	if err != nil {
		fmt.Println("❌ Failed to read AOT module")
		panic(err)
	}
	// Load AOT module
	module, err := taichi.LoadAotModule(runtime, buf)
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

	fmt.Println("✅ Got kernel: add_kernel\n", kernel)

	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	c, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()
	defer c.Release()

	// Initialize data
	err = taichi.NdArrayAsFloat32(func(arrays ...[]float32) error {
		for i := range arrays[0] {
			arrays[0][i] = float32(i)
			arrays[1][i] = float32(i) * 2
		}
		return nil
	}, a, b)
	if err != nil {
		panic(err)
	}

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
	err = c.WithFloat32(func(dataC []float32) error {
		fmt.Printf("\nFirst 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
			dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
		fmt.Printf("Expected results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
			float32(0)+float32(0)*2, float32(1)+float32(1)*2, float32(2)+float32(2)*2, float32(3)+float32(3)*2, float32(4)+float32(4)*2)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
