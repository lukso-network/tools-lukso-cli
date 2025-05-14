package installer

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	// If the overwrite is set to true, it will overwrite any existing files.
	InstallFile(url, dest string, overwrite bool) (err error)

	// InstallTar uses HTTP GET to fetch data from the url and writes decompressed tar archive to dest. It also creates any path to it if necessary.
	InstallTar(url, dest, archiveName, pattern string, overwrite bool) (err error)

	// InstallZip uses HTTP GET to fetch data from the url and writes decompressed zip archive to dest. It also creates any path to it if necessary.
	InstallZip(url, dest string, overwrite bool) (err error)
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
			"invalid response when downloading on file url: %s. Response code: %s",
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

func (i *installer) InstallFile(url, dest string, overwrite bool) (err error) {
	if !overwrite {
		return errors.ErrFileExists
	}

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

// InstallTar
func (i *installer) InstallTar(url, dest, archiveName, pattern string, overwrite bool) (err error) {
	if !overwrite {
		return errors.ErrFileExists
	}

	var (
		b   []byte
		buf *bytes.Reader
		g   *gzip.Reader
		t   *tar.Reader
	)

	b, err = i.Fetch(url)
	if err != nil {
		return
	}

	buf = bytes.NewReader(b)
	g, err = gzip.NewReader(buf)

	t = tar.NewReader(g)
	if err != nil {
		return
	}

	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var (
			path       string
			headerName = header.Name
		)

		// for the sake of compatibility with updated versions remove the tag from the tarred file - teku/teku-xx.x.x => teku/teku, same with jdk
		if strings.Contains(header.Name, pattern) {
			headerName = replaceRootFolderName(header.Name, archiveName)
		}

		path = filepath.Join(dest, headerName)

		info := header.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, info.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *installer) InstallZip(url, dest string, overwrite bool) (err error) {
	if !overwrite {
		return errors.ErrFileExists
	}

	var (
		b   []byte
		buf *bytes.Reader
		r   *zip.Reader
	)

	b, err = i.Fetch(url)
	if err != nil {
		return
	}

	buf = bytes.NewReader(b)
	r, err = zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		return
	}

	for _, header := range r.File {
		path := filepath.Join(dest, header.Name)

		info := header.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, info.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		defer file.Close()

		readFile, err := header.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, readFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func replaceRootFolderName(folder, targetRootName string) (path string) {
	splitHeader := strings.Split(folder, "/") // this assumes no / at the beginning of folder - not the case in tarred files we are interested in

	switch len(splitHeader) {
	case 0:
		return
	case 1:
		return targetRootName
	default:
		break
	}

	strippedHeader := splitHeader[1:]
	splitResolvedPath := append([]string{targetRootName}, strippedHeader...)

	return strings.Join(splitResolvedPath, "/")
}
