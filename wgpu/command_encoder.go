package wgpu

/*

#include <stdlib.h>

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"
import "unsafe"

type CommandEncoder struct{ ref C.WGPUCommandEncoder }

type ComputePassTimestampWrite struct {
	QuerySet   QuerySet
	QueryIndex uint32
	Location   ComputePassTimestampLocation
}

type ComputePassDescriptor struct {
	Label string

	// unused in wgpu
	// TimestampWrites []ComputePassTimestampWrite
}

func (p *CommandEncoder) BeginComputePass(descriptor ComputePassDescriptor) *ComputePassEncoder {
	var desc C.WGPUComputePassDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, &desc)
	if ref == nil {
		return nil
	}

	return &ComputePassEncoder{ref}
}

type Color struct {
	R, G, B, A float64
}

type RenderPassColorAttachment struct {
	View          *TextureView
	ResolveTarget *TextureView
	LoadOp        LoadOp
	StoreOp       StoreOp
	ClearColor    Color
}

type RenderPassDepthStencilAttachment struct {
	View            *TextureView
	DepthLoadOp     LoadOp
	DepthStoreOp    StoreOp
	ClearDepth      float32
	DepthReadOnly   bool
	StencilLoadOp   LoadOp
	StencilStoreOp  StoreOp
	ClearStencil    uint32
	StencilReadOnly bool
}

type RenderPassTimestampWrite struct {
	QuerySet   QuerySet
	QueryIndex uint32
	Location   RenderPassTimestampLocation
}

type RenderPassDescriptor struct {
	Label                  string
	ColorAttachments       []RenderPassColorAttachment
	DepthStencilAttachment *RenderPassDepthStencilAttachment

	// unused in wgpu
	// 	OcclusionQuerySet      QuerySet
	// 	TimestampWrites        []RenderPassTimestampWrite
}

func (p *CommandEncoder) BeginRenderPass(descriptor RenderPassDescriptor) *RenderPassEncoder {
	var desc C.WGPURenderPassDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	colorAttachmentCount := len(descriptor.ColorAttachments)
	if colorAttachmentCount > 0 {
		colorAttachments := C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassColorAttachment{})) * C.size_t(colorAttachmentCount))
		defer C.free(colorAttachments)

		colorAttachmentsSlice := unsafe.Slice((*C.WGPURenderPassColorAttachment)(colorAttachments), colorAttachmentCount)

		for i, v := range descriptor.ColorAttachments {
			colorAttachment := C.WGPURenderPassColorAttachment{
				loadOp:  C.WGPULoadOp(v.LoadOp),
				storeOp: C.WGPUStoreOp(v.StoreOp),
				clearColor: C.WGPUColor{
					r: C.double(v.ClearColor.R),
					g: C.double(v.ClearColor.G),
					b: C.double(v.ClearColor.B),
					a: C.double(v.ClearColor.A),
				},
			}
			if v.View != nil {
				colorAttachment.view = v.View.ref
			}
			if v.ResolveTarget != nil {
				colorAttachment.resolveTarget = v.ResolveTarget.ref
			}

			colorAttachmentsSlice[i] = colorAttachment
		}

		desc.colorAttachmentCount = C.uint32_t(colorAttachmentCount)
		desc.colorAttachments = (*C.WGPURenderPassColorAttachment)(colorAttachments)
	}

	if descriptor.DepthStencilAttachment != nil {
		depthStencilAttachment := (*C.WGPURenderPassDepthStencilAttachment)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassDepthStencilAttachment{}))))
		defer C.free(unsafe.Pointer(depthStencilAttachment))

		if descriptor.DepthStencilAttachment.View != nil {
			depthStencilAttachment.view = descriptor.DepthStencilAttachment.View.ref
		}
		depthStencilAttachment.depthLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.DepthLoadOp)
		depthStencilAttachment.depthStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.DepthStoreOp)
		depthStencilAttachment.clearDepth = C.float(descriptor.DepthStencilAttachment.ClearDepth)
		depthStencilAttachment.depthReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)
		depthStencilAttachment.stencilLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.StencilLoadOp)
		depthStencilAttachment.stencilStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.StencilStoreOp)
		depthStencilAttachment.clearStencil = C.uint32_t(descriptor.DepthStencilAttachment.ClearStencil)
		depthStencilAttachment.stencilReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)

		desc.depthStencilAttachment = depthStencilAttachment
	}

	ref := C.wgpuCommandEncoderBeginRenderPass(p.ref, &desc)
	if ref == nil {
		return nil
	}
	return &RenderPassEncoder{ref}
}

func (p *CommandEncoder) CopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinatonOffset uint64, size uint64) {
	C.wgpuCommandEncoderCopyBufferToBuffer(
		p.ref,
		source.ref,
		C.uint64_t(sourceOffset),
		destination.ref,
		C.uint64_t(destinatonOffset),
		C.uint64_t(size),
	)
}

type TextureDataLayout struct {
	Offset       uint64
	BytesPerRow  uint32
	RowsPerImage uint32
}

type ImageCopyBuffer struct {
	Layout TextureDataLayout
	Buffer *Buffer
}

type Origin3D struct {
	X, Y, Z uint32
}

type ImageCopyTexture struct {
	Texture  *Texture
	MipLevel uint32
	Origin   Origin3D
	Aspect   TextureAspect
}

func (p *CommandEncoder) CopyBufferToTexture(source ImageCopyBuffer, destination ImageCopyTexture, copySize Extent3D) {
	C.wgpuCommandEncoderCopyBufferToTexture(
		p.ref,
		&C.WGPUImageCopyBuffer{
			layout: C.WGPUTextureDataLayout{
				offset:       C.uint64_t(source.Layout.Offset),
				bytesPerRow:  C.uint32_t(source.Layout.BytesPerRow),
				rowsPerImage: C.uint32_t(source.Layout.RowsPerImage),
			},
			buffer: source.Buffer.ref,
		},
		&C.WGPUImageCopyTexture{
			texture:  destination.Texture.ref,
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		},
		&C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		},
	)
}

func (p *CommandEncoder) CopyTextureToBuffer(source ImageCopyTexture, destination ImageCopyBuffer, copySize Extent3D) {
	C.wgpuCommandEncoderCopyTextureToBuffer(
		p.ref,
		&C.WGPUImageCopyTexture{
			texture:  source.Texture.ref,
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		},
		&C.WGPUImageCopyBuffer{
			layout: C.WGPUTextureDataLayout{
				offset:       C.uint64_t(destination.Layout.Offset),
				bytesPerRow:  C.uint32_t(destination.Layout.BytesPerRow),
				rowsPerImage: C.uint32_t(destination.Layout.RowsPerImage),
			},
			buffer: destination.Buffer.ref,
		},
		&C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		},
	)
}

func (p *CommandEncoder) CopyTextureToTexture(source ImageCopyTexture, destination ImageCopyTexture, copySize Extent3D) {
	C.wgpuCommandEncoderCopyTextureToTexture(
		p.ref,
		&C.WGPUImageCopyTexture{
			texture:  source.Texture.ref,
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		},
		&C.WGPUImageCopyTexture{
			texture:  destination.Texture.ref,
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		},
		&C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		},
	)
}

type CommandBufferDescriptor struct {
	Label string
}

func (p *CommandEncoder) Finish(descriptor CommandBufferDescriptor) *CommandBuffer {
	var desc C.WGPUCommandBufferDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuCommandEncoderFinish(p.ref, &desc)
	if ref == nil {
		return nil
	}
	return &CommandBuffer{ref}
}

func (p *CommandEncoder) Drop() {
	C.wgpuCommandEncoderDrop(p.ref)
}
