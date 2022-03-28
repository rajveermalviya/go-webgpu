package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"
import "runtime"

type Surface struct{ ref C.WGPUSurface }

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format := C.wgpuSurfaceGetPreferredFormat(p.ref, adapter.ref)
	runtime.KeepAlive(p)
	runtime.KeepAlive(adapter)
	return TextureFormat(format)
}
