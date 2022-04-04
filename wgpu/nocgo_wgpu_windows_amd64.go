//go:build windows && amd64

package wgpu

import _ "embed"

//go:embed lib/windows/amd64/libwgpu.dll.gz
var libwgpuDllCompressed []byte
