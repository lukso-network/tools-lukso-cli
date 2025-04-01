package commands

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/model"
)

const jwtSecretPath = configs.ConfigRootDir + "/shared/secrets/jwt.hex"

// InitializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func (c *commander) Init(ctx *cli.Context) error {
	if clients.IsAnyRunning() {
		return nil
	}

	req := model.CtxToApiInit(ctx)

	resp := c.handler.Init(req)
	err := resp.Error()
	if err != nil {
		return err
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
