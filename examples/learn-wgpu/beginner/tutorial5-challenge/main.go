package main

import (
	_ "embed"
	"fmt"
	"strings"
	"unsafe"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shader.wgsl
var shaderCode string

//go:embed happy-tree.png
var happyTreePng []byte

//go:embed happy-tree-cartoon.png
var happyTreeCartoonPng []byte

type Vertex struct {
	position  [3]float32
	texCoords [2]float32
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Offset:         0,
			ShaderLocation: 0,
			Format:         wgpu.VertexFormat_Float32x3,
		},
		{
			Offset:         uint64(unsafe.Sizeof([3]float32{})),
			ShaderLocation: 1,
			Format:         wgpu.VertexFormat_Float32x2,
		},
	},
}

var VERTICES = [...]Vertex{
	{
		position:  [3]float32{-0.0868241, 0.49240386, 0.0},
		texCoords: [2]float32{0.4131759, 0.00759614},
	}, // A
	{
		position:  [3]float32{-0.49513406, 0.06958647, 0.0},
		texCoords: [2]float32{0.0048659444, 0.43041354},
	}, // B
	{
		position:  [3]float32{-0.21918549, -0.44939706, 0.0},
		texCoords: [2]float32{0.28081453, 0.949397},
	}, // C
	{
		position:  [3]float32{0.35966998, -0.3473291, 0.0},
		texCoords: [2]float32{0.85967, 0.84732914},
	}, // D
	{
		position:  [3]float32{0.44147372, 0.2347359, 0.0},
		texCoords: [2]float32{0.9414737, 0.2652641},
	}, // E
}

var INDICES = [...]uint16{0, 1, 4, 1, 2, 4, 2, 3, 4}

type State struct {
	surface          *wgpu.Surface
	swapChain        *wgpu.SwapChain
	device           *wgpu.Device
	queue            *wgpu.Queue
	config           *wgpu.SwapChainDescriptor
	size             dpi.PhysicalSize[uint32]
	renderPipeline   *wgpu.RenderPipeline
	vertexBuffer     *wgpu.Buffer
	indexBuffer      *wgpu.Buffer
	numIndices       uint32
	diffuseTexture   *Texture
	diffuseBindGroup *wgpu.BindGroup

	cartoonTexture   *Texture
	cartoonBindGroup *wgpu.BindGroup
	isSpacePressed   bool
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

	textureBindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Fragment,
				Texture: wgpu.TextureBindingLayout{
					Multisampled:  false,
					ViewDimension: wgpu.TextureViewDimension_2D,
					SampleType:    wgpu.TextureSampleType_Float,
				},
			},
			{
				Binding:    1,
				Visibility: wgpu.ShaderStage_Fragment,
				Sampler: wgpu.SamplerBindingLayout{
					Type: wgpu.SamplerBindingType_Filtering,
				},
			},
		},
		Label: "TextureBindGroupLayout",
	})
	if err != nil {
		return nil, err
	}

	diffuseTexture, err := TextureFromPNGBytes(device, queue, happyTreePng, "happy-tree.png")
	if err != nil {
		return nil, err
	}

	diffuseBindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: textureBindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: diffuseTexture.view,
			},
			{
				Binding: 1,
				Sampler: diffuseTexture.sampler,
			},
		},
		Label: "DiffuseBindGroup",
	})
	if err != nil {
		return nil, err
	}

	cartoonTexture, err := TextureFromPNGBytes(device, queue, happyTreeCartoonPng, "happy-tree-cartoon.png")
	if err != nil {
		return nil, err
	}

	cartoonBindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: textureBindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding:     0,
				TextureView: cartoonTexture.view,
			},
			{
				Binding: 1,
				Sampler: cartoonTexture.sampler,
			},
		},
		Label: "CartoonBindGroup",
	})
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
		BindGroupLayouts: []*wgpu.BindGroupLayout{
			textureBindGroupLayout,
		},
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
			Buffers:    []wgpu.VertexBufferLayout{VertexBufferLayout},
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

	vertexBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(VERTICES[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return nil, err
	}

	indexBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(INDICES[:]),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		return nil, err
	}
	numIndices := uint32(len(INDICES))

	return &State{
		surface:          surface,
		swapChain:        swapChain,
		device:           device,
		queue:            queue,
		config:           config,
		size:             size,
		renderPipeline:   renderPipeline,
		vertexBuffer:     vertexBuffer,
		indexBuffer:      indexBuffer,
		numIndices:       numIndices,
		diffuseTexture:   diffuseTexture,
		diffuseBindGroup: diffuseBindGroup,
		cartoonTexture:   cartoonTexture,
		cartoonBindGroup: cartoonBindGroup,
		isSpacePressed:   false,
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
	renderPass.SetPipeline(s.renderPipeline)
	if s.isSpacePressed {
		renderPass.SetBindGroup(0, s.cartoonBindGroup, nil)
	} else {
		renderPass.SetBindGroup(0, s.diffuseBindGroup, nil)
	}
	renderPass.SetVertexBuffer(0, s.vertexBuffer, 0, 0)
	renderPass.SetIndexBuffer(s.indexBuffer, wgpu.IndexFormat_Uint16, 0, 0)
	renderPass.DrawIndexed(s.numIndices, 1, 0, 0, 0)
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
			s.isSpacePressed = state == events.ButtonStatePressed
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
