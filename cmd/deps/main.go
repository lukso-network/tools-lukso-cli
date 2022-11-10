package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"
)

// We want to spin also 3 libraries at once, and secretly rule them by orchestrator. It matches for me somehow

// This binary will also support only some of the possible networks.
// Make a pull request to attach your network.
// We are also very open to any improvements. Please make some issue or hackmd proposal to make it better.
// Join our lukso discord https://discord.gg/E2rJPP4 to ask some questions

const (
	ubuntu  = "linux"
	macos   = "darwin"
	windows = "windows"
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

func init() {
	allFlags := make([]cli.Flag, 0)
	allFlags = append(allFlags, gethFlags...)
	allFlags = append(allFlags, validatorFlags...)
	allFlags = append(allFlags, prysmFlags...)
	appFlags = append(appFlags, allFlags...)
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.Usage = "Spins all lukso ecosystem components"
	app.Flags = appFlags
	app.Action = downloadAndRunBinaries

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

func downloadAndRunBinaries(ctx *cli.Context) (err error) {
	// Get os, then download all binaries into datadir matching desired system
	// After successful download run binary with desired arguments spin and connect them
	// Orchestrator can be run from-memory
	err = downloadGenesis(ctx)

	if nil != err {
		return
	}

	err = downloadGeth(ctx)

	if nil != err {
		return
	}

	err = downloadValidator(ctx)

	if nil != err {
		return
	}

	err = downloadPrysm(ctx)

	if nil != err {
		return
	}

	err = downloadConfig(ctx)

	return
}

func downloadGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("I am downloading Geth")
	gethDataDir := ctx.String(gethDatadirFlag)
	err = clientDependencies[gethDependencyName].Download(gethTag, gethDataDir, gethCommitHash)

	return
}

func downloadGenesis(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("I am downloading geth genesis")
	gethDataDir := ctx.String(gethDatadirFlag)
	err = clientDependencies[gethGenesisDependencyName].Download(gethTag, gethDataDir, "")

	if nil != err {
		return
	}

	log.WithField("dependencyTag", prysmTag).Info("I am downloading prysm genesis")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmGenesisDependencyName].Download(prysmTag, prysmDataDir, "")

	return
}

func downloadConfig(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("I am downloading prysm config")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmConfigDependencyName].Download(prysmTag, prysmDataDir, "")

	return
}

func downloadPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("I am downloading prysm")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmDependencyName].Download(prysmTag, prysmDataDir, "")

	return
}

func downloadValidator(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", validatorTag).Info("I am downloading validator")
	validatorDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[validatorDependencyName].Download(prysmTag, validatorDataDir, "")

	return
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
