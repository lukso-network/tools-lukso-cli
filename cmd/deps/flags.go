package main

import (
	"github.com/urfave/cli/v2"
)

const (
	// Pandora related flag names
	gethTagFlag        = "geth-tag"
	gethCommitHashFlag = "geth-commit-hash"
	gethDatadirFlag    = "geth-datadir"

	// Common for prysm client
	prysmChainConfigFlag = "prysm-chain-config"

	// Validator related flag names
	validatorTagFlag     = "validator-tag"
	validatorDatadirFlag = "validator-datadir"

	// Prysm related flag names
	prysmTagFlag     = "prysm-tag"
	prysmDatadirFlag = "prysm-datadir"

	acceptTermsOfUseFlagName = "accept-terms-of-use"
)

var (
	appFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  acceptTermsOfUseFlagName,
			Usage: "Accept terms of use. Default: true",
			Value: true,
		},
	}

	gethFlags = []cli.Flag{
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
	validatorFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  validatorTagFlag,
			Usage: "provide tag for validator binary",
			Value: "v3.1.2",
		},
	}
	prysmFlags = []cli.Flag{
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
)
