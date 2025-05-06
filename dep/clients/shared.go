package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep"
	"github.com/lukso-network/tools-lukso-cli/dep/types/apitypes"
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
	erigonGithubLocation        = "erigontech/erigon"
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

	VersionNotAvailable = "Not available"
)

var (
	ExecutionClients = map[string]dep.ExecutionClient{
		gethDependencyName:       Geth,
		erigonDependencyName:     Erigon,
		nethermindDependencyName: Nethermind,
		besuDependencyName:       Besu,
	}

	ConsensusClients = map[string]dep.ConsensusClient{
		prysmDependencyName:      Prysm,
		lighthouseDependencyName: Lighthouse,
		tekuDependencyName:       Teku,
		nimbus2DependencyName:    Nimbus2,
	}

	ValidatorClients = map[string]dep.ValidatorClient{
		prysmValidatorDependencyName:      PrysmValidator,
		lighthouseValidatorDependencyName: LighthouseValidator,
		tekuValidatorDependencyName:       TekuValidator,
		nimbus2ValidatorDependencyName:    Nimbus2Validator,
	}

	AllClients = map[string]dep.Client{
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

	// for ordered actions
	AllClientNames = []string{
		gethDependencyName,
		erigonDependencyName,
		nethermindDependencyName,
		besuDependencyName,
		prysmDependencyName,
		prysmValidatorDependencyName,
		lighthouseDependencyName,
		lighthouseValidatorDependencyName,
		tekuDependencyName,
		tekuValidatorDependencyName,
		nimbus2DependencyName,
		nimbus2ValidatorDependencyName,
	}
)

// clientBinary represents a shared logic between all clients that they share.
// It is not, however, a client. What cannot be implemented generically (e.g. install), should be implemented explicitly by the given client.
// The implementation should rely on build - if app is able to build, it means every method was implemented correctly.
type clientBinary struct {
	name           string
	fileName       string
	baseUrl        string
	githubLocation string // user + repo, f.e. prysmaticlabs/prysm
	buildInfo      buildInfo

	log       logger.Logger
	file      file.Manager
	installer installer.Installer
	pid       pid.Pid
}

func (client *clientBinary) Start(ctx *cli.Context, arguments []string) (err error) {
	if client.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", client.Name())

		err = client.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", client.Name())
	}

	command := exec.Command(client.FilePath(), arguments...)

	err = client.logFile(ctx.String(flags.LogFolderFlag), command)
	if err != nil {
		return
	}

	log.Infof("🔄  Starting %s", client.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.FileName())
	err = client.pid.Create(pidLocation, command.Process.Pid)

	log.Infof("✅  %s started!", client.Name())

	return
}

func (client *clientBinary) Stop() (err error) {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.FileName())

	pidVal, err := client.pid.Load(pidLocation)
	if err != nil {
		log.Warnf("⏭️  %s is not running - skipping...", client.Name())

		return nil
	}

	err = client.pid.Kill(pidLocation, pidVal)
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

func (client *clientBinary) Reset(dataDir string) (err error) {
	if dataDir == "" {
		return utils.Exit(fmt.Sprintf("%v", errors.ErrFlagMissing), 1)
	}

	return client.file.RemoveAll(dataDir)
}

func (client *clientBinary) IsRunning() bool {
	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, client.FileName())

	return pid.Exists(pidLocation)
}

func (client *clientBinary) ParseUrl(tag, commitHash string) (url string) {
	url = client.baseUrl

	url = strings.ReplaceAll(url, "|TAG|", tag)
	url = strings.ReplaceAll(url, "|COMMIT|", commitHash)
	url = strings.ReplaceAll(url, "|OS|", client.Os())
	url = strings.ReplaceAll(url, "|ARCH|", client.Arch())

	return
}

func (client *clientBinary) ParseUserFlags(ctx *cli.Context) (startFlags []string) {
	name := client.FileName()
	args := ctx.Args()
	argsLen := args.Len()
	flagsToSkip := []string{
		flags.ValidatorFlag,
		flags.GethConfigFileFlag,
		flags.ErigonConfigFileFlag,
		flags.NethermindConfigFileFlag,
		flags.BesuConfigFileFlag,
		flags.PrysmConfigFileFlag,
		flags.LighthouseConfigFileFlag,
		flags.TekuConfigFileFlag,
		flags.Nimbus2ConfigFileFlag,
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

func (client *clientBinary) Name() string {
	return client.name
}

func (client *clientBinary) FileName() string {
	return client.fileName
}

func (client *clientBinary) FileDir() string {
	return file.ClientsDir + "/" + client.FileName()
}

func (client *clientBinary) FilePath() string {
	return client.FileDir() + client.FileName()
}

func (client *clientBinary) Version() (v string) {
	return client.Tag()
}

// Since most clients don't need to init and it's more of a side effect, we Init nothing by default.
func (client *clientBinary) Init() error {
	return nil
}

func (client *clientBinary) Tag() string {
	return ClientVersions[client.Name()]
}

func (client *clientBinary) Commit() string {
	return ""
}

func (client *clientBinary) Os() string {
	return ""
}

func (client *clientBinary) Arch() string {
	return ""
}

func execVersionCmd(cmd string) (ver string) {
	buf := new(bytes.Buffer)

	args := []string{
		"--version",
	}
	versionCommand := exec.Command(cmd, args...)
	versionCommand.Stdout = buf
	versionCommand.Stderr = buf

	err := versionCommand.Run()
	if err != nil {
		ver = VersionNotAvailable

		return
	}

	return buf.String()
}

func removePrefix(arg, name string) string {
	prefix := fmt.Sprintf("--%s-", name)

	arg = strings.TrimPrefix(arg, prefix)

	return fmt.Sprintf("--%s", strings.Trim(arg, "- "))
}

func RunningClients() (clients []string) {
	for _, client := range AllClients {
		if client.IsRunning() {
			clients = append(clients, client.Name())
		}
	}

	return
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

func isJdkInstalled() bool {
	// JDK installed outside of the CLI
	_, isInstalled := os.LookupEnv(system.JavaHomeEnv)
	if isInstalled {
		return true
	}

	_, err := os.Stat(jdkFolder)

	return err == nil
}

func (client *clientBinary) logFile(dir string, command *exec.Cmd) (err error) {
	var (
		logFile  *os.File
		fullPath string
	)

	if dir == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.TimestampedFile(dir, client.FileName())
	if err != nil {
		return
	}

	err = client.file.Create(fullPath)
	if err != nil {
		return
	}

	logFile, err = client.file.Open(fullPath)
	if err != nil {
		return
	}

	command.Stdout = logFile
	command.Stderr = logFile

	return
}

// keystoreListWalk walks through the given directory (and its subdirectories) and prints all found pubkeys.
func keystoreListWalk(walletDir string) (err error) {
	validatorIndex := 0

	walkFunc := func(path string, d fs.DirEntry, entryError error) (err error) {
		if d == nil {
			return nil
		}

		keystoreExt := filepath.Ext(d.Name())
		if !strings.Contains(keystoreExt, "json") {
			return nil
		}

		keystoreFile, err := os.Open(path)
		if err != nil {
			return
		}
		defer keystoreFile.Close()

		keystoreFileBytes, err := io.ReadAll(keystoreFile)
		if err != nil {
			return
		}

		keystore := struct {
			Pubkey string `json:"pubkey"`
		}{}

		err = json.Unmarshal(keystoreFileBytes, &keystore)
		if err != nil {
			return
		}

		log.Infof("Validator #%d: %s", validatorIndex, keystore.Pubkey)
		validatorIndex++

		return
	}

	err = filepath.WalkDir(walletDir, walkFunc)
	if err != nil {
		return
	}

	if validatorIndex == 0 {
		log.Info("No validator keys listed. To import your validator keys run lukso validator import")
	}

	return
}

// Setup populates clients with external dependencies: file management, logger etc.
func Setup(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) {
	// Execution
	Geth = NewGethClient(log, file, installer, pid)
	Erigon = NewErigonClient(log, file, installer, pid)
	Nethermind = NewNethermindClient(log, file, installer, pid)
	Besu = NewBesuClient(log, file, installer, pid)

	// Consensus
	Prysm = NewPrysmClient(log, file, installer, pid)
	Lighthouse = NewLighthouseClient(log, file, installer, pid)
	Teku = NewTekuClient(log, file, installer, pid)
	Nimbus2 = NewNimbus2Client(log, file, installer, pid)

	// Validators
	PrysmValidator = NewPrysmValidatorClient(log, file, installer, pid)
	LighthouseValidator = NewLighthouseValidatorClient(log, file, installer, pid)
	TekuValidator = NewTekuValidatorClient(log, file, installer, pid)
	Nimbus2Validator = NewNimbus2ValidatorClient(log, file, installer, pid)
}
