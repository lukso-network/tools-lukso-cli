package input

import tea "github.com/charmbracelet/bubbletea"

type Handler interface {
	Handle(msg tea.Msg) tea.Cmd

	Show()
	Hide()
	Visible() bool

	View() string

	// Get should block the caller until the handler is able to Send the response of the input actions.
	Get() (resp any)

	// Send sends a response of a handler, that you can get with Get.
	Send(resp any)
}
