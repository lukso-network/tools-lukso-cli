package clients

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type TekuValidatorClient struct {
	*clientBinary
}

func NewTekuValidatorClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *TekuValidatorClient {
	return &TekuValidatorClient{
		&clientBinary{
			name:           tekuValidatorDependencyName,
			fileName:       "validator_tk",
			baseUrl:        "",
			githubLocation: tekuGithubLocation,
			buildInfo:      tekuBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var TekuValidator ValidatorBinaryDependency

var _ ValidatorBinaryDependency = &TekuValidatorClient{}

func (t *TekuValidatorClient) Install(version string, isUpdate bool) error {
	return nil
}

func (t *TekuValidatorClient) Update() (err error) {
	return nil
}

func (t *TekuValidatorClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = t.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.TekuValidatorConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--validators-proposer-default-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))

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
	command := exec.Command(fmt.Sprintf("./%s/bin/teku", tekuFolder), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.TimestampedFile(logFolder, t.FileName())
	if err != nil {
		return
	}

	err = os.WriteFile(fullPath, []byte{}, 0o750)
	if err != nil {
		return
	}

	logFile, err = os.OpenFile(fullPath, os.O_RDWR, 0o750)
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

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, t.FileName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", t.Name())

	return
}

func (t *TekuValidatorClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return
}

func (t *TekuValidatorClient) Version() (version string) {
	return Teku.Version()
}

func (t *TekuValidatorClient) Import(ctx *cli.Context) (err error) {
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	message := fmt.Sprintf("To run a Teku validator the LUKSO CLI needs to prepare a separate password file for each of your validator key files."+
		"Please be aware that this process will create a lot of unencrypted password files in %s folder.\n"+
		"Do you want to continue? [y/N]\n>", walletDir)

	input := utils.RegisterInputWithMessage(message)

	if !strings.EqualFold(input, "y") {
		log.Info("❌  Aborting...")

		return
	}

	err = os.MkdirAll(walletDir, 0o750)
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
	if walletDir == "" {
		return utils.Exit("❌  Wallet directory not provided - please provide a --validator-wallet-dir flag containing your keys directory", 1)
	}

	err = keystoreListWalk(walletDir)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while list validators: %v", err), 1)
	}

	return
}

func (t *TekuValidatorClient) Exit(ctx *cli.Context) (err error) {
	wallet := ctx.String(flags.ValidatorWalletDirFlag)
	if wallet == "" {
		return utils.Exit("❌  Wallet directory not provided - please provide a --validator-wallet-dir flag containing your keys directory", 1)
	}

	if !utils.FileExists(wallet) {
		return utils.Exit("❌  Wallet directory missing - please provide a --validator-wallet-dir flag containing your keys directory or use a network flag", 1)
	}

	args := []string{"voluntary-exit", "--validator-keys", fmt.Sprintf("%s:%s", wallet, wallet)}

	exitCommand := exec.Command(fmt.Sprintf("./%s/bin/teku", tekuFolder), args...)

	exitCommand.Stdout = os.Stdout
	exitCommand.Stderr = os.Stderr
	exitCommand.Stdin = os.Stdin

	err = exitCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while exiting validator: %v", err), 1)
	}

	return
}
