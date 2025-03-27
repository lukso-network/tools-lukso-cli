package types

type Response interface {
	Error() string
}

type InitResponse struct {
	ErrorResponse
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
