package clients

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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

const nimbus2Folder = common.ClientDepsFolder + "/nimbus2"

type Nimbus2Client struct {
	*clientBinary
}

func NewNimbus2Client(
	log logger.Logger,
	file file.Manager,
	installer installer.Installer,
	pid pid.Pid,
) *Nimbus2Client {
	return &Nimbus2Client{
		&clientBinary{
			name:           nimbus2DependencyName,
			fileName:       nimbus2FileName,
			commandPath:    nimbus2CommandPath,
			baseUrl:        "https://github.com/status-im/nimbus-eth2/releases/download/v|TAG|/nimbus-eth2_|OS|_|ARCH|_|TAG|_|COMMIT|.tar.gz",
			githubLocation: nimbus2GithubLocation,
			buildInfo:      nimbus2BuildInfo,
			log:            log,
			file:           file,
			installer:      installer,
			pid:            pid,
		},
	}
}

var (
	Nimbus2 dep.ConsensusClient
	_       dep.ConsensusClient = &Nimbus2Client{}
)

func (n *Nimbus2Client) Install(version string, isUpdate bool) (err error) {
	url := n.ParseUrl(version, n.Commit())

	return n.installer.InstallTar(
		url,
		file.ClientsDir,
		n.FileName(),
		"nimbus-eth2_",
		isUpdate,
	)
}

func (n *Nimbus2Client) Update() (err error) {
	tag := n.Tag()

	log.WithField("dependencyTag", tag).Infof("⬇️  Updating %s", n.name)

	return n.Install(tag, true)
}

func (n *Nimbus2Client) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = n.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.Nimbus2ConfigFileFlag)))
	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	}

	return
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

	err = n.logFile(n.FileName(), command)
	if err != nil {
		log.Errorf("There was an error while preparing a log file for %s: %v", n.Name(), err)
	}

	log.Infof("🔄  Starting %s", n.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, n.FileName())
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
		n.CommandPath(),
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
