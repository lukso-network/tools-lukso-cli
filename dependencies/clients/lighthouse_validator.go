package clients

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type LighthouseValidatorClient struct {
	*clientBinary
}

func NewLighthouseValidatorClient() *LighthouseValidatorClient {
	return &LighthouseValidatorClient{
		&clientBinary{
			name:           lighthouseValidatorDependencyName,
			commandName:    "validator_lh", // we run it using lighthouse bin
			baseUrl:        "",             // no separate client for lighthouse validator - lighthouse_beacon for reference
			githubLocation: "",
		},
	}
}

var LighthouseValidator = NewLighthouseValidatorClient()

var _ ValidatorBinaryDependency = &LighthouseValidatorClient{}

func (l *LighthouseValidatorClient) Update() error {
	return nil
}

func (l *LighthouseValidatorClient) Start(ctx *cli.Context, args []string) (err error) {
	if l.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", l.Name())

		err = l.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", l.Name())
	}

	command := exec.Command(Lighthouse.commandName, args...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, l.CommandName())
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

	log.Infof("🔄  Starting %s", l.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, l.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", l.Name())

	return
}

func (l *LighthouseValidatorClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	validatorConfigExists := utils.FlagFileExists(ctx, flags.ValidatorConfigFileFlag)
	chainConfigExists := utils.FlagFileExists(ctx, flags.PrysmChainConfigFileFlag)
	if !validatorConfigExists || !chainConfigExists {
		err = errors.ErrFlagPathInvalid

		return
	}

	defaults, err := config.LoadLighthouseConfig(ctx.String(flags.LighthouseValidatorConfigFileFlag))
	if err != nil {
		return
	}

	defaults = append(defaults, "--suggested-fee-recipient", ctx.String(flags.TransactionFeeRecipientFlag))

	userFlags := l.ParseUserFlags(ctx)

	startFlags = mergeFlags(userFlags, defaults)

	startFlags = append([]string{"vc"}, startFlags...)

	return
}

func (l *LighthouseValidatorClient) Import(ctx *cli.Context) (err error) {
	args := []string{
		"am",
		"validator",
		"import",
		"--directory",
		ctx.String(flags.ValidatorKeysFlag),
		"--datadir",
		ctx.String(flags.ValidatorWalletDirFlag),
	}

	passwordFile := ctx.String(flags.ValidatorPasswordFlag)
	if passwordFile != "" {
		args = append(args, "--password-file", passwordFile, "--reuse-password")
	}

	initCommand := exec.Command(Lighthouse.CommandName(), args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err = initCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while importing validator accounts: %v", err), 1)
	}

	return nil
}

func (l *LighthouseValidatorClient) List(ctx *cli.Context) (err error) {
	return
}
