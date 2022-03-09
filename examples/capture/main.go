package main

import (
	"image"
	"image/png"
	"os"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type BufferDimensions struct {
	width               uint64
	height              uint64
	unpaddedBytesPerRow uint64
	paddedBytesPerRow   uint64
}

func newBufferDimensions(width uint64, height uint64) BufferDimensions {
	const bytesPerPixel = unsafe.Sizeof(uint32(0))

	unpaddedBytesPerRow := width * uint64(bytesPerPixel)
	align := wgpu.CopyBytesPerRowAlignment

	paddedBytesPerRowPadding := (align - int(unpaddedBytesPerRow)%align) % align
	paddedBytesPerRow := unpaddedBytesPerRow + uint64(paddedBytesPerRowPadding)

	return BufferDimensions{
		width,
		height,
		unpaddedBytesPerRow,
		paddedBytesPerRow,
	}
}

func main() {
	width := 100
	height := 200

	adapter, err := wgpu.RequestAdapter(nil)
	if err != nil {
		// fallback to cpu
		adapter, err = wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{ForceFallbackAdapter: true})
		if err != nil {
			panic(err)
		}
	}

	device, err := adapter.RequestDevice(&wgpu.DeviceDescriptor{
		RequiredLimits: &wgpu.RequiredLimits{
			Limits: wgpu.Limits{MaxBindGroups: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	defer device.Drop()

	bufferDimensions := newBufferDimensions(uint64(width), uint64(height))

	bufferSize := bufferDimensions.paddedBytesPerRow * bufferDimensions.height
	// The output buffer lets us retrieve the data as an array
	outputBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  bufferSize,
		Usage: wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		panic(err)
	}

	textureExtent := wgpu.Extent3D{
		Width:              uint32(bufferDimensions.width),
		Height:             uint32(bufferDimensions.height),
		DepthOrArrayLayers: 1,
	}

	// The render pipeline renders data into this texture
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Size:          textureExtent,
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension_2D,
		Format:        wgpu.TextureFormat_RGBA8UnormSrgb,
		Usage:         wgpu.TextureUsage_RenderAttachment | wgpu.TextureUsage_CopySrc,
	})
	if err != nil {
		panic(err)
	}

	// Set the background to be red
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		panic(err)
	}

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:       texture.CreateView(nil),
			LoadOp:     wgpu.LoadOp_Clear,
			StoreOp:    wgpu.StoreOp_Store,
			ClearColor: wgpu.Color_Red,
		}},
	})
	renderPass.EndPass()

	// Copy the data from the texture to the buffer
	encoder.CopyTextureToBuffer(
		texture.AsImageCopy(),
		&wgpu.ImageCopyBuffer{
			Buffer: outputBuffer,
			Layout: wgpu.TextureDataLayout{
				Offset:      0,
				BytesPerRow: uint32(bufferDimensions.paddedBytesPerRow),
			},
		},
		&textureExtent,
	)

	queue := device.GetQueue()
	queue.Submit(encoder.Finish(nil))

	outputBuffer.MapAsync(wgpu.MapMode_Read, 0, bufferSize, func(status wgpu.BufferMapAsyncStatus) {
		if status != wgpu.BufferMapAsyncStatus_Success {
			panic("failed to map buffer")
		}
	})
	device.Poll(true)
	defer outputBuffer.Unmap()

	data := outputBuffer.GetMappedRange(0, bufferSize)

	// Save png
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	imageEncoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = imageEncoder.Encode(f, &image.NRGBA{
		Pix:    data,
		Stride: int(bufferDimensions.paddedBytesPerRow),
		Rect:   image.Rect(0, 0, width, height),
	})
	if err != nil {
		panic(err)
	}
}
