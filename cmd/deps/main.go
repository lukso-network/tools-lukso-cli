package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	joonix "github.com/joonix/log"
	"github.com/lukso-network/lukso-orchestrator/orchestrator/node"
	"github.com/lukso-network/lukso-orchestrator/orchestrator/vanguardchain/client"
	"github.com/lukso-network/lukso-orchestrator/shared/cmd"
	"github.com/lukso-network/lukso-orchestrator/shared/journald"
	"github.com/lukso-network/lukso-orchestrator/shared/logutil"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"
	"sync"
	"time"
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
	appName               = "lukso"
	gethTag               string
	gethCommitHash        string
	validatorTag          string
	prysmTag              string
	log                   = logrus.WithField("prefix", appName)
	systemOs              string
	gethRuntimeFlags      []string
	validatorRuntimeFlags []string
	prysmRuntimeFlags     []string
)

func init() {
	allFlags := make([]cli.Flag, 0)
	allFlags = append(allFlags, gethFlags...)
	allFlags = append(allFlags, validatorFlags...)
	allFlags = append(allFlags, prysmFlags...)
	appFlags = allFlags
}

func main() {
	app := cli.App{}
	app.Name = appName
	app.Usage = "Spins all lukso ecosystem components"
	app.Flags = appFlags
	app.Action = downloadAndRunBinaries

	app.Before = func(ctx *cli.Context) error {
		format := ctx.String(cmd.LogFormat.Name)
		switch format {
		case "text":
			formatter := new(prefixed.TextFormatter)
			formatter.TimestampFormat = "2006-01-02 15:04:05"
			formatter.FullTimestamp = true
			// If persistent log files are written - we disable the log messages coloring because
			// the colors are ANSI codes and seen as gibberish in the log files.
			formatter.DisableColors = ctx.String(cmd.LogFileName.Name) != ""
			logrus.SetFormatter(formatter)
		case "fluentd":
			f := joonix.NewFormatter()
			if err := joonix.DisableTimestampFormat(f); err != nil {
				panic(err)
			}
			logrus.SetFormatter(f)
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "journald":
			if err := journald.Enable(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown log format %s", format)
		}

		logFileName := ctx.String(cmd.LogFileName.Name)
		if logFileName != "" {
			if err := logutil.ConfigurePersistentLogging(logFileName); err != nil {
				log.WithError(err).Error("Failed to configuring logging to disk.")
			}
		}

		runtime.GOMAXPROCS(runtime.NumCPU())

		setupOperatingSystem()

		// Geth related parsing
		gethTag = ctx.String(gethTagFlag)
		gethCommitHash = ctx.String(gethCommitHashFlag)
		gethRuntimeFlags = prepareGethFlags(ctx)

		// Validator related parsing
		validatorTag = ctx.String(validatorTagFlag)
		validatorRuntimeFlags = prepareValidatorFlags(ctx)

		// Prysm related parsing
		prysmTag = ctx.String(prysmTagFlag)
		prysmRuntimeFlags = preparePrysmFlags(ctx)

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

func startGeth(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", gethTag).Info("I am running genesis.json init")
	gethDataDir := ctx.String(gethDatadirFlag)
	gethGenesisArguments := []string{
		"init",
		clientDependencies[gethGenesisDependencyName].ResolveBinaryPath(gethTag, gethDataDir),
		"--datadir",
		gethDataDir,
	}

	err = clientDependencies[gethDependencyName].Run(
		gethTag,
		gethDataDir,
		gethGenesisArguments,
		ctx.Bool(gethOutputFlag),
	)

	if nil != err {
		return
	}

	time.Sleep(time.Second * 3)

	log.WithField("dependencyTag", gethTag).Info("I am running execution engine")
	err = clientDependencies[gethDependencyName].Run(
		gethTag,
		gethDataDir,
		gethRuntimeFlags,
		ctx.Bool(gethOutputFlag),
	)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	go func() {
		for {
			_, currentErr := os.Stat(cmd.DefaultGethRPCEndpoint)
			if nil == currentErr {
				log.Info("Geth ipc is up")
				waitGroup.Done()

				return
			}

			if os.IsNotExist(currentErr) {
				time.Sleep(time.Millisecond * 50)
				log.Info("Geth ipc is dead")

				continue
			}

			panic(err)
		}
	}()

	waitGroup.Wait()

	return
}

func startPrysm(ctx *cli.Context) (err error) {
	log.WithField("dependencyTag", prysmTag).Info("I am running prysm")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[prysmDependencyName].Run(
		prysmTag,
		prysmDataDir,
		prysmRuntimeFlags,
		ctx.Bool(prysmOutputFlag),
	)

	if nil != err {
		return
	}

	for {
		log.Info("prysm readiness check")
		time.Sleep(time.Millisecond * 250)
		vanClient, currentErr := client.Dial(
			ctx.Context,
			ctx.String(cmd.VanguardGRPCEndpoint.Name),
			time.Second,
			32,
			math.MaxInt32,
		)

		if nil != currentErr {
			log.WithField("cause", "prysm not ready yet").Error(currentErr)
			continue
		}

		log.Info("prysm is ready")
		vanClient.Close()
		break
	}

	return
}

func startValidator(ctx *cli.Context) (err error) {
	// First command should be to create wallet or prompt to do this by your own. This is one-time
	log.WithField("dependencyTag", validatorTag).Info("I am running prysm")
	prysmDataDir := ctx.String(prysmDatadirFlag)
	err = clientDependencies[validatorDependencyName].Run(
		validatorTag,
		prysmDataDir,
		validatorRuntimeFlags,
		false,
	)

	return
}

func startOrchestrator(ctx *cli.Context) (err error) {
	verbosity := ctx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	log.WithField("gethFlags", gethRuntimeFlags).
		WithField("prysmFlags", prysmRuntimeFlags).
		WithField("validatorFlags", validatorRuntimeFlags).Info("\n I will try to run setup with this additional flags \n")

	orchestrator, err := node.New(ctx)
	if err != nil {
		return err
	}
	orchestrator.Start()
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
