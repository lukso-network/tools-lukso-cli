package dep

import "github.com/urfave/cli/v2"

// UrlBuilder is responsible for parsing a URL for multiple OS options.
type UrlBuilder interface {
	ParseUrl(tag, commitHash string) string
	Tag() string
	Commit() string
	Os() string
	Arch() string
}

// Installer represents a dependency that can be installed and updated.
type Installer interface {
	UrlBuilder
	Install(version string, isUpdate bool) error
	Update() error
}

// FileIdentifier ensures that the installed dependency can be located in the file system.
type FileIdentifier interface {
	FileName() string
	FileDir() string
	FilePath() string
}

// Client is a generic interface for a client dependency. For a specific layer,
// use ExecutionClient, ConsensusClient or ValidatorClient
type Client interface {
	Installer
	FileIdentifier
	Start(ctx *cli.Context, arguments []string) error
	Stop() error
	Logs(logsDirPath string) error
	Reset(datadir string) error
	IsRunning() bool
	ParseUserFlags(ctx *cli.Context) []string
	PrepareStartFlags(ctx *cli.Context) ([]string, error)
	Version() string
	Name() string
}

type ExecutionClient interface {
	Client
	Peers(ctx *cli.Context) (outbound int, inbound int, err error)
	Init() error
}

type ConsensusClient interface {
	Client
	Peers(ctx *cli.Context) (outbound int, inbound int, err error)
}

type ValidatorClient interface {
	Client
	Import(ctx *cli.Context) error
	List(ctx *cli.Context) error
	Exit(ctx *cli.Context) error
}
