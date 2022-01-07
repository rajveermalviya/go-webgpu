package wgpu

/*

#include "wrapper.h"

*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Buffer struct{ ref C.WGPUBuffer }

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	bufSlice := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return C.GoBytes(bufSlice, C.int(size))
}

func (p *Buffer) Unmap() {
	C.wgpuBufferUnmap(p.ref)
}

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
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
	l := len(src)
	if l%4 != 0 {
		panic("invalid src")
	}
	l /= 4
	return (*[1 << 30]uint32)(unsafe.Pointer(&src[0]))[:l:l]
}

func Uint32StoByteS(src []uint32) []byte {
	l := len(src) * 4
	return (*[1 << 30]byte)(unsafe.Pointer(&src[0]))[:l:l]
}
