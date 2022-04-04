# go-webgpu

Go bindings for [wgpu-native](https://github.com/gfx-rs/wgpu-native).

For more info check:
- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu-native](https://github.com/webgpu-native/webgpu-headers)

## cgo

- on windows cgo **is not** used (i.e. works with `CGO_ENABLED=0`). so only Go compiler is needed.

- on unix (linux, darwin, android) cgo **is** used, so you will need C toolchain (`gcc` or `clang`) installed.

Included static libs & windows dll are built via [Github Actions](./.github/workflows/build-wgpu.yml).

## Check out examples

### [compute](./examples/compute/main.go)

```shell
go run github.com/rajveermalviya/go-webgpu/examples/compute@latest
```

### [capture](./examples/capture/main.go)

Creates `./image.png` with all pixels red and size 100x200

```shell
go run github.com/rajveermalviya/go-webgpu/examples/capture@latest
```

### [triangle](./examples/triangle/main.go)

This example uses [go-glfw](https://github.com/go-gl/glfw) so it will use cgo on **_all platforms_**, you will also need
[some libraries installed](https://github.com/go-gl/glfw#installation) to run the example.

```shell
go run github.com/rajveermalviya/go-webgpu/examples/triangle@latest
```
