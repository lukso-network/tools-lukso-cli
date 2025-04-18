package logger

import (
	log "github.com/sirupsen/logrus"

	"github.com/lukso-network/tools-lukso-cli/common/progress"
)

var _ Logger = ConsoleLogger{}

type ConsoleLogger struct{}

func (c ConsoleLogger) Debug(msg string) {
	log.Debug(msg)
}

func (c ConsoleLogger) Info(msg string) {
	log.Info(msg)
}

func (c ConsoleLogger) Warn(msg string) {
	log.Warn(msg)
}

func (c ConsoleLogger) Error(msg string) {
	log.Error(msg)
}

// Not meant for std console logging.
func (c ConsoleLogger) Clear() {}
func (c ConsoleLogger) Close() {}

func (c ConsoleLogger) Progress() progress.Progress {
	return progress.NewStubProgress()
}
