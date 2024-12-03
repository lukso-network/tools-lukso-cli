package clients

import (
	"bytes"
	"fmt"
	"os/exec"
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
			baseUrl:        "https://github.com/sigp/lighthouse/releases/download/|TAG|/lighthouse-|TAG|-|ARCH|-|OS-NAME|-|OS|-portable.tar.gz",
			githubLocation: lighthouseGithubLocation,
		},
	}
}

var Lighthouse = NewLighthouseClient()

var _ ClientBinaryDependency = &LighthouseClient{}

func (l *LighthouseClient) Update() (err error) {
	tag := ClientVersions[l.Name()]
	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", l.name)

	url := l.ParseUrl(tag, "")

	return l.Install(url, true)
}

func (l *LighthouseClient) ParseUrl(tag, commitHash string) (url string) {
	var (
		systemName string
		urlSystem  = system.Os
		arch       string
	)

	fallback := func() {
		log.Info("⚠️  Unknown OS detected: proceeding with x86_64 as a default arch")
		arch = "x86_64"
	}

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

	switch system.Os {
	case system.Ubuntu, system.Macos:
		buf := new(bytes.Buffer)

		uname := exec.Command("uname", "-m")
		uname.Stdout = buf

		err := uname.Run()
		if err != nil {
			fallback()

			break
		}

		arch = strings.Trim(buf.String(), "\n\t ")

	default:
		fallback()
	}

	if arch != "x86_64" && arch != "aarch64" {
		fallback()
	}

	url = l.baseUrl
	url = strings.ReplaceAll(url, "|TAG|", tag)
	url = strings.ReplaceAll(url, "|OS|", urlSystem)
	url = strings.ReplaceAll(url, "|OS-NAME|", systemName) // for lighthouse
	url = strings.ReplaceAll(url, "|COMMIT|", commitHash)
	url = strings.ReplaceAll(url, "|ARCH|", arch)

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

	isSlasher := !ctx.Bool(flags.NoSlasherFlag)
	isValidator := ctx.Bool(flags.ValidatorFlag)
	if isSlasher && isValidator {
		startFlags = append(startFlags, "--slasher")
		startFlags = append(startFlags, "--slasher-max-db-size=16")
		startFlags = append(startFlags, "--slasher-history-length=256")
		startFlags = append(startFlags, fmt.Sprintf("--slasher-dir=%s", ctx.String(flags.ConsensusDatadirFlag)))
	}

	return
}

func (l *LighthouseClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 4000)
}

func (l *LighthouseClient) Version() (version string) {
	cmdVer := execVersionCmd(
		l.CommandName(),
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
