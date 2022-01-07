package wgpu

/*

#include "wrapper.h"

*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Adapter struct{ ref C.WGPUAdapter }

type SupportedLimits struct {
	Limits Limits
}

func (p *Adapter) GetLimits() SupportedLimits {
	var limits C.WGPUSupportedLimits

	C.wgpuAdapterGetLimits(p.ref, &limits)

	return SupportedLimits{limitsFromC(limits.limits)}
}

type AdapterProperties struct {
	VendorID          uint32
	DeviceID          uint32
	Name              string
	DriverDescription string
	AdapterType       AdapterType
	BackendType       BackendType
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

type DeviceExtras struct {
	NativeFeatures NativeFeature
	Label          string
	TracePath      string
}

type RequiredLimits struct {
	Limits Limits
}

type DeviceDescriptor struct {
	// unused in wgpu
	// Label     string
	// RequiredFeatures      []FeatureName
	// RequiredFeaturesCount uint32

	RequiredLimits *RequiredLimits

	// WGPUChainedStruct -> WGPUDeviceExtras
	DeviceExtras *DeviceExtras
}

func (p *Adapter) RequestDevice(descriptor DeviceDescriptor) (*Device, error) {
	var desc C.WGPUDeviceDescriptor
	defer runtime.KeepAlive(desc)

	requiredLimits := (*C.WGPURequiredLimits)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPURequiredLimits{}))))
	defer C.free(unsafe.Pointer(requiredLimits))

	requiredLimits.nextInChain = nil
	if descriptor.RequiredLimits != nil {
		requiredLimits.limits = descriptor.RequiredLimits.Limits.toC()
	}

	desc.requiredLimits = requiredLimits

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

	res := C.request_device(p.ref, &desc)
	if res.status != C.WGPURequestDeviceStatus_Success || res.device == nil {
		return nil, errors.New("failed to request device")
	}

	return &Device{res.device}, nil
}
