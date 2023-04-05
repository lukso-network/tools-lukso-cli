package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// prepareTimestampedFile concatenates directory of logs, prefix of file name and timestamp.
// It also creates directory for logs of given client.
// This directory should be given without prefixing/suffixing slashes.
func prepareTimestampedFile(logDir, logFileName string) (logFile string, err error) {
	err = os.MkdirAll(logDir, 0750)
	if err != nil {
		return
	}

	t := time.Now().Format("2006-01-02_15-04-05")

	logFile = fmt.Sprintf("%s/%s_%s.log", logDir, logFileName, t)

	return
}

func prepareLogfileFlag(logDir, dependencyName string) string {
	prysmFullLogPath, err := prepareTimestampedFile(logDir, dependencyName)
	if err != nil {
		log.Warnf("Couldn't prepare log folder for %s client. Warning: continuing without log files being saved", dependencyName)

		return ""
	}

	return fmt.Sprintf("--log-file=%s", prysmFullLogPath)
}

// The lint says this function is unused
// I'll ignore it for now as I'm missing context
func createJwtSecret(dest string) error { //nolint:all
	log.Info("Creating new JWT secret")
	jwtDir := truncateFileFromDir(dest)

	err := os.MkdirAll(jwtDir, 0750)
	if err != nil {
		return err
	}

	secretBytes := make([]byte, 32)

	_, err = rand.Read(secretBytes)
	if err != nil {
		return err
	}

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

	return strings.Join(segments[:len(segments)-1], "/")
}

// getLastFile returns name of last file from given directory in alphabetical order.
// In case of log files last one is also the newest (format increments similar to typical number - YYYY-MM-DD_HH:MM:SS)
func getLastFile(dir string, dependency string) (string, error) {
	var (
		commandName string
		lsBuf       = new(bytes.Buffer)
		grepBuf     = new(bytes.Buffer)
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
	command.Stdout = lsBuf

	grepCommand := exec.Command("grep", dependency)
	grepCommand.Stdin = lsBuf
	grepCommand.Stdout = grepBuf

	err := command.Run()
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)

		return "", err
	}

	err = grepCommand.Run()
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)

		return "", err
	}

	scan := bufio.NewScanner(grepBuf)
	for scan.Scan() {
		files = append(files, scan.Text())
	}
	if len(files) < 1 {
		log.Infof("No log files found in %s", dir)

		return "", nil
	}

	lastFile := files[len(files)-1]

	message := fmt.Sprintf("(NOTE: PATH TO FILE: %s)\nDo you want to show log file %s?"+
		" Doing so can print lots of text on your screen [Y/n]: ", dir+"/"+lastFile, lastFile)

	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Info("Aborting...") // there is no error, just a possible change of mind - shouldn't return err

		return "", nil
	}

	return lastFile, nil
}

func isRunning(dependency string) bool {
	isRunning := clientDependencies[dependency].Stat()

	return isRunning
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func registerInputWithMessage(message string) (input string) {
	fmt.Print(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

// parseFlags takes care of parsing flags that are skipped if SkipFlagParsing is set to true - if --help or -h is found we display help and stop execution
func parseFlags(ctx *cli.Context) (err error) {
	args := ctx.Args()
	argsLen := args.Len()
	for i := 0; i < argsLen; i++ {
		arg := args.Get(i)

		if arg == "--help" || arg == "-h" {
			cli.ShowSubcommandHelpAndExit(ctx, 0)
		}

		if strings.HasPrefix(arg, "--") {
			if i+1 == argsLen {
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, "true")
				if err != nil && strings.Contains(err.Error(), noSuchFlag) {
					err = nil

					return
				}

				return
			}

			// we found a flag for our client - now we need to check if it's a value or bool flag
			nextArg := args.Get(i + 1)
			if strings.HasPrefix(nextArg, "--") { // we found a next flag, so current one is a bool
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, "true")
				if err == nil || (err != nil && strings.Contains(err.Error(), noSuchFlag)) {
					err = nil

					continue
				}

				return
			} else {
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, nextArg)
				if err == nil || (err != nil && strings.Contains(err.Error(), noSuchFlag)) {
					err = nil

					continue
				}

				return
			}
		}
	}

	return
}

func isRoot() (isRoot bool, err error) {
	command := exec.Command("id", "-u")

	output, err := command.Output()
	if err != nil {
		return
	}

	usr, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		return
	}

	if usr == 0 {
		return true, nil
	}

	return false, nil
}
