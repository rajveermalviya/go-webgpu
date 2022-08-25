package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

func init() {
	runtime.LockOSThread()
}

//go:embed shader.wgsl
var shader string

type app struct {
	window          display.Window
	adapter         *wgpu.Adapter
	device          *wgpu.Device
	surface         *wgpu.Surface
	shader          *wgpu.ShaderModule
	pipeline        *wgpu.RenderPipeline
	swapChainFormat wgpu.TextureFormat
	swapChain       *wgpu.SwapChain
	config          *wgpu.SwapChainDescriptor

	hasInit        bool
	hasSurfaceInit bool
}

func (a *app) init() {
	var err error

	a.adapter, err = wgpu.RequestAdapter(nil)
	if err != nil {
		panic(err)
	}

	a.device, err = a.adapter.RequestDevice(nil)
	if err != nil {
		panic(err)
	}

	a.shader, err = a.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
	})
	if err != nil {
		panic(err)
	}

	a.hasInit = true
}

func (a *app) deinit() {
	a.hasInit = false

	if a.shader != nil {
		a.shader.Drop()
		a.shader = nil
	}
	if a.device != nil {
		a.device.Drop()
		a.device = nil
	}
	if a.adapter != nil {
		a.adapter.Drop()
		a.adapter = nil
	}
}

func (a *app) surfaceInit() {
	var err error

	a.surface = wgpu.CreateSurface(getSurfaceDescriptor(a.window))
	if a.surface == nil {
		panic("got nil surface")
	}

	a.swapChainFormat = a.surface.GetPreferredFormat(a.adapter)

	a.pipeline, err = a.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label: "Render Pipeline",
		Vertex: wgpu.VertexState{
			Module:     a.shader,
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
			Module:     a.shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    a.swapChainFormat,
					Blend:     &wgpu.BlendState_Replace,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	size := a.window.InnerSize()
	a.config = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      a.swapChainFormat,
		PresentMode: wgpu.PresentMode_Fifo,
		Width:       size.Width,
		Height:      size.Height,
	}

	a.swapChain, err = a.device.CreateSwapChain(a.surface, a.config)
	if err != nil {
		panic(err)
	}

	a.hasSurfaceInit = true
}

func (a *app) surfaceDeinit() {
	a.hasSurfaceInit = false

	if a.swapChain != nil {
		a.swapChain = nil
	}
	if a.config != nil {
		a.config = nil
	}
	if a.pipeline != nil {
		a.pipeline.Drop()
		a.pipeline = nil
	}
	if a.surface != nil {
		a.surface.Drop()
		a.surface = nil
	}
}

func (a *app) resize(width, height uint32) {
	if !a.hasInit || !a.hasSurfaceInit {
		return
	}

	if width > 0 && height > 0 {
		a.config.Width = width
		a.config.Height = height

		var err error
		a.swapChain, err = a.device.CreateSwapChain(a.surface, a.config)
		if err != nil {
			panic(err)
		}
	}
}

func (a *app) render() error {
	if !a.hasInit || !a.hasSurfaceInit {
		return nil
	}

	nextTexture, err := a.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer nextTexture.Drop()

	encoder, err := a.device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Command Encoder",
	})
	if err != nil {
		return err
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

	renderPass.SetPipeline(a.pipeline)
	renderPass.Draw(3, 1, 0, 0)
	renderPass.End()

	queue := a.device.GetQueue()
	queue.Submit(encoder.Finish(nil))
	a.swapChain.Present()

	return err
}

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer d.Destroy()

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	w.SetTitle("go-webgpu gamen-windowing example")

	a := &app{window: w}
	a.init()
	defer a.deinit()

	if w, ok := w.(display.AndroidWindowExt); ok {
		w.SetSurfaceCreatedCallback(func() { a.surfaceInit() })
		w.SetSurfaceDestroyedCallback(func() { a.surfaceDeinit() })
	} else {
		a.surfaceInit()
		defer a.surfaceDeinit()
	}

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		println(fmt.Sprintf("Resized: physicalWidth=%v physicalHeight=%v scaleFactor=%v", physicalWidth, physicalHeight, scaleFactor))

		a.resize(physicalWidth, physicalHeight)
	})

	w.SetCloseRequestedCallback(func() { d.Destroy() })

	for {
		if !d.Poll() {
			break
		}

		err := a.render()
		if err != nil {
			errstr := err.Error()
			fmt.Println(errstr)

			switch {
			case strings.Contains(errstr, "Lost"):
				size := w.InnerSize()
				a.resize(size.Width, size.Height)
			case strings.Contains(errstr, "Outdated"):
				size := w.InnerSize()
				a.resize(size.Width, size.Height)
			case strings.Contains(errstr, "Timeout"):
			default:
				panic(err)
			}
		}
	}
}
