// ============================================================
// Verify memory layout of C structures
// ============================================================
// Compile: gcc -I../taichi/c_api/include .\99_struct_alignment.c -o verify_structs
// Run: ./verify_structs
// ============================================================

#include <stdio.h>
#include <stddef.h>
#include "../taichi/c_api/include/taichi/taichi_core.h"

#define PRINT_SECTION(name) printf("\n【%s】\n", name)
#define PRINT_BASIC(name, type) printf("  %-10s size=%2zu, align=%zu\n", name ":", sizeof(type), _Alignof(type))
#define PRINT_STRUCT(type) printf("  Total size:  %3zu bytes\n  Alignment:   %3zu bytes\n", sizeof(type), _Alignof(type))
#define PRINT_FIELD(type, field) printf("  %-20s offset: %3zu\n", #field, offsetof(type, field))

int main() {
    printf("=== Verify memory layout of C structures ===\n");
    printf("\n======================================================================\n");

    // Basic Types
    PRINT_SECTION("Basic Types");
    PRINT_BASIC("uintptr", void*);
    PRINT_BASIC("uint32", uint32_t);
    PRINT_BASIC("uint64", uint64_t);
    PRINT_BASIC("int32", int32_t);
    PRINT_BASIC("float32", float);
    PRINT_BASIC("*byte", char*);

    // Image related
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

    // Sampler
    PRINT_SECTION("TiSamplerCreateInfo");
    PRINT_STRUCT(TiSamplerCreateInfo);
    PRINT_FIELD(TiSamplerCreateInfo, mag_filter);
    PRINT_FIELD(TiSamplerCreateInfo, min_filter);
    PRINT_FIELD(TiSamplerCreateInfo, address_mode);
    PRINT_FIELD(TiSamplerCreateInfo, max_anisotropy);

    // Memory related
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

    // Array related
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

    // Texture
    PRINT_SECTION("TiTexture");
    PRINT_STRUCT(TiTexture);
    PRINT_FIELD(TiTexture, image);
    PRINT_FIELD(TiTexture, sampler);
    PRINT_FIELD(TiTexture, dimension);
    PRINT_FIELD(TiTexture, extent);
    PRINT_FIELD(TiTexture, format);

    // Scalar
    PRINT_SECTION("TiScalarValue (Union)");
    PRINT_STRUCT(TiScalarValue);
    printf("  Note: Union size should equal largest member (8 bytes: i64/f64)\n");

    PRINT_SECTION("TiScalar");
    PRINT_STRUCT(TiScalar);
    PRINT_FIELD(TiScalar, type);
    PRINT_FIELD(TiScalar, value);
    if (sizeof(TiScalar) == 12) {
        printf("  ⚠️  Warning: Size is 12 bytes, may need padding to 16 bytes\n");
    }

    // Arguments
    PRINT_SECTION("TiArgumentValue (Union)");
    PRINT_STRUCT(TiArgumentValue);
    printf("  Note: Union size should equal largest member TiNdArray (152 bytes)\n");

    PRINT_SECTION("TiArgument");
    PRINT_STRUCT(TiArgument);
    PRINT_FIELD(TiArgument, type);
    PRINT_FIELD(TiArgument, value);
    if (offsetof(TiArgument, value) == 8) {
        printf("  ✅ 4 bytes padding between Type and Value\n");
    }

    PRINT_SECTION("TiNamedArgument");
    PRINT_STRUCT(TiNamedArgument);
    PRINT_FIELD(TiNamedArgument, name);
    PRINT_FIELD(TiNamedArgument, argument);

    printf("\n======================================================================\n");
    printf("\n=== Verification Complete ===\n");
    return 0;
}