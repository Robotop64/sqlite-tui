package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Tab interface {
	GetName() string
	Init() tea.Cmd
	View(width, height int) string
	Update(msg tea.Msg) (Tab, tea.Cmd)
}
