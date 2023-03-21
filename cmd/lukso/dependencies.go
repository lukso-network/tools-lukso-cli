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
	gethMainnetConfigName             = "geth_mainnet_config"
	prysmMainnetGenesisDependencyName = "prysm_mainnet_genesis_state"
	prysmMainnetConfigDependencyName  = "prysm_mainnet_config"

	gethTestnetGenesisDependencyName  = "geth_testnet_genesis"
	gethTestnetConfigName             = "geth_testnet_config"
	prysmTestnetGenesisDependencyName = "prysm_testnet_genesis_state"
	prysmTestnetConfigDependencyName  = "prysm_testnet_config"

	gethDevnetGenesisDependencyName  = "geth_devnet_genesis"
	gethDevnetConfigName             = "geth_devnet_config"
	prysmDevnetGenesisDependencyName = "prysm_devnet_genesis_state"
	prysmDevnetConfigDependencyName  = "prysm_devnet_config"
)

var (
	gethSelectedGenesis  = gethMainnetGenesisDependencyName
	gethSelectedConfig   = gethMainnetConfigName
	prysmSelectedGenesis = prysmMainnetGenesisDependencyName
	prysmSelectedConfig  = prysmMainnetConfigDependencyName

	clientDependencies = map[string]*ClientDependency{
		// ----- BINARIES -----
		gethDependencyName: {
			baseUrl:  "https://gethstore.blob.core.windows.net/builds/geth-%s-amd64-%s-%s.tar.gz",
			name:     gethDependencyName,
			filePath: "", // binary dir selected during runtime
		},
		prysmDependencyName: {
			baseUrl:  "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-%s-amd64",
			name:     prysmDependencyName,
			filePath: "", // binary dir selected during runtime
		},
		validatorDependencyName: {
			baseUrl:  "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-%s-amd64",
			name:     validatorDependencyName,
			filePath: "", // binary dir selected during runtime
		},

		// ----- CONFIGS -----
		// ----- MAINNET -----
		gethMainnetGenesisDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     gethMainnetGenesisDependencyName,
			filePath: "./config/mainnet/geth/genesis.json",
		},
		gethMainnetConfigName: {
			baseUrl:  "NOT SUPPORTED",
			name:     gethMainnetConfigName,
			filePath: "./config/mainnet/geth/geth.toml",
		},
		prysmMainnetGenesisDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     prysmMainnetGenesisDependencyName,
			filePath: "./config/mainnet/shared/genesis.ssz",
		},
		prysmMainnetConfigDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     prysmMainnetConfigDependencyName,
			filePath: "./config/mainnet/shared/config.yml",
		},
		// ----- TESTNET -----
		gethTestnetGenesisDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     gethTestnetGenesisDependencyName,
			filePath: "./config/testnet/geth/genesis.json",
		},
		gethTestnetConfigName: {
			baseUrl:  "NOT SUPPORTED",
			name:     gethTestnetConfigName,
			filePath: "./config/testnet/geth/geth.toml",
		},
		prysmTestnetGenesisDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     prysmTestnetGenesisDependencyName,
			filePath: "./config/testnet/shared/genesis.ssz",
		},
		prysmTestnetConfigDependencyName: {
			baseUrl:  "NOT SUPPORTED",
			name:     prysmTestnetConfigDependencyName,
			filePath: "./config/testnet/shared/config.yml",
		},
		// ----- DEVNET -----
		gethDevnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/2022/geth/genesis.json",
			name:     gethDevnetGenesisDependencyName,
			filePath: "./config/devnet/geth/genesis.json",
		},
		gethDevnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/2022/geth/geth.toml",
			name:     gethDevnetConfigName,
			filePath: "./config/devnet/geth/geth.toml",
		},
		prysmDevnetGenesisDependencyName: {
			baseUrl:  "https://github.com/lukso-network/network-configs/raw/main/devnets/2022/prysm/genesis.ssz",
			name:     prysmDevnetGenesisDependencyName,
			filePath: "./config/devnet/shared/genesis.ssz",
		},
		prysmDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/2022/prysm/config.yml",
			name:     prysmDevnetConfigDependencyName,
			filePath: "./config/devnet/shared/config.yml",
		},
	}
)

type ClientDependency struct {
	baseUrl  string
	name     string
	filePath string
}

func (dependency *ClientDependency) ParseUrl(tag, commitHash string) (url string) {
	// do not parse when no occurrences
	sprintOccurrences := strings.Count(dependency.baseUrl, "%s")

	baseUrl := dependency.baseUrl

	switch sprintOccurrences {
	case 3:
		if commitHash != "" {
			return fmt.Sprintf(baseUrl, systemOs, tag, commitHash)
		}
		return fmt.Sprintf(baseUrl, tag, tag, systemOs)
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
