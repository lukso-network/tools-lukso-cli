package main

import (
	"fmt"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api"
	"github.com/lukso-network/tools-lukso-cli/commands"
	"github.com/lukso-network/tools-lukso-cli/common/display"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/progress"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/dep/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

var Version = "develop"

var appName = "lukso"

func main() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	// Bottom to top dependencies
	// display for commander
	dispCh := make(chan tea.Msg, 10000)
	prg := progress.NewProgress(dispCh)
	hndlLogger := logger.NewMsgLogger(dispCh, prg)
	cmdDisplay := display.NewCmdDisplay(dispCh, nil)

	// file/installation management
	hndlFile := file.NewManager()
	hndlInstaller := installer.NewInstaller(hndlFile)
	cfg := config.NewConfigurator(config.Path, hndlFile)
	p := pid.NewPid(hndlFile)

	config.UseConfigurator(cfg)
	if config.Exists() {
		err := config.Read()
		if err != nil {
			panic("unable to read config file: %v")
		}
	}

	hndl := api.NewHandler(
		hndlFile,
		hndlLogger,
		hndlInstaller,
	)

	clients.Setup(
		hndlLogger,
		hndlFile,
		hndlInstaller,
		p,
	)

	cmd := commands.NewCommander(hndl, cmdDisplay, hndlLogger)

	app := cli.App{}
	app.Name = appName

	cli.AppHelpTemplate = fmt.Sprintf(`%s

DOCS: https://docs.lukso.tech/
REPO: https://github.com/lukso-network/tools-lukso-cli

`, cli.AppHelpTemplate)
	app.ExitErrHandler = func(c *cli.Context, err error) {} // this (somehow) produces error logs instead of standard output
	app.Usage = "The LUKSO CLI is a command line tool to install, manage and set up validators of different clients for the LUKSO Blockchain."
	app.Flags = flags.AppFlags
	app.Commands = []*cli.Command{
		{
			Name:            "install",
			Flags:           flags.InstallFlags,
			Before:          cmd.Before,
			Action:          cmd.Install,
			After:           cmd.After,
			Usage:           "Installs chosen LUKSO clients (Execution, Consensus, Validator) and their binary dependencies",
			HideHelpCommand: true,
		},
		{
			Name:            "init",
			Usage:           "Initializes the working directory, its structure, and network configuration",
			Flags:           flags.InitFlags,
			Before:          cmd.Before,
			Action:          cmd.Init,
			After:           cmd.After,
			HideHelpCommand: true,
		},
		{
			Name:            "update",
			Usage:           "Updates all or specific clients to the latest version",
			Action:          cmd.Update,
			Flags:           flags.InstallFlags,
			HideHelpCommand: true,
			Subcommands: cli.Commands{
				&cli.Command{
					Name:            "configs",
					Usage:           "Updates chain configuration and CLI configuration files, without overwriting client configuration files",
					Flags:           flags.UpdateConfigFlags,
					Action:          cmd.UpdateConfigs,
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "start",
			Usage:           "Starts all or specific clients and connects to the specified network",
			Action:          commands.SelectNetworkFor(cmd.Start),
			SkipFlagParsing: true,
			Flags:           flags.StartFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "stop",
			Usage:           "Stops all or specific client(s) that are currently running",
			Action:          cmd.Stop,
			Flags:           flags.StopFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "logs",
			Usage:           "Streams all logs from a specific client in the current terminal window",
			Action:          cmd.Logs,
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "execution",
					Usage:           "Outputs selected execution client's logs, add the network flag, if not mainnet",
					Flags:           flags.ExecutionLogsFlags,
					Action:          commands.SelectNetworkFor(cmd.LogsLayer(configs.ExecutionLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "consensus",
					Usage:           "Outputs selected consensus client's logs, add the network flag, if not mainnet",
					Flags:           flags.ConsensusLogsFlags,
					Action:          commands.SelectNetworkFor(cmd.LogsLayer(configs.ConsensusLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "validator",
					Usage:           "Outputs selected validator client's logs, add the network flag, if not mainnet",
					Flags:           flags.ValidatorLogsFlags,
					Action:          commands.SelectNetworkFor(cmd.LogsLayer(configs.ValidatorLayer)), // named as a layer for sake of
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "status",
			Usage:           "Shows the client processes that are currently running",
			Action:          cmd.Status,
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:   "peers",
					Usage:  "Displays number of peers that your clients have",
					Flags:  flags.StatusPeersFlags,
					Action: cmd.StatusPeers,
				},
			},
		},
		{
			Name:            "reset",
			Usage:           "Resets all or specific client data directories and logs excluding the validator keys",
			Flags:           flags.ResetFlags,
			Action:          commands.SelectNetworkFor(cmd.Reset),
			HideHelpCommand: true,
		},
		{
			Name:            "validator",
			Usage:           "Manages the LUKSO validator keys including their initialization",
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "import",
					Usage:           "Import your validator keys in the client wallet",
					Flags:           flags.ValidatorImportFlags,
					Action:          commands.SelectNetworkFor(cmd.ValidatorImport),
					HideHelpCommand: true,
				},
				{
					Name:            "list",
					Usage:           "List your imported validator keys from the client wallet",
					Flags:           flags.ValidatorListFlags,
					Action:          commands.SelectNetworkFor(cmd.ValidatorList),
					SkipFlagParsing: true,
					HideHelpCommand: true,
				},
				{
					Name:            "exit",
					Usage:           "Issue an exit for your validator",
					Flags:           flags.ValidatorExitFlags,
					Action:          commands.SelectNetworkFor(cmd.ValidatorExit),
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "version",
			Usage:           "Display the version of the LUKSO CLI that is currently installed",
			Action:          cmd.Version(Version),
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "clients",
					Usage:           "Display supported clients' installed versions.",
					Action:          cmd.VersionClients,
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:   "ui",
			Action: cmd.Ui,
		},
	}

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())

		return nil
	}

	defer func() {
		if x := recover(); x != nil {
			log.Errorf("Runtime panic: %v\n%v", x, string(runtimeDebug.Stack()))
			panic(x)
		}
	}()

	err := app.Run(os.Args)

	if nil != err {
		log.Error(err.Error())
	}
}
