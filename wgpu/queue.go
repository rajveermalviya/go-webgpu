package wgpu

/*

#include <stdlib.h>

#include "./lib/wgpu.h"

*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Queue struct{ ref C.WGPUQueue }

func (p *Queue) Submit(commands ...*CommandBuffer) {
	commandCount := len(commands)
	if commandCount == 0 {
		C.wgpuQueueSubmit(p.ref, 0, nil)
	} else {
		commandRefs := C.malloc(C.size_t(commandCount) * C.size_t(unsafe.Sizeof(C.WGPUCommandBuffer(nil))))
		defer C.free(commandRefs)

		commandRefsSlice := unsafe.Slice((*C.WGPUCommandBuffer)(commandRefs), commandCount)

		for i, v := range commands {
			commandRefsSlice[i] = v.ref
		}

		C.wgpuQueueSubmit(
			p.ref,
			C.uint32_t(commandCount),
			(*C.WGPUCommandBuffer)(commandRefs),
		)
	}

	runtime.KeepAlive(p)
	runtime.KeepAlive(commands)
	for _, v := range commands {
		runtime.KeepAlive(v)
	}
}

func (p *Queue) WriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	C.wgpuQueueWriteBuffer(p.ref, buffer.ref, C.uint64_t(bufferOffset), buf, C.size_t(size))
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	var dst C.WGPUImageCopyTexture
	if destination != nil {
		dst = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var layout C.WGPUTextureDataLayout
	if dataLayout != nil {
		layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(dataLayout.Offset),
			bytesPerRow:  C.uint32_t(dataLayout.BytesPerRow),
			rowsPerImage: C.uint32_t(dataLayout.RowsPerImage),
		}
	}

	var writeExtent C.WGPUExtent3D
	if writeSize != nil {
		writeExtent = C.WGPUExtent3D{
			width:              C.uint32_t(writeSize.Width),
			height:             C.uint32_t(writeSize.Height),
			depthOrArrayLayers: C.uint32_t(writeSize.DepthOrArrayLayers),
		}
	}

	C.wgpuQueueWriteTexture(p.ref, &dst, buf, C.size_t(size), &layout, &writeExtent)
	runtime.KeepAlive(p)
	if destination != nil {
		runtime.KeepAlive(destination.Texture)
	}
}
