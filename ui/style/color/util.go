package color

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ColorStr colors the whole text to a given color.
func ColorStr(str string, color lipgloss.Color) string {
	colorStyle := lipgloss.NewStyle().Foreground(color)

	return colorStyle.Render(str)
}

// ColorChar colors all characters matching char in the str, to a specified color.
// If the len(char) > 1, all of the characters in the char are colored.
func ColorChar(str, char string, color lipgloss.Color) (colored string) {
	strChars := strings.Split(str, "")
	chars := strings.Split(char, "")
	for i, strChar := range strChars {
		if slices.Contains(chars, strChar) {
			strChars[i] = ColorStr(strChar, color)
		}
	}

	return strings.Join(strChars, "")
}
