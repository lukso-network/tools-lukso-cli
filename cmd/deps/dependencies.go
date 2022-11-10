package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// TODO: consider to move it to common/shared
const (
	gethDependencyName         = "geth"
	gethGenesisDependencyName  = "geth_private_testnet_genesis.json"
	prysmDependencyName        = "prysm"
	validatorDependencyName    = "validator"
	prysmGenesisDependencyName = "prysm_private_testnet_genesis.ssz"
	prysmConfigDependencyName  = "config.yml"
)

var (
	clientDependencies = map[string]*ClientDependency{
		gethDependencyName: {
			baseUnixUrl:   "https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.26-e5eb32ac.tar.gz",
			baseDarwinUrl: "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.10.26-e5eb32ac.tar.gz",
			name:          gethDependencyName,
		},
		gethGenesisDependencyName: {
			baseUnixUrl:   "https://storage.googleapis.com/l16-common/pandora/pandora_private_testnet_genesis.json",
			baseDarwinUrl: "https://storage.googleapis.com/l16-common/pandora/pandora_private_testnet_genesis.json",
			name:          gethGenesisDependencyName,
		},
		prysmDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/v3.1.2/beacon-chain-v3.1.2-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/v3.1.2/beacon-chain-v3.1.2-darwin-amd64",
			name:          prysmDependencyName,
		},
		validatorDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/v3.1.2/validator-v3.1.2-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/v3.1.2/validator-v3.1.2-darwin-amd64",
			name:          validatorDependencyName,
		},
		prysmGenesisDependencyName: {
			baseUnixUrl:   "https://storage.googleapis.com/l16-common/vanguard/vanguard_private_testnet_genesis.ssz",
			baseDarwinUrl: "https://storage.googleapis.com/l16-common/vanguard/vanguard_private_testnet_genesis.ssz",
			name:          prysmGenesisDependencyName,
		},
		prysmConfigDependencyName: {
			baseUnixUrl:   "https://storage.googleapis.com/l16-common/vanguard/chain-config.yaml",
			baseDarwinUrl: "https://storage.googleapis.com/l16-common/vanguard/chain-config.yaml",
			name:          prysmConfigDependencyName,
		},
	}
)

type ClientDependency struct {
	baseUnixUrl   string
	baseDarwinUrl string
	name          string
}

func (dependency *ClientDependency) ParseUrl(tagName, commitHash string) (url string) {
	// do not parse when no occurrences
	sprintOccurrences := strings.Count(dependency.baseUnixUrl, "%s")
	currentOs := systemOs
	baseUrl := ""

	switch currentOs {
	case ubuntu:
		baseUrl = dependency.baseUnixUrl
	case macos:
		baseUrl = dependency.baseDarwinUrl
	default:
		baseUrl = dependency.baseUnixUrl // defaulting to unix, as in previous implementation
	}

	switch sprintOccurrences {
	case 1:
		return fmt.Sprintf(baseUrl, tagName)
	case 2:
		if commitHash != "" {
			return fmt.Sprintf(baseUrl, tagName, commitHash)
		}
		return fmt.Sprintf(baseUrl, tagName, tagName)
	default:
		return baseUrl
	}
}

func (dependency *ClientDependency) ResolveDirPath(tagName string, datadir string) (location string) {
	location = fmt.Sprintf("%s/%s", datadir, tagName)

	return
}

func (dependency *ClientDependency) ResolveBinaryPath(tagName string, datadir string) (location string) {
	location = fmt.Sprintf("%s/%s", dependency.ResolveDirPath(tagName, datadir), dependency.name)

	return
}

func (dependency *ClientDependency) Run(
	tagName string,
	destination string,
	arguments []string,
	attachStdInAndErr bool,
) (err error) {
	binaryPath := dependency.ResolveBinaryPath(tagName, destination)
	command := exec.Command(binaryPath, arguments...)

	if attachStdInAndErr {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}

	err = command.Start()

	return
}

func (dependency *ClientDependency) Download(tagName, destination, commitHash string) (err error) {
	dependencyTagPath := dependency.ResolveDirPath(tagName, destination)
	err = os.MkdirAll(dependencyTagPath, 0750)

	if nil != err {
		return
	}

	dependencyLocation := dependency.ResolveBinaryPath(tagName, destination)

	fileUrl := dependency.ParseUrl(tagName, commitHash)

	if fileExists(dependencyLocation) {
		log.Warningf("I am not downloading %s, file %s already exists", fileUrl, dependencyLocation)

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

func setupOperatingSystem() {
	systemOs = runtime.GOOS
}
