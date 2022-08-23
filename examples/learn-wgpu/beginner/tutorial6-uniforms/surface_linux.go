//go:build linux && !android

package main

import (
	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(w display.Window) *wgpu.SurfaceDescriptor {
	switch w := w.(type) {
	case display.WaylandWindow:
		return &wgpu.SurfaceDescriptor{
			WaylandSurface: &wgpu.SurfaceDescriptorFromWaylandSurface{
				Display: w.WlDisplay(),
				Surface: w.WlSurface(),
			},
		}

	case display.XcbWindow:
		return &wgpu.SurfaceDescriptor{
			XcbWindow: &wgpu.SurfaceDescriptorFromXcbWindow{
				Connection: w.XcbConnection(),
				Window:     w.XcbWindow(),
			},
		}

	case display.XlibWindow:
		return &wgpu.SurfaceDescriptor{
			XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
				Display: w.XlibDisplay(),
				Window:  w.XlibWindow(),
			},
		}
	}

	panic("unsupported window")
}
