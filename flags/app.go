package flags

import "github.com/urfave/cli/v2"

var (
	AppFlags []cli.Flag

	// Layer flags
	consensusSelectedFlag = &cli.BoolFlag{
		Name:  ConsensusFlag,
		Usage: "Run for consensus",
		Value: false,
	}
	executionSelectedFlag = &cli.BoolFlag{
		Name:  ExecutionFlag,
		Usage: "Run for execution",
		Value: false,
	}
	validatorSelectedFlag = &cli.BoolFlag{
		Name:  ValidatorFlag,
		Usage: "Run for validator",
		Value: false,
	}
)
