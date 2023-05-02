package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	configPerms = 0750
	binaryPerms = int(os.ModePerm)
)

func (dependency *ClientDependency) Download(tag, commitHash string, isUpdate bool, permissions int) (err error) {
	log.Infof("⬇️  Downloading %s...", dependency.name)

	err = dependency.createDir()
	if err != nil {
		return
	}

	fileUrl := dependency.ParseUrl(tag, commitHash)

	if fileExists(dependency.filePath) {
		switch dependency.isBinary {
		case true:
			if isUpdate {
				break
			}

			message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", dependency.name)
			input := registerInputWithMessage(message)
			if !strings.EqualFold(input, "y") && input != "" {
				log.Info("⏭️  Skipping installation...")

				return nil
			}

		case false:
			log.Infof("❌  Downloading %s file aborted: already exists", dependency.filePath)

			return
		}
	}

	response, err := http.Get(fileUrl)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", fileUrl)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			fileUrl,
			response.Status,
		)
	}

	var responseReader io.Reader = response.Body

	// this means that we are fetching tared client
	switch dependency.name {
	case gethDependencyName, erigonDependencyName, lighthouseDependencyName:
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

			var targetHeader string
			switch dependency.name {
			case gethDependencyName:
				targetHeader = "/" + gethDependencyName
			case erigonDependencyName, lighthouseDependencyName:
				targetHeader = dependency.name
			}

			if header.Typeflag == tar.TypeReg && strings.Contains(header.Name, targetHeader) {
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

	if err != nil && strings.Contains(err.Error(), "Permission denied") {
		return errNeedRoot
	}

	if err != nil {
		log.Infof("❌  Couldn't save file: %v", err)

		return
	}

	log.Infof("✅ Downloaded %s!\n\n", dependency.name)

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
	if isAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return cli.Exit(folderNotInitialized, 1)
	}

	isRoot, err := isRoot()
	if err != nil {
		return cli.Exit(fmt.Sprintf("There was an error while checking user privileges: %v", err), 1)
	}
	if !isRoot {
		return cli.Exit(errNeedRoot, 1)
	}

	var (
		consensusInput    string
		executionInput    string
		selectedConsensus string
		selectedExecution string
		consensusTag      string
		executionTag      string
		commitHash        string
	)

	consensusMessage := "\nWhich consensus client do you want to install?\n" +
		"1: prysm\n2: lighthouse\n> "
	executionMessage := "\nWhich execution client do you want to install?\n" +
		"1: geth\n2: erigon\n> "

	consensusInput = registerInputWithMessage(consensusMessage)
	for consensusInput != "1" && consensusInput != "2" {
		consensusInput = registerInputWithMessage("Please provide a valid option\n> ")
	}

	switch consensusInput {
	case "1":
		selectedConsensus = prysmDependencyName
		consensusTag = ctx.String(prysmTagFlag)
	case "2":
		selectedConsensus = lighthouseDependencyName
		consensusTag = ctx.String(lighthouseTagFlag)
	}

	executionInput = registerInputWithMessage(executionMessage)
	for executionInput != "1" && executionInput != "2" {
		executionInput = registerInputWithMessage("Please provide a valid option\n> ")
	}

	switch executionInput {
	case "1":
		selectedExecution = gethDependencyName
		executionTag = ctx.String(gethTagFlag)
		commitHash = ctx.String(gethCommitHashFlag)
	case "2":
		selectedExecution = erigonDependencyName
		executionTag = ctx.String(erigonTagFlag)
	}

	termsAgreed := ctx.Bool(agreeTermsFlag)
	if !termsAgreed {
		fmt.Println("")
		accepted := acceptTermsInteractive()
		if !accepted {
			return cli.Exit("❌  Terms of use not accepted - aborting...", 1)
		}
	} else {
		log.Info("✅  You accepted Prysm's Terms of Use: https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md")
	}

	err = clientDependencies[selectedExecution].Download(executionTag, commitHash, false, binaryPerms)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedExecution, err), 1)
	}

	err = clientDependencies[selectedConsensus].Download(consensusTag, "", false, binaryPerms)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedConsensus, err), 1)
	}

	// for now, we also need to download a validator client
	// when other validator clients will be implemented we will download the one bound to consensus clients
	err = clientDependencies[validatorDependencyName].Download(prysmTag, "", false, binaryPerms)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while downloading validator: %v", err), 1)
	}

	err = cfg.Create(selectedExecution, selectedConsensus)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while creating configration file: %v", err), 1)
	}

	log.Info("✅  Configuration files created!")
	log.Info("✅  Clients have been successfully installed.")
	log.Info("➡️   Start your node using 'lukso start'")

	return
}

func acceptTermsInteractive() bool {
	message := "You are about to download the clients that are necessary to run the LUKSO Blockchain.\n" +
		"To install Prysm you are required to accept the terms of use:\n" +
		"https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md\n\n" +
		"Do you agree? [Y/n]: "

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Error("You need to type Y to continue.")
		return false
	}

	return true
}
