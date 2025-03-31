package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
)

var (
	LogsFlags          []cli.Flag
	ExecutionLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: configs.MainnetLogs,
		},
	}
	ConsensusLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: configs.MainnetLogs,
		},
	}
	ValidatorLogsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  LogFolderFlag,
			Usage: "Directory to output logs into",
			Value: configs.MainnetLogs,
		},
	}
)

func init() {
	ExecutionLogsFlags = append(ExecutionLogsFlags, NetworkFlags...)
	ConsensusLogsFlags = append(ConsensusLogsFlags, NetworkFlags...)
	ValidatorLogsFlags = append(ValidatorLogsFlags, NetworkFlags...)
}
