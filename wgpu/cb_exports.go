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

type requestAdapterCB func(status RequestAdapterStatus, adapter *Adapter, message string)

//export requestAdapterCallback
func requestAdapterCallback(status C.WGPURequestAdapterStatus, adapter C.WGPUAdapter, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(requestAdapterCB)
	if ok {
		cb(RequestAdapterStatus(status), &Adapter{adapter}, C.GoString(message))
	}
}

type requestDeviceCB func(status RequestDeviceStatus, device *Device, message string)

//export requestDeviceCallback
func requestDeviceCallback(status C.WGPURequestDeviceStatus, device C.WGPUDevice, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(requestDeviceCB)
	if ok {
		cb(RequestDeviceStatus(status), &Device{ref: device}, C.GoString(message))
	}
}

type deviceUncapturedErrorCB func(typ ErrorType, message string)

//export deviceUncapturedErrorCallback
func deviceUncapturedErrorCallback(typ C.WGPUErrorType, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)

	device, ok := handle.Value().(*Device)
	if ok {
		device.storeErr(ErrorType(typ), C.GoString(message))
	}
}
