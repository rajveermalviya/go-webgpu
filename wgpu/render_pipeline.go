package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type RenderPipeline struct{ ref C.WGPURenderPipeline }

func renderPipelineFinalizer(p *RenderPipeline) {
	C.wgpuRenderPipelineDrop(p.ref)
}
