package file

import "os"

var (
	filePerms         = os.FileMode(0o0750)
	_         Manager = &manager{}
)

type Manager interface {
	Write(body []byte, dst string, perm os.FileMode) (err error)
	Mkdir(dst string, perm os.FileMode) (err error)
}

type manager struct{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Write(body []byte, dst string, perm os.FileMode) (err error) {
	return os.WriteFile(dst, body, filePerms)
}

func (m *manager) Mkdir(dst string, perm os.FileMode) (err error) {
	return os.MkdirAll(dst, perm)
}
