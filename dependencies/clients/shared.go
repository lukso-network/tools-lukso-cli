package clients

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
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

	jdkFolder = file.ClientsDir + "/jdk"

	VersionNotAvailable = "Not available"
)

var (
	AllClients = map[string]Client{
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

	url = strings.ReplaceAll(url, "|TAG|", client.tag())
	url = strings.ReplaceAll(url, "|OS|", client.os())
	url = strings.ReplaceAll(url, "|COMMIT|", client.commit())
	url = strings.ReplaceAll(url, "|ARCH|", client.arch())

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
	return client.tag()
}

// Since most clients don't need to init and it's more of a side effect, we Init nothing by default.
func (client *clientBinary) Init() error {
	return nil
}

func (client *clientBinary) tag() string {
	return ClientVersions[client.Name()]
}

func (client *clientBinary) commit() string {
	return ""
}

func (client *clientBinary) os() string {
	return ""
}

func (client *clientBinary) arch() string {
	return ""
}

func initClient(ctx *cli.Context, client Client) (err error) {
	log.Infof("⚙️  Running %s init...", client.Name())

	if client.IsRunning() {
		return errors.ErrAlreadyRunning
	}

	var (
		dataDir string
		cmdPath string
	)

	switch client.Name() {
	case gethDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(flags.GethDatadirFlag))
		cmdPath = client.FilePath()

	case erigonDependencyName:
		dataDir = fmt.Sprintf("--datadir=%s", ctx.String(flags.ErigonDatadirFlag))
		cmdPath = client.FilePath()
	}

	command := exec.Command(cmdPath, "init", dataDir, ctx.String(flags.GenesisJsonFlag))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
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

func (client *clientBinary) untarDir(dst, pattern string, t *tar.Reader) error {
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var (
			path       string
			headerName = header.Name
		)

		// for the sake of compatibility with updated versions remove the tag from the tarred file - teku/teku-xx.x.x => teku/teku, same with jdk
		if strings.Contains(header.Name, pattern) {
			headerName = replaceRootFolderName(header.Name, client.FileName())
		}

		path = filepath.Join(dst, headerName)

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

func (client *clientBinary) unzipDir(dst string, r *zip.Reader) (err error) {
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
func Setup() {
	// Execution
	Geth = NewGethClient()
	Erigon = NewErigonClient()
	Nethermind = NewNethermindClient()
	Besu = NewBesuClient()

	// Consensus
	Prysm = NewPrysmClient()
	Lighthouse = NewLighthouseClient()
	Teku = NewTekuClient()
	Nimbus2 = NewNimbus2Client()

	// Validators
	PrysmValidator = NewPrysmValidatorClient()
	LighthouseValidator = NewLighthouseValidatorClient()
	TekuValidator = NewTekuValidatorClient()
	Nimbus2Validator = NewNimbus2ValidatorClient()
}
