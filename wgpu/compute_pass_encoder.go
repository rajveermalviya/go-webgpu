package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"
import "unsafe"

type ComputePassEncoder struct{ ref C.WGPUComputePassEncoder }

func (p *ComputePassEncoder) Dispatch(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	C.wgpuComputePassEncoderDispatch(p.ref, C.uint32_t(workgroupCountX), C.uint32_t(workgroupCountY), C.uint32_t(workgroupCountZ))
}

func (p *ComputePassEncoder) DispatchIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuComputePassEncoderDispatchIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
}

func (p *ComputePassEncoder) End() {
	C.wgpuComputePassEncoderEnd(p.ref)
}

func (p *ComputePassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		C.wgpuComputePassEncoderSetBindGroup(p.ref, C.uint32_t(groupIndex), group.ref, 0, nil)
	} else {
		C.wgpuComputePassEncoderSetBindGroup(
			p.ref, C.uint32_t(groupIndex), group.ref,
			C.uint32_t(dynamicOffsetCount), (*C.uint32_t)(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
}

func (p *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	C.wgpuComputePassEncoderSetPipeline(p.ref, pipeline.ref)
}
