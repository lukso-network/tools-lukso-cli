package commands

import (
	"fmt"
	"github.com/lukso-network/tools-lukso-cli/dependencies/types"
	"gopkg.in/yaml.v2"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

var cfg = config.NewConfig(config.Path)

// installConfigGroup takes map of config dependencies and downloads them.
func installConfigGroup(configDependencies map[string]configs.ClientConfigDependency, isUpdate bool) error {
	for _, dependency := range configDependencies {
		err := dependency.Install(isUpdate)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was error while downloading %s file: %v", dependency.Name(), err), 1)
		}
	}

	return nil
}

func displayHardforkTimestamps(network, srcConfig string, epochZeroTimestamp uint64) (err error) {
	log.Infof("Your LUKSO Node has configs installed to run %s Hardforks on the following dates:", network)

	f, err := os.ReadFile(srcConfig)
	if err != nil {
		return
	}

	var (
		clconfig types.CLConfig

		shapellaTime    time.Time
		dencunTime      time.Time
		shapellaMessage string
		dencunMessage   string
		isValid         bool
	)

	err = yaml.Unmarshal(f, &clconfig)
	if err != nil {
		return
	}

	shapellaTime, isValid = utils.EthEpochToTimestamp(clconfig.ShapellaEpoch, epochZeroTimestamp)
	if !isValid {
		shapellaMessage = "TBA"
	} else {
		shapellaMessage = shapellaTime.Format(time.RFC1123Z)
	}

	dencunTime, isValid = utils.EthEpochToTimestamp(clconfig.DencunEpoch, epochZeroTimestamp)
	if !isValid {
		dencunMessage = "TBA"
	} else {
		dencunMessage = dencunTime.Format(time.RFC1123Z)
	}

	log.Infof("- Shapella: %v", shapellaMessage)
	log.Infof("- Dencun: %v\n\n", dencunMessage)

	return
}
