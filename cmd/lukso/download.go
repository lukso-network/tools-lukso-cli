package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/m8b-dev/lukso-cli/config"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	configPerms = 0750
	binaryPerms = int(os.ModePerm)
)

func (dependency *ClientDependency) Download(tag, commitHash string, isUpdate bool, permissions int) (err error) {
	log.Infof("Downloading %s", dependency.name)

	err = dependency.createDir()
	if err != nil {
		return
	}

	fileUrl := dependency.ParseUrl(tag, commitHash)

	switch dependency.isBinary {
	case true:
		if isUpdate {
			break
		}

	case false:
		log.Infof("Downloading %s file aborted: already exists", dependency.filePath)

		return
	}

	response, err := http.Get(fileUrl)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("File under URL %s not found - aborting...", fileUrl)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"invalid response when downloading on file url: %s. Response code: %s",
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

	err = os.WriteFile(dependency.filePath, buf.Bytes(), os.FileMode(permissions))

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		return errNeedRoot
	}

	if err != nil {
		log.Infof("I am in download section: error: %v", err)

		return
	}

	log.Infof("Downloaded %s!", dependency.name)

	return
}

// createDir creates directory structure for given dependency if it's not a binary file
func (dependency *ClientDependency) createDir() error {
	if strings.Contains(dependency.filePath, binDir) {
		return nil
	}

	err := os.MkdirAll(truncateFileFromDir(dependency.filePath), configPerms)
	if errors.Is(err, os.ErrExist) {
		log.Errorf("%s already exists!", dependency.name)
	}

	if errors.Is(err, os.ErrPermission) {
		return errNeedRoot
	}

	return err
}

func installBinaries(ctx *cli.Context) (err error) {
	isRoot, err := isRoot()
	if err != nil {
		return err
	}
	if !isRoot {
		return errNeedRoot
	}

	var (
		consensusInput    string
		executionInput    string
		selectedConsensus string
		selectedExecution string
	)

	consensusMessage := "Which consensus client you want to install?\n" +
		"1: prysm\n2: lighthouse\n> "
	executionMessage := "Which execution client you want to install?\n" +
		"1: geth\n2: erigon\n> "

	consensusInput = registerInputWithMessage(consensusMessage)
	for consensusInput != "1" && consensusInput != "2" {
		consensusInput = registerInputWithMessage("Please provide a valid option\n> ")
	}

	switch consensusInput {
	case "1":
		selectedConsensus = prysmDependencyName
	case "2":
		selectedConsensus = lighthouseDependencyName
	}

	if selectedConsensus == lighthouseDependencyName {
		log.Error("Please select different consensus client")

		return nil
	}

	executionInput = registerInputWithMessage(executionMessage)
	for executionInput != "1" && executionInput != "2" {
		executionInput = registerInputWithMessage("Please provide a valid option\n> ")
	}

	switch executionInput {
	case "1":
		selectedExecution = gethDependencyName
	case "2":
		selectedExecution = erigonDependencyName
	}

	if selectedExecution == erigonDependencyName {
		log.Error("Please select different execution client")

		return nil
	}

	termsAgreed := ctx.Bool(agreeTermsFlag)
	if !termsAgreed {
		accepted := acceptTermsInteractive()
		if !accepted {
			log.Info("Terms of use not accepted - aborting...")

			return nil
		}
	} else {
		log.Info("You accepted terms of use of accepted clients - read more here: https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md")
	}

	err = clientDependencies[selectedExecution].Download(ctx.String(gethTagFlag), ctx.String(gethCommitHashFlag), false, binaryPerms)
	if err != nil {
		return err
	}

	err = clientDependencies[selectedConsensus].Download(ctx.String(prysmTagFlag), "", false, binaryPerms)
	if err != nil {
		return err
	}

	cfg := config.NewConfig(config.Path)
	err = cfg.Create(selectedExecution, selectedConsensus)
	if err != nil {
		return err
	}

	log.Info("Config created!")

	return
}

func downloadGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("Downloading Geth")

	err = clientDependencies[gethDependencyName].Download(gethTag, gethCommitHash, false, binaryPerms)

	return
}

func downloadPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("Downloading Prysm")

	err = clientDependencies[prysmDependencyName].Download(prysmTag, "", false, binaryPerms)

	return
}

func downloadValidator(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", validatorTag).Info("Downloading Validator")
	err = clientDependencies[validatorDependencyName].Download(prysmTag, "", false, binaryPerms)

	return
}

func acceptTermsInteractive() bool {
	message := "You are about to download clients necessary to run LUKSO CLI. " +
		"By proceeding further you accept Terms of Use of provided clients, you can read more here: " +
		"https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md\n" +
		"Do you wish to continue? [Y/n]: "

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") {
		log.Error("You need to type Y to continue.")
		return false
	}

	return true
}
