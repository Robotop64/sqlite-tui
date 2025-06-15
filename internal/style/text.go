package style

import (
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	"github.com/charmbracelet/lipgloss"
)

var Normal = lipgloss.NewStyle().
	Bold(false).
	Foreground(color.TextUnselected)

var Selected = lipgloss.NewStyle().
	Bold(false).
	Foreground(color.TextSelected).
	Background(lipgloss.Color("20"))

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(color.TextHighlight).
	Underline(true)

var Button = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Background(lipgloss.Color("#444444"))
