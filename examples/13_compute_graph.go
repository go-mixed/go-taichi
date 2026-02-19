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
	runtime, err := taichi.NewRuntime(taichi.ArchVulkan, "")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Load AOT module
	module, err := taichi.LoadAotModule(runtime, "./exmaples/10_aot_module.tcm")
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
	dataA, _ := a.AsSliceFloat32()
	dataB, _ := b.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i)
		dataB[i] = float32(i) * 2
	}
	a.Unmap()
	b.Unmap()

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
	dataC, _ := c.AsSliceFloat32()
	fmt.Printf("\nFirst 10 results: ")
	for i := 0; i < 10 && i < len(dataC); i++ {
		fmt.Printf("%.1f ", dataC[i])
	}
	fmt.Println()
	c.Unmap()

	// Asynchronous execution example
	fmt.Println("\n--- Asynchronous Execution ---")

	// Reset data
	dataA, _ = a.AsSliceFloat32()
	for i := range dataA {
		dataA[i] = float32(i) * 0.1
	}
	a.Unmap()

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
	dataC, _ = c.AsSliceFloat32()
	fmt.Printf("\nAsync execution first 5 results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
		dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
	c.Unmap()

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Compute Graph vs Kernel:")
	fmt.Println("   • Kernel: Single computation function")
	fmt.Println("   • Compute Graph: Directed acyclic graph of multiple kernels")
	fmt.Println("   • Compute Graph can optimize entire computation pipeline")
	fmt.Println("   • Uses named parameters for more flexible parameter passing")
}
