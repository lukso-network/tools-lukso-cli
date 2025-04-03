package installer

import (
	"net/http"

	"github.com/lukso-network/tools-lukso-cli/common/file"
)

// Installer is used for fetching all kinds of external files - clients, configs etc.
// It's primarly used to encapsulate fetching logic for easier installations and mocking.
type Installer interface {
	// Fetch uses HTTP GET to fetch data from the url and returns the response body as bytes.
	Fetch(url string) (body []byte, err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes it to dest.
	InstallFile(url, dest string) (err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes it to dest.
	InstallTar(url, dest string) (err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes it to dest.
	InstallZip(url, dest string) (err error)
}

type installer struct {
	httpClient *http.Client
	file       file.Manager
}

var _ Installer = &installer{}

func NewInstaller(mng file.Manager) Installer {
	return &installer{
		httpClient: http.DefaultClient,
		file:       mng,
	}
}

func (i *installer) Fetch(url string) (body []byte, err error) {
	return
}

func (i *installer) InstallFile(url, dest string) (err error) {
	return
}

func (i *installer) InstallTar(url, dest string) (err error) {
	return
}

func (i *installer) InstallZip(url, dest string) (err error) {
	return
}
