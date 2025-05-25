package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Box lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))
