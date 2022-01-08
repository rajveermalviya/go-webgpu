package wgpu

/*

#include "wrapper.h"

*/
import "C"

type BindGroup struct{ ref C.WGPUBindGroup }

func (p *BindGroup) Drop() {
	C.wgpuBindGroupDrop(p.ref)
}
