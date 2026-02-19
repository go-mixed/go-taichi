# Image API Reference

Low-level image and texture management functions.

## Image Allocation

### AllocateImage

```go
func AllocateImage(runtime TiRuntime, allocInfo *TiImageAllocateInfo) TiImage
```

Allocate image/texture.

**Example**:
```go
allocInfo := c_api.TiImageAllocateInfo{
    Dimension: c_api.TI_IMAGE_DIMENSION_2D,
    Extent:    c_api.TiImageExtent{Width: 512, Height: 512, Depth: 1},
    Format:    c_api.TI_FORMAT_RGBA8,
    Usage:     c_api.TiImageUsageFlags(c_api.TI_IMAGE_USAGE_STORAGE_BIT | c_api.TI_IMAGE_USAGE_SAMPLED_BIT),
}
image := c_api.AllocateImage(runtime, &allocInfo)
defer c_api.FreeImage(runtime, image)
```

### FreeImage

```go
func FreeImage(runtime TiRuntime, image TiImage)
```

Free allocated image.

---

## Image Operations

### TransitionImage

```go
func TransitionImage(runtime TiRuntime, image TiImage, layout TiImageLayout)
```

Transition image layout.

**Common Layouts**:
- `TI_IMAGE_LAYOUT_UNDEFINED` - Initial state
- `TI_IMAGE_LAYOUT_SHADER_READ` - For sampling
- `TI_IMAGE_LAYOUT_SHADER_WRITE` - For writing
- `TI_IMAGE_LAYOUT_TRANSFER_DST` - For copying to
- `TI_IMAGE_LAYOUT_TRANSFER_SRC` - For copying from

**Example**:
```go
c_api.TransitionImage(runtime, image, c_api.TI_IMAGE_LAYOUT_SHADER_READ)
```

### CopyImageDeviceToDevice

```go
func CopyImageDeviceToDevice(runtime TiRuntime, dstSlice *TiImageSlice, srcSlice *TiImageSlice)
```

Copy between images.

---

## Sampler

### CreateSampler

```go
func CreateSampler(runtime TiRuntime, createInfo *TiSamplerCreateInfo) TiSampler
```

Create texture sampler.

**Note**: Most backends don't support custom samplers. Use `TI_NULL_HANDLE` for default.

### DestroySampler

```go
func DestroySampler(runtime TiRuntime, sampler TiSampler)
```

Destroy sampler.
