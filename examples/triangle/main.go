package main

import (
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

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		panic(err)
	}

	x11Display := glfw.GetX11Display()
	x11Window := window.GetX11Window()

	surface := wgpu.CreateSurface(wgpu.SurfaceDescriptor{
		Xlib: &wgpu.SurfaceDescriptorFromXlib{
			Display: unsafe.Pointer(x11Display),
			Window:  uint32(x11Window),
		},
	})
	if surface == nil {
		panic("got nil surface")
	}

	adapter, err := wgpu.RequestAdapter(wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})
	if err != nil {
		panic(err)
	}

	device, err := adapter.RequestDevice(wgpu.DeviceDescriptor{
		DeviceExtras: &wgpu.DeviceExtras{
			Label: "Device",
		},
		RequiredLimits: &wgpu.RequiredLimits{
			Limits: wgpu.Limits{
				MaxBindGroups: 1,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	shader := device.CreateShaderModule(wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})

	pipelineLayout := device.CreatePipelineLayout(wgpu.PipelineLayoutDescriptor{})

	swapChainFormat := surface.GetPreferredFormat(adapter)

	mask := ^0
	pipeline := device.CreateRenderPipeline(wgpu.RenderPipelineDescriptor{
		Label:  "Render pipeline",
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
			Mask:                   uint32(mask),
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format: swapChainFormat,
				Blend: &wgpu.BlendState{
					Color: wgpu.BlendComponent{
						SrcFactor: wgpu.BlendFactor_One,
						DstFactor: wgpu.BlendFactor_Zero,
						Operation: wgpu.BlendOperation_Add,
					},
					Alpha: wgpu.BlendComponent{
						SrcFactor: wgpu.BlendFactor_One,
						DstFactor: wgpu.BlendFactor_Zero,
						Operation: wgpu.BlendOperation_Add,
					},
				},
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
	})

	prevWidth, prevHeight := window.GetSize()

	swapChain := device.CreateSwapChain(surface, wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      swapChainFormat,
		Width:       uint32(prevWidth),
		Height:      uint32(prevHeight),
		PresentMode: wgpu.PresentMode_Fifo,
	})

	for !window.ShouldClose() {
		width, height := window.GetSize()

		if width != prevWidth || height != prevHeight {
			prevWidth = width
			prevHeight = height

			swapChain = device.CreateSwapChain(surface, wgpu.SwapChainDescriptor{
				Usage:       wgpu.TextureUsage_RenderAttachment,
				Format:      swapChainFormat,
				Width:       uint32(prevWidth),
				Height:      uint32(prevHeight),
				PresentMode: wgpu.PresentMode_Fifo,
			})
		}

		nextTexture := swapChain.GetCurrentTextureView()
		if nextTexture == nil {
			panic("Cannot acquire next swap chain texture")
		}

		encoder := device.CreateCommandEncoder(wgpu.CommandEncoderDescriptor{
			Label: "Command Encoder",
		})

		renderPass := encoder.BeginRenderPass(wgpu.RenderPassDescriptor{
			ColorAttachments: []wgpu.RenderPassColorAttachment{{
				View:    nextTexture,
				LoadOp:  wgpu.LoadOp_Clear,
				StoreOp: wgpu.StoreOp_Store,
				ClearColor: wgpu.Color{
					R: 0,
					G: 1,
					B: 0,
					A: 1,
				},
			}},
		})

		renderPass.SetPipeline(pipeline)
		renderPass.Draw(3, 1, 0, 0)
		renderPass.EndPass()

		queue := device.GetQueue()
		cmdBuffer := encoder.Finish(wgpu.CommandBufferDescriptor{})
		queue.Submit([]wgpu.CommandBuffer{cmdBuffer})
		swapChain.Present()

		glfw.PollEvents()
	}

	window.Destroy()
	glfw.Terminate()
}
