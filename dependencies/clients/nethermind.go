package clients

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

const (
	nethermindDepsFolder = "nethermind"
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

	err = installAndExtractFromURL(url, n.name, n.FilePath(), zipFormat, isUpdate)
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

func (n *NethermindClient) ParseUrl(tag, commitHash string) (url string) {
	url = n.baseUrl
	osName := system.Os
	archName := getUnameArch()

	if osName == system.Macos {
		osName = "macos"
	}

	if archName == "x86_64" {
		archName = "x64"
	}

	if n.name == gethDependencyName && system.Os == system.Macos {
		url = strings.Replace(url, "|ARCH|", "amd64", -1)
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
	url := n.ParseUrl(tag, common.GethCommitHash)

	return n.Install(url, true)
}

func (n *NethermindClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.GethConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	startFlags = n.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.GethConfigFileFlag)))

	return
}

func (n *NethermindClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}
