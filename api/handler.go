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
	Install(types.InstallArgs) types.Response
	Init(types.InitArgs) types.Response
	Update(types.UpdateArgs) types.Response
	Start(types.StartArgs) types.Response
	Stop(types.StopArgs) types.Response
	Status(types.StatusArgs) types.Response
	Logs(types.LogsArgs) types.Response
	Reset(types.ResetArgs) types.Response
	ValidatorImport(types.ValidatorImportArgs) types.Response
	ValidatorList(types.ValidatorListArgs) types.Response
	ValidatorExit(types.ValidatorExitArgs) types.Response
	Version(types.VersionArgs) types.Response
	VersionClients(types.VersionClientsArgs) types.Response
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
