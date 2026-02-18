package c_api

import "github.com/ebitengine/purego"

// ===== 核心函数指针 =====

var (
	tiGetVersion        func() uint32
	tiGetAvailableArchs func(archCount *uint32, archs *TiArch)
	tiGetLastError      func(messageSize *uint64, message *byte) TiError
	tiSetLastError      func(error TiError, message *byte)
	tiCreateRuntime     func(arch TiArch, deviceIndex uint32) TiRuntime
	tiDestroyRuntime    func(runtime TiRuntime)
)

// registerCoreFunctions 注册核心函数
func registerCoreFunctions() error {
	purego.RegisterLibFunc(&tiGetVersion, libHandle, "ti_get_version")
	purego.RegisterLibFunc(&tiGetAvailableArchs, libHandle, "ti_get_available_archs")
	purego.RegisterLibFunc(&tiGetLastError, libHandle, "ti_get_last_error")
	purego.RegisterLibFunc(&tiSetLastError, libHandle, "ti_set_last_error")
	purego.RegisterLibFunc(&tiCreateRuntime, libHandle, "ti_create_runtime")
	purego.RegisterLibFunc(&tiDestroyRuntime, libHandle, "ti_destroy_runtime")
	return nil
}

// ===== 导出的核心函数 =====

// GetVersion 获取Taichi C-API版本
//
// 返回值与taichi_core.h中定义的TI_C_API_VERSION相同。
//
// 示例:
//
//	version := taichi.GetVersion()
//	fmt.Printf("Taichi版本: %d\n", version)
func GetVersion() uint32 {
	return tiGetVersion()
}

// GetAvailableArchs 获取当前平台上可用的架构列表
//
// 架构只有在以下情况下才可用:
// 1. Runtime库编译时支持该架构
// 2. 当前平台安装了相应的硬件或模拟软件
//
// 可用架构至少有一个设备可用,即设备索引0始终可用。
//
// 警告:返回架构的顺序未定义。
//
// 示例:
//
//	archs := taichi.GetAvailableArchs()
//	for _, arch := range archs {
//	    fmt.Printf("可用架构: %d\n", arch)
//	}
func GetAvailableArchs() []TiArch {
	var count uint32
	tiGetAvailableArchs(&count, nil)

	if count == 0 {
		return nil
	}

	archs := make([]TiArch, count)
	tiGetAvailableArchs(&count, &archs[0])
	return archs
}

// GetLastError 获取Taichi C-API调用引发的最后一个错误
//
// 返回语义错误代码和文本错误消息。
//
// 示例:
//
//	errCode, errMsg := taichi.GetLastError()
//	if errCode != taichi.TI_ERROR_SUCCESS {
//	    fmt.Printf("错误: %d - %s\n", errCode, errMsg)
//	}
func GetLastError() (TiError, string) {
	var size uint64
	err := tiGetLastError(&size, nil)

	if size == 0 {
		return err, ""
	}

	msg := make([]byte, size)
	err = tiGetLastError(&size, &msg[0])
	return err, string(msg[:size-1]) // 去掉null terminator
}

// SetLastError 将提供的错误设置为Taichi C-API调用引发的最后一个错误
//
// 这在Taichi C-API包装器和辅助库的扩展验证程序中很有用。
//
// 参数:
//   - error: 语义错误代码
//   - message: 文本错误消息的null结尾字符串,或空字符串表示空错误消息
func SetLastError(error TiError, message string) {
	if message == "" {
		tiSetLastError(error, nil)
		return
	}
	msg := append([]byte(message), 0)
	tiSetLastError(error, &msg[0])
}

// CreateRuntime 使用指定的架构创建Taichi运行时
//
// 参数:
//   - arch: Taichi运行时的架构
//   - deviceIndex: 要在其上创建Taichi运行时的设备索引
//
// 返回:
//   - 运行时句柄,如果创建失败则返回TI_NULL_HANDLE
//
// 示例:
//
//	runtime := taichi.CreateRuntime(taichi.TI_ARCH_VULKAN, 0)
//	if runtime == taichi.TI_NULL_HANDLE {
//	    errCode, errMsg := taichi.GetLastError()
//	    log.Fatalf("创建运行时失败: %d - %s", errCode, errMsg)
//	}
//	defer taichi.DestroyRuntime(runtime)
func CreateRuntime(arch TiArch, deviceIndex uint32) TiRuntime {
	return tiCreateRuntime(arch, deviceIndex)
}

// DestroyRuntime 销毁Taichi运行时
//
// 参数:
//   - runtime: 要销毁的运行时句柄
//
// 注意:销毁运行时之前,必须先销毁所有相关资源。
//
// 示例:
//
//	runtime := taichi.CreateRuntime(taichi.TI_ARCH_VULKAN, 0)
//	defer taichi.DestroyRuntime(runtime)
func DestroyRuntime(runtime TiRuntime) {
	tiDestroyRuntime(runtime)
}
