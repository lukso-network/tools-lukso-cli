package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func LogClients(ctx *cli.Context) error {
	log.Info("⚠️  Please specify your client - run 'lukso logs --help' for more info")

	return nil
}

func LogLayer(layer string) func(*cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		return
	}
}
