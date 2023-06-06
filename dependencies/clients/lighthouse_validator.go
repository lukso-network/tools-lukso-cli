package clients

type LighthouseValidatorClient struct {
	*clientBinary
}

func NewLighthouseValidatorClient() *LighthouseClient {
	return &LighthouseClient{
		&clientBinary{
			name:           validatorDependencyName,
			commandName:    "lighthouse", // we run it using lighthouse bin
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var LighthouseValidator = NewLighthouseValidatorClient()

var _ ValidatorBinaryDependency = &LighthouseValidatorClient{}

func (l LighthouseValidatorClient) Import() {
	//TODO implement me
	panic("implement me")
}

func (l LighthouseValidatorClient) List() {
	//TODO implement me
	panic("implement me")
}
