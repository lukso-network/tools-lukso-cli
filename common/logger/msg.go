package logger

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lukso-network/tools-lukso-cli/common/progress"
)

type msgLogger struct {
	ch  chan<- tea.Msg
	prg progress.Progress
}

type LogMsg struct {
	Msg      string
	IsClear  bool
	Progress progress.Progress
}

var _ Logger = &msgLogger{}

func Log(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func NewMsgLogger(ch chan tea.Msg, prg progress.Progress) Logger {
	return &msgLogger{
		ch:  ch,
		prg: prg,
	}
}

func (l *msgLogger) Debug(msg string) {
	l.sendLeveledMsg(msg, LevelDebug)
}

func (l *msgLogger) Info(msg string) {
	l.sendLeveledMsg(msg, LevelInfo)
}

func (l *msgLogger) Warn(msg string) {
	l.sendLeveledMsg(msg, LevelWarn)
}

func (l *msgLogger) Error(msg string) {
	l.sendLeveledMsg(msg, LevelError)
}

func (l *msgLogger) Clear() {
	l.ch <- LogMsg{IsClear: true, Progress: l.prg}
}

func (l *msgLogger) Close() {
	l.ch <- tea.QuitMsg{}
	close(l.ch)
}

func (l *msgLogger) Progress() progress.Progress {
	l.ch <- LogMsg{Progress: l.prg}
	return l.prg
}

func (l *msgLogger) sendLeveledMsg(msg string, lvl logLvl) {
	msg = lvlFormatters[lvl](msg)
	l.ch <- LogMsg{
		Msg:      msg,
		IsClear:  false,
		Progress: l.prg,
	}
}
