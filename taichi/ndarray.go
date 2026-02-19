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

// AsSliceFloat32 maps the array as a float32 slice (float32 type only)
func (arr *NdArray) AsSliceFloat32() ([]float32, error) {
	if arr.elemType != DataTypeF32 {
		return nil, fmt.Errorf("array type is not float32")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*float32)(ptr), length), nil
}

// AsSliceInt32 maps the array as an int32 slice (int32 type only)
func (arr *NdArray) AsSliceInt32() ([]int32, error) {
	if arr.elemType != DataTypeI32 {
		return nil, fmt.Errorf("array type is not int32")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*int32)(ptr), length), nil
}

// AsSliceUint8 maps the array as a uint8 slice (uint8 type only)
func (arr *NdArray) AsSliceUint8() ([]uint8, error) {
	if arr.elemType != DataTypeU8 {
		return nil, fmt.Errorf("array type is not uint8")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*uint8)(ptr), length), nil
}

// Fill fills the array (float32)
func (arr *NdArray) Fill(value float32) error {
	data, err := arr.AsSliceFloat32()
	if err != nil {
		return err
	}

	for i := range data {
		data[i] = value
	}

	return nil
}

// FillInt32 fills the array (int32)
func (arr *NdArray) FillInt32(value int32) error {
	data, err := arr.AsSliceInt32()
	if err != nil {
		return err
	}

	for i := range data {
		data[i] = value
	}

	return nil
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
