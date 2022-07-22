package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"sync"

	_ "embed"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	switch os.Getenv("WGPU_LOG_LEVEL") {
	case "OFF":
		wgpu.SetLogLevel(wgpu.LogLevel_Off)
	case "ERROR":
		wgpu.SetLogLevel(wgpu.LogLevel_Error)
	case "WARN":
		wgpu.SetLogLevel(wgpu.LogLevel_Warn)
	case "INFO":
		wgpu.SetLogLevel(wgpu.LogLevel_Info)
	case "DEBUG":
		wgpu.SetLogLevel(wgpu.LogLevel_Debug)
	case "TRACE":
		wgpu.SetLogLevel(wgpu.LogLevel_Trace)
	}
}

//go:embed shader.wgsl
var shader string

func main() {
	go func() {
		w := app.NewWindow(app.CustomRenderer(true))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) (err error) {
	var ops op.Ops

	var init sync.Once

	var surface *wgpu.Surface
	var device *wgpu.Device
	var queue *wgpu.Queue
	var swapChainFormat wgpu.TextureFormat
	var swapChain *wgpu.SwapChain
	var pipeline *wgpu.RenderPipeline

	var size image.Point

	for e := range w.Events() {
		switch e := e.(type) {
		case app.ViewEvent:
			init.Do(func() {
				surface = wgpu.CreateSurface(getSurfaceDescriptor(e))
				if surface == nil {
					panic("got nil surface")
				}

				adapter, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
					CompatibleSurface: surface,
				})
				if err != nil {
					panic(err)
				}

				device, err = adapter.RequestDevice(&wgpu.DeviceDescriptor{
					DeviceExtras: &wgpu.DeviceExtras{
						Label: "Device",
					},
					RequiredLimits: &wgpu.RequiredLimits{
						Limits: wgpu.Limits{MaxBindGroups: 1},
					},
				})
				if err != nil {
					panic(err)
				}
				queue = device.GetQueue()
				swapChainFormat = surface.GetPreferredFormat(adapter)

				shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
					Label:          "shader.wgsl",
					WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
				})
				if err != nil {
					panic(err)
				}

				pipelineLayout, err := device.CreatePipelineLayout(nil)
				if err != nil {
					panic(err)
				}

				pipeline, err = device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
					Label:  "Render Pipeline",
					Layout: pipelineLayout,
					Vertex: wgpu.VertexState{
						Module:     shader,
						EntryPoint: "vs_main",
					},
					Primitive: wgpu.PrimitiveState{
						Topology:         wgpu.PrimitiveTopology_TriangleList,
						StripIndexFormat: wgpu.IndexFormat_Undefined,
						FrontFace:        wgpu.FrontFace_CCW,
						CullMode:         wgpu.CullMode_None,
					},
					Multisample: wgpu.MultisampleState{
						Count:                  1,
						Mask:                   0xFFFFFFFF,
						AlphaToCoverageEnabled: false,
					},
					Fragment: &wgpu.FragmentState{
						Module:     shader,
						EntryPoint: "fs_main",
						Targets: []wgpu.ColorTargetState{
							{
								Format:    swapChainFormat,
								Blend:     &wgpu.BlendState_Replace,
								WriteMask: wgpu.ColorWriteMask_All,
							},
						},
					},
				})
				if err != nil {
					panic(err)
				}
			})
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			fmt.Println("frame")
			if size != e.Size {
				size = e.Size

				swapChain, err = device.CreateSwapChain(surface, &wgpu.SwapChainDescriptor{
					Usage:       wgpu.TextureUsage_RenderAttachment,
					Format:      swapChainFormat,
					Width:       uint32(size.X),
					Height:      uint32(size.Y),
					PresentMode: wgpu.PresentMode_Fifo,
				})
				if err != nil {
					panic(err)
				}
			}

			func() {
				nextTexture, err := swapChain.GetCurrentTextureView()
				if err != nil {
					fmt.Printf("err: %v\n", err)
				}
				if nextTexture == nil {
					panic("Cannot acquire next swap chain texture")
				}
				defer nextTexture.Drop()

				encoder, err := device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
					Label: "Command Encoder",
				})
				if err != nil {
					panic(err)
				}

				renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
					ColorAttachments: []wgpu.RenderPassColorAttachment{
						{
							View:       nextTexture,
							LoadOp:     wgpu.LoadOp_Clear,
							StoreOp:    wgpu.StoreOp_Store,
							ClearValue: wgpu.Color_Green,
						},
					},
				})
				renderPass.SetPipeline(pipeline)
				renderPass.Draw(3, 1, 0, 0)
				renderPass.End()

				queue.Submit(encoder.Finish(nil))
				swapChain.Present()
			}()

			gtx := layout.NewContext(&ops, e)
			e.Frame(gtx.Ops)
		}
	}
	return
}
