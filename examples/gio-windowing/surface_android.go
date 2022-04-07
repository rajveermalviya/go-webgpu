//go:build android

package main

import (
	"gioui.org/app"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(ve app.ViewEvent) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		AndroidNativeWindow: &wgpu.SurfaceDescriptorFromAndroidNativeWindow{
			Window: ve.ANativeWindow,
		},
	}
}
