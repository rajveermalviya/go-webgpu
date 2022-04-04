package wgpu

const (
	CopyBytesPerRowAlignment    = 256
	QueryResolveBufferAlignment = 256
	CopyBufferAlignment         = 4
	MapAlignment                = 8
	VertexStrideAlignment       = 4
	PushConstantAlignment       = 4
	QuerySetMaxQueries          = 8192
	QuerySize                   = 8
)

var (
	Color_Transparent = Color{0, 0, 0, 0}
	Color_Black       = Color{0, 0, 0, 1}
	Color_White       = Color{1, 1, 1, 1}
	Color_Red         = Color{1, 0, 0, 1}
	Color_Green       = Color{0, 1, 0, 1}
	Color_Blue        = Color{0, 0, 1, 1}

	BlendComponent_Replace = BlendComponent{
		SrcFactor: BlendFactor_One,
		DstFactor: BlendFactor_Zero,
		Operation: BlendOperation_Add,
	}
	BlendComponent_Over = BlendComponent{
		SrcFactor: BlendFactor_One,
		DstFactor: BlendFactor_OneMinusSrcAlpha,
		Operation: BlendOperation_Add,
	}

	BlendState_Replace = BlendState{
		Color: BlendComponent_Replace,
		Alpha: BlendComponent_Replace,
	}
	BlendState_AlphaBlending = BlendState{
		Color: BlendComponent{
			SrcFactor: BlendFactor_SrcAlpha,
			DstFactor: BlendFactor_OneMinusSrcAlpha,
			Operation: BlendOperation_Add,
		},
		Alpha: BlendComponent_Over,
	}
	BlendState_PremultipliedAlphaBlending = BlendState{
		Color: BlendComponent_Over,
		Alpha: BlendComponent_Over,
	}
)

func (v VertexFormat) Size() uint64 {
	switch v {
	case VertexFormat_Uint8x2,
		VertexFormat_Sint8x2,
		VertexFormat_Unorm8x2,
		VertexFormat_Snorm8x2:
		return 2

	case VertexFormat_Uint8x4,
		VertexFormat_Sint8x4,
		VertexFormat_Unorm8x4,
		VertexFormat_Snorm8x4,
		VertexFormat_Uint16x2,
		VertexFormat_Sint16x2,
		VertexFormat_Unorm16x2,
		VertexFormat_Snorm16x2,
		VertexFormat_Float16x2,
		VertexFormat_Float32,
		VertexFormat_Uint32,
		VertexFormat_Sint32:
		return 4

	case VertexFormat_Uint16x4,
		VertexFormat_Sint16x4,
		VertexFormat_Unorm16x4,
		VertexFormat_Snorm16x4,
		VertexFormat_Float16x4,
		VertexFormat_Float32x2,
		VertexFormat_Uint32x2,
		VertexFormat_Sint32x2:
		return 8

	case VertexFormat_Float32x3,
		VertexFormat_Uint32x3,
		VertexFormat_Sint32x3:
		return 12

	case VertexFormat_Float32x4,
		VertexFormat_Uint32x4,
		VertexFormat_Sint32x4:
		return 16

	default:
		return 0
	}
}
