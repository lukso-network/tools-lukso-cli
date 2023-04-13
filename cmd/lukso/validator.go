package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func importValidator(ctx *cli.Context) error {
	if len(os.Args) < 3 {
		return errNotEnoughArguments
	}

	args := []string{
		"accounts",
		"import",
	}

	// we don't want to pass those flags
	mainnet := fmt.Sprintf("--%s", mainnetFlag)
	testnet := fmt.Sprintf("--%s", testnetFlag)
	devnet := fmt.Sprintf("--%s", devnetFlag)
	walletDir := "--wallet-dir"

	for _, osArg := range os.Args[3:] {
		if osArg == mainnet || osArg == testnet || osArg == devnet {
			continue
		}

		args = append(args, osArg)
	}

	isWalletProvided := false
	walletDefault := ctx.String(validatorWalletDirFlag)

	for _, arg := range args {
		if arg == walletDir {
			isWalletProvided = true
		}
	}

	if !isWalletProvided {
		args = append(args, walletDir, walletDefault)
	}

	initCommand := exec.Command("validator", args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err := initCommand.Run()
	if err != nil {
		return err
	}

	return nil
}
