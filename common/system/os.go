package system

import (
	"bytes"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	Ubuntu  = "linux"
	Macos   = "darwin"
	Windows = "windows"

	UnixBinDir = "/usr/local/bin"

	JavaHomeEnv = "JAVA_HOME"
)

var (
	Arch           = runtime.GOARCH
	Os             = runtime.GOOS
	SupportedArchs = []string{"x86_64", "aarch64", "arm", "i686"}
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

func GetArch() (arch string) {
	fallback := func() {
		log.Info("⚠️  Unknown OS detected: proceeding with x86_64 as a default arch")
		arch = "x86_64"
	}

	switch Os {
	case Ubuntu, Macos:
		buf := new(bytes.Buffer)

		uname := exec.Command("uname", "-m")
		uname.Stdout = buf

		err := uname.Run()
		if err != nil {
			fallback()

			break
		}

		arch = strings.Trim(buf.String(), "\n\t ")

	default:
		fallback()
	}

	if !slices.Contains(SupportedArchs, arch) {
		fallback()
	}

	return
}
