package types

type InitResponse struct {
	Error error
}

type (
	InstallResponse         struct{}
	UpdateResponse          struct{}
	StartResponse           struct{}
	StopResponse            struct{}
	StatusResponse          struct{}
	LogsResponse            struct{}
	ResetResponse           struct{}
	ValidatorImportResponse struct{}
	ValidatorListResponse   struct{}
	ValidatorExitResponse   struct{}
	VersionResponse         struct{}
	VersionClientsResponse  struct{}
)

type Response interface {
	InitResponse |
		InstallResponse |
		UpdateResponse |
		StartResponse |
		StopResponse |
		StatusResponse |
		LogsResponse |
		ResetResponse |
		ValidatorImportResponse |
		ValidatorListResponse |
		ValidatorExitResponse |
		VersionResponse |
		VersionClientsResponse
}
