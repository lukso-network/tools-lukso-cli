package page

import tea "github.com/charmbracelet/bubbletea"

type InstallPage struct{}

var _ Page = &InstallPage{}

func (p *InstallPage) Title() string {
	return "Install"
}

func (p *InstallPage) HandleInput(msg tea.KeyMsg) {
}

func (p *InstallPage) View() string {
	return "Hello world! From the Install page"
}
