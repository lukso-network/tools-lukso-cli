package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/m8b-dev/lukso-cli/config"
	"github.com/m8b-dev/lukso-cli/pid"
)

const jwtSecretPath = configsRootDir + "/shared/secrets/jwt.hex"

// initializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func initializeDirectory(ctx *cli.Context) error {
	if isAnyRunning() {
		return nil
	}

	if cfg.Exists() {
		message := "⚠️  This folder has already been initialized. Do you want to re-initialize it? Please note that configs in this folder will NOT be overwritten [y/N]:\n> "
		input := registerInputWithMessage(message)
		if !strings.EqualFold(input, "y") {
			log.Info("Aborting...")

			return nil
		}
	}

	log.Info("⬇️  Downloading shared configuration files...")
	_ = initConfigGroup(sharedConfigDependencies) // we can omit errors - all errors are catched by cli.Exit()

	log.Info("⬇️  Downloading geth configuration files...")
	_ = initConfigGroup(gethConfigDependencies)

	log.Info("⬇️  Downloading erigon configuration files...")
	_ = initConfigGroup(erigonConfigDependencies)

	log.Info("⬇️  Downloading prysm configuration files...")
	_ = initConfigGroup(prysmConfigDependencies)

	log.Info("⬇️  Downloading lighthouse configuration files...")
	_ = initConfigGroup(lighthouseConfigDependencies)

	log.Info("⬇️  Downloading prysm validator configuration files...")
	_ = initConfigGroup(validatorConfigDependencies)

	err := createJwtSecret(jwtSecretPath)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while creating JWT secret file: %v", err), 1)
	}

	err = os.MkdirAll(pid.FileDir, configPerms)
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌  There was an error while preparing PID directory: %v", err), 1)
	}

	switch cfg.Exists() {
	case true:
		log.Info("⚙️   LUKSO configuration already exists - continuing...")
	case false:
		log.Info("⚙️   Creating LUKSO configuration file...")

		err = cfg.Create("", "")
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was an error while preparing LUKSO configuration: %v", err), 1)
		}

		log.Infof("✅  LUKSO configuration created under %s", config.Path)
	}

	log.Info("✅  Working directory initialized! \n1. ⚙️  Use 'lukso install' to install clients. \n2. ▶️  Use 'lukso start' to start your node.")

	return nil
}

// initConfigGroup takes map of config dependencies and downloads them.
func initConfigGroup(configDependencies map[string]*ClientDependency) error {
	for _, dependency := range configDependencies {
		err := dependency.Download("", "", false, configPerms)
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was error while downloading %s file: %v", dependency.name, err), 1)
		}
	}

	return nil
}
