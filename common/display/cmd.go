package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/progress"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
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
	return nil
}

func (d *cmdDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logger.LogMsg:
		if msg.IsClear {
			d.logs = make([]string, 0)
		}
		if msg.Msg != "" {
			d.logs = append(d.logs, msg.Msg)
		}

		d.prg = msg.Progress

	case tea.QuitMsg:
		return d, tea.Interrupt

		// In some cases, we want to execute a Cmd that is sent on the channel
	case tea.Cmd:
		return d, tea.Batch(msg, d.listenForMsg())
	}

	return d, d.listenForMsg()
}

func (d *cmdDisplay) View() string {
	return d.Render()
}

func (d *cmdDisplay) Listen() {
	p := tea.NewProgram(d, tea.WithInput(nil))
	_, err := p.Run() // TODO: handle all possible terminal input SIGs
	if err != nil {
		utils.Exit("Command cancelled", 1)
	}
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
	return ChanCmd(d.ch)
}

func ChanCmd(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
