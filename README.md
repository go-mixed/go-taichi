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

## Important Notes

### Runtime Files

Taichi C-API internal backends require runtime files which **must be located in the directory specified by `TI_LIB_DIR` environment variable**.

1. **Directory Structure**:

   ```
   your_project/
   └── lib/
       ├── windows/
       │   ├── taichi_c_api.dll
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
   - download "runtime.7z" from [Taichi GitHub Release](https://github.com/go-mixed/go-taichi/releases) , and extract it to your project directory. 
   - Keep only the relevant directory for your operating system

2. **Set TI_LIB_DIR environment variable** (required for all backends):

```powershell
# Windows PowerShell
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

## Quick Start

Complete AOT kernel example demonstrating core features:

```go
package main

import (
    "fmt"
    "os"
    "github.com/go-mixed/go-taichi/taichi"
)

func main() {

    // 1. Create runtime (auto-select best backend)
    runtime, err := taichi.NewRuntimeAuto()
    if err != nil {
        panic(err)
    }
    defer runtime.Release()
    fmt.Printf("Backend: %s\n", runtime.ArchName())

    // 2. Load precompiled AOT module from .tcm file
    tcmData, err := os.ReadFile("./module.tcm")
    if err != nil {
        panic(err)
    }
    module, err := taichi.LoadAotModule(runtime, tcmData)
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
    taichi.MapNdArray(func(arrays ...taichi.NdArrayPtr) error {
        dataA := arrays[0].AsFloat32()
        dataB := arrays[1].AsFloat32()
        for i := range dataA {
            dataA[i] = float32(i)
            dataB[i] = float32(i) * 2
        }
        return nil
    }, a, b)

    // 6. Execute kernel: c = a + b
    kernel.Launch().
        ArgNdArray(a).
        ArgNdArray(b).
        ArgNdArray(c).
        Run()

    // 7. Read results
    c.MapFloat32(func(dataC []float32) error {
        fmt.Printf("Results: [%.1f, %.1f, %.1f, %.1f, %.1f]\n",
            dataC[0], dataC[1], dataC[2], dataC[3], dataC[4])
        return nil
    })
}
```

**Output**:
```
Backend: Vulkan
Results: [0.0, 3.0, 6.0, 9.0, 12.0]
```

This example demonstrates:
- ✅ Runtime creation and backend selection
- ✅ AOT module loading from .tcm file
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

## Examples

See [examples/](examples/) directory for complete working examples.

| Example | Feature | Level |
|---------|---------|-------|
| `01_runtime.go` | Runtime creation and management | ⭐ |
| `02_ndarray_1d.go` | 1D array operations | ⭐ |
| `03_ndarray_2d.go` | 2D matrix operations | ⭐ |
| `04_image.go` | Image operations | ⭐ |
| `10_aot_kernel.go` | AOT kernel execution | ⭐⭐ |
| `11_aot_async.go` | Async kernel execution | ⭐⭐ |
| `12_aot_batch.go` | Batch kernel execution | ⭐⭐ |
| `13_compute_graph.go` | Compute graph execution | ⭐⭐⭐ |
| `20_memory_cpu.go` | CPU memory import | ⭐⭐⭐ |
| `21_memory_cuda.go` | CUDA memory import | ⭐⭐⭐⭐ |

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

## License

Apache 2.0 License - see [LICENSE](LICENSE) file.

## Acknowledgments

- [Taichi](https://github.com/taichi-dev/taichi) - High-performance parallel programming language
- [purego](https://github.com/ebitengine/purego) - Pure Go FFI library

---

**Version**: v1.0.0 | **Taichi**: v1.7.4 | **Go**: 1.25+ | **Updated**: 2026-03-23
