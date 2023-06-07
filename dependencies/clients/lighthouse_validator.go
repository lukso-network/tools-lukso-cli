package clients

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

type LighthouseValidatorClient struct {
	*clientBinary
}

func NewLighthouseValidatorClient() *LighthouseValidatorClient {
	return &LighthouseValidatorClient{
		&clientBinary{
			name:           lighthouseValidatorDependencyName,
			commandName:    "lighthouse_validator", // we run it using lighthouse bin
			baseUrl:        "",                     // no separate client for lighthouse validator - lighthouse_beacon for reference
			githubLocation: "",
		},
	}
}

var LighthouseValidator = NewLighthouseValidatorClient()

var _ ValidatorBinaryDependency = &LighthouseValidatorClient{}

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

	logFilePath, err := utils.PrepareTimestampedFile(ctx.String(flags.LogFolderFlag), l.CommandName())
	if err != nil {
		return
	}

	defaults = append(defaults, "--logfile", logFilePath)
	defaults = append(defaults, "--logfile-debug-level", "info")
	defaults = append(defaults, "--logfile-max-number", "1")
	defaults = append(defaults, "--suggested-fee-recipient", ctx.String(flags.TransactionFeeRecipientFlag))

	userFlags := l.ParseUserFlags(ctx)

	startFlags = mergeFlags(userFlags, defaults)

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

	initCommand := exec.Command(lighthouseDependencyName, args...)

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
