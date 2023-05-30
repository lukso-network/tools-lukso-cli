package clients

import "github.com/urfave/cli/v2"

type ClientBinaryDependency interface {
	Start(*cli.Context) error
	Stop() error
	Status()
	Logs()
	Reset()
	Install()
	Update()

	IsRunning() bool
	ParseUserFlags()
	PrepareStartFlags()
	// Name is a user-readable name utility, f.e. in logs etc.
	// Should be uppercase and match CommandName (non-case-sensitively)
	Name() string
	// CommandName identifies client in all sorts of technical aspects - commands, files etc.
	// Should be lowercase and match Name (non-case-sensitively)
	CommandName() string
	// ParseUrl replaces any missing information in install link with matching system info
	ParseUrl()
}
