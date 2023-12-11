package clients

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type ErigonClient struct {
	*clientBinary
}

func NewErigonClient() *ErigonClient {
	return &ErigonClient{
		&clientBinary{
			name:           erigonDependencyName,
			commandName:    "erigon",
			baseUrl:        "https://github.com/ledgerwatch/erigon/releases/download/v|TAG|/erigon_|TAG|_|OS|_|ARCH|.tar.gz",
			githubLocation: erigonGithubLocation,
		},
	}
}

var Erigon = NewErigonClient()

var _ ClientBinaryDependency = &ErigonClient{}

func (e *ErigonClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	if !utils.FlagFileExists(ctx, flags.ErigonConfigFileFlag) {
		err = errors.ErrFlagMissing

		return
	}

	startFlags = e.ParseUserFlags(ctx)
	startFlags = append(startFlags, fmt.Sprintf("--config=%s", ctx.String(flags.ErigonConfigFileFlag)))

	return
}

func (e *ErigonClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultExecutionPeers(ctx, 8545)
}
