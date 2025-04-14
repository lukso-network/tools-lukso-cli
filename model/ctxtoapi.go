package model

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common/progress"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func CtxToApiInit(ctx *cli.Context) (api types.InitRequest) {
	return types.InitRequest{
		Reinit:   ctx.Bool(flags.ReinitFlag),
		Ip:       ctx.String(flags.IpFlag),
		Progress: progress.NewProgress(),
	}
}
