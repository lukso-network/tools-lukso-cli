package flags

import "github.com/urfave/cli/v2"

var StatusPeersFlags = []cli.Flag{
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
