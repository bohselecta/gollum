//go:build darwin

package metal

/*
#cgo darwin CFLAGS: -fobjc-arc
#cgo darwin LDFLAGS: -framework Metal -framework Foundation -framework MetalPerformanceShaders
#include "shim.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

func Init() bool {
    ok := C.AltMetal_Init()
    return ok == 1
}

func DeviceName() string {
    s := C.AltMetal_DeviceName()
    if s == nil { return "" }
    defer C.AltMetal_FreeStr(s)
    return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func Prefill(batch int) int {
    return int(C.AltMetal_Prefill(C.int(batch)))
}

func DecodeStep(batch, seqlen int) int {
    return int(C.AltMetal_DecodeStep(C.int(batch), C.int(seqlen)))
}

// MatMul performs C = A(MxK) * B(KxN) and returns a length M*N slice (row-major)
func MatMul(A []float32, B []float32, M, N, K int) ([]float32, bool) {
    if M <= 0 || N <= 0 || K <= 0 { return nil, false }
    if len(A) != M*K || len(B) != K*N { return nil, false }
    out := make([]float32, M*N)

    // Allocate C buffers
    aBytes := C.size_t(len(A) * 4)
    bBytes := C.size_t(len(B) * 4)
    cBytes := C.size_t(len(out) * 4)

    aPtr := C.malloc(aBytes)
    bPtr := C.malloc(bBytes)
    cPtr := C.malloc(cBytes)
    if aPtr == nil || bPtr == nil || cPtr == nil {
        if aPtr != nil { C.free(aPtr) }
        if bPtr != nil { C.free(bPtr) }
        if cPtr != nil { C.free(cPtr) }
        return nil, false
    }
    defer C.free(aPtr); defer C.free(bPtr); defer C.free(cPtr)

    C.memcpy(aPtr, unsafe.Pointer(&A[0]), aBytes)
    C.memcpy(bPtr, unsafe.Pointer(&B[0]), bBytes)

    ok := C.AltMetal_MatMul((*C.float)(aPtr), (*C.float)(bPtr), (*C.float)(cPtr), C.int(M), C.int(N), C.int(K))
    if ok != 1 { return nil, false }

    C.memcpy(unsafe.Pointer(&out[0]), cPtr, cBytes)
    return out, true
}
