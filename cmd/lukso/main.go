package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"
)

const (
	ubuntu  = "linux"
	macos   = "darwin"
	windows = "windows"

	unixBinDir = "/usr/local/bin"
	// should be a user-created path, like C:\bin,
	// but since it is not guaranteed that all users vahe it we can leave it as is
	windowsBinDir = "/Windows/System32"

	// command names
	installCommand = "install"
	updateCommand  = "update"
	startCommand   = "start"
	stopCommand    = "stop"
	initCommand    = "init"
	logCommand     = "log"
	statusCommand  = "status"
	resetCommand   = "reset"
)

var (
	appName        = "lukso"
	binDir         string
	gethTag        string
	gethCommitHash string
	validatorTag   string
	prysmTag       string
	log            = logrus.WithField("prefix", appName)
	systemOs       string
)

// do not confuse init func with init subcommand - for now we just pass init flags, there are more to come
func init() {
	downloadFlags = make([]cli.Flag, 0)
	downloadFlags = append(downloadFlags, gethDownloadFlags...)
	downloadFlags = append(downloadFlags, validatorDownloadFlags...)
	downloadFlags = append(downloadFlags, prysmDownloadFlags...)
	downloadFlags = append(downloadFlags, appFlags...)

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
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.Usage = "Spins all lukso ecosystem components"
	app.Flags = appFlags
	app.Commands = []*cli.Command{
		{
			Name:   installCommand,
			Usage:  "Downloads lukso binary dependencies - needs root privileges",
			Action: downloadBinaries,
			Flags:  downloadFlags,
			Before: initializeFlags,
		},
		{
			Name: initCommand,
			Usage: "Initializes your lukso working directory, it's structure and configurations for all of your clients. " +
				"Make sure that you have your clients installed before initializing",
			Action: downloadConfigs,
			Flags:  networkFlags,
			Before: selectNetworkFor(initializeFlags),
		},
		{
			Name:   updateCommand,
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
			Name:   startCommand,
			Usage:  "Start all lukso clients",
			Action: selectNetworkFor(startClients),
			Flags:  startFlags,
			Before: initializeFlags,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Start Geth client",
					Flags:  gethStartFlags,
					Before: initializeFlags,
					Action: selectNetworkFor(startGeth),
				},
				{
					Name:   "prysm",
					Usage:  "Start Prysm client",
					Flags:  prysmStartFlags,
					Before: initializeFlags,
					Action: selectNetworkFor(startPrysm),
				},
				{
					Name:   "validator",
					Usage:  "Start Validator client",
					Flags:  validatorStartFlags,
					Before: initializeFlags,
					Action: selectNetworkFor(startValidator),
				},
			},
		},
		{
			Name:   stopCommand,
			Usage:  "Stops all lukso clients",
			Action: stopClients,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Stop Geth client",
					Action: stopClient(clientDependencies[gethDependencyName]),
				},
				{
					Name:   "prysm",
					Usage:  "Stop Prysm client",
					Action: stopClient(clientDependencies[prysmDependencyName]),
				},
				{
					Name:   "validator",
					Usage:  "Stop Validator client",
					Action: stopClient(clientDependencies[validatorDependencyName]),
				},
			},
		},
		{
			Name:   logCommand,
			Usage:  "Outputs log file of given client",
			Action: logClients,
			Flags:  logsFlags,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Outputs Geth client logs",
					Flags:  networkFlags,
					Action: selectNetworkFor(logClient(gethDependencyName, gethLogDirFlag)),
				},
				{
					Name:   "prysm",
					Usage:  "Outputs Prysm client logs",
					Flags:  networkFlags,
					Action: selectNetworkFor(logClient(prysmDependencyName, prysmLogDirFlag)),
				},
				{
					Name:   "validator",
					Usage:  "Outputs Validator client logs",
					Flags:  networkFlags,
					Action: selectNetworkFor(logClient(validatorDependencyName, validatorLogDirFlag)),
				},
			},
		},
		{
			Name:   statusCommand,
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
			Name:   resetCommand,
			Usage:  "Reset data directories of all clients alongside with their log files",
			Flags:  resetFlags,
			Action: selectNetworkFor(resetClients),
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

		log = log.WithField("os", systemOs)

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
