package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type ComputePipeline struct{ ref C.WGPUComputePipeline }

func computePipelineFinalizer(p *ComputePipeline) {
	C.wgpuComputePipelineDrop(p.ref)
}
