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
