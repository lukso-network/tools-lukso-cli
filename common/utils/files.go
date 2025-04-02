package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"

	"github.com/lukso-network/tools-lukso-cli/common"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateJwtSecret() (jwt []byte, err error) {
	secretBytes := make([]byte, 32)

	_, err = rand.Read(secretBytes)
	if err != nil {
		return
	}

	jwt = []byte(hex.EncodeToString(secretBytes))

	return
}

func PrepareTimestampedFile(logDir, logFileName string) (logFile string, err error) {
	err = os.MkdirAll(logDir, common.ConfigPerms)
	if err != nil {
		return
	}

	t := time.Now().Format("2006-01-02_15-04-05")

	logFile = fmt.Sprintf("%s/%s_%s.log", logDir, logFileName, t)

	return
}

// GetLastFile returns name of last file from given directory in alphabetical order.
// In case of log files last one is also the newest (format increments similar to typical number - YYYY-MM-DD_HH:MM:SS)
func GetLastFile(dir string, dependency string) (string, error) {
	var (
		commandName string
		lsBuf       = new(bytes.Buffer)
		grepBuf     = new(bytes.Buffer)
		files       []string
	)

	commandName = "ls"

	command := exec.Command(commandName, dir)
	command.Stdout = lsBuf

	grepCommand := exec.Command("grep", dependency)
	grepCommand.Stdin = lsBuf
	grepCommand.Stdout = grepBuf

	err := command.Run()
	if err != nil {
		log.Error("❌  Unable to retrieve logs for specified network. Did you forget to add network flag? i.e. --testnet")

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

// TruncateFileFromDir removes file name from its path.
// Example: /path/to/foo/foo.txt => /path/to/foo
func TruncateFileFromDir(filePath string) string {
	segments := strings.Split(filePath, "/")

	return strings.Join(segments[:len(segments)-1], "/")
}

// FlagFileExists check whether a path under given flag exists
func FlagFileExists(ctx *cli.Context, flag string) bool {
	flagPath := ctx.String(flag)
	if !FileExists(flagPath) {
		log.Errorf("⚠️  Path '%s' in --%s flag doesn't exist - please provide a valid file path", flagPath, flag)

		return false
	}

	return true
}

// ReadValidatorPassword is responsible for creating a secure way to
// pass a validator password from lukso CLI to a client of user's choosing
func ReadValidatorPassword(ctx *cli.Context) (f *os.File, err error) {
	var password []byte
	fmt.Print("\nPlease enter your wallet password: ")
	password, err = term.ReadPassword(0)
	fmt.Println("")

	if err != nil {
		err = Exit(fmt.Sprintf("❌  Couldn't read password: %v", err), 1)

		return
	}

	b := make([]byte, 4)
	_, err = rand.Read(b)
	if err != nil {
		err = Exit(fmt.Sprintf("❌  Couldn't create random byte array: %v", err), 1)

		return
	}

	randPipe := hex.EncodeToString(b)

	passwordPipePath := fmt.Sprintf("./.%s", randPipe)
	err = syscall.Mkfifo(passwordPipePath, 0o600)
	if err != nil {
		err = Exit(fmt.Sprintf("❌  Couldn't create password pipe: %v", err), 1)

		return
	}

	f, err = os.OpenFile(passwordPipePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		err = Exit(fmt.Sprintf("❌  Couldn't open password pipe: %v", err), 1)

		return
	}
	_, err = f.Write(password)
	if err != nil {
		err = Exit(fmt.Sprintf("❌  Couldn't write password to pipe: %v", err), 1)

		return
	}

	return
}
