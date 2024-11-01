package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const jwtSecretPath = configs.ConfigRootDir + "/shared/secrets/jwt.hex"

// InitializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func InitializeDirectory(ctx *cli.Context) error {
	if clients.IsAnyRunning() {
		return nil
	}

	if cfg.Exists() {
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

	switch cfg.Exists() {
	case true:
		log.Info("⚙️   LUKSO configuration already exists - continuing...")
	case false:
		log.Info("⚙️   Creating LUKSO configuration file...")

		err = cfg.Create("", "", "")
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while preparing LUKSO configuration: %v", err), 1)
		}

		log.Infof("✅  LUKSO configuration created under %s", config.Path)
	}

	log.Info("⚙️   LUKSO CLI can replace your p2p communication flags with your public IP for better connection (optional)")
	message := "Do you want to proceed? [Y/n]\n> "

	input := utils.RegisterInputWithMessage(message)
	if strings.EqualFold(input, "y") || input == "" {
		err = setIPInConfigs()
		if err != nil {
			return err
		}
	}

	displayNetworksHardforkTimestamps()

	log.Info("✅  Working directory initialized! \n1. ⚙️  Use 'lukso install' to install clients. \n2. ▶️  Use 'lukso start' to start your node.")

	return nil
}

func setIPInConfigs() (err error) {
	ip, err := getPublicIP()
	if err != nil {
		return err
	}

	type tempConfig struct {
		path          string
		name          string
		fileType      string
		flagsToChange []string
	}

	prysmFlagsToChange := []string{"p2p-host-ip"}
	lighthouseFlagsToChange := []string{"enr-address"}
	tekuFlagsToChange := []string{"p2p-advertised-ip"}
	configs := []tempConfig{
		{
			path:          "./configs/testnet/prysm",
			name:          "prysm",
			fileType:      "yaml",
			flagsToChange: prysmFlagsToChange,
		},
		{
			path:          "./configs/testnet/lighthouse",
			name:          "lighthouse",
			fileType:      "toml",
			flagsToChange: lighthouseFlagsToChange,
		},
		{
			path:          "./configs/testnet/teku",
			name:          "teku",
			fileType:      "yaml",
			flagsToChange: tekuFlagsToChange,
		},
		{
			path:          "./configs/mainnet/prysm",
			name:          "prysm",
			fileType:      "yaml",
			flagsToChange: prysmFlagsToChange,
		},
		{
			path:          "./configs/mainnet/lighthouse",
			name:          "lighthouse",
			fileType:      "toml",
			flagsToChange: lighthouseFlagsToChange,
		},
		{
			path:          "./configs/mainnet/teku",
			name:          "teku",
			fileType:      "yaml",
			flagsToChange: tekuFlagsToChange,
		},
	}

	for _, cfg := range configs {
		v := viper.New()

		v.AddConfigPath(cfg.path)
		v.SetConfigName(cfg.name)
		v.SetConfigType(cfg.fileType)

		err = v.ReadInConfig()
		if err != nil {
			return
		}

		for _, flag := range cfg.flagsToChange {
			v.Set(flag, ip)
		}

		err = v.WriteConfig()
		if err != nil {
			return
		}
	}

	log.Info("✅  IP Address updated!\n\n")

	return
}

func getPublicIP() (ip string, err error) {
	url := "https://ipv4.ident.me"
	res, err := http.Get(url)
	if err != nil {
		return
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	ip = string(respBody)

	return
}
