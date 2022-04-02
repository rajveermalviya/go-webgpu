package wgpu

/*

#include "./lib/wgpu.h"

*/
import "C"

type (
	PipelineLayout      struct{ ref C.WGPUPipelineLayout }
	QuerySet            struct{ ref C.WGPUQuerySet }
	RenderBundleEncoder struct{ ref C.WGPURenderBundleEncoder }
	CommandBuffer       struct{ ref C.WGPUCommandBuffer }
)

type Limits struct {
	MaxTextureDimension1D                     uint32
	MaxTextureDimension2D                     uint32
	MaxTextureDimension3D                     uint32
	MaxTextureArrayLayers                     uint32
	MaxBindGroups                             uint32
	MaxDynamicUniformBuffersPerPipelineLayout uint32
	MaxDynamicStorageBuffersPerPipelineLayout uint32
	MaxSampledTexturesPerShaderStage          uint32
	MaxSamplersPerShaderStage                 uint32
	MaxStorageBuffersPerShaderStage           uint32
	MaxStorageTexturesPerShaderStage          uint32
	MaxUniformBuffersPerShaderStage           uint32
	MaxUniformBufferBindingSize               uint64
	MaxStorageBufferBindingSize               uint64
	MinUniformBufferOffsetAlignment           uint32
	MinStorageBufferOffsetAlignment           uint32
	MaxVertexBuffers                          uint32
	MaxVertexAttributes                       uint32
	MaxVertexBufferArrayStride                uint32
	MaxInterStageShaderComponents             uint32
	MaxComputeWorkgroupStorageSize            uint32
	MaxComputeInvocationsPerWorkgroup         uint32
	MaxComputeWorkgroupSizeX                  uint32
	MaxComputeWorkgroupSizeY                  uint32
	MaxComputeWorkgroupSizeZ                  uint32
	MaxComputeWorkgroupsPerDimension          uint32
}

func limitsFromC(limits C.WGPULimits) Limits {
	return Limits{
		MaxTextureDimension1D:                     uint32(limits.maxTextureDimension1D),
		MaxTextureDimension2D:                     uint32(limits.maxTextureDimension2D),
		MaxTextureDimension3D:                     uint32(limits.maxTextureDimension3D),
		MaxTextureArrayLayers:                     uint32(limits.maxTextureArrayLayers),
		MaxBindGroups:                             uint32(limits.maxBindGroups),
		MaxDynamicUniformBuffersPerPipelineLayout: uint32(limits.maxDynamicUniformBuffersPerPipelineLayout),
		MaxDynamicStorageBuffersPerPipelineLayout: uint32(limits.maxDynamicStorageBuffersPerPipelineLayout),
		MaxSampledTexturesPerShaderStage:          uint32(limits.maxSampledTexturesPerShaderStage),
		MaxSamplersPerShaderStage:                 uint32(limits.maxSamplersPerShaderStage),
		MaxStorageBuffersPerShaderStage:           uint32(limits.maxStorageBuffersPerShaderStage),
		MaxStorageTexturesPerShaderStage:          uint32(limits.maxStorageTexturesPerShaderStage),
		MaxUniformBuffersPerShaderStage:           uint32(limits.maxUniformBuffersPerShaderStage),
		MaxUniformBufferBindingSize:               uint64(limits.maxUniformBufferBindingSize),
		MaxStorageBufferBindingSize:               uint64(limits.maxStorageBufferBindingSize),
		MinUniformBufferOffsetAlignment:           uint32(limits.minUniformBufferOffsetAlignment),
		MinStorageBufferOffsetAlignment:           uint32(limits.minStorageBufferOffsetAlignment),
		MaxVertexBuffers:                          uint32(limits.maxVertexBuffers),
		MaxVertexAttributes:                       uint32(limits.maxVertexAttributes),
		MaxVertexBufferArrayStride:                uint32(limits.maxVertexBufferArrayStride),
		MaxInterStageShaderComponents:             uint32(limits.maxInterStageShaderComponents),
		MaxComputeWorkgroupStorageSize:            uint32(limits.maxComputeWorkgroupStorageSize),
		MaxComputeInvocationsPerWorkgroup:         uint32(limits.maxComputeInvocationsPerWorkgroup),
		MaxComputeWorkgroupSizeX:                  uint32(limits.maxComputeWorkgroupSizeX),
		MaxComputeWorkgroupSizeY:                  uint32(limits.maxComputeWorkgroupSizeY),
		MaxComputeWorkgroupSizeZ:                  uint32(limits.maxComputeWorkgroupSizeZ),
		MaxComputeWorkgroupsPerDimension:          uint32(limits.maxComputeWorkgroupsPerDimension),
	}
}

func (l Limits) toC() C.WGPULimits {
	return C.WGPULimits{
		maxTextureDimension1D:                     C.uint32_t(l.MaxTextureDimension1D),
		maxTextureDimension2D:                     C.uint32_t(l.MaxTextureDimension2D),
		maxTextureDimension3D:                     C.uint32_t(l.MaxTextureDimension3D),
		maxTextureArrayLayers:                     C.uint32_t(l.MaxTextureArrayLayers),
		maxBindGroups:                             C.uint32_t(l.MaxBindGroups),
		maxDynamicUniformBuffersPerPipelineLayout: C.uint32_t(l.MaxDynamicUniformBuffersPerPipelineLayout),
		maxDynamicStorageBuffersPerPipelineLayout: C.uint32_t(l.MaxDynamicStorageBuffersPerPipelineLayout),
		maxSampledTexturesPerShaderStage:          C.uint32_t(l.MaxSampledTexturesPerShaderStage),
		maxSamplersPerShaderStage:                 C.uint32_t(l.MaxSamplersPerShaderStage),
		maxStorageBuffersPerShaderStage:           C.uint32_t(l.MaxStorageBuffersPerShaderStage),
		maxStorageTexturesPerShaderStage:          C.uint32_t(l.MaxStorageTexturesPerShaderStage),
		maxUniformBuffersPerShaderStage:           C.uint32_t(l.MaxUniformBuffersPerShaderStage),
		maxUniformBufferBindingSize:               C.uint64_t(l.MaxUniformBufferBindingSize),
		maxStorageBufferBindingSize:               C.uint64_t(l.MaxStorageBufferBindingSize),
		minUniformBufferOffsetAlignment:           C.uint32_t(l.MinUniformBufferOffsetAlignment),
		minStorageBufferOffsetAlignment:           C.uint32_t(l.MinStorageBufferOffsetAlignment),
		maxVertexBuffers:                          C.uint32_t(l.MaxVertexBuffers),
		maxVertexAttributes:                       C.uint32_t(l.MaxVertexAttributes),
		maxVertexBufferArrayStride:                C.uint32_t(l.MaxVertexBufferArrayStride),
		maxInterStageShaderComponents:             C.uint32_t(l.MaxInterStageShaderComponents),
		maxComputeWorkgroupStorageSize:            C.uint32_t(l.MaxComputeWorkgroupStorageSize),
		maxComputeInvocationsPerWorkgroup:         C.uint32_t(l.MaxComputeInvocationsPerWorkgroup),
		maxComputeWorkgroupSizeX:                  C.uint32_t(l.MaxComputeWorkgroupSizeX),
		maxComputeWorkgroupSizeY:                  C.uint32_t(l.MaxComputeWorkgroupSizeY),
		maxComputeWorkgroupSizeZ:                  C.uint32_t(l.MaxComputeWorkgroupSizeZ),
		maxComputeWorkgroupsPerDimension:          C.uint32_t(l.MaxComputeWorkgroupsPerDimension),
	}
}

func (v VertexFormat) Size() uint64 {
	switch v {
	case VertexFormat_Uint8x2,
		VertexFormat_Sint8x2,
		VertexFormat_Unorm8x2,
		VertexFormat_Snorm8x2:
		return 2

	case VertexFormat_Uint8x4,
		VertexFormat_Sint8x4,
		VertexFormat_Unorm8x4,
		VertexFormat_Snorm8x4,
		VertexFormat_Uint16x2,
		VertexFormat_Sint16x2,
		VertexFormat_Unorm16x2,
		VertexFormat_Snorm16x2,
		VertexFormat_Float16x2,
		VertexFormat_Float32,
		VertexFormat_Uint32,
		VertexFormat_Sint32:
		return 4

	case VertexFormat_Uint16x4,
		VertexFormat_Sint16x4,
		VertexFormat_Unorm16x4,
		VertexFormat_Snorm16x4,
		VertexFormat_Float16x4,
		VertexFormat_Float32x2,
		VertexFormat_Uint32x2,
		VertexFormat_Sint32x2:
		return 8

	case VertexFormat_Float32x3,
		VertexFormat_Uint32x3,
		VertexFormat_Sint32x3:
		return 12

	case VertexFormat_Float32x4,
		VertexFormat_Uint32x4,
		VertexFormat_Sint32x4:
		return 16

	default:
		return 0
	}
}
