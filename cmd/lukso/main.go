package main

import (
	"fmt"
	"github.com/m8b-dev/lukso-cli/config"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Version = "develop"

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
	installFlags = append(installFlags, appFlags...)

	updateFlags = append(updateFlags, gethUpdateFlags...)
	updateFlags = append(updateFlags, prysmUpdateFlags...)
	updateFlags = append(updateFlags, validatorUpdateFlags...)

	startFlags = append(startFlags, gethStartFlags...)
	startFlags = append(startFlags, prysmStartFlags...)
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

	validatorInitFlags = append(validatorInitFlags, networkFlags...)
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.Usage = "The LUKSO CLI is a command line tool to install, manage and set up validators of different types of nodes for the LUKSO network."
	app.Flags = appFlags
	app.Commands = []*cli.Command{
		{
			Name:   "install",
			Action: installBinaries,
			Flags:  installFlags,
			Usage:  "Installs choosen LUKSO node clients (Execution, Beacon, Validator) and their binary dependencies",
			Before: initializeFlags,
		},
		{
			Name: "init",
			Usage: "Initializes the node working directory, it's structure, and network configuration",
			Action: initializeDirectory,
			Flags:  networkFlags,
			Before: initializeFlags,
		},
		{
			Name:   "update",
			Usage:  "Updates all or specific LUKSO node clients in the working directory to the newest version",
			Action: updateClients,
			Before: initializeFlags,
			Flags:  updateFlags,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Updates your geth client to newest version",
					Action: updateGethToSpec,
					Flags:  gethUpdateFlags,
					Before: initializeFlags,
				},
				{
					Name:   "prysm",
					Usage:  "Updates your prysm client to newest version",
					Action: updatePrysmToSpec,
					Flags:  prysmUpdateFlags,
					Before: initializeFlags,
				},
				{
					Name:   "validator",
					Usage:  "Updates your validator client to newest version",
					Action: updateValidatorToSpec,
					Flags:  validatorUpdateFlags,
					Before: initializeFlags,
				},
			},
		},
		{
			Name:            "start",
			Usage:           "Starts all or specific LUKSO node clients and connects to the specified network",
			Action:          selectNetworkFor(startClients),
			SkipFlagParsing: true,
			Flags:           startFlags,
			Before:          initializeFlags,
		},
		{
			Name:   "stop",
			Usage:  "Stops all or specific LUKSO node clients that are currently running",
			Action: stopClients,
			Flags:  stopFlags,
		},
		{
			Name:   "log",
			Usage:  "Listens to all log events from a specific client in the current terminal window",
			Action: logClients,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Outputs Geth client logs",
					Flags:  gethLogsFlags,
					Action: selectNetworkFor(logClient(gethDependencyName)),
				},
				{
					Name:   "prysm",
					Usage:  "Outputs Prysm client logs",
					Flags:  prysmLogsFlags,
					Action: selectNetworkFor(logClient(prysmDependencyName)),
				},
				{
					Name:   "validator",
					Usage:  "Outputs Validator client logs",
					Flags:  validatorLogsFlags,
					Action: selectNetworkFor(logClient(validatorDependencyName)),
				},
			},
		},
		{
			Name:   "status",
			Usage:  "Shows the LUKSO node client processes that are currently running",
			Action: statClients,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Displays Geth's running status",
					Action: statClient(gethDependencyName),
				},
				{
					Name:   "prysm",
					Usage:  "Displays Prysm's running status",
					Action: statClient(prysmDependencyName),
				},
				{
					Name:   "validator",
					Usage:  "Displays Validator's running status",
					Action: statClient(validatorDependencyName),
				},
			},
		},
		{
			Name:   "reset",
			Usage:  "Resets all or specific client data directories and logs excluding the validator keys",
			Flags:  resetFlags,
			Action: selectNetworkFor(resetClients),
		},
		{
			Name:  "validator",
			Usage: "Manages the LUKSO validator keys including their initialization and deposits",
			Subcommands: []*cli.Command{
				{
					Name:   "init",
					Usage:  "Initializes your validator with keys generated using keygen tools",
					Flags:  validatorInitFlags,
					Action: selectNetworkFor(initValidator),
				},
				{
					Name:   "deposit",
					Usage:  "Sends deposits for your validator keys",
					Flags:  validatorDepositFlags,
					Action: sendDeposit,
				},
			},
		},
		{
			Name:   "version",
			Usage:  "Display the version of the LUKSO CLI Tool that is currently installed",
			Action: displayVersion,
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
