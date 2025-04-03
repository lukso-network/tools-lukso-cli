package api

import (
	"fmt"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

func (h *handler) Init(args types.InitArgs) (resp types.Response) {
	if clients.IsAnyRunning() {
		return types.Error(errors.ErrClientRunning)
	}

	if h.cfg.Exists() && !args.Reinit {
		return types.Error(errors.ErrCfgExists)
	}

	h.installInitFiles(args.Reinit)

	// TODO: make sure there is a global identifier for selected directory
	// Create shared folder
	err := h.file.Mkdir("."+file.SecretsDir, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create secrets directory: %w", err)
		return types.Error(err)
	}

	// Create and Write JWT
	jwt, err := utils.CreateJwtSecret()
	if err != nil {
		err = fmt.Errorf("unable to generate JWT secret: %w", err)
		return types.Error(err)
	}

	err = h.file.Write("."+file.JwtSecretPath, jwt, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create JWT secret file: %w", err)
		return types.Error(err)
	}

	// Create dir for PIDs
	err = h.file.Mkdir(file.PidDir, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create PID directory: %w", err)
		return types.Error(err)
	}

	switch h.cfg.Exists() {
	case true:
		h.log.Info("⚙️   LUKSO configuration already exists - continuing...")
	case false:
		h.log.Info("⚙️   Creating LUKSO configuration file...")

		err = h.cfg.Create("", "", "", args.Ip)
		if err != nil {
			err = fmt.Errorf("unable to create LUKSO config: %w", err)
			return types.Error(err)
		}

		h.log.Info(fmt.Sprintf("✅  LUKSO configuration created in %s", config.Path))
	}

	return types.InitResponse{}
}

func (h *handler) installInitFiles(isUpdate bool) {
	h.log.Info("⬇️  Downloading shared configuration files...")
	_ = h.installConfigGroup(configs.SharedConfigDependencies, isUpdate) // we can omit errors - all errors are catched by cli.Exit()
	h.log.Info("✅  Shared configuration files downloaded!\n\n")

	h.installClientConfigFiles(isUpdate)
}

func (h *handler) installConfigGroup(configDependencies map[string]configs.ClientConfigDependency, isUpdate bool) error {
	for _, dependency := range configDependencies {
		err := dependency.Install(isUpdate)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was error while downloading %s file: %v", dependency.Name(), err), 1)
		}
	}

	return nil
}

func (h *handler) installClientConfigFiles(isUpdate bool) {
	h.log.Info("⬇️  Downloading geth configuration files...")
	_ = h.installConfigGroup(configs.GethConfigDependencies, isUpdate)
	h.log.Info("✅  Geth configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading erigon configuration files...")
	_ = h.installConfigGroup(configs.ErigonConfigDependencies, isUpdate)
	h.log.Info("✅  Erigon configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading nethermind configuration files...")
	_ = h.installConfigGroup(configs.NethermindConfigDependencies, isUpdate)
	h.log.Info("✅  Nethermind configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading besu configuration files...")
	_ = h.installConfigGroup(configs.BesuConfigDependencies, isUpdate)
	h.log.Info("✅  Besu configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading prysm configuration files...")
	_ = h.installConfigGroup(configs.PrysmConfigDependencies, isUpdate)
	h.log.Info("✅  Prysm configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading lighthouse configuration files...")
	_ = h.installConfigGroup(configs.LighthouseConfigDependencies, isUpdate)
	h.log.Info("✅  Lighthouse configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading prysm validator configuration files...")
	_ = h.installConfigGroup(configs.PrysmValidatorConfigDependencies, isUpdate)
	h.log.Info("✅  Prysm validator configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading teku configuration files...")
	_ = h.installConfigGroup(configs.TekuConfigDependencies, isUpdate)
	h.log.Info("✅  Teku configuration files downloaded!\n\n")

	h.log.Info("⬇️  Downloading nimbus (eth2) configuration files...")
	_ = h.installConfigGroup(configs.Nimbus2ConfigDependencies, isUpdate)
	h.log.Info("✅  Nimbus configuration files downloaded!\n\n")
}
