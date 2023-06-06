package clients

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	gethDependencyName       = "Geth"
	prysmDependencyName      = "Prysm"
	validatorDependencyName  = "Validator"
	lighthouseDependencyName = "Lighthouse"
	erigonDependencyName     = "Erigon"
)

var (
	AllClients = []ClientBinaryDependency{Geth, Erigon, Prysm, Lighthouse, PrysmValidator, LighthouseValidator}
)

type clientBinary struct {
	name           string
	commandName    string
	baseUrl        string
	githubLocation string // user + repo, f.e. prysmaticlabs/prysm
}

var _ ClientBinaryDependency = &clientBinary{}

func (client *clientBinary) Start(ctx *cli.Context, arguments []string) (err error) {
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

func (client *clientBinary) Stop() (err error) {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())

	pidVal, err := pid.Load(pidLocation)
	if err != nil {
		log.Warnf("⏭️  %s is not running - skipping...", client.CommandName())

		return nil
	}

	err = pid.Kill(pidLocation, pidVal)
	if err != nil {
		return errors.ErrProcessNotFound
	}

	return
}

func (client *clientBinary) Logs(logFilePath string) (err error) {
	var commandName string
	var commandArgs []string

	commandName = "tail"
	commandArgs = []string{"-f", "-n", "+0"}

	command := exec.Command(commandName, append(commandArgs, logFilePath)...)

	command.Stdout = os.Stdout

	err = command.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Error("No error logs found")

		return
	}

	// error unrelated to command execution
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)
	}

	return
}

func (client *clientBinary) Reset(dataDirPath string) (err error) {
	if dataDirPath == "" {
		return utils.Exit(fmt.Sprintf("%v", errors.ErrFlagMissing), 1)
	}

	return os.RemoveAll(dataDirPath)
}

func (client *clientBinary) Install(tag, commitHash string) (err error) {
	fileUrl := client.ParseUrl(tag, commitHash)

	if utils.FileExists(client.FilePath()) {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", client.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	response, err := http.Get(fileUrl)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", fileUrl)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			fileUrl,
			response.Status,
		)
	}

	var responseReader io.Reader = response.Body

	// this means that we are fetching tared client
	switch client.Name() {
	case gethDependencyName, erigonDependencyName, lighthouseDependencyName:
		g, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}

		defer func() {
			_ = g.Close()
		}()

		t := tar.NewReader(g)
		for {
			header, err := t.Next()

			switch {
			case err == io.EOF:
				break

			case err != nil:
				return err

			default:

			}

			var targetHeader string
			switch client.Name() {
			case gethDependencyName:
				targetHeader = "/" + gethDependencyName
			case erigonDependencyName, lighthouseDependencyName:
				targetHeader = client.Name()
			}

			if header.Typeflag == tar.TypeReg && strings.Contains(header.Name, targetHeader) {
				responseReader = t

				break
			}
		}
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(responseReader)
	if err != nil {
		return
	}

	err = os.WriteFile(client.FilePath(), buf.Bytes(), configs.BinaryPerms)

	if err != nil && strings.Contains(err.Error(), "Permission denied") {
		return errors.ErrNeedRoot
	}

	if err != nil {
		log.Infof("❌  Couldn't save file: %v", err)

		return
	}

	log.Infof("✅  %s downloaded!\n\n", client.Name())

	return

}

func (client *clientBinary) Update() {

}

func (client *clientBinary) IsRunning() bool {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())

	return pid.Exists(pidLocation)
}

func initClient(ctx *cli.Context, client ClientBinaryDependency) (err error) {
	if client.IsRunning() {
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

func (client *clientBinary) ParseUrl(tag, commitHash string) (url string) {
	url = client.baseUrl

	url = strings.Replace(url, "|TAG|", tag, -1)
	url = strings.Replace(url, "|OS|", system.Os, -1)
	url = strings.Replace(url, "|COMMIT|", commitHash, -1)
	url = strings.Replace(url, "|ARCH|", system.Arch, -1)

	return
}

func (client *clientBinary) FilePath() string {
	return system.UnixBinDir + "/" + client.CommandName()
}

func (client *clientBinary) Name() string {
	return client.name
}

func (client *clientBinary) CommandName() string {
	return client.commandName
}

func (client *clientBinary) ParseUserFlags(ctx *cli.Context) (startFlags []string) {
	name := client.name
	args := ctx.Args()
	argsLen := args.Len()
	flagsToSkip := []string{
		flags.ValidatorFlag,
		flags.GethConfigFileFlag,
		flags.PrysmConfigFileFlag,
		flags.ValidatorConfigFileFlag,
		flags.ValidatorWalletPasswordFileFlag,
	}

	for i := 0; i < argsLen; i++ {
		skip := false
		arg := args.Get(i)
		for _, flagToSkip := range flagsToSkip {
			if arg == fmt.Sprintf("--%s", flagToSkip) {
				skip = true
			}
		}
		if skip {
			continue
		}

		if strings.HasPrefix(arg, fmt.Sprintf("--%s", name)) {
			if i+1 == argsLen {
				startFlags = append(startFlags, removePrefix(arg, name))

				return
			}

			// we found a flag for our client - now we need to check if it's a value or bool flag
			nextArg := args.Get(i + 1)
			if strings.HasPrefix(nextArg, "--") { // we found a next flag, so current one is a bool
				startFlags = append(startFlags, removePrefix(arg, name))

				continue
			}

			startFlags = append(startFlags, removePrefix(arg, name), nextArg)
		}
	}

	return
}

func (client *clientBinary) PrepareStartFlags() (startFlags []string) {
	_ = cli.Exit(fmt.Sprintf("FATAL: START FLAGS NOT CONFIGURED FOR %s CLIENT - PLEASE MARK THIS ISSUE TO THE LUKSO TEAM", client.Name()), 1)

	return
}

func removePrefix(arg, name string) string {
	prefix := fmt.Sprintf("--%s-", name)

	arg = strings.TrimPrefix(arg, prefix)

	return fmt.Sprintf("--%s", strings.Trim(arg, "- "))
}
