package clients

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type GethClient struct {
	*clientBinary
}

func NewGethClient() *GethClient {
	return &GethClient{
		&clientBinary{
			name:           gethDependencyName,
			commandName:    "geth",
			baseUrl:        "https://gethstore.blob.core.windows.net/builds/geth-|OS|-|ARCH|-|TAG|-|COMMIT|.tar.gz",
			githubLocation: gethGithubLocation,
		},
	}
}

var Geth = NewGethClient()

var _ ClientBinaryDependency = &GethClient{}

func (g *GethClient) ParseUrl(tag, commitHash string) (url string) {
	url = g.baseUrl

	if g.name == gethDependencyName && system.Os == system.Macos {
		url = strings.Replace(url, "|ARCH|", "amd64", -1)
	}

	url = strings.Replace(url, "|TAG|", tag, -1)
	url = strings.Replace(url, "|OS|", system.Os, -1)
	url = strings.Replace(url, "|COMMIT|", commitHash, -1)
	url = strings.Replace(url, "|ARCH|", system.Arch, -1)

	return
}

func (g *GethClient) Update() (err error) {
	log.Info("⬇️  Fetching latest release for Geth")

	latestGethTag, latestGethCommitHash, err := fetchTagAndCommitHash(g.githubLocation)
	if err != nil {
		return err
	}

	log.Infof("✅  Fetched latest release: %s", latestGethTag)

	// since geth needs standard x.y.z semantic version to download (without "v" at the beginning) we need to strip it
	strippedTag := strings.TrimPrefix(latestGethTag, "v")
	// and we need to take only 8 first characters from commit hash
	chars := []rune(latestGethCommitHash)[:8]
	latestGethCommitHash = string(chars)

	log.WithField("dependencyTag", latestGethTag).Info("⬇️  Updating Geth")

	url := g.ParseUrl(strippedTag, latestGethCommitHash)

	return g.Install(url, true)
}

func (g *GethClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.GethConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	startFlags = g.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.GethConfigFileFlag)))

	return
}

func (g *GethClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}
