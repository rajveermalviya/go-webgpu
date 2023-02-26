package wgpu

/*

#include "./lib/wgpu.h"

void gowebgpu_buffer_map_callback_c(WGPUBufferMapAsyncStatus status, void *userdata) {
  extern void gowebgpu_buffer_map_callback_go(WGPUBufferMapAsyncStatus status, void *userdata);
  gowebgpu_buffer_map_callback_go(status, userdata);
}

void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, char const *message, void *userdata) {
  extern void gowebgpu_request_adapter_callback_go(WGPURequestAdapterStatus status, WGPUAdapter adapter, char const *message, void *userdata);
  gowebgpu_request_adapter_callback_go(status, adapter, message, userdata);
}

void gowebgpu_request_device_callback_c(WGPURequestDeviceStatus status, WGPUDevice device, char const *message, void *userdata) {
  extern void gowebgpu_request_device_callback_go(WGPURequestDeviceStatus status, WGPUDevice device, char const *message, void *userdata);
  gowebgpu_request_device_callback_go(status, device, message, userdata);
}

void gowebgpu_device_uncaptured_error_callback_c(WGPUErrorType type, char const * message, void * userdata) {
  extern void gowebgpu_device_uncaptured_error_callback_go(WGPUErrorType type, char const * message, void * userdata);
  gowebgpu_device_uncaptured_error_callback_go(type, message, userdata);
}

*/
import "C"
