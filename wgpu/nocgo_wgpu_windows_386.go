//go:build windows && 386

package wgpu

import _ "embed"

//go:embed lib/windows/386/libwgpu.dll.gz
var libwgpuDllCompressed []byte
