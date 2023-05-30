package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateJwtSecret(dest string) error {
	log.Info("🔄  Creating new JWT secret")
	jwtDir := truncateFileFromDir(dest)

	err := os.MkdirAll(jwtDir, configs.ConfigPerms)
	if err != nil {
		return err
	}

	secretBytes := make([]byte, 32)

	_, err = rand.Read(secretBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, []byte(hex.EncodeToString(secretBytes)), configs.ConfigPerms)
	if err != nil {
		return err
	}

	log.Infof("✅  New JWT secret created in %s", dest)

	return nil
}

func PrepareTimestampedFile(logDir, logFileName string) (logFile string, err error) {
	err = os.MkdirAll(logDir, configs.ConfigPerms)
	if err != nil {
		return
	}

	t := time.Now().Format("2006-01-02_15-04-05")

	logFile = fmt.Sprintf("%s/%s_%s.log", logDir, logFileName, t)

	return
}

// truncateFileFromDir removes file name from its path.
// Example: /path/to/foo/foo.txt => /path/to/foo
func truncateFileFromDir(filePath string) string {
	segments := strings.Split(filePath, "/")

	return strings.Join(segments[:len(segments)-1], "/")
}
