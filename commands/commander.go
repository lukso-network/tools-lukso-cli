package commands

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api"
)

// Commander is an interface holding commands that can be executed via CLI
type Commander interface {
	Install(ctx *cli.Context) error
	Init(ctx *cli.Context) error
	Update(ctx *cli.Context) error
	UpdateConfigs(ctx *cli.Context) error
	Start(ctx *cli.Context) error
	Stop(ctx *cli.Context) error
	Status(ctx *cli.Context) error
	StatusPeers(ctx *cli.Context) error
	Logs(ctx *cli.Context) error
	LogsLayer(layer string) func(*cli.Context) error
	Reset(ctx *cli.Context) error
	ValidatorImport(ctx *cli.Context) error
	ValidatorList(ctx *cli.Context) error
	ValidatorExit(ctx *cli.Context) error
	Ui(ctx *cli.Context) error
	Version(version string) func(ctx *cli.Context) error
	VersionClients(ctx *cli.Context) error
}

type commander struct {
	handler api.Handler
}

var _ Commander = &commander{}

func NewCommander() Commander {
	return &commander{}
}
