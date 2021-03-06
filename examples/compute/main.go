package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
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

//go:embed shader.wgsl
var shader string

func main() {
	numbers := []uint32{1, 2, 3, 4}
	numbersSize := len(numbers) * int(unsafe.Sizeof(uint32(0)))

	adapter, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
	})
	if err != nil {
		panic(err)
	}

	device, err := adapter.RequestDevice(&wgpu.DeviceDescriptor{
		DeviceExtras: &wgpu.DeviceExtras{
			Label: "Device",
		},
		RequiredLimits: &wgpu.RequiredLimits{
			Limits: wgpu.Limits{
				MaxBindGroups: 1,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	defer device.Drop()
	queue := device.GetQueue()

	shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})
	if err != nil {
		panic(err)
	}
	defer shader.Drop()

	stagingBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "StagingBuffer",
		Usage: wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
		Size:  uint64(numbersSize),
	})
	if err != nil {
		panic(err)
	}
	defer stagingBuffer.Drop()

	storageBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "StorageBuffer",
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopyDst | wgpu.BufferUsage_CopySrc,
		Size:  uint64(numbersSize),
	})
	if err != nil {
		panic(err)
	}
	defer storageBuffer.Drop()

	bindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{{
			Binding:    0,
			Visibility: wgpu.ShaderStage_Compute,
			Buffer: wgpu.BufferBindingLayout{
				Type: wgpu.BufferBindingType_Storage,
			},
			Sampler: wgpu.SamplerBindingLayout{
				Type: wgpu.SamplerBindingType_Undefined,
			},
			Texture: wgpu.TextureBindingLayout{
				SampleType: wgpu.TextureSampleType_Undefined,
			},
			StorageTexture: wgpu.StorageTextureBindingLayout{
				Access: wgpu.StorageTextureAccess_Undefined,
			},
		}},
	})
	if err != nil {
		panic(err)
	}
	defer bindGroupLayout.Drop()

	bindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  "Bind Group",
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{{
			Binding: 0,
			Buffer:  storageBuffer,
			Offset:  0,
			Size:    uint64(numbersSize),
		}},
	})
	if err != nil {
		panic(err)
	}
	defer bindGroup.Drop()

	pipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		BindGroupLayouts: []*wgpu.BindGroupLayout{bindGroupLayout},
	})
	if err != nil {
		panic(err)
	}
	defer pipelineLayout.Drop()

	computePipeline, err := device.CreateComputePipeline(&wgpu.ComputePipelineDescriptor{
		Layout: pipelineLayout,
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     shader,
			EntryPoint: "main",
		},
	})
	if err != nil {
		panic(err)
	}
	defer computePipeline.Drop()

	encoder, err := device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Command Encoder",
	})
	if err != nil {
		panic(err)
	}

	computePass := encoder.BeginComputePass(&wgpu.ComputePassDescriptor{
		Label: "Compute Pass",
	})

	computePass.SetPipeline(computePipeline)
	computePass.SetBindGroup(0, bindGroup, nil)
	computePass.DispatchWorkgroups(uint32(len(numbers)), 1, 1)
	computePass.End()

	encoder.CopyBufferToBuffer(storageBuffer, 0, stagingBuffer, 0, uint64(numbersSize))

	cmdBuffer := encoder.Finish(nil)
	queue.WriteBuffer(storageBuffer, 0, wgpu.ToBytes(numbers))
	index := queue.Submit(cmdBuffer)

	stagingBuffer.MapAsync(wgpu.MapMode_Read, 0, uint64(numbersSize), func(status wgpu.BufferMapAsyncStatus) {
		fmt.Println("MapAsync status:", status)
	})
	defer stagingBuffer.Unmap()

	device.Poll(true, &wgpu.WrappedSubmissionIndex{
		Queue:           queue,
		SubmissionIndex: wgpu.SubmissionIndex(index),
	})

	times := stagingBuffer.GetMappedRange(0, uint64(numbersSize))
	fmt.Println(wgpu.FromBytes[uint32](times))
}
