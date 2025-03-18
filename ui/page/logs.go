package page

import tea "github.com/charmbracelet/bubbletea"

type LogsPage struct{}

var _ Page = &LogsPage{}

func (p *LogsPage) Title() string {
	return "Logs"
}

func (p *LogsPage) HandleInput(msg tea.KeyMsg) {
}

func (p *LogsPage) View() string {
	return "Hello world! From the Logs page"
}
