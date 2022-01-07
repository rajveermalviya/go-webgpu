package wgpu

/*

#include "wrapper.h"

*/
import "C"

type AdapterType C.WGPUAdapterType

const (
	AdapterType_DiscreteGPU   AdapterType = C.WGPUAdapterType_DiscreteGPU
	AdapterType_IntegratedGPU AdapterType = C.WGPUAdapterType_IntegratedGPU
	AdapterType_CPU           AdapterType = C.WGPUAdapterType_CPU
	AdapterType_Unknown       AdapterType = C.WGPUAdapterType_Unknown
	AdapterType_Force32       AdapterType = C.WGPUAdapterType_Force32
)

type AddressMode C.WGPUAddressMode

const (
	AddressMode_Repeat       AddressMode = C.WGPUAddressMode_Repeat
	AddressMode_MirrorRepeat AddressMode = C.WGPUAddressMode_MirrorRepeat
	AddressMode_ClampToEdge  AddressMode = C.WGPUAddressMode_ClampToEdge
	AddressMode_Force32      AddressMode = C.WGPUAddressMode_Force32
)

type BackendType C.WGPUBackendType

const (
	BackendType_Null     BackendType = C.WGPUBackendType_Null
	BackendType_WebGPU   BackendType = C.WGPUBackendType_WebGPU
	BackendType_D3D11    BackendType = C.WGPUBackendType_D3D11
	BackendType_D3D12    BackendType = C.WGPUBackendType_D3D12
	BackendType_Metal    BackendType = C.WGPUBackendType_Metal
	BackendType_Vulkan   BackendType = C.WGPUBackendType_Vulkan
	BackendType_OpenGL   BackendType = C.WGPUBackendType_OpenGL
	BackendType_OpenGLES BackendType = C.WGPUBackendType_OpenGLES
	BackendType_Force32  BackendType = C.WGPUBackendType_Force32
)

type BlendFactor C.WGPUBlendFactor

const (
	BlendFactor_Zero              BlendFactor = C.WGPUBlendFactor_Zero
	BlendFactor_One               BlendFactor = C.WGPUBlendFactor_One
	BlendFactor_Src               BlendFactor = C.WGPUBlendFactor_Src
	BlendFactor_OneMinusSrc       BlendFactor = C.WGPUBlendFactor_OneMinusSrc
	BlendFactor_SrcAlpha          BlendFactor = C.WGPUBlendFactor_SrcAlpha
	BlendFactor_OneMinusSrcAlpha  BlendFactor = C.WGPUBlendFactor_OneMinusSrcAlpha
	BlendFactor_Dst               BlendFactor = C.WGPUBlendFactor_Dst
	BlendFactor_OneMinusDst       BlendFactor = C.WGPUBlendFactor_OneMinusDst
	BlendFactor_DstAlpha          BlendFactor = C.WGPUBlendFactor_DstAlpha
	BlendFactor_OneMinusDstAlpha  BlendFactor = C.WGPUBlendFactor_OneMinusDstAlpha
	BlendFactor_SrcAlphaSaturated BlendFactor = C.WGPUBlendFactor_SrcAlphaSaturated
	BlendFactor_Constant          BlendFactor = C.WGPUBlendFactor_Constant
	BlendFactor_OneMinusConstant  BlendFactor = C.WGPUBlendFactor_OneMinusConstant
	BlendFactor_Force32           BlendFactor = C.WGPUBlendFactor_Force32
)

type BlendOperation C.WGPUBlendOperation

const (
	BlendOperation_Add             BlendOperation = C.WGPUBlendOperation_Add
	BlendOperation_Subtract        BlendOperation = C.WGPUBlendOperation_Subtract
	BlendOperation_ReverseSubtract BlendOperation = C.WGPUBlendOperation_ReverseSubtract
	BlendOperation_Min             BlendOperation = C.WGPUBlendOperation_Min
	BlendOperation_Max             BlendOperation = C.WGPUBlendOperation_Max
	BlendOperation_Force32         BlendOperation = C.WGPUBlendOperation_Force32
)

type BufferBindingType C.WGPUBufferBindingType

const (
	BufferBindingType_Undefined       BufferBindingType = C.WGPUBufferBindingType_Undefined
	BufferBindingType_Uniform         BufferBindingType = C.WGPUBufferBindingType_Uniform
	BufferBindingType_Storage         BufferBindingType = C.WGPUBufferBindingType_Storage
	BufferBindingType_ReadOnlyStorage BufferBindingType = C.WGPUBufferBindingType_ReadOnlyStorage
	BufferBindingType_Force32         BufferBindingType = C.WGPUBufferBindingType_Force32
)

type BufferMapAsyncStatus C.WGPUBufferMapAsyncStatus

const (
	BufferMapAsyncStatus_Success                 BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_Success
	BufferMapAsyncStatus_Error                   BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_Error
	BufferMapAsyncStatus_Unknown                 BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_Unknown
	BufferMapAsyncStatus_DeviceLost              BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_DeviceLost
	BufferMapAsyncStatus_DestroyedBeforeCallback BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_DestroyedBeforeCallback
	BufferMapAsyncStatus_UnmappedBeforeCallback  BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_UnmappedBeforeCallback
	BufferMapAsyncStatus_Force32                 BufferMapAsyncStatus = C.WGPUBufferMapAsyncStatus_Force32
)

type CompareFunction C.WGPUCompareFunction

const (
	CompareFunction_Undefined    CompareFunction = C.WGPUCompareFunction_Undefined
	CompareFunction_Never        CompareFunction = C.WGPUCompareFunction_Never
	CompareFunction_Less         CompareFunction = C.WGPUCompareFunction_Less
	CompareFunction_LessEqual    CompareFunction = C.WGPUCompareFunction_LessEqual
	CompareFunction_Greater      CompareFunction = C.WGPUCompareFunction_Greater
	CompareFunction_GreaterEqual CompareFunction = C.WGPUCompareFunction_GreaterEqual
	CompareFunction_Equal        CompareFunction = C.WGPUCompareFunction_Equal
	CompareFunction_NotEqual     CompareFunction = C.WGPUCompareFunction_NotEqual
	CompareFunction_Always       CompareFunction = C.WGPUCompareFunction_Always
	CompareFunction_Force32      CompareFunction = C.WGPUCompareFunction_Force32
)

type CompilationInfoRequestStatus C.WGPUCompilationInfoRequestStatus

const (
	CompilationInfoRequestStatus_Success    CompilationInfoRequestStatus = C.WGPUCompilationInfoRequestStatus_Success
	CompilationInfoRequestStatus_Error      CompilationInfoRequestStatus = C.WGPUCompilationInfoRequestStatus_Error
	CompilationInfoRequestStatus_DeviceLost CompilationInfoRequestStatus = C.WGPUCompilationInfoRequestStatus_DeviceLost
	CompilationInfoRequestStatus_Unknown    CompilationInfoRequestStatus = C.WGPUCompilationInfoRequestStatus_Unknown
	CompilationInfoRequestStatus_Force32    CompilationInfoRequestStatus = C.WGPUCompilationInfoRequestStatus_Force32
)

type CompilationMessageType C.WGPUCompilationMessageType

const (
	CompilationMessageType_Error   CompilationMessageType = C.WGPUCompilationMessageType_Error
	CompilationMessageType_Warning CompilationMessageType = C.WGPUCompilationMessageType_Warning
	CompilationMessageType_Info    CompilationMessageType = C.WGPUCompilationMessageType_Info
	CompilationMessageType_Force32 CompilationMessageType = C.WGPUCompilationMessageType_Force32
)

type ComputePassTimestampLocation C.WGPUComputePassTimestampLocation

const (
	ComputePassTimestampLocation_Beginning ComputePassTimestampLocation = C.WGPUComputePassTimestampLocation_Beginning
	ComputePassTimestampLocation_End       ComputePassTimestampLocation = C.WGPUComputePassTimestampLocation_End
	ComputePassTimestampLocation_Force32   ComputePassTimestampLocation = C.WGPUComputePassTimestampLocation_Force32
)

type CreatePipelineAsyncStatus C.WGPUCreatePipelineAsyncStatus

const (
	CreatePipelineAsyncStatus_Success         CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_Success
	CreatePipelineAsyncStatus_Error           CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_Error
	CreatePipelineAsyncStatus_DeviceLost      CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_DeviceLost
	CreatePipelineAsyncStatus_DeviceDestroyed CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_DeviceDestroyed
	CreatePipelineAsyncStatus_Unknown         CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_Unknown
	CreatePipelineAsyncStatus_Force32         CreatePipelineAsyncStatus = C.WGPUCreatePipelineAsyncStatus_Force32
)

type CullMode C.WGPUCullMode

const (
	CullMode_None    CullMode = C.WGPUCullMode_None
	CullMode_Front   CullMode = C.WGPUCullMode_Front
	CullMode_Back    CullMode = C.WGPUCullMode_Back
	CullMode_Force32 CullMode = C.WGPUCullMode_Force32
)

type DeviceLostReason C.WGPUDeviceLostReason

const (
	DeviceLostReason_Undefined DeviceLostReason = C.WGPUDeviceLostReason_Undefined
	DeviceLostReason_Destroyed DeviceLostReason = C.WGPUDeviceLostReason_Destroyed
	DeviceLostReason_Force32   DeviceLostReason = C.WGPUDeviceLostReason_Force32
)

type ErrorFilter C.WGPUErrorFilter

const (
	ErrorFilter_None        = C.WGPUErrorFilter_None
	ErrorFilter_Validation  = C.WGPUErrorFilter_Validation
	ErrorFilter_OutOfMemory = C.WGPUErrorFilter_OutOfMemory
	ErrorFilter_Force32     = C.WGPUErrorFilter_Force32
)

type ErrorType C.WGPUErrorType

const (
	ErrorType_NoError     ErrorType = C.WGPUErrorType_NoError
	ErrorType_Validation  ErrorType = C.WGPUErrorType_Validation
	ErrorType_OutOfMemory ErrorType = C.WGPUErrorType_OutOfMemory
	ErrorType_Unknown     ErrorType = C.WGPUErrorType_Unknown
	ErrorType_DeviceLost  ErrorType = C.WGPUErrorType_DeviceLost
	ErrorType_Force32     ErrorType = C.WGPUErrorType_Force32
)

type FeatureName C.WGPUFeatureName

const (
	FeatureName_Undefined               FeatureName = C.WGPUFeatureName_Undefined
	FeatureName_DepthClipControl        FeatureName = C.WGPUFeatureName_DepthClipControl
	FeatureName_Depth24UnormStencil8    FeatureName = C.WGPUFeatureName_Depth24UnormStencil8
	FeatureName_Depth32FloatStencil8    FeatureName = C.WGPUFeatureName_Depth32FloatStencil8
	FeatureName_TimestampQuery          FeatureName = C.WGPUFeatureName_TimestampQuery
	FeatureName_PipelineStatisticsQuery FeatureName = C.WGPUFeatureName_PipelineStatisticsQuery
	FeatureName_TextureCompressionBC    FeatureName = C.WGPUFeatureName_TextureCompressionBC
	FeatureName_TextureCompressionETC2  FeatureName = C.WGPUFeatureName_TextureCompressionETC2
	FeatureName_TextureCompressionASTC  FeatureName = C.WGPUFeatureName_TextureCompressionASTC
	FeatureName_IndirectFirstInstance   FeatureName = C.WGPUFeatureName_IndirectFirstInstance
	FeatureName_Force32                 FeatureName = C.WGPUFeatureName_Force32
)

type FilterMode C.WGPUFilterMode

const (
	FilterMode_Nearest FilterMode = C.WGPUFilterMode_Nearest
	FilterMode_Linear  FilterMode = C.WGPUFilterMode_Linear
	FilterMode_Force32 FilterMode = C.WGPUFilterMode_Force32
)

type FrontFace C.WGPUFrontFace

const (
	FrontFace_CCW     FrontFace = C.WGPUFrontFace_CCW
	FrontFace_CW      FrontFace = C.WGPUFrontFace_CW
	FrontFace_Force32 FrontFace = C.WGPUFrontFace_Force32
)

type IndexFormat C.WGPUIndexFormat

const (
	IndexFormat_Undefined IndexFormat = C.WGPUIndexFormat_Undefined
	IndexFormat_Uint16    IndexFormat = C.WGPUIndexFormat_Uint16
	IndexFormat_Uint32    IndexFormat = C.WGPUIndexFormat_Uint32
	IndexFormat_Force32   IndexFormat = C.WGPUIndexFormat_Force32
)

type LoadOp C.WGPULoadOp

const (
	LoadOp_Clear   LoadOp = C.WGPULoadOp_Clear
	LoadOp_Load    LoadOp = C.WGPULoadOp_Load
	LoadOp_Force32 LoadOp = C.WGPULoadOp_Force32
)

type PipelineStatisticName C.WGPUPipelineStatisticName

const (
	PipelineStatisticName_VertexShaderInvocations   PipelineStatisticName = C.WGPUPipelineStatisticName_VertexShaderInvocations
	PipelineStatisticName_ClipperInvocations        PipelineStatisticName = C.WGPUPipelineStatisticName_ClipperInvocations
	PipelineStatisticName_ClipperPrimitivesOut      PipelineStatisticName = C.WGPUPipelineStatisticName_ClipperPrimitivesOut
	PipelineStatisticName_FragmentShaderInvocations PipelineStatisticName = C.WGPUPipelineStatisticName_FragmentShaderInvocations
	PipelineStatisticName_ComputeShaderInvocations  PipelineStatisticName = C.WGPUPipelineStatisticName_ComputeShaderInvocations
	PipelineStatisticName_Force32                   PipelineStatisticName = C.WGPUPipelineStatisticName_Force32
)

type PowerPreference C.WGPUPowerPreference

const (
	PowerPreference_Undefined       PowerPreference = C.WGPUPowerPreference_Undefined
	PowerPreference_LowPower        PowerPreference = C.WGPUPowerPreference_LowPower
	PowerPreference_HighPerformance PowerPreference = C.WGPUPowerPreference_HighPerformance
	PowerPreference_Force32         PowerPreference = C.WGPUPowerPreference_Force32
)

type PresentMode C.WGPUPresentMode

const (
	PresentMode_Immediate PresentMode = C.WGPUPresentMode_Immediate
	PresentMode_Mailbox   PresentMode = C.WGPUPresentMode_Mailbox
	PresentMode_Fifo      PresentMode = C.WGPUPresentMode_Fifo
	PresentMode_Force32   PresentMode = C.WGPUPresentMode_Force32
)

type PrimitiveTopology C.WGPUPrimitiveTopology

const (
	PrimitiveTopology_PointList     PrimitiveTopology = C.WGPUPrimitiveTopology_PointList
	PrimitiveTopology_LineList      PrimitiveTopology = C.WGPUPrimitiveTopology_LineList
	PrimitiveTopology_LineStrip     PrimitiveTopology = C.WGPUPrimitiveTopology_LineStrip
	PrimitiveTopology_TriangleList  PrimitiveTopology = C.WGPUPrimitiveTopology_TriangleList
	PrimitiveTopology_TriangleStrip PrimitiveTopology = C.WGPUPrimitiveTopology_TriangleStrip
	PrimitiveTopology_Force32       PrimitiveTopology = C.WGPUPrimitiveTopology_Force32
)

type QueryType C.WGPUQueryType

const (
	QueryType_Occlusion          QueryType = C.WGPUQueryType_Occlusion
	QueryType_PipelineStatistics QueryType = C.WGPUQueryType_PipelineStatistics
	QueryType_Timestamp          QueryType = C.WGPUQueryType_Timestamp
	QueryType_Force32            QueryType = C.WGPUQueryType_Force32
)

type QueueWorkDoneStatus C.WGPUQueueWorkDoneStatus

const (
	QueueWorkDoneStatus_Success    QueueWorkDoneStatus = C.WGPUQueueWorkDoneStatus_Success
	QueueWorkDoneStatus_Error      QueueWorkDoneStatus = C.WGPUQueueWorkDoneStatus_Error
	QueueWorkDoneStatus_Unknown    QueueWorkDoneStatus = C.WGPUQueueWorkDoneStatus_Unknown
	QueueWorkDoneStatus_DeviceLost QueueWorkDoneStatus = C.WGPUQueueWorkDoneStatus_DeviceLost
	QueueWorkDoneStatus_Force32    QueueWorkDoneStatus = C.WGPUQueueWorkDoneStatus_Force32
)

type RenderPassTimestampLocation C.WGPURenderPassTimestampLocation

const (
	RenderPassTimestampLocation_Beginning RenderPassTimestampLocation = C.WGPURenderPassTimestampLocation_Beginning
	RenderPassTimestampLocation_End       RenderPassTimestampLocation = C.WGPURenderPassTimestampLocation_End
	RenderPassTimestampLocation_Force32   RenderPassTimestampLocation = C.WGPURenderPassTimestampLocation_Force32
)

type RequestAdapterStatus C.WGPURequestAdapterStatus

const (
	RequestAdapterStatus_Success     RequestAdapterStatus = C.WGPURequestAdapterStatus_Success
	RequestAdapterStatus_Unavailable RequestAdapterStatus = C.WGPURequestAdapterStatus_Unavailable
	RequestAdapterStatus_Error       RequestAdapterStatus = C.WGPURequestAdapterStatus_Error
	RequestAdapterStatus_Unknown     RequestAdapterStatus = C.WGPURequestAdapterStatus_Unknown
	RequestAdapterStatus_Force32     RequestAdapterStatus = C.WGPURequestAdapterStatus_Force32
)

type RequestDeviceStatus C.WGPURequestDeviceStatus

const (
	RequestDeviceStatus_Success RequestDeviceStatus = C.WGPURequestDeviceStatus_Success
	RequestDeviceStatus_Error   RequestDeviceStatus = C.WGPURequestDeviceStatus_Error
	RequestDeviceStatus_Unknown RequestDeviceStatus = C.WGPURequestDeviceStatus_Unknown
	RequestDeviceStatus_Force32 RequestDeviceStatus = C.WGPURequestDeviceStatus_Force32
)

// type SType C.WGPUSType

// const (
// 	SType_Invalid                                 SType = C.WGPUSType_Invalid
// 	SType_SurfaceDescriptorFromMetalLayer         SType = C.WGPUSType_SurfaceDescriptorFromMetalLayer
// 	SType_SurfaceDescriptorFromWindowsHWND        SType = C.WGPUSType_SurfaceDescriptorFromWindowsHWND
// 	SType_SurfaceDescriptorFromXlib               SType = C.WGPUSType_SurfaceDescriptorFromXlib
// 	SType_SurfaceDescriptorFromCanvasHTMLSelector SType = C.WGPUSType_SurfaceDescriptorFromCanvasHTMLSelector
// 	SType_ShaderModuleSPIRVDescriptor             SType = C.WGPUSType_ShaderModuleSPIRVDescriptor
// 	SType_ShaderModuleWGSLDescriptor              SType = C.WGPUSType_ShaderModuleWGSLDescriptor
// 	SType_PrimitiveDepthClipControl               SType = C.WGPUSType_PrimitiveDepthClipControl
// 	SType_Force32                                 SType = C.WGPUSType_Force32
// )

type SamplerBindingType C.WGPUSamplerBindingType

const (
	SamplerBindingType_Undefined    SamplerBindingType = C.WGPUSamplerBindingType_Undefined
	SamplerBindingType_Filtering    SamplerBindingType = C.WGPUSamplerBindingType_Filtering
	SamplerBindingType_NonFiltering SamplerBindingType = C.WGPUSamplerBindingType_NonFiltering
	SamplerBindingType_Comparison   SamplerBindingType = C.WGPUSamplerBindingType_Comparison
	SamplerBindingType_Force32      SamplerBindingType = C.WGPUSamplerBindingType_Force32
)

type StencilOperation C.WGPUStencilOperation

const (
	StencilOperation_Keep           StencilOperation = C.WGPUStencilOperation_Keep
	StencilOperation_Zero           StencilOperation = C.WGPUStencilOperation_Zero
	StencilOperation_Replace        StencilOperation = C.WGPUStencilOperation_Replace
	StencilOperation_Invert         StencilOperation = C.WGPUStencilOperation_Invert
	StencilOperation_IncrementClamp StencilOperation = C.WGPUStencilOperation_IncrementClamp
	StencilOperation_DecrementClamp StencilOperation = C.WGPUStencilOperation_DecrementClamp
	StencilOperation_IncrementWrap  StencilOperation = C.WGPUStencilOperation_IncrementWrap
	StencilOperation_DecrementWrap  StencilOperation = C.WGPUStencilOperation_DecrementWrap
	StencilOperation_Force32        StencilOperation = C.WGPUStencilOperation_Force32
)

type StorageTextureAccess C.WGPUStorageTextureAccess

const (
	StorageTextureAccess_Undefined StorageTextureAccess = C.WGPUStorageTextureAccess_Undefined
	StorageTextureAccess_WriteOnly StorageTextureAccess = C.WGPUStorageTextureAccess_WriteOnly
	StorageTextureAccess_Force32   StorageTextureAccess = C.WGPUStorageTextureAccess_Force32
)

type StoreOp C.WGPUStoreOp

const (
	StoreOp_Store   = C.WGPUStoreOp_Store
	StoreOp_Discard = C.WGPUStoreOp_Discard
	StoreOp_Force32 = C.WGPUStoreOp_Force32
)

type TextureAspect C.WGPUTextureAspect

const (
	TextureAspect_All         TextureAspect = C.WGPUTextureAspect_All
	TextureAspect_StencilOnly TextureAspect = C.WGPUTextureAspect_StencilOnly
	TextureAspect_DepthOnly   TextureAspect = C.WGPUTextureAspect_DepthOnly
	TextureAspect_Force32     TextureAspect = C.WGPUTextureAspect_Force32
)

type TextureComponentType C.WGPUTextureComponentType

const (
	TextureComponentType_Float           TextureComponentType = C.WGPUTextureComponentType_Float
	TextureComponentType_Sint            TextureComponentType = C.WGPUTextureComponentType_Sint
	TextureComponentType_Uint            TextureComponentType = C.WGPUTextureComponentType_Uint
	TextureComponentType_DepthComparison TextureComponentType = C.WGPUTextureComponentType_DepthComparison
	TextureComponentType_Force32         TextureComponentType = C.WGPUTextureComponentType_Force32
)

type TextureDimension C.WGPUTextureDimension

const (
	TextureDimension_1D      TextureDimension = C.WGPUTextureDimension_1D
	TextureDimension_2D      TextureDimension = C.WGPUTextureDimension_2D
	TextureDimension_3D      TextureDimension = C.WGPUTextureDimension_3D
	TextureDimension_Force32 TextureDimension = C.WGPUTextureDimension_Force32
)

type TextureFormat C.WGPUTextureFormat

const (
	TextureFormat_Undefined            TextureFormat = C.WGPUTextureFormat_Undefined
	TextureFormat_R8Unorm              TextureFormat = C.WGPUTextureFormat_R8Unorm
	TextureFormat_R8Snorm              TextureFormat = C.WGPUTextureFormat_R8Snorm
	TextureFormat_R8Uint               TextureFormat = C.WGPUTextureFormat_R8Uint
	TextureFormat_R8Sint               TextureFormat = C.WGPUTextureFormat_R8Sint
	TextureFormat_R16Uint              TextureFormat = C.WGPUTextureFormat_R16Uint
	TextureFormat_R16Sint              TextureFormat = C.WGPUTextureFormat_R16Sint
	TextureFormat_R16Float             TextureFormat = C.WGPUTextureFormat_R16Float
	TextureFormat_RG8Unorm             TextureFormat = C.WGPUTextureFormat_RG8Unorm
	TextureFormat_RG8Snorm             TextureFormat = C.WGPUTextureFormat_RG8Snorm
	TextureFormat_RG8Uint              TextureFormat = C.WGPUTextureFormat_RG8Uint
	TextureFormat_RG8Sint              TextureFormat = C.WGPUTextureFormat_RG8Sint
	TextureFormat_R32Float             TextureFormat = C.WGPUTextureFormat_R32Float
	TextureFormat_R32Uint              TextureFormat = C.WGPUTextureFormat_R32Uint
	TextureFormat_R32Sint              TextureFormat = C.WGPUTextureFormat_R32Sint
	TextureFormat_RG16Uint             TextureFormat = C.WGPUTextureFormat_RG16Uint
	TextureFormat_RG16Sint             TextureFormat = C.WGPUTextureFormat_RG16Sint
	TextureFormat_RG16Float            TextureFormat = C.WGPUTextureFormat_RG16Float
	TextureFormat_RGBA8Unorm           TextureFormat = C.WGPUTextureFormat_RGBA8Unorm
	TextureFormat_RGBA8UnormSrgb       TextureFormat = C.WGPUTextureFormat_RGBA8UnormSrgb
	TextureFormat_RGBA8Snorm           TextureFormat = C.WGPUTextureFormat_RGBA8Snorm
	TextureFormat_RGBA8Uint            TextureFormat = C.WGPUTextureFormat_RGBA8Uint
	TextureFormat_RGBA8Sint            TextureFormat = C.WGPUTextureFormat_RGBA8Sint
	TextureFormat_BGRA8Unorm           TextureFormat = C.WGPUTextureFormat_BGRA8Unorm
	TextureFormat_BGRA8UnormSrgb       TextureFormat = C.WGPUTextureFormat_BGRA8UnormSrgb
	TextureFormat_RGB10A2Unorm         TextureFormat = C.WGPUTextureFormat_RGB10A2Unorm
	TextureFormat_RG11B10Ufloat        TextureFormat = C.WGPUTextureFormat_RG11B10Ufloat
	TextureFormat_RGB9E5Ufloat         TextureFormat = C.WGPUTextureFormat_RGB9E5Ufloat
	TextureFormat_RG32Float            TextureFormat = C.WGPUTextureFormat_RG32Float
	TextureFormat_RG32Uint             TextureFormat = C.WGPUTextureFormat_RG32Uint
	TextureFormat_RG32Sint             TextureFormat = C.WGPUTextureFormat_RG32Sint
	TextureFormat_RGBA16Uint           TextureFormat = C.WGPUTextureFormat_RGBA16Uint
	TextureFormat_RGBA16Sint           TextureFormat = C.WGPUTextureFormat_RGBA16Sint
	TextureFormat_RGBA16Float          TextureFormat = C.WGPUTextureFormat_RGBA16Float
	TextureFormat_RGBA32Float          TextureFormat = C.WGPUTextureFormat_RGBA32Float
	TextureFormat_RGBA32Uint           TextureFormat = C.WGPUTextureFormat_RGBA32Uint
	TextureFormat_RGBA32Sint           TextureFormat = C.WGPUTextureFormat_RGBA32Sint
	TextureFormat_Stencil8             TextureFormat = C.WGPUTextureFormat_Stencil8
	TextureFormat_Depth16Unorm         TextureFormat = C.WGPUTextureFormat_Depth16Unorm
	TextureFormat_Depth24Plus          TextureFormat = C.WGPUTextureFormat_Depth24Plus
	TextureFormat_Depth24PlusStencil8  TextureFormat = C.WGPUTextureFormat_Depth24PlusStencil8
	TextureFormat_Depth24UnormStencil8 TextureFormat = C.WGPUTextureFormat_Depth24UnormStencil8
	TextureFormat_Depth32Float         TextureFormat = C.WGPUTextureFormat_Depth32Float
	TextureFormat_Depth32FloatStencil8 TextureFormat = C.WGPUTextureFormat_Depth32FloatStencil8
	TextureFormat_BC1RGBAUnorm         TextureFormat = C.WGPUTextureFormat_BC1RGBAUnorm
	TextureFormat_BC1RGBAUnormSrgb     TextureFormat = C.WGPUTextureFormat_BC1RGBAUnormSrgb
	TextureFormat_BC2RGBAUnorm         TextureFormat = C.WGPUTextureFormat_BC2RGBAUnorm
	TextureFormat_BC2RGBAUnormSrgb     TextureFormat = C.WGPUTextureFormat_BC2RGBAUnormSrgb
	TextureFormat_BC3RGBAUnorm         TextureFormat = C.WGPUTextureFormat_BC3RGBAUnorm
	TextureFormat_BC3RGBAUnormSrgb     TextureFormat = C.WGPUTextureFormat_BC3RGBAUnormSrgb
	TextureFormat_BC4RUnorm            TextureFormat = C.WGPUTextureFormat_BC4RUnorm
	TextureFormat_BC4RSnorm            TextureFormat = C.WGPUTextureFormat_BC4RSnorm
	TextureFormat_BC5RGUnorm           TextureFormat = C.WGPUTextureFormat_BC5RGUnorm
	TextureFormat_BC5RGSnorm           TextureFormat = C.WGPUTextureFormat_BC5RGSnorm
	TextureFormat_BC6HRGBUfloat        TextureFormat = C.WGPUTextureFormat_BC6HRGBUfloat
	TextureFormat_BC6HRGBFloat         TextureFormat = C.WGPUTextureFormat_BC6HRGBFloat
	TextureFormat_BC7RGBAUnorm         TextureFormat = C.WGPUTextureFormat_BC7RGBAUnorm
	TextureFormat_BC7RGBAUnormSrgb     TextureFormat = C.WGPUTextureFormat_BC7RGBAUnormSrgb
	TextureFormat_ETC2RGB8Unorm        TextureFormat = C.WGPUTextureFormat_ETC2RGB8Unorm
	TextureFormat_ETC2RGB8UnormSrgb    TextureFormat = C.WGPUTextureFormat_ETC2RGB8UnormSrgb
	TextureFormat_ETC2RGB8A1Unorm      TextureFormat = C.WGPUTextureFormat_ETC2RGB8A1Unorm
	TextureFormat_ETC2RGB8A1UnormSrgb  TextureFormat = C.WGPUTextureFormat_ETC2RGB8A1UnormSrgb
	TextureFormat_ETC2RGBA8Unorm       TextureFormat = C.WGPUTextureFormat_ETC2RGBA8Unorm
	TextureFormat_ETC2RGBA8UnormSrgb   TextureFormat = C.WGPUTextureFormat_ETC2RGBA8UnormSrgb
	TextureFormat_EACR11Unorm          TextureFormat = C.WGPUTextureFormat_EACR11Unorm
	TextureFormat_EACR11Snorm          TextureFormat = C.WGPUTextureFormat_EACR11Snorm
	TextureFormat_EACRG11Unorm         TextureFormat = C.WGPUTextureFormat_EACRG11Unorm
	TextureFormat_EACRG11Snorm         TextureFormat = C.WGPUTextureFormat_EACRG11Snorm
	TextureFormat_ASTC4x4Unorm         TextureFormat = C.WGPUTextureFormat_ASTC4x4Unorm
	TextureFormat_ASTC4x4UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC4x4UnormSrgb
	TextureFormat_ASTC5x4Unorm         TextureFormat = C.WGPUTextureFormat_ASTC5x4Unorm
	TextureFormat_ASTC5x4UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC5x4UnormSrgb
	TextureFormat_ASTC5x5Unorm         TextureFormat = C.WGPUTextureFormat_ASTC5x5Unorm
	TextureFormat_ASTC5x5UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC5x5UnormSrgb
	TextureFormat_ASTC6x5Unorm         TextureFormat = C.WGPUTextureFormat_ASTC6x5Unorm
	TextureFormat_ASTC6x5UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC6x5UnormSrgb
	TextureFormat_ASTC6x6Unorm         TextureFormat = C.WGPUTextureFormat_ASTC6x6Unorm
	TextureFormat_ASTC6x6UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC6x6UnormSrgb
	TextureFormat_ASTC8x5Unorm         TextureFormat = C.WGPUTextureFormat_ASTC8x5Unorm
	TextureFormat_ASTC8x5UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC8x5UnormSrgb
	TextureFormat_ASTC8x6Unorm         TextureFormat = C.WGPUTextureFormat_ASTC8x6Unorm
	TextureFormat_ASTC8x6UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC8x6UnormSrgb
	TextureFormat_ASTC8x8Unorm         TextureFormat = C.WGPUTextureFormat_ASTC8x8Unorm
	TextureFormat_ASTC8x8UnormSrgb     TextureFormat = C.WGPUTextureFormat_ASTC8x8UnormSrgb
	TextureFormat_ASTC10x5Unorm        TextureFormat = C.WGPUTextureFormat_ASTC10x5Unorm
	TextureFormat_ASTC10x5UnormSrgb    TextureFormat = C.WGPUTextureFormat_ASTC10x5UnormSrgb
	TextureFormat_ASTC10x6Unorm        TextureFormat = C.WGPUTextureFormat_ASTC10x6Unorm
	TextureFormat_ASTC10x6UnormSrgb    TextureFormat = C.WGPUTextureFormat_ASTC10x6UnormSrgb
	TextureFormat_ASTC10x8Unorm        TextureFormat = C.WGPUTextureFormat_ASTC10x8Unorm
	TextureFormat_ASTC10x8UnormSrgb    TextureFormat = C.WGPUTextureFormat_ASTC10x8UnormSrgb
	TextureFormat_ASTC10x10Unorm       TextureFormat = C.WGPUTextureFormat_ASTC10x10Unorm
	TextureFormat_ASTC10x10UnormSrgb   TextureFormat = C.WGPUTextureFormat_ASTC10x10UnormSrgb
	TextureFormat_ASTC12x10Unorm       TextureFormat = C.WGPUTextureFormat_ASTC12x10Unorm
	TextureFormat_ASTC12x10UnormSrgb   TextureFormat = C.WGPUTextureFormat_ASTC12x10UnormSrgb
	TextureFormat_ASTC12x12Unorm       TextureFormat = C.WGPUTextureFormat_ASTC12x12Unorm
	TextureFormat_ASTC12x12UnormSrgb   TextureFormat = C.WGPUTextureFormat_ASTC12x12UnormSrgb
	TextureFormat_Force32              TextureFormat = C.WGPUTextureFormat_Force32
)

type TextureSampleType C.WGPUTextureSampleType

const (
	TextureSampleType_Undefined         TextureSampleType = C.WGPUTextureSampleType_Undefined
	TextureSampleType_Float             TextureSampleType = C.WGPUTextureSampleType_Float
	TextureSampleType_UnfilterableFloat TextureSampleType = C.WGPUTextureSampleType_UnfilterableFloat
	TextureSampleType_Depth             TextureSampleType = C.WGPUTextureSampleType_Depth
	TextureSampleType_Sint              TextureSampleType = C.WGPUTextureSampleType_Sint
	TextureSampleType_Uint              TextureSampleType = C.WGPUTextureSampleType_Uint
	TextureSampleType_Force32           TextureSampleType = C.WGPUTextureSampleType_Force32
)

type TextureViewDimension C.WGPUTextureViewDimension

const (
	TextureViewDimension_Undefined TextureViewDimension = C.WGPUTextureViewDimension_Undefined
	TextureViewDimension_1D        TextureViewDimension = C.WGPUTextureViewDimension_1D
	TextureViewDimension_2D        TextureViewDimension = C.WGPUTextureViewDimension_2D
	TextureViewDimension_2DArray   TextureViewDimension = C.WGPUTextureViewDimension_2DArray
	TextureViewDimension_Cube      TextureViewDimension = C.WGPUTextureViewDimension_Cube
	TextureViewDimension_CubeArray TextureViewDimension = C.WGPUTextureViewDimension_CubeArray
	TextureViewDimension_3D        TextureViewDimension = C.WGPUTextureViewDimension_3D
	TextureViewDimension_Force32   TextureViewDimension = C.WGPUTextureViewDimension_Force32
)

type VertexFormat C.WGPUVertexFormat

const (
	VertexFormat_Undefined VertexFormat = C.WGPUVertexFormat_Undefined
	VertexFormat_Uint8x2   VertexFormat = C.WGPUVertexFormat_Uint8x2
	VertexFormat_Uint8x4   VertexFormat = C.WGPUVertexFormat_Uint8x4
	VertexFormat_Sint8x2   VertexFormat = C.WGPUVertexFormat_Sint8x2
	VertexFormat_Sint8x4   VertexFormat = C.WGPUVertexFormat_Sint8x4
	VertexFormat_Unorm8x2  VertexFormat = C.WGPUVertexFormat_Unorm8x2
	VertexFormat_Unorm8x4  VertexFormat = C.WGPUVertexFormat_Unorm8x4
	VertexFormat_Snorm8x2  VertexFormat = C.WGPUVertexFormat_Snorm8x2
	VertexFormat_Snorm8x4  VertexFormat = C.WGPUVertexFormat_Snorm8x4
	VertexFormat_Uint16x2  VertexFormat = C.WGPUVertexFormat_Uint16x2
	VertexFormat_Uint16x4  VertexFormat = C.WGPUVertexFormat_Uint16x4
	VertexFormat_Sint16x2  VertexFormat = C.WGPUVertexFormat_Sint16x2
	VertexFormat_Sint16x4  VertexFormat = C.WGPUVertexFormat_Sint16x4
	VertexFormat_Unorm16x2 VertexFormat = C.WGPUVertexFormat_Unorm16x2
	VertexFormat_Unorm16x4 VertexFormat = C.WGPUVertexFormat_Unorm16x4
	VertexFormat_Snorm16x2 VertexFormat = C.WGPUVertexFormat_Snorm16x2
	VertexFormat_Snorm16x4 VertexFormat = C.WGPUVertexFormat_Snorm16x4
	VertexFormat_Float16x2 VertexFormat = C.WGPUVertexFormat_Float16x2
	VertexFormat_Float16x4 VertexFormat = C.WGPUVertexFormat_Float16x4
	VertexFormat_Float32   VertexFormat = C.WGPUVertexFormat_Float32
	VertexFormat_Float32x2 VertexFormat = C.WGPUVertexFormat_Float32x2
	VertexFormat_Float32x3 VertexFormat = C.WGPUVertexFormat_Float32x3
	VertexFormat_Float32x4 VertexFormat = C.WGPUVertexFormat_Float32x4
	VertexFormat_Uint32    VertexFormat = C.WGPUVertexFormat_Uint32
	VertexFormat_Uint32x2  VertexFormat = C.WGPUVertexFormat_Uint32x2
	VertexFormat_Uint32x3  VertexFormat = C.WGPUVertexFormat_Uint32x3
	VertexFormat_Uint32x4  VertexFormat = C.WGPUVertexFormat_Uint32x4
	VertexFormat_Sint32    VertexFormat = C.WGPUVertexFormat_Sint32
	VertexFormat_Sint32x2  VertexFormat = C.WGPUVertexFormat_Sint32x2
	VertexFormat_Sint32x3  VertexFormat = C.WGPUVertexFormat_Sint32x3
	VertexFormat_Sint32x4  VertexFormat = C.WGPUVertexFormat_Sint32x4
	VertexFormat_Force32   VertexFormat = C.WGPUVertexFormat_Force32
)

type VertexStepMode C.WGPUVertexStepMode

const (
	VertexStepMode_Vertex   = C.WGPUVertexStepMode_Vertex
	VertexStepMode_Instance = C.WGPUVertexStepMode_Instance
	VertexStepMode_Force32  = C.WGPUVertexStepMode_Force32
)

type BufferUsage C.WGPUBufferUsage

const (
	BufferUsage_None         BufferUsage = C.WGPUBufferUsage_None
	BufferUsage_MapRead      BufferUsage = C.WGPUBufferUsage_MapRead
	BufferUsage_MapWrite     BufferUsage = C.WGPUBufferUsage_MapWrite
	BufferUsage_CopySrc      BufferUsage = C.WGPUBufferUsage_CopySrc
	BufferUsage_CopyDst      BufferUsage = C.WGPUBufferUsage_CopyDst
	BufferUsage_Index        BufferUsage = C.WGPUBufferUsage_Index
	BufferUsage_Vertex       BufferUsage = C.WGPUBufferUsage_Vertex
	BufferUsage_Uniform      BufferUsage = C.WGPUBufferUsage_Uniform
	BufferUsage_Storage      BufferUsage = C.WGPUBufferUsage_Storage
	BufferUsage_Indirect     BufferUsage = C.WGPUBufferUsage_Indirect
	BufferUsage_QueryResolve BufferUsage = C.WGPUBufferUsage_QueryResolve
	BufferUsage_Force32      BufferUsage = C.WGPUBufferUsage_Force32
)

type ColorWriteMask C.WGPUColorWriteMask

const (
	ColorWriteMask_None    ColorWriteMask = C.WGPUColorWriteMask_None
	ColorWriteMask_Red     ColorWriteMask = C.WGPUColorWriteMask_Red
	ColorWriteMask_Green   ColorWriteMask = C.WGPUColorWriteMask_Green
	ColorWriteMask_Blue    ColorWriteMask = C.WGPUColorWriteMask_Blue
	ColorWriteMask_Alpha   ColorWriteMask = C.WGPUColorWriteMask_Alpha
	ColorWriteMask_All     ColorWriteMask = C.WGPUColorWriteMask_All
	ColorWriteMask_Force32 ColorWriteMask = C.WGPUColorWriteMask_Force32
)

type MapMode C.WGPUMapMode

const (
	MapMode_None    MapMode = C.WGPUMapMode_None
	MapMode_Read    MapMode = C.WGPUMapMode_Read
	MapMode_Write   MapMode = C.WGPUMapMode_Write
	MapMode_Force32 MapMode = C.WGPUMapMode_Force32
)

type ShaderStage C.WGPUShaderStage

const (
	ShaderStage_None     = C.WGPUShaderStage_None
	ShaderStage_Vertex   = C.WGPUShaderStage_Vertex
	ShaderStage_Fragment = C.WGPUShaderStage_Fragment
	ShaderStage_Compute  = C.WGPUShaderStage_Compute
	ShaderStage_Force32  = C.WGPUShaderStage_Force32
)

type TextureUsage C.WGPUTextureUsage

const (
	TextureUsage_None             = C.WGPUTextureUsage_None
	TextureUsage_CopySrc          = C.WGPUTextureUsage_CopySrc
	TextureUsage_CopyDst          = C.WGPUTextureUsage_CopyDst
	TextureUsage_TextureBinding   = C.WGPUTextureUsage_TextureBinding
	TextureUsage_StorageBinding   = C.WGPUTextureUsage_StorageBinding
	TextureUsage_RenderAttachment = C.WGPUTextureUsage_RenderAttachment
	TextureUsage_Force32          = C.WGPUTextureUsage_Force32
)
