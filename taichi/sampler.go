package taichi

import (
	"fmt"
	"go-taichi/taichi/c_api"
)

// Sampler 采样器抽象
//
// 采样器定义了如何从纹理中采样像素值，包括过滤模式、寻址模式等。
// 注意：大部分后端不支持自定义采样器，会使用默认采样器。
type Sampler struct {
	runtime *Runtime
	handle  c_api.TiSampler
}

// SamplerCreateInfo 采样器创建信息
type SamplerCreateInfo struct {
	// 放大过滤模式
	MagFilter c_api.TiFilter
	// 缩小过滤模式
	MinFilter c_api.TiFilter
	// 寻址模式 U 方向
	AddressModeU c_api.TiAddressMode
	// 寻址模式 V 方向
	AddressModeV c_api.TiAddressMode
	// 寻址模式 W 方向
	AddressModeW c_api.TiAddressMode
}

// NewSampler 创建新的采样器
//
// 注意：大部分后端不支持自定义采样器，此函数可能会失败。
// 在这种情况下，请使用 TI_NULL_HANDLE 作为默认采样器。
//
// 参数：
//   - runtime: Taichi 运行时
//   - createInfo: 采样器创建信息
//
// 返回：
//   - *Sampler: 创建的采样器对象
//   - error: 如果创建失败
//
// 示例：
//
//	info := &taichi.SamplerCreateInfo{
//	    MagFilter:     taichi.FILTER_LINEAR,
//	    MinFilter:     taichi.FILTER_LINEAR,
//	    AddressModeU:  taichi.ADDRESS_MODE_REPEAT,
//	    AddressModeV:  taichi.ADDRESS_MODE_REPEAT,
//	    AddressModeW:  taichi.ADDRESS_MODE_REPEAT,
//	}
//	sampler, err := taichi.NewSampler(runtime, info)
//	if err != nil {
//	    // 大部分后端不支持，使用默认采样器
//	    fmt.Println("使用默认采样器")
//	    sampler = nil
//	}
func NewSampler(runtime *Runtime, createInfo *SamplerCreateInfo) (*Sampler, error) {
	cInfo := &c_api.TiSamplerCreateInfo{
		MagFilter:    createInfo.MagFilter,
		MinFilter:    createInfo.MinFilter,
		AddressModeU: createInfo.AddressModeU,
		AddressModeV: createInfo.AddressModeV,
		AddressModeW: createInfo.AddressModeW,
	}

	handle := c_api.CreateSampler(runtime.handle, cInfo)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("创建采样器失败 [%d]: %s (提示：大部分后端不支持自定义采样器)", errCode, errMsg)
	}

	return &Sampler{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// Release 释放采样器
func (s *Sampler) Release() {
	if s.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroySampler(s.runtime.handle, s.handle)
		s.handle = c_api.TI_NULL_HANDLE
	}
}

// Handle 获取底层句柄（用于创建纹理）
func (s *Sampler) Handle() c_api.TiSampler {
	if s == nil {
		return c_api.TI_NULL_HANDLE
	}
	return s.handle
}

// IsValid 检查采样器是否有效
func (s *Sampler) IsValid() bool {
	return s != nil && s.handle != c_api.TI_NULL_HANDLE
}
