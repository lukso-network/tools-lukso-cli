package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
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
		return cli.Exit(fmt.Sprintf("❌  There was an error while importing validator accounts: %v", err), 1)
	}

	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("🔄  Starting Validator")
	validatorFlags, passwordPipe, err := prepareValidatorStartFlags(ctx)
	if passwordPipe != "" {
		defer os.Remove(passwordPipe)
	}
	if err != nil {
		return err
	}
	if !fileExists(fmt.Sprintf("%s/direct/accounts/all-accounts.keystore.json", ctx.String(validatorKeysFlag))) { // path to imported keys
		log.Error("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.")

		return nil
	}

	err = clientDependencies[validatorDependencyName].Start(validatorFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Validator started! Use 'lukso logs' to see the logs.")

	return nil
}
