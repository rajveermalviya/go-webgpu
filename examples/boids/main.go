package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"unsafe"

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

const (
	// number of boid particles to simulate
	NumParticles = 1500
	// number of single-particle calculations (invocations) in each gpu work group
	ParticlesPerGroup = 64
)

//go:embed compute.wgsl
var compute string

//go:embed draw.wgsl
var draw string

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

	device, err := adapter.RequestDevice(&wgpu.DeviceDescriptor{
		DeviceExtras: &wgpu.DeviceExtras{
			Label: "Device",
		},
	})
	if err != nil {
		panic(err)
	}
	defer device.Drop()
	queue := device.GetQueue()

	width, height := window.GetSize()
	config := &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      surface.GetPreferredFormat(adapter),
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
	}

	swapChain, err := device.CreateSwapChain(surface, config)
	if err != nil {
		panic(err)
	}

	computeShader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "compute.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: compute,
		},
	})
	if err != nil {
		panic(err)
	}
	defer computeShader.Drop()

	drawShader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "draw.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: draw,
		},
	})
	if err != nil {
		panic(err)
	}
	defer drawShader.Drop()

	simParamData := []float32{
		0.04,  // deltaT
		0.1,   // rule1Distance
		0.025, // rule2Distance
		0.025, // rule3Distance
		0.02,  // rule1Scale
		0.05,  // rule2Scale
		0.005, // rule3Scale
	}

	simParamBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Simulation Param Buffer",
		Contents: wgpu.ToBytes(simParamData),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
	defer simParamBuffer.Drop()

	computeBindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Compute,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_Uniform,
					HasDynamicOffset: false,
					MinBindingSize:   uint64(len(simParamData)) * uint64(unsafe.Sizeof(float32(0))),
				},
			},
			{
				Binding:    1,
				Visibility: wgpu.ShaderStage_Compute,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_ReadOnlyStorage,
					HasDynamicOffset: false,
					MinBindingSize:   NumParticles * 16,
				},
			},
			{
				Binding:    2,
				Visibility: wgpu.ShaderStage_Compute,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_Storage,
					HasDynamicOffset: false,
					MinBindingSize:   NumParticles * 16,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	defer computeBindGroupLayout.Drop()

	computePipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:            "compute",
		BindGroupLayouts: []*wgpu.BindGroupLayout{computeBindGroupLayout},
	})
	if err != nil {
		panic(err)
	}
	defer computePipelineLayout.Drop()

	renderPipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label: "render",
	})
	if err != nil {
		panic(err)
	}
	defer renderPipelineLayout.Drop()

	renderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: renderPipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     drawShader,
			EntryPoint: "main_vs",
			Buffers: []wgpu.VertexBufferLayout{
				{
					ArrayStride: 4 * 4,
					StepMode:    wgpu.VertexStepMode_Instance,
					Attributes: []wgpu.VertexAttribute{
						{
							Format:         wgpu.VertexFormat_Float32x2,
							Offset:         0,
							ShaderLocation: 0,
						},
						{
							Format:         wgpu.VertexFormat_Float32x2,
							Offset:         0 + wgpu.VertexFormat_Float32x2.Size(),
							ShaderLocation: 1,
						},
					},
				},
				{
					ArrayStride: 2 * 4,
					StepMode:    wgpu.VertexStepMode_Vertex,
					Attributes: []wgpu.VertexAttribute{
						{
							Format:         wgpu.VertexFormat_Float32x2,
							Offset:         0,
							ShaderLocation: 2,
						},
					},
				},
			},
		},
		Fragment: &wgpu.FragmentState{
			Module:     drawShader,
			EntryPoint: "main_fs",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    config.Format,
					Blend:     nil,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopology_TriangleList,
			FrontFace: wgpu.FrontFace_CCW,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
	})
	if err != nil {
		panic(err)
	}
	defer renderPipeline.Drop()

	computePipeline, err := device.CreateComputePipeline(&wgpu.ComputePipelineDescriptor{
		Label:  "Compute pipeline",
		Layout: computePipelineLayout,
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     computeShader,
			EntryPoint: "main",
		},
	})
	if err != nil {
		panic(err)
	}
	defer computePipeline.Drop()

	vertexBufferData := []float32{-0.01, -0.02, 0.01, -0.02, 0.00, 0.02}
	verticesBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexBufferData),
		Usage:    wgpu.BufferUsage_Vertex | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
	defer verticesBuffer.Drop()

	var initialParticleData [4 * NumParticles]float32
	rng := rand.NewSource(42)

	for i := 0; i < len(initialParticleData); i += 4 {
		initialParticleData[i+0] = float32(rng.Int63())/math.MaxInt64*2 - 1
		initialParticleData[i+1] = float32(rng.Int63())/math.MaxInt64*2 - 1
		initialParticleData[i+2] = (float32(rng.Int63())/math.MaxInt64*2 - 1) * 0.1
		initialParticleData[i+3] = (float32(rng.Int63())/math.MaxInt64*2 - 1) * 0.1
	}

	particleBuffers := []*wgpu.Buffer{}
	particleBindGroups := []*wgpu.BindGroup{}

	for i := 0; i < 2; i++ {
		particleBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
			Label:    "Particle Buffer " + strconv.Itoa(i),
			Contents: wgpu.ToBytes(initialParticleData[:]),
			Usage: wgpu.BufferUsage_Vertex |
				wgpu.BufferUsage_Storage |
				wgpu.BufferUsage_CopyDst,
		})
		if err != nil {
			panic(err)
		}
		defer particleBuffer.Drop()

		particleBuffers = append(particleBuffers, particleBuffer)
	}

	for i := 0; i < 2; i++ {
		particleBindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
			Layout: computeBindGroupLayout,
			Entries: []wgpu.BindGroupEntry{
				{
					Binding: 0,
					Buffer:  simParamBuffer,
				},
				{
					Binding: 1,
					Buffer:  particleBuffers[i],
				},
				{
					Binding: 2,
					Buffer:  particleBuffers[(i+1)%2],
				},
			},
		})
		if err != nil {
			panic(err)
		}
		defer particleBindGroup.Drop()

		particleBindGroups = append(particleBindGroups, particleBindGroup)
	}

	workGroupCount := uint32(math.Ceil(float64(NumParticles) / float64(ParticlesPerGroup)))
	frameNum := uint64(0)

	for !window.ShouldClose() {
		func() {
			var nextTexture *wgpu.TextureView

			for attempt := 0; attempt < 2; attempt++ {
				width, height := window.GetSize()

				if uint32(width) != config.Width || uint32(height) != config.Height {
					config.Width = uint32(width)
					config.Height = uint32(height)

					swapChain, err = device.CreateSwapChain(surface, config)
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

			commandEncoder, err := device.CreateCommandEncoder(nil)
			if err != nil {
				panic(err)
			}

			computePass := commandEncoder.BeginComputePass(nil)
			computePass.SetPipeline(computePipeline)
			computePass.SetBindGroup(0, particleBindGroups[frameNum%2], nil)
			computePass.DispatchWorkgroups(workGroupCount, 1, 1)
			computePass.End()

			renderPass := commandEncoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
				ColorAttachments: []wgpu.RenderPassColorAttachment{
					{
						View:    nextTexture,
						LoadOp:  wgpu.LoadOp_Load,
						StoreOp: wgpu.StoreOp_Store,
					},
				},
			})
			renderPass.SetPipeline(renderPipeline)
			renderPass.SetVertexBuffer(0, particleBuffers[(frameNum+1)%2], 0, 0)
			renderPass.SetVertexBuffer(1, verticesBuffer, 0, 0)
			renderPass.Draw(3, NumParticles, 0, 0)
			renderPass.End()

			frameNum += 1

			queue.Submit(commandEncoder.Finish(nil))
			swapChain.Present()

			glfw.PollEvents()
		}()
	}
}
