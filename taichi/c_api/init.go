// Package taichi 提供Taichi C-API的Go语言绑定
// 使用purego实现跨平台支持(Windows/Linux/macOS)
package c_api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// libHandle 动态链接库句柄
var libHandle uintptr

// Initialized 检查是否已初始化
func Initialized() bool {
	return libHandle != 0
}

// Init 初始化Taichi C-API
//
// 参数:
//   - libDir: 库文件目录路径
//   - 空字符串(""): 先在当前工作目录查找，找不到则在系统PATH中查找
//   - 非空路径: 先在指定目录查找，找不到则在系统PATH中查找
//
// 自动加载动态库:
//   - Windows: taichi_c_api.dll
//   - Linux: libtaichi_c_api.so
//   - macOS: libtaichi_c_api.dylib
//
// 使用前必须设置: CGO_ENABLED=0
func Init(libDir string) error {
	// 确定库文件名
	var libName string

	switch runtime.GOOS {
	case "windows":
		libName = "taichi_c_api.dll"
	case "linux":
		libName = "libtaichi_c_api.so"
	case "darwin":
		libName = "libtaichi_c_api.dylib"
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 使用指定目录
	libPath := filepath.Join(libDir, libName)

	// 如果不存在，则在PATH中查找
	if _, err := os.Stat(libPath); err != nil {
		// 在系统PATH中查找
		libPath, err = exec.LookPath(libName)
		if err != nil {
			return fmt.Errorf("%s not found in system environment \"PATH\": %w", libName, err)
		}
	}

	// 加载库文件
	var handle uintptr

	switch runtime.GOOS {
	case "windows":
		// Windows使用syscall.LoadLibrary
		h, err := syscall.LoadLibrary(libPath)
		if err != nil {
			return fmt.Errorf("加载库失败: %w (路径: %s)", err, libPath)
		}
		handle = uintptr(h)

	case "linux", "darwin":
		// Linux/macOS使用purego.Dlopen
		h, err := openLibraryPosix(libPath)
		if err != nil {
			return fmt.Errorf("加载库失败: %w (路径: %s)", err, libPath)
		}
		handle = h
	}

	libHandle = handle

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
