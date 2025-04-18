// package logger is used to transfer logs from API handler to a specific Stdout - whether a CMD or a UI component.
// It is not used as a standard logger - it serves more as a "postman" for messages between API and the app output.
package logger

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

	Clear()
	// CLose indicates that if there are blocking message listeners, they should stop listening and unblock
	Close()
}
