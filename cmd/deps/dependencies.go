package main

import (
	"fmt"
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
			baseUnixUrl:   "https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-%s-%s.tar.gz",
			baseDarwinUrl: "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-%s-%s.tar.gz",
			name:          gethDependencyName,
		},
		gethGenesisDependencyName: {
			baseUnixUrl:   "https://storage.googleapis.com/l16-common/pandora/pandora_private_testnet_genesis.json",
			baseDarwinUrl: "https://storage.googleapis.com/l16-common/pandora/pandora_private_testnet_genesis.json",
			name:          gethGenesisDependencyName,
		},
		prysmDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-darwin-amd64",
			name:          prysmDependencyName,
		},
		validatorDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-darwin-amd64",
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
	if datadir == binDir {
		location = binDir

		return
	}

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

func setupOperatingSystem() {
	systemOs = runtime.GOOS

	switch systemOs {
	case ubuntu, macos:
		binDir = unixBinDir
	case windows:
		binDir = windowsBinDir
	default:
		log.Panicf("unexpected OS: %v", systemOs)
	}
}
