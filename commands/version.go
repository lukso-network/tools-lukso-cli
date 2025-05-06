package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
)

func (c *commander) Version(version string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		fmt.Println("Version:", version)
		fmt.Println("To display versions of installed clients, run 'lukso version clients' command.")

		return nil
	}
}

func (c *commander) VersionClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	log.Info("⬇️  Getting client versions - this may take a few seconds...")

	padding := utils.MaxLength(clients.AllClientNames) + 1
	clientVersions := make(map[string]string)

	for _, clName := range clients.AllClientNames {
		cl := clients.AllClients[clName]
		clientVersions[clName] = cl.Version()
	}

	for _, clName := range clients.AllClientNames {
		ver := clientVersions[clName]
		if ver == clients.VersionNotAvailable {
			msg := fmt.Sprintf("%-*s| Not Available", padding, clName)
			log.Warn(msg)
		} else {
			msg := fmt.Sprintf("%-*s| %s", padding, clName, ver)
			log.Info(msg)
		}
	}

	return
}
