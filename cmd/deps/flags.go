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
	prysmStdOutputFlag               = "prysm-std-output"
	prysmWeb3ProviderFlag            = "prysm-web3provider"
	prysmDepositContractFlag         = "prysm-deposit-contract"
	prysmContractDeploymentBlockFlag = "prysm-deposit-deployment"
	prysmVerbosityFlag               = "prysm-verbosity"
	prysmMinSyncPeersFlag            = "prysm-min-sync-peers"
	prysmMaxSyncPeersFlag            = "prysm-max-sync-peers"
	prysmP2pHostFlag                 = "prysm-p2p-host"
	prysmP2pLocalFlag                = "prysm-p2p-local"
	prysmOrcProviderFlag             = "prysm-orc-provider"
	prysmDisableSyncFlag             = "prysm-disable-sync"
	prysmOutputFileFlag              = "prysm-output-file"

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
			Usage: "nickname:STATS_LOGIN_SECRET@GETH_STATS_HOST",
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
			Value: "",
		},
		&cli.StringFlag{
			Name:  gethHttpOriginFlag,
			Usage: "this flag sets up http accepted origins, default not set",
			Value: "",
		},
		&cli.StringFlag{
			Name:  gethNatFlag,
			Usage: "this flag sets up http nat to assign static ip for geth, default not set. Example extip:172.16.254.4",
			Value: "",
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
	prysmStartFlags = []cli.Flag{}

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
	validatorStartFlags = []cli.Flag{}
)

func prepareGethStartFlags(ctx *cli.Context) (startFlags []string) {
	// parse all runtime-related geth flags, one by one and append them
	startFlags = append(startFlags, fmt.Sprintf("--ethstats %s", ctx.String(gethEthstatsFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--bootnodes %s", ctx.String(gethBootnodesFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--networkid %s", ctx.String(gethNetworkIDFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--port %s", ctx.String(gethPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.api %s", ctx.String(gethHttpApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.api %s", ctx.String(gethWSApiFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.port %s", ctx.String(gethWSPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--miner.etherbase %s", ctx.String(gethEtherbaseFlag)))
	if ctx.Bool(gethNotifyFlag) {
		startFlags = append(startFlags, "--miner.notify.full")
	}

	startFlags = append(startFlags, fmt.Sprintf("--verbosity %s", ctx.String(gethVerbosityFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.port %s", ctx.String(gethHttpPortFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--ws.origins %s", ctx.String(gethWsOriginFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--http.corsdomain %s", ctx.String(gethHttpOriginFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--nat %s", ctx.String(gethNatFlag)))
	return
}
