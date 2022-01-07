package wgpu

/*

#include "wrapper.h"

*/
import "C"
import "unsafe"

type Device struct{ ref C.WGPUDevice }

func (p *Device) Poll(forceWait bool) {
	C.wgpuDevicePoll(p.ref, C.bool(forceWait))
}

type BufferBindingLayout struct {
	Type             BufferBindingType
	HasDynamicOffset bool
	MinBindingSize   uint64
}

type SamplerBindingLayout struct {
	Type SamplerBindingType
}

type TextureBindingLayout struct {
	SampleType    TextureSampleType
	ViewDimension TextureViewDimension
	Multisampled  bool
}

type StorageTextureBindingLayout struct {
	Access        StorageTextureAccess
	Format        TextureFormat
	ViewDimension TextureViewDimension
}

type BindGroupLayoutEntry struct {
	Binding        uint32
	Visibility     ShaderStage
	Buffer         BufferBindingLayout
	Sampler        SamplerBindingLayout
	Texture        TextureBindingLayout
	StorageTexture StorageTextureBindingLayout
}

type BindGroupLayoutDescriptor struct {
	Label   string
	Entries []BindGroupLayoutEntry
}

func (p *Device) CreateBindGroupLayout(descriptor BindGroupLayoutDescriptor) BindGroupLayout {
	var desc C.WGPUBindGroupLayoutDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	entryCount := len(descriptor.Entries)
	if entryCount > 0 {
		entries := C.malloc(C.size_t(entryCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupLayoutEntry{})))
		defer C.free(entries)

		entriesSlice := (*[1 << 30]C.WGPUBindGroupLayoutEntry)(entries)[:entryCount:entryCount]

		for i, v := range descriptor.Entries {
			entriesSlice[i] = C.WGPUBindGroupLayoutEntry{
				nextInChain: nil,
				binding:     C.uint32_t(v.Binding),
				visibility:  C.WGPUShaderStageFlags(v.Visibility),
				buffer: C.WGPUBufferBindingLayout{
					nextInChain:      nil,
					_type:            C.WGPUBufferBindingType(v.Buffer.Type),
					hasDynamicOffset: C.bool(v.Buffer.HasDynamicOffset),
					minBindingSize:   C.uint64_t(v.Buffer.MinBindingSize),
				},
				sampler: C.WGPUSamplerBindingLayout{
					nextInChain: nil,
					_type:       C.WGPUSamplerBindingType(v.Sampler.Type),
				},
				texture: C.WGPUTextureBindingLayout{
					nextInChain:   nil,
					sampleType:    C.WGPUTextureSampleType(v.Texture.SampleType),
					viewDimension: C.WGPUTextureViewDimension(v.Texture.ViewDimension),
					multisampled:  C.bool(v.Texture.Multisampled),
				},
				storageTexture: C.WGPUStorageTextureBindingLayout{
					nextInChain:   nil,
					access:        C.WGPUStorageTextureAccess(v.StorageTexture.Access),
					format:        C.WGPUTextureFormat(v.StorageTexture.Format),
					viewDimension: C.WGPUTextureViewDimension(v.StorageTexture.ViewDimension),
				},
			}
		}

		desc.entryCount = C.uint32_t(entryCount)
		desc.entries = (*C.WGPUBindGroupLayoutEntry)(entries)
	}

	return BindGroupLayout(C.wgpuDeviceCreateBindGroupLayout(p.ref, &desc))
}

type BindGroupEntry struct {
	Binding     uint32
	Buffer      *Buffer
	Offset      uint64
	Size        uint64
	Sampler     Sampler
	TextureView TextureView
}

type BindGroupDescriptor struct {
	Label   string
	Layout  BindGroupLayout
	Entries []BindGroupEntry
}

func (p *Device) CreateBindGroup(descriptor BindGroupDescriptor) BindGroup {
	var desc C.WGPUBindGroupDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	desc.layout = C.WGPUBindGroupLayout(descriptor.Layout)

	entryCount := len(descriptor.Entries)
	if entryCount > 0 {
		entries := C.malloc(C.size_t(entryCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupEntry{})))
		defer C.free(entries)

		entriesSlice := (*[1 << 30]C.WGPUBindGroupEntry)(entries)[:entryCount:entryCount]

		for i, v := range descriptor.Entries {
			entry := C.WGPUBindGroupEntry{
				binding:     C.uint32_t(v.Binding),
				offset:      C.uint64_t(v.Offset),
				size:        C.uint64_t(v.Size),
				sampler:     C.WGPUSampler(v.Sampler),
				textureView: C.WGPUTextureView(v.TextureView),
			}

			if v.Buffer != nil {
				entry.buffer = C.WGPUBuffer(v.Buffer.ref)
			}

			entriesSlice[i] = entry
		}

		desc.entryCount = C.uint32_t(entryCount)
		desc.entries = (*C.WGPUBindGroupEntry)(entries)
	}

	return BindGroup(C.wgpuDeviceCreateBindGroup(p.ref, &desc))
}

type BufferDescriptor struct {
	Label            string
	Usage            BufferUsage
	Size             uint64
	MappedAtCreation bool
}

func (p *Device) CreateBuffer(descriptor BufferDescriptor) *Buffer {
	var desc C.WGPUBufferDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	desc.usage = C.WGPUBufferUsageFlags(descriptor.Usage)
	desc.size = C.uint64_t(descriptor.Size)
	desc.mappedAtCreation = C.bool(descriptor.MappedAtCreation)

	ref := C.wgpuDeviceCreateBuffer(p.ref, &desc)
	if ref == nil {
		return nil
	}

	return &Buffer{ref}
}

type CommandEncoderDescriptor struct {
	Label string
}

func (p *Device) CreateCommandEncoder(descriptor CommandEncoderDescriptor) *CommandEncoder {
	var desc C.WGPUCommandEncoderDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuDeviceCreateCommandEncoder(p.ref, &desc)
	if ref == nil {
		return nil
	}

	return &CommandEncoder{ref}
}

type ConstantEntry struct {
	Key   string
	Value float64
}

type ProgrammableStageDescriptor struct {
	Module     ShaderModule
	EntryPoint string

	// unused in wgpu
	// Constants  []ConstantEntry
}

type ComputePipelineDescriptor struct {
	Label   string
	Layout  PipelineLayout
	Compute ProgrammableStageDescriptor
}

func (p *Device) CreateComputePipeline(descriptor ComputePipelineDescriptor) ComputePipeline {
	var desc C.WGPUComputePipelineDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	desc.layout = C.WGPUPipelineLayout(descriptor.Layout)

	compute := C.WGPUProgrammableStageDescriptor{
		module: C.WGPUShaderModule(descriptor.Compute.Module),
	}
	if descriptor.Compute.EntryPoint != "" {
		entryPoint := C.CString(descriptor.Compute.EntryPoint)
		defer C.free(unsafe.Pointer(entryPoint))

		compute.entryPoint = entryPoint
	}
	desc.compute = compute

	return ComputePipeline(C.wgpuDeviceCreateComputePipeline(p.ref, &desc))
}

type PipelineLayoutDescriptor struct {
	Label            string
	BindGroupLayouts []BindGroupLayout
}

func (p *Device) CreatePipelineLayout(descriptor PipelineLayoutDescriptor) PipelineLayout {
	var desc C.WGPUPipelineLayoutDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	bindGroupLayoutCount := len(descriptor.BindGroupLayouts)
	if bindGroupLayoutCount > 0 {
		bindGroupLayouts := C.malloc(C.size_t(bindGroupLayoutCount) * C.size_t(unsafe.Sizeof(C.WGPUBindGroupLayout(nil))))
		defer C.free(bindGroupLayouts)

		bindGroupLayoutsSlice := (*[1 << 30]C.WGPUBindGroupLayout)(bindGroupLayouts)[:bindGroupLayoutCount:bindGroupLayoutCount]

		for i, v := range descriptor.BindGroupLayouts {
			bindGroupLayoutsSlice[i] = C.WGPUBindGroupLayout(v)
		}

		desc.bindGroupLayoutCount = C.uint32_t(bindGroupLayoutCount)
		desc.bindGroupLayouts = (*C.WGPUBindGroupLayout)(bindGroupLayouts)
	}

	return PipelineLayout(C.wgpuDeviceCreatePipelineLayout(p.ref, &desc))
}

type VertexAttribute struct {
	Format         VertexFormat
	Offset         uint64
	ShaderLocation uint32
}

type VertexBufferLayout struct {
	ArrayStride uint64
	StepMode    VertexStepMode
	Attributes  []VertexAttribute
}

type VertexState struct {
	Module     ShaderModule
	EntryPoint string
	Buffers    []VertexBufferLayout

	// unused in wgpu
	// Constants  []ConstantEntry
}

type PrimitiveState struct {
	Topology         PrimitiveTopology
	StripIndexFormat IndexFormat
	FrontFace        FrontFace
	CullMode         CullMode
}

type StencilFaceState struct {
	Compare     CompareFunction
	FailOp      StencilOperation
	DepthFailOp StencilOperation
	PassOp      StencilOperation
}

type DepthStencilState struct {
	Format              TextureFormat
	DepthWriteEnabled   bool
	DepthCompare        CompareFunction
	StencilFront        StencilFaceState
	StencilBack         StencilFaceState
	StencilReadMask     uint32
	StencilWriteMask    uint32
	DepthBias           int32
	DepthBiasSlopeScale float32
	DepthBiasClamp      float32
}

type MultisampleState struct {
	Count                  uint32
	Mask                   uint32
	AlphaToCoverageEnabled bool
}

type BlendComponent struct {
	Operation BlendOperation
	SrcFactor BlendFactor
	DstFactor BlendFactor
}

type BlendState struct {
	Color BlendComponent
	Alpha BlendComponent
}

type ColorTargetState struct {
	Format    TextureFormat
	Blend     *BlendState
	WriteMask ColorWriteMask
}

type FragmentState struct {
	Module     ShaderModule
	EntryPoint string
	Targets    []ColorTargetState

	// unused in wgpu
	// Constants  []ConstantEntry
}

type RenderPipelineDescriptor struct {
	Label        string
	Layout       PipelineLayout
	Vertex       VertexState
	Primitive    PrimitiveState
	DepthStencil *DepthStencilState
	Multisample  MultisampleState
	Fragment     *FragmentState
}

func (p *Device) CreateRenderPipeline(descriptor RenderPipelineDescriptor) RenderPipeline {
	var desc C.WGPURenderPipelineDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	desc.layout = C.WGPUPipelineLayout(descriptor.Layout)

	// vertex
	{
		vertex := descriptor.Vertex

		vert := C.WGPUVertexState{
			module: C.WGPUShaderModule(vertex.Module),
		}

		if vertex.EntryPoint != "" {
			entryPoint := C.CString(vertex.EntryPoint)
			defer C.free(unsafe.Pointer(entryPoint))

			vert.entryPoint = entryPoint
		}

		bufferCount := len(vertex.Buffers)
		if bufferCount > 0 {
			buffers := C.malloc(C.size_t(bufferCount) * C.size_t(unsafe.Sizeof(C.WGPUVertexBufferLayout{})))
			defer C.free(buffers)

			buffersSlice := (*[1 << 30]C.WGPUVertexBufferLayout)(buffers)[:bufferCount:bufferCount]

			for i, v := range vertex.Buffers {
				buffer := C.WGPUVertexBufferLayout{
					arrayStride: C.uint64_t(v.ArrayStride),
					stepMode:    C.WGPUVertexStepMode(v.StepMode),
				}

				attributeCount := len(v.Attributes)
				if attributeCount > 0 {
					attributes := C.malloc(C.size_t(attributeCount) * C.size_t(unsafe.Sizeof(C.WGPUVertexAttribute{})))
					defer C.free(attributes)

					attributesSlice := (*[1 << 30]C.WGPUVertexAttribute)(attributes)[:attributeCount:attributeCount]

					for j, attribute := range v.Attributes {
						attributesSlice[j] = C.WGPUVertexAttribute{
							format:         C.WGPUVertexFormat(attribute.Format),
							offset:         C.uint64_t(attribute.Offset),
							shaderLocation: C.uint32_t(attribute.ShaderLocation),
						}
					}

					buffer.attributeCount = C.uint32_t(attributeCount)
					buffer.attributes = (*C.WGPUVertexAttribute)(attributes)
				}

				buffersSlice[i] = buffer
			}

			vert.bufferCount = C.uint32_t(bufferCount)
			vert.buffers = (*C.WGPUVertexBufferLayout)(buffers)
		}

		desc.vertex = vert
	}

	desc.primitive = C.WGPUPrimitiveState{
		topology:         C.WGPUPrimitiveTopology(descriptor.Primitive.Topology),
		stripIndexFormat: C.WGPUIndexFormat(descriptor.Primitive.StripIndexFormat),
		frontFace:        C.WGPUFrontFace(descriptor.Primitive.FrontFace),
		cullMode:         C.WGPUCullMode(descriptor.Primitive.CullMode),
	}

	if descriptor.DepthStencil != nil {
		depthStencil := descriptor.DepthStencil

		ds := (*C.WGPUDepthStencilState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUDepthStencilState{}))))
		defer C.free(unsafe.Pointer(ds))

		ds.nextInChain = nil
		ds.format = C.WGPUTextureFormat(depthStencil.Format)
		ds.depthWriteEnabled = C.bool(depthStencil.DepthWriteEnabled)
		ds.depthCompare = C.WGPUCompareFunction(depthStencil.DepthCompare)
		ds.stencilFront = C.WGPUStencilFaceState{
			compare:     C.WGPUCompareFunction(depthStencil.StencilFront.Compare),
			failOp:      C.WGPUStencilOperation(depthStencil.StencilFront.FailOp),
			depthFailOp: C.WGPUStencilOperation(depthStencil.StencilFront.DepthFailOp),
			passOp:      C.WGPUStencilOperation(depthStencil.StencilFront.PassOp),
		}
		ds.stencilBack = C.WGPUStencilFaceState{
			compare:     C.WGPUCompareFunction(depthStencil.StencilBack.Compare),
			failOp:      C.WGPUStencilOperation(depthStencil.StencilBack.FailOp),
			depthFailOp: C.WGPUStencilOperation(depthStencil.StencilBack.DepthFailOp),
			passOp:      C.WGPUStencilOperation(depthStencil.StencilBack.PassOp),
		}
		ds.stencilReadMask = C.uint32_t(depthStencil.StencilReadMask)
		ds.stencilWriteMask = C.uint32_t(depthStencil.StencilWriteMask)
		ds.depthBias = C.int32_t(depthStencil.DepthBias)
		ds.depthBiasSlopeScale = C.float(depthStencil.DepthBiasSlopeScale)
		ds.depthBiasClamp = C.float(depthStencil.DepthBiasClamp)

		desc.depthStencil = ds
	}

	desc.multisample = C.WGPUMultisampleState{
		count:                  C.uint32_t(descriptor.Multisample.Count),
		mask:                   C.uint32_t(descriptor.Multisample.Mask),
		alphaToCoverageEnabled: C.bool(descriptor.Multisample.AlphaToCoverageEnabled),
	}

	if descriptor.Fragment != nil {
		fragment := descriptor.Fragment

		frag := (*C.WGPUFragmentState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUFragmentState{}))))
		defer C.free(unsafe.Pointer(frag))

		frag.nextInChain = nil
		if fragment.EntryPoint != "" {
			entryPoint := C.CString(fragment.EntryPoint)
			defer C.free(unsafe.Pointer(entryPoint))

			frag.entryPoint = entryPoint
		}

		frag.module = C.WGPUShaderModule(fragment.Module)

		targetCount := len(fragment.Targets)
		if targetCount > 0 {
			targets := C.malloc(C.size_t(targetCount) * C.size_t(unsafe.Sizeof(C.WGPUColorTargetState{})))
			defer C.free(targets)

			targetsSlice := (*[1 << 30]C.WGPUColorTargetState)(targets)[:targetCount:targetCount]

			for i, v := range fragment.Targets {
				target := C.WGPUColorTargetState{
					format:    C.WGPUTextureFormat(v.Format),
					writeMask: C.WGPUColorWriteMaskFlags(v.WriteMask),
				}

				if v.Blend != nil {
					blend := (*C.WGPUBlendState)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUBlendState{}))))
					defer C.free(unsafe.Pointer(blend))

					blend.color = C.WGPUBlendComponent{
						operation: C.WGPUBlendOperation(v.Blend.Color.Operation),
						srcFactor: C.WGPUBlendFactor(v.Blend.Color.SrcFactor),
						dstFactor: C.WGPUBlendFactor(v.Blend.Color.DstFactor),
					}
					blend.alpha = C.WGPUBlendComponent{
						operation: C.WGPUBlendOperation(v.Blend.Alpha.Operation),
						srcFactor: C.WGPUBlendFactor(v.Blend.Alpha.SrcFactor),
						dstFactor: C.WGPUBlendFactor(v.Blend.Alpha.DstFactor),
					}

					target.blend = blend
				}

				targetsSlice[i] = target
			}

			frag.targetCount = C.uint32_t(targetCount)
			frag.targets = (*C.WGPUColorTargetState)(targets)
		} else {
			frag.targetCount = 0
			frag.targets = nil
		}

		desc.fragment = frag
	}

	return RenderPipeline(C.wgpuDeviceCreateRenderPipeline(p.ref, &desc))
}

type SamplerDescriptor struct {
	Label          string
	AddressModeU   AddressMode
	AddressModeV   AddressMode
	AddressModeW   AddressMode
	MagFilter      FilterMode
	MinFilter      FilterMode
	MipmapFilter   FilterMode
	LodMinClamp    float32
	LodMaxClamp    float32
	Compare        CompareFunction
	MaxAnisotrophy uint16
}

func (p *Device) CreateSampler(descriptor SamplerDescriptor) Sampler {
	desc := C.WGPUSamplerDescriptor{
		addressModeU:  C.WGPUAddressMode(descriptor.AddressModeU),
		addressModeV:  C.WGPUAddressMode(descriptor.AddressModeV),
		addressModeW:  C.WGPUAddressMode(descriptor.AddressModeW),
		magFilter:     C.WGPUFilterMode(descriptor.MagFilter),
		minFilter:     C.WGPUFilterMode(descriptor.MinFilter),
		mipmapFilter:  C.WGPUFilterMode(descriptor.MipmapFilter),
		lodMinClamp:   C.float(descriptor.LodMinClamp),
		lodMaxClamp:   C.float(descriptor.LodMaxClamp),
		compare:       C.WGPUCompareFunction(descriptor.Compare),
		maxAnisotropy: C.uint16_t(descriptor.MaxAnisotrophy),
	}

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	return Sampler(C.wgpuDeviceCreateSampler(p.ref, &desc))
}

type ShaderModuleSPIRVDescriptor struct {
	Code []byte
}

type ShaderModuleWGSLDescriptor struct {
	Code string
}

type ShaderModuleDescriptor struct {
	Label string

	// ChainedStruct -> WGPUShaderModuleSPIRVDescriptor
	SPIRVDescriptor *ShaderModuleSPIRVDescriptor

	// ChainedStruct -> WGPUShaderModuleWGSLDescriptor
	WGSLDescriptor *ShaderModuleWGSLDescriptor
}

func (p *Device) CreateShaderModule(descriptor ShaderModuleDescriptor) ShaderModule {
	var desc C.WGPUShaderModuleDescriptor

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	if descriptor.SPIRVDescriptor != nil {
		spirv := (*C.WGPUShaderModuleSPIRVDescriptor)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderModuleSPIRVDescriptor{}))))
		defer C.free(unsafe.Pointer(spirv))

		codeSize := len(descriptor.SPIRVDescriptor.Code)
		if codeSize > 0 {
			code := C.CBytes(descriptor.SPIRVDescriptor.Code)
			defer C.free(code)

			spirv.codeSize = C.uint32_t(codeSize)
			spirv.code = (*C.uint32_t)(code)
		}
		spirv.chain.next = nil
		spirv.chain.sType = C.WGPUSType_ShaderModuleSPIRVDescriptor

		desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(spirv))
	}

	if descriptor.WGSLDescriptor != nil {
		wgsl := (*C.WGPUShaderModuleWGSLDescriptor)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUShaderModuleWGSLDescriptor{}))))
		defer C.free(unsafe.Pointer(wgsl))

		if descriptor.WGSLDescriptor.Code != "" {
			code := C.CString(descriptor.WGSLDescriptor.Code)
			defer C.free(unsafe.Pointer(code))

			wgsl.code = code
		}
		wgsl.chain.next = nil
		wgsl.chain.sType = C.WGPUSType_ShaderModuleWGSLDescriptor

		desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(wgsl))
	}

	return ShaderModule(C.wgpuDeviceCreateShaderModule(p.ref, &desc))
}

type SwapChainDescriptor struct {
	Usage       TextureUsage
	Format      TextureFormat
	Width       uint32
	Height      uint32
	PresentMode PresentMode

	// Unused in wgpu
	// 	Label       string
}

func (p *Device) CreateSwapChain(surface *Surface, descriptor SwapChainDescriptor) *SwapChain {
	desc := C.WGPUSwapChainDescriptor{
		usage:       C.WGPUTextureUsageFlags(descriptor.Usage),
		format:      C.WGPUTextureFormat(descriptor.Format),
		width:       C.uint32_t(descriptor.Width),
		height:      C.uint32_t(descriptor.Height),
		presentMode: C.WGPUPresentMode(descriptor.PresentMode),
	}

	ref := C.wgpuDeviceCreateSwapChain(p.ref, surface.ref, &desc)
	if ref == nil {
		return nil
	}
	return &SwapChain{ref}
}

type Extent3D struct {
	Width              uint32
	Height             uint32
	DepthOrArrayLayers uint32
}

type TextureDescriptor struct {
	Label         string
	Usage         TextureUsage
	Dimension     TextureDimension
	Size          Extent3D
	Format        TextureFormat
	MipLevelCount uint32
	SampleCount   uint32
}

func (p *Device) CreateTexture(descriptor TextureDescriptor) *Texture {
	desc := C.WGPUTextureDescriptor{
		usage:     C.WGPUTextureUsageFlags(descriptor.Usage),
		dimension: C.WGPUTextureDimension(descriptor.Dimension),
		size: C.WGPUExtent3D{
			width:              C.uint32_t(descriptor.Size.Width),
			height:             C.uint32_t(descriptor.Size.Height),
			depthOrArrayLayers: C.uint32_t(descriptor.Size.DepthOrArrayLayers),
		},
		format:        C.WGPUTextureFormat(descriptor.Format),
		mipLevelCount: C.uint32_t(descriptor.MipLevelCount),
		sampleCount:   C.uint32_t(descriptor.SampleCount),
	}

	if descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer C.free(unsafe.Pointer(label))

		desc.label = label
	}

	ref := C.wgpuDeviceCreateTexture(p.ref, &desc)
	if ref == nil {
		return nil
	}
	return &Texture{ref}
}

func (p *Device) GetLimits() SupportedLimits {
	var limits C.WGPUSupportedLimits

	C.wgpuDeviceGetLimits(p.ref, &limits)

	return SupportedLimits{limitsFromC(limits.limits)}
}

func (p *Device) GetQueue() *Queue {
	ref := C.wgpuDeviceGetQueue(p.ref)
	if ref == nil {
		return nil
	}
	return &Queue{ref}
}
