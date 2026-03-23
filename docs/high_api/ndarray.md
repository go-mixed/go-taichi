# NdArray API

N-dimensional array with automatic memory management.

## Creation Functions

### NewNdArray1D

```go
func NewNdArray1D(runtime *Runtime, size uint32, dtype DataType) (*NdArray, error)
```

Create 1D array.

**Parameters**:
- `runtime` - Runtime instance
- `size` - Number of elements
- `dtype` - Data type (e.g., `DataTypeF32`, `DataTypeI32`)

**Returns**:
- `*NdArray` - Array instance
- `error` - Error if allocation fails

**Example**:
```go
arr, err := taichi.NewNdArray1D(runtime, 1000, taichi.DataTypeF32)
if err != nil {
    panic(err)
}
defer arr.Release()
```

---

### NewNdArray2D

```go
func NewNdArray2D(runtime *Runtime, width, height uint32, dtype DataType) (*NdArray, error)
```

Create 2D array.

**Parameters**:
- `runtime` - Runtime instance
- `width` - Width dimension
- `height` - Height dimension
- `dtype` - Data type

**Example**:
```go
matrix, _ := taichi.NewNdArray2D(runtime, 100, 100, taichi.DataTypeF32)
defer matrix.Release()
```

---

### NewNdArray3D

```go
func NewNdArray3D(runtime *Runtime, width, height, depth uint32, dtype DataType) (*NdArray, error)
```

Create 3D array.

**Parameters**:
- `runtime` - Runtime instance
- `width`, `height`, `depth` - Dimensions
- `dtype` - Data type

**Example**:
```go
volume, _ := taichi.NewNdArray3D(runtime, 64, 64, 64, taichi.DataTypeU8)
defer volume.Release()
```

---

## Data Access Methods

### MapFloat32

```go
func (a *NdArray) MapFloat32(f func(data []float32) error) error
```

Map array as Go float32 slice and execute function.

**Parameters**:
- `f` - Function that receives the data slice and returns error

**Example**:
```go
err := arr.MapFloat32(func(data []float32) error {
    for i := range data {
        data[i] = float32(i) * 0.5
    }
    return nil
})
```

---

### MapNdArray

```go
func MapNdArray(fn func(arrays ...NdArrayPtr) error, arrays ...*NdArray) error
```

Map multiple NdArrays and execute function. All arrays must be float32 type.

**Parameters**:
- `fn` - Function that receives NdArrayPtr slices
- `arrays` - NdArray instances to map

**Example**:
```go
taichi.MapNdArray(func(arrays ...taichi.NdArrayPtr) error {
    dataA := arrays[0].AsFloat32()
    dataB := arrays[1].AsFloat32()
    for i := range dataA {
        dataA[i] = float32(i)
        dataB[i] = float32(i) * 2
    }
    return nil
}, a, b)
```

---

## Query Methods

### Shape

```go
func (a *NdArray) Shape() NdShape
```

Get array dimensions.

**Returns**: `NdShape` - Array shape (alias for `[]uint32`)

**Example**:
```go
shape := arr.Shape() // [1000] for 1D, [100, 100] for 2D
```

---

### TotalElements

```go
func (a *NdArray) TotalElements() uint64
```

Get total number of elements (including elemShape).

**Returns**: `uint64` - Total element count

---

### ElemType

```go
func (a *NdArray) ElemType() DataType
```

Get element data type.

**Returns**: `DataType` - Data type enum

---

### Release

```go
func (a *NdArray) Release()
```

Free array resources.

**Example**:
```go
arr, _ := taichi.NewNdArray1D(runtime, 1000, taichi.DataTypeF32)
defer arr.Release()
```

---

## Data Types

| Constant | Go Type | Description |
|----------|---------|-------------|
| `DataTypeF32` | `float32` | 32-bit float |
| `DataTypeF64` | `float64` | 64-bit float |
| `DataTypeI32` | `int32` | 32-bit signed int |
| `DataTypeI64` | `int64` | 64-bit signed int |
| `DataTypeU32` | `uint32` | 32-bit unsigned int |
| `DataTypeU64` | `uint64` | 64-bit unsigned int |
