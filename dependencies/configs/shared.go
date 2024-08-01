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

	validatorMainnetConfigDependencyName           = "validator mainnet config"
	validatorTestnetConfigDependencyName           = "validator testnet config"
	lighthouseValidatorMainnetConfigDependencyName = "lighthouse validator mainnet config"
	lighthouseValidatorTestnetConfigDependencyName = "lighthouse validator testnet config"
	tekuValidatorMainnetConfigDependencyName       = "teku validator mainnet config"
	tekuValidatorTestnetConfigDependencyName       = "teku validator testnet config"

	tekuMainnetConfigDependencyName      = "teku mainnet config"
	tekuTestnetConfigDependencyName      = "teku testnet config"
	tekuMainnetChainConfigDependencyName = "teku mainnet chain config"
	tekuTestnetChainConfigDependencyName = "teku testnet chain config"

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

	ChainConfigYamlPath         = "shared/config.yaml"
	DeployBlockPath             = "shared/deploy_block.txt"
	DepositContractBlockPath    = "shared/deposit_contract_block.txt"
	GethTomlPath                = "geth/geth.toml"
	ErigonTomlPath              = "erigon/erigon.toml"
	NethermindCfgPath           = "nethermind/nethermind.cfg"
	BesuTomlPath                = "besu/besu.toml"
	PrysmYamlPath               = "prysm/prysm.yaml"
	ValidatorYamlPath           = "prysm/validator.yaml"
	LighthouseTomlPath          = "lighthouse/lighthouse.toml"
	LighthouseValidatorTomlPath = "lighthouse/validator.toml"
	TekuYamlPath                = "teku/teku.yaml"
	TekuValidatorYamlPath       = "teku/validator.yaml"
	TekuChainConfigYamlPath     = "teku/config.yaml"
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
		// same as SharedConfigDependencies but with added Teku chain configs
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
		tekuMainnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/config.yaml",
			name:     tekuMainnetChainConfigDependencyName,
			filePath: MainnetConfig + "/" + TekuChainConfigYamlPath,
		},
		tekuTestnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/config.yaml",
			name:     tekuTestnetChainConfigDependencyName,
			filePath: TestnetConfig + "/" + TekuChainConfigYamlPath,
		},
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
	return c.name
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
