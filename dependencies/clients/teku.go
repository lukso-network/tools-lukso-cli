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

const (
	tekuDepsFolder = "teku"        // folder in which both teku and JDK are stored
	tekuFolder     = "teku-23.6.2" // folder in which teku is stored (in tekuDepsFolder)
	jdkFolder      = "jdk.20.0.2"  // folder in which JDK is stored (in tekuDepsFolder)
	jdkInstallURL  = "https://download.java.net/java/GA/jdk20.0.2/6e380f22cbe7469fa75fb448bd903d8e/9/GPL/openjdk-20.0.2_linux-x64_bin.tar.gz"
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

	err = installAndUntarFromURL(url)
	if err != nil {
		return
	}

	switch isUpdate {
	case true:
		log.Infof("✅  %s updated!\n\n", t.Name())
	case false:
		log.Infof("✅  %s downloaded!\n\n", t.Name())
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

		err = setupJava(jdkInstallURL)
		if err != nil {
			return
		}
	}

	return
}

func (t *TekuClient) Update() (err error) {
	log.Infof("⬇️  Fetching latest release for %s", t.name)

	latestTag, err := fetchTag(t.githubLocation)
	if err != nil {
		return err
	}

	log.Infof("✅  Fetched latest release: %s", latestTag)

	log.WithField("dependencyTag", latestTag).Infof("⬇️  Updating %s", t.name)

	url := t.ParseUrl(latestTag, "")

	return t.Install(url, true)
}

func (t *TekuClient) FilePath() string {
	return tekuDepsFolder
}

func untarDir(dst string, t *tar.Reader) error {
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dst, header.Name)
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

func setupJava(jdkURL string) (err error) {
	log.Info("⬇️  Downloading JDK...")

	err = installAndUntarFromURL(jdkURL)
	if err != nil {
		return err
	}

	log.Info("✅  JDK downloaded!\n\n")

	luksoNodeDir, err := os.Getwd()
	if err != nil {
		return
	}

	javaHomeVal := fmt.Sprintf("%s/%s/%s", luksoNodeDir, tekuDepsFolder, jdkFolder)

	log.Infof("⚙️  To continue working with Teku please export the JAVA_HOME environmental variable.\n"+
		"The recommended way is to add the following line:\n\n"+
		"export JAVA_HOME=%s\n\n"+
		"To the bash startup file of your choosing (like .bashrc)", javaHomeVal)

	return
}

func installAndUntarFromURL(url string) (err error) {
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
			"❌  Invalid response when downloading on file URL: %s. Response code: %s",
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

	err = untarDir(Teku.FilePath(), tarReader)
	if err != nil {
		return
	}

	return
}
