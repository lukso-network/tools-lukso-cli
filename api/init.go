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

	h.log.Info("⬇️  Downloading shared configuration files...")
	_ = installConfigGroup(configs.SharedConfigDependencies, false) // we can omit errors - all errors are catched by cli.Exit()
	h.log.Info("✅  Shared configuration files downloaded!\n\n")

	installClientConfigFiles(false)

	jwt, err := utils.CreateJwtSecret()
	if err != nil {
		err = fmt.Errorf("unable to generate JWT secret: %w", err)
		return types.Error(err)
	}

	err = h.file.Write(jwt, file.JwtSecretPath, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create JWT secret file: %w", err)
		return types.Error(err)
	}

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

	return nil
}
