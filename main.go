package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dimension struct {
	width  int
	height int
}

type model struct {
	count int
	dim   dimension
}

func (m model) Init() tea.Cmd {
	// Initialize the model with a default count
	m.count = 0
	// Set default dimensions

	return nil
}

// Update handles messages and updates state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dim = dimension{msg.Width, msg.Height}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.count++
		case "down":
			m.count--
		}
	}
	return m, nil
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))

func (m model) View() string {
	return lipgloss.Place(
		m.dim.width,
		m.dim.height,
		lipgloss.Center,
		lipgloss.Center,
		style.Render(fmt.Sprintf("Counter: %d\n[↑] Increment  [↓] Decrement  [q] Quit", m.count)),
	)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
