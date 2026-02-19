# Go-Taichi Examples

Complete working examples demonstrating Go-Taichi features.

## Examples Index

### Basic Examples (01-09)

| File | Feature | Level |
|------|---------|-------|
| `01_runtime.go` | Runtime creation and management | ⭐ |
| `02_ndarray_1d.go` | 1D array operations | ⭐ |
| `03_ndarray_2d.go` | 2D matrix operations | ⭐ |
| `04_image.go` | Image operations | ⭐ |

### AOT Examples (10-19)

| File | Feature | Level | Prerequisite |
|------|---------|-------|--------------|
| `10_aot_kernel.go` | AOT kernel execution | ⭐⭐ | AOT module |
| `11_aot_async.go` | Async kernel execution | ⭐⭐ | AOT module |
| `12_aot_batch.go` | Batch kernel execution | ⭐⭐ | AOT module |
| `13_compute_graph.go` | Compute graph execution | ⭐⭐⭐ | Compute graph module |

### Advanced Examples (20-29)

| File | Feature | Level |
|------|---------|-------|
| `20_memory_cpu.go` | CPU memory import | ⭐⭐⭐ |
| `21_memory_cuda.go` | CUDA memory import | ⭐⭐⭐⭐ |

---

## Quick Start

### 1. Basic Examples (No Prerequisites)

```bash
# Runtime management
go run ./examples/01_runtime.go

# 1D arrays
go run ./examples/02_ndarray_1d.go

# 2D matrices
go run ./examples/03_ndarray_2d.go

# Image operations
go run ./examples/04_image.go
```

### 2. AOT Examples (Requires AOT Module)

**Generate AOT module first**:

```bash
# Install Taichi
uv pip install taichi==1.7.4

# Generate AOT module
uv run ./examples/10_aot_kernel.py
```

**Run examples**:

```bash
# Basic kernel execution
go run ./examples/10_aot_kernel.go

# Async execution
go run ./examples/11_aot_async.go

# Batch execution
go run ./examples/12_aot_batch.go
```

### 3. Advanced Examples

```bash
# CPU memory import
go run ./examples/20_memory_cpu.go

# CUDA memory import (concept demo)
go run ./examples/21_memory_cuda.go
```

---

## Learning Path

### Beginner Path

1. `01_runtime.go` - Understand runtime creation
2. `02_ndarray_1d.go` - Learn 1D array operations
3. `03_ndarray_2d.go` - Learn 2D matrix operations
4. `04_image.go` - Learn image operations

### Intermediate Path

5. `10_aot_kernel.go` - Learn AOT kernel basics
6. `11_aot_async.go` - Learn async execution
7. `12_aot_batch.go` - Learn batch optimization

### Advanced Path

8. `13_compute_graph.go` - Learn compute graphs
9. `20_memory_cpu.go` - Learn memory import
10. `21_memory_cuda.go` - Learn CUDA integration

---

## Generating AOT Modules

### Basic AOT Module

The `10_aot_kernel.py` script generates a basic AOT module with kernels:

```bash
python ./examples/10_aot_kernel.py
```

This creates `aot_module.tcm` containing:
- `add_kernel` - Vector addition
- `add_and_scale_kernel` - Addition with scaling

### Custom Kernels

To create your own kernels, write a Python script:

```python
import taichi as ti

ti.init(arch=ti.vulkan)

@ti.kernel
def my_kernel(
    a: ti.types.ndarray(dtype=ti.f32, ndim=1),
    b: ti.types.ndarray(dtype=ti.f32, ndim=1),
):
    for i in a:
        b[i] = a[i] * 2.0

# Export AOT module
m = ti.aot.Module(ti.vulkan)
m.add_kernel(my_kernel)
m.archive("my_module.tcm")
```

---

## Naming Convention

- **Number Prefix**: Indicates difficulty and learning order
  - `01-09`: Basic features
  - `10-19`: AOT features
  - `20-29`: Advanced features
  - `90+`: Tools and diagnostics

- **Feature Name**: Clear description of single feature
  - `runtime` - Runtime management
  - `ndarray_1d` - 1D arrays
  - `aot_kernel` - AOT kernels
  - `memory_cpu` - CPU memory import

---

## Common Issues

### AOT Module Not Found

Ensure `aot_module.tcm` exists in the examples directory:

```bash
ls ./examples/aot_module.tcm
```

If missing, generate it:

```bash
python ./examples/10_aot_kernel.py
```

### Kernel Not Found

Ensure the kernel name in Go matches the Python definition:

```python
# Python
@ti.kernel
def add_kernel(...):  # Name: add_kernel
    pass
```

```go
// Go
kernel, _ := module.GetKernel("add_kernel")  // Same name
```

### Argument Type Mismatch

Ensure argument order and types match exactly:

```python
# Python
@ti.kernel
def kernel(arr: ti.types.ndarray(), value: ti.f32):
    pass
```

```go
// Go - Order and types must match
kernel.Launch().
    ArgNdArray(arr).     // 1st: ndarray
    ArgFloat32(value).   // 2nd: f32
    Run()
```

---

## Design Principles

1. **One Example, One Feature** - Each file demonstrates one core feature
2. **Self-Documenting Names** - Understand content from filename
3. **Numbered Ordering** - Organized by difficulty and learning sequence
4. **Standalone Runnable** - Each example is a complete main program
5. **Clear Comments** - Code begins with feature description and key points

---

## See Also

- [High-Level API Documentation](../docs/high_api/)
- [Low-Level API Documentation](../docs/low_api/)
- [Taichi AOT Documentation](https://docs.taichi-lang.org/docs/aot)
