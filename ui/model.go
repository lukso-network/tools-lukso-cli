package ui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.keyPress = msg.String()

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "left":
			m.counter--
		case "right":
			m.counter++
		}
	}

	return m, nil
}

func (m *Model) View() string {
	content := lipgloss.JoinVertical(lipgloss.Center, strconv.FormatInt(int64(m.counter), 10), m.keyPress)
	navbar := navbar.FromTitles(m.titles, m.activePageI)

	renderedContent := style.App.Render(content)

	return lipgloss.JoinHorizontal(lipgloss.Left, navbar, renderedContent)
}
