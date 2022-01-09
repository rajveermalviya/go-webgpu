package wgpu

/*

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

*/
import "C"

// type NativeSType C.WGPUNativeSType

// const (
// 	SType_DeviceExtras  NativeSType = C.WGPUSType_DeviceExtras
// 	SType_AdapterExtras NativeSType = C.WGPUSType_AdapterExtras
// 	NativeSType_Force32 NativeSType = C.WGPUNativeSType_Force32
// )

type NativeFeature C.WGPUNativeFeature

const (
	NativeFeature_TEXTURE_ADAPTER_SPECIFIC_FORMAT_FEATURES NativeFeature = C.WGPUNativeFeature_TEXTURE_ADAPTER_SPECIFIC_FORMAT_FEATURES
)

type LogLevel C.WGPULogLevel

const (
	LogLevel_Off     LogLevel = C.WGPULogLevel_Off
	LogLevel_Error   LogLevel = C.WGPULogLevel_Error
	LogLevel_Warn    LogLevel = C.WGPULogLevel_Warn
	LogLevel_Info    LogLevel = C.WGPULogLevel_Info
	LogLevel_Debug   LogLevel = C.WGPULogLevel_Debug
	LogLevel_Trace   LogLevel = C.WGPULogLevel_Trace
	LogLevel_Force32 LogLevel = C.WGPULogLevel_Force32
)
