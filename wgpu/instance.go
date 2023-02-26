package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

extern void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, char const *message, void *userdata);

*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type Instance struct {
	ref C.WGPUInstance
}

type InstanceDescriptor struct {
	Backends           InstanceBackend
	Dx12ShaderCompiler Dx12Compiler
	DxilPath           string
	DxcPath            string
}

func CreateInstance(descriptor *InstanceDescriptor) *Instance {
	var desc C.WGPUInstanceDescriptor

	if descriptor != nil {
		instanceExtras := (*C.WGPUInstanceExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUInstanceExtras{}))))
		defer C.free(unsafe.Pointer(instanceExtras))

		instanceExtras.chain.next = nil
		instanceExtras.chain.sType = C.WGPUSType_InstanceExtras
		instanceExtras.backends = C.WGPUInstanceBackendFlags(descriptor.Backends)
		instanceExtras.dx12ShaderCompiler = C.WGPUDx12Compiler(descriptor.Dx12ShaderCompiler)

		if descriptor.DxilPath != "" {
			dxilPath := C.CString(descriptor.DxilPath)
			defer C.free(unsafe.Pointer(dxilPath))
			instanceExtras.dxilPath = dxilPath
		}

		if descriptor.DxcPath != "" {
			dxcPath := C.CString(descriptor.DxcPath)
			defer C.free(unsafe.Pointer(dxcPath))
			instanceExtras.dxcPath = dxcPath
		}

		desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(instanceExtras))
	}

	ref := C.wgpuCreateInstance(&desc)
	if ref == nil {
		panic("Failed to acquire Instance")
	}

	return &Instance{ref}
}

type SurfaceDescriptorFromWindowsHWND struct {
	Hinstance unsafe.Pointer
	Hwnd      unsafe.Pointer
}

type SurfaceDescriptorFromXcbWindow struct {
	Connection unsafe.Pointer
	Window     uint32
}

type SurfaceDescriptorFromXlibWindow struct {
	Display unsafe.Pointer
	Window  uint32
}

type SurfaceDescriptorFromMetalLayer struct {
	Layer unsafe.Pointer
}

type SurfaceDescriptorFromWaylandSurface struct {
	Display unsafe.Pointer
	Surface unsafe.Pointer
}

type SurfaceDescriptorFromAndroidNativeWindow struct {
	Window unsafe.Pointer
}

type SurfaceDescriptor struct {
	Label string

	WindowsHWND         *SurfaceDescriptorFromWindowsHWND
	XcbWindow           *SurfaceDescriptorFromXcbWindow
	XlibWindow          *SurfaceDescriptorFromXlibWindow
	MetalLayer          *SurfaceDescriptorFromMetalLayer
	WaylandSurface      *SurfaceDescriptorFromWaylandSurface
	AndroidNativeWindow *SurfaceDescriptorFromAndroidNativeWindow
}

func (p *Instance) CreateSurface(descriptor *SurfaceDescriptor) *Surface {
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

	ref := C.wgpuInstanceCreateSurface(p.ref, &desc)
	if ref == nil {
		panic("Failed to acquire Surface")
	}
	return &Surface{ref}
}

type requestAdapterCb func(status RequestAdapterStatus, adapter *Adapter, message string)

//export gowebgpu_request_adapter_callback_go
func gowebgpu_request_adapter_callback_go(status C.WGPURequestAdapterStatus, adapter C.WGPUAdapter, message *C.char, userdata unsafe.Pointer) {
	handle := *(*cgo.Handle)(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(requestAdapterCb)
	if ok {
		cb(RequestAdapterStatus(status), &Adapter{ref: adapter}, C.GoString(message))
	}
}

type RequestAdapterOptions struct {
	CompatibleSurface    *Surface
	PowerPreference      PowerPreference
	ForceFallbackAdapter bool
	BackendType          BackendType
}

func (p *Instance) RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	var opts *C.WGPURequestAdapterOptions

	if options != nil {
		opts = &C.WGPURequestAdapterOptions{}

		if options.CompatibleSurface != nil {
			opts.compatibleSurface = options.CompatibleSurface.ref
		}
		opts.powerPreference = C.WGPUPowerPreference(options.PowerPreference)
		opts.forceFallbackAdapter = C.bool(options.ForceFallbackAdapter)

		if options.BackendType != BackendType_Null {
			adapterExtras := (*C.WGPUAdapterExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUAdapterExtras{}))))
			defer C.free(unsafe.Pointer(adapterExtras))

			adapterExtras.chain.next = nil
			adapterExtras.chain.sType = C.WGPUSType_AdapterExtras
			adapterExtras.backend = C.WGPUBackendType(options.BackendType)

			opts.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(adapterExtras))
		}
	}

	var status RequestAdapterStatus
	var adapter *Adapter

	var cb requestAdapterCb = func(s RequestAdapterStatus, a *Adapter, _ string) {
		status = s
		adapter = a
	}
	handle := cgo.NewHandle(cb)
	C.wgpuInstanceRequestAdapter(p.ref, opts, C.WGPURequestAdapterCallback(C.gowebgpu_request_adapter_callback_c), unsafe.Pointer(&handle))

	if status != RequestAdapterStatus_Success {
		return nil, errors.New("failed to request adapter")
	}
	return adapter, nil
}

type StorageReport struct {
	NumOccupied uint64
	NumVacant   uint64
	NumError    uint64
	ElementSize uint64
}

type HubReport struct {
	Adapters         StorageReport
	Devices          StorageReport
	PipelineLayouts  StorageReport
	ShaderModules    StorageReport
	BindGroupLayouts StorageReport
	BindGroups       StorageReport
	CommandBuffers   StorageReport
	RenderBundles    StorageReport
	RenderPipelines  StorageReport
	ComputePipelines StorageReport
	QuerySets        StorageReport
	Buffers          StorageReport
	Textures         StorageReport
	TextureViews     StorageReport
	Samplers         StorageReport
}

type GlobalReport struct {
	Surfaces StorageReport
	Vulkan   *HubReport
	Metal    *HubReport
	Dx12     *HubReport
	Dx11     *HubReport
	Gl       *HubReport
}

func (p *Instance) GenerateReport() GlobalReport {
	var r C.WGPUGlobalReport
	C.wgpuGenerateReport(p.ref, &r)

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

func (p *Instance) Drop() {
	C.wgpuInstanceDrop(p.ref)
}
