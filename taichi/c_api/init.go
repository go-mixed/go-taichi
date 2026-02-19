// Package taichi provides Go language bindings for the Taichi C-API
// Uses purego for cross-platform support (Windows/Linux/macOS)
package c_api

import (
	"fmt"
	"os"
	"os/exec"
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
func Init(libDir string) error {
	// Determine library filename
	var libName string

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

	// Use specified directory
	libPath := filepath.Join(libDir, libName)

	// If not found, search in PATH
	if _, err := os.Stat(libPath); err != nil {
		// Search in system PATH
		libPath, err = exec.LookPath(libName)
		if err != nil {
			return fmt.Errorf("%s not found in system environment \"PATH\": %w", libName, err)
		}
	}

	// Load library file
	var handle uintptr

	switch runtime.GOOS {
	case "windows":
		// Windows uses syscall.LoadLibrary
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
