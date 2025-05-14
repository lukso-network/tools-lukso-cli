package apitypes

type GithubApiReleaseResponse struct {
	TagName string `json:"tag_name"`
}

type GithubApiCommitResponse struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}

type CheckpointFinalizedBlockResponse struct {
	Data struct {
		Message struct {
			Slot string `json:"slot"`
		} `json:"message"`
	} `json:"data"`
}

type CheckpointFinalizedBlockRootResponse struct {
	Data struct {
		Root string `json:"root"`
	} `json:"data"`
}

type ExplorerFinalizedSlotsResponse struct {
	Data []struct {
		Epoch     int    `json:"epoch"`
		BlockRoot string `json:"blockroot"`
	} `json:"data"`
}

type PeersJsonRpcResponse struct {
	Result []struct {
		Network struct {
			Inbound bool `json:"inbound"`
		} `json:"network"`
	} `json:"result"`
	Error interface{} `json:"error,omitempty"`
}

type PeersBeaconAPIResponse struct {
	Data []struct {
		State     string `json:"state"`
		Direction string `json:"direction"`
	} `json:"data"`
}
