package input

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lukso-network/tools-lukso-cli/dep/clients"
)

type clientOption struct {
	name     string
	codeName string
}

type clientInstallHandler struct {
	hoverI             int
	visible            bool
	consensusOpts      []clientOption
	executionOpts      []clientOption
	selectingConsensus bool
	selectingExecution bool

	selectedConsensus string
	selectedExecution string

	ch chan any
}

func ClientInstallHandler() Handler {
	consOpts := []clientOption{
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

	execOpts := []clientOption{
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
			}

			if h.selectingExecution {
				h.selectedExecution = h.executionOpts[h.hoverI].codeName
				h.Send([]string{h.selectedConsensus, h.selectedExecution})
			}
		}
	}

	return nil
}

func (h *clientInstallHandler) View() (msg string) {
	msg += "Select the Consensus client that you want to install:\n"

	for i, opt := range h.consensusOpts {
		ind := "- "
		if i == h.hoverI && h.selectingConsensus {
			ind = "> "
		}

		msg += ind + opt.name + "\n"
	}

	for i, opt := range h.executionOpts {
		ind := "- "
		if i == h.hoverI && h.selectingExecution {
			ind = "> "
		}

		msg += ind + opt.name + "\n"
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
