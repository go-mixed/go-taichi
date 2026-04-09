package taichi

import (
	"fmt"
	"unsafe"

	"github.com/go-mixed/go-taichi/f16"
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

// GetOffset gets the element indices at the given shape indices
// For 2D array with shape [width, height] and elementShape [4] (vec4):
//   - base = (x*height + y) * 4
//   - returns [base, base+1, base+2, base+3]
//
// If shapeIndices has fewer elements, pad with zeros; if more, ignore extras
func (arr *NdArray) GetOffset(shapeIndices ...int) (offset int, elementSize int) {
	// Pad shapeIndices with zeros if shorter than shape
	indices := make([]int, len(arr.shape))
	copy(indices, shapeIndices)

	// Calculate base index in row-major order
	base := 0
	for i, idx := range indices {
		// stride[i] = product of shape[i+1:]
		stride := 1
		for j := i + 1; j < len(arr.shape); j++ {
			stride *= int(arr.shape[j])
		}
		base += idx * stride
	}

	// Multiply by element count (e.g., 4 for vec4)
	elemCount := 1
	for _, dim := range arr.elemShape {
		elemCount *= int(dim)
	}

	base *= elemCount

	return base, elemCount
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

func (arr *NdArray) MapFloat16(f func(data []f16.Float16) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsFloat16()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapFloat32(f func(data []float32) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsFloat32()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapFloat64(f func(data []float64) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsFloat64()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapInt8(f func(data []int8) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsInt8()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapUint8(f func(data []uint8) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsUInt8()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapInt16(f func(data []int16) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsInt16()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapUint16(f func(data []uint16) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsUInt16()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapInt32(f func(data []int32) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsInt32()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapUint32(f func(data []uint32) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsUint32()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapInt64(f func(data []int64) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsInt64()
		return f(data)
	}, arr)
}

func (arr *NdArray) MapUint64(f func(data []uint64) error) error {
	return MapNdArray(func(datas ...NdArrayPtr) error {
		data := datas[0].AsUint64()
		return f(data)
	}, arr)
}

type NdArrayPtr struct {
	ptr unsafe.Pointer
	arr *NdArray
}

func (p NdArrayPtr) AsInt64() []int64 {
	return unsafe.Slice((*int64)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsUInt8() []byte {
	return unsafe.Slice((*uint8)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsInt8() []int8 {
	return unsafe.Slice((*int8)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsInt16() []int16 {
	return unsafe.Slice((*int16)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsUInt16() []uint16 {
	return unsafe.Slice((*uint16)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsFloat32() []float32 {
	return unsafe.Slice((*float32)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsFloat64() []float64 {
	return unsafe.Slice((*float64)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsUint32() []uint32 {
	return unsafe.Slice((*uint32)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsUint64() []uint64 {
	return unsafe.Slice((*uint64)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsInt32() []int32 {
	return unsafe.Slice((*int32)(p.ptr), p.arr.TotalElements())
}

func (p NdArrayPtr) AsFloat16() []f16.Float16 {
	return unsafe.Slice((*f16.Float16)(p.ptr), p.arr.TotalElements())
}

// MapNdArray executes fn with multiple NdArrays as NdArrayPtr
// and then, you can AsFloat32(), ... to get the data
func MapNdArray(fn func(arrays ...NdArrayPtr) error, arrays ...*NdArray) error {
	if len(arrays) == 0 {
		return fmt.Errorf("no arrays provided")
	}

	// Extract memories
	memories := make([]*Memory, len(arrays))
	for i, arr := range arrays {
		memories[i] = arr.Memory
	}

	return MapMemory(func(ptrs ...unsafe.Pointer) error {
		args := make([]NdArrayPtr, len(arrays))
		for i, arr := range arrays {
			args[i] = NdArrayPtr{ptr: ptrs[i], arr: arr}
		}
		return fn(args...)
	}, memories...)
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
