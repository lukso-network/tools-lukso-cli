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
			isBinary: true,
		},
		erigonDependencyName: {
			baseUrl:  "https://github.com/ledgerwatch/erigon/releases/download/v%s/erigon_%s_%s_amd64.tar.gz",
			name:     erigonDependencyName,
			filePath: "",
			isBinary: true,
		},
		prysmDependencyName: {
			baseUrl:  "https://github.com/prysmaticlabs/prysm/releases/download/%s/beacon-chain-%s-%s-amd64",
			name:     prysmDependencyName,
			filePath: "", // binary dir selected during runtime
			isBinary: true,
		},
		lighthouseDependencyName: {
			baseUrl:  "https://github.com/sigp/lighthouse/releases/download/%s/lighthouse-%s-x86_64-%s-%s.tar.gz",
			name:     lighthouseDependencyName,
			filePath: "", // binary dir selected during runtime
			isBinary: true,
		},
		validatorDependencyName: {
			baseUrl:  "https://github.com/prysmaticlabs/prysm/releases/download/%s/validator-%s-%s-amd64",
			name:     validatorDependencyName,
			filePath: "", // binary dir selected during runtime
			isBinary: true,
		},

		// ----- CONFIGS -----
		// ----- MAINNET -----
		gethMainnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			name:     gethMainnetGenesisDependencyName,
			filePath: mainnetConfig + "/" + genesisJsonPath,
		},
		gethMainnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
			name:     gethMainnetConfigName,
			filePath: mainnetConfig + "/" + gethTomlPath,
		},
		prysmMainnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz", // no genesis state file for mainnet yet
			name:     prysmMainnetGenesisDependencyName,
			filePath: mainnetConfig + "/" + genesisStateFilePath,
		},
		prysmMainnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			name:     prysmMainnetChainConfigDependencyName,
			filePath: mainnetConfig + "/" + chainConfigYamlPath,
		},
		prysmMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/mainnet/prysm/prysm.yaml",
			name:     prysmMainnetConfigDependencyName,
			filePath: mainnetConfig + "/" + prysmYamlPath,
		},
		validatorMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/mainnet/prysm/validator.yaml",
			name:     validatorMainnetConfigDependencyName,
			filePath: mainnetConfig + "/" + validatorYamlPath,
		},
		// ----- TESTNET -----
		gethTestnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.json",
			name:     gethTestnetGenesisDependencyName,
			filePath: testnetConfig + "/" + genesisJsonPath,
		},
		gethTestnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
			name:     gethTestnetConfigName,
			filePath: testnetConfig + "/" + gethTomlPath,
		},
		prysmTestnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.ssz", // no genesis state file for testnet yet
			name:     prysmTestnetGenesisDependencyName,
			filePath: testnetConfig + "/" + genesisStateFilePath,
		},
		prysmTestnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/config.yaml",
			name:     prysmTestnetChainConfigDependencyName,
			filePath: testnetConfig + "/" + chainConfigYamlPath,
		},
		prysmTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/testnet/prysm/prysm.yaml",
			name:     prysmTestnetConfigDependencyName,
			filePath: testnetConfig + "/" + prysmYamlPath,
		},
		validatorTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/feature/blockchain-clients-configs/testnet/prysm/validator.yaml",
			name:     validatorTestnetConfigDependencyName,
			filePath: testnetConfig + "/" + validatorYamlPath,
		},
		// ----- DEVNET -----
		gethDevnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/shared/genesis.json",
			name:     gethDevnetGenesisDependencyName,
			filePath: devnetConfig + "/" + genesisJsonPath,
		},
		gethDevnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/geth/geth.toml",
			name:     gethDevnetConfigName,
			filePath: devnetConfig + "/" + gethTomlPath,
		},
		prysmDevnetGenesisDependencyName: {
			baseUrl:  "https://github.com/lukso-network/network-configs/raw/main/devnets/3030/shared/genesis.ssz",
			name:     prysmDevnetGenesisDependencyName,
			filePath: devnetConfig + "/" + genesisStateFilePath,
		},
		prysmDevnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/shared/config.yaml",
			name:     prysmDevnetChainConfigDependencyName,
			filePath: devnetConfig + "/" + chainConfigYamlPath,
		},
		prysmDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/prysm/prysm.yaml",
			name:     prysmDevnetConfigDependencyName,
			filePath: devnetConfig + "/" + prysmYamlPath,
		},
		validatorDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/prysm/validator.yaml",
			name:     validatorDevnetConfigDependencyName,
			filePath: devnetConfig + "/" + validatorYamlPath,
		},
	}
)

type ClientDependency struct {
	baseUrl  string
	name     string
	filePath string
	isBinary bool
}

func (dependency *ClientDependency) ParseUrl(tag, commitHash string) (url string) {
	// do not parse when no occurrences
	var (
		baseUrl           = dependency.baseUrl
		sprintOccurrences = strings.Count(dependency.baseUrl, "%s")
		systemName        string
		urlSystem         = systemOs
	)

	// for lighthouse
	switch systemOs {
	case ubuntu:
		systemName = "unknown"
		urlSystem += "-gnu"
	case macos:
		systemName = "apple"
	default:
		systemName = "unknown"
		urlSystem += "-gnu"
	}

	switch sprintOccurrences {
	case 3:
		if commitHash != "" {
			return fmt.Sprintf(baseUrl, systemOs, tag, commitHash)
		}
		return fmt.Sprintf(baseUrl, tag, tag, systemOs)
	case 4:
		return fmt.Sprintf(baseUrl, tag, tag, systemName, urlSystem)
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
	clientDependencies[erigonDependencyName].filePath = binDir + "/erigon"
	clientDependencies[prysmDependencyName].filePath = binDir + "/prysm"
	clientDependencies[lighthouseDependencyName].filePath = binDir + "/lighthouse"
	clientDependencies[validatorDependencyName].filePath = binDir + "/validator"
}
