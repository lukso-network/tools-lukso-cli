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

const besuFolder = common.ClientDepsFolder + "/besu"

type BesuClient struct {
	*clientBinary
}

func NewBesuClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *BesuClient {
	return &BesuClient{
		&clientBinary{
			name:           besuDependencyName,
			fileName:       "besu",
			baseUrl:        "https://github.com/hyperledger/besu/releases/download/|TAG|/besu-|TAG|.tar.gz",
			githubLocation: besuGithubLocation,
			buildInfo:      besuBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Besu dep.ExecutionClient
	_    dep.ExecutionClient = &BesuClient{}
)

func (b *BesuClient) Install(version string, isUpdate bool) (err error) {
	url := b.ParseUrl(version, b.Commit())

	return b.installer.InstallTar(url, file.ClientsDir, b.FileName(), "besu-")
}

func (b *BesuClient) Update() (err error) {
	tag := b.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", b.name)

	return b.Install(tag, true)
}

func (b *BesuClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = b.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.BesuConfigFileFlag)))

	return
}

func (b *BesuClient) FilePath() string {
	return besuFolder
}

func (b *BesuClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if b.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", b.Name())

		err = b.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", b.Name())
	}

	command := exec.Command(fmt.Sprintf("./%s/bin/besu", b.FilePath()), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.TimestampedFile(logFolder, b.FileName())
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

	log.Infof("🔄  Starting %s", b.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, b.FileName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", b.Name())

	return
}

func (b *BesuClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}

func (b *BesuClient) Version() (version string) {
	cmdVer := execVersionCmd(
		fmt.Sprintf("./%s/bin/besu", b.FilePath()),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Besu version output to parse:

	// besu/v24.7.0/linux-x86_64/oracle_openjdk-java-22
	return strings.Split(cmdVer, "/")[1]
}
