package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type BindGroup struct{ ref C.WGPUBindGroup }

func bindGroupFinalizer(p *BindGroup) {
	C.wgpuBindGroupDrop(p.ref)
}
