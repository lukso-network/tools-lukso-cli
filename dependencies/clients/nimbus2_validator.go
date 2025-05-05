package clients

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type Nimbus2ValidatorClient struct {
	*clientBinary
}

func NewNimbus2ValidatorClient() *Nimbus2ValidatorClient {
	return &Nimbus2ValidatorClient{
		&clientBinary{
			name:           nimbus2ValidatorDependencyName,
			fileName:       "validator_nimbus",
			baseUrl:        "",
			githubLocation: nimbus2GithubLocation,
			buildInfo:      nimbus2BuildInfo,
		},
	}
}

var Nimbus2Validator = NewNimbus2ValidatorClient()

var _ ValidatorBinaryDependency = &Nimbus2ValidatorClient{}

func (n *Nimbus2ValidatorClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = n.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.Nimbus2ValidatorConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))

	return
}

func (n *Nimbus2ValidatorClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if n.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", n.Name())

		err = n.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", n.Name())
	}

	command := exec.Command(fmt.Sprintf("./%s/build/nimbus_validator_client", nimbus2Folder), arguments...)

	err = prepareLogFile(ctx, command, n.CommandName())
	if err != nil {
		log.Errorf("There was an error while preparing a log file for %s: %v", n.Name(), err)
	}

	log.Infof("🔄  Starting %s", n.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, n.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", n.Name())

	return
}

func (n *Nimbus2ValidatorClient) Version() (version string) {
	return Nimbus2.Version()
}

func (n *Nimbus2ValidatorClient) Import(ctx *cli.Context) (err error) {
	args := []string{
		"deposits",
		"import",
		fmt.Sprintf("--data-dir=%s", ctx.String(flags.ValidatorWalletDirFlag)),
		ctx.String(flags.ValidatorKeysFlag),
	}

	validatorPass := ctx.String(flags.ValidatorPasswordFlag)
	if validatorPass != "" {
		log.Warn("Password flag not available for Nimbus2")
	}

	importCommand := exec.Command(fmt.Sprintf("./%s/build/nimbus_beacon_node", nimbus2Folder), args...)

	importCommand.Stdout = os.Stdout
	importCommand.Stderr = os.Stderr
	importCommand.Stdin = os.Stdin

	err = importCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while importing keystores: %v", err), 1)
	}

	return nil
}

func (n *Nimbus2ValidatorClient) List(ctx *cli.Context) (err error) {
	walletDir := ctx.String(flags.ValidatorWalletDirFlag)
	if walletDir == "" {
		return utils.Exit("❌  Wallet directory not provided - please provide a --validator-wallet-dir flag containing your keys directory", 1)
	}

	err = keystoreListWalk(walletDir)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while listing validators: %v", err), 1)
	}

	return
}

func (n *Nimbus2ValidatorClient) Exit(ctx *cli.Context) (err error) {
	args := []string{
		"deposits",
		"exit",
		fmt.Sprintf("--data-dir=%s", ctx.String(flags.ValidatorWalletDirFlag)),
		"--all",
	}

	rpc := ctx.String(flags.RpcAddressFlag)
	if rpc != "" {
		args = append(args, fmt.Sprintf("--rest-url=%s", rpc))
	}

	exitCommand := exec.Command(fmt.Sprintf("./%s/build/nimbus_beacon_node", nimbus2Folder), args...)

	exitCommand.Stdout = os.Stdout
	exitCommand.Stderr = os.Stderr
	exitCommand.Stdin = os.Stdin

	err = exitCommand.Run()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while exiting validators: %v", err), 1)
	}

	return nil
}
