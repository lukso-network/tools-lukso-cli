package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dep/configs"
)

var (
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
)

func init() {
	ValidatorImportFlags = append(ValidatorImportFlags, NetworkFlags...)
	ValidatorListFlags = append(ValidatorListFlags, NetworkFlags...)
	ValidatorExitFlags = append(ValidatorExitFlags, NetworkFlags...)
}
