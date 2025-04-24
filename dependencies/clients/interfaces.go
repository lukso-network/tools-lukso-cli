package clients

import "github.com/urfave/cli/v2"

type Client interface {
	// Install installs the client with given version
	Install(version string, isUpdate bool) error

	// Update updates client to specific version
	Update() error

	// Start starts the client with given flags
	Start(arguments []string, logDir string) error

	// Stop stops the client
	Stop() error

	// Logs takes the latest log file and prints it to terminal, in live mode
	Logs(logsDirPath string) error

	// Reset deletes data directories of all clients
	Reset(datadir string) error

	// IsRunning determines whether the client is already running
	IsRunning() bool

	// ParseUrl replaces any missing information in install link with matching system info
	ParseUrl(tag, commitHash string) string

	// ParseUserFlags is used to trim any client prefix from flag
	ParseUserFlags(ctx *cli.Context) []string

	// PrepareStartFlags parses arguments that are later supplied to Start
	PrepareStartFlags(ctx *cli.Context) ([]string, error)

	// Name is a user-readable name utility, f.e. in logs etc.
	// Should be uppercase and match CommandName (non-case-sensitively)
	Name() string

	// CommandName identifies client in all sorts of technical aspects - commands, files etc.
	// Should be lowercase and match Name (non-case-sensitively)
	FileName() string

	// FileDir returns a path to client's installation directory.
	FileDir() string

	// FilePath returns path to installed binary
	FilePath() string

	// Peers prints to console how many peers does the client have
	Peers(ctx *cli.Context) (outbound int, inbound int, err error)

	// Version returns a version of the given client as a string (different clients may vary in versioning).
	// For compatibility with display, all clients should be preceded by 'v', e.g. Geth version: v1.14.11
	Version() string

	// Init initializes client, using internal implementations.
	Init() error

	// Each client should be able to identify what release it should install, based on the build envs.
	tag() string
	os() string
	arch() string
	commit() string
}

type ValidatorBinaryDependency interface {
	Client
	Import(ctx *cli.Context) error
	List(ctx *cli.Context) error
	Exit(ctx *cli.Context) error
}
