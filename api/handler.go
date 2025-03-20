package api

import (
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
)

// Handler unifies commands that are used in both CLI and TUI.
// It serves as a list of what's possible to do with the node, e.g. start/stop, log, manage your validator etc.
// Since CLI and TUI will require different stdouts, each command in case of logging should has an output channel attached.
type Handler interface {
	Install(clients []clients.ClientBinaryDependency, versions map[string]common.Semver) error
	Init() error
	Update(clients []clients.ClientBinaryDependency, versions map[string]common.Semver) error
	Start(clients []clients.ClientBinaryDependency, args map[string][]string) error
	Stop(clients []clients.ClientBinaryDependency) error
	Status(clients []clients.ClientBinaryDependency) error
	Logs(clients clients.ClientBinaryDependency) error
	Reset() error
	ValidatorImport(validator clients.ValidatorBinaryDependency) error
	ValidatorList(validator clients.ValidatorBinaryDependency) error
	ValidatorExit(validator clients.ValidatorBinaryDependency) error
	Version() error
	VersionClients() error
}

type handler struct {
	cfg config.Config
}

var _ Handler = &handler{}
