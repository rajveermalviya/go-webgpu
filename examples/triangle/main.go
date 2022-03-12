package main

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

//go:embed shader.wgsl
var shader string

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface := wgpu.CreateSurface(&wgpu.SurfaceDescriptor{
		XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
			Display: unsafe.Pointer(glfw.GetX11Display()),
			Window:  uint32(window.GetX11Window()),
		},
	})
	if surface == nil {
		panic("got nil surface")
	}

	adapter, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})
	if err != nil {
		// fallback to cpu
		adapter, err = wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
			CompatibleSurface:    surface,
			ForceFallbackAdapter: true,
		})
		if err != nil {
			panic(err)
		}
	}

	device, err := adapter.RequestDevice(&wgpu.DeviceDescriptor{
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
	defer device.Drop()

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

	swapChainFormat := surface.GetPreferredFormat(adapter)

	pipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
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
			Mask:                   ^uint32(0),
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

	prevWidth, prevHeight := 0, 0
	{
		width, height := window.GetSize()
		prevWidth, prevHeight = width, height
	}

	swapChain, err := device.CreateSwapChain(surface, &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      swapChainFormat,
		Width:       uint32(prevWidth),
		Height:      uint32(prevHeight),
		PresentMode: wgpu.PresentMode_Fifo,
	})
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		var nextTexture *wgpu.TextureView

		for attempt := 0; attempt < 2; attempt++ {
			width, height := window.GetSize()

			if width != prevWidth || height != prevHeight {
				prevWidth = width
				prevHeight = height

				swapChain, err = device.CreateSwapChain(
					surface,
					&wgpu.SwapChainDescriptor{
						Usage:       wgpu.TextureUsage_RenderAttachment,
						Format:      swapChainFormat,
						Width:       uint32(prevWidth),
						Height:      uint32(prevHeight),
						PresentMode: wgpu.PresentMode_Fifo,
					})
				if err != nil {
					panic(err)
				}
			}

			nextTexture, err = swapChain.GetCurrentTextureView()
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			if attempt == 0 && nextTexture == nil {
				fmt.Printf("swapChain.GetCurrentTextureView() failed; trying to create a new swap chain...\n")
				prevWidth = 0
				prevHeight = 0
				continue
			}

			break
		}

		if nextTexture == nil {
			panic("Cannot acquire next swap chain texture")
		}

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
		nextTexture.Drop()

		queue := device.GetQueue()
		queue.Submit(encoder.Finish(nil))
		swapChain.Present()

		glfw.PollEvents()
	}
}
