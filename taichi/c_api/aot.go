package c_api

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// ===== AOT模块函数指针 =====

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

// registerAotFunctions 注册AOT模块函数
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

// ===== 导出的AOT模块函数 =====

// LoadAotModule 从文件系统加载预编译的AOT模块
//
// 参数:
//   - runtime: 运行时句柄
//   - modulePath: AOT模块路径,应指向包含metadata.json的目录
//
// 返回:
//   - AOT模块句柄,如果加载失败则返回TI_NULL_HANDLE
//
// 示例:
//
//	module := taichi.LoadAotModule(runtime, "/path/to/aot/module")
//	if module == taichi.TI_NULL_HANDLE {
//	    log.Fatal("加载AOT模块失败")
//	}
//	defer taichi.DestroyAotModule(module)
func LoadAotModule(runtime TiRuntime, modulePath string) TiAotModule {
	cPath := append([]byte(modulePath), 0)
	return tiLoadAotModule(runtime, &cPath[0])
}

// CreateAotModule 从TCM数据创建预编译的AOT模块
//
// 参数:
//   - runtime: 运行时句柄
//   - tcm: TCM数据指针
//   - size: TCM数据大小
//
// 返回:
//   - AOT模块句柄,如果创建失败则返回TI_NULL_HANDLE
//
// 示例:
//
//	tcmData := loadTCMData() // 加载TCM数据
//	module := taichi.CreateAotModule(runtime, unsafe.Pointer(&tcmData[0]), uint64(len(tcmData)))
//	defer taichi.DestroyAotModule(module)
func CreateAotModule(runtime TiRuntime, tcm unsafe.Pointer, size uint64) TiAotModule {
	return tiCreateAotModule(runtime, tcm, size)
}

// DestroyAotModule 销毁已加载的AOT模块并释放所有相关资源
//
// 参数:
//   - aotModule: 要销毁的AOT模块句柄
//
// 注意:确保没有与该模块相关的kernel或compute graph等待flush。
//
// 示例:
//
//	taichi.DestroyAotModule(module)
func DestroyAotModule(aotModule TiAotModule) {
	tiDestroyAotModule(aotModule)
}

// GetAotModuleKernel 从AOT模块检索预编译的Taichi kernel
//
// 参数:
//   - aotModule: AOT模块句柄
//   - name: kernel名称
//
// 返回:
//   - kernel句柄,如果模块没有指定名称的kernel则返回TI_NULL_HANDLE
//
// 示例:
//
//	kernel := taichi.GetAotModuleKernel(module, "my_kernel")
//	if kernel == taichi.TI_NULL_HANDLE {
//	    log.Fatal("未找到kernel: my_kernel")
//	}
func GetAotModuleKernel(aotModule TiAotModule, name string) TiKernel {
	cName := append([]byte(name), 0)
	return tiGetAotModuleKernel(aotModule, &cName[0])
}

// GetAotModuleComputeGraph 从AOT模块检索预编译的compute graph
//
// 参数:
//   - aotModule: AOT模块句柄
//   - name: compute graph名称
//
// 返回:
//   - compute graph句柄,如果模块没有指定名称的compute graph则返回TI_NULL_HANDLE
//
// 示例:
//
//	graph := taichi.GetAotModuleComputeGraph(module, "my_graph")
//	if graph == taichi.TI_NULL_HANDLE {
//	    log.Fatal("未找到compute graph: my_graph")
//	}
func GetAotModuleComputeGraph(aotModule TiAotModule, name string) TiComputeGraph {
	cName := append([]byte(name), 0)
	return tiGetAotModuleComputeGraph(aotModule, &cName[0])
}

// LaunchKernel 使用提供的参数启动Taichi kernel
//
// 参数必须与源代码中的数量、类型和顺序相同。这是一个设备命令。
//
// 参数:
//   - runtime: 运行时句柄
//   - kernel: kernel句柄
//   - args: 参数数组
//
// 示例:
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
	if len(args) == 0 {
		tiLaunchKernel(runtime, kernel, 0, nil)
		return
	}
	tiLaunchKernel(runtime, kernel, uint32(len(args)), &args[0])
}

// LaunchComputeGraph 使用提供的命名参数启动Taichi compute graph
//
// 命名参数必须与源代码中的数量、名称和类型相同。这是一个设备命令。
//
// 参数:
//   - runtime: 运行时句柄
//   - computeGraph: compute graph句柄
//   - args: 命名参数数组
//
// 示例:
//
//	args := []taichi.TiNamedArgument{
//	    taichi.NewNamedArgument("foo", taichi.NewArgumentI32(123)),
//	    taichi.NewNamedArgument("bar", taichi.NewArgumentF32(456.0)),
//	}
//	taichi.LaunchComputeGraph(runtime, graph, args)
//	taichi.Flush(runtime)
//	taichi.Wait(runtime)
func LaunchComputeGraph(runtime TiRuntime, computeGraph TiComputeGraph, args []TiNamedArgument) {
	if len(args) == 0 {
		tiLaunchComputeGraph(runtime, computeGraph, 0, nil)
		return
	}
	tiLaunchComputeGraph(runtime, computeGraph, uint32(len(args)), &args[0])
}

// Flush 将之前调用的所有设备命令提交到目标设备执行
//
// 参数:
//   - runtime: 运行时句柄
//
// 示例:
//
//	taichi.LaunchKernel(runtime, kernel, args)
//	taichi.Flush(runtime)
func Flush(runtime TiRuntime) {
	tiFlush(runtime)
}

// Wait 等待之前调用的所有设备命令执行完成
//
// 任何已调用但未提交的命令将首先提交。
//
// 参数:
//   - runtime: 运行时句柄
//
// 示例:
//
//	taichi.LaunchKernel(runtime, kernel, args)
//	taichi.Flush(runtime)
//	taichi.Wait(runtime)
func Wait(runtime TiRuntime) {
	tiWait(runtime)
}
