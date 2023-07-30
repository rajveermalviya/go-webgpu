//go:build linux && !android && wayland

package wgpuext_glfw

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func GetSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		WaylandSurface: &wgpu.SurfaceDescriptorFromWaylandSurface{
			Display: unsafe.Pointer(glfw.GetWaylandDisplay()),
			Surface: unsafe.Pointer(w.GetWaylandWindow()),
		},
	}
}
