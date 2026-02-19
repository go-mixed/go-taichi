# Low-Level API Reference

Direct C-API bindings for advanced users.

## Overview

The low-level API provides direct access to Taichi C-API functions. Use these when you need:
- Fine-grained control over resources
- Custom memory management
- Direct C-API compatibility
- Performance-critical operations

**Note**: Most users should use the [High-Level API](../high_api/) instead.

## Modules

### Core API

Runtime initialization, version info, and error handling.

**Key Functions**:
- `Init(libDir string)` - Initialize C-API
- `GetVersion()` - Get version
- `GetAvailableArchs()` - List backends
- `CreateRuntime()` / `DestroyRuntime()` - Runtime management
- `Wait()` - Wait for tasks
- `GetLastError()` - Error handling

[Full Documentation →](core.md)

---

### Memory API

Low-level memory allocation and management.

**Key Functions**:
- `AllocateMemory()` / `FreeMemory()` - Memory allocation
- `MapMemory()` / `UnmapMemory()` - Host access
- `CopyMemory*()` - Memory transfer
- `AllocateNdArray()` / `FreeNdArray()` - Array allocation

[Full Documentation →](memory.md)

---

### Image API

Image and texture operations.

**Key Functions**:
- `AllocateImage()` / `FreeImage()` - Image allocation
- `TransitionImage()` - Layout transitions
- `CopyImageDeviceToDevice()` - Image copy
- `CreateSampler()` / `DestroySampler()` - Sampler management

[Full Documentation →](image.md)

---

### AOT API

Ahead-of-time compiled module management.

**Key Functions**:
- `LoadAotModule()` / `DestroyAotModule()` - Module management
- `GetAotModuleKernel()` - Get kernel
- `LaunchKernel()` - Execute kernel
- `GetAotModuleComputeGraph()` - Get compute graph
- `LaunchComputeGraph()` - Execute graph

[Full Documentation →](aot.md)

---

## Quick Reference

| Module | Purpose | Key Operations |
|--------|---------|----------------|
| Core | Initialization | Runtime, version, errors |
| Memory | Allocation | Memory, NdArray |
| Image | Textures | Images, samplers |
| AOT | Execution | Kernels, graphs |

---

## Usage Pattern

```go
// 1. Initialize
c_api.Init()

// 2. Create runtime
archs := c_api.GetAvailableArchs()
runtime := c_api.CreateRuntime(archs[0], 0)
defer c_api.DestroyRuntime(runtime)

// 3. Allocate resources
allocInfo := c_api.TiMemoryAllocateInfo{...}
memory := c_api.AllocateMemory(runtime, &allocInfo)
defer c_api.FreeMemory(runtime, memory)

// 4. Use resources
ptr := c_api.MapMemory(runtime, memory)
// ... use ptr ...
c_api.UnmapMemory(runtime, memory)
```

---

## See Also

- [High-Level API](../high_api/) - Recommended for most users
- [Examples](../../examples/) - Complete examples
