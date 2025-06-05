package api

import (
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
)

type HandlerFunc[Rq types.Request, Rs types.Response] func(Rq) Rs

// Handler unifies commands that are used in both CLI and TUI.
// It serves as a list of what's possible to do with the node, e.g. start/stop, log, manage your validator etc.
// Since CLI and TUI will require different stdouts, each command in case of logging should has an output channel attached.
// All handlers return types.Response instead of a specific struct like in request for better error handling.
type Handler interface {
	// Handlers
	Install(types.InstallRequest) types.InstallResponse
	Init(types.InitRequest) types.InitResponse
	Update(types.UpdateRequest) types.UpdateResponse
	Start(types.StartRequest) types.StartResponse
	Stop(types.StopRequest) types.StopResponse
	Status(types.StatusRequest) types.StatusResponse
	Logs(types.LogsRequest) types.LogsResponse
	Reset(types.ResetRequest) types.ResetResponse
	ValidatorImport(types.ValidatorImportRequest) types.ValidatorImportResponse
	ValidatorList(types.ValidatorListRequest) types.ValidatorListResponse
	ValidatorExit(types.ValidatorExitRequest) types.ValidatorExitResponse
	Version(types.VersionRequest) types.VersionResponse
	VersionClients(types.VersionClientsRequest) types.VersionClientsResponse
}

type handler struct {
	file      file.Manager
	log       logger.Logger
	installer installer.Installer
}

var _ Handler = &handler{}

func NewHandler(
	file file.Manager,
	logger logger.Logger,
	installer installer.Installer,
) Handler {
	return &handler{
		file,
		logger,
		installer,
	}
}
