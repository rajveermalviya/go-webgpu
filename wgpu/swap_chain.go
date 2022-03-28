package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"
import "runtime"

type SwapChain struct {
	ref    C.WGPUSwapChain
	device *Device
}

func (p *SwapChain) GetCurrentTextureView() (*TextureView, error) {
	ref := C.wgpuSwapChainGetCurrentTextureView(p.ref)
	runtime.KeepAlive(p)

	err := p.device.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire TextureView")
	}

	textureView := &TextureView{ref}
	runtime.SetFinalizer(textureView, textureViewFinalizer)
	return textureView, nil
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
	runtime.KeepAlive(p)
}
