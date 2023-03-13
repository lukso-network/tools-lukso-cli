package pid

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

var FileDir = "/var/run/lukso"

func Exists(path string) bool {
	pidVal, err := Load(path)
	if err != nil {
		return false
	}

	p, err := os.FindProcess(pidVal)
	if err != nil {
		return false
	}

	err = p.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}

func Create(path string, pid int) error {
	strPid := fmt.Sprintf("%d", pid)

	return os.WriteFile(path, []byte(strPid), os.ModePerm)
}

func Kill(path string, pid int) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = p.Signal(syscall.Signal(0))
	if err != nil {
		return err
	}

	return p.Kill()
}

func Load(path string) (pid int, err error) {
	strPid, err := os.ReadFile(path)
	if err != nil {
		return
	}

	return strconv.Atoi(string(strPid))
}
