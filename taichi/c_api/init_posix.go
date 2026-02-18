//go:build !windows

package c_api

import "github.com/ebitengine/purego"

// openLibraryPosix 在POSIX系统(Linux/macOS)上打开动态库
func openLibraryPosix(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
}
