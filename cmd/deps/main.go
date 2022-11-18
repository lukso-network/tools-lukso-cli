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

	binDir = "/usr/local/bin"
)

var (
	appName        = "lukso"
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
	downloadFlags = append(downloadFlags, gethInitFlags...)
	downloadFlags = append(downloadFlags, validatorInitFlags...)
	downloadFlags = append(downloadFlags, prysmInitFlags...)
	downloadFlags = append(downloadFlags, appFlags...)
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
			Before: beforeDownload,
		},
		{
			Name:   "init",
			Usage:  "Initializes your lukso working directory, it's structure and configurations for all of your clients",
			Action: downloadConfigs,
			Before: beforeDownload,
		},
	}

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())

		setupOperatingSystem()

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

func beforeDownload(ctx *cli.Context) error {
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
