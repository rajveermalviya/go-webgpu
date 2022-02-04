package wgpu

/*

#include "./lib/wgpu.h"

extern void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata);

*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Buffer struct{ ref C.WGPUBuffer }

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(buf), size)
}

func (p *Buffer) Unmap() {
	C.wgpuBufferUnmap(p.ref)
}

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
}

func (p *Buffer) Drop() {
	C.wgpuBufferDrop(p.ref)
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
