package ui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/page"
)

type Model struct {
	pages      []page.Page
	activePage int

	keyPress string
	counter  int
}

var _ tea.Model = &Model{}

func NewModel() *Model {
	return &Model{
		pages: []page.Page{
			&page.WelcomePage{},
		},
		activePage: 0,
		counter:    0,
	}
}

func (m *Model) Init() (cmd tea.Cmd) {
	m.counter = 0
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
	return lipgloss.JoinVertical(lipgloss.Center, strconv.FormatInt(int64(m.counter), 10), m.keyPress)
}
