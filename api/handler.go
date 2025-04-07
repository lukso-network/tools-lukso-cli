package api

import (
	"github.com/lukso-network/tools-lukso-cli/api/logger"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/config"
)

// Handler unifies commands that are used in both CLI and TUI.
// It serves as a list of what's possible to do with the node, e.g. start/stop, log, manage your validator etc.
// Since CLI and TUI will require different stdouts, each command in case of logging should has an output channel attached.
// All handlers return types.Response instead of a specific struct like in request for better error handling.
type Handler interface {
	Install(types.InstallArgs) types.InstallResponse
	Init(types.InitArgs) types.InitResponse
	Update(types.UpdateArgs) types.UpdateResponse
	Start(types.StartArgs) types.StartResponse
	Stop(types.StopArgs) types.StopResponse
	Status(types.StatusArgs) types.StatusResponse
	Logs(types.LogsArgs) types.LogsResponse
	Reset(types.ResetArgs) types.ResetResponse
	ValidatorImport(types.ValidatorImportArgs) types.ValidatorImportResponse
	ValidatorList(types.ValidatorListArgs) types.ValidatorListResponse
	ValidatorExit(types.ValidatorExitArgs) types.ValidatorExitResponse
	Version(types.VersionArgs) types.VersionResponse
	VersionClients(types.VersionClientsArgs) types.VersionClientsResponse
}

type handler struct {
	cfg       *config.Config
	file      file.Manager
	log       logger.Logger
	installer installer.Installer
}

var _ Handler = &handler{}

func NewHandler() Handler {
	log := logger.ConsoleLogger{}
	fmng := file.NewManager()
	inst := installer.NewInstaller(fmng)

	return &handler{
		cfg:       config.NewConfig(config.Path),
		file:      fmng,
		log:       log,
		installer: inst,
	}
}
