package c_api

import (
	"encoding/binary"
	"math"

	"github.com/go-mixed/go-taichi/f16"
)

// tensorToBytes converts a slice to []byte for memory copying
func tensorToBytes(slice any) []byte {
	switch s := slice.(type) {
	case []int8:
		return tensorToBytesInt8(s)
	case []uint8:
		return tensorToBytesUint8(s)
	case []f16.Float16:
		return tensorToBytesFloat16(s)
	case []int16:
		return tensorToBytesInt16(s)
	case []uint16:
		return tensorToBytesUint16(s)
	case []float32:
		return tensorToBytesFloat32(s)
	case []int32:
		return tensorToBytesInt32(s)
	case []uint32:
		return tensorToBytesUint32(s)
	case []float64:
		return tensorToBytesFloat64(s)
	case []int64:
		return tensorToBytesInt64(s)
	case []uint64:
		return tensorToBytesUint64(s)
	default:
		return nil
	}
}

func tensorToBytesInt8(data []int8) []byte {
	result := make([]byte, len(data))
	for i, v := range data {
		result[i] = byte(v)
	}
	return result
}

func tensorToBytesUint8(data []uint8) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

func tensorToBytesFloat16(data []f16.Float16) []byte {
	result := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(result[i*2:i*2+2], v.Bits())
	}
	return result
}

func tensorToBytesInt16(data []int16) []byte {
	result := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(result[i*2:i*2+2], uint16(v))
	}
	return result
}

func tensorToBytesUint16(data []uint16) []byte {
	result := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(result[i*2:i*2+2], v)
	}
	return result
}

func tensorToBytesFloat32(data []float32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(result[i*4:i*4+4], math.Float32bits(v))
	}
	return result
}

func tensorToBytesInt32(data []int32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(result[i*4:i*4+4], uint32(v))
	}
	return result
}

func tensorToBytesUint32(data []uint32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(result[i*4:i*4+4], v)
	}
	return result
}

func tensorToBytesFloat64(data []float64) []byte {
	result := make([]byte, len(data)*8)
	for i, v := range data {
		binary.LittleEndian.PutUint64(result[i*8:i*8+8], math.Float64bits(v))
	}
	return result
}

func tensorToBytesInt64(data []int64) []byte {
	result := make([]byte, len(data)*8)
	for i, v := range data {
		binary.LittleEndian.PutUint64(result[i*8:i*8+8], uint64(v))
	}
	return result
}

func tensorToBytesUint64(data []uint64) []byte {
	result := make([]byte, len(data)*8)
	for i, v := range data {
		binary.LittleEndian.PutUint64(result[i*8:i*8+8], v)
	}
	return result
}
