package api

import (
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/config"
)

// Handler unifies commands that are used in both CLI and TUI.
// It serves as a list of what's possible to do with the node, e.g. start/stop, log, manage your validator etc.
// Since CLI and TUI will require different stdouts, each command in case of logging should has an output channel attached.
type Handler interface {
	Install(types.InstallArgs) error
	Init(types.InitArgs) error
	Update(types.UpdateArgs) error
	Start(types.StartArgs) error
	Stop(types.StopArgs) error
	Status(types.StatusArgs) error
	Logs(types.LogsArgs) error
	Reset(types.ResetArgs) error
	ValidatorImport(types.ValidatorImportArgs) error
	ValidatorList(types.ValidatorListArgs) error
	ValidatorExit(types.ValidatorExitArgs) error
	Version(types.VersionArgs) error
	VersionClients(types.VersionClientsArgs) error
}

type handler struct {
	cfg config.Config
}

var _ Handler = &handler{}
