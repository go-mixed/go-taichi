package c_api

import (
	"unsafe"
)

// ===== Argument Construction Helper Functions =====

// NewArgumentI32 creates an int32 type argument
//
// Parameters:
//   - value: int32 value
//
// Returns:
//   - TiArgument
//
// Example:
//
//	arg := taichi.NewArgumentI32(123)
func NewArgumentI32(value int32) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_I32,
	}
	*(*int32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	return arg
}

// NewArgumentF32 creates a float32 type argument
//
// Parameters:
//   - value: float32 value
//
// Returns:
//   - TiArgument
//
// Example:
//
//	arg := taichi.NewArgumentF32(456.0)
func NewArgumentF32(value float32) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_F32,
	}
	*(*float32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	return arg
}

// NewArgumentNdArray creates an NdArray type argument
//
// Parameters:
//   - ndarray: TiNdArray structure
//
// Returns:
//   - TiArgument
//
// Example:
//
//	ndarray := taichi.TiNdArray{
//	    Memory: memory,
//	    Shape: taichi.TiNdShape{
//	        DimCount: 2,
//	        Dims: [16]uint32{4, 4},
//	    },
//	    ElemType: taichi.TI_DATA_TYPE_F32,
//	}
//	arg := taichi.NewArgumentNdArray(ndarray)
func NewArgumentNdArray(ndarray TiNdArray) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_NDARRAY,
	}
	*(*TiNdArray)(unsafe.Pointer(&arg.Value.Data[0])) = ndarray
	return arg
}

// NewArgumentTexture creates a texture type argument
//
// Parameters:
//   - texture: TiTexture structure
//
// Returns:
//   - TiArgument
//
// Example:
//
//	texture := taichi.TiTexture{
//	    Image:   image,
//	    Sampler: sampler,
//	    Dimension: taichi.TI_IMAGE_DIMENSION_2D,
//	    Extent:  extent,
//	    Format:  taichi.TI_FORMAT_RGBA8,
//	}
//	arg := taichi.NewArgumentTexture(texture)
func NewArgumentTexture(texture TiTexture) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_TEXTURE,
	}
	*(*TiTexture)(unsafe.Pointer(&arg.Value.Data[0])) = texture
	return arg
}

// NewArgumentScalar creates a scalar type argument
//
// Parameters:
//   - scalar: TiScalar structure
//
// Returns:
//   - TiArgument
//
// Example:
//
//	scalar := taichi.TiScalar{
//	    Type: taichi.TI_DATA_TYPE_F64,
//	    Value: taichi.TiScalarValue{X64: math.Float64bits(3.14159)},
//	}
//	arg := taichi.NewArgumentScalar(scalar)
func NewArgumentScalar(scalar TiScalar) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_SCALAR,
	}
	*(*TiScalar)(unsafe.Pointer(&arg.Value.Data[0])) = scalar
	return arg
}

// NewArgumentTensor creates a tensor type argument
//
// Parameters:
//   - tensor: TiTensor structure
//
// Returns:
//   - TiArgument
//
// Example:
//
//	tensor := taichi.TiTensor{
//	    Type: taichi.TI_DATA_TYPE_F32,
//	    Contents: taichi.TiTensorValueWithLength{
//	        Length: 4,
//	        Data: taichi.TiTensorValue{},
//	    },
//	}
//	// Fill tensor data...
//	arg := taichi.NewArgumentTensor(tensor)
func NewArgumentTensor(tensor TiTensor) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_TENSOR,
	}
	*(*TiTensor)(unsafe.Pointer(&arg.Value.Data[0])) = tensor
	return arg
}

// NewArgumentTensorF32 creates a float32 tensor argument from a slice
//
// Parameters:
//   - data: []float32 slice
//
// Returns:
//   - TiArgument
//
// Example:
//
//	data := []float32{1.0, 2.0, 3.0, 4.0}
//	arg := taichi.NewArgumentTensorF32(data)
func NewArgumentTensorF32(data []float32) TiArgument {
	tensor := TiTensor{
		Type: TI_DATA_TYPE_F32,
		Contents: TiTensorValueWithLength{
			Length: uint32(len(data)),
		},
	}
	if len(data) > 0 {
		copy(tensor.Contents.Data.Data[:len(data)*4], tensorToBytes(data))
	}
	return NewArgumentTensor(tensor)
}

// NewArgumentTensorI32 creates an int32 tensor argument from a slice
//
// Parameters:
//   - data: []int32 slice
//
// Returns:
//   - TiArgument
//
// Example:
//
//	data := []int32{1, 2, 3, 4}
//	arg := taichi.NewArgumentTensorI32(data)
func NewArgumentTensorI32(data []int32) TiArgument {
	tensor := TiTensor{
		Type: TI_DATA_TYPE_I32,
		Contents: TiTensorValueWithLength{
			Length: uint32(len(data)),
		},
	}
	if len(data) > 0 {
		copy(tensor.Contents.Data.Data[:len(data)*4], tensorToBytes(data))
	}
	return NewArgumentTensor(tensor)
}

// NewArgumentTensorF64 creates a float64 tensor argument from a slice
//
// Parameters:
//   - data: []float64 slice
//
// Returns:
//   - TiArgument
//
// Example:
//
//	data := []float64{1.0, 2.0, 3.0, 4.0}
//	arg := taichi.NewArgumentTensorF64(data)
func NewArgumentTensorF64(data []float64) TiArgument {
	tensor := TiTensor{
		Type: TI_DATA_TYPE_F64,
		Contents: TiTensorValueWithLength{
			Length: uint32(len(data)),
		},
	}
	if len(data) > 0 {
		copy(tensor.Contents.Data.Data[:len(data)*8], tensorToBytes(data))
	}
	return NewArgumentTensor(tensor)
}

// NewArgumentTensorI64 creates an int64 tensor argument from a slice
//
// Parameters:
//   - data: []int64 slice
//
// Returns:
//   - TiArgument
//
// Example:
//
//	data := []int64{1, 2, 3, 4}
//	arg := taichi.NewArgumentTensorI64(data)
func NewArgumentTensorI64(data []int64) TiArgument {
	tensor := TiTensor{
		Type: TI_DATA_TYPE_I64,
		Contents: TiTensorValueWithLength{
			Length: uint32(len(data)),
		},
	}
	if len(data) > 0 {
		copy(tensor.Contents.Data.Data[:len(data)*8], tensorToBytes(data))
	}
	return NewArgumentTensor(tensor)
}

// ===== Named Argument Helper Functions =====

// NewNamedArgument creates a named argument
//
// Parameters:
//   - name: Argument name
//   - argument: Argument value
//
// Returns:
//   - TiNamedArgument
//
// Note: The name string will be converted to a C string (null-terminated),
// and the returned structure contains a pointer to this C string.
// The original name string must remain valid during the use of the named argument.
//
// Example:
//
//	arg := taichi.NewNamedArgument("my_param", taichi.NewArgumentI32(123))
func NewNamedArgument(name string, argument TiArgument) TiNamedArgument {
	cName := append([]byte(name), 0)
	return TiNamedArgument{
		Name:     &cName[0],
		Argument: argument,
	}
}

// NewNamedArgumentWithCString creates a named argument using a prepared C string
//
// Parameters:
//   - cName: C string pointer (null-terminated)
//   - argument: Argument value
//
// Returns:
//   - TiNamedArgument
//
// Example:
//
//	cName := append([]byte("my_param"), 0)
//	arg := taichi.NewNamedArgumentWithCString(&cName[0], taichi.NewArgumentI32(123))
func NewNamedArgumentWithCString(cName *byte, argument TiArgument) TiNamedArgument {
	return TiNamedArgument{
		Name:     cName,
		Argument: argument,
	}
}

// ===== Image Helper Functions =====

// NewTexture2D creates a 2D texture
//
// Parameters:
//   - image: Image handle
//   - sampler: Sampler handle (can be TI_NULL_HANDLE to use default sampler)
//   - width, height: Texture dimensions
//   - format: Texture format
//
// Returns:
//   - TiTexture
//
// Example:
//
//	texture := taichi.NewTexture2D(image, sampler, 1024, 1024, taichi.TI_FORMAT_RGBA8)
func NewTexture2D(image TiImage, sampler TiSampler, width, height uint32, format TiFormat) TiTexture {
	return TiTexture{
		Image:     image,
		Sampler:   sampler,
		Dimension: TI_IMAGE_DIMENSION_2D,
		Extent: TiImageExtent{
			Width:           width,
			Height:          height,
			Depth:           1,
			ArrayLayerCount: 1,
		},
		Format: format,
	}
}

// ===== Memory Slice Helper Functions =====

// NewMemorySlice creates a memory slice
//
// Parameters:
//   - memory: Memory handle
//   - offset: Offset (bytes)
//   - size: Size (bytes)
//
// Returns:
//   - TiMemorySlice
//
// Example:
//
//	slice := taichi.NewMemorySlice(memory, 0, 1024)
func NewMemorySlice(memory TiMemory, offset, size uint64) TiMemorySlice {
	return TiMemorySlice{
		Memory: memory,
		Offset: offset,
		Size:   size,
	}
}

// NewFullMemorySlice creates a full memory slice (from start to end)
//
// Parameters:
//   - memory: Memory handle
//   - size: Full size (bytes)
//
// Returns:
//   - TiMemorySlice
//
// Example:
//
//	slice := taichi.NewFullMemorySlice(memory, 4096)
func NewFullMemorySlice(memory TiMemory, size uint64) TiMemorySlice {
	return TiMemorySlice{
		Memory: memory,
		Offset: 0,
		Size:   size,
	}
}
