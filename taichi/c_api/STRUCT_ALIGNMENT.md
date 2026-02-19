# Go 与 C 结构体对齐问题

## 问题概述

在使用 `purego` 进行 Go-C FFI 调用时，必须确保 Go 结构体的内存布局与 C 结构体**完全一致**。即使字段顺序和类型正确，对齐（alignment）、填充（padding）、字段顺序或缺失字段的差异都会导致 C 代码读取错误的内存位置。

---

## 常见问题类型

### 1. Union 的大小计算

**C 语言中的 Union：**
```c
typedef union TiArgumentValue {
  int32_t i32;              // 4 bytes
  float f32;                // 4 bytes
  TiNdArray ndarray;        // 152 bytes
  TiTexture texture;        // 40 bytes
  TiScalar scalar;          // 16 bytes
} TiArgumentValue;
```

Union 的大小等于**最大成员的大小**，在这个例子中是 `TiNdArray` (152 bytes)。

**错误的 Go 实现：**
```go
type TiArgumentValue struct {
    Data [512]byte  // ❌ 错误：随意设置的大小
}
```

**正确的 Go 实现：**
```go
type TiArgumentValue struct {
    Data [152]byte  // ✅ 正确：与 C union 的实际大小匹配
}
```

---

### 2. 结构体字段间的隐式 Padding

**C 结构体：**
```c
typedef struct TiArgument {
  TiArgumentType type;      // 4 bytes, offset 0
  // 隐式 padding: 4 bytes  // C 编译器自动插入
  TiArgumentValue value;    // 152 bytes, offset 8
} TiArgument;
// 总大小：160 bytes
```

C 编译器会自动插入 padding 以满足对齐要求。`TiArgumentValue` 包含 `TiNdArray`，其中 `Memory` 字段是 `uintptr`（8字节对齐），因此 C 编译器将 `value` 字段对齐到 8 字节边界。

**错误的 Go 实现：**
```go
type TiArgument struct {
    Type  TiArgumentType   // 4 bytes, offset 0
    Value TiArgumentValue  // 152 bytes, offset 4 ❌ 错误的 offset
}
// 总大小：156 bytes ❌
```

**正确的 Go 实现：**
```go
type TiArgument struct {
    Type  TiArgumentType   // 4 bytes, offset 0
    _     [4]byte          // 显式 padding
    Value TiArgumentValue  // 152 bytes, offset 8 ✅
}
// 总大小：160 bytes ✅
```

## 验证工具

### 使用完整验证工具

项目提供了完整的验证工具来检查所有结构体：

**运行 Go 验证：**
```bash
cd examples
go run 99_struct_alignment.go > go_output.txt
```

**编译并运行 C 验证：**
```bash
cd examples
gcc -I../taichi/c_api/include 99_struct_alignment.c -o verify_structs.exe
./verify_structs.exe > c_output.txt
```

**对比结果：**
```bash
diff go_output.txt c_output.txt
```

如果有差异，说明结构体定义不匹配，需要修复！


## 相关文件

- `taichi/c_api/types.go` - Go 结构体定义
- `taichi/c_api/include/taichi/taichi_core.h` - C 结构体定义（权威来源）
- `examples/99_struct_alignment.go` - Go 验证工具
- `examples/99_struct_alignment.c` - C 验证工具
- `STRUCT_MISMATCH_ISSUE.md` - 已发现问题的详细记录

## 对齐规则速查

### 基本规则

1. **字段对齐**：每个字段必须对齐到其自然对齐边界
   - `uint32`: 4 字节对齐
   - `uint64`: 8 字节对齐
   - `uintptr`: 8 字节对齐（64位系统）
   - `float32`: 4 字节对齐

2. **结构体对齐**：结构体的对齐要求 = 最大成员的对齐要求

3. **结构体大小**：结构体总大小必须是其对齐要求的倍数

### 示例

```go
// 示例 1: 需要 padding
type Example1 struct {
    A uint32   // 4 bytes, offset 0
    _  [4]byte // padding
    B uint64   // 8 bytes, offset 8
}
// 总大小: 16 bytes

// 示例 2: 不需要 padding
type Example2 struct {
    A uint64   // 8 bytes, offset 0
    B uint32   // 4 bytes, offset 8
    C uint32   // 4 bytes, offset 12
}
// 总大小: 16 bytes
```

---

## 参考资料

- [Go unsafe 包文档](https://pkg.go.dev/unsafe)
- [C 结构体对齐规则](https://en.cppreference.com/w/c/language/object#Alignment)
- [Purego FFI 文档](https://github.com/ebitengine/purego)
- [Taichi C-API 文档](https://docs.taichi-lang.org/docs/taichi_core)
