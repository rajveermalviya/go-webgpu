package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

*/
import "C"
import "unsafe"

type Surface struct {
	ref C.WGPUSurface
}

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format := C.wgpuSurfaceGetPreferredFormat(p.ref, adapter.ref)
	return TextureFormat(format)
}

func (p *Surface) GetSupportedFormats(adapter *Adapter) []TextureFormat {
	var size C.size_t
	formatsPtr := C.wgpuSurfaceGetSupportedFormats(p.ref, adapter.ref, &size)
	defer free[C.WGPUTextureFormat](unsafe.Pointer(formatsPtr), size)

	formatsSlice := unsafe.Slice((*TextureFormat)(formatsPtr), size)
	formats := make([]TextureFormat, size)
	copy(formats, formatsSlice)
	return formats
}

func (p *Surface) GetSupportedPresentModes(adapter *Adapter) []PresentMode {
	var size C.size_t
	modesPtr := C.wgpuSurfaceGetSupportedPresentModes(p.ref, adapter.ref, &size)
	defer free[C.WGPUPresentMode](unsafe.Pointer(modesPtr), size)

	modesSlice := unsafe.Slice((*PresentMode)(modesPtr), size)
	modes := make([]PresentMode, size)
	copy(modes, modesSlice)
	return modes
}

func (p *Surface) Drop() {
	C.wgpuSurfaceDrop(p.ref)
}
