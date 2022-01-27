package wgpu

/*

// Common
#cgo CFLAGS: -fPIC -flto -Wall
#cgo LDFLAGS: -lwgpu_static

// Linux
#cgo linux LDFLAGS: -lm -ldl

#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux/amd64
#cgo linux,386 LDFLAGS: -L${SRCDIR}/lib/linux/386

// Windows
#cgo windows LDFLAGS: -luserenv -lws2_32 -lDwmapi -ldbghelp -ld3dcompiler -lD3D12 -lD3D11 -lDXGI -lsetupapi -lBcrypt

#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows/amd64
#cgo windows,386 LDFLAGS: -L${SRCDIR}/lib/windows/386

// Darwin
#cgo darwin LDFLAGS: -framework Cocoa -framework CoreVideo -framework IOKit -framework QuartzCore -framework Metal

#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin/amd64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin/arm64

#include "./lib/webgpu.h"
#include "./lib/wgpu.h"

extern void logCallback_cgo(WGPULogLevel level, char const *msg);

*/
import "C"

import (
	"fmt"
	"strconv"
)

const (
	CopyBytesPerRowAlignment    = 256
	QueryResolveBufferAlignment = 256
	CopyBufferAlignment         = 4
	MapAlignment                = 8
	VertexStrideAlignment       = 4
	PushConstantAlignment       = 4
	QuerySetMaxQueries          = 8192
	QuerySize                   = 8
)

var (
	Color_Transparent = Color{0, 0, 0, 0}
	Color_Black       = Color{0, 0, 0, 1}
	Color_White       = Color{1, 1, 1, 1}
	Color_Red         = Color{1, 0, 0, 1}
	Color_Green       = Color{0, 1, 0, 1}
	Color_Blue        = Color{0, 0, 1, 1}

	BlendComponent_Replace = BlendComponent{
		SrcFactor: BlendFactor_One,
		DstFactor: BlendFactor_Zero,
		Operation: BlendOperation_Add,
	}
	BlendComponent_Over = BlendComponent{
		SrcFactor: BlendFactor_One,
		DstFactor: BlendFactor_OneMinusSrcAlpha,
		Operation: BlendOperation_Add,
	}

	BlendState_Replace = BlendState{
		Color: BlendComponent_Replace,
		Alpha: BlendComponent_Replace,
	}
	BlendState_AlphaBlending = BlendState{
		Color: BlendComponent{
			SrcFactor: BlendFactor_SrcAlpha,
			DstFactor: BlendFactor_OneMinusSrcAlpha,
			Operation: BlendOperation_Add,
		},
		Alpha: BlendComponent_Over,
	}
	BlendState_PremultipliedAlphaBlending = BlendState{
		Color: BlendComponent_Over,
		Alpha: BlendComponent_Over,
	}
)

func init() {
	C.wgpuSetLogCallback(C.WGPULogCallback(C.logCallback_cgo))
}

type LogCallback func(level LogLevel, msg string)

func SetLogCallback(f LogCallback) {
	logCb = f
}

var logCb = func(level LogLevel, msg string) {
	var l string
	switch level {
	case C.WGPULogLevel_Error:
		l = "Error"
	case C.WGPULogLevel_Warn:
		l = "Warn"
	case C.WGPULogLevel_Info:
		l = "Info"
	case C.WGPULogLevel_Debug:
		l = "Debug"
	case C.WGPULogLevel_Trace:
		l = "Trace"
	default:
		l = "Unknown Level"
	}

	fmt.Printf("[go-webgpu] [%s] %s\n", l, msg)
}

func SetLogLevel(level LogLevel) {
	C.wgpuSetLogLevel(C.WGPULogLevel(level))
}

type Version C.uint32_t

func (v Version) String() string {
	return "0x" + strconv.FormatUint(uint64(v), 8)
}

func GetVersion() Version {
	return Version(C.wgpuGetVersion())
}
