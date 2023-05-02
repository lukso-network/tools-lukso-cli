package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m8b-dev/lukso-cli/pid"
	"github.com/urfave/cli/v2"
)

func (dependency *ClientDependency) Log(logFilePath string) (err error) {
	var commandName string
	var commandArgs []string
	switch systemOs {
	case ubuntu, macos:
		commandName = "tail"
		commandArgs = []string{"-f", "-n", "+0"}
	case windows:
		commandName = "type"
	default:
		commandName = "tail" // For reviewers - do we provide default command? Or omit and return with err?
		commandArgs = []string{"-f", "-n", "+0"}
	}

	command := exec.Command(commandName, append(commandArgs, logFilePath)...)

	command.Stdout = os.Stdout

	err = command.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Error("No error logs found")

		return
	}

	// error unrelated to command execution
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)
	}

	return
}

// Stat returns whether the client is running or not, along with PID
func (dependency *ClientDependency) Stat() (isRunning bool) {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependency.name)

	return pid.Exists(pidLocation)
}

func logClients(ctx *cli.Context) error {
	log.Info("⚠️  Please specify your client - run 'lukso logs --help' for more info")

	return nil
}

func logLayer(layer string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		logFileDir := ctx.String(logFolderFlag)
		if logFileDir == "" {
			return cli.Exit(fmt.Sprintf("%v- %s", errFlagMissing, logFolderFlag), 1)
		}

		if !cfg.Exists() {
			return cli.Exit(folderNotInitialized, 1)
		}

		err := cfg.Read()
		if err != nil {
			return cli.Exit(fmt.Sprintf("❌  There was an error while reading the configuration file: %v", err), 1)
		}

		var dependencyName string

		switch layer {
		case executionLayer:
			dependencyName = cfg.Execution()
		case consensusLayer:
			dependencyName = cfg.Consensus()
		case validatorLayer:
			dependencyName = validatorDependencyName
		default:
			return cli.Exit("❌  Unexpected error: unknown layer logged", 1)
		}

		latestFile, err := getLastFile(logFileDir, dependencyName)
		if latestFile == "" && err == nil {
			return nil
		}
		if err != nil {
			return cli.Exit(fmt.Sprintf("There was an error while getting the latest log file: %v", err), 1)
		}

		log.Infof("Selected layer: %s", layer)
		log.Infof("Selected client: %s", dependencyName)

		message := fmt.Sprintf("(NOTE: PATH TO FILE: %s)\nDo you want to show the log file %s?"+
			" Doing so can print lots of text on your screen [Y/n]: ", logFileDir+"/"+latestFile, latestFile)

		input := registerInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("❌  Aborting...") // there is no error, just a possible change of mind - shouldn't return err

			return nil
		}

		return clientDependencies[dependencyName].Log(logFileDir + "/" + latestFile)
	}
}

func statClients(ctx *cli.Context) (err error) {
	if !cfg.Exists() {
		return cli.Exit(folderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		log.Errorf("❌  There was an error while reading configuration file: %v", err)

		return nil
	}

	selectedExecution := cfg.Execution()
	selectedConsensus := cfg.Consensus()

	err = statClient(selectedExecution, "Execution")(ctx)
	if err != nil {
		return
	}

	err = statClient(selectedConsensus, "Consensus")(ctx)
	if err != nil {
		return
	}

	err = statClient(validatorDependencyName, "Validator")(ctx)
	if err != nil {
		return
	}

	return
}

func statClient(dependencyName, layer string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		isRunning := false
		if dependencyName != "" {
			isRunning = clientDependencies[dependencyName].Stat()
		}

		if isRunning {
			pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependencyName)

			// since we now that process is, in fact, running (from previous checks) this shouldn't fail
			pidVal, err := pid.Load(pidLocation)
			if err != nil {
				return errProcessNotFound
			}

			log.Infof("PID %d - %s (%s): Running 🟢", pidVal, layer, dependencyName)

			return nil
		}

		if dependencyName == "" {
			dependencyName = "none"
		}

		log.Warnf("PID None - %s (%s): Stopped 🔘", layer, dependencyName)

		return nil
	}
}
