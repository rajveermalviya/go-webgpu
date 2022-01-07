package wgpu

/*

#include "wrapper.h"

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
