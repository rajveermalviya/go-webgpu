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

type RenderPassEncoder struct{ ref C.WGPURenderPassEncoder }

func (p *RenderPassEncoder) SetPushConstants(stages ShaderStage, offset uint32, data []byte) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	C.wgpuRenderPassEncoderSetPushConstants(
		p.ref,
		C.WGPUShaderStageFlags(stages),
		C.uint32_t(offset),
		C.uint32_t(size),
		buf,
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	C.wgpuRenderPassEncoderDraw(p.ref,
		C.uint32_t(vertexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstVertex),
		C.uint32_t(firstInstance),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	C.wgpuRenderPassEncoderDrawIndexed(p.ref,
		C.uint32_t(indexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstIndex),
		C.int32_t(baseVertex),
		C.uint32_t(firstInstance),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderPassEncoderDrawIndexedIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
	runtime.KeepAlive(p)
	runtime.KeepAlive(indirectBuffer)
}

func (p *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderPassEncoderDrawIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
	runtime.KeepAlive(p)
	runtime.KeepAlive(indirectBuffer)
}

func (p *RenderPassEncoder) End() {
	C.wgpuRenderPassEncoderEnd(p.ref)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		C.wgpuRenderPassEncoderSetBindGroup(
			p.ref,
			C.uint32_t(groupIndex),
			group.ref,
			0,
			nil,
		)
	} else {
		C.wgpuRenderPassEncoderSetBindGroup(
			p.ref,
			C.uint32_t(groupIndex),
			group.ref,
			C.uint32_t(dynamicOffsetCount),
			(*C.uint32_t)(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
	runtime.KeepAlive(p)
	runtime.KeepAlive(group)
}

func (p *RenderPassEncoder) SetBlendConstant(color *Color) {
	C.wgpuRenderPassEncoderSetBlendConstant(p.ref, &C.WGPUColor{
		r: C.double(color.R),
		g: C.double(color.G),
		b: C.double(color.B),
		a: C.double(color.A),
	})
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	C.wgpuRenderPassEncoderSetIndexBuffer(
		p.ref,
		buffer.ref,
		C.WGPUIndexFormat(format),
		C.uint64_t(offset),
		C.uint64_t(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	C.wgpuRenderPassEncoderSetPipeline(p.ref, pipeline.ref)
	runtime.KeepAlive(p)
	runtime.KeepAlive(pipeline)
}

func (p *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	C.wgpuRenderPassEncoderSetScissorRect(
		p.ref,
		C.uint32_t(x),
		C.uint32_t(y),
		C.uint32_t(width),
		C.uint32_t(height),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetStencilReference(reference uint32) {
	C.wgpuRenderPassEncoderSetStencilReference(p.ref, C.uint32_t(reference))
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	C.wgpuRenderPassEncoderSetVertexBuffer(
		p.ref,
		C.uint32_t(slot),
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {
	C.wgpuRenderPassEncoderSetViewport(
		p.ref,
		C.float(x),
		C.float(y),
		C.float(width),
		C.float(height),
		C.float(minDepth),
		C.float(maxDepth),
	)
	runtime.KeepAlive(p)
}
