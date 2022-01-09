package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type CommandBuffer struct{ ref C.WGPUCommandBuffer }

func (p *CommandBuffer) Drop() {
	C.wgpuCommandBufferDrop(p.ref)
}
