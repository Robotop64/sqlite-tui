package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dimension struct {
	Width  int
	Height int
}

type model struct {
	Dim       Dimension
	Tabs      []string
	ActiveTab int
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Dim = Dimension{msg.Width, msg.Height}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "left":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		}
	}
	return m, nil
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Padding(0, 1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))

func (m model) View() string {
	sidecolumn_width := m.Dim.Width / 4
	label := m.Tabs[m.ActiveTab]
	tab := style.Render("⏴", lipgloss.PlaceHorizontal(sidecolumn_width-6, lipgloss.Center, label), "⏵")

	// explorerRatio := 0.25
	// explorerDim := Dimension{
	// 	Width:  int(float64(m.Dim.Width) * 0.25),
	// 	Height: m.Dim.Height - 3,
	// }
	// contentDim := Dimension{
	// 	Width:  int(float64(m.Dim.Width)*(1-explorerRatio)) - 3,
	// 	Height: m.Dim.Height - 3,
	// }

	// explorer := style.
	// 	Width(explorerDim.Width).
	// 	Height(explorerDim.Height).
	// 	Render("Explorer!")
	// content := style.
	// 	Width(contentDim.Width).
	// 	Height(contentDim.Height).
	// 	Render("Content Area!")

	return tab
}

func main() {
	tabs := []string{"Browser", "Builder", "Editor"}
	m := model{
		Tabs:      tabs,
		ActiveTab: 0,
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
