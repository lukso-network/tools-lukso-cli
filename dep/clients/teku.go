package clients

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	tekuFolder = common.ClientDepsFolder + "/teku" // folder in which teku is stored (in tekuDepsFolder)

	tekuInstallURL = "https://artifacts.consensys.net/public/teku/raw/names/teku.tar.gz/versions/|TAG|/teku-|TAG|.tar.gz"
	jdkInstallURL  = "https://download.java.net/java/GA/jdk22.0.1/c7ec1332f7bb44aeba2eb341ae18aca4/8/GPL/openjdk-22.0.1_|OS|-|ARCH|_bin.tar.gz"
)

type TekuClient struct {
	*clientBinary
}

func NewTekuClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *TekuClient {
	return &TekuClient{
		&clientBinary{
			name:           tekuDependencyName,
			fileName:       "teku",
			baseUrl:        tekuInstallURL,
			githubLocation: tekuGithubLocation,
			buildInfo:      tekuBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Teku dep.ConsensusClient
	_    dep.ConsensusClient = &TekuClient{}
)

func (t *TekuClient) Install(version string, isUpdate bool) (err error) {
	url := t.ParseUrl(version, t.Commit())

	return t.installer.InstallTar(url, file.ClientsDir, t.FileName(), "teku-")
}

func (t *TekuClient) Update() (err error) {
	tag := t.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", t.name)

	return t.Install(tag, true)
}

func (t *TekuClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = t.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.TekuConfigFileFlag)))
	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--validators-proposer-default-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	}

	return
}

func (t *TekuClient) FilePath() string {
	return tekuFolder
}

func (t *TekuClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if t.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", t.Name())

		err = t.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", t.Name())
	}

	command := exec.Command(fmt.Sprintf("./%s/bin/teku", t.FilePath()), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.TimestampedFile(logFolder, t.FileName())
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

	log.Infof("🔄  Starting %s", t.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, t.FileName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", t.Name())

	return
}

func (t *TekuClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 5051)
}

func (t *TekuClient) Version() (version string) {
	cmdVer := execVersionCmd(
		t.FilePath(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Teku version output to parse:

	// 2025-01-27T11:22:24.202228671Z main INFO Starting configuration XmlConfiguration[location=jar:file:/home/m8b-stan/update-test/clients/teku/lib/teku-infrastructure-logging-24.12.1.jar!/log4j2.xml, lastModified=2025-01-27T10:50:28.949Z]...
	// 2025-01-27T11:22:24.203796770Z main INFO Start watching for changes to jar:file:/home/m8b-stan/update-test/clients/teku/lib/teku-infrastructure-logging-24.12.1.jar!/log4j2.xml every 0 seconds
	// 2025-01-27T11:22:24.204058832Z main INFO Configuration XmlConfiguration[location=jar:file:/home/m8b-stan/update-test/clients/teku/lib/teku-infrastructure-logging-24.12.1.jar!/log4j2.xml, lastModified=2025-01-27T10:50:28.949Z] started.
	// 2025-01-27T11:22:24.206286375Z main INFO Stopping configuration org.apache.logging.log4j.core.config.DefaultConfiguration@1b11171f...
	// 2025-01-27T11:22:24.206960660Z main INFO Configuration org.apache.logging.log4j.core.config.DefaultConfiguration@1b11171f stopped.
	// teku/v24.12.1/linux-x86_64/oracle_openjdk-java-22

	// Find the line with teku
	expr := regexp.MustCompile(fmt.Sprintf(`teku\/%s`, common.SemverExpressionRaw))
	s := expr.FindString(cmdVer)
	if s == "" {
		return VersionNotAvailable
	}

	return strings.Split(s, "/")[1]
}

func replaceRootFolderName(folder, targetRootName string) (path string) {
	splitHeader := strings.Split(folder, "/") // this assumes no / at the beginning of folder - not the case in tarred files we are interested in

	switch len(splitHeader) {
	case 0:
		return
	case 1:
		return targetRootName
	default:
		break
	}

	strippedHeader := splitHeader[1:]
	splitResolvedPath := append([]string{targetRootName}, strippedHeader...)

	return strings.Join(splitResolvedPath, "/")
}
