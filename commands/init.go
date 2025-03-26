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

func setIPInConfigs() (ip string, err error) {
	ip, err = getPublicIP()
	if err != nil {
		return
	}

	type tempConfig struct {
		path          string
		name          string
		fileType      string
		flagsToChange []string
		// wrap is an optional func that should edit the ip flag value.
		wrap func(string) string
	}

	natWrap := func(s string) string {
		return "extip:" + s
	}

	erigonFlags := []string{"nat"}
	besuFlags := []string{"p2p-host"}
	nethermindFlags := []string{"Network:ExternalIp"} // : because of the non-standard delimiter
	prysmFlags := []string{"p2p-host-ip"}
	lighthouseFlags := []string{"enr-address"}
	tekuFlags := []string{"p2p-advertised-ip"}
	nimbusFlags := []string{"nat"}

	configs := []tempConfig{
		{
			path:          "./configs/%s/erigon",
			name:          "erigon",
			fileType:      "toml",
			flagsToChange: erigonFlags,
			wrap:          natWrap,
		},
		{
			path:          "./configs/%s/besu",
			name:          "besu",
			fileType:      "toml",
			flagsToChange: besuFlags,
			wrap:          nil,
		},
		{
			path:          "./configs/%s/nethermind",
			name:          "nethermind",
			fileType:      "json",
			flagsToChange: nethermindFlags,
			wrap:          nil,
		},
		{
			path:          "./configs/%s/prysm",
			name:          "prysm",
			fileType:      "yaml",
			flagsToChange: prysmFlags,
			wrap:          nil,
		},
		{
			path:          "./configs/%s/lighthouse",
			name:          "lighthouse",
			fileType:      "toml",
			flagsToChange: lighthouseFlags,
			wrap:          nil,
		},
		{
			path:          "./configs/%s/teku",
			name:          "teku",
			fileType:      "yaml",
			flagsToChange: tekuFlags,
			wrap:          nil,
		},
		{
			path:          "./configs/%s/nimbus2",
			name:          "nimbus",
			fileType:      "toml",
			flagsToChange: nimbusFlags,
			wrap:          natWrap,
		},
	}

	networks := []string{"mainnet", "testnet"}

	for _, network := range networks {
		for _, cfg := range configs {
			cfg.path = fmt.Sprintf(cfg.path, network)

			v := viper.NewWithOptions(viper.KeyDelimiter(":"))
			v.AddConfigPath(cfg.path)
			v.SetConfigName(cfg.name)
			v.SetConfigType(cfg.fileType)

			err = v.ReadInConfig()
			if err != nil {
				return
			}

			cfgIp := ip
			if cfg.wrap != nil {
				cfgIp = cfg.wrap(cfgIp)
			}

			for _, flag := range cfg.flagsToChange {
				v.Set(flag, cfgIp)
			}

			err = v.WriteConfig()
			if err != nil {
				return
			}
		}
	}

	names := ""
	for _, cfg := range configs {
		names += fmt.Sprintf("- %s\n", cfg.name)
	}

	log.Info("✅  IP Address updated!\n\n")
	log.Infof("▶️  Your IPv4 Address has been wrtitten into the cli-config.yaml file.\n\n"+
		"Additionally, the following clients' configs have been updated: \n%s\n", names)

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
