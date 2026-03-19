package c_api

import (
	"github.com/ebitengine/purego"
)

// ===== Core Function Pointers =====

var (
	tiGetVersion        func() uint32
	tiGetAvailableArchs func(archCount *uint32, archs *TiArch)
	tiGetLastError      func(messageSize *uint64, message *byte) TiError
	tiSetLastError      func(error TiError, message *byte)
	tiCreateRuntime     func(arch TiArch, deviceIndex uint32) TiRuntime
	tiDestroyRuntime    func(runtime TiRuntime)
)

// registerCoreFunctions registers core functions
func registerCoreFunctions() error {
	purego.RegisterLibFunc(&tiGetVersion, libHandle, "ti_get_version")
	purego.RegisterLibFunc(&tiGetAvailableArchs, libHandle, "ti_get_available_archs")
	purego.RegisterLibFunc(&tiGetLastError, libHandle, "ti_get_last_error")
	purego.RegisterLibFunc(&tiSetLastError, libHandle, "ti_set_last_error")
	purego.RegisterLibFunc(&tiCreateRuntime, libHandle, "ti_create_runtime")
	purego.RegisterLibFunc(&tiDestroyRuntime, libHandle, "ti_destroy_runtime")
	return nil
}

// ===== Exported Core Functions =====

// GetVersion gets the Taichi C-API version
//
// Returns the same value as TI_C_API_VERSION defined in taichi_core.h.
//
// Example:
//
//	version := taichi.GetVersion()
//	fmt.Printf("Taichi version: %d\n", version)
func GetVersion() uint32 {
	return SyncCall(func() uint32 {
		return tiGetVersion()
	})
}

// GetAvailableArchs gets the list of available architectures on the current platform
//
// An architecture is available only if:
// 1. The runtime library was compiled with support for that architecture
// 2. The current platform has the corresponding hardware or emulation software installed
//
// Available architectures have at least one device available, i.e., device index 0 is always available.
//
// Warning: The order of returned architectures is undefined.
//
// Example:
//
//	archs := taichi.GetAvailableArchs()
//	for _, arch := range archs {
//	    fmt.Printf("Available architecture: %d\n", arch)
//	}
func GetAvailableArchs() []TiArch {
	return SyncCall(func() []TiArch {
		var count uint32
		tiGetAvailableArchs(&count, nil)

		if count == 0 {
			return nil
		}

		archs := make([]TiArch, count)
		tiGetAvailableArchs(&count, &archs[0])
		return archs
	})
}

// GetLastError gets the last error raised by a Taichi C-API call
//
// Returns the semantic error code and text error message.
//
// Example:
//
//	errCode, errMsg := taichi.GetLastError()
//	if errCode != taichi.TI_ERROR_SUCCESS {
//	    fmt.Printf("Error: %d - %s\n", errCode, errMsg)
//	}
func GetLastError() (TiError, string) {
	var result TiError
	var resultMsg string
	SyncCallVoid(func() {
		var size uint64
		err := tiGetLastError(&size, nil)

		if size == 0 {
			result = err
			resultMsg = ""
			return
		}

		msg := make([]byte, size)
		err = tiGetLastError(&size, &msg[0])
		result = err
		resultMsg = string(msg[:size-1]) // Remove null terminator
	})
	return result, resultMsg
}

// SetLastError sets the provided error as the last error raised by a Taichi C-API call
//
// This is useful in extended validators for Taichi C-API wrappers and helper libraries.
//
// Parameters:
//   - error: Semantic error code
//   - message: Null-terminated string for text error message, or empty string for no error message
func SetLastError(error TiError, message string) {
	SyncCallVoid(func() {
		if message == "" {
			tiSetLastError(error, nil)
			return
		}
		msg := append([]byte(message), 0)
		tiSetLastError(error, &msg[0])
	})
}

// CreateRuntime creates a Taichi runtime with the specified architecture
//
// Parameters:
//   - arch: Architecture for the Taichi runtime
//   - deviceIndex: Device index on which to create the Taichi runtime
//
// Returns:
//   - Runtime handle, or TI_NULL_HANDLE if creation fails
//
// Example:
//
//	runtime := taichi.CreateRuntime(taichi.TI_ARCH_VULKAN, 0)
//	if runtime == taichi.TI_NULL_HANDLE {
//	    errCode, errMsg := taichi.GetLastError()
//	    log.Fatalf("Failed to create runtime: %d - %s", errCode, errMsg)
//	}
//	defer taichi.DestroyRuntime(runtime)
func CreateRuntime(arch TiArch, deviceIndex uint32) TiRuntime {
	if !runtimeRunning.CompareAndSwap(false, true) {
		panic("taichi runtime already running")
	}

	if arch == TI_ARCH_CUDA {
		runMainThread()
	}
	return SyncCall(func() TiRuntime {
		return tiCreateRuntime(arch, deviceIndex)
	})
}

// DestroyRuntime destroys a Taichi runtime
//
// Parameters:
//   - runtime: Runtime handle to destroy
//
// Note: All associated resources must be destroyed before destroying the runtime.
//
// Example:
//
//	runtime := taichi.CreateRuntime(taichi.TI_ARCH_VULKAN, 0)
//	defer taichi.DestroyRuntime(runtime)
func DestroyRuntime(runtime TiRuntime) {
	SyncCallVoid(func() {
		tiDestroyRuntime(runtime)
	})

	closeMainThread()
	runtimeRunning.Store(false)

}
