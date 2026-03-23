# Go-Taichi

> Pure Go bindings for Taichi C-API - High-performance GPU parallel computing


[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://go.dev/)
[![Taichi Version](https://img.shields.io/badge/Taichi-1.7.4-green)](https://www.taichi-lang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-orange)](LICENSE)

## Features

- 🚀 **High Performance** - GPU-accelerated computing powered by Taichi
- 🎯 **Simple API** - High-level abstractions with automatic resource management
- 🔧 **Cross-Platform** - Windows / Linux / macOS
- 💻 **Multiple Backends** - Vulkan / CUDA / CPU / Metal
- 📦 **Pure Go** - No CGo required, No *.c/h/hpp
- 🎨 **Type Safe** - Complete type system mapping

## Installation

```bash
go get github.com/go-mixed/go-taichi
```

### Runtime Library

Download the Taichi C-API dynamic library from [Releases](https://github.com/go-mixed/go-taichi/releases):

**Required file**:
- Windows: `taichi_c_api.dll`
- Linux: `libtaichi_c_api.so`
- macOS: `libtaichi_c_api.dylib`

**Library placement**:

1. **Current Working Directory** - Place in project root (recommended for development)
2. **System PATH Directory** - Windows: `System32` or `%PATH%`; Linux/macOS: `/usr/local/lib` or `$LD_LIBRARY_PATH`/`$DYLD_LIBRARY_PATH` (recommended for production)
3. **Custom Directory** - Any directory, specify path via `NewRuntimeAuto("./lib")`

**Note**: `NewRuntimeAuto("")` searches current working directory first, then system PATH; `NewRuntimeAuto("./lib")` searches specified directory first, then system PATH.

### NO-NEED C Header Files

The C header files in `taichi/c_api/include/` are reference files used for generating Go API bindings. They are not required at runtime.

## Quick Start

Complete AOT kernel example demonstrating core features:

```go
package main

import (
    "fmt"
    "github.com/go-mixed/go-taichi/taichi"
)

func main() {
	
    // 1. Create runtime (auto-select best backend)
    // Option 1: Load from current directory or system PATH
    runtime, err := taichi.NewRuntimeAuto("")

    // Option 2: Load from custom directory (e.g., "./lib", "/opt/taichi")
    // runtime, err := taichi.NewRuntimeAuto("./lib")

    if err != nil {
        panic(err)
    }
    defer runtime.Release()
    fmt.Printf("Backend: %s\n", runtime.ArchName())

    // 2. Load precompiled AOT module
    module, err := taichi.LoadAotModule(runtime, "./module.tcm")
    if err != nil {
        panic(err)
    }
    defer module.Release()

    // 3. Get kernel from module
    kernel, err := module.GetKernel("add_kernel")
    if err != nil {
        panic(err)
    }

    // 4. Create input/output arrays
    size := uint32(1000)
    a, _ := taichi.NewNdArray1D(runtime, size, taichi.DataTypeF32)
    b, _ := taichi.NewNdArray1D(runtime, size, taichi.DataTypeF32)
    c, _ := taichi.NewNdArray1D(runtime, size, taichi.DataTypeF32)
    defer a.Release()
    defer b.Release()
    defer c.Release()

    // 5. Fill input data
    dataA, _ := a.AsSliceFloat32()
    dataB, _ := b.AsSliceFloat32()
    for i := range dataA {
        dataA[i] = float32(i)
        dataB[i] = float32(i) * 2
    }
    a.Unmap()
    b.Unmap()

    // 6. Execute kernel: c = a + b
    kernel.Launch().
        ArgNdArray(a).
        ArgNdArray(b).
        ArgNdArray(c).
        Run()

    // 7. Read results
    dataC, _ := c.AsSliceFloat32()
    fmt.Printf("Results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
        dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
    c.Unmap()
}
```

**Output**:
```
Backend: Vulkan
Results: [0.0, 3.0, 6.0, 9.0, 12.0]
```

This example demonstrates:
- ✅ Runtime creation and backend selection
- ✅ AOT module loading
- ✅ NdArray creation and data access
- ✅ Kernel execution with builder pattern
- ✅ Automatic resource management with `defer`

## Backend Support

| Backend | Windows | Linux | macOS | Recommended |
|---------|---------|-------|-------|-------------|
| Vulkan | ✅ | ✅ | ✅ | ⭐⭐⭐⭐⭐ |
| CUDA | ✅ | ✅ | ❌ | ⭐⭐⭐⭐ |
| CPU (x64) | ✅ | ✅ | ✅ | ⭐⭐⭐ |
| CPU (ARM64) | ❌ | ✅ | ✅ | ⭐⭐⭐ |
| Metal | ❌ | ❌ | ✅ | ⭐⭐⭐⭐ |
| OpenGL | ✅ | ✅ | ✅ | ⭐⭐ |

**Recommendation**: Use Vulkan (best cross-platform) or CUDA (NVIDIA GPU).

## Important Notes

### Runtime Files

Taichi C-API internal backends require runtime files (`.bc` files) which **must be located in the directory specified by `TI_LIB_DIR` environment variable**.

1. **Directory Structure**:

   ```
   your_project/
   └── lib/
       ├── windows/
       │   ├── taichi_c_api.dll
       │   ├── taichi_c_api.lib
       │   ├── runtime_x64.bc
       │   ├── runtime_cuda.bc
       │   ├── runtime_dx12.bc
       │   └── slim_libdevice.10.bc
       ├── linux/
       │   ├── libtaichi_c_api.so
       │   ├── runtime_x64.bc
       │   ├── runtime_cuda.bc
       │   └── slim_libdevice.10.bc
       └── darwin/
           ├── libtaichi_c_api.dylib
           ├── libMoltenVK.dylib
           └── runtime_arm64.bc
   ```

2. **Set TI_LIB_DIR environment variable** (required for all backends):

```powershell
# Windows PowerShell (run before your Go program)
$env:TI_LIB_DIR = "C:\path\to\your\project\lib\windows"
go run your_program.go
```

```bash
# Linux
export TI_LIB_DIR=/path/to/your/project/lib/linux
go run your_program.go
```

```bash
# macOS
export TI_LIB_DIR=/path/to/your/project/lib/darwin
go run your_program.go
```

**Note**: `TI_LIB_DIR` must point to the platform-specific directory containing `.bc` files. The dynamic library should be in the same directory or in the system PATH.

## Examples

See [examples/](examples/) directory:

| Example | Description | Level |
|---------|-------------|-------|
| `01_basic.go` | Basic runtime and memory | ⭐ |
| `02_ndarray.go` | N-dimensional arrays | ⭐ |
| `03_image.go` | Image processing | ⭐⭐ |
| `10_aot_kernel.go` | AOT kernel execution | ⭐⭐ |
| `11_aot_async.go` | Async execution | ⭐⭐⭐ |
| `12_aot_batch.go` | Batch execution | ⭐⭐⭐ |

## Documentation

### API Reference

- **[High-Level API](docs/high_api/SUMMARY.md)** - Recommended for most users
  - Automatic resource management
  - Type-safe operations
  - Go-idiomatic interfaces

- **[Low-Level API](docs/low_api/SUMMARY.md)** - For advanced users
  - Direct C-API bindings
  - Fine-grained control
  - Performance-critical operations

### Examples

See [examples/](examples/) directory for complete working examples.

## License

Apache 2.0 License - see [LICENSE](LICENSE) file.

## Acknowledgments

- [Taichi](https://github.com/taichi-dev/taichi) - High-performance parallel programming language
- [purego](https://github.com/ebitengine/purego) - Pure Go FFI library

---

**Version**: v1.0.0 | **Taichi**: v1.7.4 | **Go**: 1.25+ | **Updated**: 2026-02-19
