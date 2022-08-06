package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
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
	runtime.LockOSThread()

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

	surface := wgpu.CreateSurface(getSurfaceDescriptor(window))
	if surface == nil {
		panic("got nil surface")
	}

	adapter, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
		CompatibleSurface:    surface,
	})
	if err != nil {
		panic(err)
	}

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		panic(err)
	}
	defer device.Drop()
	queue := device.GetQueue()

	shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
	})
	if err != nil {
		panic(err)
	}
	defer shader.Drop()

	pipelineLayout, err := device.CreatePipelineLayout(nil)
	if err != nil {
		panic(err)
	}
	defer pipelineLayout.Drop()

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
			Count:                  4,
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
	defer pipeline.Drop()

	width, height := window.GetSize()
	config := &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      swapChainFormat,
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
	}
	swapChain, err := device.CreateSwapChain(surface, config)
	if err != nil {
		panic(err)
	}

	view := getMultisampledFramebuffer(device, config.Width, config.Height, swapChainFormat)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Print resource usage on pressing 'R'
		if key == glfw.KeyR && (action == glfw.Press || action == glfw.Repeat) {
			r, _ := json.MarshalIndent(wgpu.GenerateReport(), "", "  ")
			fmt.Print(string(r))
		}
	})

	for !window.ShouldClose() {
		func() {
			var nextTexture *wgpu.TextureView

			for attempt := 0; attempt < 2; attempt++ {
				width, height := window.GetSize()

				if width != int(config.Width) || height != int(config.Height) {
					config.Width = uint32(width)
					config.Height = uint32(height)

					swapChain, err = device.CreateSwapChain(surface, config)
					if err != nil {
						panic(err)
					}

					// first drop the previous view
					view.Drop()
					view = getMultisampledFramebuffer(device, config.Width, config.Height, swapChainFormat)
				}

				nextTexture, err = swapChain.GetCurrentTextureView()
				if err != nil {
					fmt.Printf("err: %v\n", err)
				}
				if attempt == 0 && nextTexture == nil {
					fmt.Printf("swapChain.GetCurrentTextureView() failed; trying to create a new swap chain...\n")
					config.Width = 0
					config.Height = 0
					continue
				}

				break
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
						View:          view,
						ResolveTarget: nextTexture,
						LoadOp:        wgpu.LoadOp_Clear,
						StoreOp:       wgpu.StoreOp_Discard,
						ClearValue:    wgpu.Color_Green,
					},
				},
			})

			renderPass.SetPipeline(pipeline)
			renderPass.Draw(3, 1, 0, 0)
			renderPass.End()

			queue.Submit(encoder.Finish(nil))
			swapChain.Present()

			glfw.PollEvents()
		}()
	}
}

func getMultisampledFramebuffer(device *wgpu.Device, width, height uint32, format wgpu.TextureFormat) *wgpu.TextureView {
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Usage:     wgpu.TextureUsage_RenderAttachment,
		Dimension: wgpu.TextureDimension_2D,
		Size: wgpu.Extent3D{
			Width:              width,
			Height:             height,
			DepthOrArrayLayers: 1,
		},
		Format:        format,
		MipLevelCount: 1,
		SampleCount:   4,
	})
	if err != nil {
		panic(err)
	}
	defer texture.Drop()

	return texture.CreateView(nil)
}
