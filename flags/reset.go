package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

var (
	ResetFlags          []cli.Flag
	ExecutionResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   GethDatadirFlag,
			Usage:  "geth datadir",
			Value:  configs.ExecutionMainnetDatadir,
			Hidden: true,
		},
	}
	ConsensusResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   PrysmDatadirFlag,
			Usage:  "prysm datadir",
			Value:  configs.ConsensusMainnetDatadir,
			Hidden: true,
		},
	}
	ValidatorResetFlags = []cli.Flag{
		&cli.StringFlag{
			Name:   ValidatorDatadirFlag,
			Usage:  "Validator datadir",
			Value:  configs.ValidatorMainnetDatadir,
			Hidden: true,
		},
	}
)

func init() {
	ResetFlags = append(ResetFlags, ExecutionResetFlags...)
	ResetFlags = append(ResetFlags, ConsensusResetFlags...)
	ResetFlags = append(ResetFlags, ValidatorResetFlags...)
	ResetFlags = append(ResetFlags, NetworkFlags...)
}
