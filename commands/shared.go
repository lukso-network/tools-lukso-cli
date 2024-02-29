package commands

import (
	"fmt"

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
