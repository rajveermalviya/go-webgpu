package wgpu

import "unsafe"

type AdapterExtras struct {
	BackendType BackendType
}

type RequestAdapterOptions struct {
	CompatibleSurface    *Surface
	PowerPreference      PowerPreference
	ForceFallbackAdapter bool

	// ChainedStruct -> WGPUAdapterExtras
	AdapterExtras *AdapterExtras
}

type SurfaceDescriptorFromWindowsHWND struct {
	Hinstance unsafe.Pointer
	Hwnd      unsafe.Pointer
}

type SurfaceDescriptorFromXcbWindow struct {
	Connection unsafe.Pointer
	Window     uint32
}

type SurfaceDescriptorFromXlibWindow struct {
	Display unsafe.Pointer
	Window  uint32
}

type SurfaceDescriptorFromMetalLayer struct {
	Layer unsafe.Pointer
}

type SurfaceDescriptorFromWaylandSurface struct {
	Display unsafe.Pointer
	Surface unsafe.Pointer
}

type SurfaceDescriptorFromAndroidNativeWindow struct {
	Window unsafe.Pointer
}

type SurfaceDescriptor struct {
	Label string

	// ChainedStruct -> WGPUSurfaceDescriptorFromWindowsHWND
	WindowsHWND *SurfaceDescriptorFromWindowsHWND

	// ChainedStruct -> WGPUSurfaceDescriptorFromXcbWindow
	XcbWindow *SurfaceDescriptorFromXcbWindow

	// ChainedStruct -> WGPUSurfaceDescriptorFromXlibWindow
	XlibWindow *SurfaceDescriptorFromXlibWindow

	// ChainedStruct -> WGPUSurfaceDescriptorFromMetalLayer
	MetalLayer *SurfaceDescriptorFromMetalLayer

	// ChainedStruct -> WGPUSurfaceDescriptorFromWaylandSurface
	WaylandSurface *SurfaceDescriptorFromWaylandSurface

	// ChainedStruct -> WGPUSurfaceDescriptorFromAndroidNativeWindow
	AndroidNativeWindow *SurfaceDescriptorFromAndroidNativeWindow
}

type Limits struct {
	MaxTextureDimension1D                     uint32
	MaxTextureDimension2D                     uint32
	MaxTextureDimension3D                     uint32
	MaxTextureArrayLayers                     uint32
	MaxBindGroups                             uint32
	MaxDynamicUniformBuffersPerPipelineLayout uint32
	MaxDynamicStorageBuffersPerPipelineLayout uint32
	MaxSampledTexturesPerShaderStage          uint32
	MaxSamplersPerShaderStage                 uint32
	MaxStorageBuffersPerShaderStage           uint32
	MaxStorageTexturesPerShaderStage          uint32
	MaxUniformBuffersPerShaderStage           uint32
	MaxUniformBufferBindingSize               uint64
	MaxStorageBufferBindingSize               uint64
	MinUniformBufferOffsetAlignment           uint32
	MinStorageBufferOffsetAlignment           uint32
	MaxVertexBuffers                          uint32
	MaxVertexAttributes                       uint32
	MaxVertexBufferArrayStride                uint32
	MaxInterStageShaderComponents             uint32
	MaxComputeWorkgroupStorageSize            uint32
	MaxComputeInvocationsPerWorkgroup         uint32
	MaxComputeWorkgroupSizeX                  uint32
	MaxComputeWorkgroupSizeY                  uint32
	MaxComputeWorkgroupSizeZ                  uint32
	MaxComputeWorkgroupsPerDimension          uint32
	MaxPushConstantSize                       uint32
	MaxBufferSize                             uint64
}

type SupportedLimits struct {
	Limits Limits
}

type AdapterProperties struct {
	VendorID          uint32
	DeviceID          uint32
	Name              string
	DriverDescription string
	AdapterType       AdapterType
	BackendType       BackendType
}

type DeviceExtras struct {
	TracePath string
}

type RequiredLimits struct {
	Limits Limits
}

type DeviceDescriptor struct {
	Label            string
	RequiredFeatures []FeatureName
	RequiredLimits   *RequiredLimits

	// WGPUChainedStruct -> WGPUDeviceExtras
	DeviceExtras *DeviceExtras
}

type ComputePassTimestampWrite struct {
	QuerySet   QuerySet
	QueryIndex uint32
	Location   ComputePassTimestampLocation
}

type ComputePassDescriptor struct {
	Label string

	// unused in wgpu
	// TimestampWrites []ComputePassTimestampWrite
}

type Color struct {
	R, G, B, A float64
}

type RenderPassColorAttachment struct {
	View          *TextureView
	ResolveTarget *TextureView
	LoadOp        LoadOp
	StoreOp       StoreOp
	ClearValue    Color
}

type RenderPassDepthStencilAttachment struct {
	View              *TextureView
	DepthLoadOp       LoadOp
	DepthStoreOp      StoreOp
	DepthClearValue   float32
	DepthReadOnly     bool
	StencilLoadOp     LoadOp
	StencilStoreOp    StoreOp
	StencilClearValue uint32
	StencilReadOnly   bool
}

type RenderPassTimestampWrite struct {
	QuerySet   QuerySet
	QueryIndex uint32
	Location   RenderPassTimestampLocation
}

type RenderPassDescriptor struct {
	Label                  string
	ColorAttachments       []RenderPassColorAttachment
	DepthStencilAttachment *RenderPassDepthStencilAttachment

	// unused in wgpu
	// 	OcclusionQuerySet      QuerySet
	// 	TimestampWrites        []RenderPassTimestampWrite
}

type TextureDataLayout struct {
	Offset       uint64
	BytesPerRow  uint32
	RowsPerImage uint32
}

type ImageCopyBuffer struct {
	Layout TextureDataLayout
	Buffer *Buffer
}

type Origin3D struct {
	X, Y, Z uint32
}

type ImageCopyTexture struct {
	Texture  *Texture
	MipLevel uint32
	Origin   Origin3D
	Aspect   TextureAspect
}

type CommandBufferDescriptor struct {
	Label string
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

type BindGroupEntry struct {
	Binding     uint32
	Buffer      *Buffer
	Offset      uint64
	Size        uint64
	Sampler     *Sampler
	TextureView *TextureView
}

type BindGroupDescriptor struct {
	Label   string
	Layout  *BindGroupLayout
	Entries []BindGroupEntry
}

type BufferDescriptor struct {
	Label            string
	Usage            BufferUsage
	Size             uint64
	MappedAtCreation bool
}

type CommandEncoderDescriptor struct {
	Label string
}

type ConstantEntry struct {
	Key   string
	Value float64
}

type ProgrammableStageDescriptor struct {
	Module     *ShaderModule
	EntryPoint string

	// unused in wgpu
	// Constants  []ConstantEntry
}

type ComputePipelineDescriptor struct {
	Label   string
	Layout  *PipelineLayout
	Compute ProgrammableStageDescriptor
}

type PushConstantRange struct {
	Stages ShaderStage
	Start  uint32
	End    uint32
}

type PipelineLayoutExtras struct {
	PushConstantRanges []PushConstantRange
}

type PipelineLayoutDescriptor struct {
	Label            string
	BindGroupLayouts []*BindGroupLayout

	// WGPUChainedStruct -> WGPUPipelineLayoutExtras
	PipelineLayoutExtras *PipelineLayoutExtras
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
	Module     *ShaderModule
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
	Module     *ShaderModule
	EntryPoint string
	Targets    []ColorTargetState

	// unused in wgpu
	// Constants  []ConstantEntry
}

type RenderPipelineDescriptor struct {
	Label        string
	Layout       *PipelineLayout
	Vertex       VertexState
	Primitive    PrimitiveState
	DepthStencil *DepthStencilState
	Multisample  MultisampleState
	Fragment     *FragmentState
}

type SamplerDescriptor struct {
	Label          string
	AddressModeU   AddressMode
	AddressModeV   AddressMode
	AddressModeW   AddressMode
	MagFilter      FilterMode
	MinFilter      FilterMode
	MipmapFilter   MipmapFilterMode
	LodMinClamp    float32
	LodMaxClamp    float32
	Compare        CompareFunction
	MaxAnisotrophy uint16
}

type ShaderModuleSPIRVDescriptor struct {
	Code []byte
}

type ShaderModuleWGSLDescriptor struct {
	Code string
}

type ShaderModuleGLSLDescriptor struct {
	Code        string
	Defines     map[string]string
	ShaderStage ShaderStage
}

type ShaderModuleDescriptor struct {
	Label string

	// ChainedStruct -> WGPUShaderModuleSPIRVDescriptor
	SPIRVDescriptor *ShaderModuleSPIRVDescriptor

	// ChainedStruct -> WGPUShaderModuleWGSLDescriptor
	WGSLDescriptor *ShaderModuleWGSLDescriptor

	// ChainedStruct -> WGPUShaderModuleGLSLDescriptor
	GLSLDescriptor *ShaderModuleGLSLDescriptor
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

type TextureViewDescriptor struct {
	Label           string
	Format          TextureFormat
	Dimension       TextureViewDimension
	BaseMipLevel    uint32
	MipLevelCount   uint32
	BaseArrayLayer  uint32
	ArrayLayerCount uint32
	Aspect          TextureAspect
}

type SubmissionIndex uint64

type WrappedSubmissionIndex struct {
	Queue           *Queue
	SubmissionIndex SubmissionIndex
}

type RenderBundleEncoderDescriptor struct {
	Label              string
	ColorFormats       []TextureFormat
	DepthStencilFormat TextureFormat
	SampleCount        uint32
	DepthReadOnly      bool
	StencilReadOnly    bool
}

type RenderBundleDescriptor struct {
	Label string
}

type StorageReport struct {
	NumOccupied uint64
	NumVacant   uint64
	NumError    uint64
	ElementSize uint64
}

type HubReport struct {
	Adapters         StorageReport
	Devices          StorageReport
	PipelineLayouts  StorageReport
	ShaderModules    StorageReport
	BindGroupLayouts StorageReport
	BindGroups       StorageReport
	CommandBuffers   StorageReport
	RenderBundles    StorageReport
	RenderPipelines  StorageReport
	ComputePipelines StorageReport
	QuerySets        StorageReport
	Buffers          StorageReport
	Textures         StorageReport
	TextureViews     StorageReport
	Samplers         StorageReport
}

type GlobalReport struct {
	Surfaces StorageReport
	Vulkan   *HubReport
	Metal    *HubReport
	Dx12     *HubReport
	Dx11     *HubReport
	Gl       *HubReport
}
