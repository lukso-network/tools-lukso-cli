package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	// geth related flag names
	gethTagFlag            = "geth-tag"
	gethCommitHashFlag     = "geth-commit-hash"
	gethDatadirFlag        = "geth-datadir"
	gethEthstatsFlag       = "geth-ethstats"
	gethBootnodesFlag      = "geth-bootnodes"
	gethNetworkIDFlag      = "geth-networkid"
	gethPortFlag           = "geth-port"
	gethWSFlag             = "geth-ws"
	gethWSApiFlag          = "geth-ws-apis"
	gethGenesisFileFlag    = "geth-genesis"
	gethVerbosityFlag      = "geth-verbosity"
	gethHttpFlag           = "geth-http"
	gethHttpApiFlag        = "geth-http-apis"
	gethHttpAddrFlag       = "geth-http-addr"
	gethHttpPortFlag       = "geth-http-port"
	gethMineFlag           = "geth-mine"
	gethMinerGaslimitFlag  = "geth-miner-gaslimit"
	gethMinerEtherbaseFlag = "geth-miner-etherbase"
	gethMinerThreadsFlag   = "geth-miner-threads"
	gethStdOutputFlag      = "geth-std-output"
	gethAuthJWTSecretFlag  = "geth-auth-jwt-secret"
	gethNatFlag            = "geth-nat"
	gethOutputDirFlag      = "geth-output-dir"

	// Validator related flag names
	validatorTagFlag                   = "validator-tag"
	validatorDatadirFlag               = "validator-datadir"
	validatorWalletDirFlag             = "validator-wallet-dir"
	validatorWalletPasswordFileFlag    = "validator-wallet-password-file"
	validatorChainConfigFileFlag       = "validator-chain-config-file"
	validatorMonitoringHostFlag        = "validator-monitoring-host"
	validatorGrpcGatewayHostFlag       = "validator-grpc-gateway-host"
	validatorRpcHostFlag               = "validator-rpc-host"
	validatorSuggestedFeeRecipientFlag = "validator-suggested-fee-recipient"
	validatorVerbosityFlag             = "validator-verbosity"
	validatorOutputDirFlag             = "validator-output-dir"
	validatorStdOutputFlag             = "validator-std-output"

	// Prysm related flag names
	prysmTagFlag                     = "prysm-tag"
	prysmGenesisStateFlag            = "prysm-genesis-state"
	prysmDatadirFlag                 = "prysm-datadir"
	prysmExecutionEndpointFlag       = "prysm-execution-endpoint"
	prysmJWTSecretFlag               = "prysm-jwt-secret"
	prysmSuggestedFeeRecipientFlag   = "prysm-suggested-fee-recipient"
	prysmMinSyncPeersFlag            = "prysm-min-sync-peers"
	prysmContractDeploymentBlockFlag = "prysm-deposit-deployment"
	prysmP2pHostFlag                 = "prysm-p2p-host"
	prysmP2pmaxPeersFlag             = "prysm-p2p-max-peers"
	prysmChainConfigFileFlag         = "prysm-chain-config-file"
	prysmMonitoringHostFlag          = "prysm-monitoring-host"
	prysmGrpcGatewayHostFlag         = "prysm-grpc-gateway-host"
	prysmRpcHostFlag                 = "prysm-rpc-host"
	prysmSubscribeAllSubnetsFlag     = "prysm-subscribe-all-subnets"
	prysmMinimumPeersPerSubnetFlag   = "prysm-minimum-peers-per-subnet"
	prysmVerbosityFlag               = "prysm-verbosity"
	prysmOutputDirFlag               = "prysm-output-dir"
	prysmStdOutputFlag               = "prysm-std-output"
	prysmTestnetFlag                 = "prysm-testnet"

	acceptTermsOfUseFlagName = "accept-terms-of-use"
)

var (
	downloadFlags []cli.Flag
	updateFlags   []cli.Flag
	startFlags    []cli.Flag
	logsFlags     []cli.Flag
	appFlags      = []cli.Flag{
		&cli.BoolFlag{
			Name:  acceptTermsOfUseFlagName,
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
			Value: "1.10.26",
		},
		&cli.StringFlag{
			Name:  gethCommitHashFlag,
			Usage: "a hash of commit that is bound to given release tag",
			Value: "e5eb32ac",
		},
	}
	// UPDATE
	gethUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethTagFlag,
			Usage: "a tag of geth you would like to run",
			Value: "1.10.26",
		},
	}
	// START
	gethStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethDatadirFlag,
			Usage: "a path you would like to store your data",
			Value: "./geth",
		},
		&cli.BoolFlag{
			Name:  gethWSFlag,
			Usage: "enable WS server",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethWSApiFlag,
			Usage: "comma separated apis",
			Value: "eth,net",
		},
		&cli.StringFlag{
			Name:  gethNatFlag,
			Usage: "this flag sets up http nat to assign static ip for geth, default not set. Example extip:172.16.254.4",
			Value: "extip:172.16.254.4",
		},
		&cli.BoolFlag{
			Name:  gethHttpFlag,
			Usage: "enable HTTP server",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethHttpApiFlag,
			Usage: "comma separated apis",
			Value: "eth,net",
		},
		&cli.StringFlag{
			Name:  gethHttpAddrFlag,
			Usage: "address used in geth http communication",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  gethVerbosityFlag,
			Usage: "this flag sets up verbosity for geth",
			Value: "3",
		},
		&cli.StringFlag{
			Name:  gethPortFlag,
			Usage: "port for geth",
			Value: "30405",
		},
		&cli.StringFlag{
			Name:  gethHttpPortFlag,
			Usage: "port used in geth http communication",
			Value: "8565",
		},
		&cli.BoolFlag{
			Name:  gethMineFlag,
			Usage: "enable mining",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethMinerThreadsFlag,
			Usage: "number opf CPU threads used for mining",
			Value: "1",
		},
		&cli.StringFlag{
			Name:  gethMinerGaslimitFlag,
			Usage: "gas ceiling",
			Value: "60000000",
		},
		&cli.StringFlag{
			Name:  gethMinerEtherbaseFlag,
			Usage: "your ECDSA public key used to get rewards on geth chain",
			// yes, If you won't set it up, I'll get rewards ;]
			Value: "0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b",
		},
		&cli.StringFlag{
			Name:  gethAuthJWTSecretFlag,
			Usage: "path to a JWT secret used for secured endpoints authorization",
			Value: "./config/mainnet/secrets/jwt.hex",
		},
		&cli.BoolFlag{
			Name:  gethStdOutputFlag,
			Usage: "set geth output to stdout",
			Value: false,
		},
		&cli.StringFlag{
			Name:  gethOutputDirFlag,
			Usage: "Directory to output logs into",
			Value: "./logs/geth",
		},
	}
	// LOGS
	gethLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethOutputDirFlag,
			Usage: "file to output logs into",
			Value: "./logs/geth",
		},
	}

	// PRYSM FLAGS
	// DOWNLOAD
	prysmDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmTagFlag,
			Usage: "tag for prysm",
			Value: "v3.1.2",
		},
	}
	// UPDATE
	prysmUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmTagFlag,
			Usage: "tag for prysm",
			Value: "v3.1.2",
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
			Value: "./consensus_data",
		},
		&cli.StringFlag{
			Name:  prysmExecutionEndpointFlag,
			Usage: "execution endpoint",
			Value: "http://localhost:8551",
		},
		&cli.StringFlag{
			Name:  prysmJWTSecretFlag,
			Usage: "path to your jwt secret",
			Value: "./config/mainnet/secrets/jwt.hex",
		},
		&cli.StringFlag{
			Name:  prysmSuggestedFeeRecipientFlag,
			Usage: "address that receives block fees",
			Value: "0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b",
		},
		&cli.StringFlag{
			Name:  prysmMinSyncPeersFlag,
			Usage: "min sync peers for prysm, default 0",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmContractDeploymentBlockFlag,
			Usage: "deployment height of deposit contract, default 0.",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmP2pHostFlag,
			Usage: "p2p host for prysm, default empty",
			Value: "",
		},
		&cli.StringFlag{
			Name:  prysmChainConfigFileFlag,
			Usage: "path to config.yaml file",
			Value: "./config/mainnet/shared/config.yaml",
		},
		&cli.StringFlag{
			Name:  prysmMonitoringHostFlag,
			Usage: "host used for interacting with prometheus metrics",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  prysmGrpcGatewayHostFlag,
			Usage: "host for grpc gateway",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  prysmRpcHostFlag,
			Usage: "rpc server host",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  prysmVerbosityFlag,
			Usage: "verbosity for prysm",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  prysmP2pmaxPeersFlag,
			Usage: "max peers for prysm, default 250",
			Value: "250",
		},
		&cli.BoolFlag{
			Name:  prysmSubscribeAllSubnetsFlag,
			Usage: "subscribe to all possible subnets",
			Value: false,
		},
		&cli.StringFlag{
			Name:  prysmMinimumPeersPerSubnetFlag,
			Usage: "minimum peers per subnet",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmOutputDirFlag,
			Usage: "output destination folder of prysm logs",
			Value: "./logs/beacon_chain",
		},
		&cli.BoolFlag{
			Name:  prysmStdOutputFlag,
			Usage: "set prysm output to stdout",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  prysmTestnetFlag,
			Usage: "testnet",
			Value: false,
		},
	}
	// LOGS
	prysmLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmOutputDirFlag,
			Usage: "file to output logs into",
			Value: "./logs/beacon_chain",
		},
	}

	// VALIDATOR
	// DOWNLOAD
	validatorDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "tag for validator binary",
			Value: "v3.1.2",
		},
	}
	// UPDATE
	validatorUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "tag for validator binary",
			Value: "v3.1.2",
		},
	}
	// START
	validatorStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorDatadirFlag,
			Usage: "validator datadir",
			Value: "./validator_data",
		},
		&cli.StringFlag{
			Name:  validatorVerbosityFlag,
			Usage: "verbosity of validator",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  validatorWalletDirFlag,
			Usage: "location of generated wallet",
			Value: "./mainnet_keystore",
		},
		&cli.StringFlag{
			Name:  validatorWalletPasswordFileFlag,
			Usage: "location of file password that you used for generation keys from deposit-cli",
			Value: "./config/mainnet/secrets/validator-password.txt",
		},
		&cli.StringFlag{
			Name:  validatorChainConfigFileFlag,
			Usage: "prysm chain config file path",
			Value: "./config/mainnet/shared/config.yaml",
		},
		&cli.StringFlag{
			Name:  validatorMonitoringHostFlag,
			Usage: "host used for interacting with prometheus metrics",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  validatorGrpcGatewayHostFlag,
			Usage: "host for grpc gateway",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  validatorRpcHostFlag,
			Usage: "rpc server host",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  validatorSuggestedFeeRecipientFlag,
			Usage: "address that receives block fees",
			Value: "0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b",
		},
		&cli.StringFlag{
			Name:  validatorOutputDirFlag,
			Usage: "output destination folder of validator logs",
			Value: "./logs/validator",
		},
		&cli.BoolFlag{
			Name:  validatorStdOutputFlag,
			Usage: "set validator output to stdout",
			Value: false,
		},
	}
	// LOGS
	validatorLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorOutputDirFlag,
			Usage: "file to output logs into",
			Value: "./logs/validator",
		},
	}
)

func prepareGethStartFlags(ctx *cli.Context) (startFlags []string) {
	// parse all runtime-related geth flags, one by one and append them
	startFlags = append(startFlags, fmt.Sprintf("--datadir=%s", ctx.String(gethDatadirFlag)))
	if ctx.Bool(gethWSFlag) {
		startFlags = append(startFlags, "--ws")
	}
	startFlags = append(startFlags, fmt.Sprintf("--ws.api=%s", ctx.String(gethWSApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--nat=%s", ctx.String(gethNatFlag)))
	if ctx.Bool(gethHttpFlag) {
		startFlags = append(startFlags, "--http")
	}
	startFlags = append(startFlags, fmt.Sprintf("--http.api=%s", ctx.String(gethHttpApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.addr=%s", ctx.String(gethHttpAddrFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.port=%s", ctx.String(gethHttpPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(gethVerbosityFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--port=%s", ctx.String(gethPortFlag)))
	if ctx.Bool(gethMineFlag) {
		startFlags = append(startFlags, "--mine")
	}
	startFlags = append(startFlags, fmt.Sprintf("--miner.threads=%s", ctx.String(gethMinerThreadsFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.gaslimit=%s", ctx.String(gethMinerGaslimitFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.etherbase=%s", ctx.String(gethMinerEtherbaseFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--authrpc.jwtsecret=%s", ctx.String(gethAuthJWTSecretFlag)))

	return
}

func prepareValidatorStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, fmt.Sprintf("--datadir=%s", ctx.String(validatorDatadirFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--wallet-dir=%s", ctx.String(validatorWalletDirFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--wallet-password-file=%s", ctx.String(validatorWalletPasswordFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--chain-config-file=%s", ctx.String(validatorChainConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--monitoring-host=%s", ctx.String(validatorMonitoringHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--grpc-gateway-host=%s", ctx.String(validatorGrpcGatewayHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--rpc-host=%s", ctx.String(validatorRpcHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(validatorSuggestedFeeRecipientFlag)))

	logFileFlag := prepareLogfileFlag(ctx, validatorOutputDirFlag, validatorDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	return
}

func preparePrysmStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, "--force-clear-db")
	startFlags = append(startFlags, fmt.Sprintf("--genesis-state=%s", ctx.String(prysmGenesisStateFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--datadir=%s", ctx.String(prysmDatadirFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--execution-endpoint=%s", ctx.String(prysmExecutionEndpointFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--jwt-secret=%s", ctx.String(prysmJWTSecretFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--min-sync-peers=%s", ctx.String(prysmMinSyncPeersFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--contract-deployment-block=%s", ctx.String(prysmContractDeploymentBlockFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--p2p-host-ip=%s", ctx.String(prysmP2pHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--chain-config-file=%s", ctx.String(prysmChainConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--monitoring-host=%s", ctx.String(prysmMonitoringHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--grpc-gateway-host=%s", ctx.String(prysmGrpcGatewayHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(prysmSuggestedFeeRecipientFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--rpc-host=%s", ctx.String(prysmRpcHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(prysmVerbosityFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--p2p-max-peers=%s", ctx.String(prysmP2pmaxPeersFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--subscribe-all-subnets=%s", ctx.String(prysmMinimumPeersPerSubnetFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--minimum-peers-per-subnet=%s", ctx.String(prysmMinimumPeersPerSubnetFlag)))

	logFileFlag := prepareLogfileFlag(ctx, prysmOutputDirFlag, prysmDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	if ctx.Bool(prysmTestnetFlag) {
		startFlags = append(startFlags, "--goerli")
	}

	return
}
