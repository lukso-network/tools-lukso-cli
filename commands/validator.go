package commands

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func (c *commander) ValidatorImport(ctx *cli.Context) (err error) {
	if len(os.Args) < 3 {
		return errors.ErrNotEnoughArguments
	}

	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
	}

	selectedValidator := cfg.Validator()
	validatorClient, ok := clients.AllClients[selectedValidator].(clients.ValidatorBinaryDependency)
	if !ok {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	err = validatorClient.Import(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while importing validator: %v", err), 1)
	}

	return
}

func (c *commander) ValidatorList(ctx *cli.Context) (err error) {
	network := "mainnet" // Set the default network to mainnet

	if ctx.Bool(flags.TestnetFlag) {
		network = "testnet"
	} else if ctx.Bool(flags.DevnetFlag) {
		network = "devnet"
	}

	return executeValidatorList(ctx, network)
}

func (c *commander) ValidatorExit(ctx *cli.Context) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit("❌  There was an error while reading config file", 1)
	}

	selectedValidator := cfg.Validator()
	validatorClient, ok := clients.AllClients[selectedValidator].(clients.ValidatorBinaryDependency)
	if !ok {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	err = validatorClient.Exit(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while exiting validator: %v", err), 1)
	}

	return
}

func executeValidatorList(ctx *cli.Context, network string) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit("❌  There was an error while reading config file", 1)
	}

	selectedValidator := cfg.Validator()
	validatorClient, ok := clients.AllClients[selectedValidator].(clients.ValidatorBinaryDependency)
	if !ok {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	err = validatorClient.List(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while listing validator keys: %v", err), 1)
	}

	return nil
}
