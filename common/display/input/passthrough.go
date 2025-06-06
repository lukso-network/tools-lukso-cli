package input

import tea "github.com/charmbracelet/bubbletea"

// passthroughHandler serves as a default input handler that does nothing but is not a nil.
type passthroughHandler struct{}

var _ Handler = &passthroughHandler{}

func PassthroughHandler() Handler {
	return &passthroughHandler{}
}

func (h *passthroughHandler) Handle(msg tea.Msg) tea.Cmd { return nil }
func (h *passthroughHandler) Show()                      {}
func (h *passthroughHandler) Hide()                      {}
func (h *passthroughHandler) Visible() bool              { return false }
func (h *passthroughHandler) View() string               { return "" }
func (h *passthroughHandler) Send(resp any)              { return }
func (h *passthroughHandler) Get() any                   { return nil }
