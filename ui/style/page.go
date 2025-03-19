package style

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/common"
	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

var App = lipgloss.NewStyle().
	Width(common.AppWidth).
	Height(common.AppHeight).
	Border(lipgloss.RoundedBorder()).
	UnsetBorderTop().
	BorderForeground(color.LuksoPink).
	Align(lipgloss.Center, lipgloss.Center)
