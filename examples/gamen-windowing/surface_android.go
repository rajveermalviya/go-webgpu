//go:build android

package main

import (
	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(w display.Window) *wgpu.SurfaceDescriptor {
	if w, ok := w.(display.AndroidWindow); ok {
		return &wgpu.SurfaceDescriptor{
			AndroidNativeWindow: &wgpu.SurfaceDescriptorFromAndroidNativeWindow{
				Window: w.ANativeWindow(),
			},
		}
	}

	panic("unsupported window")
}
