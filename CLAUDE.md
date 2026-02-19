# Go-Taichi 项目上下文

> **项目类型**: Golang封装的Taichi C-API跨平台绑定
> **最后更新**: 2026-02-16
> **Taichi版本**: v1.7.0 (C-API v1007000)
> **Go版本**: 1.25+

---

## 项目概述

Go-Taichi是Taichi C-API的Go语言绑定，使用**purego**实现跨平台支持（Windows/Linux/macOS），**无需CGo编译**。

### 核心特性

- ✅ **跨平台支持** - Windows(DLL) / Linux(SO) / macOS(Dylib)
- ✅ **无需CGo** - 使用purego实现，纯Go编译 (`CGO_ENABLED=0`)
- ✅ **简洁API** - 高级抽象层，自动资源管理
- ✅ **完整覆盖** - 覆盖所有Taichi C-API v1.7.0功能
- ✅ **类型安全** - 完整的类型系统映射
- ✅ **自动管理** - defer Release() 模式，无需手动清理

---

## 项目架构

### 双层设计

```
┌─────────────────────────────────────┐
│   用户代码 (import "go-taichi/taichi")  │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│   高级抽象层 (taichi/)               │
│   • Runtime  - 运行时管理            │
│   • Memory   - 内存基类              │
│   • NdArray  - N维数组               │
│   • Image    - 图像处理              │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│   C-API绑定层 (taichi/c_api/)       │
│   • 纯C函数绑定                      │
│   • 类型定义                         │
│   • 跨平台加载                       │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│   Taichi C动态库                     │
│   taichi_c_api.dll/.so/.dylib       │
└─────────────────────────────────────┘
```

### 目录结构

```
go-taichi/
├── taichi/                    # 公共API包（用户导入这个）
│   ├── c_api/                # C-API绑定（内部包）
│   │   ├── lib/              # 动态库目录
│   │   │   └── taichi_c_api.dll/so/dylib
│   │   ├── include/          # C头文件参考
│   │   │   └── taichi/taichi_core.h
│   │   ├── init.go           # 跨平台初始化入口
│   │   ├── init_windows.go   # Windows加载 (syscall.LoadLibrary)
│   │   ├── init_posix.go     # Linux/macOS加载 (purego.Dlopen)
│   │   ├── types.go          # 类型定义（441行）
│   │   ├── core.go           # 核心API
│   │   ├── memory_ops.go     # 内存操作
│   │   ├── image_ops.go      # 图像操作
│   │   ├── aot.go            # AOT模块
│   │   └── helpers.go        # 辅助函数
│   ├── taichi.go             # 包入口（导出常量、类型别名）
│   ├── runtime.go            # Runtime抽象
│   ├── memory.go             # Memory基类
│   ├── ndarray.go            # NdArray抽象
│   └── image.go              # Image抽象
├── examples/                 # 示例代码
│   ├── basic.go              # 基础示例
│   └── README.md             # 示例说明
├── CLAUDE.md                 # 本文件（完整文档）
└── README.md                 # 项目入口

**设计原则**：用户只导入 `go-taichi/taichi`，无需直接使用 `c_api`
```

---

## API使用方式

### 低级API

```go
archs := c_api.GetAvailableArchs()
runtime := c_api.CreateRuntime(archs[0], 0)
allocInfo := c_api.TiMemoryAllocateInfo{...}
memory := c_api.AllocateMemory(runtime, &allocInfo)
ptr := c_api.MapMemory(runtime, memory)
// ... 需要手动释放 ...
c_api.UnmapMemory(runtime, memory)
c_api.FreeMemory(runtime, memory)
c_api.DestroyRuntime(runtime)
```

### 高级API

```go
runtime, _ := taichi.NewRuntimeAuto()
defer runtime.Release()

arr, _ := taichi.NewNdArray1D(runtime, 1000, taichi.DATA_TYPE_F32)
defer arr.Release()

data, _ := arr.AsSliceFloat32()
// ... 使用数据 ...
arr.Unmap()
```

---

## 类型系统

### C API常量命名规范

**保持C语言风格**（约定俗成，不改为Go风格）：

- `TI_TRUE` / `TI_FALSE` - 布尔值
- `TI_ARCH_VULKAN` / `TI_ARCH_CUDA` - 架构类型
- `TI_FORMAT_RGBA8` / `TI_FORMAT_R32F` - 纹理格式
- `TI_DATA_TYPE_F32` / `TI_DATA_TYPE_I32` - 数据类型
- `TI_NULL_HANDLE` - 空句柄

**原因**：与官方C API文档保持一致，便于查阅和对照

### 句柄类型（7种）

```go
type TiRuntime      uintptr  // 运行时
type TiAotModule    uintptr  // AOT模块
type TiMemory       uintptr  // 内存
type TiImage        uintptr  // 图像
type TiSampler      uintptr  // 采样器
type TiKernel       uintptr  // Kernel
type TiComputeGraph uintptr  // 计算图
```

### 架构类型（8种）

```go
TI_ARCH_VULKAN    // Vulkan (推荐，跨平台)
TI_ARCH_CUDA      // NVIDIA CUDA
TI_ARCH_METAL     // Apple Metal
TI_ARCH_X64       // x64 CPU
TI_ARCH_ARM64     // ARM64 CPU
TI_ARCH_OPENGL    // OpenGL
TI_ARCH_GLES      // OpenGL ES
```

### 数据类型（14种）

```go
TI_DATA_TYPE_F16/F32/F64    // 浮点数
TI_DATA_TYPE_I8/I16/I32/I64 // 有符号整数
TI_DATA_TYPE_U1/U8/U16/U32/U64 // 无符号整数
```

### 纹理格式（44种）

```go
TI_FORMAT_RGBA8/RGBA16F/RGBA32F  // 常用格式
TI_FORMAT_R8/R16/R32F            // 单通道
TI_FORMAT_DEPTH16/DEPTH32F       // 深度格式
// ... 完整列表见 types.go
```

---

## 高级抽象API

### Runtime - 运行时管理

```go
// 自动选择最佳架构
runtime, err := taichi.NewRuntimeAuto()
defer runtime.Release()

// 手动指定架构
runtime, err := taichi.NewRuntime(taichi.ARCH_VULKAN)

// 查询信息
arch := runtime.Arch()
name := runtime.ArchName() // "Vulkan", "CUDA" 等
```

**自动选择优先级**: Vulkan > CUDA > x64 > ARM64 > OpenGL

### Memory - 内存基类

```go
// 分配内存
mem, err := taichi.NewMemory(runtime, 4096) // 4KB
defer mem.Release()

// 映射访问
ptr, err := mem.Map()
// ... 使用 unsafe.Pointer ...
mem.Unmap()

// 查询
size := mem.Size()
mapped := mem.IsMapped()
```

### NdArray - N维数组

```go
// 1D数组
arr, _ := taichi.NewNdArray1D(runtime, 1000, taichi.DATA_TYPE_F32)
defer arr.Release()

// 访问为Go切片
data, _ := arr.AsSliceFloat32()
data[0] = 3.14
arr.Unmap()

// 2D数组
mat, _ := taichi.NewNdArray2D(runtime, 100, 100, taichi.DATA_TYPE_I32)
data2d, _ := mat.AsSliceInt32() // 展平为1D切片

// 3D数组
vol, _ := taichi.NewNdArray3D(runtime, 64, 64, 64, taichi.DATA_TYPE_U8)

// 查询
shape := arr.Shape()        // []uint32{1000}
elemCount := arr.ElemCount() // 1000
```

### Image - 图像处理

```go
// 创建2D图像
img, _ := taichi.NewImage2D(runtime, 512, 512, taichi.FORMAT_RGBA8)
defer img.Release()

// 布局转换
img.TransitionLayout(taichi.IMAGE_LAYOUT_SHADER_READ)

// 图像复制
dst, _ := taichi.NewImage2D(runtime, 512, 512, taichi.FORMAT_RGBA8)
img.CopyTo(dst)

// 查询
w := img.Width()
h := img.Height()
fmt := img.Format()
```

---

## C-API绑定层（内部）

### 核心函数（core.go）

```go
func Init() error
func GetVersion() uint32
func GetAvailableArchs() []TiArch
func CreateRuntime(arch TiArch, deviceIndex uint32) TiRuntime
func DestroyRuntime(runtime TiRuntime)
func GetLastError() (uint64, string)
func SetLastError(errCode TiError, message string)
```

### 内存操作（memory_ops.go）

```go
func AllocateMemory(runtime TiRuntime, allocInfo *TiMemoryAllocateInfo) TiMemory
func FreeMemory(runtime TiRuntime, memory TiMemory)
func MapMemory(runtime TiRuntime, memory TiMemory) unsafe.Pointer
func UnmapMemory(runtime TiRuntime, memory TiMemory)
func CopyMemoryDeviceToDevice(runtime, dstMem, srcMem, size)
func CopyMemoryHostToDevice(runtime, dstMem, srcPtr, size)
func CopyMemoryDeviceToHost(runtime, dstPtr, srcMem, size)
```

### 图像操作（image_ops.go）

```go
func AllocateImage(runtime TiRuntime, allocInfo *TiImageAllocateInfo) TiImage
func FreeImage(runtime TiRuntime, image TiImage)
func CreateSampler(runtime TiRuntime, createInfo *TiSamplerCreateInfo) TiSampler
func DestroySampler(runtime TiRuntime, sampler TiSampler)
func CopyImageDeviceToDevice(runtime, dstSlice, srcSlice)
func TransitionImage(runtime, image, layout)
```

### AOT模块（aot.go）

```go
func LoadAotModule(runtime TiRuntime, modulePath string) TiAotModule
func DestroyAotModule(aotModule TiAotModule)
func GetAotModuleKernel(aotModule TiAotModule, name string) TiKernel
func GetAotModuleComputeGraph(aotModule TiAotModule, name string) TiComputeGraph
func LaunchKernel(runtime, kernel, argsCount, args)
func LaunchComputeGraph(runtime, computeGraph, argsCount, namedArgs)
```

---

## 跨平台加载机制

### Windows (init_windows.go)

```go
//go:build windows

func loadLibrary() (uintptr, error) {
    dll, err := syscall.LoadLibrary(libPath)
    if err != nil {
        return 0, err
    }
    return uintptr(dll), nil
}
```

### Linux/macOS (init_posix.go)

```go
//go:build !windows

func loadLibrary() (uintptr, error) {
    handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
    if err != nil {
        return 0, err
    }
    return uintptr(handle), nil
}
```

### 统一入口 (init.go)

```go
var (
    libHandle uintptr
    initialized bool
)

func Init() error {
    if initialized {
        return nil
    }
    handle, err := loadLibrary() // 平台特定
    if err != nil {
        return err
    }
    libHandle = handle
    registerFunctions() // 注册所有C函数
    initialized = true
    return nil
}
```

---

## 编译运行

### 编译要求

```bash
# 必须禁用CGo
CGO_ENABLED=0 go build -o example.exe example_simple.go

# 或在Windows环境
set CGO_ENABLED=0
go build -o example.exe example_simple.go
```

### 运行要求

1. 将 `taichi_c_api.dll`/`.so`/`.dylib` 放入 `taichi/c_api/lib/` 目录
2. **不需要** `runtime/lib/` 下的任何文件
3. 运行编译后的可执行文件

### 运行输出示例

```
=== Taichi 简洁API示例 ===

✅ 运行时创建成功
📌 使用架构: Vulkan

--- 1D数组示例 ---
✅ 创建1D数组: 形状=[1000], 元素数=1000
✅ 数据总和: 249750.00

--- 2D矩阵示例 ---
✅ 创建2D矩阵: 形状=[4 4]
单位矩阵:
1 0 0 0
0 1 0 0
0 0 1 0
0 0 0 1

--- 图像处理示例 ---
✅ 创建图像: 512x512, 格式=RGBA8
✅ 图像布局转换完成

=== 示例完成 ===
```

---

## 后端兼容性

### 优先级推荐

1. **Vulkan** ⭐⭐⭐⭐⭐ - 跨平台最佳选择
2. **CUDA** ⭐⭐⭐⭐ - NVIDIA GPU专用
3. **CPU (x64/ARM64)** ⭐⭐⭐ - 通用后备方案
4. **OpenGL** ⭐⭐ - 兼容性较差

### 功能支持

| 功能 | Vulkan | CUDA | CPU | OpenGL |
|------|--------|------|-----|--------|
| 内存分配 | ✅ | ✅ | ✅ | ✅ |
| 图像分配 | ✅ | ✅ | ⚠️ | ✅ |
| 自定义采样器 | ❌ | ❌ | ❌ | ❌ |
| AOT Kernel | ✅ | ✅ | ✅ | ⚠️ |
| Compute Graph | ✅ | ✅ | ✅ | ⚠️ |

**注意**：大部分后端不支持自定义采样器，使用 `TI_NULL_HANDLE` 即可

---

## 常见问题

### ❓ 需要复制 `runtime/lib/` 下的文件吗？

**不需要！** 只需要一个 `taichi_c_api.dll` (或 `.so`/`.dylib`)。

### ❓ 为什么常量命名不符合Go规范？

因为是C API的约定俗成命名（如 `TI_TRUE`），保持原样便于与官方文档对照。

### ❓ 为什么使用purego？

- 无需C编译器
- 跨平台编译简单
- 编译速度快
- 部署方便

### ❓ 看到采样器警告怎么办？

`ti_create_sampler: not supported` 是正常的，大部分后端不支持自定义采样器。使用 `TI_NULL_HANDLE` 作为默认采样器。

### ❓ 如何获取Taichi动态库？

从Python环境获取：
```bash
pip install taichi==1.7.0
# Windows: Python\Lib\site-packages\_lib\c_api\bin\taichi_c_api.dll
# Linux: Python/lib/pythonX.X/site-packages/_lib/c_api/bin/libtaichi_c_api.so
# macOS: Python/lib/pythonX.X/site-packages/_lib/c_api/bin/libtaichi_c_api.dylib
```

---

## 开发状态

| 功能模块 | 状态 | 说明 |
|---------|------|------|
| C-API绑定 | ✅ 完成 | 所有核心C函数已绑定 |
| 高级抽象 | ✅ 完成 | Runtime、Memory、NdArray、Image、Sampler |
| 内存管理 | ✅ 完成 | 自动释放，defer模式 |
| 图像处理 | ✅ 完成 | 创建、转换、复制 |
| AOT支持 | ✅ 完成 | Kernel + ComputeGraph，Builder模式 |
| 内存导入 | ✅ 完成 | CPU/CUDA内存导入，流管理 |
| 结构体对齐 | ✅ 完成 | 修复Go-C内存布局问题 |
| 文档 | ✅ 完成 | 完整文档 + 示例 + 测试工具 |

---

## 贡献指南

### 添加新功能

1. 在 `taichi/c_api/` 中添加C函数绑定
2. 在 `taichi/` 中添加高级抽象（如需要）
3. 更新 `taichi/taichi.go` 导出必要的常量/类型
4. 在 `example_simple.go` 中添加示例
5. 更新本文档

### 代码规范

- C API绑定：直接映射，不添加额外逻辑
- 高级抽象：提供便捷接口，自动资源管理
- 错误处理：返回 `error` 类型，包含错误码和消息
- 命名：C常量保持原样，Go类型遵循Go规范

---

## 许可证

遵循Taichi项目的开源许可证。

---

## 致谢

- [Taichi](https://github.com/taichi-dev/taichi) - 高性能并行编程语言
- [purego](https://github.com/ebitengine/purego) - 纯Go FFI库

---

**最后编译验证**: 2026-02-16 ✅ 通过
**测试架构**: Vulkan (Windows)
**Go版本**: 1.25+
