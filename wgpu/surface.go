package wgpu

/*

#include "wrapper.h"

*/
import "C"

type Surface struct{ ref C.WGPUSurface }

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	return TextureFormat(C.wgpuSurfaceGetPreferredFormat(p.ref, adapter.ref))
}
