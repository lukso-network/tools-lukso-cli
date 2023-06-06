package configs

import "os"

const (
	BinaryPerms   = os.ModePerm
	ConfigPerms   = 0750
	ConfigRootDir = "./configs"

	MainnetConfig = ConfigRootDir + "/mainnet"
	TestnetConfig = ConfigRootDir + "/testnet"

	MainnetLogs = "./mainnet-logs"
	TestnetLogs = "./testnet-logs"

	MainnetKeystore = "./mainnet-keystore"
	TestnetKeystore = "./testnet-keystore"

	MainnetDatadir = "./mainnet-data"
	TestnetDatadir = "./mainnet-data"

	ExecutionMainnetDatadir = MainnetDatadir + "/execution"
	ExecutionTestnetDatadir = TestnetDatadir + "/execution"

	ConsensusMainnetDatadir = MainnetDatadir + "/consensus"
	ConsensusTestnetDatadir = TestnetDatadir + "/consensus"

	ValidatorMainnetDatadir = MainnetDatadir + "/validator"
	ValidatorTestnetDatadir = TestnetDatadir + "/validator"

	ExecutionLayer = "execution"
	ConsensusLayer = "consensus"
	ValidatorLayer = "validator"

	// genesis file names
	mainnetGenesisDependencyName    = "mainnet genesis"
	testnetGenesisDependencyName    = "testnet genesis"
	mainnetGenesis35DependencyName  = "mainnet genesis (35M)"
	mainnetGenesis42DependencyName  = "mainnet genesis (42M)"
	mainnetGenesis100DependencyName = "mainnet genesis (100M)"

	// genesis state file names
	mainnetGenesisStateDependencyName    = "mainnet genesis state"
	testnetGenesisStateDependencyName    = "testnet genesis state"
	mainnetGenesisState35DependencyName  = "mainnet genesis state (35M)"
	mainnetGenesisState42DependencyName  = "mainnet genesis state (42M)"
	mainnetGenesisState100DependencyName = "mainnet genesis state (100M)"

	GenesisJsonPath    = "shared/genesis.json"
	Genesis35JsonPath  = "shared/genesis_35.json"
	Genesis42JsonPath  = "shared/genesis_42.json"
	Genesis100JsonPath = "shared/genesis_100.json"

	GenesisStateFilePath    = "shared/genesis.ssz"
	GenesisState35FilePath  = "shared/genesis_35.ssz"
	GenesisState42FilePath  = "shared/genesis_42.ssz"
	GenesisState100FilePath = "shared/genesis_100.ssz"

	// chain configurations
	mainnetChainConfigDependencyName = "mainnet chain config"
	testnetChainConfigDependencyName = "testnet chain config"

	// client configurations
	gethMainnetConfigName = "geth mainnet config"
	gethTestnetConfigName = "geth testnet config"

	erigonMainnetConfigName = "erigon mainnet config"
	erigonTestnetConfigName = "erigon testnet config"

	prysmMainnetConfigDependencyName = "prysm mainnet config"
	prysmTestnetConfigDependencyName = "prysm testnet config"

	lighthouseMainnetConfigDependencyName  = "lighthouse mainnet config"
	lighthouseTestnetConfigDependencyName  = "lighthouse testnet config"
	deployBlockMainnetConfigDependencyName = "mainnet deploy block"
	deployBlockTestnetConfigDependencyName = "testnet deploy block"

	validatorMainnetConfigDependencyName           = "validator mainnet config"
	validatorTestnetConfigDependencyName           = "validator testnet config"
	lighthouseValidatorMainnetConfigDependencyName = "lighthouse validator mainnet config"
	lighthouseValidatorTestnetConfigDependencyName = "lighthouse validator testnet config"

	ChainConfigYamlPath         = "shared/config.yaml"
	GethTomlPath                = "geth/geth.toml"
	ErigonTomlPath              = "erigon/erigon.toml"
	PrysmYamlPath               = "prysm/prysm.yaml"
	LighthouseTomlPath          = "lighthouse/lighthouse.toml"
	LighthouseValidatorTomlPath = "lighthouse/validator.toml"
	DeployBlockPath             = "shared/deploy_block.txt"
	ValidatorYamlPath           = "prysm/validator.yaml"
)

var (
	SharedConfigDependencies = map[string]ClientConfigDependency{
		// ----- SHARED -----
		mainnetGenesisDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			name:     mainnetGenesisDependencyName,
			filePath: MainnetConfig + "/" + GenesisJsonPath,
		},
		mainnetGenesis35DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_35.json",
			name:     mainnetGenesis35DependencyName,
			filePath: MainnetConfig + "/" + Genesis35JsonPath,
		},
		mainnetGenesis42DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_42.json",
			name:     mainnetGenesis42DependencyName,
			filePath: MainnetConfig + "/" + Genesis42JsonPath,
		},
		mainnetGenesis100DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_100.json",
			name:     mainnetGenesis100DependencyName,
			filePath: MainnetConfig + "/" + Genesis100JsonPath,
		},
		mainnetGenesisStateDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz",
			name:     mainnetGenesisStateDependencyName,
			filePath: MainnetConfig + "/" + GenesisStateFilePath,
		},
		mainnetGenesisState35DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_35.ssz",
			name:     mainnetGenesisState35DependencyName,
			filePath: MainnetConfig + "/" + GenesisState35FilePath,
		},
		mainnetGenesisState42DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_42.ssz",
			name:     mainnetGenesisState42DependencyName,
			filePath: MainnetConfig + "/" + GenesisState42FilePath,
		},
		mainnetGenesisState100DependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis_100.ssz",
			name:     mainnetGenesisState100DependencyName,
			filePath: MainnetConfig + "/" + GenesisState100FilePath,
		},
		mainnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			name:     mainnetChainConfigDependencyName,
			filePath: MainnetConfig + "/" + ChainConfigYamlPath,
		},
		deployBlockMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deploy_block.txt",
			name:     deployBlockMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + DeployBlockPath,
		},
		testnetGenesisDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.json",
			name:     testnetGenesisDependencyName,
			filePath: TestnetConfig + "/" + GenesisJsonPath,
		},
		testnetGenesisStateDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/genesis.ssz",
			name:     testnetGenesisStateDependencyName,
			filePath: TestnetConfig + "/" + GenesisStateFilePath,
		},
		testnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/shared/config.yaml",
			name:     testnetChainConfigDependencyName,
			filePath: TestnetConfig + "/" + ChainConfigYamlPath,
		},
		deployBlockTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deploy_block.txt",
			name:     deployBlockTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + DeployBlockPath,
		},
	}
)

type clientConfig struct {
	url      string
	name     string
	filePath string
}

var _ ClientConfigDependency = &clientConfig{}

func (c *clientConfig) Install() (err error) {
	return
}

func (c *clientConfig) Name() string {
	return c.name
}
