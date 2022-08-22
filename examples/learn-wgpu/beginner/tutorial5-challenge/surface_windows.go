//go:build windows

package main

import (
	"unsafe"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(w display.Window) *wgpu.SurfaceDescriptor {
	if w, ok := w.(display.Win32Window); ok {
		return &wgpu.SurfaceDescriptor{
			WindowsHWND: &wgpu.SurfaceDescriptorFromWindowsHWND{
				Hinstance: unsafe.Pointer(w.Win32Hinstance()),
				Hwnd:      unsafe.Pointer(w.Win32Hwnd()),
			},
		}
	}

	panic("unsupported window")
}
