//go:build windows

package c_api

import "syscall"

// getLibName 返回Windows平台的动态库名称
func getLibName() string {
	return "taichi_c_api.dll"
}

// openLibrary 在Windows上使用syscall.LoadLibrary加载动态库
func openLibrary(path string) (uintptr, error) {
	h, err := syscall.LoadLibrary(path)
	if err != nil {
		return 0, err
	}
	return uintptr(h), nil
}
