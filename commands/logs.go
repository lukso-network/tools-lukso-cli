package commands

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func (c *commander) Logs(ctx *cli.Context) error {
	log.Info("⚠️  Please specify your client - run 'lukso logs --help' for more info")

	return nil
}

func (c *commander) LogsLayer(layer string) func(*cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		logFileDir := ctx.String(flags.LogFolderFlag)
		if logFileDir == "" {
			return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
		}

		if !cfg.Exists() {
			return utils.Exit(errors.FolderNotInitialized, 1)
		}

		err = cfg.Read()
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while reading the configuration file: %v", err), 1)
		}

		var dependencyName string

		switch layer {
		case configs.ExecutionLayer:
			dependencyName = cfg.Execution()
		case configs.ConsensusLayer:
			dependencyName = cfg.Consensus()
		case configs.ValidatorLayer:
			dependencyName = cfg.Validator()
		default:
			return utils.Exit("❌  Unexpected error: unknown layer logged", 1)
		}

		client, ok := clients.AllClients[dependencyName]
		if !ok {
			return utils.Exit(errors.ErrClientNotSupported.Error(), 1)
		}

		latestFile, err := utils.GetLastFile(logFileDir, client.CommandName())
		if latestFile == "" && err == nil {
			return nil
		}
		if err != nil {
			return utils.Exit(fmt.Sprintf("There was an error while getting the latest log file: %v", err), 1)
		}

		log.Infof("Selected layer: %s", layer)
		log.Infof("Selected client: %s", dependencyName)

		message := fmt.Sprintf("(NOTE: PATH TO FILE: %s)\nDo you want to show the log file %s?"+
			" Doing so can print lots of text on your screen [Y/n]: ", logFileDir+"/"+latestFile, latestFile)

		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("❌  Aborting...") // there is no error, just a possible change of mind - shouldn't return err

			return nil
		}

		return client.Logs(logFileDir + "/" + latestFile)
	}
}
