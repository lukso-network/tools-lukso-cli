package system

import "runtime"

const (
	Ubuntu  = "linux"
	Macos   = "darwin"
	Windows = "windows"

	UnixBinDir = "/usr/local/bin"
)

var (
	Arch = runtime.GOARCH
	Os   = runtime.GOOS
)
