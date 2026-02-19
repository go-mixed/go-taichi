package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Image image abstraction
type Image struct {
	runtime *Runtime
	handle  c_api.TiImage
	width   uint32
	height  uint32
	format  Format
}

// NewImage2D creates a 2D image
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
		return nil, fmt.Errorf("image allocation failed [%d]: %s", errCode, errMsg)
	}

	return &Image{
		runtime: runtime,
		handle:  handle,
		width:   width,
		height:  height,
		format:  format,
	}, nil
}

// Release releases the image
func (img *Image) Release() {
	if img.handle != c_api.TI_NULL_HANDLE {
		c_api.FreeImage(img.runtime.handle, img.handle)
		img.handle = c_api.TI_NULL_HANDLE
	}
}

// Width gets the width
func (img *Image) Width() uint32 {
	return img.width
}

// Height gets the height
func (img *Image) Height() uint32 {
	return img.height
}

// Format gets the format
func (img *Image) Format() Format {
	return img.format
}

// TransitionLayout transitions the image layout
func (img *Image) TransitionLayout(layout ImageLayout) {
	c_api.TransitionImage(img.runtime.handle, img.handle, layout)
}

// CopyTo copies to another image (device-side)
func (img *Image) CopyTo(dst *Image) error {
	if img.width != dst.width || img.height != dst.height {
		return fmt.Errorf("image size mismatch: %dx%d vs %dx%d",
			img.width, img.height, dst.width, dst.height)
	}
	if img.format != dst.format {
		return fmt.Errorf("image format mismatch")
	}

	// Construct image slices
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

// CopyFrom copies from another image (device-side)
func (img *Image) CopyFrom(src *Image) error {
	return src.CopyTo(img)
}
