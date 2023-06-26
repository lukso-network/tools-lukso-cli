package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func resetClients(ctx *cli.Context) error {
	if isAnyRunning() {
		return nil
	}
	if !cfg.Exists() {
		return exit(folderNotInitialized, 1)
	}

	message := fmt.Sprintf("⚠️  WARNING: THIS ACTION WILL DELETE THE CLIENTS DATA DIRECTORIES ⚠️\n"+
		"\nDirectories that will be deleted:\n"+
		"- %s\n- %s\n- %s\n\nAre you sure you want to continue? [y/N]: ", ctx.String(gethDatadirFlag), ctx.String(prysmDatadirFlag), ctx.String(validatorDatadirFlag))

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") {
		log.Info("❌  Aborting...")

		return nil
	}

	err := resetExecution(ctx)
	if err != nil {
		return err
	}

	err = resetConsensus(ctx)
	if err != nil {
		return err
	}

	err = resetValidator(ctx)
	if err != nil {
		return err
	}

	log.Info("🗑️  Data directories deleted.")

	return nil
}

func resetExecution(ctx *cli.Context) error {
	dataDirPath := ctx.String(gethDatadirFlag)
	if dataDirPath == "" {
		return exit(fmt.Sprintf("%v- %s", errFlagMissing, gethDatadirFlag), 1)
	}

	return os.RemoveAll(dataDirPath)
}

func resetConsensus(ctx *cli.Context) error {
	dataDirPath := ctx.String(prysmDatadirFlag)
	if dataDirPath == "" {
		return exit(fmt.Sprintf("%v- %s", errFlagMissing, prysmDatadirFlag), 1)
	}

	return os.RemoveAll(dataDirPath)
}

func resetValidator(ctx *cli.Context) error {
	dataDirPath := ctx.String(validatorDatadirFlag)
	if dataDirPath == "" {
		return exit(fmt.Sprintf("%v- %s", errFlagMissing, validatorDatadirFlag), 1)
	}

	return os.RemoveAll(dataDirPath)
}
