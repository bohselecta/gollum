#ifndef ALTISERVE_METAL_SHIM_H
#define ALTISERVE_METAL_SHIM_H

#ifdef __cplusplus
extern "C" {
#endif

int AltMetal_Init();
const char* AltMetal_DeviceName();
int AltMetal_Prefill(int batch);
int AltMetal_DecodeStep(int batch, int seqlen);

// Perform C = A(MxK) * B(KxN) in float32. Returns 1 on success, 0 on failure.
int AltMetal_MatMul(const float* A, const float* B, float* C, int M, int N, int K);

void AltMetal_FreeStr(const char* s);

#ifdef __cplusplus
}
#endif

#endif // ALTISERVE_METAL_SHIM_H
