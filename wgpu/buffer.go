package wgpu

/*

#include "./lib/wgpu.h"

extern void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata);

*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

type Buffer struct{ ref C.WGPUBuffer }

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	runtime.KeepAlive(p)
	return unsafe.Slice((*byte)(buf), size)
}

func (p *Buffer) Unmap() {
	C.wgpuBufferUnmap(p.ref)
	runtime.KeepAlive(p)
}

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
	runtime.KeepAlive(p)
}

type BufferMapCallback func(BufferMapAsyncStatus)

func (p *Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) {
	handle := cgo.NewHandle(callback)

	C.wgpuBufferMapAsync(
		p.ref,
		C.WGPUMapModeFlags(mode),
		C.size_t(offset),
		C.size_t(size),
		(C.WGPUBufferMapCallback)(C.bufferMapCallback_cgo),
		unsafe.Pointer(&handle),
	)
	runtime.KeepAlive(p)
}

func bufferFinalizer(p *Buffer) {
	C.wgpuBufferDrop(p.ref)
}

func FromBytes[E any](src []byte, zeroElm E) []E {
	l := uintptr(len(src))
	if l == 0 {
		return nil
	}

	elmSize := unsafe.Sizeof(zeroElm)
	if l%elmSize != 0 {
		panic("invalid src")
	}

	return unsafe.Slice((*E)(unsafe.Pointer(&src[0])), l/elmSize)
}

func ToBytes[E any](src []E) []byte {
	l := uintptr(len(src))
	if l == 0 {
		return nil
	}

	elmSize := unsafe.Sizeof(src[0])
	return unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), l*elmSize)
}
