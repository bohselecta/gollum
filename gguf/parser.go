package gguf

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
)

// GGUF File Format (simplified)
// Based on: https://github.com/ggerganov/llama.cpp/blob/master/gguf-spec.md

const (
	GGUFMagic   = 0x46554747 // "GGUF"
	GGUFCows    = 0x666F     // "fo"
	GGUFVersion = 3
)

type Tensor struct {
	Name string
	Dims []uint32
	Type uint32
	Data []float32
	Size uint64
}

type GGUFHeader struct {
	Magic       uint32
	Version     uint32
	TensorCount uint64
	KVCount     uint64
}

type GGUF struct {
	Header   *GGUFHeader
	Tensors  map[string]*Tensor
	Metadata map[string]interface{}
}

// Parse reads a GGUF file and returns the parsed structure
func Parse(r io.ReadSeeker) (*GGUF, error) {
	g := &GGUF{
		Tensors:  make(map[string]*Tensor),
		Metadata: make(map[string]interface{}),
	}

	// Read magic bytes
	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return nil, fmt.Errorf("failed to read magic: %w", err)
	}
	if magic != GGUFMagic {
		return nil, fmt.Errorf("invalid GGUF magic: 0x%08x", magic)
	}

	// Read version
	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}

	// Read tensor and KV counts
	var tensorCount, kvCount uint64
	if err := binary.Read(r, binary.LittleEndian, &tensorCount); err != nil {
		return nil, fmt.Errorf("failed to read tensor count: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &kvCount); err != nil {
		return nil, fmt.Errorf("failed to read KV count: %w", err)
	}

	g.Header = &GGUFHeader{
		Magic:       magic,
		Version:     version,
		TensorCount: tensorCount,
		KVCount:     kvCount,
	}

	// For now, skip metadata parsing and jump to tensors
	// TODO: Parse KV metadata

	// Parse tensors
	for i := uint64(0); i < tensorCount; i++ {
		tensor, err := parseTensor(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tensor %d: %w", i, err)
		}
		g.Tensors[tensor.Name] = tensor
	}

	return g, nil
}

func parseTensor(r io.Reader) (*Tensor, error) {
	t := &Tensor{}

	// Read name length
	var nameLen uint64
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return nil, err
	}

	// Read name
	nameBytes := make([]byte, nameLen)
	if _, err := io.ReadFull(r, nameBytes); err != nil {
		return nil, err
	}
	t.Name = string(nameBytes)

	// Read number of dimensions
	var nDims uint32
	if err := binary.Read(r, binary.LittleEndian, &nDims); err != nil {
		return nil, err
	}

	// Read dimensions
	t.Dims = make([]uint32, nDims)
	if err := binary.Read(r, binary.LittleEndian, &t.Dims); err != nil {
		return nil, err
	}

	// Read type
	if err := binary.Read(r, binary.LittleEndian, &t.Type); err != nil {
		return nil, err
	}

	// Read offset (skip for now, we'll read data directly)
	var offset uint64
	if err := binary.Read(r, binary.LittleEndian, &offset); err != nil {
		return nil, err
	}

	// Calculate size
	t.Size = 1
	for _, dim := range t.Dims {
		t.Size *= uint64(dim)
	}

	// Read data based on type
	if err := readTensorData(r, t); err != nil {
		return nil, err
	}

	return t, nil
}

func readTensorData(r io.Reader, t *Tensor) error {
	// Read data based on type
	// Note: Simplified implementation - real GGUF has tensors at specific offsets
	// For now, read directly (actual files need seek to offset)
	
	switch t.Type {
	case 0: // F32
		t.Data = make([]float32, t.Size)
		for i := uint64(0); i < t.Size; i++ {
			if err := binary.Read(r, binary.LittleEndian, &t.Data[i]); err != nil {
				return err
			}
		}
	case 1: // F16 - convert to F32
		t.Data = make([]float32, t.Size)
		for i := uint64(0); i < t.Size; i++ {
			var f16 uint16
			if err := binary.Read(r, binary.LittleEndian, &f16); err != nil {
				return err
			}
			t.Data[i] = float16ToFloat32(f16)
		}
	case 2, 3: // Q4_0, Q4_1 - quantized, dequantize to F32
		return readQuantizedTensor(r, t)
	case 8: // Q8_0
		return readQuantizedTensor(r, t)
	default:
		log.Printf("Warning: unknown tensor type %d, skipping", t.Type)
		// Skip the tensor data
		return nil
	}
	
	return nil
}

func readQuantizedTensor(r io.Reader, t *Tensor) error {
	// For quantized formats, we dequantize to F32
	// This is simplified - real implementation needs proper dequantization
	t.Data = make([]float32, t.Size)
	// Q4_0 uses 32 values per block (18 bytes: 2 byte scale + 16 bytes data)
	blocks := t.Size / 32
	
	for i := uint64(0); i < blocks; i++ {
		// Read scale
		var scale float32
		if err := binary.Read(r, binary.LittleEndian, &scale); err != nil {
			return err
		}
		// Read quantized data (8 bytes for Q4_0)
		buf := make([]byte, 8)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		// Dequantize
		for j := 0; j < 32 && (i*32+uint64(j)) < t.Size; j++ {
			nibble := (buf[j/2] >> ((j % 2) * 4)) & 0xf
			q := int8(nibble)
			t.Data[i*32+uint64(j)] = float32(q-8) * scale
		}
	}
	return nil
}

// float16ToFloat32 converts IEEE 754 half-precision to float32
func float16ToFloat32(f16 uint16) float32 {
	sign := float32((f16 >> 15) & 0x1)
	exp := int((f16 >> 10) & 0x1f)
	mantissa := float32(f16 & 0x3ff)

	if exp == 0 {
		return sign * mantissa * float32(math.Pow(2, -24))
	} else if exp == 31 {
		return (sign*2 - 1) * float32(math.Inf(1))
	}

	expValue := float32(math.Pow(2, float64(exp-15)))
	mantissaValue := 1.0 + mantissa/1024.0

	return sign * expValue * mantissaValue
}
