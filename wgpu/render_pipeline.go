package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type RenderPipeline struct{ ref C.WGPURenderPipeline }

func (p *RenderPipeline) Drop() {
	C.wgpuRenderPipelineDrop(p.ref)
}
