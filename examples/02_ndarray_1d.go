package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: 1D array operations
// Features: Create, initialize, read/write 1D NdArray

func main() {
	fmt.Println("=== 1D Array Operations Example ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Create 1D array (1000 float32 elements)
	arr, err := taichi.NewNdArray1D(runtime, 1000, taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}
	defer arr.Release()

	fmt.Printf("✅ Created 1D array\n")
	fmt.Printf("   Shape: %v\n", arr.Shape())
	fmt.Printf("   Elements: %d\n", arr.TotalElements())
	fmt.Printf("   Data type: F32\n\n")

	// Write data
	fmt.Println("--- Writing Data ---")
	data, err := arr.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := range data {
		data[i] = float32(i) * 0.5
	}
	arr.Unmap()

	fmt.Printf("✅ Written %d elements\n\n", len(data))

	// Read data
	fmt.Println("--- Reading Data ---")
	data, err = arr.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	// Calculate sum
	sum := float32(0)
	for _, v := range data {
		sum += v
	}

	fmt.Printf("First 10 elements: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%.1f ", data[i])
	}
	fmt.Printf("\n")
	fmt.Printf("Sum: %.2f\n", sum)

	arr.Unmap()

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • AsSliceFloat32() maps to Go slice")
	fmt.Println("   • Modifying slice directly modifies GPU memory")
	fmt.Println("   • Call Unmap() after use")
	fmt.Println("   • defer arr.Release() to free resources")
}
