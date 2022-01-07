package wgpu

/*

#include "wrapper.h"

*/
import "C"
import "unsafe"

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

func (p *Texture) CreateView(descriptor TextureViewDescriptor) TextureView {
	desc := C.WGPUTextureViewDescriptor{
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

	return TextureView(C.wgpuTextureCreateView(p.ref, &desc))
}

func (p *Texture) Destroy() {
	C.wgpuTextureDestroy(p.ref)
}
