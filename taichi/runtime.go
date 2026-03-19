package taichi

import (
	"fmt"
	"os"

	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Runtime High-level abstraction for Taichi runtime
type Runtime struct {
	handle c_api.TiRuntime
	arch   Arch
}

// NewRuntime creates a new runtime
//
// Parameters:
//   - arch: Compute architecture. If 0, automatically selects the best architecture
//   - libDir: Dynamic library directory path
//   - Empty string (""): Search in current working directory first, then in system PATH
//   - Non-empty path: Search in specified directory first, then in system PATH
func NewRuntime(arch Arch, libDir string) (*Runtime, error) {
	// Initialize C-API
	if err := initial(libDir); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	// If architecture not specified, auto-select
	if arch == 0 {
		archs := GetAvailableArchs()
		if len(archs) == 0 {
			return nil, fmt.Errorf("no available compute architectures")
		}
		arch = selectBestArch(archs)
	}

	// Check TI_LIB_DIR for non-Vulkan backends
	if arch != ArchVulkan {
		tiLibDir := os.Getenv("TI_LIB_DIR")
		if tiLibDir == "" {
			return nil, fmt.Errorf("TI_LIB_DIR environment variable is required for %s backend, but it is not set", getArchName(arch))
		}
		if _, err := os.Stat(tiLibDir); err != nil {
			return nil, fmt.Errorf("TI_LIB_DIR is set to \"%s\", but the directory does not exist: %w", tiLibDir, err)
		}
	}

	// Create runtime
	handle := c_api.CreateRuntime(arch, 0)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to create runtime [%d]: %s", errCode, errMsg)
	}

	return &Runtime{
		handle: handle,
		arch:   arch,
	}, nil
}

// NewRuntimeAuto automatically selects the best architecture to create runtime
//
// Parameters:
//   - libDir: Dynamic library directory path
//   - Empty string (""): Search in current working directory first, then in system PATH
//   - Non-empty path: Search in specified directory first, then in system PATH
//
// Architecture selection priority: Vulkan > CUDA > x64 > ARM64 > OpenGL
func NewRuntimeAuto(libDir string) (*Runtime, error) {
	return NewRuntime(0, libDir)
}

// Release releases runtime resources
func (r *Runtime) Release() {
	if r.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroyRuntime(r.handle)
		r.handle = c_api.TI_NULL_HANDLE
	}
}

// Arch gets the current architecture
func (r *Runtime) Arch() Arch {
	return r.arch
}

// ArchName gets the architecture name
func (r *Runtime) ArchName() string {
	return getArchName(r.arch)
}

// Wait waits for all submitted tasks to complete
// Used to wait for all tasks to complete after asynchronous execution
func (r *Runtime) Wait() {
	c_api.Wait(r.handle)
}

// Flush flushes the command queue
// Ensures all submitted commands are sent to the device for execution
func (r *Runtime) Flush() {
	c_api.Flush(r.handle)
}

// selectBestArch selects the best architecture
func selectBestArch(archs []Arch) Arch {
	// Priority: Vulkan > CUDA > CPU
	priority := []Arch{
		ArchVulkan,
		ArchCuda,
		ArchX64,
		ArchArm64,
		ArchOpengl,
	}

	for _, preferred := range priority {
		for _, available := range archs {
			if available == preferred {
				return available
			}
		}
	}

	// Return first available
	return archs[0]
}

// getArchName gets the architecture name
func getArchName(arch Arch) string {
	switch arch {
	case ArchVulkan:
		return "Vulkan"
	case ArchCuda:
		return "CUDA"
	case ArchX64:
		return "x64 CPU"
	case ArchArm64:
		return "ARM64 CPU"
	case ArchOpengl:
		return "OpenGL"
	case ArchMetal:
		return "Metal"
	default:
		return fmt.Sprintf("Unknown(%d)", arch)
	}
}
