package main

import "github.com/urfave/cli/v2"

// selectNetworkFor accepts a CLI func as an argument, and adjusts all values that need to be changed depending on
// network passed as a flag
func selectNetworkFor(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		mainnetEnabled := ctx.Bool(mainnetEnabledFlag)
		testnetEnabled := ctx.Bool(testnetEnabledFlag)
		devnetEnabled := ctx.Bool(devnetEnabledFlag)

		enabledCount := boolToInt(mainnetEnabled) + boolToInt(testnetEnabled) + boolToInt(devnetEnabled)

		if enabledCount > 1 {
			return errMoreNetworksSelected
		}
		if enabledCount == 0 {
			mainnetEnabled = true
		}

		return f(ctx)
	}
}
