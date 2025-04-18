// package logger is used to transfer logs from API handler to a specific Stdout - whether a CMD or a UI component.
// It is not used as a standard logger - it serves more as a "postman" for messages between API and the app output.
package logger

import "github.com/lukso-network/tools-lukso-cli/common/progress"

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	// Clear clears ther log msg buffer.
	Clear()
	// Close indicates that if there are blocking message listeners, they should stop listening and unblock.
	Close()

	// Progress returns the underlying Progress tracker.
	// If Progress is nil, then a new stub progress will be created to avoid nil references.
	Progress() progress.Progress
}
