package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type ShaderModule struct{ ref C.WGPUShaderModule }

func shaderModuleFinalizer(p *ShaderModule) {
	C.wgpuShaderModuleDrop(p.ref)
}
