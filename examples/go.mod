module github.com/rajveermalviya/go-webgpu/examples

go 1.18

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211213063430-748e38ca8aec
	github.com/go-gl/mathgl v1.0.0
	github.com/rajveermalviya/go-webgpu/wgpu v0.0.0-20220225063355-7ffafaa77e15
)

require golang.org/x/image v0.0.0-20190321063152-3fc05d484e9f // indirect

replace github.com/rajveermalviya/go-webgpu/wgpu => ../wgpu/
