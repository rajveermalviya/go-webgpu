//go:build !windows

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
#cgo linux,!android,386 LDFLAGS: -L${SRCDIR}/lib/linux/386 -lwgpu_native

#cgo linux,!android LDFLAGS: -lm -ldl

// Darwin
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin/amd64 -lwgpu_native
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin/arm64 -lwgpu_native

#cgo darwin LDFLAGS: -framework QuartzCore -framework Metal

#include <stdlib.h>
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

extern void requestAdapterCallback_cgo(WGPURequestAdapterStatus status,
                                WGPUAdapter adapter, char const *message,
                                void *userdata);

extern void requestDeviceCallback_cgo(WGPURequestDeviceStatus status,
                               WGPUDevice device, char const *message,
                               void *userdata);

extern void bufferMapCallback_cgo(WGPUBufferMapAsyncStatus status, void *userdata);

extern void deviceUncapturedErrorCallback_cgo(WGPUErrorType type,
	                           char const * message, void * userdata);


*/
import "C"

import (
	"errors"
	"fmt"
	"runtime/cgo"
	"unsafe"
)

func init() {
	C.wgpuSetLogCallback(C.WGPULogCallback(C.logCallback_cgo))
}

func SetLogLevel(level LogLevel) {
	C.wgpuSetLogLevel(C.WGPULogLevel(level))
}

func GetVersion() Version {
	return Version(C.wgpuGetVersion())
}

func GenerateReport() GlobalReport {
	var r C.WGPUGlobalReport
	C.wgpuGenerateReport(&r)

	mapStorageReport := func(creport C.WGPUStorageReport) StorageReport {
		return StorageReport{
			NumOccupied: uint64(creport.numOccupied),
			NumVacant:   uint64(creport.numVacant),
			NumError:    uint64(creport.numError),
			ElementSize: uint64(creport.elementSize),
		}
	}

	mapHubReport := func(creport C.WGPUHubReport) *HubReport {
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
	case C.WGPUBackendType_Vulkan:
		report.Vulkan = mapHubReport(r.vulkan)
	case C.WGPUBackendType_Metal:
		report.Metal = mapHubReport(r.metal)
	case C.WGPUBackendType_D3D12:
		report.Dx12 = mapHubReport(r.dx12)
	case C.WGPUBackendType_D3D11:
		report.Dx11 = mapHubReport(r.dx11)
	case C.WGPUBackendType_OpenGL:
		report.Gl = mapHubReport(r.gl)
	}

	return report
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
	Adapter             struct{ ref C.WGPUAdapter }
	BindGroup           struct{ ref C.WGPUBindGroup }
	BindGroupLayout     struct{ ref C.WGPUBindGroupLayout }
	Buffer              struct{ ref C.WGPUBuffer }
	CommandBuffer       struct{ ref C.WGPUCommandBuffer }
	CommandEncoder      struct{ ref C.WGPUCommandEncoder }
	ComputePassEncoder  struct{ ref C.WGPUComputePassEncoder }
	ComputePipeline     struct{ ref C.WGPUComputePipeline }
	PipelineLayout      struct{ ref C.WGPUPipelineLayout }
	QuerySet            struct{ ref C.WGPUQuerySet }
	Queue               struct{ ref C.WGPUQueue }
	RenderBundle        struct{ ref C.WGPURenderBundle }
	RenderBundleEncoder struct{ ref C.WGPURenderBundleEncoder }
	RenderPassEncoder   struct{ ref C.WGPURenderPassEncoder }
	RenderPipeline      struct{ ref C.WGPURenderPipeline }
	Sampler             struct{ ref C.WGPUSampler }
	ShaderModule        struct{ ref C.WGPUShaderModule }
	Surface             struct{ ref C.WGPUSurface }
	Texture             struct{ ref C.WGPUTexture }
	TextureView         struct{ ref C.WGPUTextureView }

	SwapChain struct {
		ref    C.WGPUSwapChain
		device *Device
	}

	Device struct {
		ref     C.WGPUDevice
		errChan chan *Error
		handle  cgo.Handle
	}
)

func RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	var opts C.WGPURequestAdapterOptions

	if options != nil {
		if options.CompatibleSurface != nil {
			opts.compatibleSurface = options.CompatibleSurface.ref
		}
		opts.powerPreference = C.WGPUPowerPreference(options.PowerPreference)
		opts.forceFallbackAdapter = C.bool(options.ForceFallbackAdapter)

		if options.AdapterExtras != nil {
			adapterExtras := (*C.WGPUAdapterExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUAdapterExtras{}))))
			defer C.free(unsafe.Pointer(adapterExtras))

			adapterExtras.chain.next = nil
			adapterExtras.chain.sType = C.WGPUSType_AdapterExtras
			adapterExtras.backend = C.WGPUBackendType(options.AdapterExtras.BackendType)

			opts.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(adapterExtras))
		}
	}

	var status RequestAdapterStatus
	var adapter *Adapter

	var cb requestAdapterCB = func(s RequestAdapterStatus, a *Adapter, _ string) {
		status = s
		adapter = a
	}
	handle := cgo.NewHandle(cb)
	C.wgpuInstanceRequestAdapter(nil, &opts, C.WGPURequestAdapterCallback(C.requestAdapterCallback_cgo), unsafe.Pointer(&handle))

	if status != RequestAdapterStatus_Success {
		return nil, errors.New("failed to request adapter")
	}
	return adapter, nil
}

func CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	var desc C.WGPUSurfaceDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		if descriptor.WindowsHWND != nil {
			windowsHWND := (*C.WGPUSurfaceDescriptorFromWindowsHWND)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromWindowsHWND{}))))
			defer C.free(unsafe.Pointer(windowsHWND))

			windowsHWND.chain.next = nil
			windowsHWND.chain.sType = C.WGPUSType_SurfaceDescriptorFromWindowsHWND
			windowsHWND.hinstance = descriptor.WindowsHWND.Hinstance
			windowsHWND.hwnd = descriptor.WindowsHWND.Hwnd

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(windowsHWND))
		}

		if descriptor.XcbWindow != nil {
			xcbWindow := (*C.WGPUSurfaceDescriptorFromXcbWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromXcbWindow{}))))
			defer C.free(unsafe.Pointer(xcbWindow))

			xcbWindow.chain.next = nil
			xcbWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromXcbWindow
			xcbWindow.connection = descriptor.XcbWindow.Connection
			xcbWindow.window = C.uint32_t(descriptor.XcbWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xcbWindow))
		}

		if descriptor.XlibWindow != nil {
			xlibWindow := (*C.WGPUSurfaceDescriptorFromXlibWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromXlibWindow{}))))
			defer C.free(unsafe.Pointer(xlibWindow))

			xlibWindow.chain.next = nil
			xlibWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromXlibWindow
			xlibWindow.display = descriptor.XlibWindow.Display
			xlibWindow.window = C.uint32_t(descriptor.XlibWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xlibWindow))
		}

		if descriptor.MetalLayer != nil {
			metalLayer := (*C.WGPUSurfaceDescriptorFromMetalLayer)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromMetalLayer{}))))
			defer C.free(unsafe.Pointer(metalLayer))

			metalLayer.chain.next = nil
			metalLayer.chain.sType = C.WGPUSType_SurfaceDescriptorFromMetalLayer
			metalLayer.layer = descriptor.MetalLayer.Layer

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(metalLayer))
		}

		if descriptor.WaylandSurface != nil {
			waylandSurface := (*C.WGPUSurfaceDescriptorFromWaylandSurface)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromWaylandSurface{}))))
			defer C.free(unsafe.Pointer(waylandSurface))

			waylandSurface.chain.next = nil
			waylandSurface.chain.sType = C.WGPUSType_SurfaceDescriptorFromWaylandSurface
			waylandSurface.display = descriptor.WaylandSurface.Display
			waylandSurface.surface = descriptor.WaylandSurface.Surface

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(waylandSurface))
		}

		if descriptor.AndroidNativeWindow != nil {
			androidNativeWindow := (*C.WGPUSurfaceDescriptorFromAndroidNativeWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromAndroidNativeWindow{}))))
			defer C.free(unsafe.Pointer(androidNativeWindow))

			androidNativeWindow.chain.next = nil
			androidNativeWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromAndroidNativeWindow
			androidNativeWindow.window = descriptor.AndroidNativeWindow.Window

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(androidNativeWindow))
		}
	}

	ref := C.wgpuInstanceCreateSurface(nil, &desc)
	if ref == nil {
		panic("Failed to acquire Surface")
	}
	return &Surface{ref}
}

func (p *Adapter) EnumerateFeatures() []FeatureName {
	size := C.wgpuAdapterEnumerateFeatures(p.ref, nil)
	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	C.wgpuAdapterEnumerateFeatures(p.ref, (*C.WGPUFeatureName)(unsafe.Pointer(&features[0])))
	return features
}

func (p *Adapter) HasFeature(feature FeatureName) bool {
	hasFeature := C.wgpuAdapterHasFeature(p.ref, C.WGPUFeatureName(feature))
	return bool(hasFeature)
}

func (p *Adapter) GetLimits() SupportedLimits {
	var supportedLimits C.WGPUSupportedLimits

	extras := (*C.WGPUSupportedLimitsExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSupportedLimitsExtras{}))))
	defer C.free(unsafe.Pointer(extras))
	supportedLimits.nextInChain = (*C.WGPUChainedStructOut)(unsafe.Pointer(extras))

	C.wgpuAdapterGetLimits(p.ref, &supportedLimits)

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
	var props C.WGPUAdapterProperties

	C.wgpuAdapterGetProperties(p.ref, &props)

	return AdapterProperties{
		VendorID:          uint32(props.vendorID),
		DeviceID:          uint32(props.deviceID),
		Name:              C.GoString(props.name),
		DriverDescription: C.GoString(props.driverDescription),
		AdapterType:       AdapterType(props.adapterType),
		BackendType:       BackendType(props.backendType),
	}
}

func (p *Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {
	var desc C.WGPUDeviceDescriptor

	desc.requiredLimits = (*C.WGPURequiredLimits)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURequiredLimits{}))))
	defer C.free(unsafe.Pointer(desc.requiredLimits))

	if descriptor != nil {
		if descriptor.RequiredLimits != nil {
			l := descriptor.RequiredLimits.Limits
			desc.requiredLimits.limits = C.WGPULimits{
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

			requiredLimitsExtras := (*C.WGPURequiredLimitsExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURequiredLimitsExtras{}))))
			defer C.free(unsafe.Pointer(requiredLimitsExtras))

			requiredLimitsExtras.chain.next = nil
			requiredLimitsExtras.chain.sType = C.WGPUSType_RequiredLimitsExtras
			requiredLimitsExtras.maxPushConstantSize = C.uint32_t(l.MaxPushConstantSize)
			requiredLimitsExtras.maxBufferSize = C.uint64_t(l.MaxBufferSize)

			desc.requiredLimits.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(requiredLimitsExtras))
		} else {
			desc.requiredLimits.limits = C.WGPULimits{}
			desc.requiredLimits.nextInChain = nil
		}

		if descriptor.DeviceExtras != nil {
			deviceExtras := (*C.WGPUDeviceExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUDeviceExtras{}))))
			defer C.free(unsafe.Pointer(deviceExtras))

			deviceExtras.chain.next = nil
			deviceExtras.chain.sType = C.WGPUSType_DeviceExtras
			deviceExtras.nativeFeatures = C.WGPUNativeFeature(descriptor.DeviceExtras.NativeFeatures)

			if descriptor.DeviceExtras.Label != "" {
				label := C.CString(descriptor.DeviceExtras.Label)
				defer C.free(unsafe.Pointer(label))

				deviceExtras.label = label
			} else {
				deviceExtras.label = nil
			}

			if descriptor.DeviceExtras.TracePath != "" {
				tracePath := C.CString(descriptor.DeviceExtras.TracePath)
				defer C.free(unsafe.Pointer(tracePath))

				deviceExtras.tracePath = tracePath
			} else {
				deviceExtras.tracePath = nil
			}

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(deviceExtras))
		}
	}

	var status RequestDeviceStatus
	var device *Device

	var cb requestDeviceCB = func(s RequestDeviceStatus, d *Device, _ string) {
		status = s
		device = d
	}
	handle := cgo.NewHandle(cb)
	C.wgpuAdapterRequestDevice(p.ref, &desc, C.WGPURequestDeviceCallback(C.requestDeviceCallback_cgo), unsafe.Pointer(&handle))

	if status != RequestDeviceStatus_Success {
		return nil, errors.New("failed to request device")
	}

	device.errChan = make(chan *Error, 1)
	device.handle = cgo.NewHandle(device)
	C.wgpuDeviceSetUncapturedErrorCallback(device.ref, C.WGPUErrorCallback(C.deviceUncapturedErrorCallback_cgo), unsafe.Pointer(&device.handle))

	return device, nil
}

func (p *Device) Drop() {
	C.wgpuDeviceDrop(p.ref)

loop:
	for {
		select {
		case err := <-p.errChan:
			fmt.Printf("go-webgpu: uncaptured error: %s\n", err.Error())
		default:
			break loop
		}
	}

	p.handle.Delete()
	close(p.errChan)
}

func (p *Device) EnumerateFeatures() []FeatureName {
	size := C.wgpuDeviceEnumerateFeatures(p.ref, nil)
	if size == 0 {
		return nil
	}

	features := make([]FeatureName, size)
	C.wgpuDeviceEnumerateFeatures(p.ref, (*C.WGPUFeatureName)(unsafe.Pointer(&features[0])))
	return features
}

func (p *Device) HasFeature(feature FeatureName) bool {
	hasFeature := C.wgpuDeviceHasFeature(p.ref, C.WGPUFeatureName(feature))
	return bool(hasFeature)
}

func (p *Device) Poll(wait bool, wrappedSubmissionIndex *WrappedSubmissionIndex) (queueEmpty bool) {
	if wrappedSubmissionIndex != nil {
		var index C.WGPUWrappedSubmissionIndex
		index.queue = wrappedSubmissionIndex.Queue.ref
		index.submissionIndex = C.WGPUSubmissionIndex(wrappedSubmissionIndex.SubmissionIndex)

		return bool(C.wgpuDevicePoll(p.ref, C.bool(wait), &index))
	}

	return bool(C.wgpuDevicePoll(p.ref, C.bool(wait), nil))
}

func (p *Device) CreateBindGroupLayout(descriptor *BindGroupLayoutDescriptor) (*BindGroupLayout, error) {
	var desc C.WGPUBindGroupLayoutDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		entryCount := len(descriptor.Entries)
		if entryCount > 0 {
			entries := C.malloc(C.size_t(entryCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupLayoutEntry{})))
			defer C.free(entries)

			entriesSlice := unsafe.Slice((*C.WGPUBindGroupLayoutEntry)(entries), entryCount)

			for i, v := range descriptor.Entries {
				entriesSlice[i] = C.WGPUBindGroupLayoutEntry{
					nextInChain: nil,
					binding:     C.uint32_t(v.Binding),
					visibility:  C.WGPUShaderStageFlags(v.Visibility),
					buffer: C.WGPUBufferBindingLayout{
						nextInChain:      nil,
						_type:            C.WGPUBufferBindingType(v.Buffer.Type),
						hasDynamicOffset: C.bool(v.Buffer.HasDynamicOffset),
						minBindingSize:   C.uint64_t(v.Buffer.MinBindingSize),
					},
					sampler: C.WGPUSamplerBindingLayout{
						nextInChain: nil,
						_type:       C.WGPUSamplerBindingType(v.Sampler.Type),
					},
					texture: C.WGPUTextureBindingLayout{
						nextInChain:   nil,
						sampleType:    C.WGPUTextureSampleType(v.Texture.SampleType),
						viewDimension: C.WGPUTextureViewDimension(v.Texture.ViewDimension),
						multisampled:  C.bool(v.Texture.Multisampled),
					},
					storageTexture: C.WGPUStorageTextureBindingLayout{
						nextInChain:   nil,
						access:        C.WGPUStorageTextureAccess(v.StorageTexture.Access),
						format:        C.WGPUTextureFormat(v.StorageTexture.Format),
						viewDimension: C.WGPUTextureViewDimension(v.StorageTexture.ViewDimension),
					},
				}
			}

			desc.entryCount = C.uint32_t(entryCount)
			desc.entries = (*C.WGPUBindGroupLayoutEntry)(entries)
		}
	}

	ref := C.wgpuDeviceCreateBindGroupLayout(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire BindGroupLayout")
	}

	return &BindGroupLayout{ref}, nil
}

func (p *Device) CreateBindGroup(descriptor *BindGroupDescriptor) (*BindGroup, error) {
	var desc C.WGPUBindGroupDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		entryCount := len(descriptor.Entries)
		if entryCount > 0 {
			entries := C.malloc(C.size_t(entryCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupEntry{})))
			defer C.free(entries)

			entriesSlice := unsafe.Slice((*C.WGPUBindGroupEntry)(entries), entryCount)

			for i, v := range descriptor.Entries {
				entry := C.WGPUBindGroupEntry{
					binding: C.uint32_t(v.Binding),
					offset:  C.uint64_t(v.Offset),
					size:    C.uint64_t(v.Size),
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

				entriesSlice[i] = entry
			}

			desc.entryCount = C.uint32_t(entryCount)
			desc.entries = (*C.WGPUBindGroupEntry)(entries)
		}
	}

	ref := C.wgpuDeviceCreateBindGroup(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire BindGroup")
	}

	return &BindGroup{ref}, nil
}

func (p *Device) CreateBuffer(descriptor *BufferDescriptor) (*Buffer, error) {
	var desc C.WGPUBufferDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		desc.usage = C.WGPUBufferUsageFlags(descriptor.Usage)
		desc.size = C.uint64_t(descriptor.Size)
		desc.mappedAtCreation = C.bool(descriptor.MappedAtCreation)
	}

	ref := C.wgpuDeviceCreateBuffer(p.ref, &desc)

	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire Buffer")
	}

	return &Buffer{ref}, nil
}

func (p *Device) CreateCommandEncoder(descriptor *CommandEncoderDescriptor) (*CommandEncoder, error) {
	var desc C.WGPUCommandEncoderDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuDeviceCreateCommandEncoder(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire CommandEncoder")
	}

	return &CommandEncoder{ref}, nil
}

func (p *Device) CreateComputePipeline(descriptor *ComputePipelineDescriptor) (*ComputePipeline, error) {
	var desc C.WGPUComputePipelineDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		var compute C.WGPUProgrammableStageDescriptor
		if descriptor.Compute.Module != nil {
			compute.module = descriptor.Compute.Module.ref
		}
		if descriptor.Compute.EntryPoint != "" {
			entryPoint := C.CString(descriptor.Compute.EntryPoint)
			defer C.free(unsafe.Pointer(entryPoint))

			compute.entryPoint = entryPoint
		}
		desc.compute = compute
	}

	ref := C.wgpuDeviceCreateComputePipeline(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire ComputePipeline")
	}

	return &ComputePipeline{ref}, nil
}

func (p *Device) CreatePipelineLayout(descriptor *PipelineLayoutDescriptor) (*PipelineLayout, error) {
	var desc C.WGPUPipelineLayoutDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		bindGroupLayoutCount := len(descriptor.BindGroupLayouts)
		if bindGroupLayoutCount > 0 {
			bindGroupLayouts := C.malloc(C.size_t(bindGroupLayoutCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupLayout(nil))))
			defer C.free(bindGroupLayouts)

			bindGroupLayoutsSlice := unsafe.Slice((*C.WGPUBindGroupLayout)(bindGroupLayouts), bindGroupLayoutCount)

			for i, v := range descriptor.BindGroupLayouts {
				bindGroupLayoutsSlice[i] = v.ref
			}

			desc.bindGroupLayoutCount = C.uint32_t(bindGroupLayoutCount)
			desc.bindGroupLayouts = (*C.WGPUBindGroupLayout)(bindGroupLayouts)
		}

		if descriptor.PipelineLayoutExtras != nil {
			pipelineLayoutExtras := (*C.WGPUPipelineLayoutExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUPipelineLayoutExtras{}))))
			defer C.free(unsafe.Pointer(pipelineLayoutExtras))

			pipelineLayoutExtras.chain.next = nil
			pipelineLayoutExtras.chain.sType = C.WGPUSType_PipelineLayoutExtras

			pushConstantRangeCount := len(descriptor.PipelineLayoutExtras.PushConstantRanges)
			if pushConstantRangeCount > 0 {
				pushConstantRanges := C.malloc(C.size_t(pushConstantRangeCount) * C.size_t(unsafe.Sizeof(C.WGPUPushConstantRange{})))
				defer C.free(pushConstantRanges)

				pushConstantRangesSlice := unsafe.Slice((*C.WGPUPushConstantRange)(pushConstantRanges), pushConstantRangeCount)

				for i, v := range descriptor.PipelineLayoutExtras.PushConstantRanges {
					pushConstantRangesSlice[i] = C.WGPUPushConstantRange{
						stages: C.WGPUShaderStageFlags(v.Stages),
						start:  C.uint32_t(v.Start),
						end:    C.uint32_t(v.End),
					}
				}

				pipelineLayoutExtras.pushConstantRangeCount = C.uint32_t(pushConstantRangeCount)
				pipelineLayoutExtras.pushConstantRanges = (*C.WGPUPushConstantRange)(pushConstantRanges)
			}

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(pipelineLayoutExtras))
		} else {
			desc.nextInChain = nil
		}
	}

	ref := C.wgpuDeviceCreatePipelineLayout(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire PipelineLayout")
	}

	return &PipelineLayout{ref}, nil
}

func (p *Device) CreateRenderPipeline(descriptor *RenderPipelineDescriptor) (*RenderPipeline, error) {
	var desc C.WGPURenderPipelineDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		if descriptor.Layout != nil {
			desc.layout = descriptor.Layout.ref
		}

		// vertex
		{
			vertex := descriptor.Vertex

			var vert C.WGPUVertexState

			if vertex.Module != nil {
				vert.module = vertex.Module.ref
			}

			if vertex.EntryPoint != "" {
				entryPoint := C.CString(vertex.EntryPoint)
				defer C.free(unsafe.Pointer(entryPoint))

				vert.entryPoint = entryPoint
			}

			bufferCount := len(vertex.Buffers)
			if bufferCount > 0 {
				buffers := C.malloc(C.size_t(bufferCount) * C.size_t(unsafe.Sizeof(C.WGPUVertexBufferLayout{})))
				defer C.free(buffers)

				buffersSlice := unsafe.Slice((*C.WGPUVertexBufferLayout)(buffers), bufferCount)

				for i, v := range vertex.Buffers {
					buffer := C.WGPUVertexBufferLayout{
						arrayStride: C.uint64_t(v.ArrayStride),
						stepMode:    C.WGPUVertexStepMode(v.StepMode),
					}

					attributeCount := len(v.Attributes)
					if attributeCount > 0 {
						attributes := C.malloc(C.size_t(attributeCount) * C.size_t(unsafe.Sizeof(C.WGPUVertexAttribute{})))
						defer C.free(attributes)

						attributesSlice := unsafe.Slice((*C.WGPUVertexAttribute)(attributes), attributeCount)

						for j, attribute := range v.Attributes {
							attributesSlice[j] = C.WGPUVertexAttribute{
								format:         C.WGPUVertexFormat(attribute.Format),
								offset:         C.uint64_t(attribute.Offset),
								shaderLocation: C.uint32_t(attribute.ShaderLocation),
							}
						}

						buffer.attributeCount = C.uint32_t(attributeCount)
						buffer.attributes = (*C.WGPUVertexAttribute)(attributes)
					}

					buffersSlice[i] = buffer
				}

				vert.bufferCount = C.uint32_t(bufferCount)
				vert.buffers = (*C.WGPUVertexBufferLayout)(buffers)
			}

			desc.vertex = vert
		}

		desc.primitive = C.WGPUPrimitiveState{
			topology:         C.WGPUPrimitiveTopology(descriptor.Primitive.Topology),
			stripIndexFormat: C.WGPUIndexFormat(descriptor.Primitive.StripIndexFormat),
			frontFace:        C.WGPUFrontFace(descriptor.Primitive.FrontFace),
			cullMode:         C.WGPUCullMode(descriptor.Primitive.CullMode),
		}

		if descriptor.DepthStencil != nil {
			depthStencil := descriptor.DepthStencil

			ds := (*C.WGPUDepthStencilState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUDepthStencilState{}))))
			defer C.free(unsafe.Pointer(ds))

			ds.nextInChain = nil
			ds.format = C.WGPUTextureFormat(depthStencil.Format)
			ds.depthWriteEnabled = C.bool(depthStencil.DepthWriteEnabled)
			ds.depthCompare = C.WGPUCompareFunction(depthStencil.DepthCompare)
			ds.stencilFront = C.WGPUStencilFaceState{
				compare:     C.WGPUCompareFunction(depthStencil.StencilFront.Compare),
				failOp:      C.WGPUStencilOperation(depthStencil.StencilFront.FailOp),
				depthFailOp: C.WGPUStencilOperation(depthStencil.StencilFront.DepthFailOp),
				passOp:      C.WGPUStencilOperation(depthStencil.StencilFront.PassOp),
			}
			ds.stencilBack = C.WGPUStencilFaceState{
				compare:     C.WGPUCompareFunction(depthStencil.StencilBack.Compare),
				failOp:      C.WGPUStencilOperation(depthStencil.StencilBack.FailOp),
				depthFailOp: C.WGPUStencilOperation(depthStencil.StencilBack.DepthFailOp),
				passOp:      C.WGPUStencilOperation(depthStencil.StencilBack.PassOp),
			}
			ds.stencilReadMask = C.uint32_t(depthStencil.StencilReadMask)
			ds.stencilWriteMask = C.uint32_t(depthStencil.StencilWriteMask)
			ds.depthBias = C.int32_t(depthStencil.DepthBias)
			ds.depthBiasSlopeScale = C.float(depthStencil.DepthBiasSlopeScale)
			ds.depthBiasClamp = C.float(depthStencil.DepthBiasClamp)

			desc.depthStencil = ds
		}

		desc.multisample = C.WGPUMultisampleState{
			count:                  C.uint32_t(descriptor.Multisample.Count),
			mask:                   C.uint32_t(descriptor.Multisample.Mask),
			alphaToCoverageEnabled: C.bool(descriptor.Multisample.AlphaToCoverageEnabled),
		}

		if descriptor.Fragment != nil {
			fragment := descriptor.Fragment

			frag := (*C.WGPUFragmentState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUFragmentState{}))))
			defer C.free(unsafe.Pointer(frag))

			frag.nextInChain = nil
			if fragment.EntryPoint != "" {
				entryPoint := C.CString(fragment.EntryPoint)
				defer C.free(unsafe.Pointer(entryPoint))

				frag.entryPoint = entryPoint
			}

			if fragment.Module != nil {
				frag.module = fragment.Module.ref
			}

			targetCount := len(fragment.Targets)
			if targetCount > 0 {
				targets := C.malloc(C.size_t(targetCount) * C.size_t(unsafe.Sizeof(C.WGPUColorTargetState{})))
				defer C.free(targets)

				targetsSlice := unsafe.Slice((*C.WGPUColorTargetState)(targets), targetCount)

				for i, v := range fragment.Targets {
					target := C.WGPUColorTargetState{
						format:    C.WGPUTextureFormat(v.Format),
						writeMask: C.WGPUColorWriteMaskFlags(v.WriteMask),
					}

					if v.Blend != nil {
						blend := (*C.WGPUBlendState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUBlendState{}))))
						defer C.free(unsafe.Pointer(blend))

						blend.color = C.WGPUBlendComponent{
							operation: C.WGPUBlendOperation(v.Blend.Color.Operation),
							srcFactor: C.WGPUBlendFactor(v.Blend.Color.SrcFactor),
							dstFactor: C.WGPUBlendFactor(v.Blend.Color.DstFactor),
						}
						blend.alpha = C.WGPUBlendComponent{
							operation: C.WGPUBlendOperation(v.Blend.Alpha.Operation),
							srcFactor: C.WGPUBlendFactor(v.Blend.Alpha.SrcFactor),
							dstFactor: C.WGPUBlendFactor(v.Blend.Alpha.DstFactor),
						}

						target.blend = blend
					}

					targetsSlice[i] = target
				}

				frag.targetCount = C.uint32_t(targetCount)
				frag.targets = (*C.WGPUColorTargetState)(targets)
			} else {
				frag.targetCount = 0
				frag.targets = nil
			}

			desc.fragment = frag
		}
	}

	ref := C.wgpuDeviceCreateRenderPipeline(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire RenderPipeline")
	}

	return &RenderPipeline{ref}, nil
}

func (p *Device) CreateRenderBundleEncoder(descriptor *RenderBundleEncoderDescriptor) (*RenderBundleEncoder, error) {
	var desc C.WGPURenderBundleEncoderDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))
			desc.label = label
		}

		colorFormatsCount := len(descriptor.ColorFormats)
		if colorFormatsCount > 0 {
			colorFormats := C.malloc(C.size_t(colorFormatsCount) * C.size_t(unsafe.Sizeof(TextureFormat(0))))
			defer C.free(colorFormats)

			colorFormatsSlice := unsafe.Slice((*TextureFormat)(colorFormats), colorFormatsCount)
			copy(colorFormatsSlice, descriptor.ColorFormats)

			desc.colorFormatsCount = C.uint32_t(colorFormatsCount)
			desc.colorFormats = (*C.WGPUTextureFormat)(colorFormats)
		}

		desc.depthStencilFormat = C.WGPUTextureFormat(descriptor.DepthStencilFormat)
		desc.sampleCount = C.uint32_t(descriptor.SampleCount)
		desc.depthReadOnly = C.bool(descriptor.DepthReadOnly)
		desc.stencilReadOnly = C.bool(descriptor.StencilReadOnly)
	}

	ref := C.wgpuDeviceCreateRenderBundleEncoder(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire RenderPipeline")
	}

	return &RenderBundleEncoder{ref}, nil
}

func (p *Device) CreateSampler(descriptor *SamplerDescriptor) (*Sampler, error) {
	var desc C.WGPUSamplerDescriptor

	if descriptor != nil {
		desc = C.WGPUSamplerDescriptor{
			addressModeU:  C.WGPUAddressMode(descriptor.AddressModeU),
			addressModeV:  C.WGPUAddressMode(descriptor.AddressModeV),
			addressModeW:  C.WGPUAddressMode(descriptor.AddressModeW),
			magFilter:     C.WGPUFilterMode(descriptor.MagFilter),
			minFilter:     C.WGPUFilterMode(descriptor.MinFilter),
			mipmapFilter:  C.WGPUMipmapFilterMode(descriptor.MipmapFilter),
			lodMinClamp:   C.float(descriptor.LodMinClamp),
			lodMaxClamp:   C.float(descriptor.LodMaxClamp),
			compare:       C.WGPUCompareFunction(descriptor.Compare),
			maxAnisotropy: C.uint16_t(descriptor.MaxAnisotrophy),
		}

		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}
	}

	ref := C.wgpuDeviceCreateSampler(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire Sampler")
	}

	return &Sampler{ref}, nil
}

func (p *Device) CreateShaderModule(descriptor *ShaderModuleDescriptor) (*ShaderModule, error) {
	var desc C.WGPUShaderModuleDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		switch {
		case descriptor.SPIRVDescriptor != nil:
			spirv := (*C.WGPUShaderModuleSPIRVDescriptor)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderModuleSPIRVDescriptor{}))))
			defer C.free(unsafe.Pointer(spirv))

			codeSize := len(descriptor.SPIRVDescriptor.Code)
			if codeSize > 0 {
				code := C.CBytes(descriptor.SPIRVDescriptor.Code)
				defer C.free(code)

				spirv.codeSize = C.uint32_t(codeSize)
				spirv.code = (*C.uint32_t)(code)
			} else {
				spirv.code = nil
				spirv.codeSize = 0
			}

			spirv.chain.next = nil
			spirv.chain.sType = C.WGPUSType_ShaderModuleSPIRVDescriptor

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(spirv))

		case descriptor.WGSLDescriptor != nil:
			wgsl := (*C.WGPUShaderModuleWGSLDescriptor)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderModuleWGSLDescriptor{}))))
			defer C.free(unsafe.Pointer(wgsl))

			if descriptor.WGSLDescriptor.Code != "" {
				code := C.CString(descriptor.WGSLDescriptor.Code)
				defer C.free(unsafe.Pointer(code))

				wgsl.code = code
			} else {
				wgsl.code = nil
			}

			wgsl.chain.next = nil
			wgsl.chain.sType = C.WGPUSType_ShaderModuleWGSLDescriptor

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(wgsl))

		case descriptor.GLSLDescriptor != nil:
			glsl := (*C.WGPUShaderModuleGLSLDescriptor)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderModuleGLSLDescriptor{}))))
			defer C.free(unsafe.Pointer(glsl))

			if descriptor.GLSLDescriptor.Code != "" {
				code := C.CString(descriptor.GLSLDescriptor.Code)
				defer C.free(unsafe.Pointer(code))

				glsl.code = code
			} else {
				glsl.code = nil
			}

			defineCount := len(descriptor.GLSLDescriptor.Defines)
			if defineCount > 0 {
				shaderDefines := C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderDefine{})) * C.size_t(defineCount))
				defer C.free(shaderDefines)

				shaderDefinesSlice := unsafe.Slice((*C.WGPUShaderDefine)(shaderDefines), defineCount)
				index := 0

				for name, value := range descriptor.GLSLDescriptor.Defines {
					namePtr := C.CString(name)
					defer C.free(unsafe.Pointer(namePtr))
					valuePtr := C.CString(value)
					defer C.free(unsafe.Pointer(valuePtr))

					shaderDefinesSlice[index] = C.WGPUShaderDefine{
						name:  namePtr,
						value: valuePtr,
					}
					index++
				}

				glsl.defineCount = C.uint32_t(defineCount)
				glsl.defines = (*C.WGPUShaderDefine)(shaderDefines)
			} else {
				glsl.defineCount = 0
				glsl.defines = nil
			}

			glsl.stage = C.WGPUShaderStage(descriptor.GLSLDescriptor.ShaderStage)
			glsl.chain.next = nil
			glsl.chain.sType = C.WGPUSType_ShaderModuleGLSLDescriptor

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(glsl))
		}
	}

	ref := C.wgpuDeviceCreateShaderModule(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire ShaderModule")
	}

	return &ShaderModule{ref}, nil
}

func (p *Device) CreateSwapChain(surface *Surface, descriptor *SwapChainDescriptor) (*SwapChain, error) {
	var desc C.WGPUSwapChainDescriptor

	if descriptor != nil {
		desc = C.WGPUSwapChainDescriptor{
			usage:       C.WGPUTextureUsageFlags(descriptor.Usage),
			format:      C.WGPUTextureFormat(descriptor.Format),
			width:       C.uint32_t(descriptor.Width),
			height:      C.uint32_t(descriptor.Height),
			presentMode: C.WGPUPresentMode(descriptor.PresentMode),
		}
	}

	ref := C.wgpuDeviceCreateSwapChain(p.ref, surface.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire SwapChain")
	}

	return &SwapChain{ref: ref, device: p}, nil
}

func (p *Device) CreateTexture(descriptor *TextureDescriptor) (*Texture, error) {
	var desc C.WGPUTextureDescriptor

	if descriptor != nil {
		desc = C.WGPUTextureDescriptor{
			usage:     C.WGPUTextureUsageFlags(descriptor.Usage),
			dimension: C.WGPUTextureDimension(descriptor.Dimension),
			size: C.WGPUExtent3D{
				width:              C.uint32_t(descriptor.Size.Width),
				height:             C.uint32_t(descriptor.Size.Height),
				depthOrArrayLayers: C.uint32_t(descriptor.Size.DepthOrArrayLayers),
			},
			format:        C.WGPUTextureFormat(descriptor.Format),
			mipLevelCount: C.uint32_t(descriptor.MipLevelCount),
			sampleCount:   C.uint32_t(descriptor.SampleCount),
		}

		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}
	}

	ref := C.wgpuDeviceCreateTexture(p.ref, &desc)
	err := p.getErr()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		panic("Failed to acquire Texture")
	}

	texture := &Texture{ref}
	return texture, nil
}

func (p *Device) GetLimits() SupportedLimits {
	var supportedLimits C.WGPUSupportedLimits

	extras := (*C.WGPUSupportedLimitsExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSupportedLimitsExtras{}))))
	defer C.free(unsafe.Pointer(extras))
	supportedLimits.nextInChain = (*C.WGPUChainedStructOut)(unsafe.Pointer(extras))

	C.wgpuDeviceGetLimits(p.ref, &supportedLimits)

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
	ref := C.wgpuDeviceGetQueue(p.ref)
	if ref == nil {
		panic("Failed to acquire Queue")
	}
	return &Queue{ref}
}

func (p *Buffer) GetMappedRange(offset uint64, size uint64) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(buf), size)
}

func (p *Buffer) Unmap()   { C.wgpuBufferUnmap(p.ref) }
func (p *Buffer) Destroy() { C.wgpuBufferDestroy(p.ref) }

type BufferMapCallback func(BufferMapAsyncStatus)

func (p *Buffer) MapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) {
	handle := cgo.NewHandle(callback)

	C.wgpuBufferMapAsync(
		p.ref,
		C.WGPUMapModeFlags(mode),
		C.size_t(offset),
		C.size_t(size),
		(C.WGPUBufferMapCallback)(C.bufferMapCallback_cgo),
		unsafe.Pointer(&handle),
	)
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc C.WGPUComputePassDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire ComputePassEncoder")
	}

	return &ComputePassEncoder{ref}
}

func (p *CommandEncoder) BeginRenderPass(descriptor *RenderPassDescriptor) *RenderPassEncoder {
	var desc C.WGPURenderPassDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		colorAttachmentCount := len(descriptor.ColorAttachments)
		if colorAttachmentCount > 0 {
			colorAttachments := C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassColorAttachment{})) * C.size_t(colorAttachmentCount))
			defer C.free(colorAttachments)

			colorAttachmentsSlice := unsafe.Slice((*C.WGPURenderPassColorAttachment)(colorAttachments), colorAttachmentCount)

			for i, v := range descriptor.ColorAttachments {
				colorAttachment := C.WGPURenderPassColorAttachment{
					loadOp:  C.WGPULoadOp(v.LoadOp),
					storeOp: C.WGPUStoreOp(v.StoreOp),
					clearValue: C.WGPUColor{
						r: C.double(v.ClearValue.R),
						g: C.double(v.ClearValue.G),
						b: C.double(v.ClearValue.B),
						a: C.double(v.ClearValue.A),
					},
				}
				if v.View != nil {
					colorAttachment.view = v.View.ref
				}
				if v.ResolveTarget != nil {
					colorAttachment.resolveTarget = v.ResolveTarget.ref
				}

				colorAttachmentsSlice[i] = colorAttachment
			}

			desc.colorAttachmentCount = C.uint32_t(colorAttachmentCount)
			desc.colorAttachments = (*C.WGPURenderPassColorAttachment)(colorAttachments)
		}

		if descriptor.DepthStencilAttachment != nil {
			depthStencilAttachment := (*C.WGPURenderPassDepthStencilAttachment)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURenderPassDepthStencilAttachment{}))))
			defer C.free(unsafe.Pointer(depthStencilAttachment))

			if descriptor.DepthStencilAttachment.View != nil {
				depthStencilAttachment.view = descriptor.DepthStencilAttachment.View.ref
			}
			depthStencilAttachment.depthLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.DepthLoadOp)
			depthStencilAttachment.depthStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.DepthStoreOp)
			depthStencilAttachment.depthClearValue = C.float(descriptor.DepthStencilAttachment.DepthClearValue)
			depthStencilAttachment.depthReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)
			depthStencilAttachment.stencilLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.StencilLoadOp)
			depthStencilAttachment.stencilStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.StencilStoreOp)
			depthStencilAttachment.stencilClearValue = C.uint32_t(descriptor.DepthStencilAttachment.StencilClearValue)
			depthStencilAttachment.stencilReadOnly = C.bool(descriptor.DepthStencilAttachment.DepthReadOnly)

			desc.depthStencilAttachment = depthStencilAttachment
		}
	}

	ref := C.wgpuCommandEncoderBeginRenderPass(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire RenderPassEncoder")
	}

	return &RenderPassEncoder{ref}
}

func (p *CommandEncoder) ClearBuffer(buffer *Buffer, offset uint64, size uint64) {
	C.wgpuCommandEncoderClearBuffer(
		p.ref,
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)
}

func (p *CommandEncoder) CopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinatonOffset uint64, size uint64) {
	C.wgpuCommandEncoderCopyBufferToBuffer(
		p.ref,
		source.ref,
		C.uint64_t(sourceOffset),
		destination.ref,
		C.uint64_t(destinatonOffset),
		C.uint64_t(size),
	)
}

func (p *CommandEncoder) CopyBufferToTexture(source *ImageCopyBuffer, destination *ImageCopyTexture, copySize *Extent3D) {
	var src C.WGPUImageCopyBuffer
	if source != nil {
		if source.Buffer != nil {
			src.buffer = source.Buffer.ref
		}
		src.layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(source.Layout.Offset),
			bytesPerRow:  C.uint32_t(source.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(source.Layout.RowsPerImage),
		}
	}

	var dst C.WGPUImageCopyTexture
	if destination != nil {
		dst = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyBufferToTexture(p.ref, &src, &dst, &cpySize)
}

func (p *CommandEncoder) CopyTextureToBuffer(source *ImageCopyTexture, destination *ImageCopyBuffer, copySize *Extent3D) {
	var src C.WGPUImageCopyTexture
	if source != nil {
		src = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUImageCopyBuffer
	if destination != nil {
		if destination.Buffer != nil {
			dst.buffer = destination.Buffer.ref
		}
		dst.layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(destination.Layout.Offset),
			bytesPerRow:  C.uint32_t(destination.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(destination.Layout.RowsPerImage),
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyTextureToBuffer(p.ref, &src, &dst, &cpySize)
}

func (p *CommandEncoder) CopyTextureToTexture(source *ImageCopyTexture, destination *ImageCopyTexture, copySize *Extent3D) {
	var src C.WGPUImageCopyTexture
	if source != nil {
		src = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUImageCopyTexture
	if destination != nil {
		dst = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	C.wgpuCommandEncoderCopyTextureToTexture(p.ref, &src, &dst, &cpySize)
}

func (p *CommandEncoder) Finish(descriptor *CommandBufferDescriptor) *CommandBuffer {
	var desc C.WGPUCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuCommandEncoderFinish(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire CommandBuffer")
	}

	return &CommandBuffer{ref}
}

func (p *CommandEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	C.wgpuCommandEncoderInsertDebugMarker(p.ref, markerLabelStr)
}

func (p *CommandEncoder) PopDebugGroup() {
	C.wgpuCommandEncoderPopDebugGroup(p.ref)
}

func (p *CommandEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	C.wgpuCommandEncoderPushDebugGroup(p.ref, groupLabelStr)
}

func (p *ComputePassEncoder) DispatchWorkgroups(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	C.wgpuComputePassEncoderDispatchWorkgroups(p.ref, C.uint32_t(workgroupCountX), C.uint32_t(workgroupCountY), C.uint32_t(workgroupCountZ))
}

func (p *ComputePassEncoder) DispatchWorkgroupsIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuComputePassEncoderDispatchWorkgroupsIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
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

func (p *ComputePassEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	C.wgpuComputePassEncoderInsertDebugMarker(p.ref, markerLabelStr)
}

func (p *ComputePassEncoder) PopDebugGroup() {
	C.wgpuComputePassEncoderPopDebugGroup(p.ref)
}

func (p *ComputePassEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	C.wgpuComputePassEncoderPushDebugGroup(p.ref, groupLabelStr)
}

func (p *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref := C.wgpuComputePipelineGetBindGroupLayout(p.ref, C.uint32_t(groupIndex))
	if ref == nil {
		panic("Failed to accquire BindGroupLayout")
	}

	return &BindGroupLayout{ref}
}

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r := C.wgpuQueueSubmitForIndex(p.ref, 0, nil)
		return SubmissionIndex(r)
	}

	commandRefs := C.malloc(C.size_t(commandCount) * C.size_t(unsafe.Sizeof(C.WGPUCommandBuffer(nil))))
	defer C.free(commandRefs)

	commandRefsSlice := unsafe.Slice((*C.WGPUCommandBuffer)(commandRefs), commandCount)
	for i, v := range commands {
		commandRefsSlice[i] = v.ref
	}

	r := C.wgpuQueueSubmitForIndex(
		p.ref,
		C.uint32_t(commandCount),
		(*C.WGPUCommandBuffer)(commandRefs),
	)
	return SubmissionIndex(r)
}

func (p *Queue) WriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	C.wgpuQueueWriteBuffer(p.ref, buffer.ref, C.uint64_t(bufferOffset), buf, C.size_t(size))
}

func (p *Queue) WriteTexture(destination *ImageCopyTexture, data []byte, dataLayout *TextureDataLayout, writeSize *Extent3D) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	var dst C.WGPUImageCopyTexture
	if destination != nil {
		dst = C.WGPUImageCopyTexture{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var layout C.WGPUTextureDataLayout
	if dataLayout != nil {
		layout = C.WGPUTextureDataLayout{
			offset:       C.uint64_t(dataLayout.Offset),
			bytesPerRow:  C.uint32_t(dataLayout.BytesPerRow),
			rowsPerImage: C.uint32_t(dataLayout.RowsPerImage),
		}
	}

	var writeExtent C.WGPUExtent3D
	if writeSize != nil {
		writeExtent = C.WGPUExtent3D{
			width:              C.uint32_t(writeSize.Width),
			height:             C.uint32_t(writeSize.Height),
			depthOrArrayLayers: C.uint32_t(writeSize.DepthOrArrayLayers),
		}
	}

	C.wgpuQueueWriteTexture(p.ref, &dst, buf, C.size_t(size), &layout, &writeExtent)
}

func (p *RenderPassEncoder) SetPushConstants(stages ShaderStage, offset uint32, data []byte) {
	size := len(data)
	buf := C.CBytes(data)
	defer C.free(buf)

	C.wgpuRenderPassEncoderSetPushConstants(
		p.ref,
		C.WGPUShaderStageFlags(stages),
		C.uint32_t(offset),
		C.uint32_t(size),
		buf,
	)
}

func (p *RenderPassEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	C.wgpuRenderPassEncoderDraw(p.ref,
		C.uint32_t(vertexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstVertex),
		C.uint32_t(firstInstance),
	)
}

func (p *RenderPassEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	C.wgpuRenderPassEncoderDrawIndexed(p.ref,
		C.uint32_t(indexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstIndex),
		C.int32_t(baseVertex),
		C.uint32_t(firstInstance),
	)
}

func (p *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderPassEncoderDrawIndexedIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
}

func (p *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderPassEncoderDrawIndirect(p.ref, indirectBuffer.ref, C.uint64_t(indirectOffset))
}

func (p *RenderPassEncoder) End() {
	C.wgpuRenderPassEncoderEnd(p.ref)
}

func (p *RenderPassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		C.wgpuRenderPassEncoderSetBindGroup(
			p.ref,
			C.uint32_t(groupIndex),
			group.ref,
			0,
			nil,
		)
	} else {
		C.wgpuRenderPassEncoderSetBindGroup(
			p.ref,
			C.uint32_t(groupIndex),
			group.ref,
			C.uint32_t(dynamicOffsetCount),
			(*C.uint32_t)(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
}

func (p *RenderPassEncoder) SetBlendConstant(color *Color) {
	C.wgpuRenderPassEncoderSetBlendConstant(p.ref, &C.WGPUColor{
		r: C.double(color.R),
		g: C.double(color.G),
		b: C.double(color.B),
		a: C.double(color.A),
	})
}

func (p *RenderPassEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	C.wgpuRenderPassEncoderSetIndexBuffer(
		p.ref,
		buffer.ref,
		C.WGPUIndexFormat(format),
		C.uint64_t(offset),
		C.uint64_t(size),
	)
}

func (p *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	C.wgpuRenderPassEncoderSetPipeline(p.ref, pipeline.ref)
}

func (p *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	C.wgpuRenderPassEncoderSetScissorRect(
		p.ref,
		C.uint32_t(x),
		C.uint32_t(y),
		C.uint32_t(width),
		C.uint32_t(height),
	)
}

func (p *RenderPassEncoder) SetStencilReference(reference uint32) {
	C.wgpuRenderPassEncoderSetStencilReference(p.ref, C.uint32_t(reference))
}

func (p *RenderPassEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	C.wgpuRenderPassEncoderSetVertexBuffer(
		p.ref,
		C.uint32_t(slot),
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)
}

func (p *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {
	C.wgpuRenderPassEncoderSetViewport(
		p.ref,
		C.float(x),
		C.float(y),
		C.float(width),
		C.float(height),
		C.float(minDepth),
		C.float(maxDepth),
	)
}

func (p *RenderPassEncoder) ExecuteBundles(bundles ...*RenderBundle) {
	bundlesCount := len(bundles)
	if bundlesCount == 0 {
		C.wgpuRenderPassEncoderExecuteBundles(p.ref, 0, nil)
		return
	}

	bundlesPtr := C.malloc(C.size_t(bundlesCount) * C.size_t(unsafe.Sizeof(C.WGPURenderBundle(nil))))
	defer C.free(bundlesPtr)

	bundlesSlice := unsafe.Slice((*C.WGPURenderBundle)(bundlesPtr), bundlesCount)
	for i, v := range bundles {
		bundlesSlice[i] = v.ref
	}

	C.wgpuRenderPassEncoderExecuteBundles(p.ref, C.uint32_t(bundlesCount), (*C.WGPURenderBundle)(bundlesPtr))
}

func (p *RenderPassEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	C.wgpuRenderPassEncoderInsertDebugMarker(p.ref, markerLabelStr)
}

func (p *RenderPassEncoder) PopDebugGroup() {
	C.wgpuRenderPassEncoderPopDebugGroup(p.ref)
}

func (p *RenderPassEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	C.wgpuRenderPassEncoderPushDebugGroup(p.ref, groupLabelStr)
}

func (p *RenderPipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	ref := C.wgpuRenderPipelineGetBindGroupLayout(p.ref, C.uint32_t(groupIndex))
	if ref == nil {
		panic("Failed to accquire BindGroupLayout")
	}

	return &BindGroupLayout{ref}
}

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format := C.wgpuSurfaceGetPreferredFormat(p.ref, adapter.ref)
	return TextureFormat(format)
}

func (p *Surface) GetSupportedFormats(adapter *Adapter) []TextureFormat {
	var size C.size_t
	formatsPtr := C.wgpuSurfaceGetSupportedFormats(p.ref, adapter.ref, &size)
	defer free[C.WGPUTextureFormat](unsafe.Pointer(formatsPtr), size)

	formatsSlice := unsafe.Slice((*TextureFormat)(formatsPtr), size)
	formats := make([]TextureFormat, size)
	copy(formats, formatsSlice)
	return formats
}

func (p *Surface) GetSupportedPresentModes(adapter *Adapter) []PresentMode {
	var size C.size_t
	modesPtr := C.wgpuSurfaceGetSupportedPresentModes(p.ref, adapter.ref, &size)
	defer free[C.WGPUPresentMode](unsafe.Pointer(modesPtr), size)

	modesSlice := unsafe.Slice((*PresentMode)(modesPtr), size)
	modes := make([]PresentMode, size)
	copy(modes, modesSlice)
	return modes
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

	textureView := &TextureView{ref}
	return textureView, nil
}

func (p *SwapChain) Present() {
	C.wgpuSwapChainPresent(p.ref)
}

func (p *Texture) CreateView(descriptor *TextureViewDescriptor) *TextureView {
	var desc C.WGPUTextureViewDescriptor

	if descriptor != nil {
		desc = C.WGPUTextureViewDescriptor{
			format:          C.WGPUTextureFormat(descriptor.Format),
			dimension:       C.WGPUTextureViewDimension(descriptor.Dimension),
			baseMipLevel:    C.uint32_t(descriptor.BaseMipLevel),
			mipLevelCount:   C.uint32_t(descriptor.MipLevelCount),
			baseArrayLayer:  C.uint32_t(descriptor.BaseArrayLayer),
			arrayLayerCount: C.uint32_t(descriptor.ArrayLayerCount),
			aspect:          C.WGPUTextureAspect(descriptor.Aspect),
		}

		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}
	}

	ref := C.wgpuTextureCreateView(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire TextureView")
	}

	return &TextureView{ref}
}

func (p *Texture) Destroy() {
	C.wgpuTextureDestroy(p.ref)
}

func (p *RenderBundleEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	C.wgpuRenderBundleEncoderDraw(
		p.ref,
		C.uint32_t(vertexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstVertex),
		C.uint32_t(firstInstance),
	)
}
func (p *RenderBundleEncoder) DrawIndexed(indexCount, instanceCount, firstIndex, baseVertex, firstInstance uint32) {
	C.wgpuRenderBundleEncoderDrawIndexed(
		p.ref,
		C.uint32_t(indexCount),
		C.uint32_t(instanceCount),
		C.uint32_t(firstIndex),
		C.int32_t(baseVertex),
		C.uint32_t(firstInstance),
	)
}

func (p *RenderBundleEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderBundleEncoderDrawIndexedIndirect(
		p.ref,
		indirectBuffer.ref,
		C.uint64_t(indirectOffset),
	)
}

func (p *RenderBundleEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	C.wgpuRenderBundleEncoderDrawIndirect(
		p.ref,
		indirectBuffer.ref,
		C.uint64_t(indirectOffset),
	)
}

func (p *RenderBundleEncoder) Finish(descriptor *RenderBundleDescriptor) *RenderBundle {
	var desc C.WGPURenderBundleDescriptor

	if descriptor != nil {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))
		desc.label = label
	}

	ref := C.wgpuRenderBundleEncoderFinish(p.ref, &desc)
	if ref == nil {
		panic("Failed to accquire RenderBundle")
	}
	return &RenderBundle{ref}
}

func (p *RenderBundleEncoder) InsertDebugMarker(markerLabel string) {
	markerLabelStr := C.CString(markerLabel)
	defer C.free(unsafe.Pointer(markerLabelStr))

	C.wgpuRenderBundleEncoderInsertDebugMarker(p.ref, markerLabelStr)
}

func (p *RenderBundleEncoder) PopDebugGroup() {
	C.wgpuRenderBundleEncoderPopDebugGroup(p.ref)
}

func (p *RenderBundleEncoder) PushDebugGroup(groupLabel string) {
	groupLabelStr := C.CString(groupLabel)
	defer C.free(unsafe.Pointer(groupLabelStr))

	C.wgpuRenderBundleEncoderPushDebugGroup(p.ref, groupLabelStr)
}

func (p *RenderBundleEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	dynamicOffsetCount := len(dynamicOffsets)
	if dynamicOffsetCount == 0 {
		C.wgpuRenderBundleEncoderSetBindGroup(p.ref, C.uint32_t(groupIndex), group.ref, 0, nil)
	} else {
		C.wgpuRenderBundleEncoderSetBindGroup(
			p.ref, C.uint32_t(groupIndex), group.ref,
			C.uint32_t(dynamicOffsetCount), (*C.uint32_t)(unsafe.Pointer(&dynamicOffsets[0])),
		)
	}
}

func (p *RenderBundleEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	C.wgpuRenderBundleEncoderSetIndexBuffer(
		p.ref,
		buffer.ref,
		C.WGPUIndexFormat(format),
		C.uint64_t(offset),
		C.uint64_t(size),
	)
}

func (p *RenderBundleEncoder) SetPipeline(pipeline *RenderPipeline) {
	C.wgpuRenderBundleEncoderSetPipeline(p.ref, pipeline.ref)
}

func (p *RenderBundleEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	C.wgpuRenderBundleEncoderSetVertexBuffer(
		p.ref,
		C.uint32_t(slot),
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)
}

func (p *Adapter) Drop()         { C.wgpuAdapterDrop(p.ref) }
func (p *BindGroup) Drop()       { C.wgpuBindGroupDrop(p.ref) }
func (p *BindGroupLayout) Drop() { C.wgpuBindGroupLayoutDrop(p.ref) }
func (p *Buffer) Drop()          { C.wgpuBufferDrop(p.ref) }
func (p *CommandBuffer) Drop()   { C.wgpuCommandBufferDrop(p.ref) }
func (p *CommandEncoder) Drop()  { C.wgpuCommandEncoderDrop(p.ref) }
func (p *ComputePipeline) Drop() { C.wgpuComputePipelineDrop(p.ref) }
func (p *PipelineLayout) Drop()  { C.wgpuPipelineLayoutDrop(p.ref) }
func (p *QuerySet) Drop()        { C.wgpuQuerySetDrop(p.ref) }
func (p *RenderBundle) Drop()    { C.wgpuRenderBundleDrop(p.ref) }
func (p *RenderPipeline) Drop()  { C.wgpuRenderPipelineDrop(p.ref) }
func (p *Sampler) Drop()         { C.wgpuSamplerDrop(p.ref) }
func (p *ShaderModule) Drop()    { C.wgpuShaderModuleDrop(p.ref) }
func (p *Surface) Drop()         { C.wgpuSurfaceDrop(p.ref) }
func (p *Texture) Drop()         { C.wgpuTextureDrop(p.ref) }
func (p *TextureView) Drop()     { C.wgpuTextureViewDrop(p.ref) }
