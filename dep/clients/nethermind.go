package clients

import (
	"fmt"
	"regexp"
	"strings"

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
	nethermindFolder = common.ClientDepsFolder + "/nethermind"
)

type NethermindClient struct {
	*clientBinary
}

func NewNethermindClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *NethermindClient {
	return &NethermindClient{
		&clientBinary{
			name:           nethermindDependencyName,
			fileName:       nethermindFileName,
			commandPath:    nethermindCommandPath,
			baseUrl:        "https://github.com/NethermindEth/nethermind/releases/download/|TAG|/nethermind-|TAG|-|COMMIT|-|OS|-|ARCH|.zip",
			githubLocation: nethermindGithubLocation,
			buildInfo:      nethermindBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Nethermind dep.ExecutionClient
	_          dep.ExecutionClient = &NethermindClient{}
)

func (n *NethermindClient) Install(version string, isUpdate bool) (err error) {
	url := n.ParseUrl(version, n.Commit())

	return n.installer.InstallZip(url, n.FileDir(), isUpdate)
}

func (n *NethermindClient) Update() (err error) {
	tag := n.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", n.name)

	return n.Install(tag, true)
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
		n.CommandPath(),
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
