package wgpu

/*

#include "./lib/webgpu.h"
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
	bufSlice := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(bufSlice), size)
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

func ByteStoUint32S(src []byte) []uint32 {
	const s = int(unsafe.Sizeof(uint32(0)))

	l := len(src)
	if l == 0 {
		return nil
	}
	if l%s != 0 {
		panic("invalid src")
	}

	return unsafe.Slice((*uint32)(unsafe.Pointer(&src[0])), l/s)
}

func Uint32StoByteS(src []uint32) []byte {
	const s = int(unsafe.Sizeof(uint32(0)))

	l := len(src)
	if l == 0 {
		return nil
	}

	return unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), l*s)
}
