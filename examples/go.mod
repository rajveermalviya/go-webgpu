module github.com/rajveermalviya/go-webgpu/examples

go 1.18

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20220712193148-63cf1f4ef61f
	github.com/rajveermalviya/go-webgpu/wgpu v0.1.1
)

require golang.org/x/sys v0.0.0-20220731174439-a90be440212d // indirect

retract [v0.0.0, v0.1.1]
