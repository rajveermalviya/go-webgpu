package main

import (
	_ "embed"
	"fmt"
	"strings"
	"unsafe"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/go-webgpu/examples/internal/glm"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

//go:embed shader.wgsl
var shaderCode string

//go:embed happy-tree.png
var happyTreePng []byte

const NumInstancesPerRow = 10

var InstanceDisplacement = glm.Vec3[float32]{
	NumInstancesPerRow * 0.5,
	0.0,
	NumInstancesPerRow * 0.5,
}

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

var OpenGlToWgpuMatrix = glm.Mat4[float32]{
	1.0, 0.0, 0.0, 0.0,
	0.0, 1.0, 0.0, 0.0,
	0.0, 0.0, 0.5, 0.0,
	0.0, 0.0, 0.5, 1.0,
}

type Camera struct {
	eye     glm.Vec3[float32]
	target  glm.Vec3[float32]
	up      glm.Vec3[float32]
	aspect  float32
	fovYRad float32
	znear   float32
	zfar    float32
}

func (c *Camera) buildViewProjectionMatrix() glm.Mat4[float32] {
	view := glm.LookAtRH(c.eye, c.target, c.up)
	proj := glm.Perspective(c.fovYRad, c.aspect, c.znear, c.zfar)
	return proj.Mul4(view)
}

type CameraUniform struct {
	viewProj glm.Mat4[float32]
}

func NewCameraUnifrom() *CameraUniform {
	return &CameraUniform{
		viewProj: glm.Mat4[float32]{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

func (c *CameraUniform) UpdateViewProj(camera *Camera) {
	c.viewProj = OpenGlToWgpuMatrix.Mul4(camera.buildViewProjectionMatrix())
}

type CameraController struct {
	speed             float32
	isUpPressed       bool
	isDownPressed     bool
	isForwardPressed  bool
	isBackwardPressed bool
	isLeftPressed     bool
	isRightPressed    bool
}

func NewCameraController(speed float32) *CameraController {
	return &CameraController{speed: speed}
}

func (c *CameraController) UpdateCamera(camera *Camera) {
	forward := camera.target.Sub(camera.eye).Normalize()

	if c.isForwardPressed {
		camera.eye = camera.eye.Add(forward.MulScalar(c.speed))
	}
	if c.isBackwardPressed {
		camera.eye = camera.eye.Sub(forward.MulScalar(c.speed))
	}

	right := forward.Cross(camera.up)

	if c.isRightPressed {
		camera.eye = camera.eye.Add(right.MulScalar(c.speed))
	}
	if c.isLeftPressed {
		camera.eye = camera.eye.Sub(right.MulScalar(c.speed))
	}
}

type Instance struct {
	position glm.Vec3[float32]
	rotation glm.Quaternion[float32]
}

func (i Instance) ToRaw() InstanceRaw {
	return InstanceRaw{
		model: glm.Mat4FromTranslation(i.position).Mul4(glm.Mat4FromQuaternion(i.rotation)),
	}
}

type InstanceRaw struct {
	model glm.Mat4[float32]
}

var InstanceBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(InstanceRaw{})),
	StepMode:    wgpu.VertexStepMode_Instance,
	Attributes: []wgpu.VertexAttribute{
		{
			Offset:         0,
			ShaderLocation: 5,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([4]float32{})),
			ShaderLocation: 6,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([8]float32{})),
			ShaderLocation: 7,
			Format:         wgpu.VertexFormat_Float32x4,
		},
		{
			Offset:         uint64(unsafe.Sizeof([12]float32{})),
			ShaderLocation: 8,
			Format:         wgpu.VertexFormat_Float32x4,
		},
	},
}

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
	camera           *Camera
	cameraController *CameraController
	cameraUniform    *CameraUniform
	cameraBuffer     *wgpu.Buffer
	cameraBindGroup  *wgpu.BindGroup

	instances      [NumInstancesPerRow * NumInstancesPerRow]Instance
	instanceBuffer *wgpu.Buffer
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

	diffuseTexture, err := TextureFromPNGBytes(device, queue, happyTreePng, "happy-tree.png")
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

	camera := &Camera{
		eye:     glm.Vec3[float32]{0, 5, 10},
		target:  glm.Vec3[float32]{0, 0, 0},
		up:      glm.Vec3[float32]{0, 1, 0},
		aspect:  float32(size.Width) / float32(size.Height),
		fovYRad: glm.DegToRad[float32](45),
		znear:   0.1,
		zfar:    100.0,
	}
	cameraController := NewCameraController(0.2)
	cameraUniform := NewCameraUnifrom()
	cameraUniform.UpdateViewProj(camera)

	cameraBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Camera Buffer",
		Contents: wgpu.ToBytes(cameraUniform.viewProj[:]),
		Usage:    wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	var instances [NumInstancesPerRow * NumInstancesPerRow]Instance
	{
		index := 0
		for z := 0; z < NumInstancesPerRow; z++ {
			for x := 0; x < NumInstancesPerRow; x++ {
				position := glm.Vec3[float32]{float32(x), 0, float32(z)}.Sub(InstanceDisplacement)

				var rotation glm.Quaternion[float32]
				if position == (glm.Vec3[float32]{}) {
					rotation = glm.QuaternionFromAxisAngle(glm.Vec3[float32]{0, 0, 1}, 0)
				} else {
					rotation = glm.QuaternionFromAxisAngle(position.Normalize(), glm.DegToRad[float32](45))
				}

				instances[index] = Instance{
					position: position,
					rotation: rotation,
				}
				index++
			}
		}
	}

	var instanceData [NumInstancesPerRow * NumInstancesPerRow]InstanceRaw
	for i, v := range instances {
		instanceData[i] = v.ToRaw()
	}
	instanceBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Instance Buffer",
		Contents: wgpu.ToBytes(instanceData[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return nil, err
	}

	cameraBindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "CameraBindGroupLayout",
		Entries: []wgpu.BindGroupLayoutEntry{{
			Binding:    0,
			Visibility: wgpu.ShaderStage_Vertex,
			Buffer: wgpu.BufferBindingLayout{
				Type:             wgpu.BufferBindingType_Uniform,
				HasDynamicOffset: false,
				MinBindingSize:   0,
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	cameraBindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  "CameraBindGroup",
		Layout: cameraBindGroupLayout,
		Entries: []wgpu.BindGroupEntry{{
			Binding: 0,
			Buffer:  cameraBuffer,
		}},
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
			textureBindGroupLayout, cameraBindGroupLayout,
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
			Buffers:    []wgpu.VertexBufferLayout{VertexBufferLayout, InstanceBufferLayout},
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
		camera:           camera,
		cameraController: cameraController,
		cameraUniform:    cameraUniform,
		cameraBuffer:     cameraBuffer,
		cameraBindGroup:  cameraBindGroup,
		instances:        instances,
		instanceBuffer:   instanceBuffer,
	}, nil
}

func (s *State) Update() {
	s.cameraController.UpdateCamera(s.camera)
	s.cameraUniform.UpdateViewProj(s.camera)
	s.queue.WriteBuffer(s.cameraBuffer, 0, wgpu.ToBytes(s.cameraUniform.viewProj[:]))
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

		s.camera.aspect = float32(newSize.Width) / float32(newSize.Height)
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
	renderPass.SetBindGroup(0, s.diffuseBindGroup, nil)
	renderPass.SetBindGroup(1, s.cameraBindGroup, nil)
	renderPass.SetVertexBuffer(0, s.vertexBuffer, 0, 0)
	renderPass.SetVertexBuffer(1, s.instanceBuffer, 0, 0)
	renderPass.SetIndexBuffer(s.indexBuffer, wgpu.IndexFormat_Uint16, 0, 0)
	renderPass.DrawIndexed(s.numIndices, uint32(len(s.instances)), 0, 0, 0)
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

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		s.Resize(dpi.PhysicalSize[uint32]{
			Width:  physicalWidth,
			Height: physicalHeight,
		})
	})

	w.SetKeyboardInputCallback(func(state events.ButtonState, scanCode events.ScanCode, virtualKeyCode events.VirtualKey) {
		isPressed := state == events.ButtonStatePressed

		switch virtualKeyCode {
		case events.VirtualKeySpace:
			s.cameraController.isUpPressed = isPressed
		case events.VirtualKeyLShift:
			s.cameraController.isDownPressed = isPressed
		case events.VirtualKeyW, events.VirtualKeyUp:
			s.cameraController.isForwardPressed = isPressed
		case events.VirtualKeyA, events.VirtualKeyLeft:
			s.cameraController.isLeftPressed = isPressed
		case events.VirtualKeyS, events.VirtualKeyDown:
			s.cameraController.isBackwardPressed = isPressed
		case events.VirtualKeyD, events.VirtualKeyRight:
			s.cameraController.isRightPressed = isPressed
		}
	})

	w.SetCloseRequestedCallback(func() {
		d.Destroy()
	})

	for {
		if !d.Poll() {
			break
		}

		s.Update()
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
