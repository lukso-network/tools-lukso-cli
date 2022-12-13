package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	// geth related flag names
	gethTagFlag         = "geth-tag"
	gethCommitHashFlag  = "geth-commit-hash"
	gethDatadirFlag     = "geth-datadir"
	gethEthstatsFlag    = "geth-ethstats"
	gethBootnodesFlag   = "geth-bootnodes"
	gethNetworkIDFlag   = "geth-networkid"
	gethPortFlag        = "geth-port"
	gethHttpApiFlag     = "geth-http-apis"
	gethWSApiFlag       = "geth-ws-apis"
	gethWSPortFlag      = "geth-websocket-port"
	gethEtherbaseFlag   = "geth-etherbase"
	gethGenesisFileFlag = "geth-genesis"
	gethNotifyFlag      = "geth-notify"
	gethVerbosityFlag   = "geth-verbosity"
	gethHttpPortFlag    = "geth-http-port"
	gethStdOutputFlag   = "geth-std-output"
	gethWsOriginFlag    = "geth-ws-origin"
	gethHttpOriginFlag  = "geth-http-origin"
	gethNatFlag         = "geth-nat"

	// Common for prysm client
	prysmChainConfigFlag = "prysm-chain-config"

	// Validator related flag names
	validatorTagFlag                = "validator-tag"
	validatorDatadirFlag            = "validator-datadir"
	validatorPrysmRpcProviderFlag   = "validator-prysm-rpc"
	validatorVerbosityFlag          = "validator-verbosity"
	validatorTrustedGethFlag        = "validator-trusted-geth"
	validatorWalletPasswordFileFlag = "validator-wallet-password-file"
	validatorOutputFileFlag         = "validator-output-file"
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
	prysmOutputFileFlag              = "prysm-output-file"
	prysmStdOutputFlag               = "prysm-std-output"

	acceptTermsOfUseFlagName = "accept-terms-of-use"
)

var (
	downloadFlags []cli.Flag
	updateFlags   []cli.Flag
	startFlags    []cli.Flag
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
		&cli.StringFlag{
			Name:  gethDatadirFlag,
			Usage: "provide a path you would like to store your data",
			Value: "./geth",
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
			Name:  gethEthstatsFlag,
			Usage: "URL of an ethstats service, should be of scheme: nickname:STATS_LOGIN_SECRET@GETH_STATS_HOST",
			Value: "",
		},
		&cli.StringFlag{
			Name:  gethBootnodesFlag,
			Usage: "Default value should be ok for test network. Otherwise provide Comma separated enode urls, see at https://geth.ethereum.org/docs/getting-started/private-net.",
			Value: "enode://967db4f56ad0a1a35e3d30632fa600565329a23aff50c9762181810166f3c15b078cca522f930d1a2747778893232336bffd1ea5d2ca60543f1801d4360ea63a@35.204.255.172:0?discport=30301",
		},
		&cli.StringFlag{
			Name:  gethNetworkIDFlag,
			Usage: "provide network id if must be different than default",
			Value: "4004181",
		},
		&cli.StringFlag{
			Name:  gethPortFlag,
			Usage: "provide port for geth",
			Value: "30405",
		},
		&cli.StringFlag{
			Name:  gethHttpApiFlag,
			Usage: "comma separated apis",
			Value: "eth,net",
		},
		&cli.StringFlag{
			Name:  gethHttpPortFlag,
			Usage: "port used in geth http communication",
			Value: "8565",
		},
		&cli.StringFlag{
			Name:  gethWSApiFlag,
			Usage: "comma separated apis",
			Value: "eth,net",
		},
		&cli.StringFlag{
			Name:  gethWSPortFlag,
			Usage: "port for geth api",
			Value: "8546",
		},
		&cli.StringFlag{
			Name:  gethEtherbaseFlag,
			Usage: "your ECDSA public key used to get rewards on geth chain",
			// yes, If you won't set it up, I'll get rewards ;]
			Value: "0x59E3dADc83af3c127a2e29B12B0E86109Bb6d838",
		},
		&cli.StringFlag{
			Name:  gethGenesisFileFlag,
			Usage: "remote genesis file that will be downloaded to spin up the network",
			// yes, If you won't set it up, I'll get rewards ;]
			Value: "https://storage.googleapis.com/l16-common/geth/geth_private_testnet_genesis.json",
		},
		&cli.StringFlag{
			Name:  gethNotifyFlag,
			Usage: "this flag is used to geth engine to notify validator and orchestrator",
			Value: "ws://127.0.0.1:7878,http://127.0.0.1:7877",
		},
		&cli.StringFlag{
			Name:  gethVerbosityFlag,
			Usage: "this flag sets up verbosity for geth",
			Value: "3",
		},
		&cli.StringFlag{
			Name:  gethWsOriginFlag,
			Usage: "this flag sets up websocket accepted origins, default not set",
			Value: "ws://127.0.0.1:7878",
		},
		&cli.StringFlag{
			Name:  gethHttpOriginFlag,
			Usage: "this flag sets up http accepted origins, default not set",
			Value: "http://127.0.0.1:8008",
		},
		&cli.StringFlag{
			Name:  gethNatFlag,
			Usage: "this flag sets up http nat to assign static ip for geth, default not set. Example extip:172.16.254.4",
			Value: "extip:172.16.254.4",
		},
		&cli.BoolFlag{
			Name:  gethStdOutputFlag,
			Usage: "set geth output to stdout",
			Value: false,
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
			Value: "./prysm/v0.0.18-delta/vanguard_private_testnet_genesis.ssz",
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
			Name:  prysmOutputFileFlag,
			Usage: "provide output destination of prysm",
			Value: "./prysm.log",
		},
		&cli.BoolFlag{
			Name:  prysmStdOutputFlag,
			Usage: "set prysm output to stdout",
			Value: false,
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
			Name:  validatorOutputFileFlag,
			Usage: "provide output destination of validator",
			Value: "./validator.log",
		},
		&cli.BoolFlag{
			Name:  validatorStdOutputFlag,
			Usage: "set validator output to stdout",
			Value: false,
		},
	}
)

func prepareGethStartFlags(ctx *cli.Context) (startFlags []string) {
	// parse all runtime-related geth flags, one by one and append them
	startFlags = append(startFlags, fmt.Sprintf("--ethstats=%s", ctx.String(gethEthstatsFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--bootnodes=%s", ctx.String(gethBootnodesFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--networkid=%s", ctx.String(gethNetworkIDFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--port=%s", ctx.String(gethPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.api=%s", ctx.String(gethHttpApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.api=%s", ctx.String(gethWSApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.port=%s", ctx.String(gethWSPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.etherbase=%s", ctx.String(gethEtherbaseFlag)))
	if ctx.Bool(gethNotifyFlag) {
		startFlags = append(startFlags, "--miner.notify.full")
	}

	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(gethVerbosityFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.port=%s", ctx.String(gethHttpPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.origins=%s", ctx.String(gethWsOriginFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.corsdomain=%s", ctx.String(gethHttpOriginFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--nat=%s", ctx.String(gethNatFlag)))

	return
}

func prepareValidatorStartFlags(ctx *cli.Context) (startFlags []string) {
	startFlags = append(startFlags, fmt.Sprintf("--beacon-rpc-provider=%s", ctx.String(validatorPrysmRpcProviderFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--verbosity=%s", ctx.String(validatorVerbosityFlag)))
	//TODO: provide flag: see --grpc-gateway-corsdomain + --grpc-gateway-port | startFlags = append(startFlags, fmt.Sprintf("--beacon-rpc-provider=%s", ctx.String(validatorTrustedGethFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--wallet-password-file=%s", ctx.String(validatorWalletPasswordFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--log-file=%s", ctx.String(validatorOutputFileFlag)))

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
	startFlags = append(startFlags, fmt.Sprintf("--log-file=%s", ctx.String(prysmOutputFileFlag)))

	return
}
