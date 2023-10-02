package commands

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/utils"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

func StatClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return cli.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		log.Errorf("❌  There was an error while reading configuration file: %v", err)

		return nil
	}

	selectedExecution := cfg.Execution()
	selectedConsensus := cfg.Consensus()
	validatorDependencyName := cfg.Validator()

	err = statClient(selectedExecution, "Execution")(ctx)
	if err != nil {
		return
	}

	err = statClient(selectedConsensus, "Consensus")(ctx)
	if err != nil {
		return
	}

	err = statClient(validatorDependencyName, "Validator")(ctx)
	if err != nil {
		return
	}

	return
}

func statClient(dependencyName, layer string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if dependencyName == "" {
			dependencyName = "none"

			log.Warnf("PID ----- - %s (%s): Stopped 🔘", layer, dependencyName)

			return nil
		}

		client, ok := clients.AllClients[dependencyName]
		if !ok {
			return cli.Exit(errors.ErrClientNotSupported, 1)
		}

		isRunning := false
		if dependencyName != "" {
			isRunning = client.IsRunning()
		}

		if isRunning {
			pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())

			// since we now that process is, in fact, running (from previous checks) this shouldn't fail
			pidVal, err := pid.Load(pidLocation)
			if err != nil {
				return errors.ErrProcessNotFound
			}

			log.Infof("PID %d - %s (%s): Running 🟢", pidVal, layer, dependencyName)

			return nil
		}

		log.Warnf("PID ----- - %s (%s): Stopped 🔘", layer, dependencyName)

		return nil
	}
}

func DisplayPeers(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	executionClient := cfg.Execution()
	consensusClient := cfg.Consensus()

	clientPeers(ctx, executionClient, "Execution")
	clientPeers(ctx, consensusClient, "Consensus")

	return
}

func clientPeers(ctx *cli.Context, clientName, layer string) {
	if clientName == "" {
		log.Warnf("%s (none): Stopped 🔘", layer)

		return
	}

	client, ok := clients.AllClients[clientName]
	if !ok {
		log.Error(errors.ErrClientNotSupported)

		return
	}

	if !client.IsRunning() {
		log.Warnf("%s (%s): Stopped 🔘", layer, client.Name())

		return
	}

	outbound, inbound, err := client.Peers(ctx)
	if err != nil {
		log.Warnf("⚠ ️Unable to get peers from %s - please make sure that you have the right API enabled in your client configuration", client.Name())

		return
	}

	log.Infof("%s (%s): Outbound: %d | Inbound: %d 🟢", layer, client.Name(), outbound, inbound)

	return
}
