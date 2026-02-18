package taichi

import (
	"fmt"
	"go-taichi/taichi/c_api"
	"os"
	"unsafe"
)

// AotModule AOT模块抽象
type AotModule struct {
	runtime *Runtime
	handle  c_api.TiAotModule
}

// LoadAotModule 从文件系统加载AOT模块
// modulePath 应指向包含 metadata.json 的目录
func LoadAotModule(runtime *Runtime, modulePath string) (*AotModule, error) {
	handle := c_api.LoadAotModule(runtime.handle, modulePath)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("加载AOT模块失败 [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// LoadAotModuleFromTCM 从 .tcm 文件加载AOT模块
func LoadAotModuleFromTCM(runtime *Runtime, tcmPath string) (*AotModule, error) {
	// 读取 TCM 文件
	tcmData, err := os.ReadFile(tcmPath)
	if err != nil {
		return nil, fmt.Errorf("读取TCM文件失败: %w", err)
	}

	// 创建 AOT 模块
	handle := c_api.CreateAotModule(runtime.handle, unsafe.Pointer(&tcmData[0]), uint64(len(tcmData)))
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("从TCM创建模块失败 [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// Release 释放AOT模块
func (m *AotModule) Release() {
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroyAotModule(m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
}

// GetKernel 获取指定名称的kernel
func (m *AotModule) GetKernel(name string) (*Kernel, error) {
	handle := c_api.GetAotModuleKernel(m.handle, name)
	if handle == c_api.TI_NULL_HANDLE {
		return nil, fmt.Errorf("未找到kernel: %s", name)
	}

	return &Kernel{
		runtime: m.runtime,
		handle:  handle,
		name:    name,
	}, nil
}

// GetComputeGraph 获取指定名称的compute graph
func (m *AotModule) GetComputeGraph(name string) (*ComputeGraph, error) {
	handle := c_api.GetAotModuleComputeGraph(m.handle, name)
	if handle == c_api.TI_NULL_HANDLE {
		return nil, fmt.Errorf("未找到compute graph: %s", name)
	}

	return &ComputeGraph{
		runtime: m.runtime,
		handle:  handle,
		name:    name,
	}, nil
}

// Kernel Taichi kernel抽象
type Kernel struct {
	runtime *Runtime
	handle  c_api.TiKernel
	name    string
}

// Name 获取kernel名称
func (k *Kernel) Name() string {
	return k.name
}

// Launch 启动kernel（使用builder模式构建参数）
func (k *Kernel) Launch() *KernelLauncher {
	return &KernelLauncher{
		kernel: k,
		args:   make([]c_api.TiArgument, 0),
	}
}

// KernelLauncher Kernel启动器（Builder模式）
type KernelLauncher struct {
	kernel *Kernel
	args   []c_api.TiArgument
}

// ArgInt32 添加int32参数
func (kl *KernelLauncher) ArgInt32(value int32) *KernelLauncher {
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_I32,
	}
	*(*int32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	kl.args = append(kl.args, arg)
	return kl
}

// ArgFloat32 添加float32参数
func (kl *KernelLauncher) ArgFloat32(value float32) *KernelLauncher {
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_F32,
	}
	*(*float32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	kl.args = append(kl.args, arg)
	return kl
}

// ArgNdArray 添加NdArray参数
func (kl *KernelLauncher) ArgNdArray(arr *NdArray) *KernelLauncher {
	// 使用 c_api 包的辅助函数创建正确的 TiNdArray
	var ndarray c_api.TiNdArray

	switch len(arr.shape) {
	case 1:
		ndarray = c_api.NewNdArray1D(arr.Memory.handle, arr.shape[0], arr.elemType)
	case 2:
		ndarray = c_api.NewNdArray2D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.elemType)
	case 3:
		ndarray = c_api.NewNdArray3D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.shape[2], arr.elemType)
	default:
		// 对于更高维度，手动构造
		var dims [16]uint32
		for i, dim := range arr.shape {
			if i < 16 {
				dims[i] = dim
			}
		}
		ndarray = c_api.TiNdArray{
			Memory:    arr.Memory.handle,
			Shape:     c_api.TiNdShape{DimCount: uint32(len(arr.shape)), Dims: dims},
			ElemShape: c_api.TiNdShape{DimCount: 0, Dims: [16]uint32{}},
			ElemType:  arr.elemType,
		}
	}

	arg := c_api.NewArgumentNdArray(ndarray)
	kl.args = append(kl.args, arg)
	return kl
}

// Run 执行kernel并等待完成
func (kl *KernelLauncher) Run() {
	c_api.LaunchKernel(kl.kernel.runtime.handle, kl.kernel.handle, kl.args)
	c_api.Wait(kl.kernel.runtime.handle)
}

// RunAsync 异步执行kernel（不等待）
func (kl *KernelLauncher) RunAsync() {
	c_api.LaunchKernel(kl.kernel.runtime.handle, kl.kernel.handle, kl.args)
	c_api.Flush(kl.kernel.runtime.handle)
}

// ComputeGraph Compute Graph抽象
type ComputeGraph struct {
	runtime *Runtime
	handle  c_api.TiComputeGraph
	name    string
}

// Name 获取compute graph名称
func (cg *ComputeGraph) Name() string {
	return cg.name
}

// Launch 启动compute graph（使用builder模式构建参数）
func (cg *ComputeGraph) Launch() *GraphLauncher {
	return &GraphLauncher{
		graph: cg,
		args:  make([]c_api.TiNamedArgument, 0),
	}
}

// GraphLauncher Compute Graph启动器（Builder模式）
type GraphLauncher struct {
	graph *ComputeGraph
	args  []c_api.TiNamedArgument
}

// ArgInt32 添加int32命名参数
func (gl *GraphLauncher) ArgInt32(name string, value int32) *GraphLauncher {
	nameBytes := append([]byte(name), 0)
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_I32,
	}
	*(*int32)(unsafe.Pointer(&arg.Value.Data[0])) = value

	namedArg := c_api.TiNamedArgument{
		Name:     &nameBytes[0],
		Argument: arg,
	}
	gl.args = append(gl.args, namedArg)
	return gl
}

// ArgFloat32 添加float32命名参数
func (gl *GraphLauncher) ArgFloat32(name string, value float32) *GraphLauncher {
	nameBytes := append([]byte(name), 0)
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_F32,
	}
	*(*float32)(unsafe.Pointer(&arg.Value.Data[0])) = value

	namedArg := c_api.TiNamedArgument{
		Name:     &nameBytes[0],
		Argument: arg,
	}
	gl.args = append(gl.args, namedArg)
	return gl
}

// ArgNdArray 添加NdArray命名参数
func (gl *GraphLauncher) ArgNdArray(name string, arr *NdArray) *GraphLauncher {
	nameBytes := append([]byte(name), 0)

	// 使用 c_api 包的辅助函数创建正确的 TiNdArray
	var ndarray c_api.TiNdArray

	switch len(arr.shape) {
	case 1:
		ndarray = c_api.NewNdArray1D(arr.Memory.handle, arr.shape[0], arr.elemType)
	case 2:
		ndarray = c_api.NewNdArray2D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.elemType)
	case 3:
		ndarray = c_api.NewNdArray3D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.shape[2], arr.elemType)
	default:
		// 对于更高维度，手动构造
		var dims [16]uint32
		for i, dim := range arr.shape {
			if i < 16 {
				dims[i] = dim
			}
		}
		ndarray = c_api.TiNdArray{
			Memory:    arr.Memory.handle,
			Shape:     c_api.TiNdShape{DimCount: uint32(len(arr.shape)), Dims: dims},
			ElemShape: c_api.TiNdShape{DimCount: 0, Dims: [16]uint32{}},
			ElemType:  arr.elemType,
		}
	}

	arg := c_api.NewArgumentNdArray(ndarray)
	namedArg := c_api.TiNamedArgument{
		Name:     &nameBytes[0],
		Argument: arg,
	}
	gl.args = append(gl.args, namedArg)
	return gl
}

// Run 执行compute graph并等待完成
func (gl *GraphLauncher) Run() {
	c_api.LaunchComputeGraph(gl.graph.runtime.handle, gl.graph.handle, gl.args)
	c_api.Wait(gl.graph.runtime.handle)
}

// RunAsync 异步执行compute graph（不等待）
func (gl *GraphLauncher) RunAsync() {
	c_api.LaunchComputeGraph(gl.graph.runtime.handle, gl.graph.handle, gl.args)
	c_api.Flush(gl.graph.runtime.handle)
}
