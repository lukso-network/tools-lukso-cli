package clients

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dep"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type LighthouseValidatorClient struct {
	*clientBinary
}

func NewLighthouseValidatorClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *LighthouseValidatorClient {
	return &LighthouseValidatorClient{
		&clientBinary{
			name:           lighthouseValidatorDependencyName,
			fileName:       lighthouseValidatorFileName, // we run it using lighthouse bin
			commandPath:    lighthouseValidatorCommandPath,
			baseUrl:        "", // no separate client for lighthouse validator - lighthouse_beacon for reference
			githubLocation: "",
			buildInfo:      lighthouseBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	LighthouseValidator dep.ValidatorClient
	_                   dep.ValidatorClient = &LighthouseValidatorClient{}
)

func (l *LighthouseValidatorClient) Install(version string, isUpdate bool) error {
	return nil
}

func (l *LighthouseValidatorClient) Update() error {
	return nil
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

	initCommand := exec.Command(Lighthouse.FilePath(), args...)

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
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	cmd := exec.Command(Lighthouse.FilePath(), "am", "validator", "list", "--datadir", walletDir)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not fetch imported keys within the validator wallet: %v ", err)
	}

	return nil
}

func (l *LighthouseValidatorClient) Exit(ctx *cli.Context) (err error) {
	keystore := ctx.String(flags.KeystoreFlag)
	if keystore == "" {
		return utils.Exit("❌  Keystore not provided - please provide a --keystore flag containing path to keystore", 1)
	}

	args := []string{"a", "validator", "exit", "--keystore", keystore, "--testnet-dir", ctx.String(flags.TestnetDirFlag)}

	rpc := ctx.String(flags.RpcAddressFlag)
	if rpc == "" {
		rpc = "http://localhost:4000" // because we use 4000 in configs we don't want users to use default 5052
	}

	args = append(args, "--beacon-node", rpc)

	exitCommand := exec.Command(Lighthouse.FilePath(), args...)

	exitCommand.Stdout = os.Stdout
	exitCommand.Stderr = os.Stderr
	exitCommand.Stdin = os.Stdin

	err = exitCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while exiting validator: %v", err), 1)
	}

	return
}

func (l *LighthouseValidatorClient) Version() (version string) {
	return Lighthouse.Version() // same binary
}
