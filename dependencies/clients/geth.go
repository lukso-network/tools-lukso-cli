package clients

type GethClient struct {
	*clientBinary
}

func NewGethClient() *GethClient {
	return &GethClient{
		&clientBinary{
			name:           gethDependencyName,
			commandName:    "geth",
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var Geth = NewGethClient()

var _ ClientBinaryDependency = &GethClient{}
