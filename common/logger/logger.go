// package logger is used to transfer logs from API handler to a specific Stdout - whether a CMD or a UI component.
// It is not used as a standard logger - it serves more as a "postman" for messages between API and the app output.
package logger

import (
	"fmt"

	"github.com/lukso-network/tools-lukso-cli/common/progress"
)

const (
	LevelDebug logLvl = iota
	LevelInfo
	LevelWarn
	LevelError
)

var lvlFormatters = map[logLvl]msgFormat{
	LevelDebug: formatDebug,
	LevelInfo:  formatInfo,
	LevelWarn:  formatWarn,
	LevelError: formatError,
}

var lvlId = map[logLvl]string{
	LevelDebug: "DBG",
	LevelInfo:  "INF",
	LevelWarn:  "WRN",
	LevelError: "ERR",
}

var lvlColor = map[logLvl]int{
	LevelDebug: 35,
	LevelInfo:  36,
	LevelWarn:  33,
	LevelError: 31,
}

type logLvl int

type msgFormat func(msg string) string

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

func formatLog(msg string, lvl logLvl) string {
	return fmt.Sprintf("[\x1b[%dm%s\x1b[0m] %s", lvlColor[lvl], lvlId[lvl], msg)
}

func formatDebug(msg string) string {
	return formatLog(msg, LevelDebug)
}

func formatInfo(msg string) string {
	return formatLog(msg, LevelInfo)
}

func formatWarn(msg string) string {
	return formatLog(msg, LevelWarn)
}

func formatError(msg string) string {
	return formatLog(msg, LevelError)
}
