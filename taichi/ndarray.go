package taichi

import (
	"fmt"
)

// NdArray N-dimensional array abstraction
type NdArray struct {
	*Memory
	shape     NdShape
	elemType  DataType
	elemSize  int
	elemShape NdShape // Element shape (e.g., [4] for vec4), nil for scalar
}

type NdShape []uint32

func Shape(vals ...uint32) NdShape {
	return vals
}

// NewNdArray1D creates a 1D array
func NewNdArray1D(runtime *Runtime, length uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, NdShape{length}, elemType, nil)
}

// NewNdArray2D creates a 2D array
func NewNdArray2D(runtime *Runtime, dim0, dim1 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, NdShape{dim0, dim1}, elemType, nil)
}

// NewNdArray3D creates a 3D array
func NewNdArray3D(runtime *Runtime, dim0, dim1, dim2 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, NdShape{dim0, dim1, dim2}, elemType, nil)
}

// NewNdArray2DWithElemShape creates a 2D array with element shape (e.g., vec4)
// For textures with element_shape=(4,), dtype=f32, ndim=2
func NewNdArray2DWithElemShape(runtime *Runtime, dim0, dim1 uint32, elemShape NdShape, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, NdShape{dim0, dim1}, elemType, elemShape)
}

// NewNdArray creates an N-dimensional array
// elemShape is optional (nil for scalar elements, e.g., [4] for vec4)
func NewNdArray(runtime *Runtime, shape NdShape, elemType DataType, elemShape NdShape) (*NdArray, error) {
	if len(shape) == 0 || len(shape) > 16 {
		return nil, fmt.Errorf("number of dimensions must be between 1 and 16")
	}

	// Calculate total element count and memory size
	elemSize := getElemSize(elemType)
	totalElems := uint64(1)
	for _, dim := range shape {
		if dim == 0 {
			return nil, fmt.Errorf("dimension size cannot be 0")
		}
		totalElems *= uint64(dim)
	}

	// If elemShape is provided, multiply by element count
	elemCount := uint64(1)
	if elemShape != nil {
		for _, dim := range elemShape {
			elemCount *= uint64(dim)
		}
	}
	totalElems *= elemCount
	size := totalElems * uint64(elemSize)

	// Allocate memory
	memory, err := NewMemory(runtime, size)
	if err != nil {
		return nil, err
	}

	return &NdArray{
		Memory:    memory,
		shape:     shape,
		elemType:  elemType,
		elemSize:  elemSize,
		elemShape: elemShape,
	}, nil
}

// Shape gets the array shape
func (arr *NdArray) Shape() NdShape {
	return arr.shape
}

// Ndim gets the number of dimensions
func (arr *NdArray) Ndim() int {
	return len(arr.shape)
}

// TotalElements gets the total number of elements (including elemShape)
func (arr *NdArray) TotalElements() uint64 {
	total := uint64(1)
	for _, dim := range arr.shape {
		total *= uint64(dim)
	}
	// Multiply by elemShape elements if present
	if arr.elemShape != nil {
		for _, dim := range arr.elemShape {
			total *= uint64(dim)
		}
	}
	return total
}

// ElemType gets the element type
func (arr *NdArray) ElemType() DataType {
	return arr.elemType
}

// ElemSize gets the element size (bytes)
func (arr *NdArray) ElemSize() int {
	return arr.elemSize
}

// NdArrayAs executes fn with multiple NdArrays as T slices
// All arrays must be float32 type. Uses MapMemory for thread safety.
//func NdArrayAs[T any](fn func(float32Arrays ...T) error, arrays ...*NdArray) error {
//	if len(arrays) == 0 {
//		return fmt.Errorf("no arrays provided")
//	}
//
//	// Validate all arrays are float32 type
//	for _, arr := range arrays {
//		if arr.elemType != DataTypeF32 {
//			return fmt.Errorf("array type is not float32")
//		}
//	}
//
//	// Extract memories
//	memories := make([]*Memory, len(arrays))
//	for i, arr := range arrays {
//		memories[i] = arr.Memory
//	}
//
//	return MapMemory(func(ptrs ...unsafe.Pointer) error {
//		args := make([]T, len(arrays))
//		for i, arr := range arrays {
//			args[i] = unsafe.Slice((*T)(ptrs[i]), arr.TotalElements())
//		}
//		return fn(args...)
//	}, memories...)
//}

// getElemSize gets the element size
func getElemSize(elemType DataType) int {
	switch elemType {
	case DataTypeI8, DataTypeU8:
		return 1
	case DataTypeI16, DataTypeU16, DataTypeF16:
		return 2
	case DataTypeI32, DataTypeU32, DataTypeF32:
		return 4
	case DataTypeI64, DataTypeU64, DataTypeF64:
		return 8
	default:
		return 4 // Default 4 bytes
	}
}
