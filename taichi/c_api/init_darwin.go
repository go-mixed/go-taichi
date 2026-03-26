//go:build darwin

package c_api

import "github.com/ebitengine/purego"

// getLibName 返回macOS平台的动态库名称
func getLibName() string {
	return "libtaichi_c_api.dylib"
}

// openLibrary 在macOS上使用purego.Dlopen打开动态库
func openLibrary(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
}
