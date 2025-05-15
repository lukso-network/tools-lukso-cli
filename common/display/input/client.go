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
	selectedClients    []string

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
		hoverI:        0,
		consensusOpts: consOpts,
		executionOpts: execOpts,
		visible:       false,
		ch:            make(chan any),
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
}

func (h *clientInstallHandler) Get() any {
	return <-h.ch
}
