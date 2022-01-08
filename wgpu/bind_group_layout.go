package wgpu

/*

#include "wrapper.h"

*/
import "C"

type BindGroupLayout struct{ ref C.WGPUBindGroupLayout }

func (p *BindGroupLayout) Drop() {
	C.wgpuBindGroupLayoutDrop(p.ref)
}
