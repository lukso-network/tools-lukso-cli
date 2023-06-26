package system

import (
	"os/exec"
	"runtime"
	"strconv"
)

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

func IsRoot() (isRoot bool, err error) {
	command := exec.Command("id", "-u")

	output, err := command.Output()
	if err != nil {
		return
	}

	usr, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		return
	}

	if usr == 0 {
		return true, nil
	}

	return false, nil
}
