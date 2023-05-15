package main

import (
	"fmt"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"

	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Version = "0.6.0"

const (
	ubuntu  = "linux"
	macos   = "darwin"
	windows = "windows"

	unixBinDir = "/usr/local/bin"

	// should be a user-created path, like C:\bin,
	// but since it is not guaranteed that all users vahe it we can leave it as is
	windowsBinDir = "/Windows/System32"
)

var (
	appName        = "lukso"
	binDir         string
	gethTag        string
	gethCommitHash string
	validatorTag   string
	prysmTag       string
	log            = logrus.StandardLogger()
	systemOs       string
	cfg            *config.Config
)

// do not confuse init func with init subcommand - for now we just pass init flags, there are more to come
func init() {
	installFlags = make([]cli.Flag, 0)
	installFlags = append(installFlags, gethDownloadFlags...)
	installFlags = append(installFlags, validatorDownloadFlags...)
	installFlags = append(installFlags, prysmDownloadFlags...)
	installFlags = append(installFlags, erigonDownloadFlags...)
	installFlags = append(installFlags, lighthouseDownloadFlags...)
	installFlags = append(installFlags, appFlags...)

	updateFlags = append(updateFlags, gethUpdateFlags...)
	updateFlags = append(updateFlags, prysmUpdateFlags...)
	updateFlags = append(updateFlags, validatorUpdateFlags...)

	startFlags = append(startFlags, gethStartFlags...)
	startFlags = append(startFlags, erigonStartFlags...)
	startFlags = append(startFlags, prysmStartFlags...)
	startFlags = append(startFlags, lighthouseStartFlags...)
	startFlags = append(startFlags, validatorStartFlags...)
	startFlags = append(startFlags, networkFlags...)

	logsFlags = append(logsFlags, gethLogsFlags...)
	logsFlags = append(logsFlags, prysmLogsFlags...)
	logsFlags = append(logsFlags, validatorLogsFlags...)
	logsFlags = append(logsFlags, networkFlags...)

	resetFlags = append(resetFlags, gethResetFlags...)
	resetFlags = append(resetFlags, prysmResetFlags...)
	resetFlags = append(resetFlags, validatorResetFlags...)
	resetFlags = append(resetFlags, networkFlags...)

	// after we exported flags from each command we can update them
	gethStartFlags = append(gethStartFlags, networkFlags...)
	gethLogsFlags = append(gethLogsFlags, networkFlags...)
	gethResetFlags = append(gethResetFlags, networkFlags...)

	prysmStartFlags = append(prysmStartFlags, networkFlags...)
	prysmLogsFlags = append(prysmLogsFlags, networkFlags...)
	prysmResetFlags = append(prysmResetFlags, networkFlags...)

	validatorStartFlags = append(validatorStartFlags, networkFlags...)
	validatorLogsFlags = append(validatorLogsFlags, networkFlags...)
	validatorResetFlags = append(validatorResetFlags, networkFlags...)

	validatorImportFlags = append(validatorImportFlags, networkFlags...)
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.ExitErrHandler = func(c *cli.Context, err error) {} // this (somehow) produces error logs instead of standard output

	cli.AppHelpTemplate = fmt.Sprintf(`%s

DOCS: https://docs.lukso.tech/
REPO: https://github.com/lukso-network/tools-lukso-cli

`, cli.AppHelpTemplate)

	app.Usage = "The LUKSO CLI is a command line tool to install, manage and set up validators of different clients for the LUKSO Blockchain."
	app.Flags = appFlags
	app.Commands = []*cli.Command{
		{
			Name:            "install",
			Action:          installBinaries,
			Flags:           installFlags,
			Usage:           "Installs chosen LUKSO clients (Execution, Consensus, Validator) and their binary dependencies",
			Before:          initializeFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "init",
			Usage:           "Initializes the working directory, its structure, and network configuration",
			Action:          initializeDirectory,
			Before:          initializeFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "update",
			Usage:           "Updates all or specific clients in the working directory to the latest version",
			Action:          updateClients,
			Before:          initializeFlags,
			Flags:           updateFlags,
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "geth",
					Usage:           "Updates your geth client to latest version",
					Action:          updateGethToSpec,
					Flags:           gethUpdateFlags,
					Before:          initializeFlags,
					HideHelpCommand: true,
				},
				{
					Name:            "prysm",
					Usage:           "Updates your prysm client to latest version",
					Action:          updatePrysmToSpec,
					Flags:           prysmUpdateFlags,
					Before:          initializeFlags,
					HideHelpCommand: true,
				},
				{
					Name:            "validator",
					Usage:           "Updates your validator client to latest version",
					Action:          updateValidatorToSpec,
					Flags:           validatorUpdateFlags,
					Before:          initializeFlags,
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "start",
			Usage:           "Starts all or specific clients and connects to the specified network",
			Action:          selectNetworkFor(startClients),
			SkipFlagParsing: true,
			Flags:           startFlags,
			Before:          initializeFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "stop",
			Usage:           "Stops all or specific client(s) that are currently running",
			Action:          stopClients,
			Flags:           stopFlags,
			HideHelpCommand: true,
		},
		{
			Name:            "logs",
			Usage:           "Streams all logs from a specific client in the current terminal window",
			Action:          logClients,
			HideHelpCommand: true,
			Subcommands: []*cli.Command{
				{
					Name:            "execution",
					Usage:           "Outputs selected execution client's logs, add the network flag, if not mainnet",
					Flags:           gethLogsFlags,
					Action:          selectNetworkFor(logLayer(executionLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "consensus",
					Usage:           "Outputs selected consensus client's logs, add the network flag, if not mainnet",
					Flags:           prysmLogsFlags,
					Action:          selectNetworkFor(logLayer(consensusLayer)),
					HideHelpCommand: true,
				},
				{
					Name:            "validator",
					Usage:           "Outputs selected validator client's logs, add the network flag, if not mainnet",
					Flags:           validatorLogsFlags,
					Action:          selectNetworkFor(logLayer(validatorLayer)), // named as a layer for sake of
					HideHelpCommand: true,
				},
			},
		},
		{
			Name:            "status",
			Usage:           "Shows the client processes that are currently running",
			Action:          statClients,
			HideHelpCommand: true,
		},
		{
			Name:            "reset",
			Usage:           "Resets all or specific client data directories and logs excluding the validator keys",
			Flags:           resetFlags,
			Action:          selectNetworkFor(resetClients),
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
					Flags:           validatorImportFlags,
					Action:          selectNetworkFor(importValidator),
					SkipFlagParsing: true,
					HideHelpCommand: true,
				},
				{
					Name:            "list",
					Usage:           "List your imported validator keys from the client wallet",
					Flags:           validatorImportFlags,
					Action:          selectNetworkFor(listValidator),
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
		cfg = config.NewConfig(config.Path)

		setupOperatingSystem()

		// Geth related parsing
		gethTag = ctx.String(gethTagFlag)
		gethCommitHash = ctx.String(gethCommitHashFlag)

		// Validator related parsing
		validatorTag = ctx.String(validatorTagFlag)

		// Prysm related parsing
		prysmTag = ctx.String(prysmTagFlag)

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

func initializeFlags(ctx *cli.Context) error {
	// Geth related parsing
	gethTag = ctx.String(gethTagFlag)
	gethCommitHash = ctx.String(gethCommitHashFlag)

	// Validator related parsing
	validatorTag = ctx.String(validatorTagFlag)

	// Prysm related parsing
	prysmTag = ctx.String(prysmTagFlag)

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func displayVersion(ctx *cli.Context) error {
	fmt.Println("Version:", Version)
	return nil
}
