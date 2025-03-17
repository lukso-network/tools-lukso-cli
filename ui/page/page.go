package page

import tea "github.com/charmbracelet/bubbletea"

// Page represents what a ui page needs in order to correctly take input and render back output.
type Page interface {
	Title() string
	HandleInput(msg tea.KeyMsg)
	View() string
}
