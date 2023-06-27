package clients

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type PrysmValidatorClient struct {
	*clientBinary
}

func NewPrysmValidatorClient() *PrysmValidatorClient {
	return &PrysmValidatorClient{
		&clientBinary{
			name:           prysmValidatorDependencyName,
			commandName:    "validator",
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

	// terms of use already accepted during installation
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.ValidatorConfigFileFlag)))
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

func (p *PrysmValidatorClient) Exit(ctx *cli.Context) (err error) {
	wallet := ctx.String(flags.ValidatorWalletDirFlag)
	if wallet == "" {
		return utils.Exit("❌  Wallet directory not provided - please provide a --validator-wallet-dir flag containing your keys directory", 1)
	}

	if !utils.FileExists(wallet) {
		return utils.Exit("❌  Wallet directory missing - please provide a --validator-wallet-dir flag containing your keys directory or use a network flag", 1)
	}

	_, err = os.Stat(fmt.Sprintf("%s/prysmctl", system.UnixBinDir))
	if err != nil {
		err = installPrysmctl()
	}
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while installing prysmctl: %v", err), 1)
	}

	args := []string{"validator", "exit", "--wallet-dir", wallet}
	rpc := ctx.String(flags.RpcAddressFlag)
	if rpc != "" {
		args = append(args, "--beacon-rpc-provider", rpc)
	}

	exitCommand := exec.Command("prysmctl", args...)

	exitCommand.Stdout = os.Stdout
	exitCommand.Stderr = os.Stderr
	exitCommand.Stdin = os.Stdin

	err = exitCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while exiting validator: %v", err), 1)
	}

	return
}

func installPrysmctl() (err error) {
	message := "In order to exit the validator from the network you will need to install additional tool - prysmctl\n" +
		"Do you want to proceed with instalation? [Y/n]\n>"

	input := utils.RegisterInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Info("Aborting...")

		os.Exit(0)
	}

	prysmctlBin := clientBinary{
		name:        "prysmctl",
		commandName: "prysmctl",
		baseUrl:     "https://github.com/prysmaticlabs/prysm/releases/download/|TAG|/prysmctl-|TAG|-|OS|-|ARCH|",
	}

	versionCommand := exec.Command(PrysmValidator.CommandName(), "--version")
	buf := new(bytes.Buffer)

	versionCommand.Stdout = buf

	err = versionCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while getting prysm version: %v", err), 1)
	}

	versionOutput := string(buf.Bytes())
	version := strings.Split(versionOutput, "/")[1]

	url := prysmctlBin.ParseUrl(version, "")

	log.Info("⬇️  Downloading prysmctl...")
	err = prysmctlBin.Install(url, false)

	return
}
