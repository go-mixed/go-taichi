// ============================================================
// C 结构体验证代码
// ============================================================
// 将此代码保存为 verify_structs.c
// 编译命令：gcc -I../taichi/c_api/include .\99_struct_alignment.c -o verify_structs
// 运行：./verify_structs
// ============================================================

#include <stdio.h>
#include <stddef.h>
#include "../taichi/c_api/include/taichi/taichi_core.h"

#define PRINT_SECTION(name) printf("\n【%s】\n", name)
#define PRINT_BASIC(name, type) printf("  %-10s size=%2zu, align=%zu\n", name ":", sizeof(type), _Alignof(type))
#define PRINT_STRUCT(type) printf("  Total size:  %3zu bytes\n  Alignment:   %3zu bytes\n", sizeof(type), _Alignof(type))
#define PRINT_FIELD(type, field) printf("  %-20s offset: %3zu\n", #field, offsetof(type, field))

int main() {
    printf("=== C 结构体内存布局完整验证 ===\n");
    printf("\n======================================================================\n");

    // 基础类型
    PRINT_SECTION("基础类型");
    PRINT_BASIC("uintptr", void*);
    PRINT_BASIC("uint32", uint32_t);
    PRINT_BASIC("uint64", uint64_t);
    PRINT_BASIC("int32", int32_t);
    PRINT_BASIC("float32", float);
    PRINT_BASIC("*byte", char*);

    // 图像相关
    PRINT_SECTION("TiImageExtent");
    PRINT_STRUCT(TiImageExtent);
    PRINT_FIELD(TiImageExtent, width);
    PRINT_FIELD(TiImageExtent, height);
    PRINT_FIELD(TiImageExtent, depth);
    PRINT_FIELD(TiImageExtent, array_layer_count);

    PRINT_SECTION("TiImageOffset");
    PRINT_STRUCT(TiImageOffset);
    PRINT_FIELD(TiImageOffset, x);
    PRINT_FIELD(TiImageOffset, y);
    PRINT_FIELD(TiImageOffset, z);
    PRINT_FIELD(TiImageOffset, array_layer_offset);

    PRINT_SECTION("TiImageAllocateInfo");
    PRINT_STRUCT(TiImageAllocateInfo);
    PRINT_FIELD(TiImageAllocateInfo, dimension);
    PRINT_FIELD(TiImageAllocateInfo, extent);
    PRINT_FIELD(TiImageAllocateInfo, mip_level_count);
    PRINT_FIELD(TiImageAllocateInfo, format);
    PRINT_FIELD(TiImageAllocateInfo, export_sharing);
    PRINT_FIELD(TiImageAllocateInfo, usage);

    PRINT_SECTION("TiImageSlice");
    PRINT_STRUCT(TiImageSlice);
    PRINT_FIELD(TiImageSlice, image);
    PRINT_FIELD(TiImageSlice, offset);
    PRINT_FIELD(TiImageSlice, extent);
    PRINT_FIELD(TiImageSlice, mip_level);

    // 采样器
    PRINT_SECTION("TiSamplerCreateInfo");
    PRINT_STRUCT(TiSamplerCreateInfo);
    PRINT_FIELD(TiSamplerCreateInfo, mag_filter);
    PRINT_FIELD(TiSamplerCreateInfo, min_filter);
    PRINT_FIELD(TiSamplerCreateInfo, address_mode);
    PRINT_FIELD(TiSamplerCreateInfo, max_anisotropy);

    // 内存相关
    PRINT_SECTION("TiMemoryAllocateInfo");
    PRINT_STRUCT(TiMemoryAllocateInfo);
    PRINT_FIELD(TiMemoryAllocateInfo, size);
    PRINT_FIELD(TiMemoryAllocateInfo, host_write);
    PRINT_FIELD(TiMemoryAllocateInfo, host_read);
    PRINT_FIELD(TiMemoryAllocateInfo, export_sharing);
    PRINT_FIELD(TiMemoryAllocateInfo, usage);

    PRINT_SECTION("TiMemorySlice");
    PRINT_STRUCT(TiMemorySlice);
    PRINT_FIELD(TiMemorySlice, memory);
    PRINT_FIELD(TiMemorySlice, offset);
    PRINT_FIELD(TiMemorySlice, size);

    // 数组相关
    PRINT_SECTION("TiNdShape");
    PRINT_STRUCT(TiNdShape);
    PRINT_FIELD(TiNdShape, dim_count);
    PRINT_FIELD(TiNdShape, dims);

    PRINT_SECTION("TiNdArray");
    PRINT_STRUCT(TiNdArray);
    PRINT_FIELD(TiNdArray, memory);
    PRINT_FIELD(TiNdArray, shape);
    PRINT_FIELD(TiNdArray, elem_shape);
    PRINT_FIELD(TiNdArray, elem_type);

    // 纹理
    PRINT_SECTION("TiTexture");
    PRINT_STRUCT(TiTexture);
    PRINT_FIELD(TiTexture, image);
    PRINT_FIELD(TiTexture, sampler);
    PRINT_FIELD(TiTexture, dimension);
    PRINT_FIELD(TiTexture, extent);
    PRINT_FIELD(TiTexture, format);

    // 标量
    PRINT_SECTION("TiScalarValue (Union)");
    PRINT_STRUCT(TiScalarValue);
    printf("  注意：Union 大小应等于最大成员 (8 bytes: i64/f64)\n");

    PRINT_SECTION("TiScalar");
    PRINT_STRUCT(TiScalar);
    PRINT_FIELD(TiScalar, type);
    PRINT_FIELD(TiScalar, value);
    if (sizeof(TiScalar) == 12) {
        printf("  ⚠️  警告：大小为 12 bytes，可能需要 padding 到 16 bytes\n");
    }

    // 参数
    PRINT_SECTION("TiArgumentValue (Union)");
    PRINT_STRUCT(TiArgumentValue);
    printf("  注意：Union 大小应等于最大成员 TiNdArray (152 bytes)\n");

    PRINT_SECTION("TiArgument");
    PRINT_STRUCT(TiArgument);
    PRINT_FIELD(TiArgument, type);
    PRINT_FIELD(TiArgument, value);
    if (offsetof(TiArgument, value) == 8) {
        printf("  ✅ Type 和 Value 之间有 4 bytes padding\n");
    }

    PRINT_SECTION("TiNamedArgument");
    PRINT_STRUCT(TiNamedArgument);
    PRINT_FIELD(TiNamedArgument, name);
    PRINT_FIELD(TiNamedArgument, argument);

    printf("\n======================================================================\n");
    printf("\n=== 验证完成 ===\n");
    return 0;
}