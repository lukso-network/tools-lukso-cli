package flags

import "github.com/urfave/cli/v2"

var UpdateConfigFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  AllFlag,
		Usage: "Updates all configuration files, including configurations for each client",
		Value: false,
	},
}
