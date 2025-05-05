package clients

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const nimbus2Folder = common.ClientDepsFolder + "/nimbus2"

type Nimbus2Client struct {
	*clientBinary
}

func NewNimbus2Client() *Nimbus2Client {
	return &Nimbus2Client{
		&clientBinary{
			name:           nimbus2DependencyName,
			fileName:       "nimbus2",
			baseUrl:        "https://github.com/status-im/nimbus-eth2/releases/download/v|TAG|/nimbus-eth2_|OS|_|ARCH|_|TAG|_|COMMIT|.tar.gz",
			githubLocation: nimbus2GithubLocation,
			buildInfo:      nimbus2BuildInfo,
		},
	}
}

var Nimbus2 = NewNimbus2Client()

var _ Client = &Nimbus2Client{}

func (n *Nimbus2Client) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = n.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.Nimbus2ConfigFileFlag)))
	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	}

	return
}

func (n *Nimbus2Client) ParseUrl(tag, commitHash string) (url string) {
	var (
		urlSystem string
		arch      = system.GetArch()
	)

	switch system.Os {
	case system.Ubuntu:
		urlSystem = "Linux"
	case system.Macos:
		urlSystem = "macOS"
	default:
		urlSystem = "Linux"
	}

	if (arch == "arm" || arch == "aarch64") && system.Os == system.Ubuntu {
		arch = "arm64v8"
	}
	if arch == "x86_64" {
		// x86_64 not present as a released tar - default to amd64 instead
		arch = "amd64"
	}

	url = n.baseUrl
	url = strings.ReplaceAll(url, "|TAG|", tag)
	url = strings.ReplaceAll(url, "|OS|", urlSystem)
	url = strings.ReplaceAll(url, "|COMMIT|", commitHash)
	url = strings.ReplaceAll(url, "|ARCH|", arch)

	return
}

func (n *Nimbus2Client) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(n.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", n.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	err = installAndExtractFromURL(url, n.name, common.ClientDepsFolder, tarFormat, isUpdate)
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

func (n *Nimbus2Client) Update() (err error) {
	tag := n.getVersion()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", n.name)

	url := n.ParseUrl(tag, common.Nimbus2CommitHash)

	return n.Install(url, true)
}

func (n *Nimbus2Client) FilePath() string {
	return nimbus2Folder
}

func (n *Nimbus2Client) Start(ctx *cli.Context, arguments []string) (err error) {
	if n.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", n.Name())

		err = n.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", n.Name())
	}

	// in Nimbus case, we need to prepare the data dir using a syncing command:
	// nimbus_beacon_node trustedNodeSync
	if ctx.Bool(flags.CheckpointSyncFlag) {
		log.Info("⚙️  Downloading checkpointed state for Nimbus")
		network := "mainnet"
		if ctx.Bool(flags.TestnetFlag) {
			network = "testnet"
		}

		checkpointURL := fmt.Sprintf("https://checkpoints.%s.lukso.network", network)
		args := []string{
			"trustedNodeSync",
			fmt.Sprintf("--trusted-node-url=%s", checkpointURL),
			fmt.Sprintf("--data-dir=%s", ctx.String(flags.Nimbus2DatadirFlag)),
			fmt.Sprintf("--network=%s", ctx.String(flags.Nimbus2NetworkFlag)),
		}

		syncCommand := exec.Command(fmt.Sprintf("./%s/build/nimbus_beacon_node", n.FilePath()), args...)

		syncCommand.Stdout = os.Stdout
		syncCommand.Stderr = os.Stderr

		err = syncCommand.Run()
		if err != nil {
			return
		}
	}

	command := exec.Command(fmt.Sprintf("./%s/build/nimbus_beacon_node", n.FilePath()), arguments...)

	err = prepareLogFile(ctx, command, n.CommandName())
	if err != nil {
		log.Errorf("There was an error while preparing a log file for %s: %v", n.Name(), err)
	}

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

func (n *Nimbus2Client) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 5052)
}

func (n *Nimbus2Client) Version() (version string) {
	cmdVer := execVersionCmd(
		fmt.Sprintf("./%s/build/nimbus_beacon_node", n.FilePath()),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Nimbus-eth2 version output to parse:

	// Nimbus beacon node v24.10.0-c4037d-stateofus
	// Copyright (c) 2019-2024 Status Research & Development GmbH
	//
	// eth2 specification v1.5.0-alpha.8
	//
	// Nim Compiler Version 2.0.10 [Linux: amd64]

	// -> ...|v24.10.0-c4037d-stateofus|...
	s := strings.Split(cmdVer, " ")[3]
	// -> ...|v24.10.0|...
	return strings.Split(s, "-")[0]
}
