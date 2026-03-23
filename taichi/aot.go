package taichi

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/go-mixed/go-taichi/taichi/c_api"
)

// AotModule AOT module abstraction
type AotModule struct {
	runtime *Runtime
	handle  c_api.TiAotModule
	tempDir string // temporary directory for extracted TCM (non-Vulkan backends)
}

// LoadAotModuleFile loads an AOT module from the filesystem
// modulePath should point to a directory containing metadata.json
func LoadAotModuleFile(runtime *Runtime, moduleDir string) (*AotModule, error) {
	handle := c_api.LoadAotModule(runtime.handle, moduleDir)
	if handle == c_api.TI_NULL_HANDLE {
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to load AOT module [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
	}, nil
}

// LoadAotModule loads an AOT module data from a .tcm file
//
// For Vulkan backend: passes TCM data directly to C-API
// For other backends (CPU/CUDA): extracts TCM zip to temp directory and loads from it
func LoadAotModule(runtime *Runtime, tcmData []byte) (*AotModule, error) {
	if len(tcmData) == 0 {
		return nil, fmt.Errorf("empty TCM file")
	}

	// Vulkan backend supports direct TCM data loading
	if runtime.arch == ArchVulkan {
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

	// Other backends (CPU/CUDA) require extracted TCM files
	// Extract TCM zip to temp directory
	useCache := runtime.options.cacheTcm
	tempDir, err := extractTCMToDir(tcmData, useCache)
	if err != nil {
		return nil, fmt.Errorf("failed to extract TCM: %w", err)
	}

	// Load from extracted directory
	handle := c_api.LoadAotModule(runtime.handle, tempDir)
	if handle == c_api.TI_NULL_HANDLE {
		os.RemoveAll(tempDir)
		errCode, errMsg := c_api.GetLastError()
		return nil, fmt.Errorf("failed to load AOT module [%d]: %s", errCode, errMsg)
	}

	return &AotModule{
		runtime: runtime,
		handle:  handle,
		tempDir: tempDir,
	}, nil
}

// extractTCMToDir extracts TCM zip data to a temporary directory
// Parameters:
//   - tcmData: TCM zip data
//   - useCache: if true, identical tcmData (by md5) reuses cached temp directory; if false, always creates new temp directory
func extractTCMToDir(tcmData []byte, useCache bool) (string, error) {
	var tempDir string
	var err error

	if useCache {
		// Use md5 as directory name for deterministic caching
		hash := md5.Sum(tcmData)
		hashStr := fmt.Sprintf("%x", hash)
		tempDir = filepath.Join(os.TempDir(), fmt.Sprintf("taichi_tcm_%s", hashStr))

		// Check if cached directory already exists
		if _, statErr := os.Stat(tempDir); statErr == nil {
			// Directory exists, reuse it
			return tempDir, nil
		}

		// Create the cache directory
		if mkErr := os.MkdirAll(tempDir, 0755); mkErr != nil {
			return "", fmt.Errorf("failed to create cache directory: %w", mkErr)
		}
	} else {
		// Create a unique temp directory each time
		tempDir, err = os.MkdirTemp("", "taichi_tcm_*")
		if err != nil {
			return "", fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	shouldRemove := func() {
		if !useCache {
			os.RemoveAll(tempDir)
		}
	}

	reader, err := zip.NewReader(bytes.NewReader(tcmData), int64(len(tcmData)))
	if err != nil {
		shouldRemove()
		return "", fmt.Errorf("failed to read TCM zip: %w", err)
	}

	for _, file := range reader.File {
		path := filepath.Join(tempDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				shouldRemove()
				return "", fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			shouldRemove()

			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}

		src, err := file.Open()
		if err != nil {
			shouldRemove()

			return "", fmt.Errorf("failed to open zip entry: %w", err)
		}

		dst, err := os.Create(path)
		if err != nil {
			src.Close()
			shouldRemove()
			return "", fmt.Errorf("failed to create file: %w", err)
		}

		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			dst.Close()
			shouldRemove()
			return "", fmt.Errorf("failed to extract file: %w", err)
		}

		src.Close()
		dst.Close()
	}

	return tempDir, nil
}

// Release releases the AOT module
// Note: tempDir is not deleted because it may be shared with other AotModule instances via TCM cache
func (m *AotModule) Release() {
	if m.handle != c_api.TI_NULL_HANDLE {
		c_api.DestroyAotModule(m.handle)
		m.handle = c_api.TI_NULL_HANDLE
	}
	// tempDir is intentionally not deleted - it's cached by md5 hash for reuse
	m.tempDir = ""
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
	// Build TiNdShape for shape
	shape := c_api.ToTiNdShape(arr.shape)

	// Build TiNdShape for elemShape (nil means scalar, DimCount=0)
	var elemShape c_api.TiNdShape
	if arr.elemShape != nil {
		elemShape = c_api.ToTiNdShape(arr.elemShape)
	}

	ndarray := c_api.TiNdArray{
		Memory:    arr.Memory.handle,
		Shape:     shape,
		ElemShape: elemShape,
		ElemType:  arr.elemType,
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

	// Build TiNdShape for shape
	var shapeDims [16]uint32
	for i, dim := range arr.shape {
		if i < 16 {
			shapeDims[i] = dim
		}
	}
	shape := c_api.TiNdShape{
		DimCount: uint32(len(arr.shape)),
		Dims:     shapeDims,
	}

	// Build TiNdShape for elemShape (nil means scalar, DimCount=0)
	var elemShape c_api.TiNdShape
	if arr.elemShape != nil {
		var elemDims [16]uint32
		for i, dim := range arr.elemShape {
			if i < 16 {
				elemDims[i] = dim
			}
		}
		elemShape = c_api.TiNdShape{
			DimCount: uint32(len(arr.elemShape)),
			Dims:     elemDims,
		}
	}

	ndarray = c_api.TiNdArray{
		Memory:    arr.Memory.handle,
		Shape:     shape,
		ElemShape: elemShape,
		ElemType:  arr.elemType,
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
