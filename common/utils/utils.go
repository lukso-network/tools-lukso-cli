package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func Exit(message string, exitCode int) error {
	log.Error(message)

	os.Exit(exitCode)

	return nil // so we can return commands with this func
}

func IsAnyRunning() bool {
	//gethRunning := isRunning(gethDependencyName)
	//erigonRunning := isRunning(erigonDependencyName)
	//prysmRunning := isRunning(prysmDependencyName)
	//lighthouseRunning := isRunning(lighthouseDependencyName)
	//validatorRunning := isRunning(validatorDependencyName)
	//
	//if gethRunning || prysmRunning || validatorRunning || erigonRunning || lighthouseRunning {
	//	message := "⚠️  Please stop the following clients before continuing: "
	//	if gethRunning {
	//		message += "geth "
	//	}
	//	if erigonRunning {
	//		message += "erigon "
	//	}
	//	if prysmRunning {
	//		message += "prysm "
	//	}
	//	if lighthouseRunning {
	//		message += "lighthouse "
	//	}
	//	if validatorRunning {
	//		message += "validator "
	//	}
	//
	//	message += "\n➡️  You can use 'lukso stop' to stop clients."
	//
	//	log.Warn(message)
	//
	//	return true
	//}
	//
	//return false

	return false
}
