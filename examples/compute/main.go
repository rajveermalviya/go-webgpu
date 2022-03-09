package main

import (
	"fmt"
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

//go:embed shader.wgsl
var shader string

func main() {
	numbers := []uint32{1, 2, 3, 4}
	numbersSize := len(numbers) * int(unsafe.Sizeof(uint32(0)))
	numbersLength := len(numbers)

	adapter, err := wgpu.RequestAdapter(nil)
	if err != nil {
		// fallback to cpu
		adapter, err = wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{ForceFallbackAdapter: true})
		if err != nil {
			panic(err)
		}
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

	shader, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})
	if err != nil {
		panic(err)
	}

	stagingBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "StagingBuffer",
		Usage: wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
		Size:  uint64(numbersSize),
	})
	if err != nil {
		panic(err)
	}

	storageBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "StorageBuffer",
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopyDst | wgpu.BufferUsage_CopySrc,
		Size:  uint64(numbersSize),
	})
	if err != nil {
		panic(err)
	}

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

	pipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		BindGroupLayouts: []*wgpu.BindGroupLayout{bindGroupLayout},
	})
	if err != nil {
		panic(err)
	}

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
	computePass.Dispatch(uint32(numbersLength), 1, 1)
	computePass.EndPass()

	encoder.CopyBufferToBuffer(storageBuffer, 0, stagingBuffer, 0, uint64(numbersSize))

	queue := device.GetQueue()
	cmdBuffer := encoder.Finish(nil)
	queue.WriteBuffer(storageBuffer, 0, wgpu.ToBytes(numbers))
	queue.Submit(cmdBuffer)

	stagingBuffer.MapAsync(wgpu.MapMode_Read, 0, uint64(numbersSize), func(status wgpu.BufferMapAsyncStatus) {
		fmt.Println("MapAsync status:", status)
	})
	device.Poll(true)

	times := stagingBuffer.GetMappedRange(0, uint64(numbersSize))
	fmt.Println(wgpu.FromBytes(times, uint32(0)))

	stagingBuffer.Unmap()
}
