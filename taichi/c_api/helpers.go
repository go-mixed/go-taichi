package c_api

import "unsafe"

// ===== 参数构造辅助函数 =====

// NewArgumentI32 创建int32类型的参数
//
// 参数:
//   - value: int32值
//
// 返回:
//   - TiArgument
//
// 示例:
//
//	arg := taichi.NewArgumentI32(123)
func NewArgumentI32(value int32) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_I32,
	}
	*(*int32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	return arg
}

// NewArgumentF32 创建float32类型的参数
//
// 参数:
//   - value: float32值
//
// 返回:
//   - TiArgument
//
// 示例:
//
//	arg := taichi.NewArgumentF32(456.0)
func NewArgumentF32(value float32) TiArgument {
	arg := TiArgument{
		Type: TI_ARGUMENT_TYPE_F32,
	}
	*(*float32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	return arg
}

// NewArgumentNdArray 创建NdArray类型的参数
//
// 参数:
//   - ndarray: TiNdArray结构
//
// 返回:
//   - TiArgument
//
// 示例:
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

// NewArgumentTexture 创建纹理类型的参数
//
// 参数:
//   - texture: TiTexture结构
//
// 返回:
//   - TiArgument
//
// 示例:
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

// NewArgumentScalar 创建标量类型的参数
//
// 参数:
//   - scalar: TiScalar结构
//
// 返回:
//   - TiArgument
//
// 示例:
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

// ===== 命名参数辅助函数 =====

// NewNamedArgument 创建命名参数
//
// 参数:
//   - name: 参数名称
//   - argument: 参数值
//
// 返回:
//   - TiNamedArgument
//
// 注意:name字符串会被转换为C字符串(null结尾),
// 返回的结构体中包含指向这个C字符串的指针。
// 在使用命名参数期间,原始的name string必须保持有效。
//
// 示例:
//
//	arg := taichi.NewNamedArgument("my_param", taichi.NewArgumentI32(123))
func NewNamedArgument(name string, argument TiArgument) TiNamedArgument {
	cName := append([]byte(name), 0)
	return TiNamedArgument{
		Name:     &cName[0],
		Argument: argument,
	}
}

// NewNamedArgumentWithCString 使用已经准备好的C字符串创建命名参数
//
// 参数:
//   - cName: C字符串指针(null结尾)
//   - argument: 参数值
//
// 返回:
//   - TiNamedArgument
//
// 示例:
//
//	cName := append([]byte("my_param"), 0)
//	arg := taichi.NewNamedArgumentWithCString(&cName[0], taichi.NewArgumentI32(123))
func NewNamedArgumentWithCString(cName *byte, argument TiArgument) TiNamedArgument {
	return TiNamedArgument{
		Name:     cName,
		Argument: argument,
	}
}

// ===== NdArray辅助函数 =====

// NewNdArray1D 创建1维NdArray
//
// 参数:
//   - memory: 内存句柄
//   - length: 数组长度
//   - elemType: 元素类型
//
// 返回:
//   - TiNdArray
//
// 示例:
//
//	ndarray := taichi.NewNdArray1D(memory, 256, taichi.TI_DATA_TYPE_F32)
func NewNdArray1D(memory TiMemory, length uint32, elemType TiDataType) TiNdArray {
	return TiNdArray{
		Memory: memory,
		Shape: TiNdShape{
			DimCount: 1,
			Dims:     [16]uint32{length},
		},
		ElemType: elemType,
	}
}

// NewNdArray2D 创建2维NdArray
//
// 参数:
//   - memory: 内存句柄
//   - rows: 行数
//   - cols: 列数
//   - elemType: 元素类型
//
// 返回:
//   - TiNdArray
//
// 示例:
//
//	ndarray := taichi.NewNdArray2D(memory, 16, 16, taichi.TI_DATA_TYPE_F32)
func NewNdArray2D(memory TiMemory, rows, cols uint32, elemType TiDataType) TiNdArray {
	return TiNdArray{
		Memory: memory,
		Shape: TiNdShape{
			DimCount: 2,
			Dims:     [16]uint32{rows, cols},
		},
		ElemType: elemType,
	}
}

// NewNdArray3D 创建3维NdArray
//
// 参数:
//   - memory: 内存句柄
//   - dim0, dim1, dim2: 三个维度的大小
//   - elemType: 元素类型
//
// 返回:
//   - TiNdArray
//
// 示例:
//
//	ndarray := taichi.NewNdArray3D(memory, 8, 8, 8, taichi.TI_DATA_TYPE_F32)
func NewNdArray3D(memory TiMemory, dim0, dim1, dim2 uint32, elemType TiDataType) TiNdArray {
	return TiNdArray{
		Memory: memory,
		Shape: TiNdShape{
			DimCount: 3,
			Dims:     [16]uint32{dim0, dim1, dim2},
		},
		ElemType: elemType,
	}
}

// ===== 图像辅助函数 =====

// NewTexture2D 创建2D纹理
//
// 参数:
//   - image: 图像句柄
//   - sampler: 采样器句柄(可以是TI_NULL_HANDLE使用默认采样器)
//   - width, height: 纹理尺寸
//   - format: 纹理格式
//
// 返回:
//   - TiTexture
//
// 示例:
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

// ===== 内存切片辅助函数 =====

// NewMemorySlice 创建内存切片
//
// 参数:
//   - memory: 内存句柄
//   - offset: 偏移量(字节)
//   - size: 大小(字节)
//
// 返回:
//   - TiMemorySlice
//
// 示例:
//
//	slice := taichi.NewMemorySlice(memory, 0, 1024)
func NewMemorySlice(memory TiMemory, offset, size uint64) TiMemorySlice {
	return TiMemorySlice{
		Memory: memory,
		Offset: offset,
		Size:   size,
	}
}

// NewFullMemorySlice 创建完整内存切片(从头到尾)
//
// 参数:
//   - memory: 内存句柄
//   - size: 完整大小(字节)
//
// 返回:
//   - TiMemorySlice
//
// 示例:
//
//	slice := taichi.NewFullMemorySlice(memory, 4096)
func NewFullMemorySlice(memory TiMemory, size uint64) TiMemorySlice {
	return TiMemorySlice{
		Memory: memory,
		Offset: 0,
		Size:   size,
	}
}
