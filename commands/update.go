package commands

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func UpdateClients(ctx *cli.Context) (err error) {
	isRoot, err := system.IsRoot()
	if err != nil {
		return utils.Exit(fmt.Sprintf("There was an error while checking user privileges: %v", err), 1)
	}
	if !isRoot {
		return utils.Exit(errors.ErrNeedRoot.Error(), 1)
	}

	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return cli.Exit(errors.FolderNotInitialized, 1)
	}
	err = cfg.Read()
	if err != nil {
		return cli.Exit(fmt.Sprintf("There was an error reading config file: %v", err), 1)
	}

	execution, ok1 := clients.AllClients[cfg.Execution()]
	consensus, ok2 := clients.AllClients[cfg.Consensus()]
	validator, ok3 := clients.AllClients[cfg.Validator()]

	if !ok1 || !ok2 || !ok3 {
		return cli.Exit(errors.ErrClientNotSupported, 1)
	}

	toUpdate := []clients.ClientBinaryDependency{execution, consensus, validator}

	for _, client := range toUpdate {
		if client == clients.LighthouseValidator || client == clients.TekuValidator || client == clients.Nimbus2Validator {
			continue
		}

		err = client.Update()
		if err != nil {
			return cli.Exit(fmt.Sprintf("There was an error while updating %s: %v", client.Name(), err), 1)
		}
	}

	return
}

func UpdateConfigs(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return cli.Exit(errors.FolderNotInitialized, 1)
	}

	log.Info("⬇️  Updating IPv4 Addres...")

	ip, err := getPublicIP()
	if err != nil {
		log.Warn("⚠️  There was an error while getting the IPv4 Address: continuing...")
	} else {
		err = cfg.WriteIPv4(ip)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while writing IPv4 Address: %v", err), 1)
		}
	}

	if ctx.Bool(flags.AllFlag) {
		message := "⚠️  Warning - this action will overwrite all of your client configuration files, are you sure you want to continue? [y/N]\n> "

		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") {
			log.Info("❌  Aborting...")

			return nil
		}

		installClientConfigFiles(true)
	}

	log.Info("⬇️  Updating shared configuration files...")
	_ = installConfigGroup(configs.UpdateConfigDependencies, true)
	log.Info("✅  Shared configuration files updated!\n\n")

	displayNetworksHardforkTimestamps()

	return
}
