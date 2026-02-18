# Go 与 C 结构体对齐问题

## 问题概述

在使用 `purego` 进行 Go-C FFI 调用时，必须确保 Go 结构体的内存布局与 C 结构体**完全一致**。即使字段顺序和类型正确，对齐（alignment）和填充（padding）的差异也会导致 C 代码读取错误的内存位置。

## 关键问题

### 1. Union 的大小计算

**C 语言中的 Union：**
```c
typedef union TiArgumentValue {
  int32_t i32;              // 4 bytes
  float f32;                // 4 bytes
  TiNdArray ndarray;        // 152 bytes
  TiTexture texture;        // 40 bytes
  TiScalar scalar;          // 16 bytes
  TiTensor tensor;          // 144 bytes
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

## 诊断方法

### 使用 C 编译器验证

创建一个 C 程序来检查实际的结构体大小和字段偏移：

```c
#include <stdio.h>
#include <stddef.h>
#include "taichi/taichi_core.h"

int main() {
    printf("sizeof(TiArgumentValue) = %zu\n", sizeof(TiArgumentValue));
    printf("sizeof(TiArgument) = %zu\n", sizeof(TiArgument));
    printf("offset of type:  %zu\n", offsetof(TiArgument, type));
    printf("offset of value: %zu\n", offsetof(TiArgument, value));
    return 0;
}
```

编译并运行：
```bash
gcc -I./taichi/c_api/include check_struct.c -o check_struct
./check_struct
```

### 使用 Go 的 unsafe 包验证

```go
package main

import (
    "fmt"
    "unsafe"
    "go-taichi/taichi/c_api"
)

func main() {
    var arg c_api.TiArgument
    fmt.Printf("sizeof(TiArgument) = %d\n", unsafe.Sizeof(arg))
    fmt.Printf("offset of Type:  %d\n", unsafe.Offsetof(arg.Type))
    fmt.Printf("offset of Value: %d\n", unsafe.Offsetof(arg.Value))
}
```

**两者的输出必须完全一致！**

## 解决方案步骤

1. **使用 C 编译器确定真实大小**
   - 编译一个 C 程序，打印所有关键结构体的 `sizeof()` 和 `offsetof()`
   - 记录每个结构体的大小和字段偏移

2. **在 Go 中验证大小**
   - 使用 `unsafe.Sizeof()` 检查 Go 结构体大小
   - 使用 `unsafe.Offsetof()` 检查字段偏移

3. **添加显式 padding**
   - 如果 offset 不匹配，使用 `_ [N]byte` 添加显式 padding
   - 确保总大小和每个字段的 offset 都与 C 一致

4. **验证数据传递**
   - 创建测试用例，打印原始字节
   - 确认数据在 Go 和 C 之间正确传递

## 常见陷阱

### 1. 假设对齐规则相同
❌ **错误假设**：Go 和 C 的对齐规则总是相同的
✅ **正确做法**：始终用 C 编译器验证实际布局

### 2. 忽略 Union 的特性
❌ **错误做法**：为 union 分配"足够大"的空间（如 512 bytes）
✅ **正确做法**：精确计算 union 的实际大小（最大成员的大小）

### 3. 依赖隐式 padding
❌ **错误做法**：让 Go 编译器自动添加 padding
✅ **正确做法**：显式声明 padding 字段（`_ [N]byte`）

### 4. 未验证数据传递
❌ **错误做法**：假设编译通过就是正确的
✅ **正确做法**：打印原始字节，验证内存布局

## 测试检查清单

创建新的 C 结构体绑定时，请遵循以下清单：

- [ ] 使用 C 编译器检查 `sizeof(struct)`
- [ ] 使用 C 编译器检查所有字段的 `offsetof(struct, field)`
- [ ] 在 Go 中使用 `unsafe.Sizeof()` 验证大小匹配
- [ ] 在 Go 中使用 `unsafe.Offsetof()` 验证所有字段偏移匹配
- [ ] 如果是 union，确认使用最大成员的大小
- [ ] 添加必要的显式 padding（`_ [N]byte`）
- [ ] 创建测试用例，打印原始字节验证数据正确性
- [ ] 使用实际的 C API 调用测试功能

## 相关文件

- `taichi/c_api/types.go` - Go 结构体定义
- `taichi/c_api/include/taichi/taichi_core.h` - C 结构体定义
- `examples/test_struct_alignment.go` - 结构体对齐测试工具

## 参考资料

- [Go unsafe 包文档](https://pkg.go.dev/unsafe)
- [C 结构体对齐规则](https://en.cppreference.com/w/c/language/object#Alignment)
- [Purego FFI 文档](https://github.com/ebitengine/purego)
