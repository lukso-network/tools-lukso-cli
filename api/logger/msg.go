package logger

import (
	tea "github.com/charmbracelet/bubbletea"
)

type msgLogger struct {
	ch chan<- tea.Msg
}

type LogMsg struct {
	Msg     string
	IsClear bool
}

var _ Logger = &msgLogger{}

func Log(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func NewMsgLogger(ch chan tea.Msg) Logger {
	return &msgLogger{
		ch: ch,
	}
}

func (l msgLogger) Debug(msg string) {
	l.ch <- LogMsg{Msg: msg, IsClear: false}
}

func (l msgLogger) Info(msg string) {
	l.ch <- LogMsg{Msg: msg, IsClear: false}
}

func (l msgLogger) Warn(msg string) {
	l.ch <- LogMsg{Msg: msg, IsClear: false}
}

func (l msgLogger) Error(msg string) {
	l.ch <- LogMsg{Msg: msg, IsClear: false}
}

func (l msgLogger) Clear() {
	l.ch <- LogMsg{IsClear: true}
}

func (l msgLogger) Close() {
	l.ch <- tea.QuitMsg{}
	close(l.ch)
}
