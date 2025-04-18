package commands

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api"
	"github.com/lukso-network/tools-lukso-cli/common/display"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
)

// Commander is an interface holding commands that can be executed via CLI
type Commander interface {
	// Setup/utilities
	Before(ctx *cli.Context) error
	After(ctx *cli.Context) error

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
	display display.Display
	log     logger.Logger
}

var _ Commander = &commander{}

func NewCommander(handler api.Handler, display display.Display, log logger.Logger) Commander {
	return &commander{
		handler: handler,
		display: display,
		log:     log,
	}
}

// TODO: remove sleeps and rely on other form of completion assertion.
func (c *commander) Before(ctx *cli.Context) error {
	go c.display.Listen()
	return nil
}

func (c *commander) After(ctx *cli.Context) error {
	c.log.Close()
	time.Sleep(time.Millisecond * 50)
	return nil
}
