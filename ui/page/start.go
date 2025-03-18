package page

import tea "github.com/charmbracelet/bubbletea"

type StartPage struct{}

var _ Page = &StartPage{}

func (p *StartPage) Title() string {
	return "Start/Stop"
}

func (p *StartPage) HandleInput(msg tea.KeyMsg) {
}

func (p *StartPage) View() string {
	return "Hello world! From the start page"
}
