package clients

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const (
	erigonFolder = common.ClientDepsFolder + "/erigon"
)

type ErigonClient struct {
	*clientBinary
}

func NewErigonClient() *ErigonClient {
	return &ErigonClient{
		&clientBinary{
			name:           erigonDependencyName,
			commandName:    "erigon",
			baseUrl:        "https://github.com/erigontech/erigon/releases/download/|TAG|/erigon_|TAG|_|OS|_|ARCH|.tar.gz",
			githubLocation: erigonGithubLocation,
		},
	}
}

var Erigon = NewErigonClient()

var _ ClientBinaryDependency = &ErigonClient{}

func (e *ErigonClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if e.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", e.Name())

		err = e.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", e.Name())
	}

	err = initClient(ctx, e)
	if err != nil && err != errors.ErrAlreadyRunning { // if it is already running it will be caught during start
		log.Errorf("❌  There was an error while initalizing %s. Error: %v", e.Name(), err)

		return err
	}

	command := exec.Command(fmt.Sprintf("./%s/erigon", e.FilePath()), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, e.CommandName())
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

	log.Infof("🔄  Starting %s", e.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, e.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", e.Name())

	return
}

func (e *ErigonClient) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(e.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", e.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	err = installAndExtractFromURL(url, e.name, common.ClientDepsFolder, tarFormat, isUpdate)
	if err != nil {
		return
	}

	permFunc := func(path string, d fs.DirEntry, err error) error {
		return os.Chmod(path, fs.ModePerm)
	}

	err = filepath.WalkDir(e.FilePath(), permFunc)
	if err != nil {
		return
	}

	return
}

func (e *ErigonClient) Update() (err error) {
	tag := e.getVersion()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", e.name)

	url := e.ParseUrl(tag, "")

	return e.Install(url, true)
}

func (e *ErigonClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.ErigonConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	startFlags = e.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.ErigonConfigFileFlag)))

	return
}

func (e *ErigonClient) FilePath() string {
	return erigonFolder
}

func (e *ErigonClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}

func (e *ErigonClient) Version() (version string) {
	cmdVer := execVersionCmd(
		fmt.Sprintf("./%s/%s", e.FilePath(), e.CommandName()),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Erigon version output to parse:

	// erigon version 2.60.10-3afee08b

	// -> ...|2.60.10-3afee08b
	s := strings.Split(cmdVer, " ")[2]
	// -> 2.60.10|...
	return fmt.Sprintf("v%s", strings.Split(s, "-")[0])
}
