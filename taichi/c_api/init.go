// Package taichi 提供Taichi C-API的Go语言绑定
// 使用purego实现跨平台支持(Windows/Linux/macOS)
//
// 编译时请设置: CGO_ENABLED=0
package c_api

import (
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
)

// libHandle 动态链接库句柄
var libHandle uintptr

// Init 初始化Taichi C-API
//
// 自动从c_api/lib目录加载动态库:
//   - Windows: taichi_c_api.dll
//   - Linux: libtaichi_c_api.so
//   - macOS: libtaichi_c_api.dylib
//
// 使用前必须设置: CGO_ENABLED=0
func Init() error {
	// 获取当前文件所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("无法获取当前文件路径")
	}

	dir := filepath.Dir(filename)

	// 根据操作系统选择库文件
	var libPath string
	switch runtime.GOOS {
	case "windows":
		libPath = filepath.Join(dir, "lib", "taichi_c_api.dll")
		// Windows需要使用syscall.LoadLibrary
		handle, err := syscall.LoadLibrary(libPath)
		if err != nil {
			return fmt.Errorf("加载库失败: %w (路径: %s)", err, libPath)
		}
		libHandle = uintptr(handle)

	case "linux":
		libPath = filepath.Join(dir, "lib", "libtaichi_c_api.so")
		// Linux使用purego.Dlopen
		handle, err := openLibraryPosix(libPath)
		if err != nil {
			return fmt.Errorf("加载库失败: %w (路径: %s)", err, libPath)
		}
		libHandle = handle

	case "darwin":
		libPath = filepath.Join(dir, "lib", "libtaichi_c_api.dylib")
		// macOS使用purego.Dlopen
		handle, err := openLibraryPosix(libPath)
		if err != nil {
			return fmt.Errorf("加载库失败: %w (路径: %s)", err, libPath)
		}
		libHandle = handle

	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 注册所有函数
	if err := registerAllFunctions(); err != nil {
		return fmt.Errorf("注册函数失败: %w", err)
	}

	return nil
}

// registerAllFunctions 注册所有C API函数
func registerAllFunctions() error {
	if err := registerCoreFunctions(); err != nil {
		return err
	}
	if err := registerMemoryFunctions(); err != nil {
		return err
	}
	if err := registerAotFunctions(); err != nil {
		return err
	}
	if err := registerImageFunctions(); err != nil {
		return err
	}
	if err := registerMemoryImportFunctions(libHandle); err != nil {
		return err
	}
	return nil
}
