package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/apitypes"
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

func (e *ErigonClient) Update() (err error) {
	log.Info("⬇️  Fetching latest release for Erigon")

	latestErigonTag, err := fetchTag(e.githubLocation)
	if err != nil {
		return err
	}

	log.Infof("✅  Fetched latest release: %s", latestErigonTag)

	// since geth needs standard x.y.z semantic version to download (without "v" at the beginning) we need to strip it
	strippedTag := strings.TrimPrefix(latestErigonTag, "v")

	log.WithField("dependencyTag", latestErigonTag).Info("⬇️  Updating Erigon")

	url := e.ParseUrl(strippedTag, "")

	return e.Install(url, true)
}

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
	host := ctx.String(flags.ExecutionClientHost)
	port := ctx.Int(flags.ExecutionClientPort)
	if port == 0 {
		port = 8545 // default for LUKSO erigon config
	}

	url := fmt.Sprintf("http://%s:%d", host, port)

	reqBodyBytes, err := json.Marshal(apitypes.JsonRpcRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  "admin_peers",
		Params:  []string{},
	})
	if err != nil {
		return
	}

	reqBody := bytes.NewReader(reqBodyBytes)

	req, err := http.NewRequest(http.MethodGet, url, reqBody)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	peersResp := &apitypes.PeersJsonRpcResponse{}
	err = json.Unmarshal(respBodyBytes, peersResp)

	outbound = len(peersResp.Result)

	return
}
