# AOT Module API

Precompiled kernel and compute graph management.

## Module Functions

### LoadAotModule

```go
func LoadAotModule(runtime *Runtime, tcmData []byte) (*AotModule, error)
```

Load precompiled AOT module from .tcm file data.

**Parameters**:
- `runtime` - Runtime instance
- `tcmData` - Raw bytes of .tcm file (read with `os.ReadFile`)

**Returns**:
- `*AotModule` - Module instance
- `error` - Error if invalid

**Example**:
```go
tcmData, err := os.ReadFile("./module.tcm")
if err != nil {
    panic(err)
}
module, err := taichi.LoadAotModule(runtime, tcmData)
if err != nil {
    panic(err)
}
defer module.Release()
```

---

### LoadAotModuleFile

```go
func LoadAotModuleFile(runtime *Runtime, moduleDir string) (*AotModule, error)
```

Load precompiled AOT module from directory containing `metadata.json`.

**Parameters**:
- `runtime` - Runtime instance
- `moduleDir` - Directory path containing extracted TCM files and metadata.json

**Returns**:
- `*AotModule` - Module instance
- `error` - Error if invalid

---

## Module Methods

### GetKernel

```go
func (m *AotModule) GetKernel(name string) (*Kernel, error)
```

Get kernel by name from module.

**Parameters**:
- `name` - Kernel name (as defined in Python)

**Returns**:
- `*Kernel` - Kernel instance
- `error` - Error if kernel not found

**Example**:
```go
kernel, err := module.GetKernel("add_kernel")
if err != nil {
    panic(err)
}
```

---

### GetComputeGraph

```go
func (m *AotModule) GetComputeGraph(name string) (*ComputeGraph, error)
```

Get compute graph by name from module.

**Parameters**:
- `name` - Graph name

**Returns**:
- `*ComputeGraph` - Graph instance
- `error` - Error if graph not found

**Example**:
```go
graph, err := module.GetComputeGraph("my_graph")
if err != nil {
    panic(err)
}
```

---

### Release

```go
func (m *AotModule) Release()
```

Free module resources.

**Example**:
```go
module, _ := taichi.LoadAotModule(runtime, tcmData)
defer module.Release()
```

---

## Kernel Execution

### Launch

```go
func (k *Kernel) Launch() *KernelLauncher
```

Start kernel execution builder.

**Returns**: `*KernelLauncher` - Builder for adding arguments

---

### KernelLauncher Methods

#### ArgNdArray

```go
func (kl *KernelLauncher) ArgNdArray(arr *NdArray) *KernelLauncher
```

Add NdArray argument.

---

#### ArgInt32

```go
func (kl *KernelLauncher) ArgInt32(value int32) *KernelLauncher
```

Add int32 scalar argument.

---

#### ArgFloat32

```go
func (kl *KernelLauncher) ArgFloat32(value float32) *KernelLauncher
```

Add float32 scalar argument.

---

#### Run

```go
func (kl *KernelLauncher) Run()
```

Execute kernel synchronously (blocks until complete).

---

#### RunAsync

```go
func (kl *KernelLauncher) RunAsync()
```

Execute kernel asynchronously (returns immediately).

---

## Complete Example

```go
// Load module from .tcm file
tcmData, _ := os.ReadFile("./module.tcm")
module, _ := taichi.LoadAotModule(runtime, tcmData)
defer module.Release()

// Get kernel
kernel, _ := module.GetKernel("add_kernel")

// Prepare data
a, _ := taichi.NewNdArray1D(runtime, 100, taichi.DataTypeF32)
b, _ := taichi.NewNdArray1D(runtime, 100, taichi.DataTypeF32)
c, _ := taichi.NewNdArray1D(runtime, 100, taichi.DataTypeF32)
defer a.Release()
defer b.Release()
defer c.Release()

// Fill input data
taichi.MapNdArray(func(arrays ...taichi.NdArrayPtr) error {
    dataA := arrays[0].AsFloat32()
    dataB := arrays[1].AsFloat32()
    for i := range dataA {
        dataA[i] = float32(i)
        dataB[i] = float32(i) * 2
    }
    return nil
}, a, b)

// Execute kernel: c = a + b
kernel.Launch().
    ArgNdArray(a).
    ArgNdArray(b).
    ArgNdArray(c).
    Run()

// Read results
c.MapFloat32(func(dataC []float32) error {
    fmt.Printf("Result: %v\n", dataC[:5])
    return nil
})
```

---

## Generating TCM Files

Use Python to compile Taichi kernels to .tcm files:

```python
import taichi as ti

ti.init(arch=ti.vulkan)

@ti.kernel
def add_kernel(
    a: ti.types.ndarray(dtype=ti.f32, ndim=1),
    b: ti.types.ndarray(dtype=ti.f32, ndim=1),
    c: ti.types.ndarray(dtype=ti.f32, ndim=1),
):
    for i in c:
        c[i] = a[i] + b[i]

# Create AOT module
m = ti.aot.Module(ti.vulkan)
m.add_kernel(add_kernel)
m.archive("module.tcm")
```
