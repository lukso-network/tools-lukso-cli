package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/common/display/input"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/common/progress"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
)

// cmdDisplay is a bubbletea program wrapper.
// It is able to register an input handler to make the CLI more interactive.
type cmdDisplay struct {
	msgCh        <-chan tea.Msg
	inputCh      chan any
	inputHandler input.Handler
	logs         []logger.LogMsg
	prg          progress.Progress
}

func NewCmdDisplay(ch chan tea.Msg, inputCh chan any) Display {
	return &cmdDisplay{
		msgCh:        ch,
		inputCh:      inputCh,
		inputHandler: input.PassthroughHandler(),
		logs:         make([]logger.LogMsg, 0),
	}
}

func (d *cmdDisplay) AddInputHandler(h input.Handler) {
	d.inputHandler = h
	d.inputHandler.Show()
}

func (d *cmdDisplay) InputHandler() input.Handler {
	return d.inputHandler
}

func (d *cmdDisplay) Init() tea.Cmd {
	return nil
}

func (d *cmdDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logger.LogMsg:
		if msg.IsClear {
			d.logs = make([]logger.LogMsg, 0)
		}
		if msg.Msg != "" {
			d.logs = insertSeq(msg, d.logs)
		}

		d.prg = msg.Progress

	case tea.KeyMsg:
		return d.handleInput(msg)

	case tea.QuitMsg:
		return d, tea.Interrupt

		// In some cases, we want to execute a Cmd that is sent on the channel
	case tea.Cmd:
		return d, msg
	}

	return d, d.listenForMsg()
}

func (d *cmdDisplay) View() string {
	return d.Render()
}

func (d *cmdDisplay) Listen() {
	p := tea.NewProgram(d)
	_, err := p.Run() // TODO: handle all possible terminal input SIGs
	if err != nil {
		utils.Exit("Command cancelled", 1)
	}
}

func (d *cmdDisplay) Render() (out string) {
	if d.inputHandler.Visible() {
		out = d.inputHandler.View()
	} else {
		for _, log := range d.logs {
			out += log.Msg + "\n"
		}
	}

	if d.prg == nil || !d.prg.Visible() {
		return
	}

	out = lipgloss.JoinVertical(lipgloss.Left, out, d.prg.Render())

	return
}

func (d *cmdDisplay) handleInput(msg tea.KeyMsg) (m tea.Model, cmd tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return d, tea.Interrupt

	default:
		return d, tea.Batch(d.inputHandler.Handle(msg), d.listenForMsg())
	}
}

func (d *cmdDisplay) listenForMsg() tea.Cmd {
	return ChanCmd(d.msgCh)
}

func ChanCmd(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

// insertSeq ensures the sequential order for displayed logs.
func insertSeq(log logger.LogMsg, logs []logger.LogMsg) []logger.LogMsg {
	for i, presLog := range logs {
		if log.Seq < presLog.Seq {
			return utils.InsertAt(log, i, logs)
		}
	}

	return append(logs, log)
}
