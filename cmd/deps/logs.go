package main

import (
	"bufio"
	"bytes"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

func (dependency *ClientDependency) Log(logFileDir string) (err error) {
	var commandName string
	switch systemOs {
	case ubuntu, macos:
		commandName = "cat"
	case windows:
		commandName = "type"
	default:
		commandName = "cat" // For reviewers - do we provide default command? Or omit and return with err?
	}

	command := exec.Command(commandName, logFileDir)

	command.Stdout = os.Stdout

	err = command.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Error("No error logs found")

		return
	}

	// error unrelated to command execution
	if err != nil {
		log.Errorf("There was an error while executing logs command. Error: %v", err)
	}

	return
}

// Stat returns whether the client is running or not
func (dependency *ClientDependency) Stat() (err error) {
	var (
		commandName string
		buf         = new(bytes.Buffer)
	)

	switch systemOs {
	case ubuntu, macos:
		commandName = "ps"
	case windows:
		commandName = "tasklist"
	default:
		commandName = "ps"
	}

	command := exec.Command(commandName)
	command.Stdout = buf

	err = command.Run()
	if err != nil {
		log.Errorf("There was an error while executing logs command. Error: %v", err)

		return
	}

	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		if strings.Contains(scan.Text(), dependency.name) {
			log.Infof("%s: Up and running", dependency.name)

			return
		}
	}

	log.Warnf("%s: inactive", dependency.name)

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
		return clientDependencies[dependencyName].Stat()
	}
}
