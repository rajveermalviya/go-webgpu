package main

import (
	"fmt"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

//go:embed shader.wgsl
var shader string

func main() {
	numbers := []uint32{1, 2, 3, 4}
	const numbersSize = 4 * 4
	const numbersLength = 4

	adapter, err := wgpu.RequestAdapter(wgpu.RequestAdapterOptions{})
	if err != nil {
		panic(err)
	}

	device, err := adapter.RequestDevice(wgpu.DeviceDescriptor{
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

	shader := device.CreateShaderModule(wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})

	stagingBuffer := device.CreateBuffer(wgpu.BufferDescriptor{
		Label: "StagingBuffer",
		Usage: wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
		Size:  numbersSize,
	})

	storageBuffer := device.CreateBuffer(wgpu.BufferDescriptor{
		Label: "StorageBuffer",
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopyDst | wgpu.BufferUsage_CopySrc,
		Size:  numbersSize,
	})

	bindGroupLayout := device.CreateBindGroupLayout(wgpu.BindGroupLayoutDescriptor{
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

	bindGroup := device.CreateBindGroup(wgpu.BindGroupDescriptor{
		Label:  "Bind Group",
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{{
			Binding: 0,
			Buffer:  storageBuffer,
			Offset:  0,
			Size:    numbersSize,
		}},
	})

	pipelineLayout := device.CreatePipelineLayout(wgpu.PipelineLayoutDescriptor{
		BindGroupLayouts: []wgpu.BindGroupLayout{bindGroupLayout},
	})

	computePipeline := device.CreateComputePipeline(wgpu.ComputePipelineDescriptor{
		Layout: pipelineLayout,
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     shader,
			EntryPoint: "main",
		},
	})

	encoder := device.CreateCommandEncoder(wgpu.CommandEncoderDescriptor{
		Label: "Command Encoder",
	})

	computePass := encoder.BeginComputePass(wgpu.ComputePassDescriptor{
		Label: "Compute Pass",
	})

	computePass.SetPipeline(computePipeline)
	computePass.SetBindGroup(0, bindGroup, nil)
	computePass.Dispatch(numbersLength, 1, 1)
	computePass.EndPass()

	encoder.CopyBufferToBuffer(storageBuffer, 0, stagingBuffer, 0, numbersSize)

	queue := device.GetQueue()
	cmdBuffer := encoder.Finish(wgpu.CommandBufferDescriptor{})

	queue.WriteBuffer(storageBuffer, 0, wgpu.Uint32StoByteS(numbers))

	queue.Submit([]wgpu.CommandBuffer{cmdBuffer})

	stagingBuffer.MapAsync(wgpu.MapMode_Read, 0, numbersSize, func(status wgpu.BufferMapAsyncStatus) {
		fmt.Println("MapAsync status: ", status)
	})
	device.Poll(true)

	times := stagingBuffer.GetMappedRange(0, numbersSize)
	fmt.Println(wgpu.ByteStoUint32S(times))

	stagingBuffer.Unmap()
}
