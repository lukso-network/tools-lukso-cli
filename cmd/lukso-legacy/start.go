package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/lukso-network/tools-lukso-cli/pid"
	"github.com/urfave/cli/v2"
)

func (dependency *ClientDependency) Start(
	arguments []string,
	ctx *cli.Context,
) (err error) {
	if isRunning(dependency.name) {
		log.Infof("🔄️  %s is already running - stopping first...", dependency.name)

		err = dependency.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", dependency.name)
	}

	command := exec.Command(dependency.name, arguments...)

	if dependency.name == gethDependencyName || dependency.name == erigonDependencyName {
		log.Infof("⚙️  Running %s init...", dependency.name)

		err = initClient(dependency.name, ctx)
		if err != nil && !errors.Is(err, errAlreadyRunning) { // if it is already running it will be caught during start
			log.Errorf("❌  There was an error while initalizing %s. Error: %v", dependency.name, err)

			return err
		}

		var (
			logFile  *os.File
			fullPath string
		)

		logFolder := ctx.String(logFolderFlag)
		if logFolder == "" {
			return exit(fmt.Sprintf("%v- %s", errFlagMissing, logFolderFlag), 1)
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

	log.Infof("🔄  Starting %s", dependency.name)
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

	return nil
}

func startClients(ctx *cli.Context) error {
	log.Info("🔎  Looking for client configuration file...")
	if !cfg.Exists() {
		return exit(folderNotInitialized, 1)
	}

	err := cfg.Read()
	if err != nil {
		return exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		return exit(selectedClientsNotFound, 1)
	}

	log.Info("🔄  Starting all clients")

	if ctx.Bool(validatorFlag) && ctx.String(transactionFeeRecipientFlag) == "" || ctx.Bool(transactionFeeRecipientFlag) { // this means that we provided flag without value
		return exit(fmt.Sprintf("❌  %s flag is required but wasn't provided", transactionFeeRecipientFlag), 1)
	}

	switch executionClient {
	case gethDependencyName:
		err = startGeth(ctx)
	case erigonDependencyName:
		err = startErigon(ctx)
	}
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting %s: %v", executionClient, err), 1)
	}

	switch consensusClient {
	case prysmDependencyName:
		err = startPrysm(ctx)
	case lighthouseDependencyName:
		err = startLighthouse(ctx)
	}
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting %s: %v", consensusClient, err), 1)
	}

	if ctx.Bool(validatorFlag) {
		err = startValidator(ctx)
	}

	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	log.Info("🎉  Clients have been started. Checking status:")
	log.Info("")

	_ = statClients(ctx)

	log.Info("")
	log.Info("If execution and consensus clients are Running 🟢, your node is now 🆙.")
	log.Info("👉 Please check the logs with 'lukso logs' to make sure everything is running correctly.")

	return nil
}

func startGeth(ctx *cli.Context) error {
	gethFlags, ok := prepareGethStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err := clientDependencies[gethDependencyName].Start(gethFlags, ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting geth: %v", err), 1)
	}

	log.Info("✅  Geth started! Use 'lukso logs' to see the logs.")

	return nil
}

func startErigon(ctx *cli.Context) error {
	erigonFlags, ok := prepareErigonStartFlags(ctx)
	if !ok {
		return errFlagPathInvalid
	}

	err := clientDependencies[erigonDependencyName].Start(erigonFlags, ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting erigon: %v", err), 1)
	}

	log.Info("✅  Erigon started! Use 'lukso log' to see logs.")

	return nil
}

func startPrysm(ctx *cli.Context) error {
	prysmFlags, err := preparePrysmStartFlags(ctx)
	if err != nil {
		return err
	}

	err = clientDependencies[prysmDependencyName].Start(prysmFlags, ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting prysm: %v", err), 1)
	}

	log.Info("✅  Prysm started! Use 'lukso logs' to see the logs.")

	return nil
}

func startLighthouse(ctx *cli.Context) error {
	lighthouseFlags, err := prepareLighthouseStartFlags(ctx)
	if err != nil {
		return err
	}

	lighthouseFlags = append([]string{"beacon_node"}, lighthouseFlags...)

	err = clientDependencies[lighthouseDependencyName].Start(lighthouseFlags, ctx)
	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting lighthouse: %v", err), 1)
	}

	log.Info("✅  Lighthouse started! Use 'lukso log' to see logs.")

	return nil
}

func stopClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return exit(folderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()
	if executionClient == "" || consensusClient == "" {
		return exit(selectedClientsNotFound, 1)
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
			return exit(fmt.Sprintf("❌  There was an error while stopping geth: %v", err), 1)
		}
	}

	if stopConsensus {
		log.Infof("⚙️  Stopping consensus [%s]", consensusClient)

		err = stopClient(clientDependencies[consensusClient])
		if err != nil {
			return exit(fmt.Sprintf("❌  There was an error while stopping prysm: %v", err), 1)
		}
	}

	if stopValidator {
		log.Info("⚙️  Stopping validator")

		err = stopClient(clientDependencies[validatorDependencyName])
		if err != nil {
			return exit(fmt.Sprintf("❌  There was an error while stopping validator: %v", err), 1)
		}
	}

	return nil
}

func stopClient(dependency *ClientDependency) error {
	err := dependency.Stop()
	if err != nil {
		return err
	}

	log.Infof("🛑  Stopped %s", dependency.name)

	return nil
}

func initClient(client string, ctx *cli.Context) (err error) {
	if isRunning(client) {
		return errAlreadyRunning
	}

	if !fileExists(ctx.String(genesisJsonFlag)) {
		if ctx.Bool(testnetFlag) || ctx.Bool(devnetFlag) {
			return errors.New("❌  Genesis JSON not found")
		}

		message := `Choose your preferred initial LYX supply!
If you are a Genesis Validator, we recommend to choose the supply, which the majority of the Genesis Validators would choose,
to prevent your node from running on a network that can not finalize due to missing validators!
🗳️ See the voting results at https://deposit.mainnet.lukso.network

For more information read:
👉 https://medium.com/lukso/genesis-validators-deposit-smart-contract-freeze-and-testnet-launch-c5f7b568b1fc

Which initial LYX supply do you choose?
1: 35M LYX
2: 42M LYX
3: 100M LYX
> `
		var input string
		for input != "1" && input != "2" && input != "3" {
			input = registerInputWithMessage(message)
			switch input {
			case "1":
				err = ctx.Set(genesisJsonFlag, mainnetConfig+"/"+genesis35JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(genesisStateFlag, mainnetConfig+"/"+genesisState35FilePath)

			case "2":
				err = ctx.Set(genesisJsonFlag, mainnetConfig+"/"+genesis42JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(genesisStateFlag, mainnetConfig+"/"+genesisState42FilePath)

			case "3":
				err = ctx.Set(genesisJsonFlag, mainnetConfig+"/"+genesis100JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(genesisStateFlag, mainnetConfig+"/"+genesisState100FilePath)

			default:
				log.Warn("Please select a valid option\n\n")
			}

			if err != nil {
				return
			}
		}
	}

	var dataDir string
	switch client {
	case gethDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(gethDatadirFlag))
	case erigonDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(erigonDatadirFlag))
	}

	command := exec.Command(client, "init", dataDir, ctx.String(genesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
