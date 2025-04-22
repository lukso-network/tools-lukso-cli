package types

type InitRequest struct {
	Directory string
	Reinit    bool
	Ip        string
}

type (
	InstallRequest         struct{}
	UpdateRequest          struct{}
	StartRequest           struct{}
	StopRequest            struct{}
	StatusRequest          struct{}
	LogsRequest            struct{}
	ResetRequest           struct{}
	ValidatorImportRequest struct{}
	ValidatorListRequest   struct{}
	ValidatorExitRequest   struct{}
	VersionRequest         struct{}
	VersionClientsRequest  struct{}
)

type Request interface {
	InitRequest |
		InstallRequest |
		UpdateRequest |
		StartRequest |
		StopRequest |
		StatusRequest |
		LogsRequest |
		ResetRequest |
		ValidatorImportRequest |
		ValidatorListRequest |
		ValidatorExitRequest |
		VersionRequest |
		VersionClientsRequest
}
