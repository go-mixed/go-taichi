package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Sampler sampler abstraction
//
// Samplers define how to sample pixel values from textures, including filter modes, addressing modes, etc.
// Note: Most backends do not support custom samplers and will use the default sampler.
type Sampler struct {
	runtime *Runtime
	handle  c_api.TiSampler
}

// SamplerCreateInfo sampler creation information
type SamplerCreateInfo struct {
	// Magnification filter mode
	MagFilter c_api.TiFilter
	// Minification filter mode
	MinFilter c_api.TiFilter
	// Addressing mode U direction
	AddressModeU c_api.TiAddressMode
	// Addressing mode V direction
	AddressModeV c_api.TiAddressMode
	// Addressing mode W direction
	AddressModeW c_api.TiAddressMode
}

// NewSampler creates a new sampler
//
// Note: Most backends do not support custom samplers, this function may fail.
// In that case, please use TI_NULL_HANDLE as the default sampler.
//
// Parameters:
//   - runtime: Taichi runtime
//   - createInfo: Sampler creation information
//
// Returns:
//   - *Sampler: Created sampler object
//   - error: If creation fails
//
// Example:
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
//	    // Most backends don't support it, use default sampler
//	    fmt.Println("Using default sampler")
//	    sampler = nil
//	}
func NewSampler(runtime *Runtime, createInfo *SamplerCreateInfo) (*Sampler, error) {
	cInfo := &c_api.TiSamplerCreateInfo{
		MagFilter:   createInfo.MagFilter,
		MinFilter:   createInfo.MinFilter,
		AddressMode: createInfo.AddressModeU, // Use U as unified address mode
	}

	handle := c_api.CreateSampler(runtime.handle, cInfo)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to create sampler [%d]: %s (note: most backends do not support custom samplers)", errCode, errMsg)
	}

	return &Sampler{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// Release releases the sampler
func (s *Sampler) Release() {
	if s.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroySampler(s.runtime.handle, s.handle)
		s.handle = c_api.TI_NULL_HANDLE
	}
}

// Handle gets the underlying handle (for creating textures)
func (s *Sampler) Handle() c_api.TiSampler {
	if s == nil {
		return c_api.TI_NULL_HANDLE
	}
	return s.handle
}

// IsValid checks if the sampler is valid
func (s *Sampler) IsValid() bool {
	return s != nil && s.handle != c_api.TI_NULL_HANDLE
}
