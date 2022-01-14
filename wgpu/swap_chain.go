package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

type SwapChain struct{ ref C.WGPUSwapChain }

func (p *SwapChain) GetCurrentTextureView() *TextureView {
	ref := C.wgpuSwapChainGetCurrentTextureView(p.ref)
	if ref == nil {
		panic("Failed to acquire TextureView")
	}
	return &TextureView{ref}
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}
