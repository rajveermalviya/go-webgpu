package wgpu

/*

// Android
#cgo android,amd64 LDFLAGS: -L${SRCDIR}/lib/android/amd64 -lwgpu_native
#cgo android,386 LDFLAGS: -L${SRCDIR}/lib/android/386 -lwgpu_native
#cgo android,arm64 LDFLAGS: -L${SRCDIR}/lib/android/arm64 -lwgpu_native
#cgo android,arm LDFLAGS: -L${SRCDIR}/lib/android/arm -lwgpu_native

#cgo android LDFLAGS: -landroid -lm -llog

// Linux
#cgo linux,!android,amd64 LDFLAGS: -L${SRCDIR}/lib/linux/amd64 -lwgpu_native

#cgo linux,!android LDFLAGS: -lm -ldl

// Darwin
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin/amd64 -lwgpu_native
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin/arm64 -lwgpu_native

#cgo darwin LDFLAGS: -framework QuartzCore -framework Metal

// Windows
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows/amd64 -lwgpu_native

#cgo windows LDFLAGS: -ld3dcompiler_47 -lws2_32 -luserenv -lbcrypt

#include <stdio.h>
#include "./lib/wgpu.h"

#ifdef __ANDROID__
#include <android/log.h>
void logCallback_cgo(WGPULogLevel level, char const *msg) {
	switch (level) {
	case WGPULogLevel_Error:
		__android_log_write(ANDROID_LOG_ERROR, "wgpu", msg);
		break;
	case WGPULogLevel_Warn:
		__android_log_write(ANDROID_LOG_WARN, "wgpu", msg);
		break;
	default:
		__android_log_write(ANDROID_LOG_INFO, "wgpu", msg);
		break;
	}
}
#else
void logCallback_cgo(WGPULogLevel level, char const *msg) {
	char const *level_str;
	switch (level) {
	case WGPULogLevel_Error:
		level_str = "Error";
		break;
	case WGPULogLevel_Warn:
		level_str = "Warn";
		break;
	case WGPULogLevel_Info:
		level_str = "Info";
		break;
	case WGPULogLevel_Debug:
		level_str = "Debug";
		break;
	case WGPULogLevel_Trace:
		level_str = "Trace";
		break;
	default:
		level_str = "Unknown Level";
	}
	fprintf(stderr, "[wgpu] [%s] %s\n", level_str, msg);
}
#endif


*/
import "C"
import (
	"strconv"
	"unsafe"
)

const (
	ArrayLayerCountUndefined        = 0xffffffff
	CopyStrideUndefined             = 0xffffffff
	LimitU32Undefined        uint32 = 0xffffffff
	LimitU64Undefined        uint64 = 0xffffffffffffffff
	MipLevelCountUndefined          = 0xffffffff
	WholeMapSize                    = ^uint(0)
	WholeSize                       = 0xffffffffffffffff
)

type Version uint32

func (v Version) String() string {
	return "0x" + strconv.FormatUint(uint64(v), 8)
}

func init() {
	C.wgpuSetLogCallback(C.WGPULogCallback(C.logCallback_cgo), nil)
}

func SetLogLevel(level LogLevel) {
	C.wgpuSetLogLevel(C.WGPULogLevel(level))
}

func GetVersion() Version {
	return Version(C.wgpuGetVersion())
}

func free[T any](ptr unsafe.Pointer, len C.size_t) {
	var v T
	C.wgpuFree(
		unsafe.Pointer(ptr),
		len*C.size_t(unsafe.Sizeof(v)),
		C.size_t(unsafe.Alignof(v)),
	)
}

type (
	BindGroup       struct{ ref C.WGPUBindGroup }
	BindGroupLayout struct{ ref C.WGPUBindGroupLayout }
	CommandBuffer   struct{ ref C.WGPUCommandBuffer }
	PipelineLayout  struct{ ref C.WGPUPipelineLayout }
	QuerySet        struct{ ref C.WGPUQuerySet }
	RenderBundle    struct{ ref C.WGPURenderBundle }
	Sampler         struct{ ref C.WGPUSampler }
	ShaderModule    struct{ ref C.WGPUShaderModule }
	TextureView     struct{ ref C.WGPUTextureView }
)

func (p *BindGroup) Drop()       { C.wgpuBindGroupDrop(p.ref) }
func (p *BindGroupLayout) Drop() { C.wgpuBindGroupLayoutDrop(p.ref) }
func (p *CommandBuffer) Drop()   { C.wgpuCommandBufferDrop(p.ref) }
func (p *PipelineLayout) Drop()  { C.wgpuPipelineLayoutDrop(p.ref) }
func (p *QuerySet) Drop()        { C.wgpuQuerySetDrop(p.ref) }
func (p *RenderBundle) Drop()    { C.wgpuRenderBundleDrop(p.ref) }
func (p *Sampler) Drop()         { C.wgpuSamplerDrop(p.ref) }
func (p *ShaderModule) Drop()    { C.wgpuShaderModuleDrop(p.ref) }
func (p *TextureView) Drop()     { C.wgpuTextureViewDrop(p.ref) }

// common types

type Limits struct {
	MaxTextureDimension1D                     uint32
	MaxTextureDimension2D                     uint32
	MaxTextureDimension3D                     uint32
	MaxTextureArrayLayers                     uint32
	MaxBindGroups                             uint32
	MaxBindingsPerBindGroup                   uint32
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
	MaxBufferSize                             uint64
	MaxVertexAttributes                       uint32
	MaxVertexBufferArrayStride                uint32
	MaxInterStageShaderComponents             uint32
	MaxInterStageShaderVariables              uint32
	MaxColorAttachments                       uint32
	MaxComputeWorkgroupStorageSize            uint32
	MaxComputeInvocationsPerWorkgroup         uint32
	MaxComputeWorkgroupSizeX                  uint32
	MaxComputeWorkgroupSizeY                  uint32
	MaxComputeWorkgroupSizeZ                  uint32
	MaxComputeWorkgroupsPerDimension          uint32

	MaxPushConstantSize uint32
}

type Color struct {
	R, G, B, A float64
}

type Origin3D struct {
	X, Y, Z uint32
}

type ImageCopyTexture struct {
	Texture  *Texture
	MipLevel uint32
	Origin   Origin3D
	Aspect   TextureAspect
}

type TextureDataLayout struct {
	Offset       uint64
	BytesPerRow  uint32
	RowsPerImage uint32
}

type Extent3D struct {
	Width              uint32
	Height             uint32
	DepthOrArrayLayers uint32
}
