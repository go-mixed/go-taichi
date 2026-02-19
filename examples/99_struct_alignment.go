package main

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
	"strings"
	"unsafe"
)

// Complete struct alignment verification tool
// Used to verify that Go struct memory layout is completely consistent with C struct

func main() {
	fmt.Println("=== Go Struct Memory Layout Complete Verification ===")
	fmt.Println("\nPlease compile and run the C code below, compare if output is consistent")
	fmt.Println("\n" + strings.Repeat("=", 70) + "\n")

	// ===== Basic Types =====
	printSection("Basic Types")
	printBasicType("uintptr", unsafe.Sizeof(uintptr(0)), unsafe.Alignof(uintptr(0)))
	printBasicType("uint32", unsafe.Sizeof(uint32(0)), unsafe.Alignof(uint32(0)))
	printBasicType("uint64", unsafe.Sizeof(uint64(0)), unsafe.Alignof(uint64(0)))
	printBasicType("int32", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	printBasicType("float32", unsafe.Sizeof(float32(0)), unsafe.Alignof(float32(0)))
	printBasicType("*byte", unsafe.Sizeof((*byte)(nil)), unsafe.Alignof((*byte)(nil)))

	// ===== Image Related Structs =====
	printSection("TiImageExtent")
	var imageExtent c_api.TiImageExtent
	printStructInfo(unsafe.Sizeof(imageExtent), unsafe.Alignof(imageExtent))
	printField("Width", unsafe.Offsetof(imageExtent.Width))
	printField("Height", unsafe.Offsetof(imageExtent.Height))
	printField("Depth", unsafe.Offsetof(imageExtent.Depth))
	printField("ArrayLayerCount", unsafe.Offsetof(imageExtent.ArrayLayerCount))

	printSection("TiImageOffset")
	var imageOffset c_api.TiImageOffset
	printStructInfo(unsafe.Sizeof(imageOffset), unsafe.Alignof(imageOffset))
	printField("X", unsafe.Offsetof(imageOffset.X))
	printField("Y", unsafe.Offsetof(imageOffset.Y))
	printField("Z", unsafe.Offsetof(imageOffset.Z))
	printField("ArrayLayerOffset", unsafe.Offsetof(imageOffset.ArrayLayerOffset))

	printSection("TiImageAllocateInfo")
	var imageAllocInfo c_api.TiImageAllocateInfo
	printStructInfo(unsafe.Sizeof(imageAllocInfo), unsafe.Alignof(imageAllocInfo))
	printField("Dimension", unsafe.Offsetof(imageAllocInfo.Dimension))
	printField("Extent", unsafe.Offsetof(imageAllocInfo.Extent))
	printField("MipLevelCount", unsafe.Offsetof(imageAllocInfo.MipLevelCount))
	printField("Format", unsafe.Offsetof(imageAllocInfo.Format))
	printField("Export", unsafe.Offsetof(imageAllocInfo.Export))
	printField("Usage", unsafe.Offsetof(imageAllocInfo.Usage))

	printSection("TiImageSlice")
	var imgSlice c_api.TiImageSlice
	printStructInfo(unsafe.Sizeof(imgSlice), unsafe.Alignof(imgSlice))
	printField("Image", unsafe.Offsetof(imgSlice.Image))
	printField("Offset", unsafe.Offsetof(imgSlice.Offset))
	printField("Extent", unsafe.Offsetof(imgSlice.Extent))
	printField("MipLevel", unsafe.Offsetof(imgSlice.MipLevel))

	// ===== Sampler Related =====
	printSection("TiSamplerCreateInfo")
	var samplerInfo c_api.TiSamplerCreateInfo
	printStructInfo(unsafe.Sizeof(samplerInfo), unsafe.Alignof(samplerInfo))
	printField("MagFilter", unsafe.Offsetof(samplerInfo.MagFilter))
	printField("MinFilter", unsafe.Offsetof(samplerInfo.MinFilter))
	printField("AddressMode", unsafe.Offsetof(samplerInfo.AddressMode))
	printField("MaxAnisotropy", unsafe.Offsetof(samplerInfo.MaxAnisotropy))

	// ===== Memory Related Structs =====
	printSection("TiMemoryAllocateInfo")
	var memAllocInfo c_api.TiMemoryAllocateInfo
	printStructInfo(unsafe.Sizeof(memAllocInfo), unsafe.Alignof(memAllocInfo))
	printField("Size", unsafe.Offsetof(memAllocInfo.Size))
	printField("HostWrite", unsafe.Offsetof(memAllocInfo.HostWrite))
	printField("HostRead", unsafe.Offsetof(memAllocInfo.HostRead))
	printField("Export", unsafe.Offsetof(memAllocInfo.Export))
	printField("Usage", unsafe.Offsetof(memAllocInfo.Usage))

	printSection("TiMemorySlice")
	var memSlice c_api.TiMemorySlice
	printStructInfo(unsafe.Sizeof(memSlice), unsafe.Alignof(memSlice))
	printField("Memory", unsafe.Offsetof(memSlice.Memory))
	printField("Offset", unsafe.Offsetof(memSlice.Offset))
	printField("Size", unsafe.Offsetof(memSlice.Size))

	// ===== Array Related Structs =====
	printSection("TiNdShape")
	var shape c_api.TiNdShape
	printStructInfo(unsafe.Sizeof(shape), unsafe.Alignof(shape))
	printField("DimCount", unsafe.Offsetof(shape.DimCount))
	printField("Dims", unsafe.Offsetof(shape.Dims))

	printSection("TiNdArray")
	var ndarray c_api.TiNdArray
	printStructInfo(unsafe.Sizeof(ndarray), unsafe.Alignof(ndarray))
	printField("Memory", unsafe.Offsetof(ndarray.Memory))
	printField("Shape", unsafe.Offsetof(ndarray.Shape))
	printField("ElemShape", unsafe.Offsetof(ndarray.ElemShape))
	printField("ElemType", unsafe.Offsetof(ndarray.ElemType))

	// ===== Texture Related =====
	printSection("TiTexture")
	var texture c_api.TiTexture
	printStructInfo(unsafe.Sizeof(texture), unsafe.Alignof(texture))
	printField("Image", unsafe.Offsetof(texture.Image))
	printField("Sampler", unsafe.Offsetof(texture.Sampler))
	printField("Dimension", unsafe.Offsetof(texture.Dimension))
	printField("Extent", unsafe.Offsetof(texture.Extent))
	printField("Format", unsafe.Offsetof(texture.Format))

	// ===== Scalar Related =====
	printSection("TiScalarValue (Union)")
	var scalarVal c_api.TiScalarValue
	printStructInfo(unsafe.Sizeof(scalarVal), unsafe.Alignof(scalarVal))
	printField("Data", unsafe.Offsetof(scalarVal.Data))
	fmt.Println("  Note: Union size should equal largest member (8 bytes: i64/f64)")

	printSection("TiScalar")
	var scalar c_api.TiScalar
	printStructInfo(unsafe.Sizeof(scalar), unsafe.Alignof(scalar))
	printField("Type", unsafe.Offsetof(scalar.Type))
	printField("Value", unsafe.Offsetof(scalar.Value))
	if unsafe.Sizeof(scalar) == 12 {
		fmt.Println("  ⚠️  Warning: Size is 12 bytes, may need padding to 16 bytes")
	}

	// ===== Argument Related =====
	printSection("TiArgumentValue (Union)")
	var argVal c_api.TiArgumentValue
	printStructInfo(unsafe.Sizeof(argVal), unsafe.Alignof(argVal))
	printField("Data", unsafe.Offsetof(argVal.Data))
	fmt.Println("  Note: Union size should equal largest member TiNdArray (152 bytes)")

	printSection("TiArgument")
	var arg c_api.TiArgument
	printStructInfo(unsafe.Sizeof(arg), unsafe.Alignof(arg))
	printField("Type", unsafe.Offsetof(arg.Type))
	printField("Value", unsafe.Offsetof(arg.Value))
	if unsafe.Offsetof(arg.Value) == 8 {
		fmt.Println("  ✅ 4 bytes padding between Type and Value")
	} else {
		fmt.Println("  ❌ Warning: Value offset is not 8, may be missing padding")
	}

	printSection("TiNamedArgument")
	var namedArg c_api.TiNamedArgument
	printStructInfo(unsafe.Sizeof(namedArg), unsafe.Alignof(namedArg))
	printField("Name", unsafe.Offsetof(namedArg.Name))
	printField("Argument", unsafe.Offsetof(namedArg.Argument))

	// ===== Data Integrity Test =====
	printSection("Data Integrity Test")
	testDataIntegrity()

	// ===== Raw Byte Check =====
	printSection("Raw Byte Check - TiArgument first 32 bytes")
	testArg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_NDARRAY,
	}
	testNdArray := c_api.TiNdArray{
		Memory: c_api.TiMemory(0x123456789ABC),
		Shape: c_api.TiNdShape{
			DimCount: 2,
			Dims:     [16]uint32{10, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		ElemType: c_api.TI_DATA_TYPE_F32,
	}
	*(*c_api.TiNdArray)(unsafe.Pointer(&testArg.Value.Data[0])) = testNdArray

	argBytes := (*[32]byte)(unsafe.Pointer(&testArg))
	for i := 0; i < 32; i += 16 {
		fmt.Printf("  %04X: ", i)
		for j := 0; j < 16 && i+j < 32; j++ {
			fmt.Printf("%02X ", argBytes[i+j])
		}
		if i == 0 {
			fmt.Printf(" <- Type (4 bytes) + Padding (4 bytes)")
		} else if i == 16 {
			fmt.Printf(" <- Part of Memory")
		}
		fmt.Println()
	}

	// ===== Summary =====
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n=== Verification Complete ===")
	fmt.Println("\nIf all sizes and offsets match the C version, the memory layout is correct.")
	fmt.Println("\n⚠️  Pay special attention to these structs:")
	fmt.Println("   • TiScalar - check if padding is needed")
	fmt.Println("   • TiArgument - confirm Value offset is 8")
	fmt.Println("   • TiNdArray - confirm ElemShape offset is 76")
}

// Helper functions
func printSection(name string) {
	fmt.Printf("\n【%s】\n", name)
}

func printBasicType(name string, size, align uintptr) {
	fmt.Printf("  %-10s size=%2d, align=%d\n", name+":", size, align)
}

func printStructInfo(size, align uintptr) {
	fmt.Printf("  Total size:  %3d bytes\n", size)
	fmt.Printf("  Alignment:   %3d bytes\n", align)
}

func printField(name string, offset uintptr) {
	fmt.Printf("  %-20s offset: %3d\n", name, offset)
}

func testDataIntegrity() {
	testNdArray := c_api.TiNdArray{
		Memory: c_api.TiMemory(0x123456789ABC),
		Shape: c_api.TiNdShape{
			DimCount: 2,
			Dims:     [16]uint32{10, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		ElemShape: c_api.TiNdShape{
			DimCount: 0,
		},
		ElemType: c_api.TI_DATA_TYPE_F32,
	}

	testArg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_NDARRAY,
	}

	// 复制到 ArgumentValue
	*(*c_api.TiNdArray)(unsafe.Pointer(&testArg.Value.Data[0])) = testNdArray

	// 读回验证
	readBack := *(*c_api.TiNdArray)(unsafe.Pointer(&testArg.Value.Data[0]))

	fmt.Printf("  Original Memory:    0x%X\n", testNdArray.Memory)
	fmt.Printf("  Read back Memory:   0x%X\n", readBack.Memory)
	fmt.Printf("  Original DimCount:  %d\n", testNdArray.Shape.DimCount)
	fmt.Printf("  Read back DimCount: %d\n", readBack.Shape.DimCount)
	fmt.Printf("  Original Dims[0]:   %d\n", testNdArray.Shape.Dims[0])
	fmt.Printf("  Read back Dims[0]:  %d\n", readBack.Shape.Dims[0])
	fmt.Printf("  Original ElemType:  %d\n", testNdArray.ElemType)
	fmt.Printf("  Read back ElemType: %d\n", readBack.ElemType)

	if readBack.Memory == testNdArray.Memory &&
		readBack.Shape.DimCount == testNdArray.Shape.DimCount &&
		readBack.Shape.Dims[0] == testNdArray.Shape.Dims[0] &&
		readBack.ElemType == testNdArray.ElemType {
		fmt.Println("\n  ✅ Data integrity test passed!")
	} else {
		fmt.Println("\n  ❌ Data integrity test failed!")
	}
}
