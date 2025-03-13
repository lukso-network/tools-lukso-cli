package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/ui"
)

func Ui(ctx *cli.Context) (err error) {
	m := &ui.Model{}
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err = p.Run()
	return
}
