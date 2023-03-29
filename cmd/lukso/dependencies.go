package main

import (
	"fmt"
	"runtime"
	"strings"
)

// TODO: consider to move it to common/shared
const (
	gethDependencyName       = "geth"
	prysmDependencyName      = "prysm"
	validatorDependencyName  = "validator"
	lighthouseDependencyName = "lighthouse"
	erigonDependencyName     = "erigon"

	gethMainnetGenesisDependencyName      = "geth_mainnet_genesis"
	gethMainnetConfigName                 = "geth_mainnet_config"
	prysmMainnetGenesisDependencyName     = "prysm_mainnet_genesis_state"
	prysmMainnetChainConfigDependencyName = "prysm_mainnet_chain_config"
	prysmMainnetConfigDependencyName      = "prysm_mainnet_config"
	validatorMainnetConfigDependencyName  = "validator_mainnet_config"

	gethTestnetGenesisDependencyName      = "geth_testnet_genesis"
	gethTestnetConfigName                 = "geth_testnet_config"
	prysmTestnetGenesisDependencyName     = "prysm_testnet_genesis_state"
	prysmTestnetChainConfigDependencyName = "prysm_testnet_chain_config"
	prysmTestnetConfigDependencyName      = "prysm_testnet_config"
	validatorTestnetConfigDependencyName  = "validator_testnet_config"

	gethDevnetGenesisDependencyName      = "geth_devnet_genesis"
	gethDevnetConfigName                 = "geth_devnet_config"
	prysmDevnetGenesisDependencyName     = "prysm_devnet_genesis_state"
	prysmDevnetChainConfigDependencyName = "prysm_devnet_chain_config"
	prysmDevnetConfigDependencyName      = "prysm_devnet_config"
	validatorDevnetConfigDependencyName  = "validator_devnet_config"
)

var (
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
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			name:     gethMainnetGenesisDependencyName,
			filePath: "./config/mainnet/shared/genesis.json",
		},
		gethMainnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
			name:     gethMainnetConfigName,
			filePath: "./config/mainnet/geth/geth.toml",
		},
		prysmMainnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz", // no genesis state file for mainnet yet
			name:     prysmMainnetGenesisDependencyName,
			filePath: "./config/mainnet/shared/genesis.ssz",
		},
		prysmMainnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			name:     prysmMainnetChainConfigDependencyName,
			filePath: "./config/mainnet/shared/config.yaml",
		},
		prysmMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/mainnet/prysm/prysm.yaml",
			name:     prysmMainnetConfigDependencyName,
			filePath: "./config/mainnet/prysm/prysm.yaml",
		},
		validatorMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/mainnet/prysm/validator.yaml",
			name:     validatorMainnetConfigDependencyName,
			filePath: "./config/mainnet/prysm/validator.yaml",
		},
		// ----- TESTNET -----
		gethTestnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.json",
			name:     gethTestnetGenesisDependencyName,
			filePath: "./config/testnet/shared/genesis.json",
		},
		gethTestnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
			name:     gethTestnetConfigName,
			filePath: "./config/testnet/geth/geth.toml",
		},
		prysmTestnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.ssz", // no genesis state file for testnet yet
			name:     prysmTestnetGenesisDependencyName,
			filePath: "./config/testnet/shared/genesis.ssz",
		},
		prysmTestnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/config.yaml",
			name:     prysmTestnetChainConfigDependencyName,
			filePath: "./config/testnet/shared/config.yaml",
		},
		prysmTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/testnet/prysm/prysm.yaml",
			name:     prysmTestnetConfigDependencyName,
			filePath: "./config/testnet/prysm/prysm.yaml",
		},
		validatorTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/testnet/prysm/validator.yaml",
			name:     validatorTestnetConfigDependencyName,
			filePath: "./config/testnet/prysm/validator.yaml",
		},
		// ----- DEVNET -----
		gethDevnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/2022/geth/genesis.json",
			name:     gethDevnetGenesisDependencyName,
			filePath: "./config/devnet/shared/genesis.json",
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
		prysmDevnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/2022/prysm/config.yml",
			name:     prysmDevnetChainConfigDependencyName,
			filePath: "./config/devnet/shared/config.yaml",
		},
		prysmDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/devnets/3030/prysm/prysm.yaml",
			name:     prysmDevnetConfigDependencyName,
			filePath: "./config/devnet/prysm/prysm.yaml",
		},
		validatorDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/devnets/3030/prysm/validator.yaml",
			name:     validatorDevnetConfigDependencyName,
			filePath: "./config/devnet/prysm/validator.yaml",
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
