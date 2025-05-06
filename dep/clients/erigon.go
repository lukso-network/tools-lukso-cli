package clients

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const (
	erigonFolder = common.ClientDepsFolder + "/erigon"
)

type ErigonClient struct {
	*clientBinary
}

func NewErigonClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *ErigonClient {
	return &ErigonClient{
		&clientBinary{
			name:           erigonDependencyName,
			fileName:       "erigon",
			baseUrl:        "https://github.com/erigontech/erigon/releases/download/|TAG|/erigon_|TAG|_|OS|_|ARCH|.tar.gz",
			githubLocation: erigonGithubLocation,
			buildInfo:      erigonBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Erigon dep.ExecutionClient
	_      dep.ExecutionClient = &ErigonClient{}
)

func (e *ErigonClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if e.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", e.Name())

		err = e.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", e.Name())
	}

	err = e.Init()
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

	fullPath, err = utils.TimestampedFile(logFolder, e.FileName())
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

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, e.FileName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", e.Name())

	return
}

func (e *ErigonClient) Install(version string, isUpdate bool) (err error) {
	url := e.ParseUrl(version, e.Commit())

	return e.installer.InstallTar(url, file.ClientsDir, e.FileName(), "erigon-")
}

func (e *ErigonClient) Update() (err error) {
	tag := e.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", e.name)

	return e.Install(tag, true)
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
		e.FilePath(),
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
