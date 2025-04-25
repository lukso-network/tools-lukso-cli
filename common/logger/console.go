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

func (c ConsoleLogger) Debugf(msg string, args ...string) {
	log.Debugf(msg, args)
}

func (c ConsoleLogger) Info(msg string) {
	log.Info(msg)
}

func (c ConsoleLogger) Infof(msg string, args ...string) {
	log.Infof(msg, args)
}

func (c ConsoleLogger) Warn(msg string) {
	log.Warn(msg)
}

func (c ConsoleLogger) Warnf(msg string, args ...string) {
	log.Warnf(msg, args)
}

func (c ConsoleLogger) Error(msg string) {
	log.Error(msg)
}

func (c ConsoleLogger) Errorf(msg string, args ...string) {
	log.Errorf(msg, args)
}

// Not meant for std console logging.
func (c ConsoleLogger) Clear() {}
func (c ConsoleLogger) Close() {}

func (c ConsoleLogger) Progress() progress.Progress {
	return progress.NewStubProgress()
}
