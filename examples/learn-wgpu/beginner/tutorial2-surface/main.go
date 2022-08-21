package main

import (
	"fmt"
	"strings"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type State struct {
	surface   *wgpu.Surface
	swapChain *wgpu.SwapChain
	device    *wgpu.Device
	queue     *wgpu.Queue
	config    *wgpu.SwapChainDescriptor
	size      dpi.PhysicalSize[uint32]
}

func InitState(window display.Window) (*State, error) {
	size := window.InnerSize()

	surface := wgpu.CreateSurface(getSurfaceDescriptor(window))

	adaper, err := wgpu.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface,
	})
	if err != nil {
		return nil, err
	}
	device, err := adaper.RequestDevice(nil)
	if err != nil {
		return nil, err
	}
	queue := device.GetQueue()

	config := &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      surface.GetPreferredFormat(adaper),
		Width:       size.Width,
		Height:      size.Height,
		PresentMode: wgpu.PresentMode_Fifo,
	}
	swapChain, err := device.CreateSwapChain(surface, config)
	if err != nil {
		return nil, err
	}

	return &State{
		surface:   surface,
		swapChain: swapChain,
		device:    device,
		queue:     queue,
		config:    config,
		size:      size,
	}, nil
}

func (s *State) Resize(newSize dpi.PhysicalSize[uint32]) {
	if newSize.Width > 0 && newSize.Height > 0 {
		s.size = newSize
		s.config.Width = newSize.Width
		s.config.Height = newSize.Height

		var err error
		s.swapChain, err = s.device.CreateSwapChain(s.surface, s.config)
		if err != nil {
			panic(err)
		}
	}
}

func (s *State) Render() error {
	view, err := s.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer view.Drop()

	encoder, err := s.device.CreateCommandEncoder(nil)
	if err != nil {
		return err
	}

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:   view,
			LoadOp: wgpu.LoadOp_Clear,
			ClearValue: wgpu.Color{
				R: 0.1,
				G: 0.2,
				B: 0.3,
				A: 1.0,
			},
			StoreOp: wgpu.StoreOp_Store,
		}},
	})
	renderPass.End()

	s.queue.Submit(encoder.Finish(nil))
	s.swapChain.Present()

	return nil
}

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}

	s, err := InitState(w)
	if err != nil {
		panic(err)
	}

	rerender := true
	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		s.Resize(dpi.PhysicalSize[uint32]{
			Width:  physicalWidth,
			Height: physicalHeight,
		})
		rerender = true
	})

	w.SetCloseRequestedCallback(func() {
		d.Destroy()
	})

	for {
		if rerender {
			rerender = false

			err := s.Render()
			if err != nil {
				errstr := err.Error()
				fmt.Println(errstr)

				switch {
				case strings.Contains(errstr, "Lost"):
					s.Resize(s.size)
				case strings.Contains(errstr, "Outdated"):
					s.Resize(s.size)
				case strings.Contains(errstr, "Timeout"):
				default:
					panic(err)
				}
			}
		}

		if !d.Wait() {
			break
		}
	}
}
