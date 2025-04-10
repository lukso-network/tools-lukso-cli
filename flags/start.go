package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

var (
	StartFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  ValidatorFlag,
			Usage: "Run LUKSO node with validator",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  NoSlasherFlag,
			Usage: "Disables slasher",
			Value: false, // default is true, we change it to false only when running validator
		},
		&cli.StringFlag{
			Name:  TransactionFeeRecipientFlag,
			Usage: "Address to receive fees from blocks",
		},
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: configs.MainnetLogs,
		},
		&cli.BoolFlag{
			Name:  CheckpointSyncFlag,
			Usage: "Run a node with checkpoint sync feature",
		},
	}
	GethStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   GethDatadirFlag,
			Usage:  "Geth datadir",
			Value:  configs.ExecutionMainnetDatadir,
			Hidden: true,
		},
		&cli.StringFlag{
			Name:  GethConfigFileFlag,
			Usage: "Path to geth.toml config file",
			Value: configs.MainnetConfig + "/" + configs.GethTomlPath,
		},
		&cli.StringFlag{
			Name:  GenesisJsonFlag,
			Usage: "Path to genesis.json file",
			Value: configs.MainnetConfig + "/" + configs.GenesisJsonPath,
		},
	}
	ErigonStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ErigonConfigFileFlag,
			Usage: "Path to erigon.toml config file",
			Value: configs.MainnetConfig + "/" + configs.ErigonTomlPath,
		},
		&cli.StringFlag{
			Name:  ErigonDatadirFlag,
			Usage: "Erigon datadir",
			Value: configs.ExecutionMainnetDatadir,
		},
	}
	NethermindStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  NethermindConfigFileFlag,
			Usage: "Path to nethermind.json config file",
			Value: configs.MainnetConfig + "/" + configs.NethermindJsonPath,
		},
		&cli.StringFlag{
			Name:  NethermindDatadirFlag,
			Usage: "Nethermind datadir",
			Value: configs.ExecutionMainnetDatadir,
		},
	}
	BesuStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  BesuConfigFileFlag,
			Usage: "Path to besu.toml config file",
			Value: configs.MainnetConfig + "/" + configs.BesuTomlPath,
		},
		&cli.StringFlag{
			Name:  BesuDatadirFlag,
			Usage: "Besu datadir",
			Value: configs.ExecutionMainnetDatadir,
		},
	}
	PrysmStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  GenesisStateFlag,
			Usage: "Path to genesis.ssz file",
			Value: configs.MainnetConfig + "/" + configs.GenesisStateFilePath,
		},
		&cli.StringFlag{
			Name:   PrysmDatadirFlag,
			Usage:  "Prysm datadir",
			Value:  configs.ConsensusMainnetDatadir,
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   PrysmChainConfigFileFlag,
			Usage:  "Path to chain config file",
			Value:  configs.MainnetConfig + "/" + configs.ChainConfigYamlPath,
			Hidden: true,
		},
		&cli.StringFlag{
			Name:  PrysmConfigFileFlag,
			Usage: "Path to prysm.yaml config file",
			Value: configs.MainnetConfig + "/" + configs.PrysmYamlPath,
		},
	}
	LighthouseStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LighthouseConfigFileFlag,
			Usage: "Path to lighthouse.toml config file",
			Value: configs.MainnetConfig + "/" + configs.LighthouseTomlPath,
		},
		&cli.StringFlag{
			Name:  LighthouseDatadirFlag,
			Usage: "Lighthouse datadir",
			Value: configs.ConsensusMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  LighthouseValidatorConfigFileFlag,
			Usage: "Path to validator.toml config file",
			Value: configs.MainnetConfig + "/" + configs.LighthouseValidatorTomlPath,
		},
	}
	TekuStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  TekuConfigFileFlag,
			Usage: "Path to teku.yaml config file",
			Value: configs.MainnetConfig + "/" + configs.TekuYamlPath,
		},
		&cli.StringFlag{
			Name:  TekuValidatorConfigFileFlag,
			Usage: "Path to validator.yaml config file",
			Value: configs.MainnetConfig + "/" + configs.TekuValidatorYamlPath,
		},
		&cli.StringFlag{
			Name:  TekuDatadirFlag,
			Usage: "Path to teku datadir",
			Value: configs.MainnetDatadir,
		},
	}
	Nimbus2StartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  Nimbus2ConfigFileFlag,
			Usage: "Path to nimbus.toml config file",
			Value: configs.MainnetConfig + "/" + configs.Nimbus2TomlPath,
		},
		&cli.StringFlag{
			Name:  Nimbus2ValidatorConfigFileFlag,
			Usage: "Path to validator.toml config file",
			Value: configs.MainnetConfig + "/" + configs.Nimbus2ValidatorTomlPath,
		},
		&cli.StringFlag{
			Name:  Nimbus2DatadirFlag,
			Usage: "Path to nimbus2 datadir",
			Value: configs.MainnetDatadir,
		},
		&cli.StringFlag{
			Name:  Nimbus2NetworkFlag,
			Usage: "Path to nimbus2 config directory, useful when using the --network nimbus flag",
			Value: configs.MainnetConfig + "/" + configs.Nimbus2Path,
		},
	}
	ValidatorStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   ValidatorDatadirFlag,
			Usage:  "Validator datadir",
			Value:  configs.ValidatorMainnetDatadir,
			Hidden: true,
		},
		&cli.StringFlag{
			Name:  ValidatorKeysFlag,
			Usage: "Location of generated wallet",
			Value: configs.MainnetKeystore,
		},
		&cli.StringFlag{
			Name:  ValidatorWalletPasswordFileFlag,
			Usage: "Location of the password file that you used to generate keys",
			Value: "",
		},
		&cli.StringFlag{
			Name:  ValidatorConfigFileFlag,
			Usage: "Path to validator.yaml config file",
			Value: configs.MainnetConfig + "/" + configs.ValidatorYamlPath,
		},
		&cli.StringFlag{
			Name:   ValidatorChainConfigFileFlag,
			Usage:  "Prysm chain config file path",
			Value:  configs.ChainConfigYamlPath,
			Hidden: true,
		},
	}
)

func init() {
	// Execution
	StartFlags = append(StartFlags, GethStartFlags...)
	StartFlags = append(StartFlags, ErigonStartFlags...)
	StartFlags = append(StartFlags, NethermindStartFlags...)
	StartFlags = append(StartFlags, BesuStartFlags...)
	// Consensus
	StartFlags = append(StartFlags, PrysmStartFlags...)
	StartFlags = append(StartFlags, LighthouseStartFlags...)
	StartFlags = append(StartFlags, TekuStartFlags...)
	StartFlags = append(StartFlags, Nimbus2StartFlags...)
	// Validator
	StartFlags = append(StartFlags, ValidatorStartFlags...)
	// Network
	StartFlags = append(StartFlags, NetworkFlags...)
}
