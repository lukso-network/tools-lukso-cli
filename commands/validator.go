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

func ImportValidator(ctx *cli.Context) (err error) {
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

func ListValidator(ctx *cli.Context) (err error) {
	network := "mainnet" // Set the default network to mainnet

	if ctx.Bool(flags.TestnetFlag) {
		network = "testnet"
	} else if ctx.Bool(flags.DevnetFlag) {
		network = "devnet"
	}

	return executeValidatorList(ctx, network)
}

func ExitValidator(ctx *cli.Context) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while reading config file: %v", err), 1)
	}

	var (
		selectedValidator clients.ValidatorBinaryDependency
		selectedConsensus clients.ClientBinaryDependency
	)

	switch cfg.Validator() {
	case clients.PrysmValidator.Name():
		selectedValidator = clients.PrysmValidator
		selectedConsensus = clients.Prysm
	case clients.LighthouseValidator.Name():
		selectedValidator = clients.LighthouseValidator
		selectedConsensus = clients.Lighthouse
	}

	if !selectedConsensus.IsRunning() {
		return utils.Exit("⚠️  Please make sure that your validator client is running before exiting", 1)
	}

	err = selectedValidator.Exit(ctx)
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
