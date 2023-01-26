package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
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

func createJwtSecret(dest string) error {
	log.Info("Creating new JWT secret")
	jwtDir := truncateFileFromDir(dest)

	err := os.MkdirAll(jwtDir, 0750)
	if err != nil {
		return err
	}

	secretBytes := make([]byte, 32)

	_, err = rand.Read(secretBytes)

	err = os.WriteFile(dest, []byte(hex.EncodeToString(secretBytes)), configPerms)
	if err != nil {
		return err
	}

	log.Infof("New JWT secret created in %s", dest)

	return nil
}

// truncateFileFromDir removes file name from its path.
// Example: /path/to/foo/foo.txt => /path/to/foo
func truncateFileFromDir(filePath string) string {
	segments := strings.Split(filePath, "/")

	return strings.TrimRight(filePath, segments[len(segments)-1])
}
