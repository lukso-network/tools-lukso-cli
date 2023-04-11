package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/m8b-dev/lukso-cli/pid"
	"github.com/urfave/cli/v2"
)

func (dependency *ClientDependency) Start(
	arguments []string,
	ctx *cli.Context,
) (err error) {
	if isRunning(dependency.name) {
		log.Infof("⏭️  %s is already running - skipping...", dependency.name)

		return nil
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

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependency.name)
	err = pid.Create(pidLocation, command.Process.Pid)

	return
}

func (dependency *ClientDependency) Stop() error {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependency.name)

	pidVal, err := pid.Load(pidLocation)
	if err != nil {
		log.Warnf("⏭️  %s is not running - skipping...", dependency.name)

		return nil
	}

	err = pid.Kill(pidLocation, pidVal)
	if err != nil {
		return errProcessNotFound
	}

	log.Infof("🛑 Stopped %s", dependency.name)

	return nil
}

func startClients(ctx *cli.Context) error {
	log.Info("🔎  Looking for client configuration file...")
	if !cfg.Exists() {
		log.Error(folderNotInitialized)

		return nil
	}

	err := cfg.Read()
	if err != nil {
		log.Errorf("❌  Couldn't read from config file: %v", err)

		return err
	}

	// TODO for now just check if installed - when multiple clients will be supported we can run it generically
	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		log.Error(selectedClientsNotFound)

		return nil
	}

	log.Info("🔄  Starting all clients")

	if ctx.Bool(validatorFlag) && ctx.String(transactionFeeRecipientFlag) == "" {
		log.Errorf("❌  %s flag is required but wasn't provided", transactionFeeRecipientFlag)

		return errFlagMissing
	}

	err = startGeth(ctx)
	if err != nil {
		log.Errorf("❌  There was an error while starting geth: %v", err)

		return nil
	}

	err = startPrysm(ctx)
	if err != nil {
		log.Errorf("❌  There was an error while starting prysm: %v", err)

		return nil
	}

	if ctx.Bool(validatorFlag) {
		err = startValidator(ctx)
	}

	if err != nil {
		log.Errorf("❌  There was an error while starting validator: %v", err)

		return nil
	}

	log.Info("🎉  Clients have been started. Your node is now running 🆙.")

	return nil
}

func startGeth(ctx *cli.Context) error {
	log.Info("⚙️  Running geth init first...")

	err := initGeth(ctx)
	if err != nil {
		log.Errorf("❌  There was an error while initalizing geth. Error: %v", err)

		return err
	}

	log.Info("🔄  Starting Geth")
	gethFlags, ok := prepareGethStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err = clientDependencies[gethDependencyName].Start(gethFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Geth started! Use 'lukso log' to see logs.")

	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("🔄  Starting Prysm")
	prysmFlags, ok := preparePrysmStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err := clientDependencies[prysmDependencyName].Start(prysmFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Prysm started! Use 'lukso log' to see logs.")

	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("🔄  Starting Validator")
	validatorFlags, ok := prepareValidatorStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err := clientDependencies[validatorDependencyName].Start(validatorFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Validator started! Use 'lukso log' to see logs.")

	return nil
}

func stopClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		log.Error(folderNotInitialized)

		return
	}

	err = cfg.Read()
	if err != nil {
		log.Errorf("❌  Couldn't read from config: %v", err)

		return nil
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		log.Error(selectedClientsNotFound)

		return nil
	}

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
		log.Infof("⚙️  Stopping execution [%s]", executionClient)

		err = stopClient(clientDependencies[gethDependencyName])
		if err != nil {
			return err
		}
	}

	if stopConsensus {
		log.Infof("⚙️  Stopping consensus [%s]", consensusClient)

		err = stopClient(clientDependencies[prysmDependencyName])
		if err != nil {
			return err
		}
	}

	if stopValidator {
		log.Info("⚙️  Stopping validator")

		err = stopClient(clientDependencies[validatorDependencyName])
	}

	return err
}

func stopClient(dependency *ClientDependency) error {
	err := dependency.Stop()

	return err
}

func initGeth(ctx *cli.Context) (err error) {
	if isRunning(gethDependencyName) {
		return errAlreadyRunning
	}

	if !flagFileExists(ctx, genesisJsonFlag) {
		return errors.New("❌  Genesis JSON not found")
	}

	dataDir := fmt.Sprintf("--datadir=%s", ctx.String(gethDatadirFlag))
	command := exec.Command("geth", "init", dataDir, ctx.String(genesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
