package clients

type GethClient struct {
	*clientBinary
}

func NewGethClient() *GethClient {
	return &GethClient{}
}

var Geth = NewGethClient()

var _ ClientBinaryDependency = &GethClient{}
