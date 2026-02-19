package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// Runtime Taichi运行时的高级抽象
type Runtime struct {
	handle c_api.TiRuntime
	arch   Arch
}

// NewRuntime 创建新的运行时
// 如果arch为0，自动选择最佳架构
func NewRuntime(arch Arch) (*Runtime, error) {
	// 初始化C-API
	if err := Init(); err != nil {
		return nil, fmt.Errorf("初始化失败: %w", err)
	}

	// 如果未指定架构，自动选择
	if arch == 0 {
		archs := GetAvailableArchs()
		if len(archs) == 0 {
			return nil, fmt.Errorf("没有可用的计算架构")
		}
		arch = selectBestArch(archs)
	}

	// 创建运行时
	handle := c_api.CreateRuntime(arch, 0)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("创建运行时失败 [%d]: %s", errCode, errMsg)
	}

	return &Runtime{
		handle: handle,
		arch:   arch,
	}, nil
}

// NewRuntimeAuto 自动选择最佳架构创建运行时
func NewRuntimeAuto() (*Runtime, error) {
	return NewRuntime(0)
}

// Release 释放运行时资源
func (r *Runtime) Release() {
	if r.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroyRuntime(r.handle)
		r.handle = c_api.TI_NULL_HANDLE
	}
}

// Arch 获取当前架构
func (r *Runtime) Arch() Arch {
	return r.arch
}

// ArchName 获取架构名称
func (r *Runtime) ArchName() string {
	return getArchName(r.arch)
}

// Wait 等待所有提交的任务完成
// 用于异步执行后等待所有任务完成
func (r *Runtime) Wait() {
	c_api.Wait(r.handle)
}

// Flush 刷新命令队列
// 确保所有提交的命令被发送到设备执行
func (r *Runtime) Flush() {
	c_api.Flush(r.handle)
}

// selectBestArch 选择最佳架构
func selectBestArch(archs []Arch) Arch {
	// 优先级: Vulkan > CUDA > CPU
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

	// 返回第一个可用的
	return archs[0]
}

// getArchName 获取架构名称
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
