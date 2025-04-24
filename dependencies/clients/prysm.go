package clients

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

type PrysmClient struct {
	*clientBinary
}

func NewPrysmClient() *PrysmClient {
	return &PrysmClient{
		&clientBinary{
			name:           prysmDependencyName,
			commandName:    "prysm",
			baseUrl:        "https://github.com/prysmaticlabs/prysm/releases/download/|TAG|/beacon-chain-|TAG|-|OS|-|ARCH|",
			githubLocation: prysmaticLabsGithubLocation,
		},
	}
}

var Prysm = NewPrysmClient()

var _ Client = &PrysmClient{}

func (p *PrysmClient) PrepareStartFlags(ctx *cli.Context) (startFlags []string, err error) {
	genesisExists := utils.FlagFileExists(ctx, flags.GenesisStateFlag)
	prysmConfigExists := utils.FlagFileExists(ctx, flags.PrysmConfigFileFlag)
	chainConfigExists := utils.FlagFileExists(ctx, flags.PrysmChainConfigFileFlag)
	if !genesisExists || !prysmConfigExists || !chainConfigExists {
		err = errors.ErrFlagPathInvalid

		return
	}

	startFlags = p.ParseUserFlags(ctx)

	// terms of use already accepted during installation
	startFlags = append(startFlags, "--accept-terms-of-use")
	startFlags = append(startFlags, fmt.Sprintf("--config-file=%s", ctx.String(flags.PrysmConfigFileFlag)))
	startFlags = append(startFlags, "--contract-deployment-block=0")

	if ctx.String(flags.TransactionFeeRecipientFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--suggested-fee-recipient=%s", ctx.String(flags.TransactionFeeRecipientFlag)))
	}
	if ctx.String(flags.GenesisStateFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--genesis-state=%s", ctx.String(flags.GenesisStateFlag)))
	}
	if ctx.String(flags.PrysmChainConfigFileFlag) != "" {
		startFlags = append(startFlags, fmt.Sprintf("--chain-config-file=%s", ctx.String(flags.PrysmChainConfigFileFlag)))
	}

	isSlasher := !ctx.Bool(flags.NoSlasherFlag)
	isValidator := ctx.Bool(flags.ValidatorFlag)
	if isSlasher && isValidator {
		startFlags = append(startFlags, "--slasher")
	}

	return
}

func (p *PrysmClient) Peers(ctx *cli.Context) (outbound, inbound int, err error) {
	return defaultConsensusPeers(ctx, 3500)
}

func (p *PrysmClient) Version() (version string) {
	cmdVer := execVersionCmd(
		p.CommandName(),
	)

	if cmdVer == VersionNotAvailable {
		return VersionNotAvailable
	}

	// Prysm version output to parse:

	// beacon-chain version Prysm/v5.0.4/3b184f43c86baf6c36478f65a5113e7cf0836d41. Built at: 2024-06-21 00:26:00+00:00

	// -> ...|Prysm/v5.0.4/3b184f43c86baf6c36478f65a5113e7cf0836d41.|...
	s := strings.Split(cmdVer, " ")[2]
	// -> ...|v5.0.4|...
	return strings.Split(s, "/")[1]
}
