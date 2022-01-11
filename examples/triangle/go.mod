module github.com/rajveermalviya/go-webgpu/examples/triangle

go 1.17

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211213063430-748e38ca8aec
	github.com/rajveermalviya/go-webgpu/wgpu v0.0.0-20220109190449-e2909fe328c3
)

replace github.com/rajveermalviya/go-webgpu/wgpu => ../../wgpu/
