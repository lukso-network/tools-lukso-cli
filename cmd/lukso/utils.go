package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

// prepareTimestampedFile concatenates directory of logs, prefix of file name and timestamp.
// It also creates directory for logs of given client.
// This directory should be given without prefixing/suffixing slashes.
func prepareTimestampedFile(logDir, logFileName string) (logFile string, err error) {
	err = os.MkdirAll(logDir, 0750)
	if err != nil {
		return
	}

	t := time.Now().Format("02-01-2006_15:04:05")

	logFile = fmt.Sprintf("%s/%s_%s.log", logDir, logFileName, t)

	return
}

func prepareLogfileFlag(ctx *cli.Context, logDirFlag, dependencyName string) string {
	logDir := ctx.String(logDirFlag)
	prysmFullLogPath, err := prepareTimestampedFile(logDir, dependencyName)
	if err != nil {
		log.Warnf("Couldn't prepare log folder for %s client. Warning: continuing without log files being saved", dependencyName)

		return ""
	}

	return fmt.Sprintf("--log-file=%s", prysmFullLogPath)
}
