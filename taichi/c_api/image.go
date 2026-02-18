package c_api

import "github.com/ebitengine/purego"

// ===== 图像处理函数指针 =====

var (
	tiAllocateImage           func(runtime TiRuntime, allocateInfo *TiImageAllocateInfo) TiImage
	tiFreeImage               func(runtime TiRuntime, image TiImage)
	tiCreateSampler           func(runtime TiRuntime, createInfo *TiSamplerCreateInfo) TiSampler
	tiDestroySampler          func(runtime TiRuntime, sampler TiSampler)
	tiCopyImageDeviceToDevice func(runtime TiRuntime, dstImage *TiImageSlice, srcImage *TiImageSlice)
	tiTrackImageExt           func(runtime TiRuntime, image TiImage, layout TiImageLayout)
	tiTransitionImage         func(runtime TiRuntime, image TiImage, layout TiImageLayout)
)

// registerImageFunctions 注册图像处理函数
func registerImageFunctions() error {
	purego.RegisterLibFunc(&tiAllocateImage, libHandle, "ti_allocate_image")
	purego.RegisterLibFunc(&tiFreeImage, libHandle, "ti_free_image")
	purego.RegisterLibFunc(&tiCreateSampler, libHandle, "ti_create_sampler")
	purego.RegisterLibFunc(&tiDestroySampler, libHandle, "ti_destroy_sampler")
	purego.RegisterLibFunc(&tiCopyImageDeviceToDevice, libHandle, "ti_copy_image_device_to_device")
	purego.RegisterLibFunc(&tiTrackImageExt, libHandle, "ti_track_image_ext")
	purego.RegisterLibFunc(&tiTransitionImage, libHandle, "ti_transition_image")
	return nil
}

// ===== 导出的图像处理函数 =====

// AllocateImage 使用提供的参数分配设备图像
//
// 参数:
//   - runtime: 运行时句柄
//   - allocateInfo: 图像分配信息
//
// 返回:
//   - 图像句柄,如果分配失败则返回TI_NULL_HANDLE
//
// 示例:
//
//	allocInfo := taichi.TiImageAllocateInfo{
//	    Dimension: taichi.TI_IMAGE_DIMENSION_2D,
//	    Extent: taichi.TiImageExtent{
//	        Width:  1024,
//	        Height: 1024,
//	        Depth:  1,
//	        ArrayLayerCount: 1,
//	    },
//	    MipLevelCount: 1,
//	    Format:        taichi.TI_FORMAT_RGBA8,
//	    Usage:         taichi.TiImageUsageFlags(taichi.TI_IMAGE_USAGE_STORAGE_BIT | taichi.TI_IMAGE_USAGE_SAMPLED_BIT),
//	}
//	image := taichi.AllocateImage(runtime, &allocInfo)
//	defer taichi.FreeImage(runtime, image)
func AllocateImage(runtime TiRuntime, allocateInfo *TiImageAllocateInfo) TiImage {
	return tiAllocateImage(runtime, allocateInfo)
}

// FreeImage 释放图像分配
//
// 参数:
//   - runtime: 运行时句柄
//   - image: 要释放的图像句柄
//
// 示例:
//
//	taichi.FreeImage(runtime, image)
func FreeImage(runtime TiRuntime, image TiImage) {
	tiFreeImage(runtime, image)
}

// CreateSampler 创建图像采样器
//
// 参数:
//   - runtime: 运行时句柄
//   - createInfo: 采样器创建信息
//
// 返回:
//   - 采样器句柄,如果创建失败则返回TI_NULL_HANDLE
//
// 示例:
//
//	createInfo := taichi.TiSamplerCreateInfo{
//	    MagFilter:     taichi.TI_FILTER_LINEAR,
//	    MinFilter:     taichi.TI_FILTER_LINEAR,
//	    AddressMode:   taichi.TI_ADDRESS_MODE_CLAMP_TO_EDGE,
//	    MaxAnisotropy: 1.0,
//	}
//	sampler := taichi.CreateSampler(runtime, &createInfo)
//	defer taichi.DestroySampler(runtime, sampler)
func CreateSampler(runtime TiRuntime, createInfo *TiSamplerCreateInfo) TiSampler {
	return tiCreateSampler(runtime, createInfo)
}

// DestroySampler 销毁采样器
//
// 参数:
//   - runtime: 运行时句柄
//   - sampler: 要销毁的采样器句柄
//
// 示例:
//
//	taichi.DestroySampler(runtime, sampler)
func DestroySampler(runtime TiRuntime, sampler TiSampler) {
	tiDestroySampler(runtime, sampler)
}

// CopyImageDeviceToDevice 在设备内复制图像的连续子部分
//
// 两个子部分不能重叠。这是一个设备命令。
//
// 参数:
//   - runtime: 运行时句柄
//   - dst: 目标图像切片
//   - src: 源图像切片
//
// 示例:
//
//	srcSlice := &taichi.TiImageSlice{
//	    Image:    srcImage,
//	    MipLevel: 0,
//	    Offset:   taichi.TiImageOffset{X: 0, Y: 0, Z: 0},
//	    Extent:   taichi.TiImageExtent{Width: 512, Height: 512, Depth: 1},
//	}
//	dstSlice := &taichi.TiImageSlice{
//	    Image:    dstImage,
//	    MipLevel: 0,
//	    Offset:   taichi.TiImageOffset{X: 0, Y: 0, Z: 0},
//	    Extent:   taichi.TiImageExtent{Width: 512, Height: 512, Depth: 1},
//	}
//	taichi.CopyImageDeviceToDevice(runtime, dstSlice, srcSlice)
func CopyImageDeviceToDevice(runtime TiRuntime, dst *TiImageSlice, src *TiImageSlice) {
	tiCopyImageDeviceToDevice(runtime, dst, src)
}

// TrackImageExt 使用提供的图像布局跟踪设备图像
//
// 由于Taichi在内部跟踪图像布局,因此仅在通知Taichi图像已被外部过程转换为新布局时有用。
//
// 参数:
//   - runtime: 运行时句柄
//   - image: 图像句柄
//   - layout: 新的图像布局
//
// 示例:
//
//	taichi.TrackImageExt(runtime, image, taichi.TI_IMAGE_LAYOUT_SHADER_READ)
func TrackImageExt(runtime TiRuntime, image TiImage, layout TiImageLayout) {
	tiTrackImageExt(runtime, image, layout)
}

// TransitionImage 将图像转换为提供的图像布局
//
// 这是一个设备命令。由于Taichi在内部跟踪图像布局,
// 因此仅在强制执行图像布局供外部过程使用时有用。
//
// 参数:
//   - runtime: 运行时句柄
//   - image: 图像句柄
//   - layout: 目标图像布局
//
// 示例:
//
//	taichi.TransitionImage(runtime, image, taichi.TI_IMAGE_LAYOUT_SHADER_READ)
func TransitionImage(runtime TiRuntime, image TiImage, layout TiImageLayout) {
	tiTransitionImage(runtime, image, layout)
}
