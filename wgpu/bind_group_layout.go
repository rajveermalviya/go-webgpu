package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type BindGroupLayout struct{ ref C.WGPUBindGroupLayout }

func (p *BindGroupLayout) Drop() {
	C.wgpuBindGroupLayoutDrop(p.ref)
}
