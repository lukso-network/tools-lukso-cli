package configs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
)

const (
	BinaryPerms   = os.ModePerm
	ConfigRootDir = "./configs"

	// NAMES
	ExecutionLayer = "execution"
	ConsensusLayer = "consensus"
	ValidatorLayer = "validator"

	// genesis file names
	mainnetGenesisDependencyName = "mainnet_genesis"
	testnetGenesisDependencyName = "testnet_genesis"

	// nethermind chainspec
	mainnetGenesisChainspecDependencyName = "mainnet_genesis_chainspec"
	testnetGenesisChainspecDependencyName = "testnet_genesis_chainspec"

	// genesis state file names
	mainnetGenesisStateDependencyName = "mainnet_genesis_state"
	testnetGenesisStateDependencyName = "testnet genesis_state"

	// chain configurations
	mainnetChainConfigDependencyName = "mainnet_chain_config"
	testnetChainConfigDependencyName = "testnet_chain_config"

	// client configurations
	gethMainnetConfigName = "geth_mainnet_config"
	gethTestnetConfigName = "geth_testnet_config"

	erigonMainnetConfigName = "erigon_mainnet_config"
	erigonTestnetConfigName = "erigon_testnet_config"

	besuMainnetConfigName = "besu_mainnet_config"
	besuTestnetConfigName = "besu_testnet_config"

	nethermindMainnetConfigName = "nethermind_mainnet_config"
	nethermindTestnetConfigName = "nethermind_testnet_config"

	prysmMainnetConfigDependencyName = "prysm_mainnet_config"
	prysmTestnetConfigDependencyName = "prysm_testnet_config"

	lighthouseMainnetConfigDependencyName           = "lighthouse_mainnet_config"
	lighthouseTestnetConfigDependencyName           = "lighthouse_testnet_config"
	deployBlockMainnetConfigDependencyName          = "mainnet_deploy_block"
	deployBlockTestnetConfigDependencyName          = "testnet_deploy_block"
	depositContractBlockMainnetConfigDependencyName = "mainnet_deposit_contract_block"
	depositContractBlockTestnetConfigDependencyName = "testnet_deposit_contract_block"

	tekuMainnetConfigDependencyName      = "teku_mainnet_config"
	tekuTestnetConfigDependencyName      = "teku_testnet_config"
	tekuMainnetChainConfigDependencyName = "teku_mainnet_chain_config"
	tekuTestnetChainConfigDependencyName = "teku_testnet_chain_config"

	nimbus2MainnetConfigDependencyName                   = "nimbus2_mainnet_config"
	nimbus2TestnetConfigDependencyName                   = "nimbus2_testnet_config"
	nimbus2MainnetChainConfigDependencyName              = "nimbus2_mainnet_chain_config"
	nimbus2TestnetChainConfigDependencyName              = "nimbus2_testnet_chain_config"
	nimbus2MainnetDepositContractBlockHashDependencyName = "nimbus2_mainnet_deposit_contract_block_hash"
	nimbus2TestnetDepositContractBlockHashDependencyName = "nimbus2_testnet_deposit_contract_block_hash"
	nimbus2MainnetDeployBlockDependencyName              = "nimbus2_mainnet_deploy_block "
	nimbus2TestnetDeployBlockDependencyName              = "nimbus2_testnet_deploy_block "
	nimbus2MainnetBootnodesDependencyName                = "nimbus2_mainnet_bootnodes"
	nimbus2TestnetBootnodesDependencyName                = "nimbus2_testnet_bootnodes"

	prysmValidatorMainnetConfigDependencyName      = "prysm_validator_mainnet_config"
	prysmValidatorTestnetConfigDependencyName      = "prysm_validator_testnet_config"
	lighthouseValidatorMainnetConfigDependencyName = "lighthouse_validator_mainnet_config"
	lighthouseValidatorTestnetConfigDependencyName = "lighthouse_validator_testnet_config"
	tekuValidatorMainnetConfigDependencyName       = "teku_validator_mainnet_config"
	tekuValidatorTestnetConfigDependencyName       = "teku_validator_testnet_config"
	nimbus2ValidatorMainnetConfigDependencyName    = "nimbus2_validator_mainnet_config"
	nimbus2ValidatorTestnetConfigDependencyName    = "nimbus2_validator_testnet_config"

	// PATHS
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

	// shared/network configs
	GenesisJsonPath          = "shared/genesis.json"
	GenesisStateFilePath     = "shared/genesis.ssz"
	ChainConfigYamlPath      = "shared/config.yaml"
	DeployBlockPath          = "shared/deploy_block.txt"
	DepositContractBlockPath = "shared/deposit_contract_block.txt"
	GenesisChainspecPath     = "nethermind/chainspec.json"

	GethTomlPath                        = "geth/geth.toml"
	ErigonTomlPath                      = "erigon/erigon.toml"
	NethermindJsonPath                  = "nethermind/nethermind.json"
	BesuTomlPath                        = "besu/besu.toml"
	PrysmYamlPath                       = "prysm/prysm.yaml"
	LighthouseTomlPath                  = "lighthouse/lighthouse.toml"
	TekuYamlPath                        = "teku/teku.yaml"
	TekuChainConfigYamlPath             = "teku/config.yaml"
	Nimbus2Path                         = "nimbus2" // useful when the path to the nimbus configs is needed, e.g. in --network flag
	Nimbus2TomlPath                     = "nimbus2/nimbus.toml"
	Nimbus2ChainConfigYamlPath          = "nimbus2/config.yaml"
	Nimbus2DeployBlockPath              = "nimbus2/deploy_block.txt"
	Nimbus2DepositContractBlockHashPath = "nimbus2/deposit_contract_block_hash.txt"
	Nimbus2Bootnodes                    = "nimbus2/bootnodes.txt"
	PrysmValidatorYamlPath              = "prysm/validator.yaml"
	LighthouseValidatorTomlPath         = "lighthouse/validator.toml"
	TekuValidatorYamlPath               = "teku/validator.yaml"
	Nimbus2ValidatorTomlPath            = "nimbus2/validator.toml"
)

var (
	SharedConfigDependencies = map[string]ClientConfigDependency{
		// ----- SHARED -----
		mainnetGenesisDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.json",
			name:     mainnetGenesisDependencyName,
			filePath: MainnetConfig + "/" + GenesisJsonPath,
		},
		mainnetGenesisStateDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/genesis.ssz",
			name:     mainnetGenesisStateDependencyName,
			filePath: MainnetConfig + "/" + GenesisStateFilePath,
		},
		mainnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/shared/config.yaml",
			name:     mainnetChainConfigDependencyName,
			filePath: MainnetConfig + "/" + ChainConfigYamlPath,
		},
		mainnetGenesisChainspecDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nethermind/chainspec.json",
			name:     mainnetGenesisChainspecDependencyName,
			filePath: MainnetConfig + "/" + GenesisChainspecPath,
		},
		deployBlockMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deploy_block.txt",
			name:     deployBlockMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + DeployBlockPath,
		},
		depositContractBlockMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/deposit_contract_block.txt",
			name:     depositContractBlockMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + DepositContractBlockPath,
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
		testnetGenesisChainspecDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nethermind/chainspec.json",
			name:     mainnetGenesisChainspecDependencyName,
			filePath: TestnetConfig + "/" + GenesisChainspecPath,
		},
		deployBlockTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deploy_block.txt",
			name:     deployBlockTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + DeployBlockPath,
		},
		depositContractBlockTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/deposit_contract_block.txt",
			name:     depositContractBlockTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + DepositContractBlockPath,
		},
	}
	UpdateConfigDependencies = map[string]ClientConfigDependency{
		// copied existing configs
		// mainnet
		mainnetGenesisDependencyName:                    SharedConfigDependencies[mainnetGenesisDependencyName],
		mainnetGenesisStateDependencyName:               SharedConfigDependencies[mainnetGenesisStateDependencyName],
		mainnetChainConfigDependencyName:                SharedConfigDependencies[mainnetChainConfigDependencyName],
		mainnetGenesisChainspecDependencyName:           SharedConfigDependencies[mainnetGenesisChainspecDependencyName],
		deployBlockMainnetConfigDependencyName:          SharedConfigDependencies[deployBlockMainnetConfigDependencyName],
		depositContractBlockMainnetConfigDependencyName: SharedConfigDependencies[depositContractBlockMainnetConfigDependencyName],
		tekuMainnetChainConfigDependencyName:            TekuConfigDependencies[tekuMainnetChainConfigDependencyName],
		nimbus2MainnetChainConfigDependencyName:         Nimbus2ConfigDependencies[nimbus2MainnetChainConfigDependencyName],
		// testnet
		testnetGenesisDependencyName:                    SharedConfigDependencies[testnetGenesisDependencyName],
		testnetGenesisStateDependencyName:               SharedConfigDependencies[testnetGenesisStateDependencyName],
		testnetChainConfigDependencyName:                SharedConfigDependencies[testnetChainConfigDependencyName],
		testnetGenesisChainspecDependencyName:           SharedConfigDependencies[testnetGenesisChainspecDependencyName],
		deployBlockTestnetConfigDependencyName:          SharedConfigDependencies[deployBlockTestnetConfigDependencyName],
		depositContractBlockTestnetConfigDependencyName: SharedConfigDependencies[depositContractBlockTestnetConfigDependencyName],
		tekuTestnetChainConfigDependencyName:            TekuConfigDependencies[tekuTestnetChainConfigDependencyName],
		nimbus2TestnetChainConfigDependencyName:         Nimbus2ConfigDependencies[nimbus2TestnetChainConfigDependencyName],
	}
)

type clientConfig struct {
	url      string
	name     string
	filePath string
}

var _ ClientConfigDependency = &clientConfig{}

func (c *clientConfig) Install(isUpdate bool) (err error) {
	err = c.createDir()
	if err != nil {
		return
	}

	if utils.FileExists(c.filePath) && !isUpdate {
		log.Infof("  ⏩️  Skipping file %s: the file already exists", c.filePath)

		return
	}

	response, err := http.Get(c.url)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", c.url)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			c.url,
			response.Status,
		)
	}

	var responseReader io.Reader = response.Body

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(responseReader)
	if err != nil {
		return
	}

	err = os.WriteFile(c.filePath, buf.Bytes(), common.ConfigPerms)

	if err != nil && strings.Contains(err.Error(), "Permission denied") {
		return errors.ErrNeedRoot
	}

	if err != nil {
		log.Infof("❌  Couldn't save file: %v", err)

		return
	}

	return
}

func (c *clientConfig) Name() string {
	return strings.ReplaceAll(c.name, "_", " ")
}

func (c *clientConfig) createDir() error {
	err := os.MkdirAll(utils.TruncateFileFromDir(c.filePath), common.ConfigPerms)
	if err == os.ErrExist {
		log.Errorf("%s already exists!", c.name)
	}

	if err == os.ErrPermission {
		return errors.ErrNeedRoot
	}

	return err
}
