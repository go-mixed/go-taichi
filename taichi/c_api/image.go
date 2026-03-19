package c_api

import "github.com/ebitengine/purego"

// ===== Image Processing Function Pointers =====

var (
	tiAllocateImage           func(runtime TiRuntime, allocateInfo *TiImageAllocateInfo) TiImage
	tiFreeImage               func(runtime TiRuntime, image TiImage)
	tiCreateSampler           func(runtime TiRuntime, createInfo *TiSamplerCreateInfo) TiSampler
	tiDestroySampler          func(runtime TiRuntime, sampler TiSampler)
	tiCopyImageDeviceToDevice func(runtime TiRuntime, dstImage *TiImageSlice, srcImage *TiImageSlice)
	tiTrackImageExt           func(runtime TiRuntime, image TiImage, layout TiImageLayout)
	tiTransitionImage         func(runtime TiRuntime, image TiImage, layout TiImageLayout)
)

// registerImageFunctions registers image processing functions
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

// ===== Exported Image Processing Functions =====

// AllocateImage allocates a device image with the provided parameters
//
// Parameters:
//   - runtime: Runtime handle
//   - allocateInfo: Image allocation information
//
// Returns:
//   - Image handle, or TI_NULL_HANDLE if allocation fails
//
// Example:
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
	return SyncCall(func() TiImage {
		return tiAllocateImage(runtime, allocateInfo)
	})
}

// FreeImage frees an image allocation
//
// Parameters:
//   - runtime: Runtime handle
//   - image: Image handle to free
//
// Example:
//
//	taichi.FreeImage(runtime, image)
func FreeImage(runtime TiRuntime, image TiImage) {
	SyncCallVoid(func() {
		tiFreeImage(runtime, image)
	})
}

// CreateSampler creates an image sampler
//
// Parameters:
//   - runtime: Runtime handle
//   - createInfo: Sampler creation information
//
// Returns:
//   - Sampler handle, or TI_NULL_HANDLE if creation fails
//
// Example:
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
	return SyncCall(func() TiSampler {
		return tiCreateSampler(runtime, createInfo)
	})
}

// DestroySampler destroys a sampler
//
// Parameters:
//   - runtime: Runtime handle
//   - sampler: Sampler handle to destroy
//
// Example:
//
//	taichi.DestroySampler(runtime, sampler)
func DestroySampler(runtime TiRuntime, sampler TiSampler) {
	SyncCallVoid(func() {
		tiDestroySampler(runtime, sampler)
	})
}

// CopyImageDeviceToDevice copies a contiguous subsection of an image within the device
//
// The two subsections must not overlap. This is a device command.
//
// Parameters:
//   - runtime: Runtime handle
//   - dst: Destination image slice
//   - src: Source image slice
//
// Example:
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
	SyncCallVoid(func() {
		tiCopyImageDeviceToDevice(runtime, dst, src)
	})
}

// TrackImageExt tracks a device image with the provided image layout
//
// Since Taichi tracks image layouts internally, this is only useful for notifying Taichi
// that an image has been transitioned to a new layout by an external process.
//
// Parameters:
//   - runtime: Runtime handle
//   - image: Image handle
//   - layout: New image layout
//
// Example:
//
//	taichi.TrackImageExt(runtime, image, taichi.TI_IMAGE_LAYOUT_SHADER_READ)
func TrackImageExt(runtime TiRuntime, image TiImage, layout TiImageLayout) {
	SyncCallVoid(func() {
		tiTrackImageExt(runtime, image, layout)
	})
}

// TransitionImage transitions an image to the provided image layout
//
// This is a device command. Since Taichi tracks image layouts internally,
// this is only useful for forcing an image layout for use by external processes.
//
// Parameters:
//   - runtime: Runtime handle
//   - image: Image handle
//   - layout: Target image layout
//
// Example:
//
//	taichi.TransitionImage(runtime, image, taichi.TI_IMAGE_LAYOUT_SHADER_READ)
func TransitionImage(runtime TiRuntime, image TiImage, layout TiImageLayout) {
	SyncCallVoid(func() {
		tiTransitionImage(runtime, image, layout)
	})
}
