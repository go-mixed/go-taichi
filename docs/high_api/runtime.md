# Runtime API

Runtime management and backend selection.

## Functions

### NewRuntimeAuto

```go
func NewRuntimeAuto() (*Runtime, error)
```

Automatically select the best available backend.

**Returns**:
- `*Runtime` - Runtime instance
- `error` - Error if no backend available

**Selection Priority**: Vulkan > CUDA > x64 > ARM64 > OpenGL

**Example**:
```go
runtime, err := taichi.NewRuntimeAuto()
if err != nil {
    panic(err)
}
defer runtime.Release()
```

---

### NewRuntime

```go
func NewRuntime(arch Arch) (*Runtime, error)
```

Create runtime with specific backend.

**Parameters**:
- `arch` - Backend architecture (e.g., `ArchVulkan`, `ArchCUDA`)

**Returns**:
- `*Runtime` - Runtime instance
- `error` - Error if backend unavailable

**Example**:
```go
runtime, err := taichi.NewRuntime(taichi.ArchVulkan)
if err != nil {
    panic(err)
}
defer runtime.Release()
```

---

## Methods

### Arch

```go
func (r *Runtime) Arch() Arch
```

Get backend architecture type.

**Returns**: `Arch` - Architecture enum

---

### ArchName

```go
func (r *Runtime) ArchName() string
```

Get backend name as string.

**Returns**: `string` - Backend name (e.g., "Vulkan", "CUDA")

**Example**:
```go
fmt.Printf("Using: %s\n", runtime.ArchName())
```

---

### Wait

```go
func (r *Runtime) Wait()
```

Wait for all submitted tasks to complete.

**Example**:
```go
// Submit async operations...
runtime.Wait() // Block until all complete
```

---

### Release

```go
func (r *Runtime) Release()
```

Free runtime resources. Must be called when done.

**Example**:
```go
runtime, _ := taichi.NewRuntimeAuto()
defer runtime.Release() // Automatic cleanup
```

---

## Available Architectures

| Constant | Description | Platform |
|----------|-------------|----------|
| `ArchVulkan` | Vulkan | All |
| `ArchCUDA` | NVIDIA CUDA | Windows, Linux |
| `ArchX64` | x64 CPU | All |
| `ArchARM64` | ARM64 CPU | Linux, macOS |
| `ArchMetal` | Apple Metal | macOS |
| `ArchOpenGL` | OpenGL | All |
