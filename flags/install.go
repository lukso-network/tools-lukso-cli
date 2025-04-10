package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
)

var (
	InstallFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  AgreeTermsFlag,
			Usage: "Automatically accept Terms and Conditions",
			Value: false,
		},
	}
	GethInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  GethTagFlag,
			Usage: "A tag of geth you would like to run",
			Value: common.GethTag,
		},
		&cli.StringFlag{
			Name:  GethCommitHashFlag,
			Usage: "A hash of commit that is bound to given release tag",
			Value: common.GethCommitHash,
		},
	}
	ErigonInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ErigonTagFlag,
			Usage: "Tag for erigon",
			Value: common.ErigonTag,
		},
	}
	NethermindInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  NethermindTagFlag,
			Usage: "Tag for nethermind",
			Value: common.NethermindTag,
		},
		&cli.StringFlag{
			Name:  NethermindCommitHashFlag,
			Usage: "A hash of commit that is bound to given release tag",
			Value: common.NethermindCommitHash,
		},
	}
	BesuInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  BesuTagFlag,
			Usage: "Tag for besu",
			Value: common.BesuTag,
		},
	}
	PrysmInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  PrysmTagFlag,
			Usage: "Tag for prysm",
			Value: common.PrysmTag,
		},
	}
	LighthouseInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LighthouseTagFlag,
			Usage: "Tag for lighthouse",
			Value: common.LighthouseTag,
		},
	}
	TekuInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  TekuTagFlag,
			Usage: "Tag for teku client",
			Value: common.TekuTag,
		},
	}
	Nimbus2InstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  Nimbus2TagFlag,
			Usage: "Tag for nimbus-eth2 client",
			Value: common.Nimbus2Tag,
		},
		&cli.StringFlag{
			Name:  Nimbus2CommitHashFlag,
			Usage: "Commit hash for nimbus-eth2 client",
			Value: common.Nimbus2CommitHash,
		},
	}
	ValidatorInstallFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  ValidatorTagFlag,
			Usage: "Tag for validator binary",
			Value: common.PrysmTag,
		},
	}
)

func init() {
	// Execution
	InstallFlags = append(InstallFlags, GethInstallFlags...)
	InstallFlags = append(InstallFlags, ErigonInstallFlags...)
	InstallFlags = append(InstallFlags, NethermindInstallFlags...)
	InstallFlags = append(InstallFlags, BesuInstallFlags...)
	// Consensus
	InstallFlags = append(InstallFlags, PrysmInstallFlags...)
	InstallFlags = append(InstallFlags, LighthouseInstallFlags...)
	InstallFlags = append(InstallFlags, TekuInstallFlags...)
	InstallFlags = append(InstallFlags, Nimbus2InstallFlags...)
	// Validator
	InstallFlags = append(InstallFlags, ValidatorInstallFlags...)
}
