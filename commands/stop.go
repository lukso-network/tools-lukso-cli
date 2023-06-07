package commands

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func StopClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	selectedExecution := cfg.Execution()
	selectedConsensus := cfg.Consensus()
	selectedValidator := cfg.Validator()
	if selectedExecution == "" || selectedConsensus == "" || selectedValidator == "" {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	executionClient := clients.AllClients[selectedExecution]
	consensusClient := clients.AllClients[selectedConsensus]
	validatorClient := clients.AllClients[selectedValidator]

	stopConsensus := ctx.Bool(flags.ConsensusFlag)
	stopExecution := ctx.Bool(flags.ExecutionFlag)
	stopValidator := ctx.Bool(flags.ValidatorFlag)

	if !stopConsensus && !stopExecution && !stopValidator {
		// if none is given then we stop all
		stopConsensus = true
		stopExecution = true
		stopValidator = true
	}

	if stopExecution {
		log.Infof("⚙️  Stopping execution [%s]", selectedExecution)

		err = stopClient(executionClient)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while stopping geth: %v", err), 1)
		}
	}

	if stopConsensus {
		log.Infof("⚙️  Stopping consensus [%s]", selectedConsensus)

		err = stopClient(consensusClient)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while stopping prysm: %v", err), 1)
		}
	}

	if stopValidator {
		log.Infof("⚙️  Stopping validator [%s]", selectedValidator)

		err = stopClient(validatorClient)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while stopping validator: %v", err), 1)
		}
	}

	return nil
}

func stopClient(client clients.ClientBinaryDependency) error {
	err := client.Stop()
	if err != nil {
		return err
	}

	log.Infof("🛑  Stopped %s", client.Name())

	return nil
}
