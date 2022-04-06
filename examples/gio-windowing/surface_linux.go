package main

import (
	"gioui.org/app"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

func getSurfaceDescriptor(ve app.ViewEvent) *wgpu.SurfaceDescriptor {
	switch ve := ve.(type) {
	case app.X11ViewEvent:
		return &wgpu.SurfaceDescriptor{
			XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
				Display: ve.Display,
				Window:  uint32(ve.Window),
			},
		}

	case app.WaylandViewEvent:
		return &wgpu.SurfaceDescriptor{
			WaylandSurface: &wgpu.SurfaceDescriptorFromWaylandSurface{
				Display: ve.Display,
				Surface: ve.Surface,
			},
		}
	}

	panic("no display")
}
