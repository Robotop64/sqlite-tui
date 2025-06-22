package main

import (
	"fmt"
	"os"

	tabs "github.com/Robotop64/sqlite-tui/internal/tabs"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"
	persistent "github.com/Robotop64/sqlite-tui/internal/utils/persistent"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Dim       utils.Dimensions
	Tabs      []tabs.Tab
	ActiveTab int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Dim = utils.Dimensions{Width: msg.Width, Height: msg.Height}
		updated_tab, cmd := m.Tabs[m.ActiveTab].Update(msg)
		m.Tabs[m.ActiveTab] = updated_tab
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			persistent.SaveConfig()
			persistent.SaveData()
			return m, tea.Quit
		// case "ctrl+s":
		// 	utils.SaveConfig()
		// 	fmt.Println("Configuration saved.")
		case "ctrl+right":
			prev_active_tab := m.ActiveTab
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			if m.ActiveTab != prev_active_tab {
				m.Tabs[m.ActiveTab].Activate()
			}
			return m, nil
		case "ctrl+left":
			prev_active_tab := m.ActiveTab
			m.ActiveTab = max(m.ActiveTab-1, 0)
			if m.ActiveTab != prev_active_tab {
				m.Tabs[m.ActiveTab].Activate()
			}
			return m, nil
		default:
			updated_tab, cmd := m.Tabs[m.ActiveTab].Update(msg)
			m.Tabs[m.ActiveTab] = updated_tab
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.Tabs[m.ActiveTab].View(m.Dim.Width, m.Dim.Height)
}

func main() {
	// load main config
	if err := persistent.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	if err := persistent.LoadData(); err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}
	// load profiles
	persistent.LoadProfiles()

	m := model{}
	m.Tabs = []tabs.Tab{
		&tabs.ProfileTab{},
		&tabs.BrowserTab{},
	}
	for _, t := range m.Tabs {
		t.Setup()
	}

	m.Tabs[0].Activate()

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
