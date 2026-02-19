#!/usr/bin/env python3
# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "taichi>=1.7.4",
# ]
# ///
"""
AOT Kernel 生成器 - 为 10_aot_kernel.go 生成 TCM 模块

生成一个简单的加法 kernel：c = a + b

使用方法:
    python 10_aot_kernel.py
    或
    uv run 10_aot_kernel.py
"""

import taichi as ti

@ti.kernel
def add_kernel(
    a: ti.types.ndarray(dtype=ti.f32, ndim=1),
    b: ti.types.ndarray(dtype=ti.f32, ndim=1),
    c: ti.types.ndarray(dtype=ti.f32, ndim=1),
):
    """
    向量加法 kernel: c = a + b

    Args:
        a: 输入数组 A
        b: 输入数组 B
        c: 输出数组 C
    """
    for i in c:
        c[i] = a[i] + b[i]


def main():
    """主函数"""
    print(f"Taichi {ti.__version__} AOT 模块生成器")
    print("目标架构: Vulkan")
    print()

    try:
        # 初始化 Taichi (Vulkan 后端)
        ti.init(arch=ti.vulkan)

        # 创建 AOT 模块
        m = ti.aot.Module(ti.vulkan)

        # 添加 kernel
        print(f"添加 kernel: add_kernel")
        m.add_kernel(add_kernel)

        # 保存为 TCM 文件
        output_file = "aot_module.tcm"
        m.archive(output_file)

        print(f"\n✅ {output_file} 生成成功!")
        print(f"\n现在可以运行: go run 10_aot_kernel.go")

    except Exception as e:
        print(f"\n❌ 生成失败: {e}")
        import traceback
        traceback.print_exc()
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
