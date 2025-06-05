package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

type clientInstallHandler struct {
	hoverI             int
	visible            bool
	consensusOpts      []ClientOption
	executionOpts      []ClientOption
	selectingConsensus bool
	selectingExecution bool

	selectedConsensus string
	selectedExecution string

	ch chan any
}

func ClientInstallHandler() Handler {
	consOpts := []ClientOption{
		{
			"Prysm",
			clients.Prysm.Name(),
		},
		{
			"Lighthouse",
			clients.Lighthouse.Name(),
		},
		{
			"Teku",
			clients.Teku.Name(),
		},
		{
			"Nimbus (eth2)",
			clients.Nimbus2.Name(),
		},
	}

	execOpts := []ClientOption{
		{
			"Geth",
			clients.Geth.Name(),
		},
		{
			"Erigon",
			clients.Erigon.Name(),
		},
		{
			"Nethermind",
			clients.Nethermind.Name(),
		},
		{
			"Besu",
			clients.Besu.Name(),
		},
	}

	return &clientInstallHandler{
		hoverI:             0,
		consensusOpts:      consOpts,
		executionOpts:      execOpts,
		selectingConsensus: true,
		selectingExecution: false,
		visible:            false,
		ch:                 make(chan any),
	}
}

func (h *clientInstallHandler) Show() {
	h.visible = true
}

func (h *clientInstallHandler) Hide() {
	h.visible = false
}

func (h *clientInstallHandler) Visible() bool {
	return h.visible
}

func (h *clientInstallHandler) Handle(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			h.hoverI--
			if h.hoverI < 0 {
				h.hoverI = len(h.consensusOpts) - 1
			}

		case "down":
			h.hoverI++
			if h.hoverI >= len(h.consensusOpts) {
				h.hoverI = 0
			}

		case "enter":
			if h.selectingConsensus {
				h.selectingConsensus = false
				h.selectingExecution = true

				h.selectedConsensus = h.consensusOpts[h.hoverI].codeName
			} else if h.selectingExecution {
				h.selectedExecution = h.executionOpts[h.hoverI].codeName
				h.Send([]string{h.selectedConsensus, h.selectedExecution})
			}
		}
	}

	return nil
}

func (h *clientInstallHandler) View() (msg string) {
	var c lipgloss.TerminalColor = color.InactiveGrey
	var defC lipgloss.TerminalColor = color.InactiveGrey

	msg += "Select the client that you want to install:"
	msg += "\n\nConsensus:\n"

	if h.selectingConsensus {
		defC = lipgloss.NoColor{}
	} else {
		defC = color.InactiveGrey
	}

	for i, opt := range h.consensusOpts {
		c = defC
		ind := "- "
		if i == h.hoverI && h.selectingConsensus {
			ind = "> "
			c = color.HoverGreen
		}

		msg += lipgloss.NewStyle().Foreground(c).Render(ind + opt.name + "\n")
	}

	if h.selectingExecution {
		defC = lipgloss.NoColor{}
	} else {
		defC = color.InactiveGrey
	}

	msg += "\n\nExecution\n"
	for i, opt := range h.executionOpts {
		c = defC
		ind := "- "
		if i == h.hoverI && h.selectingExecution {
			ind = "> "
			c = color.HoverGreen
		}

		msg += lipgloss.NewStyle().Foreground(c).Render(ind + opt.name + "\n")
	}

	return
}

func (h *clientInstallHandler) Send(resp any) {
	h.ch <- resp
	h.visible = false
}

func (h *clientInstallHandler) Get() any {
	return <-h.ch
}
