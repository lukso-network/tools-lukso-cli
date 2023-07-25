package clients

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TekuClient struct {
	*clientBinary
}

func NewTekuClient() *TekuClient {
	return &TekuClient{
		&clientBinary{
			name:           tekuDependencyName,
			commandName:    "teku",
			baseUrl:        "https://artifacts.consensys.net/public/teku/raw/names/teku.tar.gz/versions/|TAG|/teku-|TAG|.tar.gz",
			githubLocation: tekuGithubLocation,
		},
	}
}

var Teku = NewTekuClient()

var _ ClientBinaryDependency = &TekuClient{}

func (t *TekuClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {

	return
}

func (t *TekuClient) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(t.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", t.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	_, isInstalled := os.LookupEnv(system.JavaHomeEnv) // means that JDk is not set up
	if !isInstalled {
		message := "Teku is written in Java. This means that to use it you need to have:\n" +
			"- JDK installed on your computer\n" +
			"- JAVA_HOME environmental variable set\n" +
			"Do you want to install and set up JDK along with Teku? [Y/n]\n>"

		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("❌  Aborting...")

			return
		}
	}

	response, err := http.Get(url)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", url)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			url,
			response.Status,
		)
	}

	g, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}

	defer func() {
		_ = g.Close()
	}()

	tarReader := tar.NewReader(g)

	err = untarDir(t.FilePath(), tarReader)
	if err != nil {
		return
	}

	switch isUpdate {
	case true:
		log.Infof("✅  %s updated!\n\n", t.Name())
	case false:
		log.Infof("✅  %s downloaded!\n\n", t.Name())
	}

	return

}

func (t *TekuClient) FilePath() string {
	return "./" + t.commandName
}

func untarDir(dst string, t *tar.Reader) error {
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dst, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
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
