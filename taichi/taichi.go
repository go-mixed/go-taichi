package taichi

import (
	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Init 初始化Taichi
// 必须在使用任何其他功能之前调用
func Init() error {
	return c_api.Init()
}

// GetVersion 获取Taichi C-API版本
func GetVersion() uint32 {
	c_api.Init()
	return c_api.GetVersion()
}

// GetAvailableArchs 获取所有可用的计算架构
func GetAvailableArchs() []c_api.TiArch {
	c_api.Init()
	return c_api.GetAvailableArchs()
}

// ===== 类型导出（供高级API使用）=====

// Arch 架构类型
type Arch = c_api.TiArch

// DataType 数据类型
type DataType = c_api.TiDataType

// Format 图像格式
type Format = c_api.TiFormat

// ImageLayout 图像布局
type ImageLayout = c_api.TiImageLayout

// Filter 过滤模式
type Filter = c_api.TiFilter

// AddressMode 寻址模式
type AddressMode = c_api.TiAddressMode

// ===== 常用常量导出 =====

// 架构常量
const (
	ArchVulkan = c_api.TI_ARCH_VULKAN
	ArchMetal  = c_api.TI_ARCH_METAL
	ArchCuda   = c_api.TI_ARCH_CUDA
	ArchX64    = c_api.TI_ARCH_X64
	ArchArm64  = c_api.TI_ARCH_ARM64
	ArchOpengl = c_api.TI_ARCH_OPENGL
)

// 数据类型常量
const (
	DataTypeF16 = c_api.TI_DATA_TYPE_F16
	DataTypeF32 = c_api.TI_DATA_TYPE_F32
	DataTypeF64 = c_api.TI_DATA_TYPE_F64
	DataTypeI8  = c_api.TI_DATA_TYPE_I8
	DataTypeI16 = c_api.TI_DATA_TYPE_I16
	DataTypeI32 = c_api.TI_DATA_TYPE_I32
	DataTypeI64 = c_api.TI_DATA_TYPE_I64
	DataTypeU8  = c_api.TI_DATA_TYPE_U8
	DataTypeU16 = c_api.TI_DATA_TYPE_U16
	DataTypeU32 = c_api.TI_DATA_TYPE_U32
	DataTypeU64 = c_api.TI_DATA_TYPE_U64
)

// 图像格式常量
const (
	FormatUnknown         = c_api.TI_FORMAT_UNKNOWN
	FormatR8              = c_api.TI_FORMAT_R8
	FormatRg8             = c_api.TI_FORMAT_RG8
	FormatRgba8           = c_api.TI_FORMAT_RGBA8
	FormatRgba8Srgb       = c_api.TI_FORMAT_RGBA8SRGB
	FormatBgra8           = c_api.TI_FORMAT_BGRA8
	FormatBgra8Srgb       = c_api.TI_FORMAT_BGRA8SRGB
	FormatR8U             = c_api.TI_FORMAT_R8U
	FormatRg8U            = c_api.TI_FORMAT_RG8U
	FormatRgba8U          = c_api.TI_FORMAT_RGBA8U
	FormatR8I             = c_api.TI_FORMAT_R8I
	FormatRg8I            = c_api.TI_FORMAT_RG8I
	FormatRgba8I          = c_api.TI_FORMAT_RGBA8I
	FormatR16             = c_api.TI_FORMAT_R16
	FormatRg16            = c_api.TI_FORMAT_RG16
	FormatRgb16           = c_api.TI_FORMAT_RGB16
	FormatRgba16          = c_api.TI_FORMAT_RGBA16
	FormatR16U            = c_api.TI_FORMAT_R16U
	FormatRg16U           = c_api.TI_FORMAT_RG16U
	FormatRgb16U          = c_api.TI_FORMAT_RGB16U
	FormatRgba16U         = c_api.TI_FORMAT_RGBA16U
	FormatR16I            = c_api.TI_FORMAT_R16I
	FormatRg16I           = c_api.TI_FORMAT_RG16I
	FormatRgb16I          = c_api.TI_FORMAT_RGB16I
	FormatRgba16I         = c_api.TI_FORMAT_RGBA16I
	FormatR16F            = c_api.TI_FORMAT_R16F
	FormatRg16F           = c_api.TI_FORMAT_RG16F
	FormatRgb16F          = c_api.TI_FORMAT_RGB16F
	FormatRgba16F         = c_api.TI_FORMAT_RGBA16F
	FormatR32U            = c_api.TI_FORMAT_R32U
	FormatRg32U           = c_api.TI_FORMAT_RG32U
	FormatRgb32U          = c_api.TI_FORMAT_RGB32U
	FormatRgba32U         = c_api.TI_FORMAT_RGBA32U
	FormatR32I            = c_api.TI_FORMAT_R32I
	FormatRg32I           = c_api.TI_FORMAT_RG32I
	FormatRgb32I          = c_api.TI_FORMAT_RGB32I
	FormatRgba32I         = c_api.TI_FORMAT_RGBA32I
	FormatR32F            = c_api.TI_FORMAT_R32F
	FormatRg32F           = c_api.TI_FORMAT_RG32F
	FormatRgb32F          = c_api.TI_FORMAT_RGB32F
	FormatRgba32F         = c_api.TI_FORMAT_RGBA32F
	FormatDepth16         = c_api.TI_FORMAT_DEPTH16
	FormatDepth24Stencil8 = c_api.TI_FORMAT_DEPTH24STENCIL8
	FormatDepth32F        = c_api.TI_FORMAT_DEPTH32F
)

// 图像布局常量
const (
	ImageLayoutUndefined       = c_api.TI_IMAGE_LAYOUT_UNDEFINED
	ImageLayoutShaderRead      = c_api.TI_IMAGE_LAYOUT_SHADER_READ
	ImageLayoutShaderWrite     = c_api.TI_IMAGE_LAYOUT_SHADER_WRITE
	ImageLayoutShaderReadWrite = c_api.TI_IMAGE_LAYOUT_SHADER_READ_WRITE
	ImageLayoutColorAttachment = c_api.TI_IMAGE_LAYOUT_COLOR_ATTACHMENT
	ImageLayoutDepthAttachment = c_api.TI_IMAGE_LAYOUT_DEPTH_ATTACHMENT
	ImageLayoutTransferDst     = c_api.TI_IMAGE_LAYOUT_TRANSFER_DST
	ImageLayoutTransferSrc     = c_api.TI_IMAGE_LAYOUT_TRANSFER_SRC
	ImageLayoutPresentSrc      = c_api.TI_IMAGE_LAYOUT_PRESENT_SRC
)

// 过滤模式常量
const (
	FilterNearest = c_api.TI_FILTER_NEAREST
	FilterLinear  = c_api.TI_FILTER_LINEAR
)

// 寻址模式常量
const (
	AddressModeRepeat         = c_api.TI_ADDRESS_MODE_REPEAT
	AddressModeMirroredRepeat = c_api.TI_ADDRESS_MODE_MIRRORED_REPEAT
	AddressModeClampToEdge    = c_api.TI_ADDRESS_MODE_CLAMP_TO_EDGE
)
