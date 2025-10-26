//go:build darwin

#import <Metal/Metal.h>
#import <Foundation/Foundation.h>
#import "shim.h"

static id<MTLDevice> gDevice;
static dispatch_once_t gOnceToken;

int AltMetal_Init() {
    __block int ok = 0;
    dispatch_once(&gOnceToken, ^{
        @autoreleasepool {
            gDevice = MTLCreateSystemDefaultDevice();
            ok = (gDevice != nil) ? 1 : 0;
        }
    });
    return ok;
}

const char* AltMetal_DeviceName() {
    @autoreleasepool {
        if (!gDevice) { AltMetal_Init(); }
        NSString *name = [gDevice name];
        if (!name) { name = @"(no device)"; }
        const char *utf8 = [name UTF8String];
        size_t len = strlen(utf8);
        char *cpy = (char*)malloc(len+1);
        memcpy(cpy, utf8, len+1);
        return cpy;
    }
}

int AltMetal_DecodeStep(int batch, int seqlen) {
    // Placeholder: this is where you'll encode a command buffer that runs a fused decode step.
    // For now, we just pretend something happened and return the number of tokens produced.
    // You could insert a trivial no-op compute pipeline here later.
    if (!gDevice) { if (!AltMetal_Init()) return 0; }
    // Simulate: tokens out == batch (one token per sequence)
    return batch > 0 ? batch : 0;
}

int AltMetal_MatMul(const float* A, const float* B, float* C, int M, int N, int K) {
    // Simple CPU fallback matrix multiplication
    // In a real implementation, this would use Metal compute shaders
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            float sum = 0.0f;
            for (int k = 0; k < K; k++) {
                sum += A[i * K + k] * B[k * N + j];
            }
            C[i * N + j] = sum;
        }
    }
    return 1; // success
}

void AltMetal_FreeStr(const char* s) {
    if (s) free((void*)s);
}
