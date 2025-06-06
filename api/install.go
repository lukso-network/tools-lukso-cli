package api

import (
	"fmt"

	apierrors "github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
)

func (h *handler) Install(args types.InstallRequest) (resp types.InstallResponse) {
	if runningClients := clients.RunningClients(); runningClients != nil {
		return types.InstallResponse{
			Error: apierrors.ErrClientsAlreadyRunning{Clients: runningClients},
		}
	}

	if !config.Exists() {
		return types.InstallResponse{
			Error: apierrors.ErrCfgMissing,
		}
	}

	for _, selectedClient := range args.Clients {
		n := selectedClient.Name
		v := selectedClient.Version

		c, ok := clients.AllClients[n]
		if !ok {
			h.log.Warn(fmt.Sprintf("Client %s not supported - skipping...", c.Name()))

			continue
		}

		h.log.Infof("Installing %s...", c.Name())
		err := c.Install(v, false)
		if err != nil {
			h.log.Warn(fmt.Sprintf("Unable to download %s file: %v - continuing...", c.Name(), err))

			continue
		}
	}

	config.Set("useclients.execution", args.Clients[0].Name)
	config.Set("useclients.consensus", args.Clients[1].Name)
	config.Set("useclients.validator", args.Clients[2].Name)

	err := config.Write()
	if err != nil {
		return types.InstallResponse{
			Error: fmt.Errorf("unable to write the config: %v", err),
		}
	}

	return
}
