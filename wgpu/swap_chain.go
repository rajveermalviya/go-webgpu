package wgpu

/*

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
	return &TextureView{ref}, nil
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}
