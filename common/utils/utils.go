package utils

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Exit(message string, exitCode int) error {
	log.Error(message)

	os.Exit(exitCode)

	return nil // so we can return commands with this func
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
