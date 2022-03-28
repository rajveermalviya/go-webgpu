package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type BindGroupLayout struct{ ref C.WGPUBindGroupLayout }

func bindGroupLayoutFinalizer(p *BindGroupLayout) {
	C.wgpuBindGroupLayoutDrop(p.ref)
}
