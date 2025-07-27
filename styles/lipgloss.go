package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var FormatStyle = NewStyles()

type LipglossStyles struct {
	Base      lipgloss.Style
	Title     lipgloss.Style
	Highlight lipgloss.Style
	Error     lipgloss.Style
}

// NewStyles returns a new LipglossStyles struct with predefined styles.
//
// Args:
//   - None
//
// Returns:
//   - LipglossStyles: A struct containing various lipgloss styles.
func NewStyles() LipglossStyles {
	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Background(lipgloss.Color("#FAFAFA")).
		Italic(true)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true).
		Italic(true)

	ErrorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F5F")).
		Background(lipgloss.Color("#1C1C1C")).
		Bold(true).
		Padding(0, 2).
		Align(lipgloss.Center)

	return LipglossStyles{
		Base:      baseStyle,
		Title:     titleStyle,
		Highlight: highlightStyle,
		Error:     ErrorStyle,
	}
}

// StyliseMessage applies a given lipgloss style to a message string.
//
// Args:
//   - message: The string message to stylise.
//   - style: The lipgloss style to apply.
//
// Returns:
//   - string: The stylised message string.
func StyliseMessage(message string, style lipgloss.Style) string {
	return style.Render(message)
}
