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
		"--wallet-dir",
		ctx.String(validatorWalletDirFlag),
	}

	validatorPass := ctx.String(validatorPasswordFlag)

	if validatorPass != "" {
		args = append(args, "--account-password-file", validatorPass)
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

func startValidator(ctx *cli.Context) (err error) {
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
		return cli.Exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	log.Info("✅  Validator started! Use 'lukso logs' to see the logs.")

	return
}

func executeValidatorList(network string) error {
	cmd := exec.Command("validator", "accounts", "list", "--wallet-dir", fmt.Sprintf("%s-keystore", network))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not fetch imported keys within the validator wallet: %v", err)
	}

	return nil
}

func listValidator(ctx *cli.Context) error {
	network := "mainnet" // Set the default network to mainnet

	if ctx.Bool(testnetFlag) {
		network = "testnet"
	} else if ctx.Bool(devnetFlag) {
		network = "devnet"
	}

	return executeValidatorList(network)
}
