package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
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

	t := time.Now().Format("2006-01-02_15:04:05")

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

// getLastFile returns name of last file from given directory in alphabetical order.
// In case of log files last one is also the newest (format increments similar to typical number - YYYY-MM-DD_HH:MM:SS)
func getLastFile(dir string) (string, error) {
	var (
		commandName string
		buf         = new(bytes.Buffer)
		files       []string
	)

	switch systemOs {
	case ubuntu, macos:
		commandName = "ls"
	case windows:
		commandName = "type"
	default:
		commandName = "ls"
	}

	command := exec.Command(commandName, dir)

	command.Stdout = buf

	err := command.Run()
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)

		return "", err
	}

	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		files = append(files, scan.Text())
	}
	if len(files) < 1 {
		log.Infof("No log files found in %s", dir)

		return "", nil
	}

	lastFile := files[len(files)-1]

	log.Infof("(NOTE: PATH TO FILE: %s)\nDo you want to show log file %s?"+
		" Doing so can print lots of text on your screen [Y/n]: ", dir+"/"+lastFile, lastFile)

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	input := scanner.Text()
	if input != "Y" {
		log.Info("Aborting...") // there is no error, just a possible change of mind - shouldn't return err

		return "", nil
	}

	return lastFile, nil
}

func isRunning(dependency string) (bool, error) {
	return clientDependencies[dependency].Stat()
}
