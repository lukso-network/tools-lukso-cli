package file

import "os"

var (
	filePerms         = os.FileMode(0o0750)
	_         Manager = &manager{}
)

type Manager interface {
	Create(dst string) error
	Write(dst string, body []byte, perm os.FileMode) (err error)
	Mkdir(dst string, perm os.FileMode) error
	Exists(path string) bool
}

type manager struct{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Create(dst string) (err error) {
	_, err = os.Create(dst)
	return
}

func (m *manager) Write(dst string, body []byte, perm os.FileMode) error {
	return os.WriteFile(dst, body, filePerms)
}

func (m *manager) Mkdir(dst string, perm os.FileMode) error {
	return os.MkdirAll(dst, perm)
}

func (m *manager) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
