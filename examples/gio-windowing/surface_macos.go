//go:build darwin && !ios

package main

import (
	"unsafe"

	"gioui.org/app"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

// /*

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework QuartzCore

// #import <Cocoa/Cocoa.h>
// #import <QuartzCore/CAMetalLayer.h>

// CFTypeRef metalLayerFromNSView(CFTypeRef nsViewRef) {
// 	NSView *ns_view = (__bridge NSView *)nsViewRef;
// 	ns_view.wantsLayer = YES;
// 	ns_view.layer = [CAMetalLayer layer];
// 	return ns_view.layer;
// }

// */
// import "C"

func getSurfaceDescriptor(ve app.ViewEvent) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		MetalLayer: &wgpu.SurfaceDescriptorFromMetalLayer{
			Layer: unsafe.Pointer(ve.Layer),
		},
	}
}
