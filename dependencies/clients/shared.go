package clients

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/dependencies/types/apitypes"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const (
	gethDependencyName                = "Geth"
	prysmDependencyName               = "Prysm"
	prysmValidatorDependencyName      = "Prysm Validator"
	lighthouseDependencyName          = "Lighthouse"
	lighthouseValidatorDependencyName = "Lighthouse Validator"
	erigonDependencyName              = "Erigon"
	tekuDependencyName                = "Teku"
	tekuValidatorDependencyName       = "Teku Validator"
	nethermindDependencyName          = "Nethermind"
	besuDependencyName                = "Besu"
	nimbus2DependencyName             = "Nimbus2"
	nimbus2ValidatorDependencyName    = "Nimbus2 Validator"

	gethGithubLocation          = "ethereum/go-ethereum"
	prysmaticLabsGithubLocation = "prysmaticlabs/prysm"
	lighthouseGithubLocation    = "sigp/lighthouse"
	erigonGithubLocation        = "ledgerwatch/erigon"
	tekuGithubLocation          = "Consensys/teku"
	nethermindGithubLocation    = "NethermindEth/nethermind"
	besuGithubLocation          = "hyperledger/besu"
	nimbus2GithubLocation       = "status-im/nimbus-eth2"

	peerDirectionInbound  = "inbound"
	peerDirectionOutbound = "outbound"
	peerStateConnected    = "connected"

	// distinct between zip and tar archives
	zipFormat = "zip"
	tarFormat = "tar"

	jdkFolder = common.ClientDepsFolder + "/jdk"
)

var (
	AllClients = map[string]ClientBinaryDependency{
		gethDependencyName:                Geth,
		erigonDependencyName:              Erigon,
		prysmDependencyName:               Prysm,
		lighthouseDependencyName:          Lighthouse,
		prysmValidatorDependencyName:      PrysmValidator,
		lighthouseValidatorDependencyName: LighthouseValidator,
		tekuDependencyName:                Teku,
		tekuValidatorDependencyName:       TekuValidator,
		nethermindDependencyName:          Nethermind,
		besuDependencyName:                Besu,
		nimbus2DependencyName:             Nimbus2,
		nimbus2ValidatorDependencyName:    Nimbus2Validator,
	}

	ClientVersions = map[string]string{
		gethDependencyName:                common.GethTag,
		erigonDependencyName:              common.ErigonTag,
		nethermindDependencyName:          common.NethermindTag,
		prysmDependencyName:               common.PrysmTag,
		lighthouseDependencyName:          common.LighthouseTag,
		prysmValidatorDependencyName:      common.PrysmTag,
		lighthouseValidatorDependencyName: common.LighthouseTag,
		tekuDependencyName:                common.TekuTag,
		besuDependencyName:                common.BesuTag,
		nimbus2DependencyName:             common.Nimbus2Tag,
	}
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

	log.Infof("🔄  Starting %s", client.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", client.Name())

	return
}

func (client *clientBinary) Stop() (err error) {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())

	pidVal, err := pid.Load(pidLocation)
	if err != nil {
		log.Warnf("⏭️  %s is not running - skipping...", client.Name())

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

func (client *clientBinary) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(client.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", client.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	response, err := http.Get(url)

	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", url)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading on file url: %s. Response code: %s",
			url,
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
				targetHeader = "/" + client.CommandName()
			case erigonDependencyName:
				targetHeader = "erigon"
			case lighthouseDependencyName:
				targetHeader = "lighthouse"
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

	switch isUpdate {
	case true:
		log.Infof("✅  %s updated!\n\n", client.Name())
	case false:
		log.Infof("✅  %s downloaded!\n\n", client.Name())
	}

	return
}

func (client *clientBinary) Update() (err error) {
	tag := client.getVersion()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", client.name)

	url := client.ParseUrl(tag, "")

	return client.Install(url, true)
}

func (client *clientBinary) IsRunning() bool {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.CommandName())

	return pid.Exists(pidLocation)
}

func initClient(ctx *cli.Context, client ClientBinaryDependency) (err error) {
	log.Infof("⚙️  Running %s init...", client.Name())

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
	name := client.commandName
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

func (client *clientBinary) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	_ = utils.Exit(fmt.Sprintf("FATAL: START FLAGS NOT CONFIGURED FOR %s CLIENT - PLEASE MARK THIS ISSUE TO THE LUKSO TEAM", client.Name()), 1)

	return
}

func (client *clientBinary) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	_ = utils.Exit(fmt.Sprintf("FATAL: STATUS PEERS NOT CONFIGURED FOR %s CLIENT - PLEASE MARK THIS ISSUE TO THE LUKSO TEAM", client.Name()), 1)

	return
}

func (client *clientBinary) getVersion() string {
	return ClientVersions[client.Name()]
}

func removePrefix(arg, name string) string {
	prefix := fmt.Sprintf("--%s-", name)

	arg = strings.TrimPrefix(arg, prefix)

	return fmt.Sprintf("--%s", strings.Trim(arg, "- "))
}

func IsAnyRunning() bool {
	var runningClients string
	for _, client := range AllClients {
		if client.IsRunning() {
			runningClients += fmt.Sprintf("- %s\n", client.Name())
		}
	}

	if runningClients == "" {
		return false
	}

	log.Warnf("⚠️  Please stop the following clients before continuing: \n%s", runningClients)

	return true
}

func defaultExecutionPeers(ctx *cli.Context, defaultPort int) (outbound, inbound int, err error) {
	host := ctx.String(flags.ExecutionClientHost)
	port := ctx.Int(flags.ExecutionClientPort)
	if port == 0 {
		port = defaultPort
	}

	url := fmt.Sprintf("http://%s:%d", host, port)

	reqBodyBytes, err := json.Marshal(apitypes.JsonRpcRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  "admin_peers",
		Params:  []string{},
	})
	if err != nil {
		return
	}

	reqBody := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodGet, url, reqBody)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	peersResp := &apitypes.PeersJsonRpcResponse{}
	err = json.Unmarshal(respBodyBytes, peersResp)
	if err != nil {
		return
	}
	if peersResp.Error != nil {
		err = errors.ErrRpcError

		return
	}

	for _, peer := range peersResp.Result {
		switch peer.Network.Inbound {
		case true:
			inbound++
		case false:
			outbound++
		}
	}

	return
}

func defaultConsensusPeers(ctx *cli.Context, defaultPort int) (outbound, inbound int, err error) {
	host := ctx.String(flags.ConsensusClientHost)
	port := ctx.Int(flags.ConsensusClientPort)
	if port == 0 {
		port = defaultPort
	}

	url := fmt.Sprintf("http://%s:%d/eth/v1/node/peers", host, port)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	peersResp := &apitypes.PeersBeaconAPIResponse{}
	err = json.Unmarshal(respBodyBytes, peersResp)
	if err != nil {
		return
	}

	for _, peer := range peersResp.Data {
		if peer.State != peerStateConnected {
			continue
		}

		switch peer.Direction {
		case peerDirectionInbound:
			inbound++
		case peerDirectionOutbound:
			outbound++
		}
	}

	return
}

func untarDir(dst string, t *tar.Reader) error {
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dst, header.Name)
		fmt.Println(header.Name)

		// for the sake of compatibility with updated versions remove the tag from the tarred file - teku/teku-xx.x.x => teku/teku, same with jdk
		switch {
		case strings.Contains(header.Name, "teku-"):
			newHeader := replaceRootFolderName(header.Name, "teku")
			path = filepath.Join(dst, newHeader)

		case strings.Contains(header.Name, "jdk-"):
			newHeader := replaceRootFolderName(header.Name, "jdk")
			path = filepath.Join(dst, newHeader)

		case strings.Contains(header.Name, "besu-"):
			newHeader := replaceRootFolderName(header.Name, "besu")
			path = filepath.Join(dst, newHeader)
		case strings.Contains(header.Name, "nimbus-eth2"):
			fmt.Println("Here!")
			newHeader := replaceRootFolderName(header.Name, "nimbus2")
			path = filepath.Join(dst, newHeader)
		}

		info := header.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, info.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func installAndExtractFromURL(url, name, dst, format string, isUpdate bool) (err error) {
	response, err := http.Get(url)
	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", url)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading file at URL: %s. Response code: %s",
			url,
			response.Status,
		)
	}

	switch format {
	case tarFormat:
		var g *gzip.Reader
		g, err = gzip.NewReader(response.Body)
		if err != nil {
			return
		}

		defer func() {
			_ = g.Close()
		}()

		tarReader := tar.NewReader(g)

		err = untarDir(dst, tarReader)
		if err != nil {
			return
		}

	case zipFormat:
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		buf := bytes.NewReader(b)
		var r *zip.Reader

		r, err = zip.NewReader(buf, int64(buf.Len()))
		if err != nil {
			return err
		}

		err = unzipDir(dst, r)
		if err != nil {
			return err
		}
	}

	switch isUpdate {
	case true:
		log.Infof("✅  %s updated!\n\n", name)
	case false:
		log.Infof("✅  %s downloaded!\n\n", name)
	}

	return
}

func unzipDir(dst string, r *zip.Reader) (err error) {
	for _, header := range r.File {
		path := filepath.Join(dst, header.Name)

		info := header.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, info.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		defer file.Close()

		readFile, err := header.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, readFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func getUnameArch() (arch string) {
	fallback := func() {
		log.Info("⚠️  Unknown OS detected: proceeding with x86_64 as a default arch")
		arch = "x86_64"
	}

	switch system.Os {
	case system.Ubuntu, system.Macos:
		buf := new(bytes.Buffer)

		uname := exec.Command("uname", "-m")
		uname.Stdout = buf

		err := uname.Run()
		if err != nil {
			fallback()

			break
		}

		arch = strings.Trim(buf.String(), "\n\t ")

	default:
		fallback()
	}

	if arch != "x86_64" && arch != "aarch64" {
		fallback()
	}

	return
}

func isJdkInstalled() bool {
	// JDK installed outside of the CLI
	_, isInstalled := os.LookupEnv(system.JavaHomeEnv)
	if isInstalled {
		return true
	}

	_, err := os.Stat(jdkFolder)

	return err == nil
}

func prepareLogFile(ctx *cli.Context, command *exec.Cmd, client string) (err error) {
	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, client)
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

	return
}
