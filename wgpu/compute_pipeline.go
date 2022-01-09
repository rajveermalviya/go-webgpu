package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type ComputePipeline struct{ ref C.WGPUComputePipeline }

func (p *ComputePipeline) Drop() {
	C.wgpuComputePipelineDrop(p.ref)
}
