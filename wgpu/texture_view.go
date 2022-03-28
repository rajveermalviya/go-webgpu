package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type TextureView struct{ ref C.WGPUTextureView }

func textureViewFinalizer(p *TextureView) {
	C.wgpuTextureViewDrop(p.ref)
}
