package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/progress"
)

// cmdDisplay is a bubbletea program wrapper.
type cmdDisplay struct {
	ch   <-chan tea.Msg
	logs []string
	prg  progress.Progress
}

func NewCmdDisplay(ch chan tea.Msg) Display {
	return &cmdDisplay{
		ch:   ch,
		logs: make([]string, 0),
	}
}

func (d *cmdDisplay) Init() tea.Cmd {
	return tea.ShowCursor
}

func (d *cmdDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logger.LogMsg:
		if msg.IsClear {
			d.logs = make([]string, 0)
		} else {
			d.logs = append(d.logs, msg.Msg)
		}

		d.prg = msg.Progress

	case tea.QuitMsg:
		return d, tea.Quit
	}

	return d, d.listenForMsg()
}

func (d *cmdDisplay) View() string {
	return d.Render()
}

func (d *cmdDisplay) Listen() {
	p := tea.NewProgram(d, tea.WithInput(nil))
	p.Run() // TODO: handle all possible terminal input SIGs
}

func (d *cmdDisplay) Render() (out string) {
	for _, log := range d.logs {
		out += log + "\n"
	}

	if d.prg == nil || !d.prg.Visible() {
		return
	}

	out = lipgloss.JoinVertical(lipgloss.Left, out, d.prg.Render())

	return
}

func (d *cmdDisplay) listenForMsg() tea.Cmd {
	return logger.Log(d.ch)
}
