package style

import (
	"github.com/Robotop64/sqlite-tui/internal/style/color"
	"github.com/charmbracelet/lipgloss"
)

var Box lipgloss.Style = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(color.BoxUnselected)
