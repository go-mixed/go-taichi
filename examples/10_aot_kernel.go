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
	err = taichi.MapNdArray(func(arrays ...taichi.NdArrayPtr) error {
		_a := arrays[0].AsFloat32()
		_b := arrays[1].AsFloat32()
		for i := range _a {
			_a[i] = float32(i)
			_b[i] = float32(i) * 2
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
	err = c.MapFloat32(func(dataC []float32) error {
		fmt.Printf("\nFirst 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
			dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
		fmt.Printf("Expected results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
			float32(0)+float32(0)*2, float32(1)+float32(1)*2, float32(2)+float32(2)*2, float32(3)+float32(3)*2, float32(4)+float32(4)*2)
		return nil
	})
	if err != nil {
		panic(err)
	}

	//d, _ := taichi.NewNdArray3D(runtime, 1024, 1024, 4, taichi.DataTypeF32)
	d, _ := taichi.NewNdArray2DWithElemShape(runtime, 1024, 1024, taichi.Shape(4), taichi.DataTypeF32)
	// Get kernel
	kernel, err = module.GetKernel("fill_texture")
	if err != nil {
		fmt.Printf("❌ Failed to get kernel: %v\n", err)
		return
	}

	kernel.Launch().
		ArgNdArray(d).
		ArgFloat32(0.5).
		ArgFloat32(0.5).
		ArgFloat32(0.5).
		ArgFloat32(1.).
		Run()

	_ = d.MapFloat32(func(data []float32) error {
		fmt.Println("x=0, y=0, rgba=", data[0], data[1], data[2], data[3])
		return nil
	})
}
