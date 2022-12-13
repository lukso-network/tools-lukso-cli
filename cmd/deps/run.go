package main

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"time"
)

var errorCouldntStart = errors.New("Couldn't start client")

func (dependency *ClientDependency) Start(
	arguments []string,
	attachStdoutAndErr bool,
) (err error) {
	command := exec.Command(dependency.name, arguments...)

	if attachStdoutAndErr {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if _, ok := err.(*exec.ExitError); ok {
			return errorCouldntStart
		}

		return
	}

	log.Info("Waiting for initialization...")
	err = command.Start()

	time.Sleep(time.Second) // sleep to wait for initialization error

	if _, ok := err.(*exec.ExitError); ok {
		return errorCouldntStart
	}

	return
}

func (dependency *ClientDependency) Stop() error {
	// since there aren't commands for stopping clients we have to kill them manually
	command := exec.Command("killall", "-HUP", dependency.name)

	return command.Start()
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

	err := clientDependencies[gethDependencyName].Start(prepareGethStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	stdAttached := ctx.Bool(prysmStdOutputFlag)

	err := clientDependencies[prysmDependencyName].Start(preparePrysmStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("Starting Validator")

	stdAttached := ctx.Bool(validatorStdOutputFlag)

	err := clientDependencies[validatorDependencyName].Start(prepareValidatorStartFlags(ctx), stdAttached)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}

func startGethDetached(ctx *cli.Context) error {
	log.Info("Starting Geth")

	err := clientDependencies[gethDependencyName].Start(prepareGethStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysmDetached(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	err := clientDependencies[prysmDependencyName].Start(preparePrysmStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidatorDetached(ctx *cli.Context) error {
	log.Info("Starting Validator")

	err := clientDependencies[validatorDependencyName].Start(prepareValidatorStartFlags(ctx), false)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}

func stopClients(ctx *cli.Context) error {
	err := stopClient(clientDependencies[gethDependencyName])(ctx)
	if err != nil {
		return err
	}

	err = stopClient(clientDependencies[prysmDependencyName])(ctx)
	if err != nil {
		return err
	}

	err = stopClient(clientDependencies[validatorDependencyName])(ctx)

	return err
}

func stopClient(dependency *ClientDependency) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Infof("Stopping %s", dependency.name)

		err := dependency.Stop()

		return err
	}
}
