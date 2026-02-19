package taichi

import (
	"fmt"
	"github.com/go-mixed/go-taichi/taichi/c_api"
	"os"
	"unsafe"
)

// AotModule AOT module abstraction
type AotModule struct {
	runtime *Runtime
	handle  c_api.TiAotModule
}

// LoadAotModule loads an AOT module from the filesystem
// modulePath should point to a directory containing metadata.json
func LoadAotModule(runtime *Runtime, modulePath string) (*AotModule, error) {
	handle := c_api.LoadAotModule(runtime.handle, modulePath)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to load AOT module [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// LoadAotModuleFromTCM loads an AOT module from a .tcm file
func LoadAotModuleFromTCM(runtime *Runtime, tcmPath string) (*AotModule, error) {
	// Read TCM file
	tcmData, err := os.ReadFile(tcmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TCM file: %w", err)
	}

	// Create AOT module
	handle := c_api.CreateAotModule(runtime.handle, unsafe.Pointer(&tcmData[0]), uint64(len(tcmData)))
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to create module from TCM [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// Release releases the AOT module
func (m *AotModule) Release() {
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroyAotModule(m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
}

// GetKernel gets a kernel with the specified name
func (m *AotModule) GetKernel(name string) (*Kernel, error) {
	handle := c_api.GetAotModuleKernel(m.handle, name)
	if handle == c_api.TI_NULL_HANDLE {
		return nil, fmt.Errorf("kernel not found: %s", name)
	}

	return &Kernel{
		runtime: m.runtime,
		handle:  handle,
		name:    name,
	}, nil
}

// GetComputeGraph gets a compute graph with the specified name
func (m *AotModule) GetComputeGraph(name string) (*ComputeGraph, error) {
	handle := c_api.GetAotModuleComputeGraph(m.handle, name)
	if handle == c_api.TI_NULL_HANDLE {
		return nil, fmt.Errorf("compute graph not found: %s", name)
	}

	return &ComputeGraph{
		runtime: m.runtime,
		handle:  handle,
		name:    name,
	}, nil
}

// Kernel Taichi kernel abstraction
type Kernel struct {
	runtime *Runtime
	handle  c_api.TiKernel
	name    string
}

// Name gets the kernel name
func (k *Kernel) Name() string {
	return k.name
}

// Launch launches the kernel (uses builder pattern to construct arguments)
func (k *Kernel) Launch() *KernelLauncher {
	return &KernelLauncher{
		kernel: k,
		args:   make([]c_api.TiArgument, 0),
	}
}

// KernelLauncher Kernel launcher (Builder pattern)
type KernelLauncher struct {
	kernel *Kernel
	args   []c_api.TiArgument
}

// ArgInt32 adds an int32 argument
func (kl *KernelLauncher) ArgInt32(value int32) *KernelLauncher {
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_I32,
	}
	*(*int32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	kl.args = append(kl.args, arg)
	return kl
}

// ArgFloat32 adds a float32 argument
func (kl *KernelLauncher) ArgFloat32(value float32) *KernelLauncher {
	arg := c_api.TiArgument{
		Type: c_api.TI_ARGUMENT_TYPE_F32,
	}
	*(*float32)(unsafe.Pointer(&arg.Value.Data[0])) = value
	kl.args = append(kl.args, arg)
	return kl
}

// ArgNdArray adds an NdArray argument
func (kl *KernelLauncher) ArgNdArray(arr *NdArray) *KernelLauncher {
	// Use c_api package helper functions to create correct TiNdArray
	var ndarray c_api.TiNdArray

	switch len(arr.shape) {
	case 1:
		ndarray = c_api.NewNdArray1D(arr.Memory.handle, arr.shape[0], arr.elemType)
	case 2:
		ndarray = c_api.NewNdArray2D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.elemType)
	case 3:
		ndarray = c_api.NewNdArray3D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.shape[2], arr.elemType)
	default:
		// For higher dimensions, construct manually
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

// Run executes the kernel and waits for completion
func (kl *KernelLauncher) Run() {
	c_api.LaunchKernel(kl.kernel.runtime.handle, kl.kernel.handle, kl.args)
	c_api.Wait(kl.kernel.runtime.handle)
}

// RunAsync executes the kernel asynchronously (does not wait)
func (kl *KernelLauncher) RunAsync() {
	c_api.LaunchKernel(kl.kernel.runtime.handle, kl.kernel.handle, kl.args)
	c_api.Flush(kl.kernel.runtime.handle)
}

// ComputeGraph Compute Graph abstraction
type ComputeGraph struct {
	runtime *Runtime
	handle  c_api.TiComputeGraph
	name    string
}

// Name gets the compute graph name
func (cg *ComputeGraph) Name() string {
	return cg.name
}

// Launch launches the compute graph (uses builder pattern to construct arguments)
func (cg *ComputeGraph) Launch() *GraphLauncher {
	return &GraphLauncher{
		graph: cg,
		args:  make([]c_api.TiNamedArgument, 0),
	}
}

// GraphLauncher Compute Graph launcher (Builder pattern)
type GraphLauncher struct {
	graph *ComputeGraph
	args  []c_api.TiNamedArgument
}

// ArgInt32 adds an int32 named argument
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

// ArgFloat32 adds a float32 named argument
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

// ArgNdArray adds an NdArray named argument
func (gl *GraphLauncher) ArgNdArray(name string, arr *NdArray) *GraphLauncher {
	nameBytes := append([]byte(name), 0)

	// Use c_api package helper functions to create correct TiNdArray
	var ndarray c_api.TiNdArray

	switch len(arr.shape) {
	case 1:
		ndarray = c_api.NewNdArray1D(arr.Memory.handle, arr.shape[0], arr.elemType)
	case 2:
		ndarray = c_api.NewNdArray2D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.elemType)
	case 3:
		ndarray = c_api.NewNdArray3D(arr.Memory.handle, arr.shape[0], arr.shape[1], arr.shape[2], arr.elemType)
	default:
		// For higher dimensions, construct manually
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

// Run executes the compute graph and waits for completion
func (gl *GraphLauncher) Run() {
	c_api.LaunchComputeGraph(gl.graph.runtime.handle, gl.graph.handle, gl.args)
	c_api.Wait(gl.graph.runtime.handle)
}

// RunAsync executes the compute graph asynchronously (does not wait)
func (gl *GraphLauncher) RunAsync() {
	c_api.LaunchComputeGraph(gl.graph.runtime.handle, gl.graph.handle, gl.args)
	c_api.Flush(gl.graph.runtime.handle)
}
