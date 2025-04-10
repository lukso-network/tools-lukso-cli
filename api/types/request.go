package types

type InitArgs struct {
	Directory string
	Reinit    bool
	Ip        string
}

type (
	InstallArgs         struct{}
	UpdateArgs          struct{}
	StartArgs           struct{}
	StopArgs            struct{}
	StatusArgs          struct{}
	LogsArgs            struct{}
	ResetArgs           struct{}
	ValidatorImportArgs struct{}
	ValidatorListArgs   struct{}
	ValidatorExitArgs   struct{}
	VersionArgs         struct{}
	VersionClientsArgs  struct{}
)
