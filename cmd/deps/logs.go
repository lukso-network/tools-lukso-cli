package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

func (dependency *ClientDependency) Log(logFileDir string) (err error) {
	command := exec.Command("cat", logFileDir)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()

	return
}

func logClients(ctx *cli.Context) error {
	log.Info("Please specify your client - run lukso logs help for more info")

	return nil
}

func logClient(dependencyName string, logFileFlag string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		logFileDir := ctx.String(logFileFlag)
		if logFileDir == "" {
			return errorFlagMissing
		}

		return clientDependencies[dependencyName].Log(logFileDir)
	}
}
