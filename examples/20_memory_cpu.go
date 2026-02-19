package main

import (
	"fmt"
	"unsafe"

	"github.com/go-mixed/go-taichi/taichi"
)

// Example: CPU Memory Import
// Features: Import Go slice memory into Taichi, avoid data copying

func main() {
	fmt.Println("=== CPU Memory Import Example ===\n")

	// Create runtime
	runtime, err := taichi.NewRuntimeAuto("")
	if err != nil {
		panic(err)
	}
	defer runtime.Release()

	fmt.Printf("✅ Runtime: %s\n\n", runtime.ArchName())

	// Create Go slice
	data := make([]float32, 1000)
	for i := range data {
		data[i] = float32(i) * 0.5
	}

	fmt.Printf("✅ Created Go slice: %d elements\n", len(data))
	fmt.Printf("   Memory size: %d bytes\n", len(data)*4)
	fmt.Printf("   Memory address: %p\n\n", &data[0])

	// Import CPU memory
	fmt.Println("--- Import Memory to Taichi ---")
	memory, err := taichi.ImportCPUMemory(runtime, unsafe.Pointer(&data[0]), uint64(len(data)*4))
	if err != nil {
		fmt.Printf("❌ CPU memory import failed: %v\n", err)
		fmt.Println("\n⚠️  Some backends do not support CPU memory import")
		fmt.Println("   Supported backends: CPU (x64/ARM64)")
		fmt.Println("   Unsupported backends: Vulkan, CUDA, OpenGL")
		return
	}
	defer memory.Release()

	fmt.Printf("✅ Successfully imported CPU memory\n")
	fmt.Printf("   Size: %d bytes\n\n", memory.Size())

	// Create NdArray from imported memory
	fmt.Println("--- Create NdArray ---")
	arr, err := taichi.NewNdArray1DFromMemory(memory, uint32(len(data)), taichi.DataTypeF32)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Created NdArray\n")
	fmt.Printf("   Shape: %v\n", arr.Shape())
	fmt.Printf("   Element count: %d\n\n", arr.TotalElements())

	// Notes
	fmt.Println("--- Notes ---")
	fmt.Println("⚠️  Imported memory may not be directly mappable")
	fmt.Println("   • Imported memory is mainly for interacting with Taichi kernels")
	fmt.Println("   • May not support AsSliceXXX() mapping")
	fmt.Println("   • Data modifications directly affect original Go slice")

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\n💡 Key Points:")
	fmt.Println("   • ImportCPUMemory() imports existing memory")
	fmt.Println("   • Avoid data copying, improve performance")
	fmt.Println("   • Imported memory lifetime must exceed Memory object")
	fmt.Println("   • Only CPU backend supported")
	fmt.Println("\n⚡ Use Cases:")
	fmt.Println("   • Integration with existing Go code")
	fmt.Println("   • Avoid copying large data volumes")
	fmt.Println("   • Zero-copy data transfer")
}
