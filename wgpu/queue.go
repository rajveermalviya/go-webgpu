package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

*/
import "C"
import "unsafe"

type Queue struct {
	ref C.WGPUQueue
}

type SubmissionIndex uint64

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r := C.wgpuQueueSubmitForIndex(p.ref, 0, nil)
		return SubmissionIndex(r)
	}

	commandRefs := C.malloc(C.size_t(commandCount) * C.size_t(unsafe.Sizeof(C.WGPUCommandBuffer(nil))))
	defer C.free(commandRefs)

	commandRefsSlice := unsafe.Slice((*C.WGPUCommandBuffer)(commandRefs), commandCount)
	for i, v := range commands {
		commandRefsSlice[i] = v.ref
	}

	r := C.wgpuQueueSubmitForIndex(
		p.ref,
		C.uint32_t(commandCount),
		(*C.WGPUCommandBuffer)(commandRefs),
	)
	return SubmissionIndex(r)
}

func (p *Queue) WriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) {
	size := len(data)
	if size == 0 {
		C.wgpuQueueWriteBuffer(
			p.ref,
			buffer.ref,
			C.uint64_t(bufferOffset),
			nil,
			0,
		)
		return
	}

	C.wgpuQueueWriteBuffer(
		p.ref,
		buffer.ref,
		C.uint64_t(bufferOffset),
		unsafe.Pointer(&data[0]),
		C.size_t(size),
	)
}

func (p *Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) {
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

	size := len(data)
	if size == 0 {
		C.wgpuQueueWriteTexture(
			p.ref,
			&dst,
			nil,
			0,
			&layout,
			&writeExtent,
		)
		return
	}

	C.wgpuQueueWriteTexture(
		p.ref,
		&dst,
		unsafe.Pointer(&data[0]),
		C.size_t(size),
		&layout,
		&writeExtent,
	)
}
