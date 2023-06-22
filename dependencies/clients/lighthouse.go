package clients

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type LighthouseClient struct {
	*clientBinary
}

func NewLighthouseClient() *LighthouseClient {
	return &LighthouseClient{
		&clientBinary{
			name:           lighthouseDependencyName,
			commandName:    "lighthouse",
			baseUrl:        "https://github.com/sigp/lighthouse/releases/download/|TAG|/lighthouse-|TAG|-x86_64-|OS-NAME|-|OS|-portable.tar.gz",
			githubLocation: lighthouseGithubLocation,
		},
	}
}

var Lighthouse = NewLighthouseClient()

var _ ClientBinaryDependency = &LighthouseClient{}

func (l *LighthouseClient) Update() (err error) {
	log.Infof("⬇️  Fetching latest release for %s", l.name)

	latestTag, err := fetchTag(l.githubLocation)
	if err != nil {
		return err
	}

	log.Infof("✅  Fetched latest release: %s", latestTag)

	log.WithField("dependencyTag", latestTag).Infof("⬇️  Updating %s", l.name)

	url := l.ParseUrl(latestTag, "")

	return l.Install(url, true)
}

func (l *LighthouseClient) ParseUrl(tag, commitHash string) (url string) {
	// for lighthouse
	var (
		systemName string
		urlSystem  = system.Os
	)
	switch system.Os {
	case system.Ubuntu:
		systemName = "unknown"
		urlSystem += "-gnu"
	case system.Macos:
		systemName = "apple"
	default:
		systemName = "unknown"
		urlSystem += "-gnu"
	}

	url = l.baseUrl
	url = strings.Replace(url, "|TAG|", tag, -1)
	url = strings.Replace(url, "|OS|", urlSystem, -1)
	url = strings.Replace(url, "|OS-NAME|", systemName, -1)
	url = strings.Replace(url, "|COMMIT|", commitHash, -1)
	url = strings.Replace(url, "|ARCH|", system.Arch, -1)

	return
}

func (l *LighthouseClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	defaults, err := config.LoadLighthouseConfig(ctx.String(flags.LighthouseConfigFileFlag))
	if err != nil {
		return
	}

	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		defaults = append(defaults, "--suggested-fee-recipient", ctx.String(flags.TransactionFeeRecipientFlag))
	}

	userFlags := l.ParseUserFlags(ctx)

	startFlags = mergeFlags(userFlags, defaults)

	startFlags = append([]string{"bn"}, startFlags...)

	return
}

// mergeFlags takes 2 arrays of flag sets, and replaces all flags present in both with user input. Appends default otherwise
func mergeFlags(userFlags, configFlags []string) (startFlags []string) {
	for cfgI, configArg := range configFlags {
		for usrI, userArg := range userFlags {
			if !strings.HasPrefix(configArg, "--") || !strings.HasPrefix(userArg, "--") {
				continue
			}

			if configArg != userArg {
				continue
			}

			if usrI == len(userFlags)-1 || cfgI == len(configFlags)-1 {
				configFlags = pop(configFlags, cfgI)

				continue
			}

			if !strings.HasPrefix(userFlags[usrI+1], "--") {
				configFlags = pop(configFlags, cfgI)
				configFlags = pop(configFlags, cfgI)

				continue
			}
		}
	}

	var mergedFlags []string
	mergedFlags = append(mergedFlags, userFlags...)
	mergedFlags = append(mergedFlags, configFlags...)

	// merge flags with values using =
	for i, arg := range mergedFlags {
		if i == len(mergedFlags)-1 {
			break
		}

		if strings.HasPrefix(arg, "--") {
			if !strings.HasPrefix(mergedFlags[i+1], "--") {
				startFlags = append(startFlags, fmt.Sprintf("%s=%s", arg, mergedFlags[i+1]))

				continue
			}

			startFlags = append(startFlags, arg)
		}
	}

	return
}

func pop(arr []string, i int) []string {
	l := arr[:i]
	r := arr[i+1:]
	l = append(l, r...)

	return append([]string{}, l...)
}
