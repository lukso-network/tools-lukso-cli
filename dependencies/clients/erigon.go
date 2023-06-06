package clients

type ErigonClient struct {
	*clientBinary
}

func NewErigonClient() *ErigonClient {
	return &ErigonClient{
		&clientBinary{
			name:           erigonDependencyName,
			commandName:    "erigon",
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var Erigon = NewErigonClient()

var _ ClientBinaryDependency = &ErigonClient{}
