package main

import (
	"fmt"
	"runtime"
	"strings"
)

// TODO: consider to move it to common/shared
// TODO: Disconnect names of dependencies from actual terminal outputs
const (
	// client names
	gethDependencyName       = "geth"
	prysmDependencyName      = "prysm"
	validatorDependencyName  = "validator"
	lighthouseDependencyName = "lighthouse"
	erigonDependencyName     = "erigon"

	// genesis files
	mainnetGenesisDependencyName    = "mainnet genesis"
	testnetGenesisDependencyName    = "testnet genesis"
	devnetGenesisDependencyName     = "devnet genesis"
	mainnetGenesis35DependencyName  = "mainnet genesis (35M)"
	mainnetGenesis42DependencyName  = "mainnet genesis (42M)"
	mainnetGenesis100DependencyName = "mainnet genesis (100M)"

	// genesis state files
	mainnetGenesisStateDependencyName    = "mainnet genesis state"
	testnetGenesisStateDependencyName    = "testnet genesis state"
	devnetGenesisStateDependencyName     = "devnet genesis state"
	mainnetGenesisState35DependencyName  = "mainnet genesis state (35M)"
	mainnetGenesisState42DependencyName  = "mainnet genesis state (42M)"
	mainnetGenesisState100DependencyName = "mainnet genesis state (100M)"

	// chain configurations
	mainnetChainConfigDependencyName = "mainnet chain config"
	testnetChainConfigDependencyName = "testnet chain config"
	devnetChainConfigDependencyName  = "devnet chain config"

	// client configurations
	gethMainnetConfigName = "geth mainnet config"
	gethTestnetConfigName = "geth testnet config"
	gethDevnetConfigName  = "geth devnet config"

	erigonMainnetConfigName = "erigon mainnet config"
	erigonTestnetConfigName = "erigon testnet config"
	erigonDevnetConfigName  = "erigon devnet config"

	prysmMainnetConfigDependencyName = "prysm mainnet config"
	prysmTestnetConfigDependencyName = "prysm testnet config"
	prysmDevnetConfigDependencyName  = "prysm devnet config"

	lighthouseMainnetConfigDependencyName  = "lighthouse mainnet config"
	lighthouseTestnetConfigDependencyName  = "lighthouse testnet config"
	lighthouseDevnetConfigDependencyName   = "lighthouse devnet config"
	deployBlockMainnetConfigDependencyName = "mainnet deploy block"
	deployBlockTestnetConfigDependencyName = "testnet deploy block"
	deployBlockDevnetConfigDependencyName  = "devnet deploy block"

	validatorMainnetConfigDependencyName = "validator mainnet config"
	validatorTestnetConfigDependencyName = "validator testnet config"
	validatorDevnetConfigDependencyName  = "validator devnet config"
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
	}

	// ----- CONFIGS -----
	sharedConfigDependencies = map[string]*ClientDependency{
		// ----- SHARED -----
		mainnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			name:     mainnetGenesisDependencyName,
			filePath: mainnetConfig + "/" + genesisJsonPath,
		},
		mainnetGenesis35DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_35.json",
			name:     mainnetGenesis35DependencyName,
			filePath: mainnetConfig + "/" + genesis35JsonPath,
		},
		mainnetGenesis42DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_42.json",
			name:     mainnetGenesis42DependencyName,
			filePath: mainnetConfig + "/" + genesis42JsonPath,
		},
		mainnetGenesis100DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_100.json",
			name:     mainnetGenesis100DependencyName,
			filePath: mainnetConfig + "/" + genesis100JsonPath,
		},
		mainnetGenesisStateDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz",
			name:     mainnetGenesisStateDependencyName,
			filePath: mainnetConfig + "/" + genesisStateFilePath,
		},
		mainnetGenesisState35DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_35.ssz",
			name:     mainnetGenesisState35DependencyName,
			filePath: mainnetConfig + "/" + genesisState35FilePath,
		},
		mainnetGenesisState42DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_42.ssz",
			name:     mainnetGenesisState42DependencyName,
			filePath: mainnetConfig + "/" + genesisState42FilePath,
		},
		mainnetGenesisState100DependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_100.ssz",
			name:     mainnetGenesisState100DependencyName,
			filePath: mainnetConfig + "/" + genesisState100FilePath,
		},
		mainnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			name:     mainnetChainConfigDependencyName,
			filePath: mainnetConfig + "/" + chainConfigYamlPath,
		},
		deployBlockMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deploy_block.txt",
			name:     deployBlockDevnetConfigDependencyName,
			filePath: mainnetConfig + "/" + deployBlockPath,
		},
		testnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.json",
			name:     testnetGenesisDependencyName,
			filePath: testnetConfig + "/" + genesisJsonPath,
		},
		testnetGenesisStateDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.ssz",
			name:     testnetGenesisStateDependencyName,
			filePath: testnetConfig + "/" + genesisStateFilePath,
		},
		testnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/config.yaml",
			name:     testnetChainConfigDependencyName,
			filePath: testnetConfig + "/" + chainConfigYamlPath,
		},
		deployBlockTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deploy_block.txt",
			name:     deployBlockDevnetConfigDependencyName,
			filePath: testnetConfig + "/" + deployBlockPath,
		},
		devnetGenesisDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/shared/genesis.json",
			name:     devnetGenesisDependencyName,
			filePath: devnetConfig + "/" + genesisJsonPath,
		},
		devnetGenesisStateDependencyName: {
			baseUrl:  "https://github.com/lukso-network/network-configs/raw/main/devnets/3030/shared/genesis.ssz",
			name:     devnetGenesisStateDependencyName,
			filePath: devnetConfig + "/" + genesisStateFilePath,
		},
		devnetChainConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/shared/config.yaml",
			name:     devnetChainConfigDependencyName,
			filePath: devnetConfig + "/" + chainConfigYamlPath,
		},
		deployBlockDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/lighthouse/deploy_block.txt",
			name:     deployBlockDevnetConfigDependencyName,
			filePath: devnetConfig + "/" + deployBlockPath,
		},
	}

	gethConfigDependencies = map[string]*ClientDependency{
		gethMainnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
			name:     gethMainnetConfigName,
			filePath: mainnetConfig + "/" + gethTomlPath,
		},
		gethTestnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
			name:     gethTestnetConfigName,
			filePath: testnetConfig + "/" + gethTomlPath,
		},
		gethDevnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/geth/geth.toml",
			name:     gethDevnetConfigName,
			filePath: devnetConfig + "/" + gethTomlPath,
		},
	}

	erigonConfigDependencies = map[string]*ClientDependency{
		erigonMainnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/erigon/erigon.toml",
			name:     erigonMainnetConfigName,
			filePath: mainnetConfig + "/" + erigonTomlPath,
		},
		erigonTestnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/erigon/erigon.toml",
			name:     erigonTestnetConfigName,
			filePath: testnetConfig + "/" + erigonTomlPath,
		},
		erigonDevnetConfigName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/erigon/erigon.toml",
			name:     erigonDevnetConfigName,
			filePath: devnetConfig + "/" + erigonTomlPath,
		},
	}

	prysmConfigDependencies = map[string]*ClientDependency{
		prysmMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/prysm.yaml",
			name:     prysmMainnetConfigDependencyName,
			filePath: mainnetConfig + "/" + prysmYamlPath,
		},
		prysmTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/prysm.yaml",
			name:     prysmTestnetConfigDependencyName,
			filePath: testnetConfig + "/" + prysmYamlPath,
		},
		prysmDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/prysm/prysm.yaml",
			name:     prysmDevnetConfigDependencyName,
			filePath: devnetConfig + "/" + prysmYamlPath,
		},
	}

	lighthouseConfigDependencies = map[string]*ClientDependency{
		lighthouseMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/lighthouse.toml",
			name:     lighthouseMainnetConfigDependencyName,
			filePath: mainnetConfig + "/" + lighthouseTomlPath,
		},
		lighthouseTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/lighthouse.toml",
			name:     lighthouseTestnetConfigDependencyName,
			filePath: testnetConfig + "/" + lighthouseTomlPath,
		},
		lighthouseDevnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/devnets/3030/lighthouse/lighthouse.toml",
			name:     lighthouseDevnetConfigDependencyName,
			filePath: devnetConfig + "/" + lighthouseTomlPath,
		},
	}

	validatorConfigDependencies = map[string]*ClientDependency{
		validatorMainnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/validator.yaml",
			name:     validatorMainnetConfigDependencyName,
			filePath: mainnetConfig + "/" + validatorYamlPath,
		},
		validatorTestnetConfigDependencyName: {
			baseUrl:  "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/validator.yaml",
			name:     validatorTestnetConfigDependencyName,
			filePath: testnetConfig + "/" + validatorYamlPath,
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
