package flags

import "github.com/urfave/cli/v2"

var InitFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  ReinitFlag,
		Usage: "Reinitialize directory",
		Value: false,
	},
	&cli.StringFlag{
		Name:  IpFlag,
		Usage: "IP used by clients for P2P communication. Can be set to 'auto' for automatic IP configuration, 'disabled' for skipping the IP configuration and '[ipv4]' for manual configuration",
		Value: "disabled",
	},
}
