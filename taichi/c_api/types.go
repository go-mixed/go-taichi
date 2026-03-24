package c_api

// ===== Go 与 C 结构体对齐 =====
//
// ⚠️ 重要：本文件中的结构体定义必须与 C 结构体的内存布局完全一致！
//
// 在使用 purego 进行 Go-C FFI 调用时，必须确保：
// 1. 结构体总大小匹配（使用 unsafe.Sizeof 验证）
// 2. 每个字段的偏移量匹配（使用 unsafe.Offsetof 验证）
// 3. Union 使用最大成员的大小
// 4. 必要时添加显式 padding（使用 _ [N]byte）
//
// 详见：STRUCT_ALIGNMENT.md
// 测试工具：examples/struct_alignment.go
//
// ===== 基础类型定义 =====

// TiBool 布尔类型
// taichi_core.h:246
//
// 可以是TI_TRUE或TI_FALSE,使用其他值会导致未定义行为。
type TiBool uint32

// TiFlags 标志位字段
// taichi_core.h:266
//
// 用于表示32个正交标志的位字段。
type TiFlags uint32

// ===== 句柄类型 =====

// TiRuntime Taichi运行时句柄
// taichi_core.h:280
//
// 表示逻辑后端及其内部动态状态的实例。
// 用户负责同步对TiRuntime的任何使用,不能在同一线程中操作多个TiRuntime。
type TiRuntime uintptr

// TiAotModule AOT模块句柄
// taichi_core.h:286
//
// 提前编译的Taichi模块,包含kernel和compute graph的集合。
type TiAotModule uintptr

// TiMemory 内存句柄
// taichi_core.h:291
//
// 设备内存的连续分配。
type TiMemory uintptr

// TiImage 图像句柄
// taichi_core.h:296
//
// 设备图像的连续分配。
type TiImage uintptr

// TiSampler 采样器句柄
// taichi_core.h:303
//
// 图像采样器。TI_NULL_HANDLE表示运行时提供的默认采样器。
type TiSampler uintptr

// TiKernel 内核句柄
// taichi_core.h:308
//
// 可在目标设备上执行的Taichi kernel。
type TiKernel uintptr

// TiComputeGraph 计算图句柄
// taichi_core.h:314
//
// 按预定义顺序在目标设备上启动的Taichi kernel集合。
type TiComputeGraph uintptr

// ===== 常量定义 =====

const (
	// TI_FALSE 表示条件或谓词不满足;语句无效
	TI_FALSE TiBool = 0

	// TI_TRUE 表示条件或谓词满足;语句有效
	TI_TRUE TiBool = 1

	// TI_NULL_HANDLE 无效句柄的哨兵值,永远不会从有效的Taichi C-API调用中产生
	TI_NULL_HANDLE = 0

	// TI_C_API_VERSION Taichi C-API版本号
	TI_C_API_VERSION = 1007000
)

// ===== 错误码 =====

// TiError Taichi C-API报告的错误
// taichi_core.h:319-353
type TiError int32

const (
	// TI_ERROR_SUCCESS Taichi C-API调用成功完成
	TI_ERROR_SUCCESS TiError = 0

	// TI_ERROR_NOT_SUPPORTED 调用的API或参数组合不被Taichi C-API支持
	TI_ERROR_NOT_SUPPORTED TiError = -1

	// TI_ERROR_CORRUPTED_DATA 提供的数据已损坏
	TI_ERROR_CORRUPTED_DATA TiError = -2

	// TI_ERROR_NAME_NOT_FOUND 提供的名称不引用任何现有项
	TI_ERROR_NAME_NOT_FOUND TiError = -3

	// TI_ERROR_INVALID_ARGUMENT 一个或多个函数参数违反了C-API文档中指定的约束
	TI_ERROR_INVALID_ARGUMENT TiError = -4

	// TI_ERROR_ARGUMENT_NULL 一个或多个按引用(指针)传递的函数参数指向null
	TI_ERROR_ARGUMENT_NULL TiError = -5

	// TI_ERROR_ARGUMENT_OUT_OF_RANGE 一个或多个函数参数超出可接受范围
	TI_ERROR_ARGUMENT_OUT_OF_RANGE TiError = -6

	// TI_ERROR_ARGUMENT_NOT_FOUND 缺少一个或多个kernel参数
	TI_ERROR_ARGUMENT_NOT_FOUND TiError = -7

	// TI_ERROR_INVALID_INTEROP 当前架构上不可能进行预期的互操作
	TI_ERROR_INVALID_INTEROP TiError = -8

	// TI_ERROR_INVALID_STATE Taichi C-API进入不可恢复的无效状态
	TI_ERROR_INVALID_STATE TiError = -9

	// TI_ERROR_INCOMPATIBLE_MODULE AOT模块与当前运行时不兼容
	TI_ERROR_INCOMPATIBLE_MODULE TiError = -10

	// TI_ERROR_OUT_OF_MEMORY 内存不足
	TI_ERROR_OUT_OF_MEMORY TiError = -11
)

// ===== 架构类型 =====

// TiArch 后端架构类型
// taichi_core.h:358-375
type TiArch uint32

const (
	// TI_ARCH_RESERVED 保留值
	TI_ARCH_RESERVED TiArch = 0

	// TI_ARCH_VULKAN Vulkan GPU后端
	TI_ARCH_VULKAN TiArch = 1

	// TI_ARCH_METAL Metal GPU后端
	TI_ARCH_METAL TiArch = 2

	// TI_ARCH_CUDA NVIDIA CUDA GPU后端
	TI_ARCH_CUDA TiArch = 3

	// TI_ARCH_X64 x64原生CPU后端
	TI_ARCH_X64 TiArch = 4

	// TI_ARCH_ARM64 ARM64原生CPU后端
	TI_ARCH_ARM64 TiArch = 5

	// TI_ARCH_OPENGL OpenGL GPU后端
	TI_ARCH_OPENGL TiArch = 6

	// TI_ARCH_GLES OpenGL ES GPU后端
	TI_ARCH_GLES TiArch = 7
)

// ===== 数据类型 =====

// TiDataType 基本(原始)数据类型
// taichi_core.h:423-450
type TiDataType uint32

const (
	// TI_DATA_TYPE_F16 16位IEEE 754半精度浮点数
	TI_DATA_TYPE_F16 TiDataType = 0

	// TI_DATA_TYPE_F32 32位IEEE 754单精度浮点数
	TI_DATA_TYPE_F32 TiDataType = 1

	// TI_DATA_TYPE_F64 64位IEEE 754双精度浮点数
	TI_DATA_TYPE_F64 TiDataType = 2

	// TI_DATA_TYPE_I8 8位补码有符号整数
	TI_DATA_TYPE_I8 TiDataType = 3

	// TI_DATA_TYPE_I16 16位补码有符号整数
	TI_DATA_TYPE_I16 TiDataType = 4

	// TI_DATA_TYPE_I32 32位补码有符号整数
	TI_DATA_TYPE_I32 TiDataType = 5

	// TI_DATA_TYPE_I64 64位补码有符号整数
	TI_DATA_TYPE_I64 TiDataType = 6

	// TI_DATA_TYPE_U1 1位无符号整数
	TI_DATA_TYPE_U1 TiDataType = 7

	// TI_DATA_TYPE_U8 8位无符号整数
	TI_DATA_TYPE_U8 TiDataType = 8

	// TI_DATA_TYPE_U16 16位无符号整数
	TI_DATA_TYPE_U16 TiDataType = 9

	// TI_DATA_TYPE_U32 32位无符号整数
	TI_DATA_TYPE_U32 TiDataType = 10

	// TI_DATA_TYPE_U64 64位无符号整数
	TI_DATA_TYPE_U64 TiDataType = 11

	// TI_DATA_TYPE_GEN 通用类型
	TI_DATA_TYPE_GEN TiDataType = 12

	// TI_DATA_TYPE_UNKNOWN 未知类型
	TI_DATA_TYPE_UNKNOWN TiDataType = 13
)

// ===== 参数类型 =====

// TiArgumentType kernel和compute graph参数类型
// taichi_core.h:455-469
type TiArgumentType uint32

const (
	// TI_ARGUMENT_TYPE_I32 32位补码有符号整数
	TI_ARGUMENT_TYPE_I32 TiArgumentType = 0

	// TI_ARGUMENT_TYPE_F32 32位IEEE 754单精度浮点数
	TI_ARGUMENT_TYPE_F32 TiArgumentType = 1

	// TI_ARGUMENT_TYPE_NDARRAY 围绕handle.memory包装的ND数组
	TI_ARGUMENT_TYPE_NDARRAY TiArgumentType = 2

	// TI_ARGUMENT_TYPE_TEXTURE 围绕handle.image包装的纹理
	TI_ARGUMENT_TYPE_TEXTURE TiArgumentType = 3

	// TI_ARGUMENT_TYPE_SCALAR 类型化标量
	TI_ARGUMENT_TYPE_SCALAR TiArgumentType = 4

	// TI_ARGUMENT_TYPE_TENSOR 类型化张量
	TI_ARGUMENT_TYPE_TENSOR TiArgumentType = 5
)

// ===== 图像相关类型 =====

// TiImageDimension 图像维度
// taichi_core.h:562-577
type TiImageDimension uint32

const (
	TI_IMAGE_DIMENSION_1D   TiImageDimension = 0
	TI_IMAGE_DIMENSION_2D   TiImageDimension = 1
	TI_IMAGE_DIMENSION_3D   TiImageDimension = 2
	TI_IMAGE_DIMENSION_CUBE TiImageDimension = 3
)

// TiImageLayout 图像布局
// taichi_core.h:580-605
type TiImageLayout uint32

const (
	TI_IMAGE_LAYOUT_UNDEFINED         TiImageLayout = 0
	TI_IMAGE_LAYOUT_SHADER_READ       TiImageLayout = 1
	TI_IMAGE_LAYOUT_SHADER_WRITE      TiImageLayout = 2
	TI_IMAGE_LAYOUT_SHADER_READ_WRITE TiImageLayout = 3
	TI_IMAGE_LAYOUT_COLOR_ATTACHMENT  TiImageLayout = 4
	TI_IMAGE_LAYOUT_DEPTH_ATTACHMENT  TiImageLayout = 5
	TI_IMAGE_LAYOUT_TRANSFER_DST      TiImageLayout = 6
	TI_IMAGE_LAYOUT_TRANSFER_SRC      TiImageLayout = 7
	TI_IMAGE_LAYOUT_PRESENT_SRC       TiImageLayout = 8
)

// TiFormat 纹理格式
// taichi_core.h:611-657
type TiFormat uint32

const (
	TI_FORMAT_UNKNOWN         TiFormat = 0
	TI_FORMAT_R8              TiFormat = 1
	TI_FORMAT_RG8             TiFormat = 2
	TI_FORMAT_RGBA8           TiFormat = 3
	TI_FORMAT_RGBA8SRGB       TiFormat = 4
	TI_FORMAT_BGRA8           TiFormat = 5
	TI_FORMAT_BGRA8SRGB       TiFormat = 6
	TI_FORMAT_R8U             TiFormat = 7
	TI_FORMAT_RG8U            TiFormat = 8
	TI_FORMAT_RGBA8U          TiFormat = 9
	TI_FORMAT_R8I             TiFormat = 10
	TI_FORMAT_RG8I            TiFormat = 11
	TI_FORMAT_RGBA8I          TiFormat = 12
	TI_FORMAT_R16             TiFormat = 13
	TI_FORMAT_RG16            TiFormat = 14
	TI_FORMAT_RGB16           TiFormat = 15
	TI_FORMAT_RGBA16          TiFormat = 16
	TI_FORMAT_R16U            TiFormat = 17
	TI_FORMAT_RG16U           TiFormat = 18
	TI_FORMAT_RGB16U          TiFormat = 19
	TI_FORMAT_RGBA16U         TiFormat = 20
	TI_FORMAT_R16I            TiFormat = 21
	TI_FORMAT_RG16I           TiFormat = 22
	TI_FORMAT_RGB16I          TiFormat = 23
	TI_FORMAT_RGBA16I         TiFormat = 24
	TI_FORMAT_R16F            TiFormat = 25
	TI_FORMAT_RG16F           TiFormat = 26
	TI_FORMAT_RGB16F          TiFormat = 27
	TI_FORMAT_RGBA16F         TiFormat = 28
	TI_FORMAT_R32U            TiFormat = 29
	TI_FORMAT_RG32U           TiFormat = 30
	TI_FORMAT_RGB32U          TiFormat = 31
	TI_FORMAT_RGBA32U         TiFormat = 32
	TI_FORMAT_R32I            TiFormat = 33
	TI_FORMAT_RG32I           TiFormat = 34
	TI_FORMAT_RGB32I          TiFormat = 35
	TI_FORMAT_RGBA32I         TiFormat = 36
	TI_FORMAT_R32F            TiFormat = 37
	TI_FORMAT_RG32F           TiFormat = 38
	TI_FORMAT_RGB32F          TiFormat = 39
	TI_FORMAT_RGBA32F         TiFormat = 40
	TI_FORMAT_DEPTH16         TiFormat = 41
	TI_FORMAT_DEPTH24STENCIL8 TiFormat = 42
	TI_FORMAT_DEPTH32F        TiFormat = 43
)

// TiImageExtent 图像尺寸
// taichi_core.h:684-702
type TiImageExtent struct {
	Width           uint32
	Height          uint32
	Depth           uint32
	ArrayLayerCount uint32
}

// TiImageUsageFlags 图像用途标志
// taichi_core.h:548-556
type TiImageUsageFlags uint32

const (
	TI_IMAGE_USAGE_STORAGE_BIT    TiImageUsageFlags = 1 << 0
	TI_IMAGE_USAGE_SAMPLED_BIT    TiImageUsageFlags = 1 << 1
	TI_IMAGE_USAGE_ATTACHMENT_BIT TiImageUsageFlags = 1 << 2
)

// TiImageAllocateInfo 图像分配信息
// taichi_core.h:707-722
type TiImageAllocateInfo struct {
	Dimension     TiImageDimension
	Extent        TiImageExtent
	MipLevelCount uint32
	Format        TiFormat
	Export        TiBool
	Usage         TiImageUsageFlags
}

// TiFilter 过滤模式
// taichi_core.h:740-744
type TiFilter uint32

const (
	TI_FILTER_NEAREST TiFilter = 0
	TI_FILTER_LINEAR  TiFilter = 1
)

// TiAddressMode 寻址模式
// taichi_core.h:747-752
type TiAddressMode uint32

const (
	TI_ADDRESS_MODE_REPEAT          TiAddressMode = 0
	TI_ADDRESS_MODE_MIRRORED_REPEAT TiAddressMode = 1
	TI_ADDRESS_MODE_CLAMP_TO_EDGE   TiAddressMode = 2
)

// TiSamplerCreateInfo 采样器创建信息
// taichi_core.h:755-760
type TiSamplerCreateInfo struct {
	MagFilter     TiFilter
	MinFilter     TiFilter
	AddressMode   TiAddressMode
	MaxAnisotropy float32
}

// TiMemoryAllocateInfo 内存分配信息
// taichi_core.h:490-503
type TiMemoryAllocateInfo struct {
	Size      uint64
	HostWrite TiBool
	HostRead  TiBool
	Export    TiBool
	Usage     TiMemoryUsageFlags
}

// TiMemoryUsageFlags 内存用途标志
// taichi_core.h:475-484
type TiMemoryUsageFlags uint32

const (
	TI_MEMORY_USAGE_STORAGE_BIT TiMemoryUsageFlags = 1 << 0
	TI_MEMORY_USAGE_UNIFORM_BIT TiMemoryUsageFlags = 1 << 1
	TI_MEMORY_USAGE_VERTEX_BIT  TiMemoryUsageFlags = 1 << 2
	TI_MEMORY_USAGE_INDEX_BIT   TiMemoryUsageFlags = 1 << 3
)

// ===== 复合结构体定义 =====

// TiArgumentValue Argument值的联合体
// taichi_core.h:840-855
// C中是union，最大成员是TiNdArray (152 bytes)
type TiArgumentValue struct {
	Data [152]byte // 与C union的实际大小匹配
}

// TiArgument Kernel参数
// taichi_core.h:860-865
type TiArgument struct {
	Type  TiArgumentType  // 4 bytes, offset 0
	_     [4]byte         // padding，确保 Value 在 offset 8
	Value TiArgumentValue // 152 bytes, offset 8
}

// TiNamedArgument 命名参数（用于Compute Graph）
// taichi_core.h:870-875
type TiNamedArgument struct {
	Name     *byte
	Argument TiArgument
}

// TiNdShape N维形状
// taichi_core.h:522-527
type TiNdShape struct {
	DimCount uint32
	Dims     [16]uint32
}

func ToTiNdShape(uints []uint32) TiNdShape {
	var shape TiNdShape = TiNdShape{
		DimCount: min(16, uint32(len(uints))),
	}
	for i, dim := range uints {
		if i < 16 {
			shape.Dims[i] = dim
		}
	}
	return shape
}

// TiNdArray N维数组
// taichi_core.h:532-542
type TiNdArray struct {
	Memory    TiMemory
	Shape     TiNdShape
	ElemShape TiNdShape
	ElemType  TiDataType
}

// TiTexture 纹理
// taichi_core.h:765-777
type TiTexture struct {
	Image     TiImage
	Sampler   TiSampler
	Dimension TiImageDimension
	Extent    TiImageExtent
	Format    TiFormat
}

// TiMemorySlice 内存切片
// taichi_core.h:509-516
type TiMemorySlice struct {
	Memory TiMemory
	Offset uint64
	Size   uint64
}

// TiImageSlice 图像切片
// taichi_core.h:728-737
type TiImageSlice struct {
	Image    TiImage
	Offset   TiImageOffset
	Extent   TiImageExtent
	MipLevel uint32
}

// TiImageOffset 图像偏移
// taichi_core.h:662-679
type TiImageOffset struct {
	X                uint32
	Y                uint32
	Z                uint32
	ArrayLayerOffset uint32
}

// TiScalar 标量值
// taichi_core.h:802-805
type TiScalar struct {
	Type  TiDataType
	Value TiScalarValue
}

// TiScalarValue 标量值（联合体）
// taichi_core.h:788-797
type TiScalarValue struct {
	Data [8]byte
}

// TiTensorValue 张量值（联合体）
// taichi_core.h:810-819
type TiTensorValue struct {
	Data [128]byte
}

// TiTensorValueWithLength 带长度信息的张量值
// taichi_core.h:824-827
// C 中需要 8 字节对齐（因为 TiTensorValue 包含 uint64_t）
type TiTensorValueWithLength struct {
	Length uint32
	_      [4]byte // padding for 8-byte alignment of Data
	Data   TiTensorValue
}

// TiTensor 类型化张量值
// taichi_core.h:832-835
// C 中需要 8 字节对齐（因为 Contents 需要 8 字节对齐）
type TiTensor struct {
	Type     TiDataType
	_        [4]byte // padding for 8-byte alignment of Contents
	Contents TiTensorValueWithLength
}
