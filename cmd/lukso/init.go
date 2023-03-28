package main

import (
	"github.com/urfave/cli/v2"
	"strings"
)

// initializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func initializeDirectory(ctx *cli.Context) error {
	for _, dependency := range clientDependencies {
		// this logci may fail when folder structure changes, but this shouldn't be the case
		if !strings.Contains(dependency.filePath, "./config") {
			continue
		}

		err := dependency.Download("", "", false, configPerms)
		if err != nil {
			log.Errorf("There was error while downloading %s file: %v", dependency.name, err.Error())
		}

	}

	log.Info("Folder initialized! To start your node run lukso start command")

	return nil
}
