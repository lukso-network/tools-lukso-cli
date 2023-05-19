package main

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/pid"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func importValidator(ctx *cli.Context) (err error) {
	if len(os.Args) < 3 {
		return errNotEnoughArguments
	}

	err = cfg.Read()
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
	}

	switch cfg.Consensus() {
	case prysmDependencyName:
		err = importPrysmValidator(ctx)
	case lighthouseDependencyName:
		err = importLighthouseValidator(ctx)
	}

	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	return
}

func importPrysmValidator(ctx *cli.Context) (err error) {
	args := []string{
		"accounts",
		"import",
		"--wallet-dir",
		ctx.String(validatorWalletDirFlag),
		"--keys-dir",
		ctx.String(validatorKeysFlag),
	}

	validatorPass := ctx.String(validatorPasswordFlag)

	if validatorPass != "" {
		args = append(args, "--account-password-file", validatorPass)
	}

	initCommand := exec.Command("validator", args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err = initCommand.Run()
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while importing validator accounts: %v", err), 1)
	}

	return nil
}

func importLighthouseValidator(ctx *cli.Context) (err error) {
	args := []string{
		"am",
		"validator",
		"import",
		"--directory",
		ctx.String(validatorKeysFlag),
		"--datadir",
		ctx.String(validatorWalletDirFlag),
	}

	passwordFile := ctx.String(validatorPasswordFlag)
	if passwordFile != "" {
		args = append(args, "--password-file", passwordFile)
	}

	initCommand := exec.Command(lighthouseDependencyName, args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err = initCommand.Run()
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while importing validator accounts: %v", err), 1)
	}

	return nil
}

func startValidator(ctx *cli.Context) (err error) {
	err = cfg.Read()
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
	}

	switch cfg.Consensus() {
	case prysmDependencyName:
		err = startPrysmValidator(ctx)
	case lighthouseDependencyName:
		err = startLighthouseValidator(ctx)
	}

	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	return
}

func startPrysmValidator(ctx *cli.Context) (err error) {
	validatorFlags, passwordPipe, err := prepareValidatorStartFlags(ctx)
	if passwordPipe != nil && passwordPipe.Name() != "" {
		defer os.Remove(passwordPipe.Name())
	}
	if err != nil {
		return err
	}
	if !fileExists(fmt.Sprintf("%s/direct/accounts/all-accounts.keystore.json", ctx.String(validatorKeysFlag))) { // path to imported keys
		return exit("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.", 1)
	}

	err = clientDependencies[validatorDependencyName].Start(validatorFlags, ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	passwordPipe.Close()

	log.Info("⚙️  Please wait a few seconds while your password is being validated...")
	time.Sleep(time.Second * 10) // should be enough

	logFile, err := getLastFile(ctx.String(logFolderFlag), validatorDependencyName)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while getting latest log file: %v", err), 1)
	}

	logs, err := os.ReadFile(ctx.String(logFolderFlag) + "/" + logFile)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while reading log file: %v", err), 1)
	}

	if strings.Contains(string(logs), wrongPassword) {
		err = cfg.Read()
		if err != nil {
			return exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
		}

		err = clientDependencies[cfg.Consensus()].Stop()
		if err != nil {
			return exit(fmt.Sprintf("❌  There was an error while stopping consensus: %v", err), 1)
		}

		err = clientDependencies[cfg.Execution()].Stop()
		if err != nil {
			return exit(fmt.Sprintf("❌  There was an error while stopping execution: %v", err), 1)
		}

		return exit("❌  Incorrect password, please restart and try again", 1)
	}

	log.Info("✅  Validator started! Use 'lukso logs validator' to see the logs.")

	return
}

func startLighthouseValidator(ctx *cli.Context) (err error) {
	lighthouseValidatorFlags, passwordPipe, err := prepareLighthouseValidatorFlags(ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while preparing lighthouse validator flags: %v", err), 1)
	}
	if passwordPipe != nil && passwordPipe.Name() != "" {
		defer func() {
			os.Remove(passwordPipe.Name())
		}()
	}

	if !fileExists(fmt.Sprintf("%s/validators", ctx.String(validatorKeysFlag))) { // path to imported keys
		return exit("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.", 1)
	}

	lighthouseValidatorFlags = append([]string{"vc"}, lighthouseValidatorFlags...)
	startCommand := exec.Command(lighthouseDependencyName, lighthouseValidatorFlags...)

	fmt.Println(startCommand.String())
	err = startCommand.Start()
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting lighthouse validator flags: %v", err), 1)
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, validatorDependencyName)
	err = pid.Create(pidLocation, startCommand.Process.Pid)

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

func exit(message string, exitCode int) error {
	log.Error(message)

	os.Exit(exitCode)

	return nil // so we can return commands with this func
}
