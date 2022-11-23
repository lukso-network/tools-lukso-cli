package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

func (dependency *ClientDependency) Run(
	commandName string,
	arguments []string,
	attachStdoutAndErr bool,
) (err error) {
	command := exec.Command(commandName, arguments...)

	if attachStdoutAndErr {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}

	err = command.Start()

	return
}

func startClients(ctx *cli.Context) error {
	log.Info("Starting all clients")

	err := startGeth(ctx)

	return err
}

func startGeth(ctx *cli.Context) error {
	log.Info("Starting Geth")

	err := clientDependencies[gethDependencyName].Run(gethDependencyName, []string{}, false)
	if err != nil {
		return err
	}

	log.Info("Geth started! Use lukso logs command to see logs")
	return nil
}
