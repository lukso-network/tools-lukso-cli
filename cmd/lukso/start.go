package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

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
	if dependency.name == gethDependencyName || dependency.name == erigonDependencyName {
		var (
			logFile  *os.File
			fullPath string
		)

		logFolder := ctx.String(logFolderFlag)
		if logFolder == "" {
			return cli.Exit(fmt.Sprintf("%v- %s", errFlagMissing, logFolderFlag), 1)
		}

		fullPath, err = prepareTimestampedFile(logFolder, dependency.name)
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

	time.Sleep(1 * time.Second)

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

	log.Infof("🛑  Stopped %s", dependency.name)

	return nil
}

func startClients(ctx *cli.Context) error {
	log.Info("🔎  Looking for client configuration file...")
	if !cfg.Exists() {
		return cli.Exit(folderNotInitialized, 1)
	}

	err := cfg.Read()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		return cli.Exit(selectedClientsNotFound, 1)
	}

	log.Info("🔄  Starting all clients")

	if ctx.Bool(validatorFlag) && ctx.String(transactionFeeRecipientFlag) == "" || ctx.Bool(transactionFeeRecipientFlag) { // this means that we provided flag without value
		return cli.Exit(fmt.Sprintf("❌  %s flag is required but wasn't provided", transactionFeeRecipientFlag), 1)
	}

	switch executionClient {
	case gethDependencyName:
		err = startGeth(ctx)
	case erigonDependencyName:
		err = startErigon(ctx)
	}
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", executionClient, err), 1)
	}

	switch consensusClient {
	case prysmDependencyName:
		err = startPrysm(ctx)
	case lighthouseDependencyName:
		err = startLighthouse(ctx)
	}
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", consensusClient, err), 1)
	}

	if ctx.Bool(validatorFlag) {
		err = startValidator(ctx)
	}

	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	log.Info("🎉  Clients have been started. Your node is now running 🆙.")

	return nil
}

func startGeth(ctx *cli.Context) error {
	log.Info("⚙️  Running geth init first...")

	err := initClient(gethDependencyName, ctx)
	if err != nil && !errors.Is(err, errAlreadyRunning) { // if it is already running it will be caught during start
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

	log.Info("✅  Geth started! Use 'lukso logs' to see the logs.")

	return nil
}

func startErigon(ctx *cli.Context) error {
	log.Info("⚙️  Running erigon init first...")

	err := initClient(erigonDependencyName, ctx)
	if err != nil && !errors.Is(err, errAlreadyRunning) { // if it is already running it will be caught during start
		log.Errorf("❌  There was an error while initalizing geth. Error: %v", err)

		return err
	}

	log.Info("🔄  Starting Erigon")

	erigonFlags, ok := prepareErigonStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err = clientDependencies[erigonDependencyName].Start(erigonFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Erigon started! Use 'lukso log' to see logs.")

	return nil
}

func startPrysm(ctx *cli.Context) error {
	log.Info("🔄  Starting Prysm")
	prysmFlags, err := preparePrysmStartFlags(ctx)
	if err != nil {
		return err
	}

	err = clientDependencies[prysmDependencyName].Start(prysmFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Prysm started! Use 'lukso logs' to see the logs.")

	return nil
}

func startLighthouse(ctx *cli.Context) error {
	log.Info("🔄  Starting Lighthouse")
	lighthouseFlags, err := prepareLighthouseStartFlags(ctx)
	if err != nil {
		return err
	}

	lighthouseFlags = append([]string{"beacon_node"}, lighthouseFlags...)

	err = clientDependencies[lighthouseDependencyName].Start(lighthouseFlags, ctx)
	if err != nil {
		return err
	}

	log.Info("✅  Lighthouse started! Use 'lukso log' to see logs.")

	return nil
}

func startValidator(ctx *cli.Context) error {
	log.Info("🔄  Starting Validator")
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
		return err
	}

	log.Info("✅  Validator started! Use 'lukso logs' to see the logs.")

	return nil
}

func stopClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return cli.Exit(folderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		return cli.Exit(selectedClientsNotFound, 1)
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

		err = stopClient(clientDependencies[executionClient])
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was an error while stopping geth: %v", err), 1)
		}
	}

	if stopConsensus {
		log.Infof("⚙️  Stopping consensus [%s]", consensusClient)

		err = stopClient(clientDependencies[consensusClient])
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was an error while stopping prysm: %v", err), 1)
		}
	}

	if stopValidator {
		log.Info("⚙️  Stopping validator")

		err = stopClient(clientDependencies[validatorDependencyName])
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was an error while stopping validator: %v", err), 1)
		}
	}

	return nil
}

func stopClient(dependency *ClientDependency) error {
	err := dependency.Stop()

	return err
}

func initClient(client string, ctx *cli.Context) (err error) {
	if isRunning(client) {
		return errAlreadyRunning
	}

	if !flagFileExists(ctx, genesisJsonFlag) {
		return errors.New("❌  Genesis JSON not found")
	}

	var dataDir string
	switch client {
	case gethDependencyName:
		dataDir = fmt.Sprintf("--datadi r=%s", ctx.String(gethDatadirFlag))
	case erigonDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(erigonDatadirFlag))
	}

	command := exec.Command(client, "init", dataDir, ctx.String(genesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
