package main

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

//go:embed shader.wgsl
var shader string

type State struct {
	swapChainDescriptor *wgpu.SwapChainDescriptor
	swapChain           *wgpu.SwapChain
	surface             *wgpu.Surface
	device              *wgpu.Device
	queue               *wgpu.Queue
	renderPipeline      *wgpu.RenderPipeline
}

func (s *State) Resize(width, height uint32) {
	if width > 0 && height > 0 {
		s.swapChainDescriptor.Width = width
		s.swapChainDescriptor.Height = height
		s.swapChain = s.device.CreateSwapChain(s.surface, s.swapChainDescriptor)
	}
}

func (s *State) Render() {
	nextTexture := s.swapChain.GetCurrentTextureView()
	if nextTexture == nil {
		panic("Failed to acquire next swap chain texture")
	}
	defer nextTexture.Drop()

	encoder := s.device.CreateCommandEncoder(nil)

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
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

	renderPass.SetPipeline(s.renderPipeline)
	renderPass.Draw(3, 1, 0, 0)
	renderPass.EndPass()

	s.queue.Submit(encoder.Finish(nil))
	s.swapChain.Present()
}

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
		Xlib: &wgpu.SurfaceDescriptorFromXlib{
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
		panic(err)
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

	swapChainFormat := surface.GetPreferredFormat(adapter)

	width, height := window.GetSize()

	s := &State{
		swapChainDescriptor: &wgpu.SwapChainDescriptor{
			Usage:       wgpu.TextureUsage_RenderAttachment,
			Format:      swapChainFormat,
			Width:       uint32(width),
			Height:      uint32(height),
			PresentMode: wgpu.PresentMode_Mailbox,
		},
		surface: surface,
		device:  device,
		queue:   device.GetQueue(),
	}

	shader := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})

	pipelineLayout := device.CreatePipelineLayout(nil)

	mask := ^0
	s.renderPipeline = device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
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

	s.swapChain = device.CreateSwapChain(surface, s.swapChainDescriptor)

	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		s.Resize(uint32(width), uint32(height))
		s.Render()
	})

	s.Render()

	for !window.ShouldClose() {
		glfw.WaitEvents()
	}
}
