package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Normal = lipgloss.NewStyle().
	Bold(false).
	Foreground(lipgloss.Color("205"))

var Selected = lipgloss.NewStyle().
	Bold(false).
	Foreground(lipgloss.Color("63")).
	Background(lipgloss.Color("20"))

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Underline(true)

var Button = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Background(lipgloss.Color("#444444"))
