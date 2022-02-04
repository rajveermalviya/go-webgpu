package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type ShaderModule struct{ ref C.WGPUShaderModule }

func (p *ShaderModule) Drop() {
	C.wgpuShaderModuleDrop(p.ref)
}
