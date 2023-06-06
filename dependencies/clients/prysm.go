package clients

type PrysmClient struct {
	*clientBinary
}

func NewPrysmClient() *PrysmClient {
	return &PrysmClient{
		&clientBinary{
			name:           prysmDependencyName,
			commandName:    "prysm",
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var Prysm = NewPrysmClient()

var _ ClientBinaryDependency = &PrysmClient{}
