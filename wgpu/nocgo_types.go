//go:build windows

package wgpu

import (
	"sync/atomic"
	"unsafe"

	"golang.org/x/exp/constraints"
	"golang.org/x/sys/windows"
)

// https://github.com/rust-lang/rust/blob/d40f24e956a698e47a209541031c4045acc5a684/library/std/src/sys/windows/alloc.rs

var (
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	_GetProcessHeap = kernel32.NewProc("GetProcessHeap")
	_HeapFree       = kernel32.NewProc("HeapFree")
)

var processHeapHandle uintptr

func getProcessHeap() uintptr {
	handle := atomic.LoadUintptr(&processHeapHandle)
	if handle == 0 {
		r, _, _ := _GetProcessHeap.Call()
		atomic.StoreUintptr(&processHeapHandle, r)
		return r
	}

	return handle
}

func free(ptr uintptr) {
	_HeapFree.Call(getProcessHeap(), 0, ptr)
}

//go:linkname gostring runtime.gostring
func gostring(*byte) string

func cstring(v string) *byte {
	s := make([]byte, len(v)+1)
	copy(s, v)
	return (*byte)(unsafe.Pointer(&s[0]))
}

func cbool[T constraints.Integer](v bool) T {
	if v {
		return 1
	}
	return 0
}

type (
	wgpuAdapter             uintptr
	wgpuBindGroup           uintptr
	wgpuBindGroupLayout     uintptr
	wgpuBuffer              uintptr
	wgpuCommandBuffer       uintptr
	wgpuCommandEncoder      uintptr
	wgpuComputePassEncoder  uintptr
	wgpuComputePipeline     uintptr
	wgpuDevice              uintptr
	wgpuPipelineLayout      uintptr
	wgpuQuerySet            uintptr
	wgpuQueue               uintptr
	wgpuRenderBundleEncoder uintptr
	wgpuRenderPassEncoder   uintptr
	wgpuRenderPipeline      uintptr
	wgpuSampler             uintptr
	wgpuShaderModule        uintptr
	wgpuSurface             uintptr
	wgpuSwapChain           uintptr
	wgpuTexture             uintptr
	wgpuTextureView         uintptr
)

type sType uint32

// webgpu.h
const (
	sType_Invalid                                  = 0x00000000
	sType_SurfaceDescriptorFromMetalLayer          = 0x00000001
	sType_SurfaceDescriptorFromWindowsHWND         = 0x00000002
	sType_SurfaceDescriptorFromXlibWindow          = 0x00000003
	sType_SurfaceDescriptorFromCanvasHTMLSelector  = 0x00000004
	sType_ShaderModuleSPIRVDescriptor              = 0x00000005
	sType_ShaderModuleWGSLDescriptor               = 0x00000006
	sType_PrimitiveDepthClipControl                = 0x00000007
	sType_SurfaceDescriptorFromWaylandSurface      = 0x00000008
	sType_SurfaceDescriptorFromAndroidNativeWindow = 0x00000009
	sType_SurfaceDescriptorFromXcbWindow           = 0x0000000A
)

// wgpu.h
const (
	sType_DeviceExtras         = 0x60000001
	sType_AdapterExtras        = 0x60000002
	sType_RequiredLimitsExtras = 0x60000003
	sType_PipelineLayoutExtras = 0x60000004
)

type wgpuChainedStruct struct {
	next  *wgpuChainedStruct
	sType sType
	_     [4]byte
}

type wgpuChainedStructOut struct {
	next  *wgpuChainedStructOut
	sType sType
	_     [4]byte
}

type wgpuRequestAdapterOptions struct {
	nextInChain          *wgpuChainedStruct
	compatibleSurface    wgpuSurface
	powerPreference      PowerPreference
	forceFallbackAdapter bool
	_                    [3]byte
}

type wgpuAdapterExtras struct {
	chain   wgpuChainedStruct
	backend BackendType
	_       [4]byte
}

type wgpuSurfaceDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
}

type wgpuSurfaceDescriptorFromWindowsHWND struct {
	chain     wgpuChainedStruct
	hinstance unsafe.Pointer
	hwnd      unsafe.Pointer
}

type wgpuLimits struct {
	maxTextureDimension1D                     uint32
	maxTextureDimension2D                     uint32
	maxTextureDimension3D                     uint32
	maxTextureArrayLayers                     uint32
	maxBindGroups                             uint32
	maxDynamicUniformBuffersPerPipelineLayout uint32
	maxDynamicStorageBuffersPerPipelineLayout uint32
	maxSampledTexturesPerShaderStage          uint32
	maxSamplersPerShaderStage                 uint32
	maxStorageBuffersPerShaderStage           uint32
	maxStorageTexturesPerShaderStage          uint32
	maxUniformBuffersPerShaderStage           uint32
	maxUniformBufferBindingSize               uint64
	maxStorageBufferBindingSize               uint64
	minUniformBufferOffsetAlignment           uint32
	minStorageBufferOffsetAlignment           uint32
	maxVertexBuffers                          uint32
	maxVertexAttributes                       uint32
	maxVertexBufferArrayStride                uint32
	maxInterStageShaderComponents             uint32
	maxComputeWorkgroupStorageSize            uint32
	maxComputeInvocationsPerWorkgroup         uint32
	maxComputeWorkgroupSizeX                  uint32
	maxComputeWorkgroupSizeY                  uint32
	maxComputeWorkgroupSizeZ                  uint32
	maxComputeWorkgroupsPerDimension          uint32
}

type wgpuSupportedLimits struct {
	nextInChain *wgpuChainedStructOut
	limits      wgpuLimits
}

type wgpuAdapterProperties struct {
	nextInChain       *wgpuChainedStructOut
	vendorID          uint32
	deviceID          uint32
	name              *byte
	driverDescription *byte
	adapterType       AdapterType
	backendType       BackendType
}

type wgpuQueueDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
}

type wgpuRequiredLimits struct {
	nextInChain *wgpuChainedStruct
	limits      wgpuLimits
}

type wgpuRequiredLimitsExtras struct {
	chain               wgpuChainedStruct
	maxPushConstantSize uint32
	_                   [4]byte
}

type wgpuDeviceDescriptor struct {
	nextInChain           *wgpuChainedStruct
	label                 *byte
	requiredFeaturesCount uint32
	requiredFeatures      *FeatureName
	requiredLimits        *wgpuRequiredLimits
	defaultQueue          wgpuQueueDescriptor
}

type wgpuDeviceExtras struct {
	chain          wgpuChainedStruct
	nativeFeatures NativeFeature
	label          *byte
	tracePath      *byte
}

type wgpuBufferBindingLayout struct {
	nextInChain      *wgpuChainedStruct
	_type            BufferBindingType
	hasDynamicOffset bool
	minBindingSize   uint64
}

type wgpuSamplerBindingLayout struct {
	nextInChain *wgpuChainedStruct
	_type       SamplerBindingType
	_           [4]byte
}

type wgpuTextureBindingLayout struct {
	nextInChain   *wgpuChainedStruct
	sampleType    TextureSampleType
	viewDimension TextureViewDimension
	multisampled  bool
	_             [7]byte
}

type wgpuStorageTextureBindingLayout struct {
	nextInChain   *wgpuChainedStruct
	access        StorageTextureAccess
	format        TextureFormat
	viewDimension TextureViewDimension
	_             [4]byte
}

type wgpuBindGroupLayoutEntry struct {
	nextInChain    *wgpuChainedStruct
	binding        uint32
	visibility     ShaderStage
	buffer         wgpuBufferBindingLayout
	sampler        wgpuSamplerBindingLayout
	texture        wgpuTextureBindingLayout
	storageTexture wgpuStorageTextureBindingLayout
}

type wgpuBindGroupLayoutDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
	entryCount  uint32
	entries     *wgpuBindGroupLayoutEntry
}

type wgpuBindGroupEntry struct {
	nextInChain *wgpuChainedStruct
	binding     uint32
	buffer      wgpuBuffer
	offset      uint64
	size        uint64
	sampler     wgpuSampler
	textureView wgpuTextureView
}

type wgpuBindGroupDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
	layout      wgpuBindGroupLayout
	entryCount  uint32
	entries     *wgpuBindGroupEntry
}

type wgpuBufferDescriptor struct {
	nextInChain      *wgpuChainedStruct
	label            *byte
	usage            BufferUsage
	size             uint64
	mappedAtCreation bool
	_                [7]byte
}

type wgpuCommandEncoderDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
}

type wgpuConstantEntry struct {
	nextInChain *wgpuChainedStruct
	key         *byte
	value       float64
}

type wgpuProgrammableStageDescriptor struct {
	nextInChain   *wgpuChainedStruct
	module        wgpuShaderModule
	entryPoint    *byte
	constantCount uint32
	constants     *wgpuConstantEntry
}

type wgpuComputePipelineDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
	layout      wgpuPipelineLayout
	compute     wgpuProgrammableStageDescriptor
}

type wgpuPipelineLayoutDescriptor struct {
	nextInChain          *wgpuChainedStruct
	label                *byte
	bindGroupLayoutCount uint32
	bindGroupLayouts     *wgpuBindGroupLayout
}

type wgpuPipelineLayoutExtras struct {
	chain                  wgpuChainedStruct
	pushConstantRangeCount uint32
	pushConstantRanges     *wgpuPushConstantRange
}

type wgpuPushConstantRange struct {
	stages ShaderStage
	start  uint32
	end    uint32
}

type wgpuVertexAttribute struct {
	format         VertexFormat
	offset         uint64
	shaderLocation uint32
	_              [4]byte
}

type wgpuVertexBufferLayout struct {
	arrayStride    uint64
	stepMode       VertexStepMode
	attributeCount uint32
	attributes     *wgpuVertexAttribute
}

type wgpuVertexState struct {
	nextInChain   *wgpuChainedStruct
	module        wgpuShaderModule
	entryPoint    *byte
	constantCount uint32
	constants     *wgpuConstantEntry
	bufferCount   uint32
	buffers       *wgpuVertexBufferLayout
}

type wgpuPrimitiveState struct {
	nextInChain      *wgpuChainedStruct
	topology         PrimitiveTopology
	stripIndexFormat IndexFormat
	frontFace        FrontFace
	cullMode         CullMode
}

type wgpuStencilFaceState struct {
	compare     CompareFunction
	failOp      StencilOperation
	depthFailOp StencilOperation
	passOp      StencilOperation
}

type wgpuDepthStencilState struct {
	nextInChain         *wgpuChainedStruct
	format              TextureFormat
	depthWriteEnabled   bool
	depthCompare        CompareFunction
	stencilFront        wgpuStencilFaceState
	stencilBack         wgpuStencilFaceState
	stencilReadMask     uint32
	stencilWriteMask    uint32
	depthBias           int32
	depthBiasSlopeScale float32
	depthBiasClamp      float32
}

type wgpuMultisampleState struct {
	nextInChain            *wgpuChainedStruct
	count                  uint32
	mask                   uint32
	alphaToCoverageEnabled bool
	_                      [7]byte
}

type wgpuBlendComponent struct {
	operation BlendOperation
	srcFactor BlendFactor
	dstFactor BlendFactor
}

type wgpuBlendState struct {
	color wgpuBlendComponent
	alpha wgpuBlendComponent
}

type wgpuColorTargetState struct {
	nextInChain *wgpuChainedStruct
	format      TextureFormat
	blend       *wgpuBlendState
	writeMask   ColorWriteMask
	_           [4]byte
}

type wgpuFragmentState struct {
	nextInChain   *wgpuChainedStruct
	module        wgpuShaderModule
	entryPoint    *byte
	constantCount uint32
	constants     *wgpuConstantEntry
	targetCount   uint32
	targets       *wgpuColorTargetState
}

type wgpuRenderPipelineDescriptor struct {
	nextInChain  *wgpuChainedStruct
	label        *byte
	layout       wgpuPipelineLayout
	vertex       wgpuVertexState
	primitive    wgpuPrimitiveState
	depthStencil *wgpuDepthStencilState
	multisample  wgpuMultisampleState
	fragment     *wgpuFragmentState
}

type wgpuSamplerDescriptor struct {
	nextInChain   *wgpuChainedStruct
	label         *byte
	addressModeU  AddressMode
	addressModeV  AddressMode
	addressModeW  AddressMode
	magFilter     FilterMode
	minFilter     FilterMode
	mipmapFilter  MipmapFilterMode
	lodMinClamp   float32
	lodMaxClamp   float32
	compare       CompareFunction
	maxAnisotropy uint16
	_             [2]byte
}

type wgpuShaderModuleCompilationHint struct {
	nextInChain *wgpuChainedStruct
	entryPoint  *byte
	layout      wgpuPipelineLayout
}

type wgpuShaderModuleDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
	hintCount   uint32
	hints       *wgpuShaderModuleCompilationHint
}

type wgpuShaderModuleSPIRVDescriptor struct {
	chain    wgpuChainedStruct
	codeSize uint32
	code     *uint32
}

type wgpuShaderModuleWGSLDescriptor struct {
	chain wgpuChainedStruct
	code  *byte
}

type wgpuSwapChainDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
	usage       TextureUsage
	format      TextureFormat
	width       uint32
	height      uint32
	presentMode PresentMode
	_           [4]byte
}

type wgpuExtent3D struct {
	width              uint32
	height             uint32
	depthOrArrayLayers uint32
}

type wgpuTextureDescriptor struct {
	nextInChain     *wgpuChainedStruct
	label           *byte
	usage           TextureUsage
	dimension       TextureDimension
	size            wgpuExtent3D
	format          TextureFormat
	mipLevelCount   uint32
	sampleCount     uint32
	viewFormatCount uint32
	viewFormats     *TextureFormat
}

type wgpuComputePassTimestampWrite struct {
	querySet   wgpuQuerySet
	queryIndex uint32
	location   ComputePassTimestampLocation
}

type wgpuComputePassDescriptor struct {
	nextInChain         *wgpuChainedStruct
	label               *byte
	timestampWriteCount uint32
	timestampWrites     *wgpuComputePassTimestampWrite
}

type wgpuColor struct {
	r float64
	g float64
	b float64
	a float64
}

type wgpuRenderPassColorAttachment struct {
	view          wgpuTextureView
	resolveTarget wgpuTextureView
	loadOp        LoadOp
	storeOp       StoreOp
	clearValue    wgpuColor
}

type wgpuRenderPassDepthStencilAttachment struct {
	view              wgpuTextureView
	depthLoadOp       LoadOp
	depthStoreOp      StoreOp
	depthClearValue   float32
	depthReadOnly     bool
	stencilLoadOp     LoadOp
	stencilStoreOp    StoreOp
	stencilClearValue uint32
	stencilReadOnly   bool
	_                 [3]byte
}

type wgpuRenderPassTimestampWrite struct {
	querySet   wgpuQuerySet
	queryIndex uint32
	location   RenderPassTimestampLocation
}

type wgpuRenderPassDescriptor struct {
	nextInChain            *wgpuChainedStruct
	label                  *byte
	colorAttachmentCount   uint32
	colorAttachments       *wgpuRenderPassColorAttachment
	depthStencilAttachment *wgpuRenderPassDepthStencilAttachment
	occlusionQuerySet      wgpuQuerySet
	timestampWriteCount    uint32
	timestampWrites        *wgpuRenderPassTimestampWrite
}

type wgpuTextureDataLayout struct {
	nextInChain  *wgpuChainedStruct
	offset       uint64
	bytesPerRow  uint32
	rowsPerImage uint32
}

type wgpuImageCopyBuffer struct {
	nextInChain *wgpuChainedStruct
	layout      wgpuTextureDataLayout
	buffer      wgpuBuffer
}

type wgpuOrigin3D struct {
	x uint32
	y uint32
	z uint32
}

type wgpuImageCopyTexture struct {
	nextInChain *wgpuChainedStruct
	texture     wgpuTexture
	mipLevel    uint32
	origin      wgpuOrigin3D
	aspect      TextureAspect
	_           [4]byte
}

type wgpuCommandBufferDescriptor struct {
	nextInChain *wgpuChainedStruct
	label       *byte
}

type wgpuTextureViewDescriptor struct {
	nextInChain     *wgpuChainedStruct
	label           *byte
	format          TextureFormat
	dimension       TextureViewDimension
	baseMipLevel    uint32
	mipLevelCount   uint32
	baseArrayLayer  uint32
	arrayLayerCount uint32
	aspect          TextureAspect
	_               [4]byte
}
