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

	logsFlags = append(logsFlags, gethLogsFlags...)
	logsFlags = append(logsFlags, prysmLogsFlags...)
	logsFlags = append(logsFlags, validatorLogsFlags...)
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.Usage = "Spins all lukso ecosystem components"
	app.Flags = appFlags
	app.Commands = []*cli.Command{
		{
			Name:   "download",
			Usage:  "Downloads lukso binary dependencies - needs root privileges",
			Action: downloadBinaries,
			Flags:  downloadFlags,
			Before: initializeFlags,
		},
		{
			Name: "init",
			Usage: "Initializes your lukso working directory, it's structure and configurations for all of your clients." +
				"Make sure that you have your clients installed before initializing",
			Action: downloadConfigs,
			Flags:  downloadFlags,
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
			Name:   "start",
			Usage:  "Start all lukso clients",
			Action: startClients,
			Flags:  startFlags,
			Before: initializeFlags,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Start Geth client",
					Flags:  gethStartFlags,
					Before: initializeFlags,
					Action: startGeth,
				},
				{
					Name:   "prysm",
					Usage:  "Start Prysm client",
					Flags:  prysmStartFlags,
					Before: initializeFlags,
					Action: startPrysm,
				},
				{
					Name:   "validator",
					Usage:  "Start Validator client",
					Flags:  validatorStartFlags,
					Before: initializeFlags,
					Action: startValidator,
				},
			},
		},
		{
			Name:   "stop",
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
			Name:   "logs",
			Usage:  "Outputs log file of given client",
			Action: logClients,
			Flags:  logsFlags,
			Subcommands: []*cli.Command{
				{
					Name:   "geth",
					Usage:  "Outputs Geth client logs",
					Action: logClient(gethDependencyName, gethOutputFileFlag),
				},
				{
					Name:   "prysm",
					Usage:  "Outputs Prysm client logs",
					Action: logClient(prysmDependencyName, prysmOutputFileFlag),
				},
				{
					Name:   "validator",
					Usage:  "Outputs Validator client logs",
					Action: logClient(validatorDependencyName, validatorOutputFileFlag),
				},
			},
		},
		{
			Name:   "stat",
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
