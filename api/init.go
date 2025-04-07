package api

import (
	"fmt"
	"net"
	"strings"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

func (h *handler) Init(args types.InitArgs) (resp types.InitResponse) {
	if runningClients := clients.RunningClients(); runningClients != nil {
		return types.InitResponse{
			Error: errors.ErrClientsAlreadyRunning{Clients: runningClients},
		}
	}

	if h.cfg.Exists() && !args.Reinit {
		return types.InitResponse{
			Error: errors.ErrCfgExists,
		}
	}

	var (
		ip  = "disabled"
		err error
	)

	switch args.Ip {
	case "disabled":
		ip = "disabled"

	case "auto":
		ip, err = h.getPublicIP()
		if err != nil {
			h.log.Warn("unable to automatically get IP - proceeding with --ip=disabled")

			break
		}

		h.log.Info("Auto IP set to " + ip)

	default:
		parsedIp := net.ParseIP(args.Ip)
		if parsedIp == nil {
			h.log.Warn(fmt.Sprintf("--ip=%s is not a valid IP - proceeding with --ip=disabled", ip))

			break
		}

		// This means it's IPv6 - only v4 is supported
		if strings.ContainsAny(args.Ip, ":%") {
			h.log.Warn("--ip is an IPv6, but IPv4 is required - proceeding with --ip=disabled")

			break
		}

		ip = args.Ip
	}

	h.installInitFiles(args.Reinit)

	// Create shared folder
	err = h.file.Mkdir(file.SecretsDir, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create secrets directory: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	// Create and Write JWT
	jwt, err := utils.CreateJwtSecret()
	if err != nil {
		err = fmt.Errorf("unable to generate JWT secret: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	err = h.file.Write(file.JwtSecretPath, jwt, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create JWT secret file: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	// Create dir for PIDs
	err = h.file.Mkdir(file.PidDir, common.ConfigPerms)
	if err != nil {
		err = fmt.Errorf("unable to create PID directory: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	h.log.Info("⚙️   Creating LUKSO configuration file...")

	err = h.cfg.Create("", "", "", ip)
	if err != nil {
		err = fmt.Errorf("unable to create LUKSO config: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	h.log.Info(fmt.Sprintf("✅  LUKSO configuration created in %s", config.Path))

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
