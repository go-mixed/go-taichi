# Memory API Reference

Low-level memory management functions.

## Memory Allocation

### AllocateMemory

```go
func AllocateMemory(runtime TiRuntime, allocInfo *TiMemoryAllocateInfo) TiMemory
```

Allocate device memory.

**Parameters**:
- `runtime` - Runtime handle
- `allocInfo` - Allocation parameters

**Returns**: `TiMemory` - Memory handle

**Example**:
```go
allocInfo := c_api.TiMemoryAllocateInfo{
    Size:      1024,
    HostWrite: c_api.TI_TRUE,
    HostRead:  c_api.TI_TRUE,
    Usage:     c_api.TiMemoryUsageFlags(c_api.TI_MEMORY_USAGE_STORAGE_BIT),
}
memory := c_api.AllocateMemory(runtime, &allocInfo)
defer c_api.FreeMemory(runtime, memory)
```

### FreeMemory

```go
func FreeMemory(runtime TiRuntime, memory TiMemory)
```

Free allocated memory.

---

## Memory Mapping

### MapMemory

```go
func MapMemory(runtime TiRuntime, memory TiMemory) unsafe.Pointer
```

Map memory to host-accessible pointer.

**Returns**: `unsafe.Pointer` - Host pointer

**Example**:
```go
ptr := c_api.MapMemory(runtime, memory)
// Use pointer...
c_api.UnmapMemory(runtime, memory)
```

### UnmapMemory

```go
func UnmapMemory(runtime TiRuntime, memory TiMemory)
```

Unmap previously mapped memory.

---

## Memory Copy

### CopyMemoryDeviceToDevice

```go
func CopyMemoryDeviceToDevice(runtime TiRuntime, dstMem TiMemory, dstOffset uint64,
                               srcMem TiMemory, srcOffset uint64, size uint64)
```

Copy between device buffers.

### CopyMemoryHostToDevice

```go
func CopyMemoryHostToDevice(runtime TiRuntime, dstMem TiMemory, dstOffset uint64,
                            srcPtr unsafe.Pointer, size uint64)
```

Copy from host to device.

### CopyMemoryDeviceToHost

```go
func CopyMemoryDeviceToHost(runtime TiRuntime, dstPtr unsafe.Pointer,
                            srcMem TiMemory, srcOffset uint64, size uint64)
```

Copy from device to host.

---

## NdArray

### AllocateNdArray

```go
func AllocateNdArray(runtime TiRuntime, allocInfo *TiNdArrayAllocateInfo) TiNdArray
```

Allocate N-dimensional array.

**Example**:
```go
allocInfo := c_api.TiNdArrayAllocateInfo{
    Shape:     []uint32{100, 100},
    ElemType:  c_api.TI_DATA_TYPE_F32,
}
array := c_api.AllocateNdArray(runtime, &allocInfo)
defer c_api.FreeNdArray(runtime, array)
```

### FreeNdArray

```go
func FreeNdArray(runtime TiRuntime, array TiNdArray)
```

Free N-dimensional array.
