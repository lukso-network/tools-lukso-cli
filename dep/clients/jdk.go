package clients

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/installer"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/dep"
)

const jdkFolder = file.ClientsDir + "/jdk"

type JdkDep struct {
	baseUrl   string
	buildInfo buildInfo
	sourceDir string

	log       logger.Logger
	installer installer.Installer
}

var (
	Jdk JdkDep
	_   dep.Installer = &JdkDep{}
)

func NewJdk() *JdkDep {
	return &JdkDep{}
}

func (j *JdkDep) ParseUrl(tag, commitHash string) (url string) {
	url = j.baseUrl

	url = strings.ReplaceAll(url, "|TAG|", tag)
	url = strings.ReplaceAll(url, "|COMMIT|", commitHash)
	url = strings.ReplaceAll(url, "|OS|", j.Os())
	url = strings.ReplaceAll(url, "|ARCH|", j.Arch())

	return
}

func (j *JdkDep) Tag() string {
	return common.JdkTag
}

func (j *JdkDep) Commit() string {
	return common.JdkCommitHash
}

func (j *JdkDep) Os() string {
	return j.buildInfo.Os()
}

func (j *JdkDep) Arch() string {
	return j.buildInfo.Arch()
}

func (j *JdkDep) Install(version string, isUpdate bool) (err error) {
	j.log.Info("⬇️  Downloading JDK...")

	url := j.ParseUrl(j.Tag(), j.Commit())
	err = j.installer.InstallTar(url, j.sourceDir)
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

func (j *JdkDep) Update() error {
	return nil
}
