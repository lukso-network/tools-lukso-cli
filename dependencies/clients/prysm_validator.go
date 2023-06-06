package clients

type PrysmValidatorClient struct {
	*clientBinary
}

func NewPrysmValidatorClient() *PrysmClient {
	return &PrysmClient{
		&clientBinary{
			name:           validatorDependencyName,
			commandName:    "validator",
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var PrysmValidator = NewPrysmValidatorClient()

var _ ValidatorBinaryDependency = &PrysmValidatorClient{}

func (p *PrysmValidatorClient) Import() {
	//TODO implement me
	panic("implement me")
}

func (p *PrysmValidatorClient) List() {
	//TODO implement me
	panic("implement me")
}
