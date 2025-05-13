package clients

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
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
			fileName:       besuFileName,
			commandPath:    besuCommandPath,
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

	return b.installer.InstallTar(
		url,
		file.ClientsDir,
		b.FileName(),
		"besu-",
		isUpdate,
	)
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

func (b *BesuClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}

func (b *BesuClient) Version() (version string) {
	cmdVer := execVersionCmd(
		b.CommandPath(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Besu version output to parse:

	// besu/v24.7.0/linux-x86_64/oracle_openjdk-java-22
	return strings.Split(cmdVer, "/")[1]
}
