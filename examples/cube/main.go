package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/examples/internal/glm"
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

const vertexSize = uint64(unsafe.Sizeof(Vertex{}))

type Vertex struct {
	pos      [4]float32
	texCoord [2]float32
}

func vertex(pos1, pos2, pos3, tc1, tc2 float32) Vertex {
	return Vertex{
		pos:      [4]float32{pos1, pos2, pos3, 1},
		texCoord: [2]float32{tc1, tc2},
	}
}

var vertexData = []Vertex{
	// top (0, 0, 1)
	vertex(-1, -1, 1, 0, 0),
	vertex(1, -1, 1, 1, 0),
	vertex(1, 1, 1, 1, 1),
	vertex(-1, 1, 1, 0, 1),
	// bottom (0, 0, -1)
	vertex(-1, 1, -1, 1, 0),
	vertex(1, 1, -1, 0, 0),
	vertex(1, -1, -1, 0, 1),
	vertex(-1, -1, -1, 1, 1),
	// right (1, 0, 0)
	vertex(1, -1, -1, 0, 0),
	vertex(1, 1, -1, 1, 0),
	vertex(1, 1, 1, 1, 1),
	vertex(1, -1, 1, 0, 1),
	// left (-1, 0, 0)
	vertex(-1, -1, 1, 1, 0),
	vertex(-1, 1, 1, 0, 0),
	vertex(-1, 1, -1, 0, 1),
	vertex(-1, -1, -1, 1, 1),
	// front (0, 1, 0)
	vertex(1, 1, -1, 1, 0),
	vertex(-1, 1, -1, 0, 0),
	vertex(-1, 1, 1, 0, 1),
	vertex(1, 1, 1, 1, 1),
	// back (0, -1, 0)
	vertex(1, -1, 1, 0, 0),
	vertex(-1, -1, 1, 1, 0),
	vertex(-1, -1, -1, 1, 1),
	vertex(1, -1, -1, 0, 1),
}

var indexData = []uint16{
	0, 1, 2, 2, 3, 0, // top
	4, 5, 6, 6, 7, 4, // bottom
	8, 9, 10, 10, 11, 8, // right
	12, 13, 14, 14, 15, 12, // left
	16, 17, 18, 18, 19, 16, // front
	20, 21, 22, 22, 23, 20, // back
}

func createTexels(size int) []uint8 {
	arr := make([]uint8, size*size)

	for id := 0; id < (size * size); id++ {
		cx := 3.0*float32(id%size)/float32(size-1) - 2.0
		cy := 2.0*float32(id/size)/float32(size-1) - 1.0
		x, y, count := float32(cx), float32(cy), uint8(0)
		for count < 0xFF && x*x+y*y < 4.0 {
			oldX := x
			x = x*x - y*y + cx
			y = 2.0*oldX*y + cy
			count += 1
		}
		arr[id] = count
	}

	return arr
}

func generateMatrix(aspectRatio float32) glm.Mat4[float32] {
	projection := glm.PerspectiveRH(math.Pi/4, aspectRatio, 1, 10)
	view := glm.LookAtRH(
		glm.Vec3[float32]{1.5, -5, 3},
		glm.Vec3[float32]{0, 0, 0},
		glm.Vec3[float32]{0, 0, 1},
	)

	return projection.Mul4(view)
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

	device, err := adapter.RequestDevice(&wgpu.DeviceDescriptor{
		Label: "Device",
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

	vertexBuf, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexData),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		panic(err)
	}
	defer vertexBuf.Drop()

	indexBuf, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(indexData),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		panic(err)
	}
	defer indexBuf.Drop()

	bindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex,
				Buffer: wgpu.BufferBindingLayout{
					Type:             wgpu.BufferBindingType_Uniform,
					HasDynamicOffset: false,
					MinBindingSize:   64,
				},
			},
			{
				Binding:    1,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					SampleType:    wgpu.TextureSampleType_Uint,
					ViewDimension: wgpu.TextureViewDimension_2D,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	defer bindGroupLayout.Drop()

	pipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		BindGroupLayouts: []*wgpu.BindGroupLayout{bindGroupLayout},
	})
	if err != nil {
		panic(err)
	}
	defer pipelineLayout.Drop()

	size := 256
	texels := createTexels(size)
	textureExtent := wgpu.Extent3D{
		Width:              uint32(size),
		Height:             uint32(size),
		DepthOrArrayLayers: 1,
	}
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_R8Uint,
		Usage:         wgpu.TextureUsage_TextureBinding | wgpu.TextureUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
	defer texture.Drop()

	textureView := texture.CreateView(nil)
	defer textureView.Drop()

	queue.WriteTexture(
		texture.AsImageCopy(),
		wgpu.ToBytes(texels),
		&wgpu.TextureDataLayout{
			Offset:       0,
			BytesPerRow:  uint32(size),
			RowsPerImage: 0,
		},
		&textureExtent,
	)

	mxTotal := generateMatrix(float32(config.Width) / float32(config.Height))
	uniformBuf, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Uniform Buffer",
		Contents: wgpu.ToBytes(mxTotal[:]),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}
	defer uniformBuf.Drop()

	bindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  uniformBuf,
				Offset:  0,
				Size:    0,
			},
			{
				Binding:     1,
				TextureView: textureView,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	defer bindGroup.Drop()

	shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
	})
	if err != nil {
		panic(err)
	}
	defer shader.Drop()

	pipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{{
				ArrayStride: vertexSize,
				StepMode:    wgpu.VertexStepMode_Vertex,
				Attributes: []wgpu.VertexAttribute{
					{
						Format:         wgpu.VertexFormat_Float32x4,
						Offset:         0,
						ShaderLocation: 0,
					},
					{
						Format:         wgpu.VertexFormat_Float32x2,
						Offset:         4 * 4,
						ShaderLocation: 1,
					},
				},
			}},
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
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
			CullMode:  wgpu.CullMode_Back,
		},
		DepthStencil: nil,
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
	})
	if err != nil {
		panic(err)
	}
	defer pipeline.Drop()

	for !window.ShouldClose() {
		func() {
			var nextTexture *wgpu.TextureView

			for attempt := 0; attempt < 2; attempt++ {
				width, height := window.GetSize()

				if width != int(config.Width) || height != int(config.Height) {
					config.Width = uint32(width)
					config.Height = uint32(height)

					mxTotal := generateMatrix(float32(config.Width) / float32(config.Height))
					queue.WriteBuffer(uniformBuf, 0, wgpu.ToBytes(mxTotal[:]))

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

			encoder, err := device.CreateCommandEncoder(nil)
			if err != nil {
				panic(err)
			}

			renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
				ColorAttachments: []wgpu.RenderPassColorAttachment{
					{
						View:       nextTexture,
						LoadOp:     wgpu.LoadOp_Clear,
						StoreOp:    wgpu.StoreOp_Store,
						ClearValue: wgpu.Color{R: 0.1, G: 0.2, B: 0.3, A: 1.0},
					},
				},
			})

			renderPass.SetPipeline(pipeline)
			renderPass.SetBindGroup(0, bindGroup, nil)
			renderPass.SetIndexBuffer(indexBuf, wgpu.IndexFormat_Uint16, 0, 0)
			renderPass.SetVertexBuffer(0, vertexBuf, 0, 0)
			renderPass.DrawIndexed(uint32(len(indexData)), 1, 0, 0, 0)
			renderPass.End()

			queue.Submit(encoder.Finish(nil))
			swapChain.Present()

			glfw.PollEvents()
		}()
	}
}
