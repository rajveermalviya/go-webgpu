package wgpu

import (
	"fmt"
	"strconv"
)

type Version uint32

func (v Version) String() string {
	return "0x" + strconv.FormatUint(uint64(v), 8)
}

func (p *Device) getErr() (err error) {
	select {
	case err = <-p.errChan:
	default:
	}
	return
}

func (p *Device) storeErr(typ ErrorType, message string) {
	err := &Error{Type: typ, Message: fmt.Sprint(message)}
	select {
	case p.errChan <- err:
	default:
		var prevErr *Error

		select {
		case prevErr = <-p.errChan:
		default:
		}

		var str string
		if prevErr != nil {
			str = fmt.Sprintf("go-webgpu: previous uncaptured error: %s\n\n", prevErr.Error())
		}
		str += fmt.Sprintf("go-webgpu: current uncaptured error: %s\n\n", err.Error())
		panic(str)
	}
}

func (p *Texture) AsImageCopy() *ImageCopyTexture {
	return &ImageCopyTexture{
		Texture:  p,
		MipLevel: 0,
		Origin:   Origin3D{},
		Aspect:   TextureAspect_All,
	}
}
