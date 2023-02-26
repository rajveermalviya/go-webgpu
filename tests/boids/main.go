package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	runtime.LockOSThread()

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

type State struct {
	surface            *wgpu.Surface
	swapChain          *wgpu.SwapChain
	device             *wgpu.Device
	queue              *wgpu.Queue
	config             *wgpu.SwapChainDescriptor
	renderPipeline     *wgpu.RenderPipeline
	computePipeline    *wgpu.ComputePipeline
	vertexBuffer       *wgpu.Buffer
	particleBindGroups []*wgpu.BindGroup
	particleBuffers    []*wgpu.Buffer
	frameNum           uint64
	workGroupCount     uint32
}

func InitState(window *glfw.Window) (s *State, err error) {
	defer func() {
		if err != nil {
			s.Destroy()
			s = nil
		}
	}()
	s = &State{}

	instance := wgpu.CreateInstance(nil)
	defer instance.Drop()

	s.surface = instance.CreateSurface(getSurfaceDescriptor(window))

	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
		CompatibleSurface:    s.surface,
	})
	if err != nil {
		return s, err
	}
	defer adapter.Drop()

	s.device, err = adapter.RequestDevice(nil)
	if err != nil {
		return s, err
	}
	s.queue = s.device.GetQueue()

	width, height := window.GetSize()
	s.config = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      s.surface.GetPreferredFormat(adapter),
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
	}

	s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
	if err != nil {
		return s, err
	}

	computeShader, err := s.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "compute.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: compute,
		},
	})
	if err != nil {
		return s, err
	}
	defer computeShader.Drop()

	drawShader, err := s.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "draw.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: draw,
		},
	})
	if err != nil {
		return s, err
	}
	defer drawShader.Drop()

	simParamData := [...]float32{
		0.04,  // deltaT
		0.1,   // rule1Distance
		0.025, // rule2Distance
		0.025, // rule3Distance
		0.02,  // rule1Scale
		0.05,  // rule2Scale
		0.005, // rule3Scale
	}

	simParamBuffer, err := s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Simulation Param Buffer",
		Contents: wgpu.ToBytes(simParamData[:]),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return s, err
	}
	defer simParamBuffer.Drop()

	s.renderPipeline, err = s.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
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
					Format:    s.config.Format,
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
		return s, err
	}

	s.computePipeline, err = s.device.CreateComputePipeline(&wgpu.ComputePipelineDescriptor{
		Label: "Compute pipeline",
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     computeShader,
			EntryPoint: "main",
		},
	})
	if err != nil {
		return s, err
	}

	vertexBufferData := [...]float32{-0.01, -0.02, 0.01, -0.02, 0.00, 0.02}
	s.vertexBuffer, err = s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexBufferData[:]),
		Usage:    wgpu.BufferUsage_Vertex | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return s, err
	}

	var initialParticleData [4 * NumParticles]float32
	rng := rand.NewSource(42)

	for i := 0; i < len(initialParticleData); i += 4 {
		initialParticleData[i+0] = float32(rng.Int63())/math.MaxInt64*2 - 1
		initialParticleData[i+1] = float32(rng.Int63())/math.MaxInt64*2 - 1
		initialParticleData[i+2] = (float32(rng.Int63())/math.MaxInt64*2 - 1) * 0.1
		initialParticleData[i+3] = (float32(rng.Int63())/math.MaxInt64*2 - 1) * 0.1
	}

	for i := 0; i < 2; i++ {
		particleBuffer, err := s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
			Label:    "Particle Buffer " + strconv.Itoa(i),
			Contents: wgpu.ToBytes(initialParticleData[:]),
			Usage: wgpu.BufferUsage_Vertex |
				wgpu.BufferUsage_Storage |
				wgpu.BufferUsage_CopyDst,
		})
		if err != nil {
			return s, err
		}

		s.particleBuffers = append(s.particleBuffers, particleBuffer)
	}

	computeBindGroupLayout := s.computePipeline.GetBindGroupLayout(0)
	defer computeBindGroupLayout.Drop()

	for i := 0; i < 2; i++ {
		particleBindGroup, err := s.device.CreateBindGroup(&wgpu.BindGroupDescriptor{
			Layout: computeBindGroupLayout,
			Entries: []wgpu.BindGroupEntry{
				{
					Binding: 0,
					Buffer:  simParamBuffer,
					Size:    wgpu.WholeSize,
				},
				{
					Binding: 1,
					Buffer:  s.particleBuffers[i],
					Size:    wgpu.WholeSize,
				},
				{
					Binding: 2,
					Buffer:  s.particleBuffers[(i+1)%2],
					Size:    wgpu.WholeSize,
				},
			},
		})
		if err != nil {
			return s, err
		}

		s.particleBindGroups = append(s.particleBindGroups, particleBindGroup)
	}

	s.workGroupCount = uint32(math.Ceil(float64(NumParticles) / float64(ParticlesPerGroup)))
	s.frameNum = uint64(0)

	return s, nil
}

func (s *State) Resize(width, height int) {
	if width > 0 && height > 0 {
		s.config.Width = uint32(width)
		s.config.Height = uint32(height)

		if s.swapChain != nil {
			s.swapChain.Drop()
		}
		var err error
		s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
		if err != nil {
			panic(err)
		}
	}
}

func (s *State) Render() error {
	nextTexture, err := s.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer nextTexture.Drop()

	commandEncoder, err := s.device.CreateCommandEncoder(nil)
	if err != nil {
		return err
	}

	computePass := commandEncoder.BeginComputePass(nil)
	computePass.SetPipeline(s.computePipeline)
	computePass.SetBindGroup(0, s.particleBindGroups[s.frameNum%2], nil)
	computePass.DispatchWorkgroups(s.workGroupCount, 1, 1)
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
	renderPass.SetPipeline(s.renderPipeline)
	renderPass.SetVertexBuffer(0, s.particleBuffers[(s.frameNum+1)%2], 0, wgpu.WholeSize)
	renderPass.SetVertexBuffer(1, s.vertexBuffer, 0, wgpu.WholeSize)
	renderPass.Draw(3, NumParticles, 0, 0)
	renderPass.End()

	s.frameNum += 1

	s.queue.Submit(commandEncoder.Finish(nil))
	s.swapChain.Present()

	return nil
}

func (s *State) Destroy() {
	if s.particleBindGroups != nil {
		for _, bg := range s.particleBindGroups {
			bg.Drop()
		}
		s.particleBindGroups = nil
	}
	if s.particleBuffers != nil {
		for _, buffers := range s.particleBuffers {
			buffers.Drop()
		}
		s.particleBuffers = nil
	}
	if s.vertexBuffer != nil {
		s.vertexBuffer.Drop()
		s.vertexBuffer = nil
	}
	if s.computePipeline != nil {
		s.computePipeline.Drop()
		s.computePipeline = nil
	}
	if s.renderPipeline != nil {
		s.renderPipeline.Drop()
		s.renderPipeline = nil
	}
	if s.swapChain != nil {
		s.swapChain.Drop()
		s.swapChain = nil
	}
	if s.config != nil {
		s.config = nil
	}
	if s.queue != nil {
		s.queue = nil
	}
	if s.device != nil {
		s.device.Drop()
		s.device = nil
	}
	if s.surface != nil {
		s.surface.Drop()
		s.surface = nil
	}
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

	s, err := InitState(window)
	if err != nil {
		panic(err)
	}
	defer s.Destroy()

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		s.Resize(width, height)
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		err := s.Render()
		if err != nil {
			fmt.Println(err)

			errstr := err.Error()
			switch {
			case strings.Contains(errstr, "Lost"):
				s.Resize(window.GetSize())
			case strings.Contains(errstr, "Outdated"):
				s.Resize(window.GetSize())
			case strings.Contains(errstr, "Timeout"):
			default:
				panic(err)
			}
		}
	}
}
