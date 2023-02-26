package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_buffer_map_callback_c(WGPUBufferMapAsyncStatus status, void *userdata);

*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

type Buffer struct {
	ref C.WGPUBuffer
}

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
}

func (p *Buffer) GetMappedRange(offset, size uint) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(buf), size)
}

type BufferMapCallback func(BufferMapAsyncStatus)

//export gowebgpu_buffer_map_callback_go
func gowebgpu_buffer_map_callback_go(status C.WGPUBufferMapAsyncStatus, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(BufferMapCallback)
	if ok {
		cb(BufferMapAsyncStatus(status))
	}
}

func (p *Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) {
	handle := cgo.NewHandle(callback)

	C.wgpuBufferMapAsync(
		p.ref,
		C.WGPUMapModeFlags(mode),
		C.size_t(offset),
		C.size_t(size),
		(C.WGPUBufferMapCallback)(C.gowebgpu_buffer_map_callback_c),
		unsafe.Pointer(&handle),
	)
}

func (p *Buffer) Unmap() {
	C.wgpuBufferUnmap(p.ref)
}

func (p *Buffer) Drop() {
	C.wgpuBufferDrop(p.ref)
}
