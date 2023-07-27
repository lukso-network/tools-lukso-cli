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

	tekuMainnetConfigDependencyName      = "teku mainnet config"
	tekuTestnetConfigDependencyName      = "teku testnet config"
	tekuMainnetChainConfigDependencyName = "teku mainnet chain config"
	tekuTestnetChainConfigDependencyName = "teku testnet chain config"

	ChainConfigYamlPath         = "shared/config.yaml"
	GethTomlPath                = "geth/geth.toml"
	ErigonTomlPath              = "erigon/erigon.toml"
	PrysmYamlPath               = "prysm/prysm.yaml"
	LighthouseTomlPath          = "lighthouse/lighthouse.toml"
	LighthouseValidatorTomlPath = "lighthouse/validator.toml"
	DeployBlockPath             = "shared/deploy_block.txt"
	ValidatorYamlPath           = "prysm/validator.yaml"
	TekuYamlPath                = "teku/teku.yaml"
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
	err = c.createDir()
	if err != nil {
		return
	}

	if utils.FileExists(c.filePath) {
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
