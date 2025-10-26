//go:build darwin

package main

import (
	"fmt"

	"github.com/haydenlabs/gollum/kernels/metal"
)

func selfTestMetal() {
	// Only prints if Metal is present (darwin build)
	defer func() { recover() }()
	name := metalDeviceName()
	if name == "" {
		return
	}
	fmt.Printf("[GoLLuM] Metal device: %s\n", name)
	// 2x3 * 3x2 => 2x2
	A := []float32{1, 2, 3, 4, 5, 6}
	B := []float32{7, 8, 9, 10, 11, 12}
	C, ok := metalMatMul(A, B, 2, 2, 3)
	if !ok {
		fmt.Println("[GoLLuM] Metal matmul failed")
		return
	}
	fmt.Printf("[GoLLuM] Metal matmul ok: %.1f %.1f | %.1f %.1f\n", C[0], C[1], C[2], C[3])
}

func metalDeviceName() string {
	return metal.DeviceName()
}

func metalMatMul(A, B []float32, M, N, K int) ([]float32, bool) {
	return metal.MatMul(A, B, M, N, K)
}
