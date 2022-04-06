package main

import (
	"unsafe"

	"gioui.org/app"
	"github.com/rajveermalviya/go-webgpu/wgpu"
	"golang.org/x/sys/windows"
)

var (
	kernel32          = windows.NewLazySystemDLL("kernel32.dll")
	_GetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

func getSurfaceDescriptor(ve app.ViewEvent) *wgpu.SurfaceDescriptor {
	hinstance, _, _ := _GetModuleHandleW.Call(0)

	return &wgpu.SurfaceDescriptor{
		WindowsHWND: &wgpu.SurfaceDescriptorFromWindowsHWND{
			Hwnd:      unsafe.Pointer(ve.HWND),
			Hinstance: unsafe.Pointer(hinstance),
		},
	}
}
