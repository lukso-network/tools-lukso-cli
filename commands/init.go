package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/urfave/cli/v2"

	apierrors "github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/model"
)

// InitializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func (c *commander) Init(ctx *cli.Context) error {
	req := model.CtxToApiInit(ctx)

	resp := c.handler.Init(req)
	err := resp.Error
	if err != nil {
		c.log.Warn(fmt.Sprintf("Unable to initialize directory: %v", err))

		if errors.Is(err, apierrors.ErrCfgExists) {
			c.log.Warn("To reinitialize directory, rerun the 'lukso init' command with '--reinit' flag")

			return nil
		}
	}

	displayNetworksHardforkTimestamps()

	c.log.Info("✅  Working directory initialized! \n1. ⚙️  Use 'lukso install' to install clients. \n2. ▶️  Use 'lukso start' to start your node.")

	return nil
}

func getPublicIP() (ip string, err error) {
	url := "https://ipv4.ident.me"
	res, err := http.Get(url)
	if err != nil {
		return
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	ip = string(respBody)

	return
}
