package page

import tea "github.com/charmbracelet/bubbletea"

type WelcomePage struct{}

var _ Page = &WelcomePage{}

func (p *WelcomePage) HandleInput(msg tea.KeyMsg) {
}

func (p *WelcomePage) View() string {
	return "Hello world!"
}
