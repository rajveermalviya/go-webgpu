package wgpu

type BufferInitDescriptor struct {
	Label    string
	Contents []byte
	Usage    BufferUsage
}

func (p *Device) CreateBufferInit(descriptor BufferInitDescriptor) *Buffer {
	if len(descriptor.Contents) == 0 {
		desc := BufferDescriptor{
			Label:            descriptor.Label,
			Size:             0,
			Usage:            descriptor.Usage,
			MappedAtCreation: false,
		}

		return p.CreateBuffer(desc)
	}

	unpaddedSize := len(descriptor.Contents)
	const alignMask = CopyBufferAlignment - 1
	paddedSize := max(((unpaddedSize + alignMask) & ^alignMask), CopyBufferAlignment)

	desc := BufferDescriptor{
		Label:            descriptor.Label,
		Size:             uint64(paddedSize),
		Usage:            descriptor.Usage,
		MappedAtCreation: true,
	}

	buffer := p.CreateBuffer(desc)
	buf := buffer.GetMappedRange(0, uint64(paddedSize))
	copy(buf, descriptor.Contents)
	buffer.Unmap()

	return buffer
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
