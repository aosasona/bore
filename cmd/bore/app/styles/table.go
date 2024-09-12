package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var BaseTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(TintColor))

func TableStyle() table.Styles {
	style := table.DefaultStyles()

	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color(TintColor)).
		BorderBottom(true).
		Bold(true)

	style.Selected = style.Selected.
		Foreground(lipgloss.Color(ForegroundColor)).
		Background(lipgloss.Color(TintColor)).
		Bold(false)

	return style
}
