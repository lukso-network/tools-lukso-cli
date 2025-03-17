package color

import "github.com/charmbracelet/lipgloss"

var LuksoPink = lipgloss.Color("#F60B5E")

func ColorStr(str string, color lipgloss.Color) string {
	colorStyle := lipgloss.NewStyle().Foreground(color)

	return colorStyle.Render(str)
}
