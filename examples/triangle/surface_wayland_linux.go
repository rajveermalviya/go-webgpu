//go:build wayland

package main

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		WaylandSurface: &wgpu.SurfaceDescriptorFromWaylandSurface{
			Display: unsafe.Pointer(glfw.GetWaylandDisplay()),
			Surface: unsafe.Pointer(w.GetWaylandWindow()),
		},
	}
}
