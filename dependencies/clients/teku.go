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

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
	"github.com/lukso-network/tools-lukso-cli/pid"
)

const (
	tekuFolder = clientDepsFolder + "/teku" // folder in which teku is stored (in tekuDepsFolder)

	tekuInstallURL = "https://artifacts.consensys.net/public/teku/raw/names/teku.tar.gz/versions/|TAG|/teku-|TAG|.tar.gz"
	jdkInstallURL  = "https://download.java.net/java/GA/jdk22.0.1/c7ec1332f7bb44aeba2eb341ae18aca4/8/GPL/openjdk-22.0.1_|OS|-|ARCH|_bin.tar.gz"
)

type TekuClient struct {
	*clientBinary
}

func NewTekuClient() *TekuClient {
	return &TekuClient{
		&clientBinary{
			name:           tekuDependencyName,
			commandName:    "teku",
			baseUrl:        tekuInstallURL,
			githubLocation: tekuGithubLocation,
		},
	}
}

var Teku = NewTekuClient()

var _ ClientBinaryDependency = &TekuClient{}

func (t *TekuClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	startFlags = t.ParseUserFlags(ctx)

	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.TekuConfigFileFlag)))
	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--validators-proposer-default-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	}

	return
}

func (t *TekuClient) Install(url string, isUpdate bool) (err error) {
	if utils.FileExists(t.FilePath()) && !isUpdate {
		message := fmt.Sprintf("You already have the %s client installed, do you want to override your installation? [Y/n]: ", t.Name())
		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return nil
		}
	}

	err = installAndExtractFromURL(url, t.name, clientDepsFolder, tarFormat, isUpdate)
	if err != nil {
		return
	}

	permFunc := func(path string, d fs.DirEntry, err error) error {
		return os.Chmod(path, fs.ModePerm)
	}

	err = filepath.WalkDir(t.FilePath(), permFunc)
	if err != nil {
		return
	}

	isInstalled := isJdkInstalled()
	if !isInstalled {
		message := "Teku is written in Java. This means that to use it you need to have:\n" +
			"- JDK installed on your computer\n" +
			"- JAVA_HOME environment variable set\n" +
			"Do you want to install and set up JDK along with Teku? [Y/n]\n>"

		input := utils.RegisterInputWithMessage(message)
		if !strings.EqualFold(input, "y") && input != "" {
			log.Info("⏭️  Skipping installation...")

			return
		}

		err = setupJava(isUpdate)
		if err != nil {
			return
		}
	}

	return
}

func (t *TekuClient) FilePath() string {
	return tekuFolder
}

func (t *TekuClient) Start(ctx *cli.Context, arguments []string) (err error) {
	if t.IsRunning() {
		log.Infof("🔄️  %s is already running - stopping first...", t.Name())

		err = t.Stop()
		if err != nil {
			return
		}

		log.Infof("🛑  Stopped %s", t.Name())
	}

	command := exec.Command(fmt.Sprintf("./%s/bin/teku", t.FilePath()), arguments...)

	var (
		logFile  *os.File
		fullPath string
	)

	logFolder := ctx.String(flags.LogFolderFlag)
	if logFolder == "" {
		return utils.Exit(fmt.Sprintf("%v- %s", errors.ErrFlagMissing, flags.LogFolderFlag), 1)
	}

	fullPath, err = utils.PrepareTimestampedFile(logFolder, t.CommandName())
	if err != nil {
		return
	}

	err = os.WriteFile(fullPath, []byte{}, 0750)
	if err != nil {
		return
	}

	logFile, err = os.OpenFile(fullPath, os.O_RDWR, 0750)
	if err != nil {
		return
	}

	command.Stdout = logFile
	command.Stderr = logFile

	log.Infof("🔄  Starting %s", t.Name())
	err = command.Start()
	if err != nil {
		return
	}

	pidLocation := fmt.Sprintf("%s/%s.pid", pid.FileDir, t.CommandName())
	err = pid.Create(pidLocation, command.Process.Pid)

	time.Sleep(1 * time.Second)

	log.Infof("✅  %s started!", t.Name())

	return
}

func (t *TekuClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 5051)
}

func setupJava(isUpdate bool) (err error) {
	log.Info("⬇️  Downloading JDK...")

	var systemOs, arch string
	switch system.Os {
	case system.Ubuntu:
		systemOs = "linux"
	case system.Macos:
		systemOs = "macos"
	}

	arch = system.GetArch()

	if arch == "x86_64" {
		arch = "x64"
	}
	if arch != "aarch64" && arch != "x64" {
		log.Warnf("⚠️  x64 or aarch64 architecture is required to continue - skipping ...")

		return
	}

	jdkURL := strings.Replace(jdkInstallURL, "|OS|", systemOs, -1)
	jdkURL = strings.Replace(jdkURL, "|ARCH|", arch, -1)

	err = installAndExtractFromURL(jdkURL, "JDK", clientDepsFolder, tarFormat, isUpdate)
	if err != nil {
		return err
	}

	luksoNodeDir, err := os.Getwd()
	if err != nil {
		return
	}

	javaHomeVal := fmt.Sprintf("%s/%s", luksoNodeDir, jdkFolder)

	permFunc := func(path string, d fs.DirEntry, err error) error {
		return os.Chmod(path, fs.ModePerm)
	}

	err = filepath.WalkDir(jdkFolder, permFunc)
	if err != nil {
		return
	}

	log.Infof("⚙️  To continue working with Java clients please export the JAVA_HOME environment variable.\n"+
		"The recommended way is to add the following line:\n\n"+
		"export JAVA_HOME=%s\n\n"+
		"To the bash startup file of your choosing (like .bashrc)", javaHomeVal)

	return
}

func replaceRootFolderName(folder, targetRootName string) (path string) {
	splitHeader := strings.Split(folder, "/") // this assumes no / at the beginning of folder - not the case in tarred files we are interested in

	switch len(splitHeader) {
	case 0:
		return
	case 1:
		return targetRootName
	default:
		break
	}

	strippedHeader := splitHeader[1:]
	splitResolvedPath := append([]string{targetRootName}, strippedHeader...)

	return strings.Join(splitResolvedPath, "/")
}
