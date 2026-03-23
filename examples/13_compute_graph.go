package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: Compute Graph Execution
// Features: Execute compute graph containing multiple kernels, use named parameters

func main() {
	fmt.Println("=== Compute Graph Example ===\n")

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
		fmt.Println("\nPlease run the following command to generate AOT module with Compute Graph:")
		fmt.Println("  python generate_compute_graph.py")
		return
	}
	defer module.Release()

	fmt.Println("✅ AOT module loaded successfully")

	// Get compute graph
	graph, err := module.GetComputeGraph("my_compute_graph")
	if err != nil {
		fmt.Printf("❌ Failed to get Compute Graph: %v\n", err)
		fmt.Println("\nCurrent AOT module may only contain kernels, no compute graph.")
		fmt.Println("Compute Graph is a combination of multiple kernels, needs to be defined in Python.")
		fmt.Println("\n💡 Tip: If you only need to use Kernel, please run 10_aot_kernel.go")
		return
	}

	fmt.Printf("✅ Got Compute Graph: %s\n\n", graph.Name())

	// Create test data
	n := uint32(100)
	a, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	b, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	c, _ := taichi.NewNdArray1D(runtime, n, taichi.DataTypeF32)
	defer a.Release()
	defer b.Release()
	defer c.Release()

	// Initialize input data
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

	// Execute Compute Graph (using named parameters)
	fmt.Println("\n--- Execute Compute Graph ---")
	graph.Launch().
		ArgNdArray("input_a", a).
		ArgNdArray("input_b", b).
		ArgNdArray("output_c", c).
		ArgFloat32("scale_factor", 1.5).
		Run()

	fmt.Println("✅ Compute Graph execution completed")

	// Check results
	err = c.MapFloat32(func(dataC []float32) error {
		fmt.Printf("\nFirst 10 results: ")
		for i := 0; i < 10 && i < len(dataC); i++ {
			fmt.Printf("%.1f ", dataC[i])
		}
		fmt.Println()
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Asynchronous execution example
	fmt.Println("\n--- Asynchronous Execution ---")

	// Reset data
	err = a.MapFloat32(func(dataA []float32) error {
		for i := range dataA {
			dataA[i] = float32(i) * 0.1
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Execute asynchronously
	graph.Launch().
		ArgNdArray("input_a", a).
		ArgNdArray("input_b", b).
		ArgNdArray("output_c", c).
		ArgFloat32("scale_factor", 2.0).
		RunAsync()

	fmt.Println("✅ Asynchronous task submitted")

	// Wait for completion
	runtime.Wait()
	fmt.Println("✅ Asynchronous task completed")

	// Check results
	err = c.MapFloat32(func(dataC []float32) error {
		fmt.Printf("\nAsync execution first 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
			dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Compute Graph vs Kernel:")
	fmt.Println("   • Kernel: Single computation function")
	fmt.Println("   • Compute Graph: Directed acyclic graph of multiple kernels")
	fmt.Println("   • Compute Graph can optimize entire computation pipeline")
	fmt.Println("   • Uses named parameters for more flexible parameter passing")
}
