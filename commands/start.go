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

func StartClients(ctx *cli.Context) (err error) {
	log.Info("🔎  Looking for client configuration file...")
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	selectedExecution := cfg.Execution()
	selectedConsensus := cfg.Consensus()
	if selectedExecution == "" || selectedConsensus == "" {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	log.Info("🔄  Starting all clients")

	if ctx.Bool(flags.ValidatorFlag) && ctx.String(flags.TransactionFeeRecipientFlag) == "" || ctx.Bool(flags.TransactionFeeRecipientFlag) { // this means that we provided flag without value
		return utils.Exit(fmt.Sprintf("❌  %s flag is required but wasn't provided", flags.TransactionFeeRecipientFlag), 1)
	}

	executionClient, ok1 := clients.AllClients[selectedExecution]
	consensusClient, ok2 := clients.AllClients[selectedConsensus]

	if !ok1 || !ok2 {
		return utils.Exit("❌  Client found in LUKSO configuration file is not supported", 1)
	}

	execArgs, err := executionClient.PrepareStartFlags(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing %s flags: %v", executionClient.Name(), err), 1)
	}
	consArgs, err := consensusClient.PrepareStartFlags(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing %s flags: %v", consensusClient.Name(), err), 1)
	}

	err = executionClient.Start(ctx, execArgs)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", executionClient.Name(), err), 1)
	}
	err = consensusClient.Start(ctx, consArgs)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", executionClient.Name(), err), 1)
	}

	if ctx.Bool(flags.ValidatorFlag) {
		err = startValidator(ctx)
	}
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	log.Info("🎉  Clients have been started. Checking status:")
	log.Info("")

	_ = StatClients(ctx)

	log.Info("")
	log.Info("If execution and consensus clients are Running 🟢, your node is now 🆙.")
	log.Info("👉 Please check the logs with 'lukso logs' to make sure everything is running correctly.")

	return nil
}

func startValidator(ctx *cli.Context) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
	}

	consensusClient, ok := clients.AllClients[cfg.Consensus()]
	if !ok {
		return utils.Exit(errors.ErrClientNotSupported.Error(), 1)
	}

	err = startPrysmValidator(ctx)

	if err != nil {
		return exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	return
}
