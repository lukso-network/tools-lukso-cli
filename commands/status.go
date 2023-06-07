package commands

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/pid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
