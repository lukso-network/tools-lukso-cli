package clients

import (
	"fmt"
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

type GethClient struct {
	*clientBinary
}

func NewGethClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *GethClient {
	return &GethClient{
		&clientBinary{
			name:           gethDependencyName,
			fileName:       gethFileName,
			commandPath:    gethCommandPath,
			baseUrl:        "https://gethstore.blob.core.windows.net/builds/geth-|OS|-|ARCH|-|TAG|-|COMMIT|.tar.gz",
			githubLocation: gethGithubLocation,
			buildInfo:      gethBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Geth dep.ExecutionClient
	_    dep.ExecutionClient = &GethClient{}
)

func (g *GethClient) Install(version string, isUpdate bool) error {
	url := g.ParseUrl(version, g.Commit())

	return g.installer.InstallTar(
		url,
		file.ClientsDir,
		g.FileName(),
		"geth-",
		isUpdate,
	)
}

func (g *GethClient) Update() (err error) {
	tag := g.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", g.name)

	// this commit hash is hardcoded, but since update should only be responsible for updating the client to
	// LUKSO supported version this is fine.

	return g.Install(tag, true)
}

func (g *GethClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.GethConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	ip := ctx.Context.Value(common.ConfigKey("ip"))

	startFlags = g.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.GethConfigFileFlag)))
	startFlags = append(startFlags, fmt.Sprintf("--nat=extip:%s", ip))

	return
}

func (g *GethClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}

func (g *GethClient) Version() (version string) {
	cmdVer := execVersionCmd(
		g.CommandPath(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Geth version output to parse:

	// geth version 1.14.7-stable-aa55f5ea

	// -> ...|1.14.7-stable-aa55f5ea
	s := strings.Split(cmdVer, " ")[2]
	// -> v1.14.7|...
	return fmt.Sprintf("v%s", strings.Split(s, "-")[0])
}
