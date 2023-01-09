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
	gethMinerGaslimit      = "geth-miner-gaslimit"
	gethMinerEtherbaseFlag = "geth-miner-etherbase"
	gethMinerThreads       = "geth-miner-threads"
	gethStdOutputFlag      = "geth-std-output"
	gethAuthJWTSecretFlag  = "geth-auth-jwt-secret"
	gethNatFlag            = "geth-nat"
	gethOutputDirFlag      = "geth-output-dir"

	// Common for prysm client
	prysmChainConfigFlag = "prysm-chain-config"

	// Validator related flag names
	validatorTagFlag                = "validator-tag"
	validatorDatadirFlag            = "validator-datadir"
	validatorPrysmRpcProviderFlag   = "validator-prysm-rpc"
	validatorVerbosityFlag          = "validator-verbosity"
	validatorTrustedGethFlag        = "validator-trusted-geth"
	validatorWalletPasswordFileFlag = "validator-wallet-password-file"
	validatorOutputDirFlag          = "validator-output-dir"
	validatorStdOutputFlag          = "validator-std-output"

	// Prysm related flag names
	prysmTagFlag                     = "prysm-tag"
	prysmDatadirFlag                 = "prysm-datadir"
	prysmGenesisStateFlag            = "prysm-genesis-state"
	prysmBootnodesFlag               = "prysm-bootnode"
	prysmPeerFlag                    = "prysm-peer"
	prysmWeb3ProviderFlag            = "prysm-web3provider"
	prysmDepositContractFlag         = "prysm-deposit-contract"
	prysmContractDeploymentBlockFlag = "prysm-deposit-deployment"
	prysmVerbosityFlag               = "prysm-verbosity"
	prysmMinSyncPeersFlag            = "prysm-min-sync-peers"
	prysmMaxSyncPeersFlag            = "prysm-max-sync-peers"
	prysmP2pHostFlag                 = "prysm-p2p-host"
	prysmP2pLocalFlag                = "prysm-p2p-local"
	prysmDisableSyncFlag             = "prysm-disable-sync"
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
			Usage: "provide a tag of geth you would like to run",
			Value: "1.10.26",
		},
		&cli.StringFlag{
			Name:  gethCommitHashFlag,
			Usage: "provide a hash of commit that is bound to given release tag",
			Value: "e5eb32ac",
		},
	}
	// UPDATE
	gethUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethTagFlag,
			Usage: "provide a tag of geth you would like to run",
			Value: "1.10.26",
		},
	}
	// START
	gethStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  gethDatadirFlag,
			Usage: "provide a path you would like to store your data",
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
			Usage: "provide port for geth",
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
			Name:  gethMinerThreads,
			Usage: "number opf CPU threads used for mining",
			Value: "1",
		},
		&cli.StringFlag{
			Name:  gethMinerGaslimit,
			Usage: "gas ceiling",
			Value: "60000000",
		},
		&cli.StringFlag{
			Name:  gethMinerEtherbaseFlag,
			Usage: "your ECDSA public key used to get rewards on geth chain",
			// yes, If you won't set it up, I'll get rewards ;]
			Value: "0x59E3dADc83af3c127a2e29B12B0E86109Bb6d838",
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
			Usage: "provide tag for prysm",
			Value: "v3.1.2",
		},
		&cli.StringFlag{
			Name:  prysmDatadirFlag,
			Usage: "provide prysm datadir",
			Value: "./prysm",
		},
	}
	// UPDATE
	prysmUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmTagFlag,
			Usage: "provide tag for prysm",
			Value: "v3.1.2",
		},
	}
	// START
	prysmStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  prysmGenesisStateFlag,
			Usage: "provide genesis.ssz file",
			Value: "./config/mainnet/prysm/prysm-beaconchain-genesis.ssz",
		},
		&cli.StringFlag{
			Name:  prysmDatadirFlag,
			Usage: "provide prysm datadir",
			Value: "./prysm",
		},
		&cli.StringFlag{
			Name:  prysmBootnodesFlag,
			Usage: `provide coma separated bootnode enr, default: "enr:-Ku4QANldTRLCRUrY9K4OAmk_ATOAyS_sxdTAaGeSh54AuDJXxOYij1fbgh4KOjD4tb2g3T-oJmMjuJyzonLYW9OmRQBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhAoABweJc2VjcDI1NmsxoQKWfbT1atCho149MGMvpgBWUymiOv9QyXYhgYEBZvPBW4N1ZHCCD6A"`,
			Value: "",
		},
		&cli.StringFlag{
			Name:  prysmPeerFlag,
			Usage: `provide coma separated peer enr address, default: ""`,
			Value: "",
		},
		&cli.StringFlag{
			Name:  prysmWeb3ProviderFlag,
			Usage: "provide web3 provider (network of deposit contract deployment)",
			Value: "http://127.0.0.1:8565",
		},
		&cli.StringFlag{
			Name:  prysmDepositContractFlag,
			Usage: "provide deposit contract address",
			Value: "0x000000000000000000000000000000000000cafe",
		},
		&cli.StringFlag{
			Name:  prysmContractDeploymentBlockFlag,
			Usage: "provide deployment height of deposit contract, default 0.",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmVerbosityFlag,
			Usage: "provide verobosity for prysm",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  prysmMinSyncPeersFlag,
			Usage: "provide min sync peers for prysm, default 0",
			Value: "0",
		},
		&cli.StringFlag{
			Name:  prysmMaxSyncPeersFlag,
			Usage: "provide max sync peers for prysm, default 25",
			Value: "25",
		},
		&cli.StringFlag{
			Name:  prysmP2pHostFlag,
			Usage: "provide p2p host for prysm, default empty",
			Value: "",
		}, &cli.StringFlag{
			Name:  prysmP2pLocalFlag,
			Usage: "provide p2p local ip for prysm, default empty",
			Value: "",
		},
		&cli.BoolFlag{
			Name:  prysmDisableSyncFlag,
			Usage: "disable initial sync phase",
			Value: false,
		},
		&cli.StringFlag{
			Name:  prysmOutputDirFlag,
			Usage: "provide output destination of prysm",
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
			Usage: "provide tag for validator binary",
			Value: "v3.1.2",
		},
	}
	// UPDATE
	validatorUpdateFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "provide tag for validator binary",
			Value: "v3.1.2",
		},
	}
	// START
	validatorStartFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorPrysmRpcProviderFlag,
			Usage: "provide url without prefix, example: localhost:4000",
			Value: "localhost:4000",
		},
		&cli.StringFlag{
			Name:  validatorVerbosityFlag,
			Usage: "provide verbosity of validator",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  validatorTrustedGethFlag,
			Usage: "provide host:port for trusted geth, default: http://127.0.0.1:8565",
			Value: "http://127.0.0.1:8565",
		},
		&cli.StringFlag{
			Name:  validatorWalletPasswordFileFlag,
			Usage: "location of file password that you used for generation keys from deposit-cli",
			Value: "./password.txt",
		},
		&cli.StringFlag{
			Name:  validatorOutputDirFlag,
			Usage: "provide output destination of validator",
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
	startFlags = append(startFlags, fmt.Sprintf("--miner.threads=%s", ctx.String(gethMinerThreads)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.gaslimit=%s", ctx.String(gethMinerGaslimit)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.etherbase=%s", ctx.String(gethMinerEtherbaseFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--authrpc.jwtsecret=%s", ctx.String(gethAuthJWTSecretFlag)))

	return
}

func prepareValidatorStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, fmt.Sprintf("--beacon-rpc-provider=%s", ctx.String(validatorPrysmRpcProviderFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(validatorVerbosityFlag)))
	//TODO: provide flag: see --grpc-gateway-corsdomain + --grpc-gateway-port | startFlags = append(startFlags, fmt.Sprintf("--beacon-rpc-provider=%s", ctx.String(validatorTrustedGethFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--wallet-password-file=%s", ctx.String(validatorWalletPasswordFileFlag)))

	logFileFlag := prepareLogfileFlag(ctx, validatorOutputDirFlag, validatorDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	return
}

func preparePrysmStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, fmt.Sprintf("--genesis-state=%s", ctx.String(prysmGenesisStateFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--bootstrap-node=%s", ctx.String(prysmBootnodesFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--peer=%s", ctx.String(prysmPeerFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http-web3provider=%s", ctx.String(prysmWeb3ProviderFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--deposit-contract=%s", ctx.String(prysmDepositContractFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--contract-deployment-block=%s", ctx.String(prysmContractDeploymentBlockFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(prysmVerbosityFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--min-sync-peers=%s", ctx.String(prysmMinSyncPeersFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--p2p-max-peers=%s", ctx.String(prysmMaxSyncPeersFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--p2p-host-ip=%s", ctx.String(prysmP2pHostFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--p2p-local-ip=%s", ctx.String(prysmP2pLocalFlag)))

	logFileFlag := prepareLogfileFlag(ctx, prysmOutputDirFlag, prysmDependencyName)
	if logFileFlag != "" {
		startFlags = append(startFlags, logFileFlag)
	}

	if ctx.Bool(prysmTestnetFlag) {
		startFlags = append(startFlags, "--goerli")
	}

	return
}
