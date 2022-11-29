package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

func (dependency *ClientDependency) Run(
	commandName string,
	arguments []string,
	attachStdoutAndErr bool,
) (err error) {
	command := exec.Command(commandName, arguments...)

	if attachStdoutAndErr {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()

		return
	}

	err = command.Start()

	return
}

func startClients(ctx *cli.Context) error {
	log.Info("Starting all clients")

	err := startGethDetached(ctx)
	if err != nil {
		return err
	}

	err = startPrysmDetached(ctx)
	if err != nil {
		return err
	}

	err = startValidatorDetached(ctx)
	if err != nil {
		return err
	}

	return err
}

func startGeth(ctx *cli.Context) error {
	log.Info("Starting Geth")

	stdAttached := ctx.Bool(gethStdOutputFlag)

	err := clientDependencies[gethDependencyName].Run(gethDependencyName, prepareGethStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	stdAttached := ctx.Bool(prysmStdOutputFlag)

	err := clientDependencies[prysmDependencyName].Run(prysmDependencyName, preparePrysmStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("Starting Validator")

	stdAttached := ctx.Bool(validatorStdOutputFlag)

	err := clientDependencies[validatorDependencyName].Run(validatorDependencyName, prepareValidatorStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}

func startGethDetached(ctx *cli.Context) error {
	log.Info("Starting Geth")

	err := clientDependencies[gethDependencyName].Run(gethDependencyName, prepareGethStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysmDetached(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	err := clientDependencies[prysmDependencyName].Run(prysmDependencyName, preparePrysmStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidatorDetached(ctx *cli.Context) error {
	log.Info("Starting Validator")

	err := clientDependencies[validatorDependencyName].Run(validatorDependencyName, prepareValidatorStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}
