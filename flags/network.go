package flags

import "github.com/urfave/cli/v2"

var (
	// Network flags
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

	NetworkFlags = []cli.Flag{
		mainnetEnabledFlag,
		testnetEnabledFlag,
		devnetEnabledFlag,
	}
)
