# AOT API Reference

Ahead-of-Time compiled module management.

## Module Management

### LoadAotModule

```go
func LoadAotModule(runtime TiRuntime, modulePath string) TiAotModule
```

Load precompiled AOT module (.tcm file).

**Parameters**:
- `runtime` - Runtime handle
- `modulePath` - Path to .tcm file

**Returns**: `TiAotModule` - Module handle

**Example**:
```go
module := c_api.LoadAotModule(runtime, "./module.tcm")
defer c_api.DestroyAotModule(module)
```

### DestroyAotModule

```go
func DestroyAotModule(aotModule TiAotModule)
```

Destroy AOT module and free resources.

---

## Kernel Execution

### GetAotModuleKernel

```go
func GetAotModuleKernel(aotModule TiAotModule, name string) TiKernel
```

Get kernel from module by name.

**Parameters**:
- `aotModule` - Module handle
- `name` - Kernel name

**Returns**: `TiKernel` - Kernel handle

**Example**:
```go
kernel := c_api.GetAotModuleKernel(module, "my_kernel")
```

### LaunchKernel

```go
func LaunchKernel(runtime TiRuntime, kernel TiKernel, args []TiArgument)
```

Execute kernel with arguments.

**Parameters**:
- `runtime` - Runtime handle
- `kernel` - Kernel handle
- `args` - Array of arguments

**Example**:
```go
args := []c_api.TiArgument{
    {Type: c_api.TI_ARGUMENT_TYPE_NDARRAY, Value: ndarray1},
    {Type: c_api.TI_ARGUMENT_TYPE_NDARRAY, Value: ndarray2},
    {Type: c_api.TI_ARGUMENT_TYPE_I32, Value: int32(42)},
}
c_api.LaunchKernel(runtime, kernel, args)
```

---

## Compute Graph

### GetAotModuleComputeGraph

```go
func GetAotModuleComputeGraph(aotModule TiAotModule, name string) TiComputeGraph
```

Get compute graph from module.

**Parameters**:
- `aotModule` - Module handle
- `name` - Graph name

**Returns**: `TiComputeGraph` - Graph handle

### LaunchComputeGraph

```go
func LaunchComputeGraph(runtime TiRuntime, computeGraph TiComputeGraph, args []TiNamedArgument)
```

Execute compute graph with named arguments.

**Parameters**:
- `runtime` - Runtime handle
- `computeGraph` - Graph handle
- `args` - Array of named arguments

**Example**:
```go
args := []c_api.TiNamedArgument{
    {Name: "input", Argument: c_api.TiArgument{Type: c_api.TI_ARGUMENT_TYPE_NDARRAY, Value: input}},
    {Name: "output", Argument: c_api.TiArgument{Type: c_api.TI_ARGUMENT_TYPE_NDARRAY, Value: output}},
}
c_api.LaunchComputeGraph(runtime, graph, args)
```
