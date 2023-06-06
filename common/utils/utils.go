package utils

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Exit(message string, exitCode int) error {
	log.Error(message)

	os.Exit(exitCode)

	return nil // so we can return commands with this func
}

func IsAnyRunning() bool {
	var runningClients string
	for _, client := range clients.AllClients {
		if client.IsRunning() {
			runningClients += fmt.Sprintf("- %s\n", client.Name())
		}
	}

	if runningClients == "" {
		return false
	}

	log.Warnf("⚠️  Please stop the following clients before continuing: \n%s", runningClients)

	return true
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func AcceptTermsInteractive() bool {
	message := "You are about to download the clients that are necessary to run the LUKSO Blockchain.\n" +
		"To install Prysm you are required to accept the terms of use:\n" +
		"https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md\n\n" +
		"Do you agree? [Y/n]: "

	input := RegisterInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Error("You need to type Y to continue.")
		return false
	}

	return true
}
