package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type Sampler struct{ ref C.WGPUSampler }

func samplerFinalizer(p *Sampler) {
	C.wgpuSamplerDrop(p.ref)
}
