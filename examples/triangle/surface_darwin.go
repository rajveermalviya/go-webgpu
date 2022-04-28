//go:build darwin

package main

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

/*

#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore

#import <Cocoa/Cocoa.h>
#import <QuartzCore/CAMetalLayer.h>

CFTypeRef metalLayerFromNSWindow(CFTypeRef nsWindowRef) {
	NSWindow *ns_window = (__bridge NSWindow *)nsWindowRef;
	id metal_layer = NULL;
    [ns_window.contentView setWantsLayer:YES];
    metal_layer = [CAMetalLayer layer];
    [ns_window.contentView setLayer:metal_layer];
	return metal_layer;
}

*/
import "C"

func getSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		MetalLayer: &wgpu.SurfaceDescriptorFromMetalLayer{
			Layer: unsafe.Pointer(C.metalLayerFromNSWindow((C.CFTypeRef)(unsafe.Pointer(w.GetCocoaWindow())))),
		},
	}
}
