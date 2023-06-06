package clients

type LighthouseClient struct {
	*clientBinary
}

func NewLighthouseClient() *LighthouseClient {
	return &LighthouseClient{
		&clientBinary{
			name:           lighthouseDependencyName,
			commandName:    "lighthouse",
			baseUrl:        "",
			githubLocation: "",
		},
	}
}

var Lighthouse = NewLighthouseClient()

var _ ClientBinaryDependency = &LighthouseClient{}
