package navbar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/common"
	"github.com/lukso-network/tools-lukso-cli/ui/style"
	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

func FromTitles(titles []string, activeI int) (navbar string) {
	navbar = color.ColorStr("╭", color.LuksoPink)

	var tabStyle lipgloss.Style

	for i, tabTitle := range titles {
		if activeI == i {
			tabStyle = style.ActiveNavbarTab
		} else {
			tabStyle = style.InactiveNavbarTab
		}

		tab := tabStyle.Render(tabTitle)
		navbar = lipgloss.JoinHorizontal(lipgloss.Bottom, navbar, tab)
	}

	tabsWidth := lipgloss.Width(navbar)
	fill := strings.Repeat("─", common.AppWidth-tabsWidth+1) + "╮"
	fill = color.ColorStr(fill, color.LuksoPink)
	navbar = lipgloss.JoinHorizontal(lipgloss.Bottom, navbar, fill)

	return
}
