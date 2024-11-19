package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
)

func ClientVersions(ctx *cli.Context) (err error) {
	log.Info("⬇️  Getting client versions - this may take a few seconds...")
	clientVersions := make(map[string]string)

	for _, clName := range clients.AllClientNames {
		cl := clients.AllClients[clName]
		clientVersions[clName] = cl.Version()
	}

	for _, clName := range clients.AllClientNames {
		ver := clientVersions[clName]
		if ver == clients.VersionNotAvailable {
			log.Warnf("%s: Not available", clName)
		} else {
			log.Infof("%s: %s", clName, ver)
		}
	}

	return
}
