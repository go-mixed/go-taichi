package taichi

import (
	"fmt"
	"unsafe"
)

// NdArray N-dimensional array abstraction
type NdArray struct {
	*Memory
	shape    []uint32
	elemType DataType
	elemSize int
}

// NewNdArray1D creates a 1D array
func NewNdArray1D(runtime *Runtime, length uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{length}, elemType)
}

// NewNdArray2D creates a 2D array
func NewNdArray2D(runtime *Runtime, dim0, dim1 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{dim0, dim1}, elemType)
}

// NewNdArray3D creates a 3D array
func NewNdArray3D(runtime *Runtime, dim0, dim1, dim2 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{dim0, dim1, dim2}, elemType)
}

// NewNdArray creates an N-dimensional array
func NewNdArray(runtime *Runtime, shape []uint32, elemType DataType) (*NdArray, error) {
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
	size := totalElems * uint64(elemSize)

	// Allocate memory
	memory, err := NewMemory(runtime, size)
	if err != nil {
		return nil, err
	}

	return &NdArray{
		Memory:   memory,
		shape:    shape,
		elemType: elemType,
		elemSize: elemSize,
	}, nil
}

// Shape gets the array shape
func (arr *NdArray) Shape() []uint32 {
	return arr.shape
}

// Ndim gets the number of dimensions
func (arr *NdArray) Ndim() int {
	return len(arr.shape)
}

// TotalElements gets the total number of elements
func (arr *NdArray) TotalElements() uint64 {
	total := uint64(1)
	for _, dim := range arr.shape {
		total *= uint64(dim)
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

// AsPtr invokes fn with the array data (internal use, uses MapMemory for thread safety)
func (arr *NdArray) AsPtr(fn func(ptr unsafe.Pointer) error) error {
	return NdArrayAsPtr(func(ptrs ...unsafe.Pointer) error {
		return fn(ptrs[0])
	}, arr)
}

// NdArrayAsPtr executes fn with multiple NdArrays as unsafe.Pointer
// All arrays must be unsafe.Pointer type. Uses MapMemory for thread safety.
func NdArrayAsPtr(fn func(ptr ...unsafe.Pointer) error, arrays ...*NdArray) error {
	// Extract memories
	memories := make([]*Memory, len(arrays))
	for i, arr := range arrays {
		memories[i] = arr.Memory
	}

	return MapMemory(func(ptrs ...unsafe.Pointer) error {
		return fn(ptrs...)
	}, memories...)
}

// NdArrayAsFloat32 executes fn with multiple NdArrays as float32 slices
// All arrays must be float32 type. Uses MapMemory for thread safety.
func NdArrayAsFloat32(fn func(float32Arrays ...[]float32) error, arrays ...*NdArray) error {
	// Validate all arrays are float32 type
	for _, arr := range arrays {
		if arr.elemType != DataTypeF32 {
			return fmt.Errorf("array type is not float32")
		}
	}

	// Extract memories
	memories := make([]*Memory, len(arrays))
	for i, arr := range arrays {
		memories[i] = arr.Memory
	}

	return MapMemory(func(ptrs ...unsafe.Pointer) error {
		float32Arrays := make([][]float32, len(arrays))
		for i, arr := range arrays {
			float32Arrays[i] = unsafe.Slice((*float32)(ptrs[i]), arr.TotalElements())
		}
		return fn(float32Arrays...)
	}, memories...)
}

// WithFloat32 executes fn with the array data as a float32 slice (float32 type only)
// The slice is only valid during the callback execution.
func (arr *NdArray) WithFloat32(fn func([]float32) error) error {
	return NdArrayAsFloat32(func(arrays ...[]float32) error {
		return fn(arrays[0])
	}, arr)
}

// NdArrayAsInt32 executes fn with multiple NdArrays as int32 slices
// All arrays must be int32 type. Uses MapMemory for thread safety.
func NdArrayAsInt32(fn func(int32Arrays ...[]int32) error, arrays ...*NdArray) error {
	// Validate all arrays are int32 type
	for _, arr := range arrays {
		if arr.elemType != DataTypeI32 {
			return fmt.Errorf("array type is not int32")
		}
	}

	// Extract memories
	memories := make([]*Memory, len(arrays))
	for i, arr := range arrays {
		memories[i] = arr.Memory
	}

	return MapMemory(func(ptrs ...unsafe.Pointer) error {
		int32Arrays := make([][]int32, len(arrays))
		for i, arr := range arrays {
			int32Arrays[i] = unsafe.Slice((*int32)(ptrs[i]), arr.TotalElements())
		}
		return fn(int32Arrays...)
	}, memories...)
}

// WithInt32 executes fn with the array data as an int32 slice (int32 type only)
// The slice is only valid during the callback execution.
func (arr *NdArray) WithInt32(fn func([]int32) error) error {
	return NdArrayAsInt32(func(arrays ...[]int32) error {
		return fn(arrays[0])
	}, arr)
}

// NdArrayAsUint8 executes fn with multiple NdArrays as uint8 slices
// All arrays must be uint8 type. Uses MapMemory for thread safety.
func NdArrayAsUint8(fn func(uint8Arrays ...[]uint8) error, arrays ...*NdArray) error {
	// Validate all arrays are uint8 type
	for _, arr := range arrays {
		if arr.elemType != DataTypeU8 {
			return fmt.Errorf("array type is not uint8")
		}
	}

	// Extract memories
	memories := make([]*Memory, len(arrays))
	for i, arr := range arrays {
		memories[i] = arr.Memory
	}

	return MapMemory(func(ptrs ...unsafe.Pointer) error {
		uint8Arrays := make([][]uint8, len(arrays))
		for i, arr := range arrays {
			uint8Arrays[i] = unsafe.Slice((*uint8)(ptrs[i]), arr.TotalElements())
		}
		return fn(uint8Arrays...)
	}, memories...)
}

// WithUint8 executes fn with the array data as a uint8 slice (uint8 type only)
// The slice is only valid during the callback execution.
func (arr *NdArray) WithUint8(fn func([]uint8) error) error {
	return NdArrayAsUint8(func(arrays ...[]uint8) error {
		return fn(arrays[0])
	}, arr)
}

// Fill fills the array (float32)
func (arr *NdArray) Fill(value float32) error {
	return arr.WithFloat32(func(data []float32) error {
		for i := range data {
			data[i] = value
		}
		return nil
	})
}

// FillInt32 fills the array (int32)
func (arr *NdArray) FillInt32(value int32) error {
	return arr.WithInt32(func(data []int32) error {
		for i := range data {
			data[i] = value
		}
		return nil
	})
}

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
