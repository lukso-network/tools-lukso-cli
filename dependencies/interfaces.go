package dependencies

import "github.com/urfave/cli/v2"

type UrlBuilder interface {
	ParseUrl(tag, commitHash string) string
	tag() string
	commit() string
	os() string
	arch() string
}

type Installer interface {
	UrlBuilder
	Install(version string, isUpdate bool) error
	Update() error
}

type FileIdentifier interface {
	FileName() string
	FileDir() string
	FilePath() string
}

type Client interface {
	Installer
	Start(ctx *cli.Context, arguments []string) error
	Stop() error
	Logs(logsDirPath string) error
	Reset(datadir string) error
	IsRunning() bool
	ParseUserFlags(ctx *cli.Context) []string
	PrepareStartFlags(ctx *cli.Context) ([]string, error)
	Version() string
}

type ExecutionClient interface {
	Peers(ctx *cli.Context) (outbound int, inbound int, err error)
	Init() error
}

type ConsensusClient interface {
	Peers(ctx *cli.Context) (outbound int, inbound int, err error)
}

type ValidatorClient interface {
	Client
	Import(ctx *cli.Context) error
	List(ctx *cli.Context) error
	Exit(ctx *cli.Context) error
}
