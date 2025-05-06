package configs

import (
	"maps"
	"os"
	"slices"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
)

const (
	BinaryPerms   = os.ModePerm
	ConfigRootDir = "./configs"

	// NAMES
	ExecutionLayer = "execution"
	ConsensusLayer = "consensus"
	ValidatorLayer = "validator"

	// genesis file names
	mainnetGenesisDependencyName    = "mainnet genesis"
	testnetGenesisDependencyName    = "testnet genesis"
	mainnetGenesis35DependencyName  = "mainnet genesis (35M)"
	mainnetGenesis42DependencyName  = "mainnet genesis (42M)"
	mainnetGenesis100DependencyName = "mainnet genesis (100M)"

	// nethermind chainspec
	mainnetGenesisChainspecDependencyName = "mainnet genesis chainspec"
	testnetGenesisChainspecDependencyName = "testnet genesis chainspec"

	// genesis state file names
	mainnetGenesisStateDependencyName    = "mainnet genesis state"
	testnetGenesisStateDependencyName    = "testnet genesis state"
	mainnetGenesisState35DependencyName  = "mainnet genesis state (35M)"
	mainnetGenesisState42DependencyName  = "mainnet genesis state (42M)"
	mainnetGenesisState100DependencyName = "mainnet genesis state (100M)"

	// chain configurations
	mainnetChainConfigDependencyName = "mainnet chain config"
	testnetChainConfigDependencyName = "testnet chain config"

	// client configurations
	gethMainnetConfigName = "geth mainnet config"
	gethTestnetConfigName = "geth testnet config"

	erigonMainnetConfigName = "erigon mainnet config"
	erigonTestnetConfigName = "erigon testnet config"

	nethermindMainnetConfigName = "nethermind mainnet config"
	nethermindTestnetConfigName = "nethermind testnet config"

	besuMainnetConfigName = "besu mainnet config"
	besuTestnetConfigName = "besu testnet config"

	prysmMainnetConfigDependencyName = "prysm mainnet config"
	prysmTestnetConfigDependencyName = "prysm testnet config"

	lighthouseMainnetConfigDependencyName           = "lighthouse mainnet config"
	lighthouseTestnetConfigDependencyName           = "lighthouse testnet config"
	deployBlockMainnetConfigDependencyName          = "mainnet deploy block"
	deployBlockTestnetConfigDependencyName          = "testnet deploy block"
	depositContractBlockMainnetConfigDependencyName = "mainnet deposit contract block"
	depositContractBlockTestnetConfigDependencyName = "testnet deposit contract block"

	tekuMainnetConfigDependencyName      = "teku mainnet config"
	tekuTestnetConfigDependencyName      = "teku testnet config"
	tekuMainnetChainConfigDependencyName = "teku mainnet chain config"
	tekuTestnetChainConfigDependencyName = "teku testnet chain config"

	nimbus2MainnetConfigDependencyName                   = "nimbus2 mainnet config"
	nimbus2TestnetConfigDependencyName                   = "nimbus2 testnet config"
	nimbus2MainnetChainConfigDependencyName              = "nimbus2 mainnet chain config"
	nimbus2TestnetChainConfigDependencyName              = "nimbus2 testnet chain config"
	nimbus2MainnetDepositContractBlockHashDependencyName = "nimbus2 mainnet deposit contract block hash"
	nimbus2TestnetDepositContractBlockHashDependencyName = "nimbus2 testnet deposit contract block hash"
	nimbus2MainnetDeployBlockDependencyName              = "nimbus2 mainnet deploy block "
	nimbus2TestnetDeployBlockDependencyName              = "nimbus2 testnet deploy block "
	nimbus2MainnetBootnodesDependencyName                = "nimbus2 mainnet bootnodes"
	nimbus2TestnetBootnodesDependencyName                = "nimbus2 testnet bootnodes"

	validatorMainnetConfigDependencyName           = "validator mainnet config"
	validatorTestnetConfigDependencyName           = "validator testnet config"
	lighthouseValidatorMainnetConfigDependencyName = "lighthouse validator mainnet config"
	lighthouseValidatorTestnetConfigDependencyName = "lighthouse validator testnet config"
	tekuValidatorMainnetConfigDependencyName       = "teku validator mainnet config"
	tekuValidatorTestnetConfigDependencyName       = "teku validator testnet config"
	nimbus2ValidatorMainnetConfigDependencyName    = "nimbus2 validator mainnet config"
	nimbus2ValidatorTestnetConfigDependencyName    = "nimbus2 validator testnet config"

	// PATHS
	GenesisJsonPath    = "shared/genesis.json"
	Genesis35JsonPath  = "shared/genesis_35.json"
	Genesis42JsonPath  = "shared/genesis_42.json"
	Genesis100JsonPath = "shared/genesis_100.json"

	GenesisStateFilePath    = "shared/genesis.ssz"
	GenesisState35FilePath  = "shared/genesis_35.ssz"
	GenesisState42FilePath  = "shared/genesis_42.ssz"
	GenesisState100FilePath = "shared/genesis_100.ssz"

	GenesisChainspecPath = "nethermind/chainspec.json"

	MainnetConfig = ConfigRootDir + "/mainnet"
	TestnetConfig = ConfigRootDir + "/testnet"

	MainnetLogs = "./mainnet-logs"
	TestnetLogs = "./testnet-logs"

	MainnetKeystore = "./mainnet-keystore"
	TestnetKeystore = "./testnet-keystore"

	MainnetDatadir = "./mainnet-data"
	TestnetDatadir = "./testnet-data"

	ExecutionMainnetDatadir = MainnetDatadir + "/execution"
	ExecutionTestnetDatadir = TestnetDatadir + "/execution"

	ConsensusMainnetDatadir = MainnetDatadir + "/consensus"
	ConsensusTestnetDatadir = TestnetDatadir + "/consensus"

	ValidatorMainnetDatadir = MainnetDatadir + "/validator"
	ValidatorTestnetDatadir = TestnetDatadir + "/validator"

	ChainConfigYamlPath                 = "shared/config.yaml"
	DeployBlockPath                     = "shared/deploy_block.txt"
	DepositContractBlockPath            = "shared/deposit_contract_block.txt"
	GethTomlPath                        = "geth/geth.toml"
	ErigonTomlPath                      = "erigon/erigon.toml"
	NethermindJsonPath                  = "nethermind/nethermind.json"
	BesuTomlPath                        = "besu/besu.toml"
	PrysmYamlPath                       = "prysm/prysm.yaml"
	ValidatorYamlPath                   = "prysm/validator.yaml"
	LighthouseTomlPath                  = "lighthouse/lighthouse.toml"
	LighthouseValidatorTomlPath         = "lighthouse/validator.toml"
	TekuYamlPath                        = "teku/teku.yaml"
	TekuValidatorYamlPath               = "teku/validator.yaml"
	TekuChainConfigYamlPath             = "teku/config.yaml"
	Nimbus2Path                         = "nimbus2" // useful when the path to the nimbus configs is needed, e.g. in --network flag
	Nimbus2TomlPath                     = "nimbus2/nimbus.toml"
	Nimbus2ValidatorTomlPath            = "nimbus2/validator.toml"
	Nimbus2ChainConfigYamlPath          = "nimbus2/config.yaml"
	Nimbus2DeployBlockPath              = "nimbus2/deploy_block.txt"
	Nimbus2DepositContractBlockHashPath = "nimbus2/deposit_contract_block_hash.txt"
	Nimbus2Bootnodes                    = "nimbus2/bootnodes.txt"
)

var (
	CfgFileManager = file.NewManager()
	CfgInstaller   = installer.NewInstaller(CfgFileManager)

	SharedConfigDependencies = map[string]ClientConfigDependency{
		// ----- SHARED -----
		mainnetGenesisDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			mainnetGenesisDependencyName,
			MainnetConfig+"/"+GenesisJsonPath,
		),
		mainnetGenesisStateDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz",
			mainnetGenesisStateDependencyName,
			MainnetConfig+"/"+GenesisStateFilePath,
		),
		mainnetChainConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			mainnetChainConfigDependencyName,
			MainnetConfig+"/"+ChainConfigYamlPath,
		),
		mainnetGenesisChainspecDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nethermind/chainspec.json",
			mainnetGenesisChainspecDependencyName,
			MainnetConfig+"/"+GenesisChainspecPath,
		),
		deployBlockMainnetConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deploy_block.txt",
			deployBlockMainnetConfigDependencyName,
			MainnetConfig+"/"+DeployBlockPath,
		),
		depositContractBlockMainnetConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deposit_contract_block.txt",
			depositContractBlockMainnetConfigDependencyName,
			MainnetConfig+"/"+DepositContractBlockPath,
		),
		testnetGenesisDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.json",
			testnetGenesisDependencyName,
			TestnetConfig+"/"+GenesisJsonPath,
		),
		testnetGenesisStateDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.ssz",
			testnetGenesisStateDependencyName,
			TestnetConfig+"/"+GenesisStateFilePath,
		),
		testnetChainConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/config.yaml",
			testnetChainConfigDependencyName,
			TestnetConfig+"/"+ChainConfigYamlPath,
		),
		testnetGenesisChainspecDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nethermind/chainspec.json",
			mainnetGenesisChainspecDependencyName,
			TestnetConfig+"/"+GenesisChainspecPath,
		),
		deployBlockTestnetConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deploy_block.txt",
			deployBlockTestnetConfigDependencyName,
			TestnetConfig+"/"+DeployBlockPath,
		),
		depositContractBlockTestnetConfigDependencyName: newClientConfig(
			"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deposit_contract_block.txt",
			depositContractBlockTestnetConfigDependencyName,
			TestnetConfig+"/"+DepositContractBlockPath,
		),
	}
	UpdateConfigDependencies = map[string]ClientConfigDependency{
		// copied existing configs
		mainnetGenesisDependencyName:                    SharedConfigDependencies[mainnetGenesisDependencyName],
		mainnetGenesisStateDependencyName:               SharedConfigDependencies[mainnetGenesisStateDependencyName],
		mainnetChainConfigDependencyName:                SharedConfigDependencies[mainnetChainConfigDependencyName],
		mainnetGenesisChainspecDependencyName:           SharedConfigDependencies[mainnetGenesisChainspecDependencyName],
		deployBlockMainnetConfigDependencyName:          SharedConfigDependencies[deployBlockMainnetConfigDependencyName],
		depositContractBlockMainnetConfigDependencyName: SharedConfigDependencies[depositContractBlockMainnetConfigDependencyName],
		testnetGenesisDependencyName:                    SharedConfigDependencies[testnetGenesisDependencyName],
		testnetGenesisStateDependencyName:               SharedConfigDependencies[testnetGenesisStateDependencyName],
		testnetChainConfigDependencyName:                SharedConfigDependencies[testnetChainConfigDependencyName],
		testnetGenesisChainspecDependencyName:           SharedConfigDependencies[testnetGenesisChainspecDependencyName],
		deployBlockTestnetConfigDependencyName:          SharedConfigDependencies[deployBlockTestnetConfigDependencyName],
		depositContractBlockTestnetConfigDependencyName: SharedConfigDependencies[depositContractBlockTestnetConfigDependencyName],
		tekuMainnetChainConfigDependencyName:            TekuConfigDependencies[tekuMainnetChainConfigDependencyName],
		tekuTestnetChainConfigDependencyName:            TekuConfigDependencies[tekuTestnetChainConfigDependencyName],
		nimbus2MainnetChainConfigDependencyName:         Nimbus2ConfigDependencies[nimbus2MainnetChainConfigDependencyName],
		nimbus2TestnetChainConfigDependencyName:         Nimbus2ConfigDependencies[nimbus2TestnetChainConfigDependencyName],
	}

	// An array of all unique configs, shared and client specific - used for progress calculations
	AllDependencies = makeAllDependencies()
)

type clientConfig struct {
	url      string
	name     string
	filePath string
}

func newClientConfig(url, name, filePath string) ClientConfigDependency {
	return &clientConfig{
		url:      url,
		name:     name,
		filePath: filePath,
	}
}

var _ ClientConfigDependency = &clientConfig{}

func (c *clientConfig) Install(isUpdate bool) (err error) {
	if CfgFileManager.Exists(c.filePath) && !isUpdate {
		return errors.ErrFileExists
	}

	err = CfgInstaller.InstallFile(c.url, c.filePath)
	if err != nil {
		return
	}

	return
}

func (c *clientConfig) Name() string {
	return c.name
}

func makeAllDependencies() (deps []ClientConfigDependency) {
	depsMap := []map[string]ClientConfigDependency{
		SharedConfigDependencies,
		GethConfigDependencies,
		ErigonConfigDependencies,
		NethermindConfigDependencies,
		BesuConfigDependencies,
		PrysmConfigDependencies,
		LighthouseConfigDependencies,
		TekuConfigDependencies,
		Nimbus2ConfigDependencies,
	}

	for _, dep := range depsMap {
		deps = append(deps, slices.Collect(maps.Values(dep))...)
	}

	return
}
