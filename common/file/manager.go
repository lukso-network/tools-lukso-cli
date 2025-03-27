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

func (m *manager) Write(body []byte, dst string, perm os.FileMode) (err error) {
	f, err := os.OpenFile(dst, os.O_RDWR, filePerms)
	if err != nil {
		return
	}

	_, err = f.Write(body)

	return
}

func (m *manager) Mkdir(dst string, perm os.FileMode) (err error) {
	return os.MkdirAll(dst, perm)
}
