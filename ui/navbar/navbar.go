package navbar

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/style"
)

func FromTitles(titles []string, activeI int) (navbar string) {
	navbar = "╭"
	var tabStyle lipgloss.Style

	for i, tabTitle := range titles {
		if activeI == i {
			tabStyle = style.ActiveNavbarTab
		} else {
			tabStyle = style.InactiveNavbarTab
		}

		tab := tabStyle.Render(tabTitle)
		navbar = lipgloss.JoinVertical(lipgloss.Bottom, navbar, tab)
	}

	return
}
