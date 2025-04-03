package commands

import (
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	apierrors "github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/model"
)

// InitializeDirectory initializes a working directory for lukso node, with all configurations for all networks
func (c *commander) Init(ctx *cli.Context) error {
	if clients.IsAnyRunning() {
		return nil
	}

	req := model.CtxToApiInit(ctx)

	resp := c.handler.Init(req)
	err := resp.Error()
	switch {
	case errors.Is(err, apierrors.ErrCfgExists):
		log.Warn("This directory has already been initialized. If you want to reinitialize it, add '--reinit' flag.")
	}

	displayNetworksHardforkTimestamps()

	log.Info("✅  Working directory initialized! \n1. ⚙️  Use 'lukso install' to install clients. \n2. ▶️  Use 'lukso start' to start your node.")

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
