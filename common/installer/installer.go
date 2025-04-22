package installer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
)

// Installer is used for fetching all kinds of external files - clients, configs etc.
// It's primarly used to encapsulate fetching logic for easier installations and mocking.
type Installer interface {
	// Fetch uses HTTP GET to fetch data from the url and returns the response body as bytes.
	Fetch(url string) (body []byte, err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes it to dest file. It also creates any path to it if necessary.
	InstallFile(url, dest string) (err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes decompressed tar archive to dest. It also creates any path to it if necessary.
	InstallTar(url, dest string) (err error)

	// InstallFile uses HTTP GET to fetch data from the url and writes decompressed zip archive to dest. It also creates any path to it if necessary.
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
	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, errors.ErrNotFound
	}

	if http.StatusOK != response.StatusCode {
		return nil, fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			url,
			response.Status,
		)
	}

	var (
		responseReader io.Reader = response.Body
		buf                      = new(bytes.Buffer)
	)
	_, err = buf.ReadFrom(responseReader)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

func (i *installer) InstallFile(url, dest string) (err error) {
	body, err := i.Fetch(url)
	if err != nil {
		return
	}

	err = i.file.Mkdir(utils.TruncateFileFromDir(dest), common.ConfigPerms)
	if err != nil {
		return
	}

	err = i.file.Write(dest, body, common.ConfigPerms)
	if err != nil && strings.Contains(err.Error(), "Permission denied") {
		return errors.ErrNeedsRoot
	}

	if err != nil {
		return
	}

	return
}

func (i *installer) InstallTar(url, dest string) (err error) {
	return
}

func (i *installer) InstallZip(url, dest string) (err error) {
	return
}
