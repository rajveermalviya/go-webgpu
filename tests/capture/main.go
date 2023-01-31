package main

import (
	"image"
	"image/png"
	"os"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"
)

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	switch os.Getenv("WGPU_LOG_LEVEL") {
	case "OFF":
		wgpu.SetLogLevel(wgpu.LogLevel_Off)
	case "ERROR":
		wgpu.SetLogLevel(wgpu.LogLevel_Error)
	case "WARN":
		wgpu.SetLogLevel(wgpu.LogLevel_Warn)
	case "INFO":
		wgpu.SetLogLevel(wgpu.LogLevel_Info)
	case "DEBUG":
		wgpu.SetLogLevel(wgpu.LogLevel_Debug)
	case "TRACE":
		wgpu.SetLogLevel(wgpu.LogLevel_Trace)
	}
}

type BufferDimensions struct {
	width               uint64
	height              uint64
	unpaddedBytesPerRow uint64
	paddedBytesPerRow   uint64
}

func newBufferDimensions(width uint64, height uint64) BufferDimensions {
	const bytesPerPixel = unsafe.Sizeof(uint32(0))
	unpaddedBytesPerRow := width * uint64(bytesPerPixel)
	align := uint64(wgpu.CopyBytesPerRowAlignment)
	paddedBytesPerRowPadding := (align - unpaddedBytesPerRow%align) % align
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

	instance := wgpu.CreateInstance(nil)
	defer instance.Drop()

	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
	})
	if err != nil {
		panic(err)
	}
	defer adapter.Drop()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		panic(err)
	}
	defer device.Drop()
	queue := device.GetQueue()

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
	defer outputBuffer.Drop()

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
	defer texture.Drop()

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
			ClearValue: wgpu.Color_Red,
		}},
	})
	renderPass.End()

	// Copy the data from the texture to the buffer
	encoder.CopyTextureToBuffer(
		texture.AsImageCopy(),
		&wgpu.ImageCopyBuffer{
			Buffer: outputBuffer,
			Layout: wgpu.TextureDataLayout{
				Offset:       0,
				BytesPerRow:  uint32(bufferDimensions.paddedBytesPerRow),
				RowsPerImage: wgpu.CopyStrideUndefined,
			},
		},
		&textureExtent,
	)

	index := queue.Submit(encoder.Finish(nil))

	outputBuffer.MapAsync(wgpu.MapMode_Read, 0, bufferSize, func(status wgpu.BufferMapAsyncStatus) {
		if status != wgpu.BufferMapAsyncStatus_Success {
			panic("failed to map buffer")
		}
	})
	defer outputBuffer.Unmap()

	device.Poll(true, &wgpu.WrappedSubmissionIndex{
		Queue:           queue,
		SubmissionIndex: index,
	})

	data := outputBuffer.GetMappedRange(0, uint(bufferSize))

	// Save png
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

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
