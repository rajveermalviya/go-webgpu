package wgpu

/*

#cgo CFLAGS: -fPIC -Wall
#cgo LDFLAGS: -lwgpu_static -lm

#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux/amd64
#cgo linux,386 LDFLAGS: -L${SRCDIR}/lib/linux/386
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows/amd64
#cgo windows,386 LDFLAGS: -L${SRCDIR}/lib/windows/386
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
