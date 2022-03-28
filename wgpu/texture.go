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

type Texture struct{ ref C.WGPUTexture }

type TextureViewDescriptor struct {
	Label           string
	Format          TextureFormat
	Dimension       TextureViewDimension
	BaseMipLevel    uint32
	MipLevelCount   uint32
	BaseArrayLayer  uint32
	ArrayLayerCount uint32
	Aspect          TextureAspect
}

func (p *Texture) CreateView(descriptor *TextureViewDescriptor) *TextureView {
	var desc C.WGPUTextureViewDescriptor

	if descriptor != nil {
		desc = C.WGPUTextureViewDescriptor{
			format:          C.WGPUTextureFormat(descriptor.Format),
			dimension:       C.WGPUTextureViewDimension(descriptor.Dimension),
			baseMipLevel:    C.uint32_t(descriptor.BaseMipLevel),
			mipLevelCount:   C.uint32_t(descriptor.MipLevelCount),
			baseArrayLayer:  C.uint32_t(descriptor.BaseArrayLayer),
			arrayLayerCount: C.uint32_t(descriptor.ArrayLayerCount),
			aspect:          C.WGPUTextureAspect(descriptor.Aspect),
		}

		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}
	}

	ref := C.wgpuTextureCreateView(p.ref, &desc)
	runtime.KeepAlive(p)

	if ref == nil {
		panic("Failed to acquire TextureView")
	}

	textureView := &TextureView{ref}
	runtime.SetFinalizer(textureView, textureViewFinalizer)
	return textureView
}

func (p *Texture) Destroy() {
	C.wgpuTextureDestroy(p.ref)
	runtime.KeepAlive(p)
}

func (p *Texture) AsImageCopy() *ImageCopyTexture {
	return &ImageCopyTexture{
		Texture:  p,
		MipLevel: 0,
		Origin:   Origin3D{},
		Aspect:   TextureAspect_All,
	}
}

func textureFinalizer(p *Texture) {
	C.wgpuTextureDrop(p.ref)
}
