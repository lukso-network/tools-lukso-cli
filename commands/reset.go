package commands

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/dep/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func (c *commander) Reset(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return nil
	}
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	dataDir := configs.MainnetDatadir
	if ctx.Bool(flags.TestnetFlag) {
		dataDir = configs.TestnetDatadir
	}

	message := fmt.Sprintf("⚠️  WARNING: THIS ACTION WILL DELETE THE CLIENTS DATA DIRECTORIES ⚠️\n"+
		"\nDirectory that will be deleted:\n"+
		"- %s\n\nAre you sure you want to continue? [y/N]: ", dataDir)

	input := utils.RegisterInputWithMessage(message)
	if !strings.EqualFold(input, "y") {
		log.Info("❌  Aborting...")

		return nil
	}

	err = os.RemoveAll(dataDir)
	if err != nil {
		log.Info("🗑️  Data directories deleted.")
	}

	log.Info("🗑️  Data directories deleted.")

	return nil
}
