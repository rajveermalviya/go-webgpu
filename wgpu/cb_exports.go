package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export logCallback
func logCallback(level C.WGPULogLevel, msg *C.char) {
	logCb(LogLevel(level), C.GoString(msg))
}

//export bufferMapCallback
func bufferMapCallback(status C.WGPUBufferMapAsyncStatus, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(BufferMapCallback)
	if ok {
		cb(BufferMapAsyncStatus(status))
	}
}

//export requestAdapterCallback
func requestAdapterCallback(status C.WGPURequestAdapterStatus, adapter C.WGPUAdapter, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(requestAdapterCB)
	if ok {
		cb(RequestAdapterStatus(status), &Adapter{adapter}, C.GoString(message))
	}
}

//export requestDeviceCallback
func requestDeviceCallback(status C.WGPURequestDeviceStatus, device C.WGPUDevice, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(requestDeviceCB)
	if ok {
		cb(RequestDeviceStatus(status), &Device{device}, C.GoString(message))
	}
}
