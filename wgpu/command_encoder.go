package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

*/
import "C"
import "unsafe"

type CommandEncoder struct {
	ref C.WGPUCommandEncoder
}

type ComputePassDescriptor struct {
	Label string

	// unused in wgpu
	// TimestampWrites []ComputePassTimestampWrite
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc *C.WGPUComputePassDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc = &C.WGPUComputePassDescriptor{
			label: label,
		}
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, desc)
	if ref == nil {
		panic("Failed to acquire ComputePassEncoder")
	}

	return &ComputePassEncoder{ref}
}

type RenderPassColorAttachment struct {
	View          *TextureView
	ResolveTarget *TextureView
	LoadOp        LoadOp
	StoreOp       StoreOp
	ClearValue    Color
}

type RenderPassDepthStencilAttachment struct {
	View              *TextureView
	DepthLoadOp       LoadOp
	DepthStoreOp      StoreOp
	DepthClearValue   float32
	DepthReadOnly     bool
	StencilLoadOp     LoadOp
	StencilStoreOp    StoreOp
	StencilClearValue uint32
	StencilReadOnly   bool
}

type RenderPassDescriptor struct {
	Label                  string
	ColorAttachments       []RenderPassColorAttachment
	DepthStencilAttachment *RenderPassDepthStencilAttachment

	// unused in wgpu
	// 	OcclusionQuerySet      QuerySet
	// 	TimestampWrites        []RenderPassTimestampWrite
}

func (p *CommandEncoder) BeginRenderPass(descriptor *RenderPassDescriptor) *RenderPassEncoder {
	var desc C.WGPURenderPassDescriptor

	if descriptor != nil {
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
					clearValue: C.WGPUColor{
						r: C.double(v.ClearValue.R),
						g: C.double(v.ClearValue.G),
						b: C.double(v.ClearValue.B),
						a: C.double(v.ClearValue.A),
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
			depthStencilAttachment.depthClearValue = C.float(descriptor.DepthStencilAttachment.DepthClearValue)
			depthStencilAttachment.depthReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)
			depthStencilAttachment.stencilLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.StencilLoadOp)
			depthStencilAttachment.stencilStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.StencilStoreOp)
			depthStencilAttachment.stencilClearValue = C.uint32_t(descriptor.DepthStencilAttachment.StencilClearValue)
			depthStencilAttachment.stencilReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)

			desc.depthStencilAttachment = depthStencilAttachment
		}
	}

	ref := C.wgpuCommandEncoderBeginRenderPass(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire RenderPassEncoder")
	}

	return &RenderPassEncoder{ref}
}

func (p *CommandEncoder) ClearBuffer(buffer *Buffer, offset uint64, size uint64) {
	C.wgpuCommandEncoderClearBuffer(
		p.ref,
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)
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

type ImageCopyBuffer struct {
	Layout TextureDataLayout
	Buffer *Buffer
}

func (p *CommandEncoder) CopyBufferToTexture(source *ImageCopyBuffer, destination *ImageCopyTexture, copySize *Extent3D) {
	var src C.WGPUImageCopyBuffer
	if source != nil {
		if source.Buffer != nil {
			src.buffer = source.Buffer.ref
		}
		src.layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(source.Layout.Offset),
			bytesPerRow:  C.uint32_t(source.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(source.Layout.RowsPerImage),
		}
	}

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

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyBufferToTexture(p.ref, &src, &dst, &cpySize)
}

func (p *CommandEncoder) CopyTextureToBuffer(source *ImageCopyTexture, destination *ImageCopyBuffer, copySize *Extent3D) {
	var src C.WGPUImageCopyTexture
	if source != nil {
		src = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUImageCopyBuffer
	if destination != nil {
		if destination.Buffer != nil {
			dst.buffer = destination.Buffer.ref
		}
		dst.layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(destination.Layout.Offset),
			bytesPerRow:  C.uint32_t(destination.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(destination.Layout.RowsPerImage),
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyTextureToBuffer(p.ref, &src, &dst, &cpySize)
}

func (p *CommandEncoder) CopyTextureToTexture(source *ImageCopyTexture, destination *ImageCopyTexture, copySize *Extent3D) {
	var src C.WGPUImageCopyTexture
	if source != nil {
		src = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

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

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyTextureToTexture(p.ref, &src, &dst, &cpySize)
}

type CommandBufferDescriptor struct {
	Label string
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) *CommandBuffer {
	var desc *C.WGPUCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc = &C.WGPUCommandBufferDescriptor{
			label: label,
		}
	}

	ref := C.wgpuCommandEncoderFinish(p.ref, desc)
	if ref == nil {
		panic("Failed to acquire CommandBuffer")
	}

	return &CommandBuffer{ref}
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	C.wgpuCommandEncoderInsertDebugMarker(p.ref, markerLabelStr)
}

func (p *CommandEncoder) PopDebugGroup() {
	C.wgpuCommandEncoderPopDebugGroup(p.ref)
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	C.wgpuCommandEncoderPushDebugGroup(p.ref, groupLabelStr)
}

func (p *CommandEncoder) Drop() {
	C.wgpuCommandEncoderDrop(p.ref)
}
