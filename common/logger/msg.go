package logger

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lukso-network/tools-lukso-cli/common/progress"
)

type msgLogger struct {
	ch  chan<- tea.Msg
	prg progress.Progress
	seq int
}

type LogMsg struct {
	Msg      string
	IsClear  bool
	Progress progress.Progress
	Seq      int
}

var _ Logger = &msgLogger{}

func NewMsgLogger(ch chan tea.Msg, prg progress.Progress) Logger {
	return &msgLogger{
		ch:  ch,
		prg: prg,
		seq: 0,
	}
}

func (l *msgLogger) Debug(msg string) {
	l.sendLeveledMsg(msg, LevelDebug)
}

func (l *msgLogger) Debugf(msg string, args ...any) {
	l.sendLeveledMsg(fmt.Sprintf(msg, args...), LevelDebug)
}

func (l *msgLogger) Info(msg string) {
	l.sendLeveledMsg(msg, LevelInfo)
}

func (l *msgLogger) Infof(msg string, args ...any) {
	l.sendLeveledMsg(fmt.Sprintf(msg, args...), LevelInfo)
}

func (l *msgLogger) Warn(msg string) {
	l.sendLeveledMsg(msg, LevelWarn)
}

func (l *msgLogger) Warnf(msg string, args ...any) {
	l.sendLeveledMsg(fmt.Sprintf(msg, args...), LevelWarn)
}

func (l *msgLogger) Error(msg string) {
	l.sendLeveledMsg(msg, LevelError)
}

func (l *msgLogger) Errorf(msg string, args ...any) {
	l.sendLeveledMsg(fmt.Sprintf(msg, args...), LevelError)
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
	l.seq++
	return l.prg
}

func (l *msgLogger) sendLeveledMsg(msg string, lvl logLvl) {
	msg = lvlFormatters[lvl](msg)
	l.ch <- LogMsg{
		Msg:      msg,
		IsClear:  false,
		Progress: l.prg,
		Seq:      l.seq,
	}
	l.seq++
}
