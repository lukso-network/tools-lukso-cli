package api

import (
	"fmt"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
)

func (h *handler) Install(args types.InstallRequest) (resp types.InstallResponse) {
	if runningClients := clients.RunningClients(); runningClients != nil {
		return types.InstallResponse{
			Error: errors.ErrClientsAlreadyRunning{Clients: runningClients},
		}
	}

	if !h.cfg.Exists() {
		return types.InstallResponse{
			Error: errors.ErrCfgMissing,
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

		err := c.Install(v, false)
		if err != nil {
			h.log.Warn(fmt.Sprintf("Unable to download %s file: %v - continuing...", c.Name(), err))

			continue
		}
	}

	err = cfg.WriteExecution(selectedExecution.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing execution client: %v", err), 1)
	}

	err = cfg.WriteConsensus(selectedConsensus.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing consensus client: %v", err), 1)
	}

	err = cfg.WriteValidator(selectedValidator.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing validator client: %v", err), 1)
	}

	return
}
