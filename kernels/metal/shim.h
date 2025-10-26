#ifndef ALTISERVE_METAL_SHIM_H
#define ALTISERVE_METAL_SHIM_H

#ifdef __cplusplus
extern "C" {
#endif

// Initialize Metal and return 1 on success, 0 on failure.
int AltMetal_Init();
// Returns a malloc'd C string with the device name (caller must free via AltMetal_FreeStr).
const char* AltMetal_DeviceName();
// Example placeholder for a fused op entrypoint.
int AltMetal_DecodeStep(int batch, int seqlen);
// Matrix multiplication: A[M*K] * B[K*N] = C[M*N]
int AltMetal_MatMul(const float* A, const float* B, float* C, int M, int N, int K);
// KV cache operations
int AltMetal_KVCreate(int capacity, int dim);
int AltMetal_KVAppend(int handle, const float* K, const float* V, int dim);
int AltMetal_AttnSingle(int handle, const float* Q, int dim, float* out);
void AltMetal_KVFree(int handle);
// Free a C string returned by this shim.
void AltMetal_FreeStr(const char* s);

#ifdef __cplusplus
}
#endif

#endif // ALTISERVE_METAL_SHIM_H
