package clients

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"time"
)

const (
	// client names - to use them in anything non-log related use CommandName() to access non-capitalized name
	gethDependencyName       = "Geth"
	prysmDependencyName      = "Prysm"
	validatorDependencyName  = "Validator"
	lighthouseDependencyName = "Lighthouse"
	erigonDependencyName     = "Erigon"
)

var (
	allClients = []ClientBinaryDependency{Geth}
)

type clientBinary struct {
	name           string
	commandName    string
	baseUrl        string
	filePath       string
	githubLocation string // user + repo, f.e. prysmaticlabs/prysm
}

func start(ctx *cli.Context, client ClientBinaryDependency, arguments []string) (err error) {
	if client.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", client.Name())

		err = client.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", client.Name())
	}

	command := exec.Command(client.CommandName(), arguments...)

	if client.Name() == gethDependencyName || client.Name() == erigonDependencyName {
		log.Infof("⚙️  Running %s init...", client.Name())

		err = initClient(ctx, client)
		if err != nil && err != errors.ErrAlreadyRunning { // if it is already running it will be caught during start
			log.Errorf("❌  There was an error while initalizing %s. Error: %v", client.Name(), err)

			return err
		}
	}

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, client.CommandName())
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

	log.Infof("🔄  Starting %s", client.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	return
}

func stop(client ClientBinaryDependency) (err error) {
	return
}

func status(client ClientBinaryDependency) {

}

func logs(client ClientBinaryDependency) {

}

func reset(client ClientBinaryDependency) {

}

func install(client ClientBinaryDependency) {

}

func update(client ClientBinaryDependency) {

}

func isRunning(client ClientBinaryDependency) bool {
	return false
}

func initClient(ctx *cli.Context, client ClientBinaryDependency) (err error) {
	if isRunning(client) {
		return errors.ErrAlreadyRunning
	}

	if !utils.FileExists(ctx.String(flags.GenesisJsonFlag)) {
		if ctx.Bool(flags.TestnetFlag) || ctx.Bool(flags.DevnetFlag) {
			return errors.ErrGenesisNotFound
		}

		message := `Choose your preferred initial LYX supply!
If you are a Genesis Validator, we recommend to choose the supply, which the majority of the Genesis Validators would choose,
to prevent your node from running on a network that can not finalize due to missing validators!
🗳️ See the voting results at https://deposit.mainnet.lukso.network

For more information read:
👉 https://medium.com/lukso/genesis-validators-deposit-smart-contract-freeze-and-testnet-launch-c5f7b568b1fc

Which initial LYX supply do you choose?
1: 35M LYX
2: 42M LYX
3: 100M LYX
> `
		var input string
		for input != "1" && input != "2" && input != "3" {
			input = utils.RegisterInputWithMessage(message)
			switch input {
			case "1":
				err = ctx.Set(flags.GenesisJsonFlag, configs.MainnetConfig+"/"+configs.Genesis35JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(flags.GenesisStateFlag, configs.MainnetConfig+"/"+configs.GenesisState35FilePath)

			case "2":
				err = ctx.Set(flags.GenesisJsonFlag, configs.MainnetConfig+"/"+configs.Genesis42JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(flags.GenesisStateFlag, configs.MainnetConfig+"/"+configs.GenesisState42FilePath)

			case "3":
				err = ctx.Set(flags.GenesisJsonFlag, configs.MainnetConfig+"/"+configs.Genesis100JsonPath)
				if err != nil {
					return
				}
				err = ctx.Set(flags.GenesisStateFlag, configs.MainnetConfig+"/"+configs.GenesisState100FilePath)

			default:
				log.Warn("Please select a valid option\n\n")
			}

			if err != nil {
				return
			}
		}
	}

	var dataDir string
	switch client.Name() {
	case gethDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(flags.GethDatadirFlag))
	case erigonDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(flags.ErigonDatadirFlag))
	}

	command := exec.Command(client.CommandName(), "init", dataDir, ctx.String(flags.GenesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
