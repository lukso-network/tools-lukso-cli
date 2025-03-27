package logger

import log "github.com/sirupsen/logrus"

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

func (c ConsoleLogger) IsConsole() bool {
	return true
}
