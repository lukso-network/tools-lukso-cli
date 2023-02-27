package main

import (
	"fmt"
	"runtime"
	"strings"
)

// TODO: consider to move it to common/shared
const (
	gethDependencyName      = "geth"
	prysmDependencyName     = "prysm"
	validatorDependencyName = "validator"

	gethMainnetGenesisDependencyName  = "geth_mainnet_genesis"
	prysmMainnetGenesisDependencyName = "prysm_mainnet_genesis_state"
	prysmMainnetConfigDependencyName  = "prysm_mainnet_config"

	gethTestnetGenesisDependencyName  = "geth_testnet_genesis"
	prysmTestnetGenesisDependencyName = "prysm_testnet_genesis_state"
	prysmTestnetConfigDependencyName  = "prysm_testnet_config"

	gethDevnetGenesisDependencyName  = "geth_devnet_genesis"
	prysmDevnetGenesisDependencyName = "prysm_devnet_genesis_state"
	prysmDevnetConfigDependencyName  = "prysm_devnet_config"
)

var (
	gethSelectedGenesis  = gethMainnetGenesisDependencyName
	prysmSelectedGenesis = prysmMainnetGenesisDependencyName
	prysmSelectedConfig  = prysmMainnetConfigDependencyName

	clientDependencies = map[string]*ClientDependency{
		// ----- BINARIES -----
		gethDependencyName: {
			baseUnixUrl:   "https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-%s-%s.tar.gz",
			baseDarwinUrl: "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-%s-%s.tar.gz",
			name:          gethDependencyName,
			filePath:      "", // binary dir selected during runtime
		},
		prysmDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-darwin-amd64",
			name:          prysmDependencyName,
			filePath:      "", // binary dir selected during runtime
		},
		validatorDependencyName: {
			baseUnixUrl:   "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-linux-amd64",
			baseDarwinUrl: "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-darwin-amd64",
			name:          validatorDependencyName,
			filePath:      "", // binary dir selected during runtime
		},

		// ----- CONFIGS -----
		// ----- MAINNET -----
		gethMainnetGenesisDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          gethMainnetGenesisDependencyName,
			filePath:      "./config/mainnet/geth/genesis.json",
		},
		prysmMainnetGenesisDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          prysmMainnetGenesisDependencyName,
			filePath:      "./config/mainnet/shared/genesis.ssz",
		},
		prysmMainnetConfigDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          prysmMainnetConfigDependencyName,
			filePath:      "./config/mainnet/shared/config.yml",
		},
		// ----- TESTNET -----
		gethTestnetGenesisDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          gethTestnetGenesisDependencyName,
			filePath:      "./config/testnet/geth/genesis.json",
		},
		prysmTestnetGenesisDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          prysmTestnetGenesisDependencyName,
			filePath:      "./config/testnet/shared/genesis.ssz",
		},
		prysmTestnetConfigDependencyName: {
			baseUnixUrl:   "NOT SUPPORTED",
			baseDarwinUrl: "NOT SUPPORTED",
			name:          prysmTestnetConfigDependencyName,
			filePath:      "./config/testnet/shared/config.yml",
		},
		// ----- DEVNET -----
		gethDevnetGenesisDependencyName: {
			baseUnixUrl:   "https://raw.githubusercontent.com/lukso-network/network-configs/devnet/dev/2022/geth/genesis.json",
			baseDarwinUrl: "https://raw.githubusercontent.com/lukso-network/network-configs/devnet/dev/2022/geth/genesis.json",
			name:          gethDevnetGenesisDependencyName,
			filePath:      "./config/devnet/geth/genesis.json",
		},
		prysmDevnetGenesisDependencyName: {
			baseUnixUrl:   "https://github.com/lukso-network/network-configs/raw/devnet/dev/2022/prysm/genesis.ssz",
			baseDarwinUrl: "https://github.com/lukso-network/network-configs/raw/devnet/dev/2022/prysm/genesis.ssz",
			name:          prysmDevnetGenesisDependencyName,
			filePath:      "./config/devnet/shared/genesis.ssz",
		},
		prysmDevnetConfigDependencyName: {
			baseUnixUrl:   "https://raw.githubusercontent.com/lukso-network/network-configs/devnet/dev/2022/prysm/config.yml",
			baseDarwinUrl: "https://raw.githubusercontent.com/lukso-network/network-configs/devnet/dev/2022/prysm/config.yml",
			name:          prysmDevnetConfigDependencyName,
			filePath:      "./config/devnet/shared/config.yml",
		},
	}
)

type ClientDependency struct {
	baseUnixUrl   string
	baseDarwinUrl string
	name          string
	filePath      string
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

	// setting PATH for binaries
	clientDependencies[gethDependencyName].filePath = binDir + "/geth"
	clientDependencies[prysmDependencyName].filePath = binDir + "/prysm"
	clientDependencies[validatorDependencyName].filePath = binDir + "/validator"
}
