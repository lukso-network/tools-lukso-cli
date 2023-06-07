package clients

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
)

type PrysmValidatorClient struct {
	*clientBinary
}

func NewPrysmValidatorClient() *PrysmValidatorClient {
	return &PrysmValidatorClient{
		&clientBinary{
			name:           prysmValidatorDependencyName,
			commandName:    "prysm_validator",
			baseUrl:        "https://github.com/prysmaticlabs/prysm/releases/download/|TAG|/validator-|TAG|-|OS|-|ARCH|",
			githubLocation: prysmaticLabsGithubLocation,
		},
	}
}

var PrysmValidator = NewPrysmValidatorClient()

var _ ValidatorBinaryDependency = &PrysmValidatorClient{}

func (p *PrysmValidatorClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	validatorConfigExists := utils.FlagFileExists(ctx, flags.ValidatorConfigFileFlag)
	chainConfigExists := utils.FlagFileExists(ctx, flags.PrysmChainConfigFileFlag)
	validatorKeysExists := utils.FlagFileExists(ctx, flags.ValidatorKeysFlag)
	if !validatorConfigExists || !chainConfigExists {
		err = errors.ErrFlagPathInvalid

		return
	}
	if !validatorKeysExists {
		err = errors.ErrValidatorNotImported

		return
	}

	var passwordPipe *os.File
	validatorPasswordPath := ctx.String(flags.ValidatorWalletPasswordFileFlag)
	if validatorPasswordPath == "" {
		passwordPipe, err = utils.ReadValidatorPassword(ctx)
		if err != nil {
			err = utils.Exit(fmt.Sprintf("❌  There was an error while reading password: %v", err), 1)
		}
	}

	defer func() {
		if err != nil {
			os.Remove(passwordPipe.Name())
		}
	}()

	startFlags = p.ParseUserFlags(ctx)

	logFilePath, err := utils.PrepareTimestampedFile(ctx.String(flags.LogFolderFlag), p.CommandName())
	if err != nil {
		return
	}

	// terms of use already accepted during installation
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.ValidatorConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--log-file=%s", logFilePath))
	startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	if ctx.String(flags.ValidatorKeysFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--wallet-dir=%s", ctx.String(flags.ValidatorKeysFlag)))
	}
	if ctx.String(flags.ValidatorWalletPasswordFileFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--wallet-password-file=%s", ctx.String(flags.ValidatorWalletPasswordFileFlag)))
	}
	if ctx.String(flags.ValidatorChainConfigFileFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--chain-config-file=%s", ctx.String(flags.ValidatorChainConfigFileFlag)))
	}
	return
}

func (p *PrysmValidatorClient) Import(ctx *cli.Context) (err error) {
	args := []string{
		"accounts",
		"import",
		"--wallet-dir",
		ctx.String(flags.ValidatorWalletDirFlag),
		"--keys-dir",
		ctx.String(flags.ValidatorKeysFlag),
	}

	validatorPass := ctx.String(flags.ValidatorPasswordFlag)

	if validatorPass != "" {
		args = append(args, "--account-password-file", validatorPass)
	}

	initCommand := exec.Command("validator", args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err = initCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while importing validator accounts: %v", err), 1)
	}

	return nil
}

func (p *PrysmValidatorClient) List(ctx *cli.Context) (err error) {
	return
}
