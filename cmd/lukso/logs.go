package main

import (
	"fmt"
	"github.com/m8b-dev/lukso-cli/pid"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
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
	log.Info("Please specify your client - run lukso logs help for more info")

	return nil
}

func logClient(dependencyName string, logFileFlag string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		logFileDir := ctx.String(logFileFlag)
		if logFileDir == "" {
			return errFlagMissing
		}

		latestFile, err := getLastFile(logFileDir, dependencyName)
		if latestFile == "" && err == nil {
			return nil
		}

		if err != nil {
			return err
		}

		return clientDependencies[dependencyName].Log(logFileDir + "/" + latestFile)
	}
}

func statClients(ctx *cli.Context) (err error) {
	err = statClient(gethDependencyName)(ctx)
	if err != nil {
		return
	}

	err = statClient(prysmDependencyName)(ctx)
	if err != nil {
		return
	}

	err = statClient(validatorDependencyName)(ctx)
	if err != nil {
		return
	}

	return
}

func statClient(dependencyName string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		isRunning := clientDependencies[dependencyName].Stat()

		if isRunning {
			pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, dependencyName)

			// since we now that process is, in fact, running (from previous checks) this shouldn't fail
			pidVal, err := pid.Load(pidLocation)
			if err != nil {
				return errProcessNotFound
			}

			log.Infof("%s: Running - PID: %d", dependencyName, pidVal)

			return nil
		}

		log.Warnf("%s: Stopped", dependencyName)

		return nil
	}
}
