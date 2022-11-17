package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"strings"
)

func (dependency *ClientDependency) Download(tagName, destination, commitHash string) (err error) {
	dependencyTagPath := dependency.ResolveDirPath(tagName, destination)
	err = os.MkdirAll(dependencyTagPath, 0750)

	if nil != err {
		return
	}

	dependencyLocation := dependency.ResolveBinaryPath(tagName, destination)

	fileUrl := dependency.ParseUrl(tagName, commitHash)

	if fileExists(dependencyLocation) {
		log.Warningf("Downloading %s aborted, file %s already exists", fileUrl, dependencyLocation)

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

	output, err := os.Create(dependencyLocation)

	if nil != err {
		return
	}

	defer func() {
		_ = output.Close()
	}()

	_, err = io.Copy(output, responseReader)

	if nil != err {
		return
	}

	err = os.Chmod(dependencyLocation, os.ModePerm)

	return
}

func downloadBinaries(ctx *cli.Context) (err error) {
	if !ctx.Bool(acceptTermsOfUseFlagName) {
		return errors.New("Terms of Use must be accepted")
	}
	// Get os, then download all binaries into datadir matching desired system
	// After successful download run binary with desired arguments spin and connect them
	// Orchestrator can be run from-memory
	//err = downloadGenesis(ctx)
	//
	//if nil != err {
	//	return
	//}

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

	//err = downloadConfig(ctx)

	return
}

func downloadGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Geth")

	err = clientDependencies[gethDependencyName].Download(gethTag, binDir, gethCommitHash)

	return
}

func downloadGenesis(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Geth Genesis")
	gethDataDir := ctx.String(gethDatadirFlag)
	err = clientDependencies[gethGenesisDependencyName].Download(gethTag, gethDataDir, "")

	if nil != err {
		return
	}

	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm Genesis")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmGenesisDependencyName].Download(prysmTag, prysmDataDir, "")

	return
}

func downloadConfig(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm Config")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmConfigDependencyName].Download(prysmTag, prysmDataDir, "")

	return
}

func downloadPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm")

	err = clientDependencies[prysmDependencyName].Download(prysmTag, binDir, "")

	return
}

func downloadValidator(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", validatorTag).Info("Downloading Validator")
	err = clientDependencies[validatorDependencyName].Download(prysmTag, binDir, "")

	return
}
