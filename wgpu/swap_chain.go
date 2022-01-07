package wgpu

/*

#include "wrapper.h"

*/
import "C"

type SwapChain struct{ ref C.WGPUSwapChain }

func (p *SwapChain) GetCurrentTextureView() TextureView {
	return TextureView(C.wgpuSwapChainGetCurrentTextureView(p.ref))
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}
