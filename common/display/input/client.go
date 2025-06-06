package input

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lukso-network/tools-lukso-cli/dep/clients"
)

type clientInstallHandler struct {
	visible bool
	sects   *SectionGroup

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

	sects := NewSectionGroup(
		"Please select the clients that you want to install on your node:",
		[]Section{
			NewSection("Consensus", consOpts),
			NewSection("Execution", execOpts),
		},
	)

	return &clientInstallHandler{
		sects:   sects,
		visible: false,
		ch:      make(chan any),
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
			h.sects.PrevOption()

		case "down":
			h.sects.NextOption()

		case "enter":
			res := h.sects.NextSection()
			if res != nil {
				h.Send(res)
			}
		}
	}

	return nil
}

func (h *clientInstallHandler) View() (msg string) {
	return h.sects.View()
}

func (h *clientInstallHandler) Send(resp any) {
	h.ch <- resp
	h.visible = false
}

func (h *clientInstallHandler) Get() any {
	return <-h.ch
}
