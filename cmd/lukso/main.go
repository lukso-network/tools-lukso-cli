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
	app.Usage = "Spins all lukso ecosystem components"
	app.Flags = appFlags
	app.Commands = []*cli.Command{
		{
			Name:   "install",
			Usage:  "Downloads lukso binary dependencies - needs root privileges",
			Action: installBinaries,
			Flags:  installFlags,
			Before: initializeFlags,
		},
		{
			Name: "init",
			Usage: "Initializes your lukso working directory, it's structure and configurations for all of your clients. " +
				"Make sure that you have your clients installed before initializing",
			Action: initializeDirectory,
			Flags:  networkFlags,
			Before: initializeFlags,
		},
		{
			Name:   "update",
			Usage:  "Updates all clients to newest versions",
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
			Usage:           "Start all lukso clients",
			Action:          selectNetworkFor(startClients),
			SkipFlagParsing: true,
			Flags:           startFlags,
			Before:          initializeFlags,
		},
		{
			Name:   "stop",
			Usage:  "Stops all lukso clients",
			Action: stopClients,
			Flags:  stopFlags,
		},
		{
			Name:   "log",
			Usage:  "Outputs log file of given client",
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
			Usage:  "Displays running status of clients",
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
			Usage:  "Reset data directories of all clients alongside with their log files",
			Flags:  resetFlags,
			Action: selectNetworkFor(resetClients),
		},
		{
			Name:  "validator",
			Usage: "Manage your lukso validator",
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
			Usage:  "Display version of the lukso command",
			Action: displayVersion,
		},
		{
			Name: "config",
			Action: func(c *cli.Context) error {
				cfg := config.NewConfig("./cli-config.yml")
				err := cfg.Create(gethDependencyName, prysmDependencyName)
				if err != nil {
					return err
				}

				err = cfg.WriteClientVersion(gethDependencyName, "1.11.4")

				return err
			},
		},
	}

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())

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
