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
	gethWSAddrFlag         = "geth-ws-addr"
	gethWSOriginsFlag      = "geth-ws-origins"
	gethVerbosityFlag      = "geth-verbosity"
	gethHttpFlag           = "geth-http"
	gethHttpApiFlag        = "geth-http-apis"
	gethHttpAddrFlag       = "geth-http-addr"
	gethHttpPortFlag       = "geth-http-port"
	gethHttpCorsDomainFlag = "geth-http-corsdomain"
	gethHttpVHostsFlag     = "geth-http-vhosts"
	gethIPCDisableFlag     = "geth-ipcdisable"
	gethMetricsFlag        = "geth-metrics"
	gethMetricsAddrFlag    = "geth-metrics-addr"
	gethSyncmodeFlag       = "geth-syncmode"
	gethGcmodeFlag         = "geth-gcmode"
	gethMineFlag           = "geth-mine"
	gethMinerGaslimitFlag  = "geth-miner-gaslimit"
	gethMinerEtherbaseFlag = "geth-miner-etherbase"
	gethMinerThreadsFlag   = "geth-miner-threads"
	gethStdOutputFlag      = "geth-std-output"
	gethAuthJWTSecretFlag  = "geth-auth-jwt-secret"
	gethNatFlag            = "geth-nat"
	gethTxLookupLimitFlag  = "geth-tx-lookup-limit"
	gethCachePreimagesFlag = "geth-cache-preimages"
	gethLogDirFlag         = "geth-log-dir"
	gethOutputFileFlag     = "geth-output-file"

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
	validatorLogDirFlag                = "validator-log-dir"
	validatorOutputFileFlag            = "validator-output-file"

	validatorStdOutputFlag = "validator-std-output"

	// Prysm related flag names
	prysmTagFlag                     = "prysm-tag"
	prysmGenesisStateFlag            = "prysm-genesis-state"
	prysmDatadirFlag                 = "prysm-datadir"
	prysmBootstrapNodesFlag          = "prysm-bootstrap-nodes"
	prysmEnableRpcDebugEndpointsFlag = "prysm-enable-rpc-debug-endpoints"
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
	prysmLogDirFlag                  = "prysm-log-dir"
	prysmStdOutputFlag               = "prysm-std-output"
	prysmOutputFileFlag              = "prysm-output-file"

	// non-specific flags
	validatorEnabledFlag = "validator"
	mainnetEnabledFlag   = "mainnet"
	testnetEnabledFlag   = "testnet"
	devnetEnabledFlag    = "devnet"

	acceptTermsOfUseFlagName = "accept-terms-of-use"

	// shared values
	jwtSecretDefaultPath = "./config/mainnet/shared/secrets/jwt.hex"

	// bootnodes
	gethBootstrapNode  = "enode://c2bb19ce658cfdf1fecb45da599ee6c7bf36e5292efb3fb61303a0b2cd07f96c20ac9b376a464d687ac456675a2e4a44aec39a0509bcb4b6d8221eedec25aca2@34.91.82.99:30303"
	prysmBootstrapNode = "enr:-MK4QOtYSPGAg5FCQRxy8_kAyrq1lSkvkqA4FXPc-myHYCdmW-U0mu_m1oFR-YL-tDbhecFo05WerA1IbFk4tBHVgC6GAYXMBqQXh2F0dG5ldHOIAAAAAAAAAACEZXRoMpDXjD-DICIABP__________gmlkgnY0gmlwhCJbUmOJc2VjcDI1NmsxoQLt3oS_p6rhGF3E8aS3UZLcMboK93av0NkFVAwwsbmoc4hzeW5jbmV0cwCDdGNwgjLIg3VkcIIu4A"

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
	configTomlPath       = "geth/config.toml"
	genesisJsonPath      = "geth/genesis.json"
)

var (
	jwtPath = jwtSecretDefaultPath

	mainnetFlag = &cli.BoolFlag{
		Name:  mainnetEnabledFlag,
		Usage: "Run for mainnet (default)",
		Value: false,
	}
	testnetFlag = &cli.BoolFlag{
		Name:  testnetEnabledFlag,
		Usage: "Run for testnet",
		Value: false,
	}
	devnetFlag = &cli.BoolFlag{
		Name:  devnetEnabledFlag,
		Usage: "Run for devnet",
		Value: false,
	}

	networkFlags = []cli.Flag{
		mainnetFlag,
		testnetFlag,
		devnetFlag,
	}

	downloadFlags []cli.Flag
	updateFlags   []cli.Flag
	startFlags    = []cli.Flag{
		&cli.BoolFlag{
			Name:  validatorEnabledFlag,
			Usage: "Run lukso node with validator",
			Value: false,
		},
	}
	logsFlags  []cli.Flag
	resetFlags []cli.Flag
	appFlags   = []cli.Flag{
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
			Value: gethMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  gethEthstatsFlag,
			Usage: "URL of ethstats service",
			Value: "local-pc-test@34.90.215.182",
		},
		&cli.StringFlag{
			Name:  gethBootnodesFlag,
			Usage: "Bootnodes for geth",
			Value: gethBootstrapNode,
		},
		&cli.StringFlag{
			Name:  gethNetworkIDFlag,
			Usage: "Network ID",
			Value: "2022",
		},
		&cli.BoolFlag{
			Name:  gethWSFlag,
			Usage: "enable WS server",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethWSApiFlag,
			Usage: "comma separated apis",
			Value: "net,eth,debug,engine",
		},
		&cli.StringFlag{
			Name:  gethWSAddrFlag,
			Usage: "WS address",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  gethWSOriginsFlag,
			Usage: "Origins from which to accept WS requests",
			Value: "*",
		},
		&cli.StringFlag{
			Name:  gethNatFlag,
			Usage: "this flag sets up http nat to assign static ip for geth, default not set. Example extip:172.16.254.4",
			Value: "extip:83.144.95.19",
		},
		&cli.BoolFlag{
			Name:  gethHttpFlag,
			Usage: "enable HTTP server",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethHttpApiFlag,
			Usage: "comma separated apis",
			Value: "net,eth,debug,engine,txlookup",
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
			Value: "30303",
		},
		&cli.StringFlag{
			Name:  gethHttpPortFlag,
			Usage: "port used in geth http communication",
			Value: "8545",
		},
		&cli.StringFlag{
			Name:  gethHttpCorsDomainFlag,
			Usage: "accepted CORS domains",
			Value: "*",
		},
		&cli.StringFlag{
			Name:  gethHttpVHostsFlag,
			Usage: "comma separated virtual hostnames",
			Value: "*",
		},
		&cli.BoolFlag{
			Name:  gethIPCDisableFlag,
			Usage: "disable IPC communication",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  gethMetricsFlag,
			Usage: "enable metrics reporting",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethMetricsAddrFlag,
			Usage: "address of metrics service",
			Value: "0.0.0.0",
		},
		&cli.StringFlag{
			Name:  gethSyncmodeFlag,
			Usage: "geth's sync mode",
			Value: "snap",
		},
		&cli.StringFlag{
			Name:  gethGcmodeFlag,
			Usage: "garbage collection mode",
			Value: "full",
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
			Value: "1600000000",
		},
		&cli.StringFlag{
			Name:  gethMinerEtherbaseFlag,
			Usage: "your ECDSA public key used to get rewards on geth chain",
			// yes, If you won't set it up, I'll get rewards ;]
			Value: "0x0000000000000000000000000000000000000000",
		},
		&cli.StringFlag{
			Name:  gethTxLookupLimitFlag,
			Usage: "number of blocks to maintain tx indexes from",
			Value: "1",
		},
		&cli.BoolFlag{
			Name:  gethCachePreimagesFlag,
			Usage: "enable preimage caching",
			Value: true,
		},
		&cli.StringFlag{
			Name:  gethAuthJWTSecretFlag,
			Usage: "path to a JWT secret used for secured endpoints authorization",
			Value: jwtSecretDefaultPath,
		},
		&cli.BoolFlag{
			Name:  gethStdOutputFlag,
			Usage: "set geth output to stdout",
			Value: false,
		},
		&cli.StringFlag{
			Name:  gethLogDirFlag,
			Usage: "Directory to output logs into",
			Value: "./mainnet-logs",
		},
	}
	// LOGS
	gethLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethLogDirFlag,
			Usage: "path file to log from",
			Value: "./mainnet-logs",
		},
	}
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
			Value: prysmMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  prysmBootstrapNodesFlag,
			Usage: "bootnodes for prysm beaconchain",
			Value: prysmBootstrapNode,
		},
		&cli.BoolFlag{
			Name:  prysmEnableRpcDebugEndpointsFlag,
			Usage: "enable debugging RPC endpoints",
			Value: true,
		},
		&cli.StringFlag{
			Name:  prysmExecutionEndpointFlag,
			Usage: "execution endpoint",
			Value: "http://localhost:8551",
		},
		&cli.StringFlag{
			Name:  prysmJWTSecretFlag,
			Usage: "path to your jwt secret",
			Value: jwtSecretDefaultPath,
		},
		&cli.StringFlag{
			Name:  prysmSuggestedFeeRecipientFlag,
			Usage: "address that receives block fees",
			Value: "0x0000000000000000000000000000000000000000",
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
			Usage: "path to config.yml file",
			Value: "./config/mainnet/shared/config.yml",
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
			Value: true,
		},
		&cli.StringFlag{
			Name:  prysmMinimumPeersPerSubnetFlag,
			Usage: "minimum peers per subnet",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmLogDirFlag,
			Usage: "output destination folder of prysm logs",
			Value: "./mainnet-logs",
		},
		&cli.BoolFlag{
			Name:  prysmStdOutputFlag,
			Usage: "set prysm output to stdout",
			Value: false,
		},
	}
	// LOGS
	prysmLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmLogDirFlag,
			Usage: "path to file to log from",
			Value: "./mainnet-logs",
		},
	}
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
			Value: validatorMainnetDatadir,
		},
		&cli.StringFlag{
			Name:  validatorVerbosityFlag,
			Usage: "verbosity of validator",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  validatorWalletDirFlag,
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
			Value: "0x0000000000000000000000000000000000000000",
		},
		&cli.StringFlag{
			Name:  validatorLogDirFlag,
			Usage: "output destination folder of validator logs",
			Value: "./mainnet-logs",
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
			Name:  validatorLogDirFlag,
			Usage: "path to file to log from",
			Value: "./mainnet-logs",
		},
	}
	// RESET
	validatorResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorDatadirFlag,
			Usage: "validator datadir",
			Value: validatorMainnetDatadir,
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
	startFlags = append(startFlags, fmt.Sprintf("--ws.addr=%s", ctx.String(gethWSAddrFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.origins=%s", ctx.String(gethWSOriginsFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--bootnodes=%s", ctx.String(gethBootnodesFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--networkid=%s", ctx.String(gethNetworkIDFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--nat=%s", ctx.String(gethNatFlag)))
	if ctx.Bool(gethHttpFlag) {
		startFlags = append(startFlags, "--http")
	}
	startFlags = append(startFlags, fmt.Sprintf("--http.api=%s", ctx.String(gethHttpApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.addr=%s", ctx.String(gethHttpAddrFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.port=%s", ctx.String(gethHttpPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.corsdomain=%s", ctx.String(gethHttpCorsDomainFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.vhosts=%s", ctx.String(gethHttpVHostsFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(gethVerbosityFlag)))
	if ctx.Bool(gethIPCDisableFlag) {
		startFlags = append(startFlags, "--ipcdisable")
	}
	if ctx.String(gethEthstatsFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--ethstats=%s", ctx.String(gethEthstatsFlag)))
	}
	if ctx.Bool(gethMetricsFlag) {
		startFlags = append(startFlags, "--metrics")
	}
	startFlags = append(startFlags, fmt.Sprintf("--metrics.addr=%s", ctx.String(gethMetricsAddrFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--syncmode=%s", ctx.String(gethSyncmodeFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--gcmode=%s", ctx.String(gethGcmodeFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--txlookuplimit=%s", ctx.String(gethTxLookupLimitFlag)))
	if ctx.Bool(gethMineFlag) {
		startFlags = append(startFlags, "--mine")
	}
	startFlags = append(startFlags, fmt.Sprintf("--port=%s", ctx.String(gethPortFlag)))
	if ctx.Bool(gethCachePreimagesFlag) {
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

	logFileFlag := prepareLogfileFlag(ctx, validatorLogDirFlag, validatorDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	return
}

func preparePrysmStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, "--force-clear-db")
	startFlags = append(startFlags, fmt.Sprintf("--bootstrap-node=%s", ctx.String(prysmBootstrapNodesFlag)))
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
	if ctx.Bool(prysmSubscribeAllSubnetsFlag) {
		startFlags = append(startFlags, "--subscribe-all-subnets")
	}
	startFlags = append(startFlags, fmt.Sprintf("--minimum-peers-per-subnet=%s", ctx.String(prysmMinimumPeersPerSubnetFlag)))
	if ctx.Bool(prysmEnableRpcDebugEndpointsFlag) {
		startFlags = append(startFlags, "--enable-debug-rpc-endpoints")
	}

	logFileFlag := prepareLogfileFlag(ctx, prysmLogDirFlag, prysmDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	return
}
