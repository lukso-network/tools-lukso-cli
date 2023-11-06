package clients

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
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
	tekuDepsFolder = "teku" // folder in which both teku and JDK are stored
	tekuFolder     = "teku" // folder in which teku is stored (in tekuDepsFolder)
	jdkFolder      = "jdk"  // folder in which JDK is stored (in tekuDepsFolder)

	tekuInstallURL = "https://artifacts.consensys.net/public/teku/raw/names/teku.tar.gz/versions/|TAG|/teku-|TAG|.tar.gz"
	jdkInstallURL  = "https://download.java.net/java/GA/jdk20.0.2/6e380f22cbe7469fa75fb448bd903d8e/9/GPL/openjdk-20.0.2_|OS|-|ARCH|_bin.tar.gz"
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

	err = installAndUntarFromURL(url, t.name, isUpdate)
	if err != nil {
		return
	}

	_, isInstalled := os.LookupEnv(system.JavaHomeEnv) // means that JDk is not set up
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

	permFunc := func(path string, d fs.DirEntry, err error) error {
		return os.Chmod(path, fs.ModePerm)
	}

	err = filepath.WalkDir(tekuDepsFolder, permFunc)
	if err != nil {
		return
	}

	return
}

func (t *TekuClient) Update() (err error) {
	log.Infof("⬇️  Fetching latest release for %s", t.name)

	latestTag, err := fetchTag(t.githubLocation)
	if err != nil {
		return err
	}

	log.Infof("✅  Fetched latest release: %s", latestTag)

	log.WithField("dependencyTag", latestTag).Infof("⬇️  Updating %s", t.name)

	url := t.ParseUrl(latestTag, "")

	return t.Install(url, true)
}

func (t *TekuClient) FilePath() string {
	return tekuDepsFolder
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

	command := exec.Command(fmt.Sprintf("./%s/%s/bin/teku", tekuDepsFolder, tekuFolder), arguments...)

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

func untarDir(dst string, t *tar.Reader) error {
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dst, header.Name)

		// for the sake of compatibility with updated versions remove the tag from the tarred file - teku/teku-xx.x.x => teku/teku, same with jdk
		switch {
		case strings.Contains(header.Name, tekuFolder):
			newHeader := replaceRootFolderName(header.Name, tekuFolder)
			path = filepath.Join(dst, newHeader)

		case strings.Contains(header.Name, jdkFolder):
			newHeader := replaceRootFolderName(header.Name, jdkFolder)
			path = filepath.Join(dst, newHeader)
		}

		info := header.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}

			continue
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, info.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, t)
		if err != nil {
			return err
		}
	}

	return nil
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

	err = installAndUntarFromURL(jdkURL, "JDK", isUpdate)
	if err != nil {
		return err
	}

	luksoNodeDir, err := os.Getwd()
	if err != nil {
		return
	}

	javaHomeVal := fmt.Sprintf("%s/%s/%s", luksoNodeDir, tekuDepsFolder, jdkFolder)

	log.Infof("⚙️  To continue working with Teku please export the JAVA_HOME environment variable.\n"+
		"The recommended way is to add the following line:\n\n"+
		"export JAVA_HOME=%s\n\n"+
		"To the bash startup file of your choosing (like .bashrc)", javaHomeVal)

	return
}

func installAndUntarFromURL(url, name string, isUpdate bool) (err error) {
	response, err := http.Get(url)
	if nil != err {
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNotFound {
		log.Warnf("⚠️  File under URL %s not found - skipping...", url)

		return nil
	}

	if http.StatusOK != response.StatusCode {
		return fmt.Errorf(
			"❌  Invalid response when downloading file at URL: %s. Response code: %s",
			url,
			response.Status,
		)
	}

	g, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}

	defer func() {
		_ = g.Close()
	}()

	tarReader := tar.NewReader(g)

	err = untarDir(Teku.FilePath(), tarReader)
	if err != nil {
		return
	}

	switch isUpdate {
	case true:
		log.Infof("✅  %s updated!\n\n", name)
	case false:
		log.Infof("✅  %s downloaded!\n\n", name)
	}

	return
}

func replaceRootFolderName(folder, targetRootName string) (path string) {
	splitHeader := strings.Split(folder, "/") //this assumes no / at the beginning of folder - not the case in tarred files we are interested in

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
