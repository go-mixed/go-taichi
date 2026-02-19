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

### AsSliceFloat32

```go
func (a *NdArray) AsSliceFloat32() ([]float32, error)
```

Map array as Go float32 slice.

**Returns**:
- `[]float32` - Go slice view of array data
- `error` - Error if type mismatch or map fails

**Example**:
```go
data, _ := arr.AsSliceFloat32()
for i := range data {
    data[i] = float32(i) * 0.5
}
arr.Unmap() // Must unmap after use
```

---

### AsSliceInt32

```go
func (a *NdArray) AsSliceInt32() ([]int32, error)
```

Map array as Go int32 slice.

---

### AsSliceUint32

```go
func (a *NdArray) AsSliceUint32() ([]uint32, error)
```

Map array as Go uint32 slice.

---

### Unmap

```go
func (a *NdArray) Unmap()
```

Unmap previously mapped memory. **Must be called** after accessing data.

**Example**:
```go
data, _ := arr.AsSliceFloat32()
// Use data...
arr.Unmap() // Required!
```

---

## Query Methods

### Shape

```go
func (a *NdArray) Shape() []uint32
```

Get array dimensions.

**Returns**: `[]uint32` - Array shape

**Example**:
```go
shape := arr.Shape() // [1000] for 1D, [100, 100] for 2D
```

---

### ElemCount

```go
func (a *NdArray) ElemCount() uint32
```

Get total number of elements.

**Returns**: `uint32` - Element count

---

### DataType

```go
func (a *NdArray) DataType() DataType
```

Get element data type.

**Returns**: `DataType` - Data type enum

---

### IsMapped

```go
func (a *NdArray) IsMapped() bool
```

Check if array is currently mapped.

**Returns**: `bool` - True if mapped

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
