package wgpu

func (v ErrorType) String() string {
	switch v {
	case ErrorType_NoError:
		return "NoError"
	case ErrorType_Validation:
		return "Validation"
	case ErrorType_OutOfMemory:
		return "OutOfMemory"
	case ErrorType_DeviceLost:
		return "DeviceLost"

	case ErrorType_Unknown:
		fallthrough
	default:
		return "Unknown"
	}
}

type Error struct {
	Type    ErrorType
	Message string
}

func (v *Error) Error() string {
	return v.Type.String() + " : " + v.Message
}
