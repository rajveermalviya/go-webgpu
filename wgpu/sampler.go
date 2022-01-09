package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type Sampler struct{ ref C.WGPUSampler }

func (p *Sampler) Drop() {
	C.wgpuSamplerDrop(p.ref)
}
