package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var errNeedRoot = errors.New("You need root privilages to perform this action")

func (dependency *ClientDependency) Download(tagName, destination, commitHash string, checkExist bool) (err error) {
	err = dependency.createDir()
	if err != nil {
		return
	}

	fileUrl := dependency.ParseUrl(tagName, commitHash)

	if fileExists(dependency.filePath) && checkExist {
		log.Warningf("Downloading %s aborted, file %s already exists", fileUrl, dependency.filePath)

		return
	}

	response, err := http.Get(fileUrl)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"invalid response when downloading on file url: %s. Response: %s",
			fileUrl,
			response.Status,
		)
	}

	var responseReader io.Reader = response.Body

	// this means that we are fetching tared geth
	if commitHash != "" {
		g, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}

		defer func() {
			_ = g.Close()
		}()

		t := tar.NewReader(g)
		for {
			header, err := t.Next()

			switch {
			case err == io.EOF:
				break

			case err != nil:
				return err

			default:

			}

			if header.Typeflag == tar.TypeReg && strings.Contains(header.Name, "/geth") {
				responseReader = t

				break
			}
		}
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(responseReader)
	if err != nil {
		return
	}

	err = os.WriteFile(dependency.filePath, buf.Bytes(), os.ModePerm)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		return errNeedRoot
	}

	if err != nil {
		log.Infof("I am in download section: error: %v", err)
		return
	}

	return
}

// createDir creates directory structure for given dependency if it's not a binary file
func (dependency *ClientDependency) createDir() error {
	if strings.Contains(dependency.filePath, binDir) {
		return nil
	}

	segments := strings.Split(dependency.filePath, "/")

	return os.MkdirAll(strings.TrimRight(dependency.filePath, segments[len(segments)-1]), 0750)
}

func downloadBinaries(ctx *cli.Context) (err error) {
	if !ctx.Bool(acceptTermsOfUseFlagName) {
		return errors.New("Terms of Use must be accepted")
	}
	// Get os, then download all binaries into datadir matching desired system
	// After successful download run binary with desired arguments spin and connect them

	err = downloadGeth(ctx)

	if nil != err {
		return
	}

	err = downloadValidator(ctx)

	if nil != err {
		return
	}

	err = downloadPrysm(ctx)

	if nil != err {
		return
	}

	return
}

func downloadConfigs(ctx *cli.Context) error {
	err := downloadGenesis(ctx)
	if nil != err {
		return err
	}

	err = downloadConfig(ctx)
	if nil != err {
		return err
	}

	// after initializing we can run geth init to make sure geth runs fine
	command := exec.Command("geth", "init", clientDependencies[gethGenesisDependencyName].filePath)

	err = command.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Error("No error logs found")

		return err
	}

	// error unrelated to command execution
	if err != nil {
		log.Errorf("There was an error while executing logs command. Error: %v", err)
	}

	return err
}

func downloadGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Geth")

	err = clientDependencies[gethDependencyName].Download(gethTag, binDir, gethCommitHash, true)

	return
}

func downloadGenesis(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Execution Genesis")
	gethDataDir := ctx.String(gethDatadirFlag)
	err = clientDependencies[gethGenesisDependencyName].Download(gethTag, gethDataDir, "", true)

	if nil != err {
		return
	}

	log.WithField("dependencyTag", prysmTag).Info("Downloading Consensus Genesis")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmGenesisDependencyName].Download(prysmTag, prysmDataDir, "", true)

	return
}

func downloadConfig(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm Config")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmConfigDependencyName].Download(prysmTag, prysmDataDir, "", true)

	return
}

func downloadPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm")

	err = clientDependencies[prysmDependencyName].Download(prysmTag, binDir, "", true)

	return
}

func downloadValidator(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", validatorTag).Info("Downloading Validator")
	err = clientDependencies[validatorDependencyName].Download(prysmTag, binDir, "", true)

	return
}
