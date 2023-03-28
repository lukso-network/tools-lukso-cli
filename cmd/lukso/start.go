package main

import (
	"fmt"
	"github.com/m8b-dev/lukso-cli/pid"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

func (dependency *ClientDependency) Start(
	arguments []string,
	ctx *cli.Context,
) (err error) {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependency.name)

	exists := pid.Exists(pidLocation)
	if exists {
		return errAlreadyRunning
	}

	command := exec.Command(dependency.name, arguments...)

	// since geth removed --logfile flag we have to manually adjust geth's stdout
	if dependency.name == gethDependencyName {
		var (
			logFile  *os.File
			fullPath string
		)

		gethLogDir := ctx.String(logFolderFlag)
		if gethLogDir == "" {
			return errFlagMissing
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
	if err != nil {
		return
	}

	err = pid.Create(pidLocation, command.Process.Pid)

	return
}

func (dependency *ClientDependency) Stop() error {
	pidLocation := fmt.Sprintf("%s/%s.pid", pidFileDir, dependency.name)

	pidVal, err := pid.Load(pidLocation)
	if err != nil {
		log.Warnf("%s is not running - skipping...", dependency.name)

		return nil
	}

	err = pid.Kill(pidLocation, pidVal)
	if err != nil {
		return errProcessNotFound
	}

	return nil
}

func startClients(ctx *cli.Context) error {
	log.Info("Starting all clients")

	if ctx.Bool(validatorFlag) && ctx.String(transactionFeeRecipientFlag) == "" {
		log.Errorf("%s flag is required but wasn't provided", transactionFeeRecipientFlag)

		return errFlagMissing
	}

	err := startGeth(ctx)
	if err != nil {
		return err
	}

	err = startPrysm(ctx)
	if err != nil {
		return err
	}

	if ctx.Bool(validatorFlag) {
		err = startValidator(ctx)
	}

	return err
}

func startGeth(ctx *cli.Context) error {
	log.Info("Running geth init first...")

	err := initGeth(ctx)
	if err != nil {
		log.Errorf("There was an error while initalizing geth. Error: %v", err)
	}

	log.Info("Starting Geth")

	err = clientDependencies[gethDependencyName].Start(prepareGethStartFlags(ctx), ctx)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("Starting Prysm")

	err := clientDependencies[prysmDependencyName].Start(preparePrysmStartFlags(ctx), ctx)
	if err != nil {
		return err
	}

	log.Info("Prysm started! Use lukso logs command to see logs")
	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("Starting Validator")

	err := clientDependencies[validatorDependencyName].Start(prepareValidatorStartFlags(ctx), ctx)
	if err != nil {
		return err
	}

	log.Info("Validator started! Use lukso logs command to see logs")
	return nil
}

func stopClients(ctx *cli.Context) (err error) {
	stopConsensus := ctx.Bool(consensusFlag)
	stopExecution := ctx.Bool(executionFlag)
	stopValidator := ctx.Bool(validatorFlag)

	if !stopConsensus && !stopExecution && !stopValidator {
		// if none is given then we stop all
		stopConsensus = true
		stopExecution = true
		stopValidator = true
	}

	if stopExecution {
		err = stopClient(clientDependencies[gethDependencyName])(ctx)
		if err != nil {
			return err
		}
	}

	if stopConsensus {
		err = stopClient(clientDependencies[prysmDependencyName])(ctx)
		if err != nil {
			return err
		}
	}

	if stopValidator {
		err = stopClient(clientDependencies[validatorDependencyName])(ctx)
	}

	return err
}

func stopClient(dependency *ClientDependency) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Infof("Stopping %s", dependency.name)

		err := dependency.Stop()

		return err
	}
}

func initGeth(ctx *cli.Context) (err error) {
	dataDir := fmt.Sprintf("--datadir=%s", ctx.String(gethDatadirFlag))
	command := exec.Command("geth", "init", dataDir, ctx.String(genesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
