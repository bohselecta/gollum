//go:build darwin

package metal

/*
#cgo darwin CFLAGS: -fobjc-arc
#cgo darwin LDFLAGS: -framework Metal -framework Foundation
#include "shim.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func Init() bool {
	ok := C.AltMetal_Init()
	return ok == 1
}

func DeviceName() string {
	s := C.AltMetal_DeviceName()
	if s == nil {
		return ""
	}
	defer C.AltMetal_FreeStr(s)
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func DecodeStep(batch, seqlen int) int {
	return int(C.AltMetal_DecodeStep(C.int(batch), C.int(seqlen)))
}

func MatMul(A, B []float32, M, N, K int) ([]float32, bool) {
	if len(A) != M*K || len(B) != K*N {
		return nil, false
	}
	C := make([]float32, M*N)
	result := C.AltMetal_MatMul((*C.float)(&A[0]), (*C.float)(&B[0]), (*C.float)(&C[0]), C.int(M), C.int(N), C.int(K))
	if result == 0 {
		return nil, false
	}
	return C, true
}

func KVCreate(capacity, dim int) int {
	return int(C.AltMetal_KVCreate(C.int(capacity), C.int(dim)))
}

func KVAppend(handle int, K, V []float32, dim int) bool {
	if len(K) != dim || len(V) != dim {
		return false
	}
	result := C.AltMetal_KVAppend(C.int(handle), (*C.float)(&K[0]), (*C.float)(&V[0]), C.int(dim))
	return result != 0
}

func AttnSingle(handle int, Q []float32, dim int) ([]float32, bool) {
	if len(Q) != dim {
		return nil, false
	}
	out := make([]float32, dim)
	result := C.AltMetal_AttnSingle(C.int(handle), (*C.float)(&Q[0]), C.int(dim), (*C.float)(&out[0]))
	return out, result != 0
}

func KVFree(handle int) {
	C.AltMetal_KVFree(C.int(handle))
}
