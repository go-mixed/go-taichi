# High-Level API Reference

High-level abstractions with automatic resource management.

## Overview

The high-level API provides convenient wrappers around the low-level C-API, featuring:
- Automatic resource cleanup with `defer Release()`
- Type-safe operations
- Simplified error handling
- Go-idiomatic interfaces

## Core Types

### Runtime

Runtime management and backend selection.

**Key Methods**:
- `NewRuntimeAuto()` - Auto-select best backend
- `NewRuntime(arch)` - Manual backend selection
- `Arch()` - Get backend type
- `ArchName()` - Get backend name
- `Wait()` - Wait for all tasks
- `Release()` - Free resources

**Example**:
```go
runtime, _ := taichi.NewRuntimeAuto()
defer runtime.Release()
fmt.Printf("Backend: %s\n", runtime.ArchName())
```

[Full Documentation →](runtime.md)

---

### NdArray

N-dimensional array with automatic memory management.

**Key Methods**:
- `NewNdArray1D/2D/3D()` - Create arrays
- `AsSliceFloat32/Int32/...()` - Access as Go slice
- `Unmap()` - Unmap memory
- `Shape()` - Get dimensions
- `ElemCount()` - Get element count
- `Release()` - Free resources

**Example**:
```go
arr, _ := taichi.NewNdArray1D(runtime, 1000, taichi.DataTypeF32)
defer arr.Release()

data, _ := arr.AsSliceFloat32()
data[0] = 3.14
arr.Unmap()
```

[Full Documentation →](ndarray.md)

---

### Image

Image and texture processing.

**Key Methods**:
- `NewImage2D/3D()` - Create images
- `TransitionLayout()` - Change layout
- `CopyTo()` - Copy to another image
- `Width/Height/Depth()` - Get dimensions
- `Format()` - Get pixel format
- `Release()` - Free resources

**Example**:
```go
img, _ := taichi.NewImage2D(runtime, 512, 512, taichi.FormatRGBA8)
defer img.Release()

img.TransitionLayout(taichi.ImageLayoutShaderRead)
```

[Full Documentation →](image.md)

---

### AotModule

Precompiled kernel execution.

**Key Methods**:
- `LoadAotModule()` - Load .tcm file
- `GetKernel()` - Get kernel by name
- `GetComputeGraph()` - Get compute graph
- `Release()` - Free resources

**Example**:
```go
module, _ := taichi.LoadAotModule(runtime, "./module.tcm")
defer module.Release()

kernel, _ := module.GetKernel("my_kernel")
```

[Full Documentation →](aot.md)

---

### Kernel

Kernel execution with builder pattern.

**Key Methods**:
- `Launch()` - Start builder
- `ArgNdArray()` - Add array argument
- `ArgInt32/Float32/...()` - Add scalar argument
- `Run()` - Execute synchronously
- `RunAsync()` - Execute asynchronously

**Example**:
```go
kernel.Launch().
    ArgNdArray(input).
    ArgNdArray(output).
    ArgInt32(42).
    Run()
```

[Full Documentation →](kernel.md)

---

## Quick Reference

| Type | Purpose | Key Feature |
|------|---------|-------------|
| `Runtime` | Backend management | Auto-selection |
| `NdArray` | N-D arrays | Go slice access |
| `Image` | Textures | Layout management |
| `AotModule` | Kernel loading | .tcm files |
| `Kernel` | Execution | Builder pattern |

---

## See Also

- [Low-Level API](../low_api/) - C-API bindings
- [Examples](../../examples/) - Complete examples
