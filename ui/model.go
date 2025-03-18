package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/binding"
	"github.com/lukso-network/tools-lukso-cli/ui/common/context"
	"github.com/lukso-network/tools-lukso-cli/ui/navbar"
	"github.com/lukso-network/tools-lukso-cli/ui/page"
	"github.com/lukso-network/tools-lukso-cli/ui/style"
)

type Model struct {
	pages       []page.Page
	activePageI int
	titles      []string

	keyPress string
	counter  int
}

var _ tea.Model = &Model{}

func NewModel() *Model {
	pages := []page.Page{
		&page.WelcomePage{},
		&page.InstallPage{},
		&page.StartPage{},
		&page.LogsPage{},
		&page.TestPage{},
	}

	titles := make([]string, len(pages))
	for i, page := range pages {
		titles[i] = page.Title()
	}

	return &Model{
		pages:       pages,
		activePageI: 0,
		titles:      titles,
		counter:     0,
	}
}

func (m *Model) Init() (cmd tea.Cmd) {
	return nil
}

func (m *Model) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) {
	ctx := context.Context{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyPress := msg.String()
		m.keyPress = keyPress

		switch {
		case binding.Matches(binding.Mapping.Quit, keyPress):
			return m, tea.Quit

		case binding.Matches(binding.Mapping.PreviousPage, keyPress):
			m.activePageI--
			if m.activePageI < 0 {
				m.activePageI = len(m.pages) - 1
			}

		case binding.Matches(binding.Mapping.NextPage, keyPress):
			m.activePageI++
			if m.activePageI > len(m.pages)-1 {
				m.activePageI = 0
			}

		default:
		}
	}

	ctx.Msg = msg

	m.pages[m.activePageI].Handle(ctx)

	return m, nil
}

func (m *Model) View() string {
	content := lipgloss.JoinVertical(lipgloss.Center, m.pages[m.activePageI].View())
	navbar := navbar.FromTitles(m.titles, m.activePageI)

	renderedContent := style.App.Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, navbar, renderedContent)
}
