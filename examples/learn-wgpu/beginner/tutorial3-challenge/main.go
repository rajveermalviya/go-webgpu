package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shader.wgsl
var shaderCode string

//go:embed challenge.wgsl
var challengeShaderCode string

type State struct {
	surface        *wgpu.Surface
	swapChain      *wgpu.SwapChain
	device         *wgpu.Device
	queue          *wgpu.Queue
	config         *wgpu.SwapChainDescriptor
	size           dpi.PhysicalSize[uint32]
	renderPipeline *wgpu.RenderPipeline

	challengeRenderPipeline *wgpu.RenderPipeline
	useColor                bool
}

func InitState(window display.Window) (*State, error) {
	size := window.InnerSize()

	surface := wgpu.CreateSurface(getSurfaceDescriptor(window))

	adaper, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})
	if err != nil {
		return nil, err
	}
	device, err := adaper.RequestDevice(nil)
	if err != nil {
		return nil, err
	}
	queue := device.GetQueue()

	config := &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      surface.GetPreferredFormat(adaper),
		Width:       size.Width,
		Height:      size.Height,
		PresentMode: wgpu.PresentMode_Fifo,
	}
	swapChain, err := device.CreateSwapChain(surface, config)
	if err != nil {
		return nil, err
	}

	shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shaderCode,
		},
	})
	if err != nil {
		return nil, err
	}

	renderPipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label: "Render Pipeline Layout",
	})
	if err != nil {
		return nil, err
	}

	renderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Render Pipeline",
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    config.Format,
				Blend:     &wgpu.BlendState_Replace,
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopology_TriangleList,
			FrontFace: wgpu.FrontFace_CCW,
			CullMode:  wgpu.CullMode_Back,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
	})
	if err != nil {
		return nil, err
	}

	challengeShader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "challenge.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: challengeShaderCode,
		},
	})
	if err != nil {
		return nil, err
	}

	challengeRenderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Render Pipeline",
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     challengeShader,
			EntryPoint: "vs_main",
		},
		Fragment: &wgpu.FragmentState{
			Module:     challengeShader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    config.Format,
				Blend:     &wgpu.BlendState_Replace,
				WriteMask: wgpu.ColorWriteMask_All,
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopology_TriangleList,
			FrontFace: wgpu.FrontFace_CCW,
			CullMode:  wgpu.CullMode_Back,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
	})
	if err != nil {
		return nil, err
	}

	return &State{
		surface:                 surface,
		swapChain:               swapChain,
		device:                  device,
		queue:                   queue,
		config:                  config,
		size:                    size,
		renderPipeline:          renderPipeline,
		challengeRenderPipeline: challengeRenderPipeline,
		useColor:                true,
	}, nil
}

func (s *State) Resize(newSize dpi.PhysicalSize[uint32]) {
	if newSize.Width > 0 && newSize.Height > 0 {
		s.size = newSize
		s.config.Width = newSize.Width
		s.config.Height = newSize.Height

		var err error
		s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
		if err != nil {
			panic(err)
		}
	}
}

func (s *State) Render() error {
	view, err := s.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer view.Drop()

	encoder, err := s.device.CreateCommandEncoder(nil)
	if err != nil {
		return err
	}

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:   view,
			LoadOp: wgpu.LoadOp_Clear,
			ClearValue: wgpu.Color{
				R: 0.1,
				G: 0.2,
				B: 0.3,
				A: 1.0,
			},
			StoreOp: wgpu.StoreOp_Store,
		}},
	})
	if s.useColor {
		renderPass.SetPipeline(s.renderPipeline)
	} else {
		renderPass.SetPipeline(s.challengeRenderPipeline)
	}
	renderPass.Draw(3, 1, 0, 0)
	renderPass.End()

	s.queue.Submit(encoder.Finish(nil))
	s.swapChain.Present()

	return nil
}

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}

	s, err := InitState(w)
	if err != nil {
		panic(err)
	}

	w.SetKeyboardInputCallback(func(state events.ButtonState, scanCode events.ScanCode, virtualKeyCode events.VirtualKey) {
		if virtualKeyCode == events.VirtualKeySpace {
			s.useColor = state == events.ButtonStateReleased
		}
	})

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		s.Resize(dpi.PhysicalSize[uint32]{
			Width:  physicalWidth,
			Height: physicalHeight,
		})
	})

	w.SetCloseRequestedCallback(func() {
		d.Destroy()
	})

	for {
		if !d.Poll() {
			break
		}

		err := s.Render()
		if err != nil {
			errstr := err.Error()
			fmt.Println(errstr)

			switch {
			case strings.Contains(errstr, "Lost"):
				s.Resize(s.size)
			case strings.Contains(errstr, "Outdated"):
				s.Resize(s.size)
			case strings.Contains(errstr, "Timeout"):
			default:
				panic(err)
			}
		}
	}
}
