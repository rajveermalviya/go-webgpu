//go:build !windows

package wgpu

/*

#include "./lib/wgpu.h"

void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata) {
  extern void bufferMapCallback(WGPUBufferMapAsyncStatus status,
                                void *userdata);
  bufferMapCallback(status, userdata);
}

void requestAdapterCallback_cgo(WGPURequestAdapterStatus status,
                                WGPUAdapter adapter, char const *message,
                                void *userdata) {
  extern void requestAdapterCallback(WGPURequestAdapterStatus status,
                                     WGPUAdapter adapter, char const *message,
                                     void *userdata);
  requestAdapterCallback(status, adapter, message, userdata);
}

void requestDeviceCallback_cgo(WGPURequestDeviceStatus status,
                               WGPUDevice device, char const *message,
                               void *userdata) {
  extern void requestDeviceCallback(WGPURequestDeviceStatus status,
                                    WGPUDevice device, char const *message,
                                    void *userdata);
  requestDeviceCallback(status, device, message, userdata);
}

void deviceUncapturedErrorCallback_cgo(WGPUErrorType type, char const * message,
                               void * userdata) {
  extern void deviceUncapturedErrorCallback(WGPUErrorType type,
                               char const * message, void * userdata);

  deviceUncapturedErrorCallback(type, message, userdata);
}

*/
import "C"
