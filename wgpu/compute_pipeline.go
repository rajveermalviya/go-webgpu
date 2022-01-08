package wgpu

/*

#include "wrapper.h"

*/
import "C"

type ComputePipeline struct{ ref C.WGPUComputePipeline }

func (p *ComputePipeline) Drop() {
	C.wgpuComputePipelineDrop(p.ref)
}
