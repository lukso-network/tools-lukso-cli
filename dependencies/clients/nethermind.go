package clients

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const (
	nethermindFolder = common.ClientDepsFolder + "/nethermind"
)

type NethermindClient struct {
	*clientBinary
}

func NewNethermindClient() *NethermindClient {
	return &NethermindClient{
		&clientBinary{
			name:           nethermindDependencyName,
			commandName:    "nethermind",
			baseUrl:        "https://github.com/NethermindEth/nethermind/releases/download/|TAG|/nethermind-|TAG|-|COMMIT|-|OS|-|ARCH|.zip",
			githubLocation: nethermindGithubLocation,
		},
	}
}

var Nethermind = NewNethermindClient()

var _ ClientBinaryDependency = &NethermindClient{}

func (n *NethermindClient) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(n.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", n.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	err = installAndExtractFromURL(url, n.name, nethermindFolder, zipFormat, isUpdate)
	if err != nil {
		return
	}

	permFunc := func(path string, d fs.DirEntry, err error) error {
		return os.Chmod(path, fs.ModePerm)
	}

	err = filepath.WalkDir(n.FilePath(), permFunc)
	if err != nil {
		return
	}

	return
}

func (n *NethermindClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if n.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", n.Name())

		err = n.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", n.Name())
	}

	command := exec.Command(fmt.Sprintf("./%s/nethermind", n.FilePath()), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, n.CommandName())
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

func (n *NethermindClient) ParseUrl(tag, commitHash string) (url string) {
	url = n.baseUrl
	osName := system.Os
	archName := system.GetArch()

	if osName == system.Macos {
		osName = "macos"
	}

	if archName == "x86_64" {
		archName = "x64"
	}
	if archName == "arm" || archName == "aarch64" {
		archName = "arm64"
	}

	url = strings.Replace(url, "|TAG|", tag, -1)
	url = strings.Replace(url, "|OS|", osName, -1)
	url = strings.Replace(url, "|COMMIT|", commitHash, -1)
	url = strings.Replace(url, "|ARCH|", archName, -1)

	return
}

func (n *NethermindClient) Update() (err error) {
	tag := n.getVersion()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", n.name)

	// this commit hash is hardcoded, but since update should only be responsible for updating the client to
	// LUKSO supported version this is fine.
	url := n.ParseUrl(tag, common.NethermindCommitHash)

	return n.Install(url, true)
}

func (n *NethermindClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.NethermindConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	startFlags = n.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.NethermindConfigFileFlag)))

	return
}

func (n *NethermindClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}

func (n *NethermindClient) Version() (version string) {
	cmdVer := execVersionCmd(
		fmt.Sprintf("./%s/nethermind", n.FilePath()),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Nethermind version output to parse:

	// < v1.30:

	// 2024-11-19 11-41-29.6737|Nethermind starting initialization.
	// 2024-11-19 11-41-29.6915|Client version: Nethermind/v1.27.0+220b5b85/linux-x64/dotnet8.0.6
	// 2024-11-19 11-41-29.7003|Loading embedded plugins
	// ...
	// 2024-11-19 11-41-29.7460|Loading assembly Nethermind.Init.Snapshot
	// 2024-11-19 11-41-29.7463|  Found plugin type Nethermind.Init.Snapshot
	//
	// Version: 1.27.0+220b5b85
	// Commit: 220b5b856b1530482e957c002c9b24148a25f075
	// Build Date: 2024-06-21 11:48:18Z
	// OS: Linux x64
	// Runtime: .NET 8.0.6

	// > v1.30:
	// Version: 1.27.0+220b5b85
	// Commit: 220b5b856b1530482e957c002c9b24148a25f075
	// Build Date: 2024-06-21 11:48:18Z
	// OS: Linux x64
	// Runtime: .NET 8.0.6

	// Regexr: https://regexr.com/8bf90
	// Find the 'Version' line (with uppercase Version)
	expr := regexp.MustCompile(fmt.Sprintf(`Version: * %s\+[a-f0-9]{8}`, common.SemverExpressionRaw))
	s := expr.FindString(cmdVer)
	if s == "" {
		return VersionNotAvailable
	}

	// s: Version: 1.27.0+220b5b85
	// -> ...|v1.27.0+220b5b85
	splits := strings.Split(s, " ")
	s = splits[len(splits)-1]

	return fmt.Sprintf("v%s", strings.Split(s, "+")[0])
}

func (n *NethermindClient) FilePath() string {
	return nethermindFolder
}
