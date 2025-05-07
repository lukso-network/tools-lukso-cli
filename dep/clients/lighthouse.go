package clients

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dep"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

type LighthouseClient struct {
	*clientBinary
}

func NewLighthouseClient(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *LighthouseClient {
	return &LighthouseClient{
		&clientBinary{
			name:           lighthouseDependencyName,
			fileName:       lighthouseFileName,
			commandPath:    lighthouseCommandPath,
			baseUrl:        "https://github.com/sigp/lighthouse/releases/download/|TAG|/lighthouse-|TAG|-|ARCH|-|OS|.tar.gz",
			githubLocation: lighthouseGithubLocation,
			buildInfo:      lighthouseBuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Lighthouse dep.ConsensusClient
	_          dep.ConsensusClient = &LighthouseClient{}
)

func (l *LighthouseClient) Install(version string, isUpdate bool) error {
	url := l.ParseUrl(version, l.Commit())

	return l.installer.InstallTar(url, file.ClientsDir, l.FileName(), "lighthouse-")
}

func (l *LighthouseClient) Update() (err error) {
	tag := l.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", l.name)

	return l.Install(tag, true)
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

	isSlasher := !ctx.Bool(flags.NoSlasherFlag)
	isValidator := ctx.Bool(flags.ValidatorFlag)
	if isSlasher && isValidator {
		startFlags = append(startFlags, "--slasher")
		startFlags = append(startFlags, "--slasher-max-db-size=16")
		startFlags = append(startFlags, "--slasher-history-length=256")
		startFlags = append(startFlags, fmt.Sprintf("--slasher-dir=%s", ctx.String(flags.LighthouseDatadirFlag)))
	}

	return
}

func (l *LighthouseClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 4000)
}

func (l *LighthouseClient) Version() (version string) {
	cmdVer := execVersionCmd(
		l.CommandPath(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Lighthouse version output to parse:

	// Lighthouse v5.2.1-9e12c21
	// BLS library: blst-portable
	// SHA256 hardware acceleration: false
	// Allocator: jemalloc
	// Profile: maxperf
	// Specs: mainnet (true), minimal (false), gnosis (true)

	// -> ...|v5.2.1-9e12c21|...
	s := strings.Split(cmdVer, " ")[1]
	// -> v5.2.1|...
	return strings.Split(s, "-")[0]
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
