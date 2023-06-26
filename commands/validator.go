package commands

import (
	"fmt"
	"os"
	"os/exec"

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

	return executeValidatorList(network)
}

func executeValidatorList(network string) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit("❌  There was an error while reading config file", 1)
	}

	var cmd *exec.Cmd
	switch cfg.Consensus() {
	case clients.Prysm.Name():
		cmd = exec.Command("validator", "accounts", "list", "--wallet-dir", fmt.Sprintf("%s-keystore", network))
	case clients.Lighthouse.Name():
		cmd = exec.Command(clients.Lighthouse.CommandName(), "am", "validator", "list", "--datadir", fmt.Sprintf("%s-keystore", network))
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not fetch imported keys within the validator wallet: %v ", err)
	}

	return nil
}
