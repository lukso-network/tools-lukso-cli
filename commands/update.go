package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
)

func UpdateClients(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return cli.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return cli.Exit(fmt.Sprintf("There was an error reading config file: %v", err), 1)
	}

	execution, ok1 := clients.AllClients[cfg.Execution()]
	consensus, ok2 := clients.AllClients[cfg.Consensus()]
	validator, ok3 := clients.AllClients[cfg.Validator()]

	if !ok1 || !ok2 || !ok3 {
		return cli.Exit(errors.ErrClientNotSupported, 1)
	}

	toUpdate := []clients.ClientBinaryDependency{execution, consensus, validator}

	for _, client := range toUpdate {
		if client == clients.LighthouseValidator || client == clients.TekuValidator {
			continue
		}

		err = client.Update()
		if err != nil {
			return cli.Exit(fmt.Sprintf("There was an error while updating %s: %v", client.Name(), err), 1)
		}
	}

	return
}
