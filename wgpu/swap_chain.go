package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_error_callback_c(WGPUErrorType type, char const * message, void * userdata);

static inline WGPUTextureView gowebgpu_swap_chain_get_current_texture_view(WGPUSwapChain swapChain, WGPUDevice device, void * error_userdata) {
	WGPUTextureView ref = NULL;
	wgpuDevicePushErrorScope(device, WGPUErrorFilter_Validation);
	ref = wgpuSwapChainGetCurrentTextureView(swapChain);
	wgpuDevicePopErrorScope(device, gowebgpu_error_callback_c, error_userdata);
	return ref;
}

static inline void gowebgpu_swap_chain_release(WGPUSwapChain swapChain, WGPUDevice device) {
	wgpuDeviceRelease(device);
	wgpuSwapChainRelease(swapChain);
}

*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type SwapChain struct {
	deviceRef C.WGPUDevice
	ref       C.WGPUSwapChain
}

func (p *SwapChain) GetCurrentTextureView() (*TextureView, error) {
	var err error = nil
	var cb errorCallback = func(_ ErrorType, message string) {
		err = errors.New("wgpu.(*SwapChain).GetCurrentTextureView(): " + message)
	}
	errorCallbackHandle := cgo.NewHandle(cb)
	defer errorCallbackHandle.Delete()

	ref := C.gowebgpu_swap_chain_get_current_texture_view(
		p.ref,
		p.deviceRef,
		unsafe.Pointer(&errorCallbackHandle),
	)
	if err != nil {
		if ref != nil {
			C.wgpuTextureViewRelease(ref)
		}
		return nil, err
	}

	return &TextureView{ref}, nil
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}

func (p *SwapChain) Release() {
	C.gowebgpu_swap_chain_release(p.ref, p.deviceRef)
}
