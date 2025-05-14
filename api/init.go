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
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/dep/configs"
)

func (h *handler) Init(args types.InitRequest) (resp types.InitResponse) {
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
			h.log.Warn(fmt.Sprintf("--ip='%s' is not a valid IP - proceeding with --ip=disabled", args.Ip))

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

	h.log.Info("Creating LUKSO configuration file...")

	nodeCfg := config.NodeConfig{
		UseClients: config.UseClients{
			ExecutionClient: "",
			ConsensusClient: "",
			ValidatorClient: "",
		},
		Ipv4: ip,
	}

	err = h.cfg.Create(nodeCfg)
	if err != nil {
		err = fmt.Errorf("unable to create LUKSO config: %w", err)
		return types.InitResponse{
			Error: err,
		}
	}

	h.log.Info(fmt.Sprintf("LUKSO configuration created in %s", config.Path))

	return types.InitResponse{}
}

func (h *handler) installInitFiles(isUpdate bool) {
	h.log.Info("Downloading shared configuration files...")
	h.installConfigGroup(configs.SharedConfigDependencies, isUpdate)
	h.log.Info("Shared configuration files downloaded!\n")

	h.installClientConfigFiles(isUpdate)
}

func (h *handler) installConfigGroup(configDependencies map[string]configs.ClientConfigDependency, isUpdate bool) {
	for _, dependency := range configDependencies {
		err := dependency.Install(isUpdate)
		if err != nil {
			h.log.Warn(fmt.Sprintf("Unable to download %s file: %v - continuing...", dependency.Name(), err))
		}
		h.log.Progress().Move(1)
	}
}

func (h *handler) installClientConfigFiles(isUpdate bool) {
	h.log.Info("Downloading geth configuration files...")
	h.installConfigGroup(configs.GethConfigDependencies, isUpdate)
	h.log.Info("Geth configuration files downloaded!\n")

	h.log.Info("Downloading erigon configuration files...")
	h.installConfigGroup(configs.ErigonConfigDependencies, isUpdate)
	h.log.Info("Erigon configuration files downloaded!\n")

	h.log.Info("Downloading nethermind configuration files...")
	h.installConfigGroup(configs.NethermindConfigDependencies, isUpdate)
	h.log.Info("Nethermind configuration files downloaded!\n")

	h.log.Info("Downloading besu configuration files...")
	h.installConfigGroup(configs.BesuConfigDependencies, isUpdate)
	h.log.Info("Besu configuration files downloaded!\n")

	h.log.Info("Downloading prysm configuration files...")
	h.installConfigGroup(configs.PrysmConfigDependencies, isUpdate)
	h.log.Info("Prysm configuration files downloaded!\n")

	h.log.Info("Downloading lighthouse configuration files...")
	h.installConfigGroup(configs.LighthouseConfigDependencies, isUpdate)
	h.log.Info("Lighthouse configuration files downloaded!\n")

	h.log.Info("Downloading prysm validator configuration files...")
	h.installConfigGroup(configs.PrysmValidatorConfigDependencies, isUpdate)
	h.log.Info("Prysm validator configuration files downloaded!\n")

	h.log.Info("Downloading teku configuration files...")
	h.installConfigGroup(configs.TekuConfigDependencies, isUpdate)
	h.log.Info("Teku configuration files downloaded!\n")

	h.log.Info("Downloading nimbus (eth2) configuration files...")
	h.installConfigGroup(configs.Nimbus2ConfigDependencies, isUpdate)
	h.log.Info("Nimbus configuration files downloaded!\n")
}
