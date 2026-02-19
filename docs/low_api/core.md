# Core API Reference

Low-level C-API bindings for Taichi.

## Initialization

### Init

```go
func Init() error
```

Initialize Taichi C-API and load the dynamic library.

**Returns**: `error` - Error if initialization fails

**Example**:
```go
if err := c_api.Init(); err != nil {
    panic(err)
}
```

**Note**: Must be called before any other API functions.

---

## Version

### GetVersion

```go
func GetVersion() uint32
```

Get Taichi C-API version number.

**Returns**: `uint32` - Version (e.g., 1007000 for v1.7.0)

**Example**:
```go
version := c_api.GetVersion()
fmt.Printf("Taichi version: %d\n", version)
```

---

## Architecture

### GetAvailableArchs

```go
func GetAvailableArchs() []TiArch
```

Get list of available compute backends.

**Returns**: `[]TiArch` - Array of available architectures

**Example**:
```go
archs := c_api.GetAvailableArchs()
for _, arch := range archs {
    fmt.Printf("Available: %v\n", arch)
}
```

**Common Architectures**:
- `TI_ARCH_VULKAN` - Vulkan (recommended)
- `TI_ARCH_CUDA` - NVIDIA CUDA
- `TI_ARCH_X64` - x64 CPU
- `TI_ARCH_ARM64` - ARM64 CPU
- `TI_ARCH_METAL` - Apple Metal
- `TI_ARCH_OPENGL` - OpenGL

---

## Runtime

### CreateRuntime

```go
func CreateRuntime(arch TiArch, deviceIndex uint32) TiRuntime
```

Create a Taichi runtime instance.

**Parameters**:
- `arch` - Architecture to use
- `deviceIndex` - Device index (usually 0)

**Returns**: `TiRuntime` - Runtime handle

**Example**:
```go
archs := c_api.GetAvailableArchs()
runtime := c_api.CreateRuntime(archs[0], 0)
defer c_api.DestroyRuntime(runtime)
```

### DestroyRuntime

```go
func DestroyRuntime(runtime TiRuntime)
```

Destroy a runtime instance and free resources.

**Parameters**:
- `runtime` - Runtime handle to destroy

### Wait

```go
func Wait(runtime TiRuntime)
```

Wait for all submitted tasks to complete.

**Parameters**:
- `runtime` - Runtime handle

**Example**:
```go
// Submit async tasks...
c_api.Wait(runtime) // Wait for completion
```

---

## Error Handling

### GetLastError

```go
func GetLastError() (uint64, string)
```

Get the last error code and message.

**Returns**:
- `uint64` - Error code
- `string` - Error message

**Example**:
```go
runtime := c_api.CreateRuntime(arch, 0)
if runtime == c_api.TI_NULL_HANDLE {
    code, msg := c_api.GetLastError()
    fmt.Printf("Error %d: %s\n", code, msg)
}
```

### SetLastError

```go
func SetLastError(errCode TiError, message string)
```

Set error code and message (for internal use).

**Parameters**:
- `errCode` - Error code
- `message` - Error message
