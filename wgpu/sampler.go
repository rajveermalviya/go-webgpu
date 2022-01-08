package wgpu

/*

#include "wrapper.h"

*/
import "C"

type Sampler struct{ ref C.WGPUSampler }

func (p *Sampler) Drop() {
	C.wgpuSamplerDrop(p.ref)
}
