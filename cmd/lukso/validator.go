package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"time"
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
	if !fileExists(fmt.Sprintf("%s/direct/accounts/all-accounts.keystore.json", ctx.String(validatorKeysFlag))) { // path to imported keys
		log.Error("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.")

		return nil
	}

	validatorFlags, passwordPipe, err := prepareValidatorStartFlags(ctx)
	if passwordPipe != "" {
		defer os.Remove(passwordPipe)
	}
	if err != nil {
		return err
	}

	// before starting validate password
	walletDir := ctx.String(validatorKeysFlag)

	validateCommand := exec.Command("validator", "accounts", "list", "--wallet-dir", walletDir, "--wallet-password-file", passwordPipe)
	buf := new(bytes.Buffer)

	validateCommand.Stdout = buf
	validateCommand.Stderr = buf

	fmt.Println(validateCommand.String())
	err = validateCommand.Start()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while validating password: %v", err), 1)
	}

	time.Sleep(time.Second * 3)

	fmt.Println(buf.String())

	log.Error("❌  Password incorrect, please restart and try again")

	err = clientDependencies[validatorDependencyName].Start(validatorFlags, ctx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

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

//err = cfg.Read()
//if err != nil {
//return cli.Exit(fmt.Sprintf("❌  There was an error while reading config file: %v", err), 1)
//}
//
//err = clientDependencies[cfg.Execution()].Stop()
//if err != nil {
//return cli.Exit(fmt.Sprintf("❌  There was an error while stopping execution: %v", err), 1)
//}
//
//err = clientDependencies[cfg.Consensus()].Stop()
//if err != nil {
//return cli.Exit(fmt.Sprintf("❌  There was an error while stopping consensus: %v", err), 1)
//}
//f, err := os.Open(passwordPipe)
//if err != nil {
//return cli.Exit(fmt.Sprintf("❌  There was an error while opening password pipe: %v", err), 1)
//}
//
//_, err = f.Write(password)
//if err != nil {
//return cli.Exit(fmt.Sprintf("❌  There was an error while writing password: %v", err), 1)
//}
