module github.com/rajveermalviya/go-webgpu/tests

go 1.20

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b
	github.com/rajveermalviya/go-webgpu/wgpu v0.17.1
	github.com/rajveermalviya/go-webgpu/wgpuext/glfw v0.0.0-00010101000000-000000000000
)

replace github.com/rajveermalviya/go-webgpu/wgpu => ../wgpu

replace github.com/rajveermalviya/go-webgpu/wgpuext/glfw => ../wgpuext/glfw
