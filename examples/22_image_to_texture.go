package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/go-mixed/go-taichi/taichi"
)

// LoadImage loads an image from file (supports PNG, JPEG, etc.)
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image: %w", err)
	}

	return img, nil
}

// LoadImageToTexture loads an image and creates a 3D NdArray texture
func LoadImageToTexture(rt *taichi.Runtime, filePath string) (*taichi.NdArray, error) {
	// Load image from file
	img, err := LoadImage(filePath)
	if err != nil {
		return nil, fmt.Errorf("load image failed: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Image loaded: %dx%d, file: %s\n", width, height, filepath.Base(filePath))

	// Create 3D NdArray: width x height x 4 (RGBA)
	texture, err := taichi.NewNdArray3D(rt, uint32(width), uint32(height), 4, taichi.DataTypeF32)
	if err != nil {
		return nil, fmt.Errorf("create texture failed: %w", err)
	}

	// Write image pixels to texture
	err = texture.MapFloat32(func(data []float32) error {
		// Copy image pixels to texture
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := (y*width + x) * 4

				// Get pixel color - RGBA returns values in [0, 65535]
				r, g, b, a := img.At(x, y).RGBA()
				// Convert to float32 [0, 1]
				data[idx] = float32(r) / 65535.0
				data[idx+1] = float32(g) / 65535.0
				data[idx+2] = float32(b) / 65535.0
				data[idx+3] = float32(a) / 65535.0
			}
		}
		return nil
	})
	if err != nil {
		texture.Release()
		return nil, fmt.Errorf("map texture failed: %w", err)
	}

	return texture, nil
}

func main() {
	fmt.Println("=== Image to Texture Example ===\n")

	// Create runtime (use Vulkan for best compatibility)
	rt, err := taichi.NewRuntime(taichi.ArchVulkan)
	if err != nil {
		panic(fmt.Sprintf("create runtime failed: %v", err))
	}
	defer rt.Release()

	// Test image file path
	imagePath := "examples/22.jpg"

	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("Test image not found at: %s\n", imagePath)
		fmt.Println("Please create a test image or update the path.")
		return
	}

	// Load image and create texture
	texture, err := LoadImageToTexture(rt, imagePath)
	if err != nil {
		panic(fmt.Sprintf("load image to texture failed: %v", err))
	}
	defer texture.Release()

	fmt.Printf("✅ Texture created: shape=%v, elemType=%v\n", texture.Shape(), texture.ElemType())

	bounds := texture.Shape()
	width := int(bounds[0])
	height := int(bounds[1])
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Read back texture data to verify
	err = texture.MapFloat32(func(data []float32) error {
		// Create output image
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := (y*width + x) * 4
				r := uint8(data[idx] * 255)
				g := uint8(data[idx+1] * 255)
				b := uint8(data[idx+2] * 255)
				a := uint8(data[idx+3] * 255)
				outImg.SetRGBA(x, y, color.RGBA{R: g, G: b, B: r, A: a}) // 故意这么写的偏色，测试用。
			}
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("read texture failed: %v", err))
	}

	// Save as PNG
	outputPath := "output_texture.png"
	f, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("create output file failed: %v", err))
	}
	defer f.Close()

	if err := png.Encode(f, outImg); err != nil {
		panic(fmt.Sprintf("encode png failed: %v", err))
	}

	fmt.Printf("✅ Output saved to: %s\n", outputPath)
	fmt.Println("\n=== Example Complete ===")
}
