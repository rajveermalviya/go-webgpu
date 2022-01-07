#include "wrapper.h"

void logCallback_cgo(WGPULogLevel level, char const *msg) {
  extern void logCallback(WGPULogLevel level, char const *msg);
  logCallback(level, msg);
}

void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata) {
  extern void bufferMapCallback(WGPUBufferMapAsyncStatus status,
                                void *userdata);
  bufferMapCallback(status, userdata);
}

void request_adapter_cb(WGPURequestAdapterStatus status, WGPUAdapter adapter,
                        char const *message, void *userdata) {
  request_adapter_result *res;
  res = userdata;

  res->status = status;
  res->adapter = adapter;
  res->message = message;
}

request_adapter_result
request_adapter(WGPURequestAdapterOptions const *options) {
  request_adapter_result res;

  wgpuInstanceRequestAdapter(NULL, options, request_adapter_cb, &res);

  return res;
}

void request_device_cb(WGPURequestDeviceStatus status, WGPUDevice device,
                       char const *message, void *userdata) {
  request_device_result *res;
  res = userdata;

  res->status = status;
  res->device = device;
  res->message = message;
}

request_device_result request_device(WGPUAdapter adapter,
                                     WGPUDeviceDescriptor const *descriptor) {
  request_device_result res;

  wgpuAdapterRequestDevice(adapter, descriptor, request_device_cb, &res);

  return res;
}
