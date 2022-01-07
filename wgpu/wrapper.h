#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

void logCallback_cgo(WGPULogLevel level, char const *msg);
void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata);

typedef struct request_adapter_result {
  WGPURequestAdapterStatus status;
  WGPUAdapter adapter;
  char const *message;
} request_adapter_result;

request_adapter_result
request_adapter(WGPURequestAdapterOptions const *options);

typedef struct request_device_result {
  WGPURequestDeviceStatus status;
  WGPUDevice device;
  char const *message;
} request_device_result;

request_device_result request_device(WGPUAdapter adapter,
                                     WGPUDeviceDescriptor const *descriptor);
