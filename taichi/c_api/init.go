// Package taichi provides Go language bindings for the Taichi C-API
// Uses purego for cross-platform support (Windows/Linux/macOS)
package c_api

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

// libHandle dynamic library handle
var libHandle uintptr

// Initialized checks if initialization is complete
func Initialized() bool {
	return libHandle != 0
}

// Init initializes the Taichi C-API
//
// Parameters:
//   - libDir: Library file directory path
//   - Empty string (""): Search in current working directory first, then in system PATH
//   - Non-empty path: Search in specified directory first, then in system PATH
//
// Automatically loads the dynamic library:
//   - Windows: taichi_c_api.dll
//   - Linux: libtaichi_c_api.so
//   - macOS: libtaichi_c_api.dylib
//
// Must set before use: CGO_ENABLED=0
func Init() error {
	// Determine library filename
	var libName string
	// 添加/osName

	libDir := os.Getenv("TI_LIB_DIR")
	switch runtime.GOOS {
	case "windows":
		libName = "taichi_c_api.dll"
	case "linux":
		libName = "libtaichi_c_api.so"
	case "darwin":
		libName = "libtaichi_c_api.dylib"
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	libPath := filepath.Join(libDir, libName)

	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("TI_LIB_DIR environment variable is required: \"%s\" does not contain a valid taichi library. (%s, *.bc must be in the TI_LIB_DIR directory)", libDir, libName)
	}

	// Load library file
	var handle uintptr

	switch runtime.GOOS {
	case "windows":
		// Windows: set TI_LIB_DIR BEFORE loading DLL
		//if err := setLibEnv(libDir); err != nil {
		//	return fmt.Errorf("failed to set TI_LIB_DIR: %w", err)
		//}
		h, err := syscall.LoadLibrary(libPath)
		if err != nil {
			return fmt.Errorf("failed to load library: %w (path: %s)", err, libPath)
		}
		handle = uintptr(h)

	case "linux", "darwin":
		// Linux/macOS uses purego.Dlopen
		h, err := openLibraryPosix(libPath)
		if err != nil {
			return fmt.Errorf("failed to load library: %w (path: %s)", err, libPath)
		}
		handle = h
	}

	libHandle = handle

	// Register all functions
	if err := registerAllFunctions(); err != nil {
		return fmt.Errorf("failed to register functions: %w", err)
	}

	return nil
}

// registerAllFunctions registers all C API functions
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
