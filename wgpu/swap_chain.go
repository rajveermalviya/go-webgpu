package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

*/
import "C"

type SwapChain struct {
	ref    C.WGPUSwapChain
	device *Device
}

func (p *SwapChain) GetCurrentTextureView() (*TextureView, error) {
	ref := C.wgpuSwapChainGetCurrentTextureView(p.ref)
	err := p.device.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire TextureView")
	}

	textureView := &TextureView{ref}
	return textureView, nil
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}

func (p *SwapChain) Drop() {
	C.wgpuSwapChainDrop(p.ref)
}
