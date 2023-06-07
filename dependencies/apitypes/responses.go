package apitypes

type GithubApiReleaseResponse struct {
	TagName string `json:"tag_name"`
}

type GithubApiCommitResponse struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}
