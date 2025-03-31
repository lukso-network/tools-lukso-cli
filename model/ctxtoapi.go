package model

import (
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/api/types"
)

func CtxToApiInit(ctx *cli.Context) (api types.InitArgs) {
	return types.InitArgs{}
}
