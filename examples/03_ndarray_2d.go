package main

import (
	"fmt"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: 2D matrix operations
// Features: Create, initialize, manipulate 2D NdArray

func main() {
	fmt.Println("=== 2D Matrix Operations Example ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Create 4x4 matrix
	matrix, err := taichi.NewNdArray2D(runtime, 4, 4, taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}
	defer matrix.Release()

	shape := matrix.Shape()
	rows, cols := shape[0], shape[1]

	fmt.Printf("✅ Created 2D matrix\n")
	fmt.Printf("   Shape: %dx%d\n", rows, cols)
	fmt.Printf("   Elements: %d\n", matrix.TotalElements())
	fmt.Printf("   Data type: F32\n\n")

	// Initialize as identity matrix
	fmt.Println("--- Initialize as Identity Matrix ---")
	data, err := matrix.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := uint32(0); i < rows; i++ {
		for j := uint32(0); j < cols; j++ {
			if i == j {
				data[i*cols+j] = 1.0
			} else {
				data[i*cols+j] = 0.0
			}
		}
	}
	matrix.Unmap()

	fmt.Println("✅ Initialization complete\n")

	// Read and print matrix
	fmt.Println("--- Matrix Content ---")
	data, err = matrix.AsSliceFloat32()
	if err != nil {
		panic(err)
	}

	for i := uint32(0); i < rows; i++ {
		for j := uint32(0); j < cols; j++ {
			fmt.Printf("%.0f ", data[i*cols+j])
		}
		fmt.Println()
	}

	matrix.Unmap()

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • 2D arrays are stored in row-major order")
	fmt.Println("   • Index calculation: data[row*cols + col]")
	fmt.Println("   • Shape() returns [rows, cols]")
}
