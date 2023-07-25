# go-webgpu

go bindings for [`wgpu-native`](https://github.com/gfx-rs/wgpu-native), a cross-platform, safe, graphics api. it runs natively on vulkan, metal, d3d12 and opengles.

for more info check:
- [webgpu](https://gpuweb.github.io/gpuweb/)
- [wgsl](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu-native](https://github.com/webgpu-native/webgpu-headers)

included static libs are built via [github actions](./.github/workflows/build-wgpu.yml).

## Examples

|[boids][b]|[cube][c]|[triangle][t]|
:-:|:-:|:-:
| [![b-i]][b] | [![c-i]][c] | [![t-i]][t] |

[b-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/boids/image-msaa.png
[b]: https://github.com/rajveermalviya/go-webgpu-examples/tree/main/boids
[c-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/cube/image-msaa.png
[c]: https://github.com/rajveermalviya/go-webgpu-examples/tree/main/cube
[t-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/triangle/image-msaa.png
[t]: https://github.com/rajveermalviya/go-webgpu-examples/tree/main/triangle

you can check out all the examples in [go-webgpu-examples repo](https://github.com/rajveermalviya/go-webgpu-examples)
