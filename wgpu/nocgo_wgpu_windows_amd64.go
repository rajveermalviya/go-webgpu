//go:build windows && amd64

package wgpu

import _ "embed"

//go:embed lib/windows/amd64/wgpu_native.dll.gz
var libwgpuDllCompressed []byte
