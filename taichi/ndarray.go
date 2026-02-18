package taichi

import (
	"fmt"
	"unsafe"
)

// NdArray N维数组抽象
type NdArray struct {
	*Memory
	shape    []uint32
	elemType DataType
	elemSize int
}

// NewNdArray1D 创建1D数组
func NewNdArray1D(runtime *Runtime, length uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{length}, elemType)
}

// NewNdArray2D 创建2D数组
func NewNdArray2D(runtime *Runtime, dim0, dim1 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{dim0, dim1}, elemType)
}

// NewNdArray3D 创建3D数组
func NewNdArray3D(runtime *Runtime, dim0, dim1, dim2 uint32, elemType DataType) (*NdArray, error) {
	return NewNdArray(runtime, []uint32{dim0, dim1, dim2}, elemType)
}

// NewNdArray 创建N维数组
func NewNdArray(runtime *Runtime, shape []uint32, elemType DataType) (*NdArray, error) {
	if len(shape) == 0 || len(shape) > 16 {
		return nil, fmt.Errorf("维度数量必须在1-16之间")
	}

	// 计算总元素数和内存大小
	elemSize := getElemSize(elemType)
	totalElems := uint64(1)
	for _, dim := range shape {
		if dim == 0 {
			return nil, fmt.Errorf("维度大小不能为0")
		}
		totalElems *= uint64(dim)
	}
	size := totalElems * uint64(elemSize)

	// 分配内存
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

// Shape 获取数组形状
func (arr *NdArray) Shape() []uint32 {
	return arr.shape
}

// Ndim 获取维度数
func (arr *NdArray) Ndim() int {
	return len(arr.shape)
}

// TotalElements 获取总元素数
func (arr *NdArray) TotalElements() uint64 {
	total := uint64(1)
	for _, dim := range arr.shape {
		total *= uint64(dim)
	}
	return total
}

// ElemType 获取元素类型
func (arr *NdArray) ElemType() DataType {
	return arr.elemType
}

// ElemSize 获取元素大小（字节）
func (arr *NdArray) ElemSize() int {
	return arr.elemSize
}

// AsSliceFloat32 将数组映射为float32切片（仅限float32类型）
func (arr *NdArray) AsSliceFloat32() ([]float32, error) {
	if arr.elemType != DataTypeF32 {
		return nil, fmt.Errorf("数组类型不是float32")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*float32)(ptr), length), nil
}

// AsSliceInt32 将数组映射为int32切片（仅限int32类型）
func (arr *NdArray) AsSliceInt32() ([]int32, error) {
	if arr.elemType != DataTypeI32 {
		return nil, fmt.Errorf("数组类型不是int32")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*int32)(ptr), length), nil
}

// AsSliceUint8 将数组映射为uint8切片（仅限uint8类型）
func (arr *NdArray) AsSliceUint8() ([]uint8, error) {
	if arr.elemType != DataTypeU8 {
		return nil, fmt.Errorf("数组类型不是uint8")
	}

	ptr, err := arr.Map()
	if err != nil {
		return nil, err
	}

	length := arr.TotalElements()
	return unsafe.Slice((*uint8)(ptr), length), nil
}

// Fill 填充数组（float32）
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

// FillInt32 填充数组（int32）
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

// getElemSize 获取元素大小
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
		return 4 // 默认4字节
	}
}
