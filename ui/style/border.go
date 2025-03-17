package style

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

var activeNavbarTabBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       " ",
	Left:         "│",
	Right:        "│",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "╯",
	BottomRight:  "╰",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

var inactiveNavbarTabBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       "─",
	Left:         "│",
	Right:        "│",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "┴",
	BottomRight:  "┴",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

var navbarBorder = lipgloss.NewStyle().
	BorderForeground(color.LuksoPink)

var ActiveNavbarTab = navbarBorder.
	Border(activeNavbarTabBorder)

var InactiveNavbarTab = navbarBorder.
	Border(inactiveNavbarTabBorder)
