//go:build windows

package wgpu

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/windows"
)

func writeDllToCacheDir(dllPath string) {
	r, err := gzip.NewReader(bytes.NewReader(libwgpuDllCompressed))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	f, err := os.OpenFile(dllPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, r)
	if err1 := f.Close(); err1 != nil && err == nil {
		panic(err1)
	}
}

var lib = func() *windows.LazyDLL {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "go-webgpu")
	dllPath := filepath.Join(dir, "wgpu_native.dll")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	if _, err := os.Stat(dllPath); err != nil {
		// dll does not already exists
		writeDllToCacheDir(dllPath)
	} else {
		// dll already exists
		f, err := os.Open(dllPath)
		if err != nil {
			panic(err)
		}
		hash := sha256.New()
		_, err = io.Copy(hash, f)
		if err1 := f.Close(); err1 != nil && err == nil {
			panic(err1)
		}
		sum := hash.Sum(nil)
		if hex.EncodeToString(sum) != strings.Split(libwgpuDllSha256, " ")[0] {
			// dll hash doesn't match
			writeDllToCacheDir(dllPath)
		}
	}

	lib := windows.NewLazyDLL(dllPath)
	if err := lib.Load(); err != nil {
		panic(err)
	}
	return lib
}()

var (
	wgpuSetLogCallback         = lib.NewProc("wgpuSetLogCallback")
	wgpuSetLogLevel            = lib.NewProc("wgpuSetLogLevel")
	wgpuGetVersion             = lib.NewProc("wgpuGetVersion")
	wgpuGetResourceUsageString = lib.NewProc("wgpuGetResourceUsageString")
	wgpuGenerateReport         = lib.NewProc("wgpuGenerateReport")
	wgpuFree                   = lib.NewProc("wgpuFree")

	wgpuDeviceDrop          = lib.NewProc("wgpuDeviceDrop")
	wgpuBindGroupLayoutDrop = lib.NewProc("wgpuBindGroupLayoutDrop")
	wgpuBindGroupDrop       = lib.NewProc("wgpuBindGroupDrop")
	wgpuBufferDrop          = lib.NewProc("wgpuBufferDrop")
	wgpuComputePipelineDrop = lib.NewProc("wgpuComputePipelineDrop")
	wgpuRenderPipelineDrop  = lib.NewProc("wgpuRenderPipelineDrop")
	wgpuSamplerDrop         = lib.NewProc("wgpuSamplerDrop")
	wgpuShaderModuleDrop    = lib.NewProc("wgpuShaderModuleDrop")
	wgpuTextureViewDrop     = lib.NewProc("wgpuTextureViewDrop")
	wgpuTextureDrop         = lib.NewProc("wgpuTextureDrop")

	wgpuInstanceRequestAdapter                       = lib.NewProc("wgpuInstanceRequestAdapter")
	wgpuInstanceCreateSurface                        = lib.NewProc("wgpuInstanceCreateSurface")
	wgpuAdapterEnumerateFeatures                     = lib.NewProc("wgpuAdapterEnumerateFeatures")
	wgpuAdapterHasFeature                            = lib.NewProc("wgpuAdapterHasFeature")
	wgpuAdapterGetLimits                             = lib.NewProc("wgpuAdapterGetLimits")
	wgpuAdapterGetProperties                         = lib.NewProc("wgpuAdapterGetProperties")
	wgpuAdapterRequestDevice                         = lib.NewProc("wgpuAdapterRequestDevice")
	wgpuDeviceSetUncapturedErrorCallback             = lib.NewProc("wgpuDeviceSetUncapturedErrorCallback")
	wgpuDeviceEnumerateFeatures                      = lib.NewProc("wgpuDeviceEnumerateFeatures")
	wgpuDeviceHasFeature                             = lib.NewProc("wgpuDeviceHasFeature")
	wgpuDevicePoll                                   = lib.NewProc("wgpuDevicePoll")
	wgpuDeviceCreateBindGroupLayout                  = lib.NewProc("wgpuDeviceCreateBindGroupLayout")
	wgpuDeviceCreateBindGroup                        = lib.NewProc("wgpuDeviceCreateBindGroup")
	wgpuDeviceCreateBuffer                           = lib.NewProc("wgpuDeviceCreateBuffer")
	wgpuDeviceCreateCommandEncoder                   = lib.NewProc("wgpuDeviceCreateCommandEncoder")
	wgpuDeviceCreateComputePipeline                  = lib.NewProc("wgpuDeviceCreateComputePipeline")
	wgpuDeviceCreatePipelineLayout                   = lib.NewProc("wgpuDeviceCreatePipelineLayout")
	wgpuDeviceCreateRenderPipeline                   = lib.NewProc("wgpuDeviceCreateRenderPipeline")
	wgpuDeviceCreateSampler                          = lib.NewProc("wgpuDeviceCreateSampler")
	wgpuDeviceCreateShaderModule                     = lib.NewProc("wgpuDeviceCreateShaderModule")
	wgpuDeviceCreateSwapChain                        = lib.NewProc("wgpuDeviceCreateSwapChain")
	wgpuDeviceCreateTexture                          = lib.NewProc("wgpuDeviceCreateTexture")
	wgpuDeviceGetLimits                              = lib.NewProc("wgpuDeviceGetLimits")
	wgpuDeviceGetQueue                               = lib.NewProc("wgpuDeviceGetQueue")
	wgpuBufferGetMappedRange                         = lib.NewProc("wgpuBufferGetMappedRange")
	wgpuBufferUnmap                                  = lib.NewProc("wgpuBufferUnmap")
	wgpuBufferDestroy                                = lib.NewProc("wgpuBufferDestroy")
	wgpuBufferMapAsync                               = lib.NewProc("wgpuBufferMapAsync")
	wgpuCommandEncoderBeginComputePass               = lib.NewProc("wgpuCommandEncoderBeginComputePass")
	wgpuCommandEncoderBeginRenderPass                = lib.NewProc("wgpuCommandEncoderBeginRenderPass")
	wgpuCommandEncoderClearBuffer                    = lib.NewProc("wgpuCommandEncoderClearBuffer")
	wgpuCommandEncoderCopyBufferToBuffer             = lib.NewProc("wgpuCommandEncoderCopyBufferToBuffer")
	wgpuCommandEncoderCopyBufferToTexture            = lib.NewProc("wgpuCommandEncoderCopyBufferToTexture")
	wgpuCommandEncoderCopyTextureToBuffer            = lib.NewProc("wgpuCommandEncoderCopyTextureToBuffer")
	wgpuCommandEncoderCopyTextureToTexture           = lib.NewProc("wgpuCommandEncoderCopyTextureToTexture")
	wgpuCommandEncoderFinish                         = lib.NewProc("wgpuCommandEncoderFinish")
	wgpuCommandEncoderInsertDebugMarker              = lib.NewProc("wgpuCommandEncoderInsertDebugMarker")
	wgpuCommandEncoderPopDebugGroup                  = lib.NewProc("wgpuCommandEncoderPopDebugGroup")
	wgpuCommandEncoderPushDebugGroup                 = lib.NewProc("wgpuCommandEncoderPushDebugGroup")
	wgpuComputePassEncoderDispatchWorkgroups         = lib.NewProc("wgpuComputePassEncoderDispatchWorkgroups")
	wgpuComputePassEncoderDispatchWorkgroupsIndirect = lib.NewProc("wgpuComputePassEncoderDispatchWorkgroupsIndirect")
	wgpuComputePassEncoderEnd                        = lib.NewProc("wgpuComputePassEncoderEnd")
	wgpuComputePassEncoderSetBindGroup               = lib.NewProc("wgpuComputePassEncoderSetBindGroup")
	wgpuComputePassEncoderSetPipeline                = lib.NewProc("wgpuComputePassEncoderSetPipeline")
	wgpuComputePassEncoderInsertDebugMarker          = lib.NewProc("wgpuComputePassEncoderInsertDebugMarker")
	wgpuComputePassEncoderPopDebugGroup              = lib.NewProc("wgpuComputePassEncoderPopDebugGroup")
	wgpuComputePassEncoderPushDebugGroup             = lib.NewProc("wgpuComputePassEncoderPushDebugGroup")
	wgpuComputePipelineGetBindGroupLayout            = lib.NewProc("wgpuComputePipelineGetBindGroupLayout")
	wgpuQueueSubmitForIndex                          = lib.NewProc("wgpuQueueSubmitForIndex")
	wgpuQueueWriteBuffer                             = lib.NewProc("wgpuQueueWriteBuffer")
	wgpuQueueWriteTexture                            = lib.NewProc("wgpuQueueWriteTexture")
	wgpuRenderPassEncoderSetPushConstants            = lib.NewProc("wgpuRenderPassEncoderSetPushConstants")
	wgpuRenderPassEncoderDraw                        = lib.NewProc("wgpuRenderPassEncoderDraw")
	wgpuRenderPassEncoderDrawIndexed                 = lib.NewProc("wgpuRenderPassEncoderDrawIndexed")
	wgpuRenderPassEncoderDrawIndexedIndirect         = lib.NewProc("wgpuRenderPassEncoderDrawIndexedIndirect")
	wgpuRenderPassEncoderDrawIndirect                = lib.NewProc("wgpuRenderPassEncoderDrawIndirect")
	wgpuRenderPassEncoderEnd                         = lib.NewProc("wgpuRenderPassEncoderEnd")
	wgpuRenderPassEncoderSetBindGroup                = lib.NewProc("wgpuRenderPassEncoderSetBindGroup")
	wgpuRenderPassEncoderSetBlendConstant            = lib.NewProc("wgpuRenderPassEncoderSetBlendConstant")
	wgpuRenderPassEncoderSetIndexBuffer              = lib.NewProc("wgpuRenderPassEncoderSetIndexBuffer")
	wgpuRenderPassEncoderSetPipeline                 = lib.NewProc("wgpuRenderPassEncoderSetPipeline")
	wgpuRenderPassEncoderSetScissorRect              = lib.NewProc("wgpuRenderPassEncoderSetScissorRect")
	wgpuRenderPassEncoderSetStencilReference         = lib.NewProc("wgpuRenderPassEncoderSetStencilReference")
	wgpuRenderPassEncoderSetVertexBuffer             = lib.NewProc("wgpuRenderPassEncoderSetVertexBuffer")
	wgpuRenderPassEncoderSetViewport                 = lib.NewProc("wgpuRenderPassEncoderSetViewport")
	wgpuRenderPassEncoderInsertDebugMarker           = lib.NewProc("wgpuRenderPassEncoderInsertDebugMarker")
	wgpuRenderPassEncoderPopDebugGroup               = lib.NewProc("wgpuRenderPassEncoderPopDebugGroup")
	wgpuRenderPassEncoderPushDebugGroup              = lib.NewProc("wgpuRenderPassEncoderPushDebugGroup")
	wgpuRenderPipelineGetBindGroupLayout             = lib.NewProc("wgpuRenderPipelineGetBindGroupLayout")
	wgpuSurfaceGetPreferredFormat                    = lib.NewProc("wgpuSurfaceGetPreferredFormat")
	wgpuSurfaceGetSupportedFormats                   = lib.NewProc("wgpuSurfaceGetSupportedFormats")
	wgpuSwapChainGetCurrentTextureView               = lib.NewProc("wgpuSwapChainGetCurrentTextureView")
	wgpuSwapChainPresent                             = lib.NewProc("wgpuSwapChainPresent")
	wgpuTextureCreateView                            = lib.NewProc("wgpuTextureCreateView")
	wgpuTextureDestroy                               = lib.NewProc("wgpuTextureDestroy")
)

func init() {
	logCb = defaultlogCallback
	wgpuSetLogCallback.Call(logCallback)
}

var logCb LogCallback

var logCallback = windows.NewCallbackCDecl(func(level LogLevel, msg *byte) (_ uintptr) {
	logCb(level, gostring(msg))
	return
})

func SetLogCallback(f LogCallback) {
	logCb = f
}

func SetLogLevel(level LogLevel) {
	wgpuSetLogLevel.Call(uintptr(level))
}

func GetVersion() Version {
	r, _, _ := wgpuGetVersion.Call()
	return Version(r)
}

func GenerateReport() GlobalReport {
	var r wgpuGlobalReport
	wgpuGenerateReport.Call(uintptr(unsafe.Pointer(&r)))

	mapStorageReport := func(creport wgpuStorageReport) StorageReport {
		return StorageReport{
			NumOccupied: uint64(creport.numOccupied),
			NumVacant:   uint64(creport.numVacant),
			NumError:    uint64(creport.numError),
			ElementSize: uint64(creport.elementSize),
		}
	}

	mapHubReport := func(creport wgpuHubReport) *HubReport {
		return &HubReport{
			Adapters:         mapStorageReport(creport.adapters),
			Devices:          mapStorageReport(creport.devices),
			PipelineLayouts:  mapStorageReport(creport.pipelineLayouts),
			ShaderModules:    mapStorageReport(creport.shaderModules),
			BindGroupLayouts: mapStorageReport(creport.bindGroupLayouts),
			BindGroups:       mapStorageReport(creport.bindGroups),
			CommandBuffers:   mapStorageReport(creport.commandBuffers),
			RenderBundles:    mapStorageReport(creport.renderBundles),
			RenderPipelines:  mapStorageReport(creport.renderPipelines),
			ComputePipelines: mapStorageReport(creport.computePipelines),
			QuerySets:        mapStorageReport(creport.querySets),
			Buffers:          mapStorageReport(creport.buffers),
			Textures:         mapStorageReport(creport.textures),
			TextureViews:     mapStorageReport(creport.textureViews),
			Samplers:         mapStorageReport(creport.samplers),
		}
	}

	report := GlobalReport{
		Surfaces: mapStorageReport(r.surfaces),
	}

	switch r.backendType {
	case BackendType_Vulkan:
		report.Vulkan = mapHubReport(r.vulkan)
	case BackendType_Metal:
		report.Metal = mapHubReport(r.metal)
	case BackendType_D3D12:
		report.Dx12 = mapHubReport(r.dx12)
	case BackendType_D3D11:
		report.Dx11 = mapHubReport(r.dx11)
	case BackendType_OpenGL:
		report.Gl = mapHubReport(r.gl)
	}

	return report
}

func free[T any](ptr uintptr, len uintptr) {
	var v T
	wgpuFree.Call(
		ptr,
		len*unsafe.Sizeof(v),
		unsafe.Alignof(v),
	)
}

type (
	Adapter             struct{ ref wgpuAdapter }
	BindGroup           struct{ ref wgpuBindGroup }
	BindGroupLayout     struct{ ref wgpuBindGroupLayout }
	Buffer              struct{ ref wgpuBuffer }
	CommandBuffer       struct{ ref wgpuCommandBuffer }
	CommandEncoder      struct{ ref wgpuCommandEncoder }
	ComputePassEncoder  struct{ ref wgpuComputePassEncoder }
	ComputePipeline     struct{ ref wgpuComputePipeline }
	PipelineLayout      struct{ ref wgpuPipelineLayout }
	QuerySet            struct{ ref wgpuQuerySet }
	Queue               struct{ ref wgpuQueue }
	RenderBundleEncoder struct{ ref wgpuRenderBundleEncoder }
	RenderPassEncoder   struct{ ref wgpuRenderPassEncoder }
	RenderPipeline      struct{ ref wgpuRenderPipeline }
	Sampler             struct{ ref wgpuSampler }
	ShaderModule        struct{ ref wgpuShaderModule }
	Surface             struct{ ref wgpuSurface }
	Texture             struct{ ref wgpuTexture }
	TextureView         struct{ ref wgpuTextureView }

	SwapChain struct {
		ref    wgpuSwapChain
		device *Device
	}

	Device struct {
		ref     wgpuDevice
		errChan chan *Error
	}
)

func bindGroupLayoutFinalizer(p *BindGroupLayout) {
	wgpuBindGroupLayoutDrop.Call(uintptr(p.ref))
}

func bindGroupFinalizer(p *BindGroup) {
	wgpuBindGroupDrop.Call(uintptr(p.ref))
}

func bufferFinalizer(p *Buffer) {
	wgpuBufferDrop.Call(uintptr(p.ref))
}

func computePipelineFinalizer(p *ComputePipeline) {
	wgpuComputePipelineDrop.Call(uintptr(p.ref))
}

func renderPipelineFinalizer(p *RenderPipeline) {
	wgpuRenderPipelineDrop.Call(uintptr(p.ref))
}

func samplerFinalizer(p *Sampler) {
	wgpuSamplerDrop.Call(uintptr(p.ref))
}

func shaderModuleFinalizer(p *ShaderModule) {
	wgpuShaderModuleDrop.Call(uintptr(p.ref))
}
func textureViewFinalizer(p *TextureView) {
	wgpuTextureViewDrop.Call(uintptr(p.ref))
}

func textureFinalizer(p *Texture) {
	wgpuTextureDrop.Call(uintptr(p.ref))
}

func RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	var opts wgpuRequestAdapterOptions

	if options != nil {
		if options.CompatibleSurface != nil {
			opts.compatibleSurface = options.CompatibleSurface.ref
		}
		opts.powerPreference = options.PowerPreference
		opts.forceFallbackAdapter = options.ForceFallbackAdapter

		if options.AdapterExtras != nil {
			var adapterExtras wgpuAdapterExtras

			adapterExtras.chain.next = nil
			adapterExtras.chain.sType = sType_AdapterExtras
			adapterExtras.backend = options.AdapterExtras.BackendType

			opts.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&adapterExtras))
		}
	}

	var status RequestAdapterStatus
	var adapter *Adapter

	cb := windows.NewCallbackCDecl(func(s RequestAdapterStatus, a wgpuAdapter, _ *byte, _ uintptr) (_ uintptr) {
		status = s
		adapter = &Adapter{ref: a}
		return
	})
	wgpuInstanceRequestAdapter.Call(0, uintptr(unsafe.Pointer(&opts)), cb, 0)

	if status != RequestAdapterStatus_Success {
		return nil, errors.New("failed to request adapter")
	}
	return adapter, nil
}

func CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	var desc wgpuSurfaceDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		if descriptor.WindowsHWND != nil {
			var windowsHWND wgpuSurfaceDescriptorFromWindowsHWND

			windowsHWND.chain.next = nil
			windowsHWND.chain.sType = sType_SurfaceDescriptorFromWindowsHWND
			windowsHWND.hinstance = descriptor.WindowsHWND.Hinstance
			windowsHWND.hwnd = descriptor.WindowsHWND.Hwnd

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&windowsHWND))
		}
	}

	ref, _, _ := wgpuInstanceCreateSurface.Call(0, uintptr(unsafe.Pointer(&desc)))
	if ref == 0 {
		panic("Failed to acquire Surface")
	}
	return &Surface{ref: wgpuSurface(ref)}
}

func (p *Adapter) EnumerateFeatures() []FeatureName {
	size, _, _ := wgpuAdapterEnumerateFeatures.Call(uintptr(p.ref), 0)
	runtime.KeepAlive(p)

	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	wgpuAdapterEnumerateFeatures.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&features[0])))
	runtime.KeepAlive(p)

	return features
}

func (p *Adapter) HasFeature(feature FeatureName) bool {
	hasFeature, _, _ := wgpuAdapterHasFeature.Call(uintptr(p.ref), uintptr(feature))
	runtime.KeepAlive(p)
	return gobool(hasFeature)
}

func (p *Adapter) GetLimits() SupportedLimits {
	var supportedLimits wgpuSupportedLimits

	wgpuAdapterGetLimits.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&supportedLimits)))
	runtime.KeepAlive(p)
	runtime.KeepAlive(supportedLimits)

	limits := supportedLimits.limits
	return SupportedLimits{
		Limits{
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
		},
	}
}

func (p *Adapter) GetProperties() AdapterProperties {
	var props wgpuAdapterProperties

	wgpuAdapterGetProperties.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&props)))
	runtime.KeepAlive(p)

	return AdapterProperties{
		VendorID:          uint32(props.vendorID),
		DeviceID:          uint32(props.deviceID),
		Name:              gostring(props.name),
		DriverDescription: gostring(props.driverDescription),
		AdapterType:       AdapterType(props.adapterType),
		BackendType:       BackendType(props.backendType),
	}
}

func (p *Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {
	var desc wgpuDeviceDescriptor
	desc.requiredLimits = &wgpuRequiredLimits{}

	if descriptor != nil {
		if descriptor.RequiredLimits != nil {
			l := descriptor.RequiredLimits.Limits
			desc.requiredLimits.limits = wgpuLimits{
				maxTextureDimension1D:                     l.MaxTextureDimension1D,
				maxTextureDimension2D:                     l.MaxTextureDimension2D,
				maxTextureDimension3D:                     l.MaxTextureDimension3D,
				maxTextureArrayLayers:                     l.MaxTextureArrayLayers,
				maxBindGroups:                             l.MaxBindGroups,
				maxDynamicUniformBuffersPerPipelineLayout: l.MaxDynamicUniformBuffersPerPipelineLayout,
				maxDynamicStorageBuffersPerPipelineLayout: l.MaxDynamicStorageBuffersPerPipelineLayout,
				maxSampledTexturesPerShaderStage:          l.MaxSampledTexturesPerShaderStage,
				maxSamplersPerShaderStage:                 l.MaxSamplersPerShaderStage,
				maxStorageBuffersPerShaderStage:           l.MaxStorageBuffersPerShaderStage,
				maxStorageTexturesPerShaderStage:          l.MaxStorageTexturesPerShaderStage,
				maxUniformBuffersPerShaderStage:           l.MaxUniformBuffersPerShaderStage,
				maxUniformBufferBindingSize:               l.MaxUniformBufferBindingSize,
				maxStorageBufferBindingSize:               l.MaxStorageBufferBindingSize,
				minUniformBufferOffsetAlignment:           l.MinUniformBufferOffsetAlignment,
				minStorageBufferOffsetAlignment:           l.MinStorageBufferOffsetAlignment,
				maxVertexBuffers:                          l.MaxVertexBuffers,
				maxVertexAttributes:                       l.MaxVertexAttributes,
				maxVertexBufferArrayStride:                l.MaxVertexBufferArrayStride,
				maxInterStageShaderComponents:             l.MaxInterStageShaderComponents,
				maxComputeWorkgroupStorageSize:            l.MaxComputeWorkgroupStorageSize,
				maxComputeInvocationsPerWorkgroup:         l.MaxComputeInvocationsPerWorkgroup,
				maxComputeWorkgroupSizeX:                  l.MaxComputeWorkgroupSizeX,
				maxComputeWorkgroupSizeY:                  l.MaxComputeWorkgroupSizeY,
				maxComputeWorkgroupSizeZ:                  l.MaxComputeWorkgroupSizeZ,
				maxComputeWorkgroupsPerDimension:          l.MaxComputeWorkgroupsPerDimension,
			}

			if descriptor.RequiredLimits.RequiredLimitsExtras != nil {
				var requiredLimitsExtras wgpuRequiredLimitsExtras

				requiredLimitsExtras.chain.next = nil
				requiredLimitsExtras.chain.sType = sType_RequiredLimitsExtras
				requiredLimitsExtras.maxPushConstantSize = descriptor.RequiredLimits.RequiredLimitsExtras.MaxPushConstantSize

				desc.requiredLimits.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&requiredLimitsExtras))
			} else {
				desc.requiredLimits.nextInChain = nil
			}
		} else {
			desc.requiredLimits.limits = wgpuLimits{}
			desc.requiredLimits.nextInChain = nil
		}

		if descriptor.DeviceExtras != nil {
			var deviceExtras wgpuDeviceExtras

			deviceExtras.chain.next = nil
			deviceExtras.chain.sType = sType_DeviceExtras
			deviceExtras.nativeFeatures = descriptor.DeviceExtras.NativeFeatures

			if descriptor.DeviceExtras.Label != "" {
				deviceExtras.label = cstring(descriptor.DeviceExtras.Label)
			} else {
				deviceExtras.label = nil
			}

			if descriptor.DeviceExtras.TracePath != "" {
				deviceExtras.tracePath = cstring(descriptor.DeviceExtras.TracePath)
			} else {
				deviceExtras.tracePath = nil
			}

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&deviceExtras))
		}
	}

	var status RequestDeviceStatus
	var device *Device

	cb := windows.NewCallbackCDecl(func(s RequestDeviceStatus, d wgpuDevice, _ *byte, _ uintptr) (_ uintptr) {
		status = s
		device = &Device{ref: d}
		return
	})

	wgpuAdapterRequestDevice.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)), cb, 0)
	runtime.KeepAlive(p)

	if status != RequestDeviceStatus_Success {
		return nil, errors.New("failed to request device")
	}

	device.errChan = make(chan *Error, 1)
	errCb := windows.NewCallbackCDecl(func(typ ErrorType, msg *byte, _ uintptr) (_ uintptr) {
		device.storeErr(typ, gostring(msg))
		return
	})
	wgpuDeviceSetUncapturedErrorCallback.Call(uintptr(device.ref), errCb)
	runtime.SetFinalizer(device, deviceFinalizer)

	return device, nil
}

func deviceFinalizer(p *Device) {
	wgpuDeviceDrop.Call(uintptr(p.ref))

loop:
	for {
		select {
		case err := <-p.errChan:
			fmt.Printf("go-webgpu: uncaptured error: %s\n", err.Error())
		default:
			break loop
		}
	}

	close(p.errChan)
}

func (p *Device) EnumerateFeatures() []FeatureName {
	size, _, _ := wgpuDeviceEnumerateFeatures.Call(uintptr(p.ref), 0)
	runtime.KeepAlive(p)

	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	wgpuDeviceEnumerateFeatures.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&features[0])))
	runtime.KeepAlive(p)

	return features
}

func (p *Device) HasFeature(feature FeatureName) bool {
	hasFeature, _, _ := wgpuDeviceHasFeature.Call(uintptr(p.ref), uintptr(feature))
	runtime.KeepAlive(p)
	return gobool(hasFeature)
}

func (p *Device) Poll(wait bool, wrappedSubmissionIndex *WrappedSubmissionIndex) (queueEmpty bool) {
	if wrappedSubmissionIndex != nil {
		var index wgpuWrappedSubmissionIndex
		index.queue = wrappedSubmissionIndex.Queue.ref
		index.submissionIndex = wgpuSubmissionIndex(wrappedSubmissionIndex.SubmissionIndex)

		r, _, _ := wgpuDevicePoll.Call(uintptr(p.ref), cbool[uintptr](wait), uintptr(unsafe.Pointer(&index)))
		runtime.KeepAlive(p)
		runtime.KeepAlive(wrappedSubmissionIndex.Queue)
		return gobool(r)
	}

	r, _, _ := wgpuDevicePoll.Call(uintptr(p.ref), cbool[uintptr](wait), 0)
	runtime.KeepAlive(p)
	return gobool(r)
}

func (p *Device) CreateBindGroupLayout(descriptor *BindGroupLayoutDescriptor) (*BindGroupLayout, error) {
	var desc wgpuBindGroupLayoutDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		entryCount := len(descriptor.Entries)
		if entryCount > 0 {
			entries := make([]wgpuBindGroupLayoutEntry, entryCount)

			for i, v := range descriptor.Entries {
				entries[i] = wgpuBindGroupLayoutEntry{
					nextInChain: nil,
					binding:     v.Binding,
					visibility:  v.Visibility,
					buffer: wgpuBufferBindingLayout{
						nextInChain:      nil,
						_type:            v.Buffer.Type,
						hasDynamicOffset: v.Buffer.HasDynamicOffset,
						minBindingSize:   v.Buffer.MinBindingSize,
					},
					sampler: wgpuSamplerBindingLayout{
						nextInChain: nil,
						_type:       v.Sampler.Type,
					},
					texture: wgpuTextureBindingLayout{
						nextInChain:   nil,
						sampleType:    v.Texture.SampleType,
						viewDimension: v.Texture.ViewDimension,
						multisampled:  v.Texture.Multisampled,
					},
					storageTexture: wgpuStorageTextureBindingLayout{
						nextInChain:   nil,
						access:        v.StorageTexture.Access,
						format:        v.StorageTexture.Format,
						viewDimension: v.StorageTexture.ViewDimension,
					},
				}
			}

			desc.entryCount = uint32(entryCount)
			desc.entries = (*wgpuBindGroupLayoutEntry)(unsafe.Pointer(&entries[0]))
		}
	}

	ref, _, _ := wgpuDeviceCreateBindGroupLayout.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire BindGroupLayout")
	}

	layout := &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}
	runtime.SetFinalizer(layout, bindGroupLayoutFinalizer)
	return layout, nil
}

func (p *Device) CreateBindGroup(descriptor *BindGroupDescriptor) (*BindGroup, error) {
	var desc wgpuBindGroupDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		entryCount := len(descriptor.Entries)
		if entryCount > 0 {
			entries := make([]wgpuBindGroupEntry, entryCount)
			for i, v := range descriptor.Entries {
				entry := wgpuBindGroupEntry{
					binding: v.Binding,
					offset:  v.Offset,
					size:    v.Size,
				}

				if v.Buffer != nil {
					entry.buffer = v.Buffer.ref
				}
				if v.Sampler != nil {
					entry.sampler = v.Sampler.ref
				}
				if v.TextureView != nil {
					entry.textureView = v.TextureView.ref
				}

				entries[i] = entry
			}

			desc.entryCount = uint32(entryCount)
			desc.entries = (*wgpuBindGroupEntry)(unsafe.Pointer(&entries[0]))
		}
	}

	ref, _, _ := wgpuDeviceCreateBindGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if descriptor != nil {
		runtime.KeepAlive(descriptor.Layout)

		for _, v := range descriptor.Entries {
			runtime.KeepAlive(v.Buffer)
			runtime.KeepAlive(v.Sampler)
			runtime.KeepAlive(v.TextureView)
		}
	}

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire BindGroup")
	}

	bindGroup := &BindGroup{ref: wgpuBindGroup(ref)}
	runtime.SetFinalizer(bindGroup, bindGroupFinalizer)
	return bindGroup, nil
}

func (p *Device) CreateBuffer(descriptor *BufferDescriptor) (*Buffer, error) {
	var desc wgpuBufferDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		desc.usage = descriptor.Usage
		desc.size = descriptor.Size
		desc.mappedAtCreation = descriptor.MappedAtCreation
	}

	ref, _, _ := wgpuDeviceCreateBuffer.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Buffer")
	}

	buffer := &Buffer{ref: wgpuBuffer(ref)}
	runtime.SetFinalizer(buffer, bufferFinalizer)
	return buffer, nil
}

func (p *Device) CreateCommandEncoder(descriptor *CommandEncoderDescriptor) (*CommandEncoder, error) {
	var desc wgpuCommandEncoderDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuDeviceCreateCommandEncoder.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire CommandEncoder")
	}

	return &CommandEncoder{ref: wgpuCommandEncoder(ref)}, nil
}

func (p *Device) CreateComputePipeline(descriptor *ComputePipelineDescriptor) (*ComputePipeline, error) {
	var desc wgpuComputePipelineDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		var compute wgpuProgrammableStageDescriptor
		if descriptor.Compute.Module != nil {
			compute.module = descriptor.Compute.Module.ref
		}
		if descriptor.Compute.EntryPoint != "" {
			compute.entryPoint = cstring(descriptor.Compute.EntryPoint)
		}
		desc.compute = compute
	}

	ref, _, _ := wgpuDeviceCreateComputePipeline.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if descriptor != nil {
		runtime.KeepAlive(descriptor.Layout)
		runtime.KeepAlive(descriptor.Compute.Module)
	}

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire ComputePipeline")
	}

	pipeline := &ComputePipeline{ref: wgpuComputePipeline(ref)}
	runtime.SetFinalizer(pipeline, computePipelineFinalizer)
	return pipeline, nil
}

func (p *Device) CreatePipelineLayout(descriptor *PipelineLayoutDescriptor) (*PipelineLayout, error) {
	var desc wgpuPipelineLayoutDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		bindGroupLayoutCount := len(descriptor.BindGroupLayouts)
		if bindGroupLayoutCount > 0 {
			bindGroupLayouts := make([]wgpuBindGroupLayout, bindGroupLayoutCount)

			for i, v := range descriptor.BindGroupLayouts {
				bindGroupLayouts[i] = v.ref
			}

			desc.bindGroupLayoutCount = uint32(bindGroupLayoutCount)
			desc.bindGroupLayouts = (*wgpuBindGroupLayout)(unsafe.Pointer(&bindGroupLayouts[0]))
		}

		if descriptor.PipelineLayoutExtras != nil {
			var pipelineLayoutExtras wgpuPipelineLayoutExtras

			pipelineLayoutExtras.chain.next = nil
			pipelineLayoutExtras.chain.sType = sType_PipelineLayoutExtras

			pushConstantRangeCount := len(descriptor.PipelineLayoutExtras.PushConstantRanges)
			if pushConstantRangeCount > 0 {
				pushConstantRanges := make([]wgpuPushConstantRange, pushConstantRangeCount)

				for i, v := range descriptor.PipelineLayoutExtras.PushConstantRanges {
					pushConstantRanges[i] = wgpuPushConstantRange{
						stages: v.Stages,
						start:  v.Start,
						end:    v.End,
					}
				}

				pipelineLayoutExtras.pushConstantRangeCount = uint32(pushConstantRangeCount)
				pipelineLayoutExtras.pushConstantRanges = (*wgpuPushConstantRange)(unsafe.Pointer(&pushConstantRanges[0]))
			}

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&pipelineLayoutExtras))
		} else {
			desc.nextInChain = nil
		}
	}

	ref, _, _ := wgpuDeviceCreatePipelineLayout.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if descriptor != nil {
		for _, v := range descriptor.BindGroupLayouts {
			runtime.KeepAlive(v)
		}
	}

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire PipelineLayout")
	}
	return &PipelineLayout{ref: wgpuPipelineLayout(ref)}, nil
}

func (p *Device) CreateRenderPipeline(descriptor *RenderPipelineDescriptor) (*RenderPipeline, error) {
	var desc wgpuRenderPipelineDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		// vertex
		{
			vertex := descriptor.Vertex

			var vert wgpuVertexState

			if vertex.Module != nil {
				vert.module = vertex.Module.ref
			}

			if vertex.EntryPoint != "" {
				vert.entryPoint = cstring(vertex.EntryPoint)
			}

			bufferCount := len(vertex.Buffers)
			if bufferCount > 0 {
				buffers := make([]wgpuVertexBufferLayout, bufferCount)

				for i, v := range vertex.Buffers {
					buffer := wgpuVertexBufferLayout{
						arrayStride: v.ArrayStride,
						stepMode:    v.StepMode,
					}

					attributeCount := len(v.Attributes)
					if attributeCount > 0 {
						attributes := make([]wgpuVertexAttribute, attributeCount)

						for j, attribute := range v.Attributes {
							attributes[j] = wgpuVertexAttribute{
								format:         attribute.Format,
								offset:         attribute.Offset,
								shaderLocation: attribute.ShaderLocation,
							}
						}

						buffer.attributeCount = uint32(attributeCount)
						buffer.attributes = (*wgpuVertexAttribute)(unsafe.Pointer(&attributes[0]))
					}

					buffers[i] = buffer
				}

				vert.bufferCount = uint32(bufferCount)
				vert.buffers = (*wgpuVertexBufferLayout)(unsafe.Pointer(&buffers[0]))
			}

			desc.vertex = vert
		}

		desc.primitive = wgpuPrimitiveState{
			topology:         descriptor.Primitive.Topology,
			stripIndexFormat: descriptor.Primitive.StripIndexFormat,
			frontFace:        descriptor.Primitive.FrontFace,
			cullMode:         descriptor.Primitive.CullMode,
		}

		if descriptor.DepthStencil != nil {
			depthStencil := descriptor.DepthStencil

			var ds wgpuDepthStencilState

			ds.nextInChain = nil
			ds.format = depthStencil.Format
			ds.depthWriteEnabled = depthStencil.DepthWriteEnabled
			ds.depthCompare = depthStencil.DepthCompare
			ds.stencilFront = wgpuStencilFaceState{
				compare:     depthStencil.StencilFront.Compare,
				failOp:      depthStencil.StencilFront.FailOp,
				depthFailOp: depthStencil.StencilFront.DepthFailOp,
				passOp:      depthStencil.StencilFront.PassOp,
			}
			ds.stencilBack = wgpuStencilFaceState{
				compare:     depthStencil.StencilBack.Compare,
				failOp:      depthStencil.StencilBack.FailOp,
				depthFailOp: depthStencil.StencilBack.DepthFailOp,
				passOp:      depthStencil.StencilBack.PassOp,
			}
			ds.stencilReadMask = depthStencil.StencilReadMask
			ds.stencilWriteMask = depthStencil.StencilWriteMask
			ds.depthBias = depthStencil.DepthBias
			ds.depthBiasSlopeScale = depthStencil.DepthBiasSlopeScale
			ds.depthBiasClamp = depthStencil.DepthBiasClamp

			desc.depthStencil = &ds
		}

		desc.multisample = wgpuMultisampleState{
			count:                  descriptor.Multisample.Count,
			mask:                   descriptor.Multisample.Mask,
			alphaToCoverageEnabled: descriptor.Multisample.AlphaToCoverageEnabled,
		}

		if descriptor.Fragment != nil {
			fragment := descriptor.Fragment

			var frag wgpuFragmentState

			frag.nextInChain = nil
			if fragment.EntryPoint != "" {
				frag.entryPoint = cstring(fragment.EntryPoint)
			}

			if fragment.Module != nil {
				frag.module = fragment.Module.ref
			}

			targetCount := len(fragment.Targets)
			if targetCount > 0 {
				targets := make([]wgpuColorTargetState, targetCount)

				for i, v := range fragment.Targets {
					target := wgpuColorTargetState{
						format:    v.Format,
						writeMask: v.WriteMask,
					}

					if v.Blend != nil {
						var blend wgpuBlendState

						blend.color = wgpuBlendComponent{
							operation: v.Blend.Color.Operation,
							srcFactor: v.Blend.Color.SrcFactor,
							dstFactor: v.Blend.Color.DstFactor,
						}
						blend.alpha = wgpuBlendComponent{
							operation: v.Blend.Alpha.Operation,
							srcFactor: v.Blend.Alpha.SrcFactor,
							dstFactor: v.Blend.Alpha.DstFactor,
						}

						target.blend = &blend
					}

					targets[i] = target
				}

				frag.targetCount = uint32(targetCount)
				frag.targets = (*wgpuColorTargetState)(unsafe.Pointer(&targets[0]))
			} else {
				frag.targetCount = 0
				frag.targets = nil
			}

			desc.fragment = &frag
		}
	}

	ref, _, _ := wgpuDeviceCreateRenderPipeline.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if descriptor != nil {
		runtime.KeepAlive(descriptor.Layout)
		runtime.KeepAlive(descriptor.Vertex.Module)
		if descriptor.Fragment != nil {
			runtime.KeepAlive(descriptor.Fragment.Module)
		}
	}

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire RenderPipeline")
	}

	renderPipeline := &RenderPipeline{ref: wgpuRenderPipeline(ref)}
	runtime.SetFinalizer(renderPipeline, renderPipelineFinalizer)
	return renderPipeline, nil
}

func (p *Device) CreateSampler(descriptor *SamplerDescriptor) (*Sampler, error) {
	var desc wgpuSamplerDescriptor

	if descriptor != nil {
		desc = wgpuSamplerDescriptor{
			addressModeU:  descriptor.AddressModeU,
			addressModeV:  descriptor.AddressModeV,
			addressModeW:  descriptor.AddressModeW,
			magFilter:     descriptor.MagFilter,
			minFilter:     descriptor.MinFilter,
			mipmapFilter:  descriptor.MipmapFilter,
			lodMinClamp:   descriptor.LodMinClamp,
			lodMaxClamp:   descriptor.LodMaxClamp,
			compare:       descriptor.Compare,
			maxAnisotropy: descriptor.MaxAnisotrophy,
		}

		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}
	}

	ref, _, _ := wgpuDeviceCreateSampler.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Sampler")
	}

	sampler := &Sampler{ref: wgpuSampler(ref)}
	runtime.SetFinalizer(sampler, samplerFinalizer)
	return sampler, nil
}

func (p *Device) CreateShaderModule(descriptor *ShaderModuleDescriptor) (*ShaderModule, error) {
	var desc wgpuShaderModuleDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		switch {
		case descriptor.SPIRVDescriptor != nil:
			var spirv wgpuShaderModuleSPIRVDescriptor

			codeSize := len(descriptor.SPIRVDescriptor.Code)
			if codeSize > 0 {
				spirv.codeSize = uint32(codeSize)
				spirv.code = (*uint32)(unsafe.Pointer(&descriptor.SPIRVDescriptor.Code[0]))
			}

			spirv.chain.next = nil
			spirv.chain.sType = sType_ShaderModuleSPIRVDescriptor

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&spirv))

		case descriptor.WGSLDescriptor != nil:
			var wgsl wgpuShaderModuleWGSLDescriptor

			if descriptor.WGSLDescriptor.Code != "" {
				wgsl.code = cstring(descriptor.WGSLDescriptor.Code)
			}
			wgsl.chain.next = nil
			wgsl.chain.sType = sType_ShaderModuleWGSLDescriptor

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&wgsl))

		case descriptor.GLSLDescriptor != nil:
			var glsl wgpuShaderModuleGLSLDescriptor

			if descriptor.GLSLDescriptor.Code != "" {
				glsl.code = cstring(descriptor.GLSLDescriptor.Code)
			}

			defineCount := len(descriptor.GLSLDescriptor.Defines)
			if defineCount > 0 {
				shaderDefinesSlice := make([]wgpuShaderDefine, defineCount)
				index := 0

				for name, value := range descriptor.GLSLDescriptor.Defines {
					shaderDefinesSlice[index] = wgpuShaderDefine{
						name:  cstring(name),
						value: cstring(value),
					}
					index++
				}

				glsl.defineCount = uint32(defineCount)
				glsl.defines = (*wgpuShaderDefine)(unsafe.Pointer(&shaderDefinesSlice[0]))
			}

			glsl.stage = descriptor.GLSLDescriptor.ShaderStage
			glsl.chain.next = nil
			glsl.chain.sType = sType_ShaderModuleGLSLDescriptor

			desc.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&glsl))
		}

	}

	ref, _, _ := wgpuDeviceCreateShaderModule.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire ShaderModule")
	}

	shaderModule := &ShaderModule{ref: wgpuShaderModule(ref)}
	runtime.SetFinalizer(shaderModule, shaderModuleFinalizer)
	return shaderModule, nil
}

func (p *Device) CreateSwapChain(surface *Surface, descriptor *SwapChainDescriptor) (*SwapChain, error) {
	var desc wgpuSwapChainDescriptor

	if descriptor != nil {
		desc = wgpuSwapChainDescriptor{
			usage:       descriptor.Usage,
			format:      descriptor.Format,
			width:       descriptor.Width,
			height:      descriptor.Height,
			presentMode: descriptor.PresentMode,
		}
	}

	ref, _, _ := wgpuDeviceCreateSwapChain.Call(uintptr(p.ref), uintptr(surface.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	runtime.KeepAlive(surface)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire SwapChain")
	}

	return &SwapChain{ref: wgpuSwapChain(ref), device: p}, nil
}

func (p *Device) CreateTexture(descriptor *TextureDescriptor) (*Texture, error) {
	var desc wgpuTextureDescriptor

	if descriptor != nil {
		desc = wgpuTextureDescriptor{
			usage:     descriptor.Usage,
			dimension: descriptor.Dimension,
			size: wgpuExtent3D{
				width:              descriptor.Size.Width,
				height:             descriptor.Size.Height,
				depthOrArrayLayers: descriptor.Size.DepthOrArrayLayers,
			},
			format:        descriptor.Format,
			mipLevelCount: descriptor.MipLevelCount,
			sampleCount:   descriptor.SampleCount,
		}

		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}
	}

	ref, _, _ := wgpuDeviceCreateTexture.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Texture")
	}

	texture := &Texture{ref: wgpuTexture(ref)}
	runtime.SetFinalizer(texture, textureFinalizer)
	return texture, nil
}

func (p *Device) GetLimits() SupportedLimits {
	var supportedLimits wgpuSupportedLimits

	wgpuDeviceGetLimits.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&supportedLimits)))
	runtime.KeepAlive(p)

	limits := supportedLimits.limits
	return SupportedLimits{
		Limits{
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
		},
	}
}

func (p *Device) GetQueue() *Queue {
	ref, _, _ := wgpuDeviceGetQueue.Call(uintptr(p.ref))
	runtime.KeepAlive(p)

	if ref == 0 {
		panic("Failed to acquire Queue")
	}
	return &Queue{ref: wgpuQueue(ref)}
}

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	buf, _, _ := wgpuBufferGetMappedRange.Call(uintptr(p.ref), uintptr(offset), uintptr(size))
	runtime.KeepAlive(p)
	return unsafe.Slice((*byte)(unsafe.Pointer(buf)), size)
}

func (p *Buffer) Unmap() {
	wgpuBufferUnmap.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *Buffer) Destroy() {
	wgpuBufferDestroy.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

type BufferMapCallback func(BufferMapAsyncStatus)

var cbCounter uintptr
var cbStoreMutex sync.RWMutex
var cbStore = map[uintptr]BufferMapCallback{}

var bufferMapCallback = windows.NewCallbackCDecl(func(status BufferMapAsyncStatus, handle uintptr) (_ uintptr) {
	cbStoreMutex.RLock()
	cb, ok := cbStore[handle]
	if ok {
		cbStoreMutex.RUnlock()
		cb(status)

		cbStoreMutex.Lock()
		cbStore[handle] = nil
		delete(cbStore, handle)
		cbStoreMutex.Unlock()
	} else {
		cbStoreMutex.RUnlock()
	}
	return
})

func (p *Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) {
	handle := atomic.AddUintptr(&cbCounter, 1)
	cbStoreMutex.Lock()
	cbStore[handle] = callback
	cbStoreMutex.Unlock()

	wgpuBufferMapAsync.Call(
		uintptr(p.ref),
		uintptr(mode),
		uintptr(offset),
		uintptr(size),
		bufferMapCallback,
		uintptr(handle),
	)
	runtime.KeepAlive(p)
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc wgpuComputePassDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuCommandEncoderBeginComputePass.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if ref == 0 {
		panic("Failed to acquire ComputePassEncoder")
	}

	return &ComputePassEncoder{ref: wgpuComputePassEncoder(ref)}
}

func (p *CommandEncoder) BeginRenderPass(descriptor *RenderPassDescriptor) *RenderPassEncoder {
	var desc wgpuRenderPassDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		colorAttachmentCount := len(descriptor.ColorAttachments)
		if colorAttachmentCount > 0 {
			colorAttachments := make([]wgpuRenderPassColorAttachment, colorAttachmentCount)

			for i, v := range descriptor.ColorAttachments {
				colorAttachment := wgpuRenderPassColorAttachment{
					loadOp:  v.LoadOp,
					storeOp: v.StoreOp,
					clearValue: wgpuColor{
						r: v.ClearValue.R,
						g: v.ClearValue.G,
						b: v.ClearValue.B,
						a: v.ClearValue.A,
					},
				}
				if v.View != nil {
					colorAttachment.view = v.View.ref
				}
				if v.ResolveTarget != nil {
					colorAttachment.resolveTarget = v.ResolveTarget.ref
				}

				colorAttachments[i] = colorAttachment
			}

			desc.colorAttachmentCount = uint32(colorAttachmentCount)
			desc.colorAttachments = (*wgpuRenderPassColorAttachment)(unsafe.Pointer(&colorAttachments[0]))
		}

		if descriptor.DepthStencilAttachment != nil {
			var depthStencilAttachment wgpuRenderPassDepthStencilAttachment

			if descriptor.DepthStencilAttachment.View != nil {
				depthStencilAttachment.view = descriptor.DepthStencilAttachment.View.ref
			}
			depthStencilAttachment.depthLoadOp = descriptor.DepthStencilAttachment.DepthLoadOp
			depthStencilAttachment.depthStoreOp = descriptor.DepthStencilAttachment.DepthStoreOp
			depthStencilAttachment.depthClearValue = descriptor.DepthStencilAttachment.DepthClearValue
			depthStencilAttachment.depthReadOnly = descriptor.DepthStencilAttachment.DepthReadOnly
			depthStencilAttachment.stencilLoadOp = descriptor.DepthStencilAttachment.StencilLoadOp
			depthStencilAttachment.stencilStoreOp = descriptor.DepthStencilAttachment.StencilStoreOp
			depthStencilAttachment.stencilClearValue = descriptor.DepthStencilAttachment.StencilClearValue
			depthStencilAttachment.stencilReadOnly = descriptor.DepthStencilAttachment.DepthReadOnly

			desc.depthStencilAttachment = &depthStencilAttachment
		}
	}

	ref, _, _ := wgpuCommandEncoderBeginRenderPass.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)
	if descriptor != nil {
		for _, v := range descriptor.ColorAttachments {
			runtime.KeepAlive(v.View)
			runtime.KeepAlive(v.ResolveTarget)
		}

		if descriptor.DepthStencilAttachment != nil {
			runtime.KeepAlive(descriptor.DepthStencilAttachment.View)
		}
	}

	if ref == 0 {
		panic("Failed to acquire RenderPassEncoder")
	}

	return &RenderPassEncoder{ref: wgpuRenderPassEncoder(ref)}
}

func (p *CommandEncoder) ClearBuffer(buffer *Buffer, offset uint64, size uint64) {
	wgpuCommandEncoderClearBuffer.Call(
		uintptr(p.ref),
		uintptr(buffer.ref),
		uintptr(offset),
		uintptr(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *CommandEncoder) CopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinatonOffset uint64, size uint64) {
	wgpuCommandEncoderCopyBufferToBuffer.Call(
		uintptr(p.ref),
		uintptr(source.ref),
		uintptr(sourceOffset),
		uintptr(destination.ref),
		uintptr(destinatonOffset),
		uintptr(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(source)
	runtime.KeepAlive(destination)
}

func (p *CommandEncoder) CopyBufferToTexture(source *ImageCopyBuffer, destination *ImageCopyTexture, copySize *Extent3D) {
	var src wgpuImageCopyBuffer
	if source != nil {
		if source.Buffer != nil {
			src.buffer = source.Buffer.ref
		}
		src.layout = wgpuTextureDataLayout{
			offset:       source.Layout.Offset,
			bytesPerRow:  source.Layout.BytesPerRow,
			rowsPerImage: source.Layout.RowsPerImage,
		}
	}

	var dst wgpuImageCopyTexture
	if destination != nil {
		dst = wgpuImageCopyTexture{
			mipLevel: destination.MipLevel,
			origin: wgpuOrigin3D{
				x: destination.Origin.X,
				y: destination.Origin.Y,
				z: destination.Origin.Z,
			},
			aspect: destination.Aspect,
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize wgpuExtent3D
	if copySize != nil {
		cpySize = wgpuExtent3D{
			width:              copySize.Width,
			height:             copySize.Height,
			depthOrArrayLayers: copySize.DepthOrArrayLayers,
		}
	}

	wgpuCommandEncoderCopyBufferToTexture.Call(
		uintptr(p.ref),
		uintptr(unsafe.Pointer(&src)),
		uintptr(unsafe.Pointer(&dst)),
		uintptr(unsafe.Pointer(&cpySize)),
	)

	runtime.KeepAlive(p)
	if source != nil {
		runtime.KeepAlive(source.Buffer)
	}
	if destination != nil {
		runtime.KeepAlive(destination.Texture)
	}
}

func (p *CommandEncoder) CopyTextureToBuffer(source *ImageCopyTexture, destination *ImageCopyBuffer, copySize *Extent3D) {
	var src wgpuImageCopyTexture
	if source != nil {
		src = wgpuImageCopyTexture{
			mipLevel: source.MipLevel,
			origin: wgpuOrigin3D{
				x: source.Origin.X,
				y: source.Origin.Y,
				z: source.Origin.Z,
			},
			aspect: source.Aspect,
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst wgpuImageCopyBuffer
	if destination != nil {
		if destination.Buffer != nil {
			dst.buffer = destination.Buffer.ref
		}
		dst.layout = wgpuTextureDataLayout{
			offset:       destination.Layout.Offset,
			bytesPerRow:  destination.Layout.BytesPerRow,
			rowsPerImage: destination.Layout.RowsPerImage,
		}
	}

	var cpySize wgpuExtent3D
	if copySize != nil {
		cpySize = wgpuExtent3D{
			width:              copySize.Width,
			height:             copySize.Height,
			depthOrArrayLayers: copySize.DepthOrArrayLayers,
		}
	}

	wgpuCommandEncoderCopyTextureToBuffer.Call(
		uintptr(p.ref),
		uintptr(unsafe.Pointer(&src)),
		uintptr(unsafe.Pointer(&dst)),
		uintptr(unsafe.Pointer(&cpySize)),
	)

	runtime.KeepAlive(p)
	if source != nil {
		runtime.KeepAlive(source.Texture)
	}
	if destination != nil {
		runtime.KeepAlive(destination.Buffer)
	}
}

func (p *CommandEncoder) CopyTextureToTexture(source *ImageCopyTexture, destination *ImageCopyTexture, copySize *Extent3D) {
	var src wgpuImageCopyTexture
	if source != nil {
		src = wgpuImageCopyTexture{
			mipLevel: source.MipLevel,
			origin: wgpuOrigin3D{
				x: source.Origin.X,
				y: source.Origin.Y,
				z: source.Origin.Z,
			},
			aspect: source.Aspect,
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst wgpuImageCopyTexture
	if destination != nil {
		dst = wgpuImageCopyTexture{
			mipLevel: destination.MipLevel,
			origin: wgpuOrigin3D{
				x: destination.Origin.X,
				y: destination.Origin.Y,
				z: destination.Origin.Z,
			},
			aspect: destination.Aspect,
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize wgpuExtent3D
	if copySize != nil {
		cpySize = wgpuExtent3D{
			width:              copySize.Width,
			height:             copySize.Height,
			depthOrArrayLayers: copySize.DepthOrArrayLayers,
		}
	}

	wgpuCommandEncoderCopyTextureToTexture.Call(
		uintptr(p.ref),
		uintptr(unsafe.Pointer(&src)),
		uintptr(unsafe.Pointer(&dst)),
		uintptr(unsafe.Pointer(&cpySize)),
	)
	runtime.KeepAlive(p)
	if source != nil {
		runtime.KeepAlive(source.Texture)
	}
	if destination != nil {
		runtime.KeepAlive(destination.Texture)
	}
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) *CommandBuffer {
	var desc wgpuCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuCommandEncoderFinish.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	if ref == 0 {
		panic("Failed to acquire CommandBuffer")
	}

	return &CommandBuffer{ref: wgpuCommandBuffer(ref)}
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) {
	wgpuCommandEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
	runtime.KeepAlive(p)
}

func (p *CommandEncoder) PopDebugGroup() {
	wgpuCommandEncoderPopDebugGroup.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) {
	wgpuCommandEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
	runtime.KeepAlive(p)
}

func (p *ComputePassEncoder) DispatchWorkgroups(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	wgpuComputePassEncoderDispatchWorkgroups.Call(
		uintptr(p.ref),
		uintptr(workgroupCountX),
		uintptr(workgroupCountY),
		uintptr(workgroupCountZ),
	)
	runtime.KeepAlive(p)
}

func (p *ComputePassEncoder) DispatchWorkgroupsIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuComputePassEncoderDispatchWorkgroupsIndirect.Call(
		uintptr(p.ref),
		uintptr(indirectBuffer.ref),
		uintptr(indirectOffset),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(indirectBuffer)
}

func (p *ComputePassEncoder) End() {
	wgpuComputePassEncoderEnd.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *ComputePassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		wgpuComputePassEncoderSetBindGroup.Call(uintptr(p.ref), uintptr(groupIndex), uintptr(group.ref), 0, 0)
	} else {
		wgpuComputePassEncoderSetBindGroup.Call(
			uintptr(p.ref), uintptr(groupIndex), uintptr(group.ref),
			uintptr(dynamicOffsetCount), uintptr((unsafe.Pointer(&dynamicOffsets[0]))),
		)
	}

	runtime.KeepAlive(p)
	runtime.KeepAlive(group)
}

func (p *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	wgpuComputePassEncoderSetPipeline.Call(uintptr(p.ref), uintptr(pipeline.ref))
	runtime.KeepAlive(p)
	runtime.KeepAlive(pipeline)
}

func (p *ComputePassEncoder) InsertDebugMarker(markerLabel string) {
	wgpuComputePassEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
	runtime.KeepAlive(p)
}

func (p *ComputePassEncoder) PopDebugGroup() {
	wgpuComputePassEncoderPopDebugGroup.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *ComputePassEncoder) PushDebugGroup(groupLabel string) {
	wgpuComputePassEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
	runtime.KeepAlive(p)
}

func (p *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref, _, _ := wgpuComputePipelineGetBindGroupLayout.Call(uintptr(p.ref), uintptr(groupIndex))
	runtime.KeepAlive(p)

	if ref == 0 {
		panic("Failed to accquire BindGroupLayout")
	}

	bindGroupLayout := &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}
	runtime.SetFinalizer(bindGroupLayout, bindGroupLayoutFinalizer)
	return bindGroupLayout
}

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r, _, _ := wgpuQueueSubmitForIndex.Call(uintptr(p.ref), 0, 0)
		runtime.KeepAlive(p)
		return SubmissionIndex(r)
	}

	commandRefs := make([]wgpuCommandBuffer, commandCount)
	for i, v := range commands {
		commandRefs[i] = v.ref
	}

	r, _, _ := wgpuQueueSubmitForIndex.Call(
		uintptr(p.ref),
		uintptr(commandCount),
		uintptr(unsafe.Pointer(&commandRefs[0])),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(commands)
	for _, v := range commands {
		runtime.KeepAlive(v)
	}
	return SubmissionIndex(r)
}

func (p *Queue) WriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) {
	size := len(data)
	if size == 0 {
		wgpuQueueWriteBuffer.Call(
			uintptr(p.ref),
			uintptr(buffer.ref),
			uintptr(bufferOffset),
			0,
			0,
		)
	} else {
		wgpuQueueWriteBuffer.Call(
			uintptr(p.ref),
			uintptr(buffer.ref),
			uintptr(bufferOffset),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(size),
		)
	}

	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) {
	size := len(data)

	var dst wgpuImageCopyTexture
	if destination != nil {
		dst = wgpuImageCopyTexture{
			mipLevel: destination.MipLevel,
			origin: wgpuOrigin3D{
				x: destination.Origin.X,
				y: destination.Origin.Y,
				z: destination.Origin.Z,
			},
			aspect: destination.Aspect,
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var layout wgpuTextureDataLayout
	if dataLayout != nil {
		layout = wgpuTextureDataLayout{
			offset:       dataLayout.Offset,
			bytesPerRow:  dataLayout.BytesPerRow,
			rowsPerImage: dataLayout.RowsPerImage,
		}
	}

	var writeExtent wgpuExtent3D
	if writeSize != nil {
		writeExtent = wgpuExtent3D{
			width:              writeSize.Width,
			height:             writeSize.Height,
			depthOrArrayLayers: writeSize.DepthOrArrayLayers,
		}
	}

	if size == 0 {
		wgpuQueueWriteTexture.Call(
			uintptr(p.ref),
			uintptr(unsafe.Pointer(&dst)),
			0,
			0,
			uintptr(unsafe.Pointer(&layout)),
			uintptr(unsafe.Pointer(&writeExtent)),
		)
	} else {
		wgpuQueueWriteTexture.Call(
			uintptr(p.ref),
			uintptr(unsafe.Pointer(&dst)),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(size),
			uintptr(unsafe.Pointer(&layout)),
			uintptr(unsafe.Pointer(&writeExtent)),
		)
	}

	runtime.KeepAlive(p)
	if destination != nil {
		runtime.KeepAlive(destination.Texture)
	}
}

func (p *RenderPassEncoder) SetPushConstants(stages ShaderStage, offset uint32, data []byte) {
	size := len(data)

	wgpuRenderPassEncoderSetPushConstants.Call(
		uintptr(p.ref),
		uintptr(stages),
		uintptr(offset),
		uintptr(size),
		uintptr(unsafe.Pointer(&data[0])),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	wgpuRenderPassEncoderDraw.Call(
		uintptr(p.ref),
		uintptr(vertexCount),
		uintptr(instanceCount),
		uintptr(firstVertex),
		uintptr(firstInstance),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	wgpuRenderPassEncoderDrawIndexed.Call(
		uintptr(p.ref),
		uintptr(indexCount),
		uintptr(instanceCount),
		uintptr(firstIndex),
		uintptr(baseVertex),
		uintptr(firstInstance),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderPassEncoderDrawIndexedIndirect.Call(uintptr(p.ref), uintptr(indirectBuffer.ref), uintptr(indirectOffset))
	runtime.KeepAlive(p)
	runtime.KeepAlive(indirectBuffer)
}

func (p *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderPassEncoderDrawIndirect.Call(uintptr(p.ref), uintptr(indirectBuffer.ref), uintptr(indirectOffset))
	runtime.KeepAlive(p)
	runtime.KeepAlive(indirectBuffer)
}

func (p *RenderPassEncoder) End() {
	wgpuRenderPassEncoderEnd.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		wgpuRenderPassEncoderSetBindGroup.Call(
			uintptr(p.ref),
			uintptr(groupIndex),
			uintptr(group.ref),
			0,
			0,
		)
	} else {
		wgpuRenderPassEncoderSetBindGroup.Call(
			uintptr(p.ref),
			uintptr(groupIndex),
			uintptr(group.ref),
			uintptr(dynamicOffsetCount),
			uintptr(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
	runtime.KeepAlive(p)
	runtime.KeepAlive(group)
}

func (p *RenderPassEncoder) SetBlendConstant(color *Color) {
	wgpuRenderPassEncoderSetBlendConstant.Call(
		uintptr(p.ref),
		uintptr(unsafe.Pointer(&wgpuColor{
			r: color.R,
			g: color.G,
			b: color.B,
			a: color.A,
		})),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	wgpuRenderPassEncoderSetIndexBuffer.Call(
		uintptr(p.ref),
		uintptr(buffer.ref),
		uintptr(format),
		uintptr(offset),
		uintptr(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	wgpuRenderPassEncoderSetPipeline.Call(uintptr(p.ref), uintptr(pipeline.ref))
	runtime.KeepAlive(p)
	runtime.KeepAlive(pipeline)
}

func (p *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	wgpuRenderPassEncoderSetScissorRect.Call(
		uintptr(p.ref),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetStencilReference(reference uint32) {
	wgpuRenderPassEncoderSetStencilReference.Call(uintptr(p.ref), uintptr(reference))
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	wgpuRenderPassEncoderSetVertexBuffer.Call(
		uintptr(p.ref),
		uintptr(slot),
		uintptr(buffer.ref),
		uintptr(offset),
		uintptr(size),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buffer)
}

func (p *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {

	wgpuRenderPassEncoderSetViewport.Call(
		uintptr(p.ref),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(minDepth)),
		uintptr(math.Float32bits(maxDepth)),
	)
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) InsertDebugMarker(markerLabel string) {
	wgpuRenderPassEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) PopDebugGroup() {
	wgpuRenderPassEncoderPopDebugGroup.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *RenderPassEncoder) PushDebugGroup(groupLabel string) {
	wgpuRenderPassEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
	runtime.KeepAlive(p)
}

func (p *RenderPipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref, _, _ := wgpuRenderPipelineGetBindGroupLayout.Call(uintptr(p.ref), uintptr(groupIndex))
	runtime.KeepAlive(p)

	if ref == 0 {
		panic("Failed to accquire BindGroupLayout")
	}

	bindGroupLayout := &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}
	runtime.SetFinalizer(bindGroupLayout, bindGroupLayoutFinalizer)
	return bindGroupLayout
}

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format, _, _ := wgpuSurfaceGetPreferredFormat.Call(uintptr(p.ref), uintptr(adapter.ref))
	runtime.KeepAlive(p)
	runtime.KeepAlive(adapter)
	return TextureFormat(format)
}

func (p *Surface) GetSupportedFormats(adapter *Adapter) []TextureFormat {
	var count uintptr

	formatsPtr, _, _ := wgpuSurfaceGetSupportedFormats.Call(
		uintptr(p.ref),
		uintptr(adapter.ref),
		uintptr(unsafe.Pointer(&count)),
	)
	runtime.KeepAlive(p)
	runtime.KeepAlive(adapter)
	defer free[TextureFormat](formatsPtr, count)

	formatsSlice := unsafe.Slice((*TextureFormat)(unsafe.Pointer(formatsPtr)), count)
	formats := make([]TextureFormat, count)
	copy(formats, formatsSlice)
	return formats
}

func (p *SwapChain) GetCurrentTextureView() (*TextureView, error) {
	ref, _, _ := wgpuSwapChainGetCurrentTextureView.Call(uintptr(p.ref))
	runtime.KeepAlive(p)

	err := p.device.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire TextureView")
	}

	textureView := &TextureView{ref: wgpuTextureView(ref)}
	runtime.SetFinalizer(textureView, textureViewFinalizer)
	return textureView, nil
}

func (p *SwapChain) Present() {
	wgpuSwapChainPresent.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}

func (p *Texture) CreateView(descriptor *TextureViewDescriptor) *TextureView {
	var desc wgpuTextureViewDescriptor

	if descriptor != nil {
		desc = wgpuTextureViewDescriptor{
			format:          descriptor.Format,
			dimension:       descriptor.Dimension,
			baseMipLevel:    descriptor.BaseMipLevel,
			mipLevelCount:   descriptor.MipLevelCount,
			baseArrayLayer:  descriptor.BaseArrayLayer,
			arrayLayerCount: descriptor.ArrayLayerCount,
			aspect:          descriptor.Aspect,
		}

		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}
	}

	ref, _, _ := wgpuTextureCreateView.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	runtime.KeepAlive(p)

	if ref == 0 {
		panic("Failed to acquire TextureView")
	}

	textureView := &TextureView{ref: wgpuTextureView(ref)}
	runtime.SetFinalizer(textureView, textureViewFinalizer)
	return textureView
}

func (p *Texture) Destroy() {
	wgpuTextureDestroy.Call(uintptr(p.ref))
	runtime.KeepAlive(p)
}
