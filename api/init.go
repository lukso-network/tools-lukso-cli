package api

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

func (h *handler) Init(args types.InitArgs) (err error) {
	if clients.IsAnyRunning() {
		return nil
	}

	if h.cfg.Exists() {
		message := "⚠️  This folder has already been initialized. Do you want to re-initialize it? Please note that configs in this folder will NOT be overwritten [y/N]:\n> "
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") {
			log.Info("Aborting...")

			return nil
		}
	}

	log.Info("⬇️  Downloading shared configuration files...")
	_ = installConfigGroup(configs.SharedConfigDependencies, false) // we can omit errors - all errors are catched by cli.Exit()
	log.Info("✅  Shared configuration files downloaded!\n\n")

	installClientConfigFiles(false)

	err := utils.CreateJwtSecret(jwtSecretPath)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while creating JWT secret file: %v", err), 1)
	}

	err = os.MkdirAll(pid.FileDir, common.ConfigPerms)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing PID directory: %v", err), 1)
	}

	switch h.cfg.Exists() {
	case true:
		log.Info("⚙️   LUKSO configuration already exists - continuing...")
	case false:
		log.Info("⚙️   LUKSO CLI can replace your p2p communication flags with your public IP for better connection (optional)")
		message := "Do you want to proceed? [Y/n]\n> "

		ip := "0.0.0.0"
		input := utils.RegisterInputWithMessage(message)
		if strings.EqualFold(input, "y") || input == "" {
			ip, err = setIPInConfigs()
			if err != nil {
				return err
			}
		}

		log.Info("⚙️   Creating LUKSO configuration file...")

		err = cfg.Create("", "", "", ip)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while preparing LUKSO configuration: %v", err), 1)
		}

		log.Infof("⚙️  IPv4 found: %s", ip)
		log.Infof("✅  LUKSO configuration created under %s", config.Path)
	}

	displayNetworksHardforkTimestamps()

	log.Info("✅  Working directory initialized! \n1. ⚙️  Use 'lukso install' to install clients. \n2. ▶️  Use 'lukso start' to start your node.")

	return nil
}
