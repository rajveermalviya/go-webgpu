//go:build windows && 386

package wgpu

import _ "embed"

//go:embed lib/windows/386/wgpu_native.dll.gz
var libwgpuDllCompressed []byte

//go:embed lib/windows/386/wgpu_native.dll.sum
var libwgpuDllSha256 string
