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

	wgpuAdapterDrop         = lib.NewProc("wgpuAdapterDrop")
	wgpuBindGroupDrop       = lib.NewProc("wgpuBindGroupDrop")
	wgpuBindGroupLayoutDrop = lib.NewProc("wgpuBindGroupLayoutDrop")
	wgpuBufferDrop          = lib.NewProc("wgpuBufferDrop")
	wgpuCommandBufferDrop   = lib.NewProc("wgpuCommandBufferDrop")
	wgpuCommandEncoderDrop  = lib.NewProc("wgpuCommandEncoderDrop")
	wgpuComputePipelineDrop = lib.NewProc("wgpuComputePipelineDrop")
	wgpuDeviceDrop          = lib.NewProc("wgpuDeviceDrop")
	wgpuPipelineLayoutDrop  = lib.NewProc("wgpuPipelineLayoutDrop")
	wgpuQuerySetDrop        = lib.NewProc("wgpuQuerySetDrop")
	wgpuRenderBundleDrop    = lib.NewProc("wgpuRenderBundleDrop")
	wgpuRenderPipelineDrop  = lib.NewProc("wgpuRenderPipelineDrop")
	wgpuSamplerDrop         = lib.NewProc("wgpuSamplerDrop")
	wgpuShaderModuleDrop    = lib.NewProc("wgpuShaderModuleDrop")
	wgpuSurfaceDrop         = lib.NewProc("wgpuSurfaceDrop")
	wgpuTextureDrop         = lib.NewProc("wgpuTextureDrop")
	wgpuTextureViewDrop     = lib.NewProc("wgpuTextureViewDrop")

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
	wgpuDeviceCreateRenderBundleEncoder              = lib.NewProc("wgpuDeviceCreateRenderBundleEncoder")
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
	wgpuRenderPassEncoderExecuteBundles              = lib.NewProc("wgpuRenderPassEncoderExecuteBundles")
	wgpuRenderPassEncoderInsertDebugMarker           = lib.NewProc("wgpuRenderPassEncoderInsertDebugMarker")
	wgpuRenderPassEncoderPopDebugGroup               = lib.NewProc("wgpuRenderPassEncoderPopDebugGroup")
	wgpuRenderPassEncoderPushDebugGroup              = lib.NewProc("wgpuRenderPassEncoderPushDebugGroup")
	wgpuRenderPipelineGetBindGroupLayout             = lib.NewProc("wgpuRenderPipelineGetBindGroupLayout")
	wgpuSurfaceGetPreferredFormat                    = lib.NewProc("wgpuSurfaceGetPreferredFormat")
	wgpuSurfaceGetSupportedFormats                   = lib.NewProc("wgpuSurfaceGetSupportedFormats")
	wgpuSurfaceGetSupportedPresentModes              = lib.NewProc("wgpuSurfaceGetSupportedPresentModes")
	wgpuSwapChainGetCurrentTextureView               = lib.NewProc("wgpuSwapChainGetCurrentTextureView")
	wgpuSwapChainPresent                             = lib.NewProc("wgpuSwapChainPresent")
	wgpuTextureCreateView                            = lib.NewProc("wgpuTextureCreateView")
	wgpuTextureDestroy                               = lib.NewProc("wgpuTextureDestroy")
	wgpuRenderBundleEncoderDraw                      = lib.NewProc("wgpuRenderBundleEncoderDraw")
	wgpuRenderBundleEncoderDrawIndexed               = lib.NewProc("wgpuRenderBundleEncoderDrawIndexed")
	wgpuRenderBundleEncoderDrawIndexedIndirect       = lib.NewProc("wgpuRenderBundleEncoderDrawIndexedIndirect")
	wgpuRenderBundleEncoderDrawIndirect              = lib.NewProc("wgpuRenderBundleEncoderDrawIndirect")
	wgpuRenderBundleEncoderFinish                    = lib.NewProc("wgpuRenderBundleEncoderFinish")
	wgpuRenderBundleEncoderInsertDebugMarker         = lib.NewProc("wgpuRenderBundleEncoderInsertDebugMarker")
	wgpuRenderBundleEncoderPopDebugGroup             = lib.NewProc("wgpuRenderBundleEncoderPopDebugGroup")
	wgpuRenderBundleEncoderPushDebugGroup            = lib.NewProc("wgpuRenderBundleEncoderPushDebugGroup")
	wgpuRenderBundleEncoderSetBindGroup              = lib.NewProc("wgpuRenderBundleEncoderSetBindGroup")
	wgpuRenderBundleEncoderSetIndexBuffer            = lib.NewProc("wgpuRenderBundleEncoderSetIndexBuffer")
	wgpuRenderBundleEncoderSetPipeline               = lib.NewProc("wgpuRenderBundleEncoderSetPipeline")
	wgpuRenderBundleEncoderSetVertexBuffer           = lib.NewProc("wgpuRenderBundleEncoderSetVertexBuffer")
)

var logCallback = windows.NewCallbackCDecl(func(level LogLevel, msg *byte) (_ uintptr) {
	var l string
	switch level {
	case LogLevel_Error:
		l = "Error"
	case LogLevel_Warn:
		l = "Warn"
	case LogLevel_Info:
		l = "Info"
	case LogLevel_Debug:
		l = "Debug"
	case LogLevel_Trace:
		l = "Trace"
	default:
		l = "Unknown Level"
	}

	fmt.Fprintf(os.Stderr, "[wgpu] [%s] %s\n", l, gostring(msg))
	return
})

func init() {
	wgpuSetLogCallback.Call(logCallback)
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
	RenderBundle        struct{ ref wgpuRenderBundle }
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
	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	wgpuAdapterEnumerateFeatures.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&features[0])))
	return features
}

func (p *Adapter) HasFeature(feature FeatureName) bool {
	hasFeature, _, _ := wgpuAdapterHasFeature.Call(uintptr(p.ref), uintptr(feature))
	return gobool(hasFeature)
}

func (p *Adapter) GetLimits() SupportedLimits {
	var supportedLimits wgpuSupportedLimits

	var extras wgpuSupportedLimitsExtras
	supportedLimits.nextInChain = (*wgpuChainedStructOut)(unsafe.Pointer(&extras))

	wgpuAdapterGetLimits.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&supportedLimits)))

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
			MaxPushConstantSize:                       uint32(extras.maxPushConstantSize),
			MaxBufferSize:                             uint64(extras.maxBufferSize),
		},
	}
}

func (p *Adapter) GetProperties() AdapterProperties {
	var props wgpuAdapterProperties

	wgpuAdapterGetProperties.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&props)))

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

			var requiredLimitsExtras wgpuRequiredLimitsExtras

			requiredLimitsExtras.chain.next = nil
			requiredLimitsExtras.chain.sType = sType_RequiredLimitsExtras
			requiredLimitsExtras.maxPushConstantSize = l.MaxPushConstantSize
			requiredLimitsExtras.maxBufferSize = l.MaxBufferSize

			desc.requiredLimits.nextInChain = (*wgpuChainedStruct)(unsafe.Pointer(&requiredLimitsExtras))
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
	if status != RequestDeviceStatus_Success {
		return nil, errors.New("failed to request device")
	}

	device.errChan = make(chan *Error, 1)
	errCb := windows.NewCallbackCDecl(func(typ ErrorType, msg *byte, _ uintptr) (_ uintptr) {
		device.storeErr(typ, gostring(msg))
		return
	})
	wgpuDeviceSetUncapturedErrorCallback.Call(uintptr(device.ref), errCb)

	return device, nil
}

func (p *Device) Drop() {
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
	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	wgpuDeviceEnumerateFeatures.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&features[0])))
	return features
}

func (p *Device) HasFeature(feature FeatureName) bool {
	hasFeature, _, _ := wgpuDeviceHasFeature.Call(uintptr(p.ref), uintptr(feature))
	return gobool(hasFeature)
}

func (p *Device) Poll(wait bool, wrappedSubmissionIndex *WrappedSubmissionIndex) (queueEmpty bool) {
	if wrappedSubmissionIndex != nil {
		var index wgpuWrappedSubmissionIndex
		index.queue = wrappedSubmissionIndex.Queue.ref
		index.submissionIndex = wgpuSubmissionIndex(wrappedSubmissionIndex.SubmissionIndex)

		r, _, _ := wgpuDevicePoll.Call(uintptr(p.ref), cbool(wait), uintptr(unsafe.Pointer(&index)))
		return gobool(r)
	}

	r, _, _ := wgpuDevicePoll.Call(uintptr(p.ref), cbool(wait), 0)
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire BindGroupLayout")
	}

	return &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire BindGroup")
	}

	return &BindGroup{ref: wgpuBindGroup(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Buffer")
	}

	return &Buffer{ref: wgpuBuffer(ref)}, nil
}

func (p *Device) CreateCommandEncoder(descriptor *CommandEncoderDescriptor) (*CommandEncoder, error) {
	var desc wgpuCommandEncoderDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuDeviceCreateCommandEncoder.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire ComputePipeline")
	}

	return &ComputePipeline{ref: wgpuComputePipeline(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire RenderPipeline")
	}

	return &RenderPipeline{ref: wgpuRenderPipeline(ref)}, nil
}

func (p *Device) CreateRenderBundleEncoder(descriptor *RenderBundleEncoderDescriptor) (*RenderBundleEncoder, error) {
	var desc wgpuRenderBundleEncoderDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			desc.label = cstring(descriptor.Label)
		}

		colorFormatsCount := len(descriptor.ColorFormats)
		if colorFormatsCount > 0 {
			desc.colorFormatsCount = uint32(colorFormatsCount)
			desc.colorFormats = (*TextureFormat)(unsafe.Pointer(&descriptor.ColorFormats[0]))
		}

		desc.depthStencilFormat = descriptor.DepthStencilFormat
		desc.sampleCount = descriptor.SampleCount
		desc.depthReadOnly = descriptor.DepthReadOnly
		desc.stencilReadOnly = descriptor.StencilReadOnly
	}

	ref, _, _ := wgpuDeviceCreateRenderBundleEncoder.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire RenderPipeline")
	}

	return &RenderBundleEncoder{wgpuRenderBundleEncoder(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Sampler")
	}

	return &Sampler{ref: wgpuSampler(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire ShaderModule")
	}

	return &ShaderModule{ref: wgpuShaderModule(ref)}, nil
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
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire Texture")
	}

	return &Texture{ref: wgpuTexture(ref)}, nil
}

func (p *Device) GetLimits() SupportedLimits {
	var supportedLimits wgpuSupportedLimits

	var extras wgpuSupportedLimitsExtras
	supportedLimits.nextInChain = (*wgpuChainedStructOut)(unsafe.Pointer(&extras))

	wgpuDeviceGetLimits.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&supportedLimits)))

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
			MaxPushConstantSize:                       uint32(extras.maxPushConstantSize),
			MaxBufferSize:                             uint64(extras.maxBufferSize),
		},
	}
}

func (p *Device) GetQueue() *Queue {
	ref, _, _ := wgpuDeviceGetQueue.Call(uintptr(p.ref))
	if ref == 0 {
		panic("Failed to acquire Queue")
	}

	return &Queue{ref: wgpuQueue(ref)}
}

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	buf, _, _ := wgpuBufferGetMappedRange.Call(uintptr(p.ref), uintptr(offset), uintptr(size))
	return unsafe.Slice((*byte)(unsafe.Pointer(buf)), size)
}

func (p *Buffer) Unmap()   { wgpuBufferUnmap.Call(uintptr(p.ref)) }
func (p *Buffer) Destroy() { wgpuBufferDestroy.Call(uintptr(p.ref)) }

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
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc wgpuComputePassDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuCommandEncoderBeginComputePass.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
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
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) *CommandBuffer {
	var desc wgpuCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuCommandEncoderFinish.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	if ref == 0 {
		panic("Failed to acquire CommandBuffer")
	}

	return &CommandBuffer{ref: wgpuCommandBuffer(ref)}
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) {
	wgpuCommandEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
}

func (p *CommandEncoder) PopDebugGroup() {
	wgpuCommandEncoderPopDebugGroup.Call(uintptr(p.ref))
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) {
	wgpuCommandEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
}

func (p *ComputePassEncoder) DispatchWorkgroups(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	wgpuComputePassEncoderDispatchWorkgroups.Call(
		uintptr(p.ref),
		uintptr(workgroupCountX),
		uintptr(workgroupCountY),
		uintptr(workgroupCountZ),
	)
}

func (p *ComputePassEncoder) DispatchWorkgroupsIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuComputePassEncoderDispatchWorkgroupsIndirect.Call(
		uintptr(p.ref),
		uintptr(indirectBuffer.ref),
		uintptr(indirectOffset),
	)
}

func (p *ComputePassEncoder) End() {
	wgpuComputePassEncoderEnd.Call(uintptr(p.ref))
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
}

func (p *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	wgpuComputePassEncoderSetPipeline.Call(uintptr(p.ref), uintptr(pipeline.ref))
}

func (p *ComputePassEncoder) InsertDebugMarker(markerLabel string) {
	wgpuComputePassEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
}

func (p *ComputePassEncoder) PopDebugGroup() {
	wgpuComputePassEncoderPopDebugGroup.Call(uintptr(p.ref))
}

func (p *ComputePassEncoder) PushDebugGroup(groupLabel string) {
	wgpuComputePassEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
}

func (p *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref, _, _ := wgpuComputePipelineGetBindGroupLayout.Call(uintptr(p.ref), uintptr(groupIndex))
	if ref == 0 {
		panic("Failed to accquire BindGroupLayout")
	}

	return &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}
}

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r, _, _ := wgpuQueueSubmitForIndex.Call(uintptr(p.ref), 0, 0)
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
}

func (p *RenderPassEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	wgpuRenderPassEncoderDraw.Call(
		uintptr(p.ref),
		uintptr(vertexCount),
		uintptr(instanceCount),
		uintptr(firstVertex),
		uintptr(firstInstance),
	)
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
}

func (p *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderPassEncoderDrawIndexedIndirect.Call(uintptr(p.ref), uintptr(indirectBuffer.ref), uintptr(indirectOffset))
}

func (p *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderPassEncoderDrawIndirect.Call(uintptr(p.ref), uintptr(indirectBuffer.ref), uintptr(indirectOffset))
}

func (p *RenderPassEncoder) End() {
	wgpuRenderPassEncoderEnd.Call(uintptr(p.ref))
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
}

func (p *RenderPassEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	wgpuRenderPassEncoderSetIndexBuffer.Call(
		uintptr(p.ref),
		uintptr(buffer.ref),
		uintptr(format),
		uintptr(offset),
		uintptr(size),
	)
}

func (p *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	wgpuRenderPassEncoderSetPipeline.Call(uintptr(p.ref), uintptr(pipeline.ref))
}

func (p *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	wgpuRenderPassEncoderSetScissorRect.Call(
		uintptr(p.ref),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
	)
}

func (p *RenderPassEncoder) SetStencilReference(reference uint32) {
	wgpuRenderPassEncoderSetStencilReference.Call(uintptr(p.ref), uintptr(reference))
}

func (p *RenderPassEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	wgpuRenderPassEncoderSetVertexBuffer.Call(
		uintptr(p.ref),
		uintptr(slot),
		uintptr(buffer.ref),
		uintptr(offset),
		uintptr(size),
	)
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
}

func (p *RenderPassEncoder) ExecuteBundles(bundles ...*RenderBundle) {
	bundlesCount := len(bundles)
	if bundlesCount == 0 {
		wgpuRenderPassEncoderExecuteBundles.Call(uintptr(p.ref), 0, 0)
		return
	}

	bundlesSlice := make([]wgpuRenderBundle, bundlesCount)
	for i, v := range bundles {
		bundlesSlice[i] = v.ref
	}

	wgpuRenderPassEncoderExecuteBundles.Call(uintptr(p.ref), uintptr(bundlesCount), uintptr(unsafe.Pointer(&bundlesSlice[0])))
}

func (p *RenderPassEncoder) InsertDebugMarker(markerLabel string) {
	wgpuRenderPassEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(markerLabel))))
}

func (p *RenderPassEncoder) PopDebugGroup() {
	wgpuRenderPassEncoderPopDebugGroup.Call(uintptr(p.ref))
}

func (p *RenderPassEncoder) PushDebugGroup(groupLabel string) {
	wgpuRenderPassEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(cstring(groupLabel))))
}

func (p *RenderPipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref, _, _ := wgpuRenderPipelineGetBindGroupLayout.Call(uintptr(p.ref), uintptr(groupIndex))
	if ref == 0 {
		panic("Failed to accquire BindGroupLayout")
	}

	return &BindGroupLayout{ref: wgpuBindGroupLayout(ref)}
}

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format, _, _ := wgpuSurfaceGetPreferredFormat.Call(uintptr(p.ref), uintptr(adapter.ref))
	return TextureFormat(format)
}

func (p *Surface) GetSupportedFormats(adapter *Adapter) []TextureFormat {
	var count uintptr

	formatsPtr, _, _ := wgpuSurfaceGetSupportedFormats.Call(
		uintptr(p.ref),
		uintptr(adapter.ref),
		uintptr(unsafe.Pointer(&count)),
	)
	defer free[TextureFormat](formatsPtr, count)

	formatsSlice := unsafe.Slice((*TextureFormat)(unsafe.Pointer(formatsPtr)), count)
	formats := make([]TextureFormat, count)
	copy(formats, formatsSlice)
	return formats
}

func (p *Surface) GetSupportedPresentModes(adapter *Adapter) []PresentMode {
	var size uintptr
	modesPtr, _, _ := wgpuSurfaceGetSupportedPresentModes.Call(uintptr(p.ref), uintptr(adapter.ref), uintptr(unsafe.Pointer(&size)))
	defer free[PresentMode](modesPtr, size)

	modesSlice := unsafe.Slice((*PresentMode)(unsafe.Pointer(modesPtr)), size)
	modes := make([]PresentMode, size)
	copy(modes, modesSlice)
	return modes
}

func (p *SwapChain) GetCurrentTextureView() (*TextureView, error) {
	ref, _, _ := wgpuSwapChainGetCurrentTextureView.Call(uintptr(p.ref))
	err := p.device.getErr()
	if err != nil {
		return nil, err
	}
	if ref == 0 {
		panic("Failed to acquire TextureView")
	}

	return &TextureView{ref: wgpuTextureView(ref)}, nil
}

func (p *SwapChain) Present() {
	wgpuSwapChainPresent.Call(uintptr(p.ref))
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
	if ref == 0 {
		panic("Failed to acquire TextureView")
	}

	return &TextureView{ref: wgpuTextureView(ref)}
}

func (p *Texture) Destroy() {
	wgpuTextureDestroy.Call(uintptr(p.ref))
}

func (p *RenderBundleEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	wgpuRenderBundleEncoderDraw.Call(
		uintptr(p.ref),
		uintptr(vertexCount),
		uintptr(instanceCount),
		uintptr(firstVertex),
		uintptr(firstInstance),
	)
}

func (p *RenderBundleEncoder) DrawIndexed(indexCount, instanceCount, firstIndex, baseVertex, firstInstance uint32) {
	wgpuRenderBundleEncoderDrawIndexed.Call(
		uintptr(p.ref),
		uintptr(indexCount),
		uintptr(instanceCount),
		uintptr(firstIndex),
		uintptr(baseVertex),
		uintptr(firstInstance),
	)
}

func (p *RenderBundleEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderBundleEncoderDrawIndexedIndirect.Call(
		uintptr(p.ref),
		uintptr(indirectBuffer.ref),
		uintptr(indirectOffset),
	)
}

func (p *RenderBundleEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	wgpuRenderBundleEncoderDrawIndirect.Call(
		uintptr(p.ref),
		uintptr(indirectBuffer.ref),
		uintptr(indirectOffset),
	)
}

func (p *RenderBundleEncoder) Finish(descriptor *RenderBundleDescriptor) *RenderBundle {
	var desc wgpuRenderBundleDescriptor

	if descriptor != nil {
		desc.label = cstring(descriptor.Label)
	}

	ref, _, _ := wgpuRenderBundleEncoderFinish.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&desc)))
	if ref == 0 {
		panic("Failed to accquire RenderBundle")
	}
	return &RenderBundle{ref: wgpuRenderBundle(ref)}
}

func (p *RenderBundleEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := cstring(markerLabel)

	wgpuRenderBundleEncoderInsertDebugMarker.Call(uintptr(p.ref), uintptr(unsafe.Pointer(&markerLabelStr)))
}

func (p *RenderBundleEncoder) PopDebugGroup() {
	wgpuRenderBundleEncoderPopDebugGroup.Call(uintptr(p.ref))
}

func (p *RenderBundleEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := cstring(groupLabel)

	wgpuRenderBundleEncoderPushDebugGroup.Call(uintptr(p.ref), uintptr(unsafe.Pointer(groupLabelStr)))
}

func (p *RenderBundleEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		wgpuRenderBundleEncoderSetBindGroup.Call(uintptr(p.ref), uintptr(groupIndex), uintptr(group.ref), 0, 0)
	} else {
		wgpuRenderBundleEncoderSetBindGroup.Call(
			uintptr(p.ref), uintptr(groupIndex), uintptr(group.ref),
			uintptr(dynamicOffsetCount), uintptr(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
}

func (p *RenderBundleEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	wgpuRenderBundleEncoderSetIndexBuffer.Call(
		uintptr(p.ref),
		uintptr(buffer.ref),
		uintptr(format),
		uintptr(offset),
		uintptr(size),
	)
}

func (p *RenderBundleEncoder) SetPipeline(pipeline *RenderPipeline) {
	wgpuRenderBundleEncoderSetPipeline.Call(uintptr(p.ref), uintptr(pipeline.ref))
}

func (p *RenderBundleEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	wgpuRenderBundleEncoderSetVertexBuffer.Call(
		uintptr(p.ref),
		uintptr(slot),
		uintptr(buffer.ref),
		uintptr(offset),
		uintptr(size),
	)
}

func (p *Adapter) Drop()         { wgpuAdapterDrop.Call(uintptr(p.ref)) }
func (p *BindGroup) Drop()       { wgpuBindGroupDrop.Call(uintptr(p.ref)) }
func (p *BindGroupLayout) Drop() { wgpuBindGroupLayoutDrop.Call(uintptr(p.ref)) }
func (p *Buffer) Drop()          { wgpuBufferDrop.Call(uintptr(p.ref)) }
func (p *CommandBuffer) Drop()   { wgpuCommandBufferDrop.Call(uintptr(p.ref)) }
func (p *CommandEncoder) Drop()  { wgpuCommandEncoderDrop.Call(uintptr(p.ref)) }
func (p *ComputePipeline) Drop() { wgpuComputePipelineDrop.Call(uintptr(p.ref)) }
func (p *PipelineLayout) Drop()  { wgpuPipelineLayoutDrop.Call(uintptr(p.ref)) }
func (p *QuerySet) Drop()        { wgpuQuerySetDrop.Call(uintptr(p.ref)) }
func (p *RenderBundle) Drop()    { wgpuRenderBundleDrop.Call(uintptr(p.ref)) }
func (p *RenderPipeline) Drop()  { wgpuRenderPipelineDrop.Call(uintptr(p.ref)) }
func (p *Sampler) Drop()         { wgpuSamplerDrop.Call(uintptr(p.ref)) }
func (p *ShaderModule) Drop()    { wgpuShaderModuleDrop.Call(uintptr(p.ref)) }
func (p *Surface) Drop()         { wgpuSurfaceDrop.Call(uintptr(p.ref)) }
func (p *Texture) Drop()         { wgpuTextureDrop.Call(uintptr(p.ref)) }
func (p *TextureView) Drop()     { wgpuTextureViewDrop.Call(uintptr(p.ref)) }
