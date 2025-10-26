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

// Simple KV cache implementation (CPU fallback for demo)
static NSMutableDictionary* gKVHandles = nil;

int AltMetal_KVCreate(int capacity, int dim) {
    if (!gKVHandles) {
        gKVHandles = [[NSMutableDictionary alloc] init];
    }
    static int nextHandle = 1;
    int handle = nextHandle++;
    NSMutableArray* kv = [[NSMutableArray alloc] init];
    [gKVHandles setObject:kv forKey:@(handle)];
    return handle;
}

int AltMetal_KVAppend(int handle, const float* K, const float* V, int dim) {
    if (!gKVHandles) return 0;
    NSMutableArray* kv = [gKVHandles objectForKey:@(handle)];
    if (!kv) return 0;
    
    // Store K and V vectors
    NSMutableData* kData = [NSMutableData dataWithBytes:K length:dim * sizeof(float)];
    NSMutableData* vData = [NSMutableData dataWithBytes:V length:dim * sizeof(float)];
    [kv addObject:@[kData, vData]];
    return 1;
}

int AltMetal_AttnSingle(int handle, const float* Q, int dim, float* out) {
    if (!gKVHandles) return 0;
    NSMutableArray* kv = [gKVHandles objectForKey:@(handle)];
    if (!kv || kv.count == 0) {
        // No KV cache, just copy Q to out
        memcpy(out, Q, dim * sizeof(float));
        return 1;
    }
    
    // Simple attention: average all V vectors weighted by similarity to Q
    float* result = (float*)calloc(dim, sizeof(float));
    float totalWeight = 0.0f;
    
    for (NSArray* kvPair in kv) {
        NSData* kData = kvPair[0];
        NSData* vData = kvPair[1];
        const float* K = (const float*)kData.bytes;
        const float* V = (const float*)vData.bytes;
        
        // Compute similarity (dot product)
        float sim = 0.0f;
        for (int i = 0; i < dim; i++) {
            sim += Q[i] * K[i];
        }
        
        // Add weighted V to result
        for (int i = 0; i < dim; i++) {
            result[i] += sim * V[i];
        }
        totalWeight += fabsf(sim);
    }
    
    // Normalize
    if (totalWeight > 0) {
        for (int i = 0; i < dim; i++) {
            result[i] /= totalWeight;
        }
    }
    
    memcpy(out, result, dim * sizeof(float));
    free(result);
    return 1;
}

void AltMetal_KVFree(int handle) {
    if (!gKVHandles) return;
    [gKVHandles removeObjectForKey:@(handle)];
}

void AltMetal_FreeStr(const char* s) {
    if (s) free((void*)s);
}