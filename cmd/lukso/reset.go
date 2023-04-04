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

	message := fmt.Sprintf("WARNING: THIS ACTION WILL REMOVE DATA DIRECTORIES FROM ALL OF RUNNING CLIENTS.\n"+
		"Are you sure you want to continue?\nDirectories that will be deleted:\n"+
		"- %s\n- %s\n- %s\n[Y/n]: ", ctx.String(gethDatadirFlag), ctx.String(prysmDatadirFlag), ctx.String(validatorDatadirFlag))

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Info("Aborting...")

		return nil
	}

	err := resetGeth(ctx)
	if err != nil {
		return err
	}

	err = resetPrysm(ctx)
	if err != nil {
		return err
	}

	err = resetValidator(ctx)

	return err
}

func resetGeth(ctx *cli.Context) error {
	dataDirPath := ctx.String(gethDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}

func resetPrysm(ctx *cli.Context) error {
	dataDirPath := ctx.String(prysmDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}

func resetValidator(ctx *cli.Context) error {
	dataDirPath := ctx.String(validatorDatadirFlag)
	if dataDirPath == "" {
		return errFlagMissing
	}

	return os.RemoveAll(dataDirPath)
}
