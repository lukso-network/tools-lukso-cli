package main

import (
	"github.com/m8b-dev/lukso-cli/pid"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

// initializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func initializeDirectory(ctx *cli.Context) error {
	for _, dependency := range clientDependencies {
		// this logic may fail when folder structure changes, but this shouldn't be the case
		if !strings.Contains(dependency.filePath, "./config") {
			continue
		}

		err := dependency.Download("", "", false, configPerms)
		if err != nil {
			log.Errorf("There was error while downloading %s file: %v", dependency.name, err.Error())
		}

	}

	err := os.MkdirAll(pid.FileDir, configPerms)
	if err != nil {
		log.Errorf("There was an error while preparing PID directory: %v", err)

		return err
	}

	log.Info("Folder initialized! To start your node run lukso start command")

	return nil
}
