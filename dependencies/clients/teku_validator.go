package clients

import "github.com/urfave/cli/v2"

type TekuValidatorClient struct {
	*clientBinary
}

func NewTekuValidatorClient() *TekuValidatorClient {
	return &TekuValidatorClient{
		&clientBinary{
			name:           tekuValidatorDependencyName,
			commandName:    "validator_tk",
			baseUrl:        "",
			githubLocation: tekuGithubLocation,
		},
	}
}

var TekuValidator = NewTekuValidatorClient()

var _ ValidatorBinaryDependency = &TekuValidatorClient{}

func (t TekuValidatorClient) Import(ctx *cli.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t TekuValidatorClient) List(ctx *cli.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t TekuValidatorClient) Exit(ctx *cli.Context) error {
	//TODO implement me
	panic("implement me")
}
