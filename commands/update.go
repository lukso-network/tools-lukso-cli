package commands

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/urfave/cli/v2"
)

func UpdateClients(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	for _, client := range clients.AllClients {
		if client == clients.LighthouseValidator {
			continue
		}

		err = client.Update()
		if err != nil {
			return cli.Exit(fmt.Sprintf("There was an error while updating %s: %v", client.Name(), err), 1)
		}
	}

	return
}
