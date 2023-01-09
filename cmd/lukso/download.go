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
	"strings"
)

var errNeedRoot = errors.New("You need root privilages to perform this action")

func (dependency *ClientDependency) Download(tagName, commitHash string, overrideFile bool) (err error) {
	err = dependency.createDir()
	if err != nil {
		return
	}

	fileUrl := dependency.ParseUrl(tagName, commitHash)

	if fileExists(dependency.filePath) && !overrideFile {
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

	err := os.MkdirAll(strings.TrimRight(dependency.filePath, segments[len(segments)-1]), 0750)
	if errors.Is(err, os.ErrExist) {
		log.Errorf("%s already exists!", dependency.name)
	}

	return err
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

	return downloadPrysmConfig(ctx)
}

func downloadGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Geth")

	err = clientDependencies[gethDependencyName].Download(gethTag, gethCommitHash, false)

	return
}

func downloadGenesis(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Execution Genesis")

	err = clientDependencies[gethGenesisDependencyName].Download(gethTag, "", false)

	if nil != err {
		return
	}

	log.WithField("dependencyTag", prysmTag).Info("Downloading Consensus Genesis")

	err = clientDependencies[prysmGenesisDependencyName].Download(prysmTag, "", false)

	return
}

func downloadPrysmConfig(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm Config")

	err = clientDependencies[prysmConfigDependencyName].Download(prysmTag, "", false)

	return
}

func downloadPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm")

	err = clientDependencies[prysmDependencyName].Download(prysmTag, "", false)

	return
}

func downloadValidator(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", validatorTag).Info("Downloading Validator")
	err = clientDependencies[validatorDependencyName].Download(prysmTag, "", false)

	return
}
