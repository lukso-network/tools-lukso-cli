package clients

import (
	"github.com/urfave/cli/v2"
	"strings"
)

type GethClient struct {
	clientBinary
}

func NewGethClient() *GethClient {
	return &GethClient{}
}

var Geth = NewGethClient()

var _ ClientBinaryDependency = &GethClient{}

func (g *GethClient) Start(ctx *cli.Context) (err error) {
	return start(ctx, g, []string{})
}

func (g *GethClient) Stop() error {
	return stop(g)
}

func (g *GethClient) Status() {
	status(g)
}

func (g *GethClient) Logs() {
	logs(g)
}

func (g *GethClient) Reset() {
	reset(g)
}

func (g *GethClient) Install() {
	install(g)
}

func (g *GethClient) Update() {
	update(g)
}

func (g *GethClient) IsRunning() bool {
	return isRunning(g)
}

func (g *GethClient) ParseUserFlags() {
}

func (g *GethClient) PrepareStartFlags() {
}

func (g *GethClient) Name() string {
	return g.name
}

func (g *GethClient) CommandName() string {
	return strings.ToLower(g.name)
}

func (g *GethClient) ParseUrl() {
}
