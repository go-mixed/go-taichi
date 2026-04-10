package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// ===== AOT Module Function Pointers =====

var (
	tiLoadAotModule            func(runtime TiRuntime, modulePath *byte) TiAotModule
	tiCreateAotModule          func(runtime TiRuntime, tcm unsafe.Pointer, size uint64) TiAotModule
	tiDestroyAotModule         func(aotModule TiAotModule)
	tiGetAotModuleKernel       func(aotModule TiAotModule, name *byte) TiKernel
	tiGetAotModuleComputeGraph func(aotModule TiAotModule, name *byte) TiComputeGraph
	tiLaunchKernel             func(runtime TiRuntime, kernel TiKernel, argCount uint32, args *TiArgument)
	tiLaunchComputeGraph       func(runtime TiRuntime, computeGraph TiComputeGraph, argCount uint32, args *TiNamedArgument)
	tiFlush                    func(runtime TiRuntime)
	tiWait                     func(runtime TiRuntime)
)

// registerAotFunctions registers AOT module functions
func registerAotFunctions() error {
	purego.RegisterLibFunc(&tiLoadAotModule, libHandle, "ti_load_aot_module")
	purego.RegisterLibFunc(&tiCreateAotModule, libHandle, "ti_create_aot_module")
	purego.RegisterLibFunc(&tiDestroyAotModule, libHandle, "ti_destroy_aot_module")
	purego.RegisterLibFunc(&tiGetAotModuleKernel, libHandle, "ti_get_aot_module_kernel")
	purego.RegisterLibFunc(&tiGetAotModuleComputeGraph, libHandle, "ti_get_aot_module_compute_graph")
	purego.RegisterLibFunc(&tiLaunchKernel, libHandle, "ti_launch_kernel")
	purego.RegisterLibFunc(&tiLaunchComputeGraph, libHandle, "ti_launch_compute_graph")
	purego.RegisterLibFunc(&tiFlush, libHandle, "ti_flush")
	purego.RegisterLibFunc(&tiWait, libHandle, "ti_wait")
	return nil
}

// ===== Exported AOT Module Functions =====

// LoadAotModule loads a precompiled AOT module from the filesystem
//
// Parameters:
//   - runtime: Runtime handle
//   - modulePath: Path to the AOT module, should point to a directory containing metadata.json
//
// Returns:
//   - AOT module handle, or TI_NULL_HANDLE if loading fails
//
// Example:
//
//	module := taichi.LoadAotModule(runtime, "/path/to/aot/module")
//	if module == taichi.TI_NULL_HANDLE {
//	    log.Fatal("Failed to load AOT module")
//	}
//	defer taichi.DestroyAotModule(module)
func LoadAotModule(runtime TiRuntime, modulePath string) TiAotModule {
	cPath := append([]byte(modulePath), 0)
	return SyncCall[TiAotModule](func() TiAotModule {
		return tiLoadAotModule(runtime, &cPath[0])
	})
}

// CreateAotModule creates a precompiled AOT module from TCM data
//
// Parameters:
//   - runtime: Runtime handle
//   - tcm: Pointer to TCM data
//   - size: Size of TCM data
//
// Returns:
//   - AOT module handle, or TI_NULL_HANDLE if creation fails
//
// Example:
//
//	tcmData := loadTCMData() // Load TCM data
//	module := taichi.CreateAotModule(runtime, unsafe.Pointer(&tcmData[0]), uint64(len(tcmData)))
//	defer taichi.DestroyAotModule(module)
func CreateAotModule(runtime TiRuntime, tcm unsafe.Pointer, size uint64) TiAotModule {
	return SyncCall[TiAotModule](func() TiAotModule {
		return tiCreateAotModule(runtime, tcm, size)
	})
}

// DestroyAotModule destroys a loaded AOT module and frees all associated resources
//
// Parameters:
//   - aotModule: AOT module handle to destroy
//
// Note: Ensure no kernels or compute graphs associated with this module are waiting to flush.
//
// Example:
//
//	taichi.DestroyAotModule(module)
func DestroyAotModule(aotModule TiAotModule) {
	SyncCallVoid(func() {
		tiDestroyAotModule(aotModule)
	})
}

// GetAotModuleKernel retrieves a precompiled Taichi kernel from an AOT module
//
// Parameters:
//   - aotModule: AOT module handle
//   - name: Kernel name
//
// Returns:
//   - Kernel handle, or TI_NULL_HANDLE if the module doesn't have a kernel with the specified name
//
// Example:
//
//	kernel := taichi.GetAotModuleKernel(module, "my_kernel")
//	if kernel == taichi.TI_NULL_HANDLE {
//	    log.Fatal("Kernel not found: my_kernel")
//	}
func GetAotModuleKernel(aotModule TiAotModule, name string) TiKernel {
	cName := append([]byte(name), 0)
	return SyncCall[TiKernel](func() TiKernel {
		return tiGetAotModuleKernel(aotModule, &cName[0])
	})
}

// GetAotModuleComputeGraph retrieves a precompiled compute graph from an AOT module
//
// Parameters:
//   - aotModule: AOT module handle
//   - name: Compute graph name
//
// Returns:
//   - Compute graph handle, or TI_NULL_HANDLE if the module doesn't have a compute graph with the specified name
//
// Example:
//
//	graph := taichi.GetAotModuleComputeGraph(module, "my_graph")
//	if graph == taichi.TI_NULL_HANDLE {
//	    log.Fatal("Compute graph not found: my_graph")
//	}
func GetAotModuleComputeGraph(aotModule TiAotModule, name string) TiComputeGraph {
	cName := append([]byte(name), 0)
	return SyncCall[TiComputeGraph](func() TiComputeGraph {
		return tiGetAotModuleComputeGraph(aotModule, &cName[0])
	})
}

// LaunchKernel launches a Taichi kernel with the provided arguments
//
// Arguments must match the number, type, and order in the source code. This is a device command.
//
// Parameters:
//   - runtime: Runtime handle
//   - kernel: Kernel handle
//   - args: Argument array
//
// Example:
//
//	args := []taichi.TiArgument{
//	    taichi.NewArgumentI32(123),
//	    taichi.NewArgumentF32(456.0),
//	    taichi.NewArgumentNdArray(ndarray),
//	}
//	taichi.LaunchKernel(runtime, kernel, args)
//	taichi.Flush(runtime)
//	taichi.Wait(runtime)
func LaunchKernel(runtime TiRuntime, kernel TiKernel, args []TiArgument) {
	SyncCallVoid(func() {
		if len(args) == 0 {
			tiLaunchKernel(runtime, kernel, 0, nil)
			return
		}
		tiLaunchKernel(runtime, kernel, uint32(len(args)), &args[0])
		asyncTasks.Add(1)
	})

}

// LaunchComputeGraph launches a Taichi compute graph with the provided named arguments
//
// Named arguments must match the number, names, and types in the source code. This is a device command.
//
// Parameters:
//   - runtime: Runtime handle
//   - computeGraph: Compute graph handle
//   - args: Named argument array
//
// Example:
//
//	args := []taichi.TiNamedArgument{
//	    taichi.NewNamedArgument("foo", taichi.NewArgumentI32(123)),
//	    taichi.NewNamedArgument("bar", taichi.NewArgumentF32(456.0)),
//	}
//	taichi.LaunchComputeGraph(runtime, graph, args)
//	taichi.Flush(runtime)
//	taichi.Wait(runtime)
func LaunchComputeGraph(runtime TiRuntime, computeGraph TiComputeGraph, args []TiNamedArgument) {
	SyncCallVoid(func() {
		if len(args) == 0 {
			tiLaunchComputeGraph(runtime, computeGraph, 0, nil)
			return
		}
		tiLaunchComputeGraph(runtime, computeGraph, uint32(len(args)), &args[0])
		asyncTasks.Add(1)
	})

}

// Flush submits all previously called device commands to the target device for execution
//
// Parameters:
//   - runtime: Runtime handle
//
// Example:
//
//	taichi.LaunchKernel(runtime, kernel, args)
//	taichi.Flush(runtime)
func Flush(runtime TiRuntime) {
	SyncCallVoid(func() {
		tiFlush(runtime)
	})
}

// Wait waits for all previously called device commands to complete execution
//
// Any called but not yet submitted commands will be submitted first.
//
// Parameters:
//   - runtime: Runtime handle
//
// Example:
//
//	taichi.LaunchKernel(runtime, kernel, args)
//	taichi.Flush(runtime)
//	taichi.Wait(runtime)
func Wait(runtime TiRuntime) {
	if asyncTasks.Swap(0) > 0 {
		SyncCallVoid(func() {
			tiWait(runtime)
		})
	}
}
