# Go and C Struct Alignment Issues

## Overview

When using `purego` for Go-C FFI calls, it is critical to ensure that Go struct memory layouts **exactly match** their C counterparts. Even with correct field order and types, differences in alignment, padding, field ordering, or missing fields will cause C code to read from incorrect memory locations.

---

## Common Issue Types

### 1. Union Size Calculation

**C Union:**
```c
typedef union TiArgumentValue {
  int32_t i32;              // 4 bytes
  float f32;                // 4 bytes
  TiNdArray ndarray;        // 152 bytes
  TiTexture texture;        // 40 bytes
  TiScalar scalar;          // 16 bytes
} TiArgumentValue;
```

Union size equals the **size of the largest member**, which in this case is `TiNdArray` (152 bytes).

**Incorrect Go Implementation:**
```go
type TiArgumentValue struct {
    Data [512]byte  // ❌ Wrong: arbitrary size
}
```

**Correct Go Implementation:**
```go
type TiArgumentValue struct {
    Data [152]byte  // ✅ Correct: matches actual C union size
}
```

---

### 2. Implicit Padding Between Struct Fields

**C Struct:**
```c
typedef struct TiArgument {
  TiArgumentType type;      // 4 bytes, offset 0
  // implicit padding: 4 bytes  // Automatically inserted by C compiler
  TiArgumentValue value;    // 152 bytes, offset 8
} TiArgument;
// Total size: 160 bytes
```

The C compiler automatically inserts padding to satisfy alignment requirements. `TiArgumentValue` contains `TiNdArray`, which has a `Memory` field of type `uintptr` (8-byte aligned), so the C compiler aligns the `value` field to an 8-byte boundary.

**Incorrect Go Implementation:**
```go
type TiArgument struct {
    Type  TiArgumentType   // 4 bytes, offset 0
    Value TiArgumentValue  // 152 bytes, offset 4 ❌ Wrong offset
}
// Total size: 156 bytes ❌
```

**Correct Go Implementation:**
```go
type TiArgument struct {
    Type  TiArgumentType   // 4 bytes, offset 0
    _     [4]byte          // Explicit padding
    Value TiArgumentValue  // 152 bytes, offset 8 ✅
}
// Total size: 160 bytes ✅
```

## Verification Tools

### Using Complete Verification Tools

The project provides complete verification tools to check all structs:

**Run Go Verification:**
```bash
cd examples
go run 99_struct_alignment.go > go_output.txt
```

**Compile and Run C Verification:**
```bash
cd examples
gcc -I../taichi/c_api/include 99_struct_alignment.c -o verify_structs.exe
./verify_structs.exe > c_output.txt
```

**Compare Results:**
```bash
diff go_output.txt c_output.txt
```

If there are differences, it means the struct definitions don't match and need to be fixed!


## Related Files

- `taichi/c_api/types.go` - Go struct definitions
- `taichi/c_api/include/taichi/taichi_core.h` - C struct definitions (authoritative source)
- `examples/99_struct_alignment.go` - Go verification tool
- `examples/99_struct_alignment.c` - C verification tool
- `STRUCT_MISMATCH_ISSUE.md` - Detailed record of discovered issues

## Alignment Rules Quick Reference

### Basic Rules

1. **Field Alignment**: Each field must be aligned to its natural alignment boundary
   - `uint32`: 4-byte aligned
   - `uint64`: 8-byte aligned
   - `uintptr`: 8-byte aligned (on 64-bit systems)
   - `float32`: 4-byte aligned

2. **Struct Alignment**: Struct alignment requirement = alignment requirement of largest member

3. **Struct Size**: Total struct size must be a multiple of its alignment requirement

### Examples

```go
// Example 1: Requires padding
type Example1 struct {
    A uint32   // 4 bytes, offset 0
    _  [4]byte // padding
    B uint64   // 8 bytes, offset 8
}
// Total size: 16 bytes

// Example 2: No padding needed
type Example2 struct {
    A uint64   // 8 bytes, offset 0
    B uint32   // 4 bytes, offset 8
    C uint32   // 4 bytes, offset 12
}
// Total size: 16 bytes
```

---

## References

- [Go unsafe Package Documentation](https://pkg.go.dev/unsafe)
- [C Struct Alignment Rules](https://en.cppreference.com/w/c/language/object#Alignment)
- [Purego FFI Documentation](https://github.com/ebitengine/purego)
- [Taichi C-API Documentation](https://docs.taichi-lang.org/docs/taichi_core)
