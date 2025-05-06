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
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type PrysmValidatorClient struct {
	*clientBinary
}

func NewPrysmValidatorClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *PrysmValidatorClient {
	return &PrysmValidatorClient{
		&clientBinary{
			name:           prysmValidatorDependencyName,
			fileName:       "validator_pr",
			baseUrl:        "https://github.com/prysmaticlabs/prysm/releases/download/|TAG|/validator-|TAG|-|OS|-|ARCH|",
			githubLocation: prysmaticLabsGithubLocation,
			buildInfo:      prysmBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var PrysmValidator ValidatorBinaryDependency

var _ ValidatorBinaryDependency = &PrysmValidatorClient{}

func (p *PrysmValidatorClient) Install(version string, isUpdate bool) error {
	url := p.ParseUrl(version, p.commit())

	return p.installer.InstallFile(url, p.FilePath())
}

func (p *PrysmValidatorClient) Update() (err error) {
	tag := p.tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", p.name)

	return p.Install(tag, true)
}

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

func (p *PrysmValidatorClient) Peers(ctx *cli.Context) (outbound int, inbound int, err error) {
	return
}

func (p *PrysmValidatorClient) Version() (version string) {
	cmdVer := execVersionCmd(
		p.FilePath(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Prysm validator version output to parse:

	// validator version Prysm/v5.0.4/3b184f43c86baf6c36478f65a5113e7cf0836d41. Built at: 2024-06-21 00:26:00+00:00

	// -> ...|Prysm/v5.0.4/3b184f43c86baf6c36478f65a5113e7cf0836d41.|...
	s := strings.Split(cmdVer, " ")[2]
	// -> ...|v5.0.4|...
	return strings.Split(s, "/")[1]
}

func (p *PrysmValidatorClient) List(ctx *cli.Context) (err error) {
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	cmd := exec.Command("validator", "accounts", "list", "--wallet-dir", walletDir)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not fetch imported keys within the validator wallet: %v ", err)
	}

	return nil
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
		name:     "prysmctl",
		fileName: "prysmctl",
		baseUrl:  "https://github.com/prysmaticlabs/prysm/releases/download/|TAG|/prysmctl-|TAG|-|OS|-|ARCH|",
	}

	versionCommand := exec.Command(PrysmValidator.FileName(), "--version")
	buf := new(bytes.Buffer)

	versionCommand.Stdout = buf

	err = versionCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while getting prysm version: %v", err), 1)
	}

	versionOutput := buf.String()
	version := strings.Split(versionOutput, "/")[1]

	url := prysmctlBin.ParseUrl(version, "")

	log.Info("⬇️  Downloading prysmctl...")
	err = prysmctlBin.Install(url, false)

	return
}
