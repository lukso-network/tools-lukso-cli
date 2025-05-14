package api

import (
	"github.com/lukso-network/tools-lukso-cli/api/types"
)

func (h *handler) Install(args types.InstallRequest) (resp types.InstallResponse) {
	// if runningClients := clients.RunningClients(); runningClients != nil {
	// 	return types.InstallResponse{
	// 		Error: errors.ErrClientsAlreadyRunning{Clients: runningClients},
	// 	}
	// }
	//
	// if !h.cfg.Exists() {
	// 	return types.InstallResponse{
	// 		Error: errors.ErrCfgMissing,
	// 	}
	// }
	//
	// selectedClients := make([]clients.Client, 0)
	//
	// for _, reqClient := range args.Clients {
	// 	c, ok := clients.AllClients[reqClient.Name]
	// 	if !ok {
	// 		continue
	// 	}
	//
	// 	selectedClients = append(selectedClients, c)
	// }
	//
	// for _, selectedClient := range selectedClients {
	// 	err := selectedClient.Install("", false)
	// 	if err != nil {
	// 		h.log.Warn(fmt.Sprintf("Unable to download %s file: %v - continuing...", selectedClient.Name(), err))
	// 	}
	// }
	//
	// return
	return
}
