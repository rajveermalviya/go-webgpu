package wgpu

/*

#include "wrapper.h"

*/
import "C"

type TextureView struct{ ref C.WGPUTextureView }

func (p *TextureView) Drop() {
	C.wgpuTextureViewDrop(p.ref)
}
