package apitypes

type JsonRpcRequest struct {
	JsonRPC string
	ID      int
	Method  string
	Params  []string
}
