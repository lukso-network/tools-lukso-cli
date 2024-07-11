package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

func init() {
	InstallFlags = append(InstallFlags, GethDownloadFlags...)
	InstallFlags = append(InstallFlags, ValidatorDownloadFlags...)
	InstallFlags = append(InstallFlags, PrysmDownloadFlags...)
	InstallFlags = append(InstallFlags, ErigonDownloadFlags...)
	InstallFlags = append(InstallFlags, LighthouseDownloadFlags...)
	InstallFlags = append(InstallFlags, TekuDownloadFlags...)
	InstallFlags = append(InstallFlags, NethermindDownloadFlags...)
	InstallFlags = append(InstallFlags, BesuDownloadFlags...)

	StartFlags = append(StartFlags, GethStartFlags...)
	StartFlags = append(StartFlags, ErigonStartFlags...)
	StartFlags = append(StartFlags, PrysmStartFlags...)
	StartFlags = append(StartFlags, LighthouseStartFlags...)
	StartFlags = append(StartFlags, ValidatorStartFlags...)
	StartFlags = append(StartFlags, TekuStartFlags...)
	StartFlags = append(StartFlags, NetworkFlags...)
	StartFlags = append(StartFlags, NethermindStartFlags...)
	StartFlags = append(StartFlags, BesuStartFlags...)

	LogsFlags = append(LogsFlags, GethLogsFlags...)
	LogsFlags = append(LogsFlags, PrysmLogsFlags...)
	LogsFlags = append(LogsFlags, ValidatorLogsFlags...)
	LogsFlags = append(LogsFlags, NetworkFlags...)

	ResetFlags = append(ResetFlags, GethResetFlags...)
	ResetFlags = append(ResetFlags, PrysmResetFlags...)
	ResetFlags = append(ResetFlags, ValidatorResetFlags...)
	ResetFlags = append(ResetFlags, NetworkFlags...)

	// after we exported flags from each command we can update them
	GethStartFlags = append(GethStartFlags, NetworkFlags...)
	GethLogsFlags = append(GethLogsFlags, NetworkFlags...)
	GethResetFlags = append(GethResetFlags, NetworkFlags...)

	PrysmStartFlags = append(PrysmStartFlags, NetworkFlags...)
	PrysmLogsFlags = append(PrysmLogsFlags, NetworkFlags...)
	PrysmResetFlags = append(PrysmResetFlags, NetworkFlags...)

	ValidatorStartFlags = append(ValidatorStartFlags, NetworkFlags...)
	ValidatorLogsFlags = append(ValidatorLogsFlags, NetworkFlags...)
	ValidatorResetFlags = append(ValidatorResetFlags, NetworkFlags...)

	ValidatorImportFlags = append(ValidatorImportFlags, NetworkFlags...)
	ValidatorListFlags = append(ValidatorListFlags, NetworkFlags...)
	ValidatorExitFlags = append(ValidatorExitFlags, NetworkFlags...)
}

var (
	mainnetEnabledFlag = &cli.BoolFlag{
		Name:  MainnetFlag,
		Usage: "Run for mainnet",
		Value: true,
	}
	testnetEnabledFlag = &cli.BoolFlag{
		Name:  TestnetFlag,
		Usage: "Run for testnet",
		Value: false,
	}
	devnetEnabledFlag = &cli.BoolFlag{
		Name:  DevnetFlag,
		Usage: "Run for devnet",
		Value: false,
	}

	consensusSelectedFlag = &cli.BoolFlag{
		Name:  ConsensusFlag,
		Usage: "Run for consensus",
		Value: false,
	}
	executionSelectedFlag = &cli.BoolFlag{
		Name:  ExecutionFlag,
		Usage: "Run for execution",
		Value: false,
	}
	validatorSelectedFlag = &cli.BoolFlag{
		Name:  ValidatorFlag,
		Usage: "Run for validator",
		Value: false,
	}

	NetworkFlags = []cli.Flag{
		mainnetEnabledFlag,
		testnetEnabledFlag,
		devnetEnabledFlag,
	}

	ValidatorImportFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   ValidatorWalletDirFlag,
			Usage:  "Selected wallet",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   ValidatorDatadirFlag,
			Usage:  "Selected wallet", // wallet for lighthouse - incorporated into datadir
			Hidden: true,
		},
		&cli.StringFlag{
			Name:     ValidatorKeysFlag,
			Usage:    "Location of your validator keys",
			Required: true,
		},
		&cli.StringFlag{
			Name:  ValidatorPasswordFlag,
			Usage: "Location of your validator keys' password",
		},
	}

	ValidatorListFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ValidatorWalletDirFlag,
			Usage: "Path to wallet containing validators that you want to exit - this flag is required when using Prysm validator",
		},
	}
	ValidatorExitFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  KeystoreFlag,
			Usage: "Path to keystore containing public key that you want to exit - this flag is required when using Lighthouse validator",
		},
		&cli.StringFlag{
			Name:  ValidatorWalletDirFlag,
			Usage: "Path to wallet containing validators that you want to exit - this flag is required when using Prysm validator",
		},
		&cli.StringFlag{
			Name:  RpcAddressFlag,
			Usage: "Address of node that is used to make an exit (defaults to the default RPC address provided by your selected client)",
		},
		&cli.StringFlag{
			Name:  TestnetDirFlag,
			Usage: "Path to network configuration folder",
			Value: configs.MainnetConfig + "/shared",
		},
	}
	InstallFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  AgreeTermsFlag,
			Usage: "Automatically accept Terms and Conditions",
			Value: false,
		},
	}
	UpdateConfigFlags []cli.Flag
	StopFlags         = []cli.Flag{
		executionSelectedFlag,
		consensusSelectedFlag,
		validatorSelectedFlag,
	}

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
	LogsFlags        []cli.Flag
	ResetFlags       []cli.Flag
	AppFlags         []cli.Flag
	StatusPeersFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ExecutionClientHost,
			Usage: "Host for execution client",
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:  ConsensusClientHost,
			Usage: "Host for consensus client",
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:  ValidatorClientHost,
			Usage: "Host for validator client",
			Value: "localhost",
		},
		&cli.IntFlag{
			Name:  ExecutionClientPort,
			Usage: "Port for execution client (Defaults to: 8545: Geth, Erigon)",
		},
		&cli.IntFlag{
			Name:  ConsensusClientPort,
			Usage: "Port for consensus client (Defaults to: 3500: Prysm | 4000: Lighthouse)",
		},
		&cli.IntFlag{
			Name:  ValidatorClientPort,
			Usage: "Port for validator client (Defaults to: 8545: Geth, Erigon)",
		},
	}

	GethDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  GethTagFlag,
			Usage: "A tag of geth you would like to run",
			Value: common.GethTag,
		},
		&cli.StringFlag{
			Name:  GethCommitHashFlag,
			Usage: "A hash of commit that is bound to given release tag",
			Value: common.GethCommitHash,
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

	GethLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: "./mainnet-logs",
		},
	}

	GethResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   GethDatadirFlag,
			Usage:  "geth datadir",
			Value:  configs.ExecutionMainnetDatadir,
			Hidden: true,
		},
	}

	ErigonDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ErigonTagFlag,
			Usage: "Tag for erigon",
			Value: common.ErigonTag,
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

	NethermindDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  NethermindTagFlag,
			Usage: "Tag for nethermind",
			Value: common.NethermindTag,
		},
		&cli.StringFlag{
			Name:  NethermindCommitHashFlag,
			Usage: "A hash of commit that is bound to given release tag",
			Value: common.NethermindCommitHash,
		},
	}
	NethermindStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  NethermindConfigFileFlag,
			Usage: "Path to nethermind.cfg config file",
			Value: configs.MainnetConfig + "/" + configs.NethermindCfgPath,
		},
		&cli.StringFlag{
			Name:  NethermindDatadirFlag,
			Usage: "Nethermind datadir",
			Value: configs.ExecutionMainnetDatadir,
		},
	}

	BesuDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  BesuTagFlag,
			Usage: "Tag for besu",
			Value: common.BesuTag,
		},
	}
	BesuStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  BesuConfigFileFlag,
			Usage: "Path to besu.yaml config file",
			Value: configs.MainnetConfig + "/" + configs.BesuYamlPath,
		},
		&cli.StringFlag{
			Name:  BesuDatadirFlag,
			Usage: "Besu datadir",
			Value: configs.ExecutionMainnetDatadir,
		},
	}

	PrysmDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  PrysmTagFlag,
			Usage: "Tag for prysm",
			Value: common.PrysmTag,
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
	PrysmLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: "./mainnet-logs",
		},
	}
	PrysmResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   PrysmDatadirFlag,
			Usage:  "prysm datadir",
			Value:  configs.ConsensusMainnetDatadir,
			Hidden: true,
		},
	}

	LighthouseDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LighthouseTagFlag,
			Usage: "Tag for lighthouse",
			Value: common.LighthouseTag,
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

	ValidatorDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ValidatorTagFlag,
			Usage: "Tag for validator binary",
			Value: common.PrysmTag,
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
	ValidatorLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: configs.MainnetLogs,
		},
	}
	ValidatorResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   ValidatorDatadirFlag,
			Usage:  "Validator datadir",
			Value:  configs.ValidatorMainnetDatadir,
			Hidden: true,
		},
	}
	TekuDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  TekuTagFlag,
			Usage: "Tag for teku client",
			Value: common.TekuTag,
		},
	}
	TekuStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  TekuConfigFileFlag,
			Usage: "Path to teku.yaml config file",
			Value: configs.MainnetConfig + configs.TekuYamlPath,
		},
		&cli.StringFlag{
			Name:  TekuValidatorConfigFileFlag,
			Usage: "Path to validator.yaml config file",
			Value: configs.MainnetConfig + configs.TekuValidatorYamlPath,
		},
	}
)
