package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

// Indicates a uint32 overflow in an intermediate Collatz value
const OVERFLOW = 0xffffffff

func main() {
	numbers := []uint32{1, 2, 3, 4}

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

	shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: shader,
		},
	})
	if err != nil {
		panic(err)
	}
	defer shaderModule.Drop()

	size := uint64(len(numbers)) * uint64(unsafe.Sizeof(uint32(0)))

	stagingBuffer, err := device.CreateBuffer(&wgpu.BufferDescriptor{
		Size:             size,
		Usage:            wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
		MappedAtCreation: false,
	})
	if err != nil {
		panic(err)
	}
	defer stagingBuffer.Drop()

	storageBuffer, err := device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Storage Buffer",
		Contents: wgpu.ToBytes(numbers),
		Usage: wgpu.BufferUsage_Storage |
			wgpu.BufferUsage_CopyDst |
			wgpu.BufferUsage_CopySrc,
	})
	if err != nil {
		panic(err)
	}
	defer storageBuffer.Drop()

	computePipeline, err := device.CreateComputePipeline(&wgpu.ComputePipelineDescriptor{
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     shaderModule,
			EntryPoint: "main",
		},
	})
	if err != nil {
		panic(err)
	}
	defer computePipeline.Drop()

	bindGroupLayout := computePipeline.GetBindGroupLayout(0)
	defer bindGroupLayout.Drop()

	bindGroup, err := device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout: bindGroupLayout,
		Entries: []wgpu.BindGroupEntry{{
			Binding: 0,
			Buffer:  storageBuffer,
			Size:    wgpu.WholeSize,
		}},
	})
	if err != nil {
		panic(err)
	}
	defer bindGroup.Drop()

	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		panic(err)
	}

	computePass := encoder.BeginComputePass(nil)
	computePass.SetPipeline(computePipeline)
	computePass.SetBindGroup(0, bindGroup, nil)
	computePass.DispatchWorkgroups(uint32(len(numbers)), 1, 1)
	computePass.End()

	encoder.CopyBufferToBuffer(storageBuffer, 0, stagingBuffer, 0, size)

	queue.Submit(encoder.Finish(nil))

	var status wgpu.BufferMapAsyncStatus
	stagingBuffer.MapAsync(wgpu.MapMode_Read, 0, size, func(s wgpu.BufferMapAsyncStatus) {
		status = s
	})

	device.Poll(true, nil)

	if status != wgpu.BufferMapAsyncStatus_Success {
		panic(status)
	}

	steps := make([]uint32, len(numbers))
	{
		data := stagingBuffer.GetMappedRange(0, uint(size))

		copy(steps, wgpu.FromBytes[uint32](data))

		data = nil
		stagingBuffer.Unmap()
	}

	dispSteps := mapSlice(steps, func(e uint32) string {
		if e == OVERFLOW {
			return "OVERFLOW"
		}
		return strconv.FormatUint(uint64(e), 10)
	})

	fmt.Printf("Steps: [%s]\n", strings.Join(dispSteps, ", "))
}

func mapSlice[E1 any, E2 any](s []E1, f func(e E1) E2) []E2 {
	rs := make([]E2, len(s))
	for i, e := range s {
		rs[i] = f(e)
	}
	return rs
}
