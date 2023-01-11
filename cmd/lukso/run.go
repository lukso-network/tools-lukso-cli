package main

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

var (
	errorCouldntStart = errors.New("Couldn't start client ")
	errorFlagMissing  = errors.New("Couldn't find given flag ")
)

func (dependency *ClientDependency) Start(
	arguments []string,
	attachStdoutAndErr bool,
	ctx *cli.Context,
) (err error) {
	command := exec.Command(dependency.name, arguments...)

	if attachStdoutAndErr {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()

		return
	}

	// since geth removed --logfile flag we have to manually adjust geth's stdout
	if dependency.name == gethDependencyName {
		var (
			logFile  *os.File
			fullPath string
		)

		gethLogDir := ctx.String(gethOutputDirFlag)
		if gethLogDir == "" {
			return errorFlagMissing
		}

		fullPath, err = prepareTimestampedFile(gethLogDir, gethDependencyName)
		if err != nil {
			return
		}

		err = os.WriteFile(fullPath, []byte{}, 0750)
		if err != nil {
			return
		}

		logFile, err = os.OpenFile(fullPath, os.O_RDWR, 0750)
		if err != nil {
			return
		}

		command.Stdout = logFile
		command.Stderr = logFile
	}

	err = command.Start()

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
	log.Info("Running geth init first...")

	err := initGeth()
	if err != nil {
		log.Errorf("There was an error while initalizing geth. Error: %v", err)
	}

	log.Info("Starting Geth")

	stdAttached := ctx.Bool(gethStdOutputFlag)

	err = clientDependencies[gethDependencyName].Start(prepareGethStartFlags(ctx), stdAttached, ctx)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	stdAttached := ctx.Bool(prysmStdOutputFlag)

	err := clientDependencies[prysmDependencyName].Start(preparePrysmStartFlags(ctx), stdAttached, ctx)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("Starting Validator")

	stdAttached := ctx.Bool(validatorStdOutputFlag)

	err := clientDependencies[validatorDependencyName].Start(prepareValidatorStartFlags(ctx), stdAttached, ctx)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}

func startGethDetached(ctx *cli.Context) error {
	log.Info("Starting Geth")

	err := clientDependencies[gethDependencyName].Start(prepareGethStartFlags(ctx), false, ctx)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysmDetached(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	err := clientDependencies[prysmDependencyName].Start(preparePrysmStartFlags(ctx), false, ctx)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidatorDetached(ctx *cli.Context) error {
	log.Info("Starting Validator")

	err := clientDependencies[validatorDependencyName].Start(prepareValidatorStartFlags(ctx), false, ctx)
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

func initGeth() (err error) {
	command := exec.Command("geth", "init", clientDependencies[gethGenesisDependencyName].filePath)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
