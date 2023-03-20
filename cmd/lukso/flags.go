package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

const (
	// geth related flag names
	gethTagFlag        = "geth-tag"
	gethCommitHashFlag = "geth-commit-hash"
	gethDatadirFlag    = "geth-datadir"
	gethConfigFileFlag = "geth-config"

	// Validator related flag names
	validatorTagFlag                = "validator-tag"
	validatorDatadirFlag            = "validator-datadir"
	validatorWalletPasswordFileFlag = "validator-password"
	validatorChainConfigFileFlag    = "validator-chain-config-file"

	// Prysm related flag names
	prysmTagFlag             = "prysm-tag"
	prysmGenesisStateFlag    = "genesis-ssz"
	prysmChainConfigFileFlag = "prysm-chain-config"
	prysmDatadirFlag         = "prysm-datadir"
	noSlasherFlag            = "no-slasher"

	// shared flags
	transactionFeeRecipientFlag = "transaction-fee-recipient"
	logFolderFlag               = "log-folder"
	jwtSecretFlag               = "jwt-secret"

	// non-specific flags
	mainnetFlag   = "mainnet"
	testnetFlag   = "testnet"
	devnetFlag    = "devnet"
	validatorFlag = "validator"
	consensusFlag = "consensus"
	executionFlag = "execution"

	acceptTermsOfUseFlag = "accept-terms-of-use"

	// shared values
	jwtSecretDefaultPath = "./config/mainnet/shared/secrets/jwt.hex"

	// flag defaults used in different contexts
	gethMainnetDatadir = "./mainnet-data/execution"
	gethTestnetDatadir = "./testnet-data/execution"
	gethDevnetDatadir  = "./devnet-data/execution"

	prysmMainnetDatadir = "./mainnet-data/consensus"
	prysmTestnetDatadir = "./testnet-data/consensus"
	prysmDevnetDatadir  = "./devnet-data/consensus"

	validatorMainnetDatadir = "./mainnet-data/validator"
	validatorTestnetDatadir = "./testnet-data/validator"
	validatorDevnetDatadir  = "./devnet-data/validator"

	mainnetLogs = "./mainnet-logs"
	testnetLogs = "./testnet-logs"
	devnetLogs  = "./devnet-logs"

	mainnetConfig = "./config/mainnet"
	testnetConfig = "./config/testnet"
	devnetConfig  = "./config/devnet"

	mainnetKeystore = "./mainnet-keystore"
	testnetKeystore = "./testnet-keystore"
	devnetKeystore  = "./devnet-keystore"

	// structure inside /config/selected-network directory.
	// we will select directory based on provided flag, by concatenating config path + file path
	genesisStateFilePath = "shared/genesis.ssz"
	configYamlPath       = "shared/config.yml"
	jwtSecretPath        = "shared/secrets/jwt.hex"
	configTomlPath       = "geth/geth.toml"
	genesisJsonPath      = "geth/genesis.json"

	// validator tool related flags
	depositDataJson    = "deposit-data-json"
	genesisDepositFlag = "genesis"
	validatorKeysFlag  = "validator-keys"
	gasPriceFlag       = "gas-price"
	rpcFlag            = "rpc"
	maxTxsPerBlock     = "max-txs-per-block"
)

var (
	jwtSelectedPath = jwtSecretDefaultPath

	mainnetEnabledFlag = &cli.BoolFlag{
		Name:  mainnetFlag,
		Usage: "Run for mainnetFlag (default)",
		Value: false,
	}
	testnetEnabledFlag = &cli.BoolFlag{
		Name:  testnetFlag,
		Usage: "Run for testnetFlag",
		Value: false,
	}
	devnetEnabledFlag = &cli.BoolFlag{
		Name:  devnetFlag,
		Usage: "Run for devnet",
		Value: false,
	}

	consensusSelectedFlag = &cli.BoolFlag{
		Name:  consensusFlag,
		Usage: "Run for consensus",
		Value: false,
	}
	executionSelectedFlag = &cli.BoolFlag{
		Name:  executionFlag,
		Usage: "Run for execution",
		Value: false,
	}
	validatorSelectedFlag = &cli.BoolFlag{
		Name:  validatorFlag,
		Usage: "Run for validator",
		Value: false,
	}

	networkFlags = []cli.Flag{
		mainnetEnabledFlag,
		testnetEnabledFlag,
		devnetEnabledFlag,
	}

	validatorDepositFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     depositDataJson,
			Usage:    "Path to your deposit file",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  genesisDepositFlag,
			Usage: "Specify if deposit is made for genesis",
			Value: false,
		},
		&cli.StringFlag{
			Name:  rpcFlag,
			Usage: "Your RPC provider",
			Value: "https://rpc.2022.l16.lukso.network",
		},
		&cli.IntFlag{
			Name:  gasPriceFlag,
			Usage: "Gas price provided by user",
			Value: 1000000000,
		},
		&cli.IntFlag{
			Name:  maxTxsPerBlock,
			Usage: "Maximum amount of txs sent per single block",
			Value: 10,
		},
	}

	validatorInitFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     validatorKeysFlag,
			Usage:    "Path to your validator keys and generated wallet",
			Required: true,
		},
		&cli.StringFlag{
			Name:  validatorWalletPasswordFileFlag,
			Usage: "Path to your password file",
			Value: "",
		},
	}

	downloadFlags []cli.Flag
	updateFlags   []cli.Flag
	stopFlags     = []cli.Flag{
		executionSelectedFlag,
		consensusSelectedFlag,
		validatorSelectedFlag,
	}
	startFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  validatorFlag,
			Usage: "Run lukso node with validator",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  noSlasherFlag,
			Usage: "disable slasher",
			Value: false, // default is true, we change it to false only when running validator
		},
		&cli.StringFlag{
			Name:     transactionFeeRecipientFlag,
			Usage:    "address to receive fees from blocks",
			Required: true,
		},
		&cli.StringFlag{
			Name:  logFolderFlag,
			Usage: "Directory to output logs into",
			Value: "./mainnet-logs",
		},
		&cli.StringFlag{
			Name:  jwtSecretFlag,
			Usage: "Path to jwt secret used for clients communication",
			Value: jwtSecretDefaultPath,
		},
	}
	logsFlags  []cli.Flag
	resetFlags []cli.Flag
	appFlags   = []cli.Flag{
		&cli.BoolFlag{
			Name:  acceptTermsOfUseFlag,
			Usage: "Accept terms of use. Default: false",
			Value: false,
		},
	}

	// GETH FLAGS
	// DOWNLOAD
	gethDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethTagFlag,
			Usage: "a tag of geth you would like to run",
			Value: "1.11.4",
		},
		&cli.StringFlag{
			Name:  gethCommitHashFlag,
			Usage: "a hash of commit that is bound to given release tag",
			Value: "7e3b149b",
		},
	}
	// UPDATE
	gethUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethTagFlag,
			Usage: "a tag of geth you would like to run",
			Value: "1.11.4",
		},
	}
	// START
	gethStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethDatadirFlag,
			Usage: "a path you would like to store your data",
			Value: gethMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  gethConfigFileFlag,
			Usage: "path to geth.toml config file",
			Value: "./config/mainnet/geth/geth.toml",
		},
	}
	// LOGS
	gethLogsFlags = []cli.Flag{}
	// RESET
	gethResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethDatadirFlag,
			Usage: "geth datadir",
			Value: gethMainnetDatadir,
		},
	}

	// PRYSM FLAGS
	// DOWNLOAD
	prysmDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmTagFlag,
			Usage: "tag for prysm",
			Value: "v3.2.2",
		},
	}
	// UPDATE
	prysmUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmTagFlag,
			Usage: "tag for prysm",
			Value: "v3.2.2",
		},
	}
	// START
	prysmStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmGenesisStateFlag,
			Usage: "genesis.ssz file path",
			Value: "./config/mainnet/shared/genesis.ssz",
		},
		&cli.StringFlag{
			Name:  prysmDatadirFlag,
			Usage: "prysm datadir",
			Value: prysmMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  prysmChainConfigFileFlag,
			Usage: "path to config.yml file",
			Value: "./config/mainnet/shared/config.yml",
		},
	}
	// LOGS
	prysmLogsFlags = []cli.Flag{}
	// RESET
	prysmResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmDatadirFlag,
			Usage: "prysm datadir",
			Value: prysmMainnetDatadir,
		},
	}

	// VALIDATOR
	// DOWNLOAD
	validatorDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "tag for validator binary",
			Value: "v3.2.2",
		},
	}
	// UPDATE
	validatorUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "tag for validator binary",
			Value: "v3.2.2",
		},
	}
	// START
	validatorStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorDatadirFlag,
			Usage: "validator datadir",
			Value: validatorMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  validatorKeysFlag,
			Usage: "location of generated wallet",
			Value: mainnetKeystore,
		},
		&cli.StringFlag{
			Name:  validatorWalletPasswordFileFlag,
			Usage: "location of file password that you used for generation keys from deposit-cli",
			Value: "./config/mainnet/shared/secrets/validator-password.txt",
		},
		&cli.StringFlag{
			Name:  validatorChainConfigFileFlag,
			Usage: "prysm chain config file path",
			Value: "./config/mainnet/shared/config.yml",
		},
	}
	// LOGS
	validatorLogsFlags = []cli.Flag{}
	// RESET
	validatorResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorDatadirFlag,
			Usage: "validator datadir",
			Value: validatorMainnetDatadir,
		},
	}
)

func (dependency *ClientDependency) ParseStartFlags(ctx *cli.Context) (startFlags []string) {
	name := dependency.name
	args := ctx.Args()
	argsLen := args.Len()
	for i := 0; i < argsLen; i++ {
		arg := args.Get(i)

		if strings.HasPrefix(arg, fmt.Sprintf("--%s", name)) {
			if i+1 == argsLen {
				startFlags = append(startFlags, arg)

				return
			}

			// we found a flag for our client - now we need to check if it's a value or bool flag
			nextArg := args.Get(i + 1)
			if strings.HasPrefix(nextArg, "--") { // we found a next flag, so current one is a bool
				startFlags = append(startFlags, arg)
			} else {
				startFlags = append(
					startFlags,
					fmt.Sprintf("--%s=%s", strings.TrimLeft(arg, fmt.Sprintf("--%s-", name)), nextArg),
				)
			}
		}
	}

	return
}

func prepareGethStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = clientDependencies[gethDependencyName].ParseStartFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(gethConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--authrpc.jwtsecret=%s", ctx.String(jwtSecretFlag)))

	return
}

func prepareValidatorStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = clientDependencies[validatorDependencyName].ParseStartFlags(ctx)

	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, fmt.Sprintf("--log-file=%s", prepareLogfileFlag(ctx, ctx.String(logFolderFlag), validatorDependencyName)))

	return
}

func preparePrysmStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = clientDependencies[prysmDependencyName].ParseStartFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--log-file=%s", prepareLogfileFlag(ctx, ctx.String(logFolderFlag), prysmDependencyName)))
	startFlags = append(startFlags, fmt.Sprintf("--jwt-secret=%s", ctx.String(jwtSecretFlag)))

	return
}
