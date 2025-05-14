package pid

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/lukso-network/tools-lukso-cli/common/file"
)

var FileDir = "/tmp" // until a script for /tmp/lukso or other dir is provided

type Pid interface {
	Exists(path string) bool
	Create(path string, pid int) error
	Kill(path string, pid int) error
	Load(path string) (pid int, err error)
}

type pid struct {
	file file.Manager
}

var _ Pid = pid{}

func NewPid(file file.Manager) pid {
	return pid{
		file: file,
	}
}

func (p pid) Exists(path string) bool {
	return Exists(path)
}

func (p pid) Create(path string, pid int) error {
	return Create(path, pid)
}

func (p pid) Kill(path string, pid int) error {
	return Kill(path, pid)
}

func (p pid) Load(path string) (pid int, err error) {
	return Load(path)
}

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

	return err == nil
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
		return nil //  we just return, if we get an error it means process is already dead - pid is cleared, we're good to go
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
