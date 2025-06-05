package commands

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common/display/input"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func (c *commander) Install(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	c.display.AddInputHandler(input.ClientInstallHandler())
	cls, ok := c.display.InputHandler().Get().([]string)
	if !ok {
		return
	}

	var (
		consensusTag    string
		executionTag    string
		consensusCommit string
		executionCommit string
		validatorName   string
	)

	switch cls[0] {
	case clients.Prysm.Name():
		consensusTag = ctx.String(flags.PrysmTagFlag)
		validatorName = clients.PrysmValidator.Name()

	case clients.Lighthouse.Name():
		consensusTag = ctx.String(flags.LighthouseTagFlag)
		validatorName = clients.PrysmValidator.Name()

	case clients.Teku.Name():
		consensusTag = ctx.String(flags.TekuTagFlag)
		validatorName = clients.PrysmValidator.Name()

	case clients.Nimbus2.Name():
		consensusTag = ctx.String(flags.Nimbus2TagFlag)
		consensusCommit = ctx.String(flags.Nimbus2CommitHashFlag)
		validatorName = clients.PrysmValidator.Name()
	}

	switch cls[1] {
	case clients.Geth.Name():
		executionTag = ctx.String(flags.GethTagFlag)
		executionCommit = ctx.String(flags.GethCommitHashFlag)

	case clients.Erigon.Name():
		executionTag = ctx.String(flags.ErigonTagFlag)

	case clients.Nethermind.Name():
		executionTag = ctx.String(flags.NethermindTagFlag)
		executionCommit = ctx.String(flags.NethermindCommitHashFlag)

	case clients.Besu.Name():
		executionTag = ctx.String(flags.BesuTagFlag)
	}

	req := types.InstallRequest{
		Clients: []types.Client{
			{
				Name:    cls[0],
				Version: consensusTag,
				Commit:  consensusCommit,
			},
			{
				Name:    cls[1],
				Version: executionTag,
				Commit:  executionCommit,
			},
			{
				Name:    validatorName,
				Version: consensusTag,
				Commit:  consensusCommit,
			},
		},
	}

	r := c.handler.Install(req)
	if r.Error != nil {
		return utils.Exit(r.Error.Error(), 1)
	}

	c.log.Info("✅  Configuration files created!")
	c.log.Info("✅  Clients have been successfully installed.")
	c.log.Info("➡️  Start your node using 'lukso start'")

	return
}
