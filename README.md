# go-webgpu

Go bindings for [wgpu-native](https://github.com/gfx-rs/wgpu-native).

For more info check:
- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu-native](https://github.com/webgpu-native/webgpu-headers)

## Uses cgo

To use this module or run any examples you will need C toolchain (`gcc` or `clang`) installed first.

Included static lib binaries are built via [Github Actions](https://github.com/rajveermalviya/wgpu-native/actions), but manually commited (have to look into automating it).

## Check out examples

### [compute](./examples/compute/main.go)

```shell
go run github.com/rajveermalviya/go-webgpu/examples/compute@latest
```
### [triangle](./examples/triangle/main.go)

This uses glfw so you will need [some libraries installed](https://github.com/go-gl/glfw#installation) to run the example.

```shell
go run github.com/rajveermalviya/go-webgpu/examples/triangle@latest
```
### [capture](./examples/capture/main.go)

Creates `./image.png` with all pixels red and size 100x200

```shell
go run github.com/rajveermalviya/go-webgpu/examples/capture@latest
```
