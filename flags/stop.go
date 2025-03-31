package flags

import "github.com/urfave/cli/v2"

var StopFlags = []cli.Flag{
	executionSelectedFlag,
	consensusSelectedFlag,
	validatorSelectedFlag,
}
