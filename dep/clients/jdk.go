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
	"github.com/lukso-network/tools-lukso-cli/common/system"
)

const jdkFolder = file.ClientsDir + "/jdk"

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

	err = installAndExtractFromURL(jdkURL, "JDK", common.ClientDepsFolder, tarFormat, isUpdate)
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
