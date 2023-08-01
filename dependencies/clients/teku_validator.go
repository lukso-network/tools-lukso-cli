package clients

import (
	"encoding/json"
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/pid"
	"golang.org/x/term"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type TekuValidatorClient struct {
	*clientBinary
}

func NewTekuValidatorClient() *TekuValidatorClient {
	return &TekuValidatorClient{
		&clientBinary{
			name:           tekuValidatorDependencyName,
			commandName:    "validator_tk",
			baseUrl:        "",
			githubLocation: tekuGithubLocation,
		},
	}
}

var TekuValidator = NewTekuValidatorClient()

var _ ValidatorBinaryDependency = &TekuValidatorClient{}

func (t *TekuValidatorClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.TekuValidatorConfigFileFlag)))

	return
}

func (t *TekuValidatorClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if t.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", t.Name())

		err = t.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", t.Name())
	}

	arguments = append([]string{"vc"}, arguments...)
	command := exec.Command(fmt.Sprintf("./%s/%s/bin/teku", tekuDepsFolder, tekuFolder), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, t.CommandName())
	if err != nil {
		return
	}

	err = os.WriteFile(fullPath, []byte{}, 0750)
	if err != nil {
		return
	}

	logFile, err = os.OpenFile(fullPath, os.O_RDWR, 0750)
	if err != nil {
		return
	}

	command.Stdout = logFile
	command.Stderr = logFile

	log.Infof("🔄  Starting %s", t.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, t.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", t.Name())

	return
}

func (t *TekuValidatorClient) Import(ctx *cli.Context) (err error) {
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	message := fmt.Sprintf("To run Teku validator LUKSO CLI needs to prepare password file for each one of your keystores. "+
		"Please be aware that this process will produce a lot of unencrypted password files in %s folder.\n"+
		"Do you want to continue? [y/N]\n>", walletDir)

	input := utils.RegisterInputWithMessage(message)

	if !strings.EqualFold(input, "y") {
		log.Info("❌  Aborting...")

		return
	}

	err = os.MkdirAll(walletDir, 0750)
	if err != nil {
		return
	}

	fmt.Print("Please enter your keystore password: ")
	passwordBytes, err := term.ReadPassword(0)
	if err != nil {
		return
	}
	fmt.Println()

	password := string(passwordBytes)

	keysImported := 0
	passwordsImported := 0

	keysDir := ctx.String(flags.ValidatorKeysFlag)
	walkFunc := func(path string, d fs.DirEntry, _ error) (err error) {
		if !strings.Contains(path, "keystore") || d.IsDir() {
			return
		}

		// open all needed files
		var keystoreFile, newKeystoreFile, passwordFile *os.File

		info, err := d.Info()
		if err != nil {
			log.Warnf("Couldn't get keystore info at %s - skipping...", path)

			return nil
		}
		keystoreExt := filepath.Ext(d.Name())
		if !strings.Contains(keystoreExt, "json") {
			log.Warnf("Couldn't copy keystore at %s (not json extension) - skipping...", path)

			return nil
		}

		keystoreFile, err = os.OpenFile(path, os.O_RDONLY, info.Mode())
		if err != nil {
			log.Warnf("Couldn't open keystore at %s - skipping...", path)

			return nil
		}
		defer keystoreFile.Close()

		newKeystorePath := fmt.Sprintf("%s/%s", walletDir, d.Name())
		if utils.FileExists(newKeystorePath) {
			log.Warnf("Keystore at %s already exists - continuing...", newKeystorePath)

			return nil
		}

		newKeystoreFile, err = os.OpenFile(newKeystorePath, os.O_CREATE|os.O_RDWR, info.Mode())
		if err != nil {
			log.Warnf("Couldn't create new keystore at %s - skipping...", newKeystorePath)

			return nil
		}
		defer newKeystoreFile.Close()

		passwordFileName := fmt.Sprintf("%s.txt", strings.TrimSuffix(d.Name(), keystoreExt))
		passwordFilePath := fmt.Sprintf("%s/%s", walletDir, passwordFileName)
		passwordFile, err = os.OpenFile(passwordFilePath, os.O_CREATE|os.O_RDWR, info.Mode())
		if err != nil {
			log.Warnf("Couldn't create new password file at %s - skipping...", passwordFilePath)

			return nil
		}
		defer passwordFile.Close()

		// copy data from old keystores to new ones and write password to password file
		_, err = io.Copy(newKeystoreFile, keystoreFile)
		if err != nil {
			log.Warnf("Couldn't create new keystore file at %s - skipping...", newKeystorePath)

			return nil
		}

		keysImported++

		_, err = passwordFile.Write([]byte(password))
		if err != nil {
			log.Warnf("Couldn't create new password file at %s - skipping...", passwordFilePath)

			return nil
		}

		passwordsImported++

		return
	}

	err = filepath.WalkDir(keysDir, walkFunc)
	if err != nil {
		return
	}

	log.Infof("✅  Successfully imported %d keystores and %d passwords!", keysImported, passwordsImported)
	log.Info("✅  To list your imported keys run: lukso validator list")

	return
}

func (t *TekuValidatorClient) List(ctx *cli.Context) (err error) {
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	validatorIndex := 0

	walkFunc := func(path string, d fs.DirEntry, entryError error) (err error) {
		if !strings.Contains(d.Name(), "json") {
			return nil
		}

		keystoreFile, err := os.Open(path)
		if err != nil {
			return
		}
		defer keystoreFile.Close()

		keystoreFileBytes, err := io.ReadAll(keystoreFile)
		if err != nil {
			return
		}

		keystore := struct {
			Pubkey string `json:"pubkey"`
		}{}

		err = json.Unmarshal(keystoreFileBytes, &keystore)
		if err != nil {
			return
		}

		log.Infof("Validator #%d: %s", validatorIndex, keystore.Pubkey)
		validatorIndex++

		return
	}

	err = filepath.WalkDir(walletDir, walkFunc)
	if err != nil {
		return
	}

	return
}

func (t *TekuValidatorClient) Exit(ctx *cli.Context) error {
	//TODO implement me
	panic("implement me")
}
