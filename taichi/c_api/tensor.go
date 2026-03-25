package c_api

import (
	"encoding/binary"
	"math"
)

// tensorToBytes converts a slice to []byte for memory copying
func tensorToBytes(slice any) []byte {
	switch s := slice.(type) {
	case []float32:
		return tensorToBytesFloat32(s)
	case []int32:
		return tensorToBytesInt32(s)
	case []float64:
		return tensorToBytesFloat64(s)
	case []int64:
		return tensorToBytesInt64(s)
	default:
		return nil
	}
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
