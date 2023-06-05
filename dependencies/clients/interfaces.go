package clients

import "github.com/urfave/cli/v2"

type ClientBinaryDependency interface {
	// Start starts the client with given flags
	Start(ctx *cli.Context, arguments []string) error

	// Stop stops the client
	Stop() error

	// Logs takes the latest log file and prints it to terminal, in live mode
	Logs(logsDirPath string) error

	// Reset deletes data directories of all clients
	Reset(datadir string) error

	// Install installs the client with given version
	Install(tag, commitHash string) error

	// Update updates client to specific version - TODO
	Update()

	// IsRunning determines whether the client is already running
	IsRunning() bool

	// ParseUserFlags is used to trim any client prefix from flag
	ParseUserFlags(ctx *cli.Context) []string

	// PrepareStartFlags parses arguments that are later supplied to Start
	PrepareStartFlags() []string

	// Name is a user-readable name utility, f.e. in logs etc.
	// Should be uppercase and match CommandName (non-case-sensitively)
	Name() string

	// CommandName identifies client in all sorts of technical aspects - commands, files etc.
	// Should be lowercase and match Name (non-case-sensitively)
	CommandName() string

	// ParseUrl replaces any missing information in install link with matching system info
	ParseUrl(tag, commitHash string) string

	// FilePath returns path to installed binary
	FilePath() string
}
