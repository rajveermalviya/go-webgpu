package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type TextureView struct{ ref C.WGPUTextureView }

func (p *TextureView) Drop() {
	C.wgpuTextureViewDrop(p.ref)
}
