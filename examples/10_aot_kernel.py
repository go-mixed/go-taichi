#!/usr/bin/env python3
# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "taichi>=1.7.4",
# ]
# ///
"""
AOT Kernel Generator - Generate TCM module for 10_aot_kernel.go

Generates a simple addition kernel: c = a + b

Usage:
    python 10_aot_kernel.py
    or
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
    Vector addition kernel: c = a + b

    Args:
        a: Input array A
        b: Input array B
        c: Output array C
    """
    for i in c:
        c[i] = a[i] + b[i]


def main():
    """Main function"""
    print(f"Taichi {ti.__version__} AOT Module Generator")
    print("Target architecture: Vulkan")
    print()

    try:
        # Initialize Taichi (Vulkan backend)
        ti.init(arch=ti.vulkan)

        # Create AOT module
        m = ti.aot.Module(ti.vulkan)

        # Add kernel
        print(f"Adding kernel: add_kernel")
        m.add_kernel(add_kernel)

        # Save as TCM file
        output_file = "aot_module.tcm"
        m.archive(output_file)

        print(f"\n✅ {output_file} generated successfully!")
        print(f"\nNow you can run: go run 10_aot_kernel.go")

    except Exception as e:
        print(f"\n❌ Generation failed: {e}")
        import traceback
        traceback.print_exc()
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
