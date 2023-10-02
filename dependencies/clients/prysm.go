package clients

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"

	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/apitypes"
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

var _ ClientBinaryDependency = &PrysmClient{}

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
	host := ctx.String(flags.ConsensusClientHost)
	port := ctx.Int(flags.ConsensusClientPort)
	if port == 0 {
		port = 3500 // default for LUKSO prysm config
	}

	url := fmt.Sprintf("http://%s:%d/eth/v1/node/peers", host, port)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	peersResp := &apitypes.PeersBeaconAPIResponse{}
	err = json.Unmarshal(respBodyBytes, peersResp)
	if err != nil {
		return
	}

	for _, peer := range peersResp.Data {
		if peer.State != peerStateConnected {
			continue
		}

		switch peer.Direction {
		case peerDirectionInbound:
			inbound++
		case peerDirectionOutbound:
			outbound++
		default:
			log.Errorf("Unknown direction for %s peer", p.name)
		}
	}

	return
}
