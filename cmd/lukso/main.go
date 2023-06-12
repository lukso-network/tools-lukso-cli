package main

import (
	"fmt"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/commands"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

var Version = "develop"

var (
	appName = "lukso"
)

func main() {
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
			Action:          commands.InstallBinaries,
			Flags:           flags.InstallFlags,
			Usage:           "Installs chosen LUKSO clients (Execution, Consensus, Validator) and their binary dependencies",
			HideHelpCommand: true,
		},
		{
			Name:            "init",
			Usage:           "Initializes the working directory, its structure, and network configuration",
			Action:          commands.InitializeDirectory,
			HideHelpCommand: true,
		},
		{
			Name:            "update",
			Usage:           "Updates all or specific clients to the latest version",
			Action:          commands.UpdateClients,
			Flags:           flags.UpdateFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "start",
			Usage:           "Starts all or specific clients and connects to the specified network",
			Action:          commands.SelectNetworkFor(commands.StartClients),
			SkipFlagParsing: true,
			Flags:           flags.StartFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "stop",
			Usage:           "Stops all or specific client(s) that are currently running",
			Action:          commands.StopClients,
			Flags:           flags.StopFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "logs",
			Usage:           "Streams all logs from a specific client in the current terminal window",
			Action:          commands.LogClients,
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "execution",
					Usage:           "Outputs selected execution client's logs, add the network flag, if not mainnet",
					Flags:           flags.GethLogsFlags,
					Action:          commands.SelectNetworkFor(commands.LogLayer(configs.ExecutionLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "consensus",
					Usage:           "Outputs selected consensus client's logs, add the network flag, if not mainnet",
					Flags:           flags.PrysmLogsFlags,
					Action:          commands.SelectNetworkFor(commands.LogLayer(configs.ConsensusLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "validator",
					Usage:           "Outputs selected validator client's logs, add the network flag, if not mainnet",
					Flags:           flags.ValidatorLogsFlags,
					Action:          commands.SelectNetworkFor(commands.LogLayer(configs.ValidatorLayer)), // named as a layer for sake of
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "status",
			Usage:           "Shows the client processes that are currently running",
			Action:          commands.StatClients,
			HideHelpCommand: true,
		},
		{
			Name:            "reset",
			Usage:           "Resets all or specific client data directories and logs excluding the validator keys",
			Flags:           flags.ResetFlags,
			Action:          commands.SelectNetworkFor(commands.ResetClients),
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
					Action:          commands.SelectNetworkFor(commands.ImportValidator),
					HideHelpCommand: true,
				},
				{
					Name:            "list",
					Usage:           "List your imported validator keys from the client wallet",
					Flags:           flags.ValidatorListFlags,
					Action:          commands.SelectNetworkFor(commands.ListValidator),
					SkipFlagParsing: true,
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "version",
			Usage:           "Display the version of the LUKSO CLI that is currently installed",
			Action:          displayVersion,
			HideHelpCommand: true,
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

func displayVersion(ctx *cli.Context) error {
	fmt.Println("Version:", Version)
	return nil
}
