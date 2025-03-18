package page

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/common"
	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

type WelcomePage struct{}

var _ Page = &WelcomePage{}

func (p *WelcomePage) Title() string {
	return "Home"
}

func (p *WelcomePage) HandleInput(msg tea.KeyMsg) {
}

func (p *WelcomePage) View() string {
	logo := common.LuksoLogo
	logo = color.ColorChar(logo, "*", color.LuksoPink)
	logo = color.ColorChar(logo, "#", color.White)

	banner := lipgloss.NewStyle().Bold(true).Render(common.LuksoBanner)

	content := lipgloss.JoinVertical(lipgloss.Center, logo, banner)

	return content
}
