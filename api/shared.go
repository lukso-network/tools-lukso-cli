package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

func installConfigGroup(configDependencies map[string]configs.ClientConfigDependency, isUpdate bool) error {
	for _, dependency := range configDependencies {
		err := dependency.Install(isUpdate)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was error while downloading %s file: %v", dependency.Name(), err), 1)
		}
	}

	return nil
}

func installClientConfigFiles(override bool) {
	log.Info("⬇️  Downloading geth configuration files...")
	_ = installConfigGroup(configs.GethConfigDependencies, override)
	log.Info("✅  Geth configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading erigon configuration files...")
	_ = installConfigGroup(configs.ErigonConfigDependencies, override)
	log.Info("✅  Erigon configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading nethermind configuration files...")
	_ = installConfigGroup(configs.NethermindConfigDependencies, override)
	log.Info("✅  Nethermind configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading besu configuration files...")
	_ = installConfigGroup(configs.BesuConfigDependencies, override)
	log.Info("✅  Besu configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading prysm configuration files...")
	_ = installConfigGroup(configs.PrysmConfigDependencies, override)
	log.Info("✅  Prysm configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading lighthouse configuration files...")
	_ = installConfigGroup(configs.LighthouseConfigDependencies, override)
	log.Info("✅  Lighthouse configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading prysm validator configuration files...")
	_ = installConfigGroup(configs.PrysmValidatorConfigDependencies, override)
	log.Info("✅  Prysm validator configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading teku configuration files...")
	_ = installConfigGroup(configs.TekuConfigDependencies, override)
	log.Info("✅  Teku configuration files downloaded!\n\n")

	log.Info("⬇️  Downloading nimbus (eth2) configuration files...")
	_ = installConfigGroup(configs.Nimbus2ConfigDependencies, override)
	log.Info("✅  Nimbus configuration files downloaded!\n\n")
}
