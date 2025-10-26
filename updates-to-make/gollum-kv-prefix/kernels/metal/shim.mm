//go:build darwin

#import <Metal/Metal.h>
#import <MetalPerformanceShaders/MetalPerformanceShaders.h>
#import <Foundation/Foundation.h>
#import "shim.h"

static id<MTLDevice> gDevice;
static id<MTLCommandQueue> gQueue;
static dispatch_once_t gOnceToken;

int AltMetal_Init() {
    __block int ok = 0;
    dispatch_once(&gOnceToken, ^{
        @autoreleasepool {
            gDevice = MTLCreateSystemDefaultDevice();
            if (gDevice) {
                gQueue = [gDevice newCommandQueue];
            }
            ok = (gDevice != nil && gQueue != nil) ? 1 : 0;
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

int AltMetal_Prefill(int batch) {
    if (!gDevice) { if (!AltMetal_Init()) return 0; }
    return batch;
}

int AltMetal_DecodeStep(int batch, int seqlen) {
    if (!gDevice) { if (!AltMetal_Init()) return 0; }
    return batch > 0 ? batch : 0;
}

int AltMetal_MatMul(const float* A, const float* B, float* C, int M, int N, int K) {
    if (!gDevice) { if (!AltMetal_Init()) return 0; }
    if (M <= 0 || N <= 0 || K <= 0) return 0;

    @autoreleasepool {
        NSUInteger aBytes = (NSUInteger)M * (NSUInteger)K * sizeof(float);
        NSUInteger bBytes = (NSUInteger)K * (NSUInteger)N * sizeof(float);
        NSUInteger cBytes = (NSUInteger)M * (NSUInteger)N * sizeof(float);

        id<MTLBuffer> aBuf = [gDevice newBufferWithBytes:A length:aBytes options:MTLResourceStorageModeShared];
        id<MTLBuffer> bBuf = [gDevice newBufferWithBytes:B length=bBytes options:MTLResourceStorageModeShared];
        id<MTLBuffer> cBuf = [gDevice newBufferWithLength:cBytes options:MTLResourceStorageModeShared];
        if (!aBuf || !bBuf || !cBuf) return 0;

        NSUInteger aRowBytes = (NSUInteger)K * sizeof(float);
        NSUInteger bRowBytes = (NSUInteger)N * sizeof(float);
        NSUInteger cRowBytes = (NSUInteger)N * sizeof(float);

        MPSMatrixDescriptor *aDesc = [MPSMatrixDescriptor matrixDescriptorWithRows:M columns:K rowBytes:aRowBytes dataType:MPSDataTypeFloat32];
        MPSMatrixDescriptor *bDesc = [MPSMatrixDescriptor matrixDescriptorWithRows:K columns:N rowBytes:bRowBytes dataType:MPSDataTypeFloat32];
        MPSMatrixDescriptor *cDesc = [MPSMatrixDescriptor matrixDescriptorWithRows:M columns:N rowBytes:cRowBytes dataType:MPSDataTypeFloat32];

        MPSMatrix *A_m = [[MPSMatrix alloc] initWithBuffer:aBuf descriptor:aDesc];
        MPSMatrix *B_m = [[MPSMatrix alloc] initWithBuffer:bBuf descriptor:bDesc];
        MPSMatrix *C_m = [[MPSMatrix alloc] initWithBuffer:cBuf descriptor:cDesc];

        MPSMatrixMultiplication *mm = [[MPSMatrixMultiplication alloc] initWithDevice:gDevice
                                                                       transposeLeft:false
                                                                      transposeRight:false
                                                                                resultRows:M
                                                                             resultColumns:N
                                                                      interiorColumns:K
                                                                                  alpha:1.0f
                                                                                   beta:0.0f];

        id<MTLCommandBuffer> cb = [gQueue commandBuffer];
        if (!cb) return 0;

        [mm encodeToCommandBuffer:cb leftMatrix:A_m rightMatrix:B_m resultMatrix:C_m];
        [cb commit];
        [cb waitUntilCompleted];

        // Copy result back
        memcpy(C, [cBuf contents], cBytes);

        return 1;
    }
}

void AltMetal_FreeStr(const char* s) {
    if (s) free((void*)s);
}
