# Go-Taichi 示例

本目录包含 Go-Taichi 的使用示例，**每个示例只演示一个功能**，从文件名即可看出功能。

---

## 📚 示例索引

### 基础示例 (01-09)

| 文件 | 功能 | 难度 |
|------|------|------|
| `01_runtime.go` | 创建和管理运行时 | ⭐ |
| `02_ndarray_1d.go` | 1D 数组操作 | ⭐ |
| `03_ndarray_2d.go` | 2D 矩阵操作 | ⭐ |
| `04_image.go` | 图像基础操作 | ⭐ |

### AOT 示例 (10-19)

| 文件 | 功能 | 难度 | 前置要求 |
|------|------|------|----------|
| `10_aot_kernel.go` | AOT Kernel 基础执行 | ⭐⭐ | AOT 模块 |
| `11_aot_async.go` | AOT Kernel 异步执行 | ⭐⭐ | AOT 模块 |
| `12_aot_batch.go` | AOT Kernel 批量执行 | ⭐⭐ | AOT 模块 |
| `13_compute_graph.go` | Compute Graph 执行 | ⭐⭐⭐ | Compute Graph 模块 |

### 高级示例 (20-29)

| 文件 | 功能 | 难度 |
|------|------|------|
| `20_memory_cpu.go` | CPU 内存导入 | ⭐⭐⭐ |
| `21_memory_cuda.go` | CUDA 内存导入和流管理 | ⭐⭐⭐⭐ |

### 工具 (90+)

| 文件 | 功能 | 说明 |
|------|------|------|
| `99_struct_alignment.go` | 结构体对齐测试 | 诊断工具，验证 Go-C FFI 正确性 |

---

## 🚀 快速开始

### 1. 基础示例（无需前置步骤）

```bash
# 运行时管理
CGO_ENABLED=0 go run 01_runtime.go

# 1D 数组
CGO_ENABLED=0 go run 02_ndarray_1d.go

# 2D 矩阵
CGO_ENABLED=0 go run 03_ndarray_2d.go

# 图像操作
CGO_ENABLED=0 go run 04_image.go
```

### 2. AOT 示例（需要生成 AOT 模块）

**前置步骤：**

```bash
# 安装 Taichi
pip install taichi==1.7.0

# 生成 AOT 模块
python generate_aot.py
```

**运行示例：**

```bash
# 基础 Kernel 执行
CGO_ENABLED=0 go run 10_aot_kernel.go

# 异步执行
CGO_ENABLED=0 go run 11_aot_async.go

# 批量执行
CGO_ENABLED=0 go run 12_aot_batch.go

# Compute Graph（需要特殊的 AOT 模块）
python generate_compute_graph.py
CGO_ENABLED=0 go run 13_compute_graph.go
```

### 3. 高级示例

```bash
# CPU 内存导入
CGO_ENABLED=0 go run 20_memory_cpu.go

# CUDA 内存导入（概念演示）
CGO_ENABLED=0 go run 21_memory_cuda.go
```

---

## 📖 学习路径

### 初学者路径

1. `01_runtime.go` - 了解如何创建运行时
2. `02_ndarray_1d.go` - 学习 1D 数组操作
3. `03_ndarray_2d.go` - 学习 2D 矩阵操作
4. `04_image.go` - 学习图像操作

### 进阶路径

5. `10_aot_kernel.go` - 学习 AOT Kernel 基础
6. `11_aot_async.go` - 学习异步执行
7. `12_aot_batch.go` - 学习批量执行优化

### 高级路径

8. `13_compute_graph.go` - 学习复杂计算图
9. `20_memory_cpu.go` - 学习内存导入
10. `21_memory_cuda.go` - 学习 CUDA 集成

---

## 🔧 AOT 模块生成

### 基础 AOT 模块

创建 `generate_aot.py`：

```python
import taichi as ti

ti.init(arch=ti.vulkan)

@ti.kernel
def add_kernel(a: ti.types.ndarray(), b: ti.types.ndarray(), c: ti.types.ndarray()):
    for i in a:
        c[i] = a[i] + b[i]

# 导出 AOT
m = ti.aot.Module(arch=ti.vulkan)
m.add_kernel(add_kernel, template_args={
    'a': ti.types.ndarray(),
    'b': ti.types.ndarray(),
    'c': ti.types.ndarray()
})
m.archive('./aot_module.tcm')
print("✅ AOT 模块已生成到 ./aot_module.tcm")
```

运行：
```bash
python generate_aot.py
```

### Compute Graph 模块

创建 `generate_compute_graph.py`：

```python
import taichi as ti

ti.init(arch=ti.vulkan)

# 定义多个 kernel
@ti.kernel
def kernel1(a: ti.types.ndarray(), b: ti.types.ndarray()):
    for i in a:
        b[i] = a[i] * 2

@ti.kernel
def kernel2(b: ti.types.ndarray(), c: ti.types.ndarray(), scale: ti.f32):
    for i in b:
        c[i] = b[i] * scale

# 创建 Compute Graph
graph_builder = ti.graph.GraphBuilder()

# 定义输入输出
input_a = graph_builder.create_ndarray_arg(ti.f32, 1)
input_b = graph_builder.create_ndarray_arg(ti.f32, 1)
output_c = graph_builder.create_ndarray_arg(ti.f32, 1)
scale_factor = graph_builder.create_scalar_arg(ti.f32)

# 添加 kernel 到图
graph_builder.dispatch(kernel1, input_a, input_b)
graph_builder.dispatch(kernel2, input_b, output_c, scale_factor)

# 构建图
graph = graph_builder.compile()

# 导出
m = ti.aot.Module(arch=ti.vulkan)
m.add_graph('my_compute_graph', graph)
m.save('./aot_module')
print("✅ Compute Graph 模块已生成")
```

---

## 💡 命名规则

- **数字前缀**：表示难度和学习顺序
  - `01-09`: 基础功能
  - `10-19`: AOT 功能
  - `20-29`: 高级功能
  - `30-39`: 渲染功能（预留）
  - `90+`: 工具和测试

- **功能名称**：清晰描述单一功能
  - `runtime` - 运行时管理
  - `ndarray_1d` - 1D 数组
  - `aot_kernel` - AOT Kernel
  - `memory_cpu` - CPU 内存导入

---

## ❓ 常见问题

### 编译失败？

确保设置 `CGO_ENABLED=0`：
```bash
CGO_ENABLED=0 go build xxx.go
```

### AOT 模块加载失败？

1. 确认 `./aot_module.tcm` 存在
2. 确认包含 `metadata.json` 文件
3. 确认 Python 和 Go 使用相同架构（如 Vulkan）
4. 确认 Taichi 版本为 1.7.0

### Kernel 找不到？

确保 Python 导出的 kernel 名称与 Go 中使用的名称一致：

```python
# Python
@ti.kernel
def add_kernel(...):  # 名称: add_kernel
    pass
```

```go
// Go
kernel, _ := module.GetKernel("add_kernel")  // 使用相同名称
```

### 参数类型不匹配？

确保参数顺序和类型完全匹配：

```python
# Python
@ti.kernel
def kernel(arr: ti.types.ndarray(), value: ti.f32):
    pass
```

```go
// Go - 顺序和类型必须匹配
kernel.Launch().
    ArgNdArray(arr).     // 第1个: ndarray
    ArgFloat32(value).   // 第2个: f32
    Run()
```

---

## 📊 示例对比

| 类别 | 示例数量 | 总难度 | 前置要求 |
|------|---------|--------|----------|
| 基础 | 4 | ⭐ | 无 |
| AOT | 4 | ⭐⭐ | AOT 模块 |
| 高级 | 2 | ⭐⭐⭐ | 无/CUDA |
| 工具 | 1 | ⭐⭐⭐⭐ | 无 |

---

## 📚 参考资料

- [Taichi AOT 文档](https://docs.taichi-lang.org/docs/aot)
- [Go-Taichi 完整文档](../CLAUDE.md)
- [Taichi Python 教程](https://docs.taichi-lang.org/)
- [结构体对齐说明](../taichi/c_api/STRUCT_ALIGNMENT.md)

---

## 🎯 设计原则

1. **一个示例一个功能** - 每个文件只演示一个核心功能
2. **从文件名看功能** - 无需打开文件即可知道内容
3. **数字前缀排序** - 按难度和学习顺序组织
4. **独立可运行** - 每个示例都是完整的 main 程序
5. **清晰的注释** - 代码开头说明功能和要点

---

**提示**：如果遇到问题，请查看 [../CLAUDE.md](../CLAUDE.md) 的常见问题部分。
