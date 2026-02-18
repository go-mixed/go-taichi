//go:build windows

package c_api

// openLibraryPosix Windows上不需要此函数,仅为了编译通过
func openLibraryPosix(path string) (uintptr, error) {
	panic("不应在Windows上调用openLibraryPosix")
}
