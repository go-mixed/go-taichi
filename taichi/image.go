package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Image 图像抽象
type Image struct {
	runtime *Runtime
	handle  c_api.TiImage
	width   uint32
	height  uint32
	format  Format
}

// NewImage2D 创建2D图像
func NewImage2D(runtime *Runtime, width, height uint32, format Format) (*Image, error) {
	allocInfo := c_api.TiImageAllocateInfo{
		Dimension: c_api.TI_IMAGE_DIMENSION_2D,
		Extent: c_api.TiImageExtent{
			Width:           width,
			Height:          height,
			Depth:           1,
			ArrayLayerCount: 1,
		},
		MipLevelCount: 1,
		Format:        format,
		Usage: c_api.TiImageUsageFlags(
			c_api.TI_IMAGE_USAGE_STORAGE_BIT |
				c_api.TI_IMAGE_USAGE_SAMPLED_BIT,
		),
	}

	handle := c_api.AllocateImage(runtime.handle, &allocInfo)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("图像分配失败 [%d]: %s", errCode, errMsg)
	}

	return &Image{
		runtime: runtime,
		handle:  handle,
		width:   width,
		height:  height,
		format:  format,
	}, nil
}

// Release 释放图像
func (img *Image) Release() {
	if img.handle != c_api.TI_NULL_HANDLE {
		c_api.FreeImage(img.runtime.handle, img.handle)
		img.handle = c_api.TI_NULL_HANDLE
	}
}

// Width 获取宽度
func (img *Image) Width() uint32 {
	return img.width
}

// Height 获取高度
func (img *Image) Height() uint32 {
	return img.height
}

// Format 获取格式
func (img *Image) Format() Format {
	return img.format
}

// TransitionLayout 转换图像布局
func (img *Image) TransitionLayout(layout ImageLayout) {
	c_api.TransitionImage(img.runtime.handle, img.handle, layout)
}

// CopyTo 复制到另一个图像（设备端）
func (img *Image) CopyTo(dst *Image) error {
	if img.width != dst.width || img.height != dst.height {
		return fmt.Errorf("图像尺寸不匹配: %dx%d vs %dx%d",
			img.width, img.height, dst.width, dst.height)
	}
	if img.format != dst.format {
		return fmt.Errorf("图像格式不匹配")
	}

	// 构造图像切片
	srcSlice := &c_api.TiImageSlice{
		Image:    img.handle,
		MipLevel: 0,
		Offset: c_api.TiImageOffset{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Extent: c_api.TiImageExtent{
			Width:           img.width,
			Height:          img.height,
			Depth:           1,
			ArrayLayerCount: 1,
		},
	}

	dstSlice := &c_api.TiImageSlice{
		Image:    dst.handle,
		MipLevel: 0,
		Offset: c_api.TiImageOffset{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Extent: c_api.TiImageExtent{
			Width:           dst.width,
			Height:          dst.height,
			Depth:           1,
			ArrayLayerCount: 1,
		},
	}

	c_api.CopyImageDeviceToDevice(img.runtime.handle, dstSlice, srcSlice)
	return nil
}

// CopyFrom 从另一个图像复制（设备端）
func (img *Image) CopyFrom(src *Image) error {
	return src.CopyTo(img)
}
