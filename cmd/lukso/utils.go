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

func createJwtSecret(dest string) error {
	log.Info("🔄  Creating new JWT secret")
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

	log.Infof("✅  New JWT secret created in %s", dest)

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
		log.Errorf("❌  There was an error while executing command: %s. Error: %v", commandName, err)

		return "", err
	}

	err = grepCommand.Run()
	if err != nil {
		log.Errorf("❌  There was an error while executing command: %s. Error: %v", commandName, err)

		return "", err
	}

	scan := bufio.NewScanner(grepBuf)
	for scan.Scan() {
		files = append(files, scan.Text())
	}
	if len(files) < 1 {
		log.Infof("❗️  No log files found in %s", dir)

		return "", nil
	}

	lastFile := files[len(files)-1]

	return lastFile, nil
}

func isRunning(dependency string) bool {
	isRunning := clientDependencies[dependency].Stat()

	return isRunning
}

func isAnyRunning() bool {
	gethRunning := isRunning(gethDependencyName)
	erigonRunning := isRunning(erigonDependencyName)
	prysmRunning := isRunning(prysmDependencyName)
	lighthouseRunning := isRunning(lighthouseDependencyName)
	validatorRunning := isRunning(validatorDependencyName)

	if gethRunning || prysmRunning || validatorRunning || erigonRunning || lighthouseRunning {
		message := "⚠️  Please stop the following clients before continuing: "
		if gethRunning {
			message += "geth "
		}
		if erigonRunning {
			message += "erigon "
		}
		if prysmRunning {
			message += "prysm "
		}
		if lighthouseRunning {
			message += "lighthouse "
		}
		if validatorRunning {
			message += "validator "
		}

		message += "\n➡️  You can use 'lukso stop' to stop clients."

		log.Warn(message)

		return true
	}

	return false
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

// flagFileExists check whether a path under given flag exists
func flagFileExists(ctx *cli.Context, flag string) bool {
	flagPath := ctx.String(flag)
	if !fileExists(flagPath) {
		log.Errorf("⚠️  Path '%s' in --%s flag doesn't exist - please provide a valid file path", flagPath, flag)

		return false
	}

	return true
}
