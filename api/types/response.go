package types

type Response interface {
	// This is not an errorer: it's just a way of unwraping an actual error from the API response
	Error() error
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
