package main

import (
	"fmt"
	"os"

	tabs "github.com/Robotop64/sqlite-tui/internal/tabs"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	cfg "github.com/spf13/viper"
)

type model struct {
	Dim       utils.Dimensions
	Tabs      []tabs.Tab
	ActiveTab int
	Profile   *cfg.Viper
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Dim = utils.Dimensions{Width: msg.Width, Height: msg.Height}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "alt+right":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "alt+left":
			m.ActiveTab = max(m.ActiveTab-1, 0)
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
	if err := utils.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	//try to load last profile config
	var profile utils.Profile = nil
	if cfg.GetInt("profiles.last_used") > 0 {
		profilePath := cfg.GetStringSlice("profiles.paths")[cfg.GetInt("profiles.last_used")-1]
		loadedProfile, err := utils.LoadProfile(profilePath)
		if err != nil {
			fmt.Println("Error loading profile:", err)
			loadedProfile = nil
		} else {
			profile = loadedProfile
		}
	}

	m := model{Profile: profile}
	m.Tabs = []tabs.Tab{
		tabs.ProfileTab{Name: "Profiles", Profiles: utils.LoadProfiles()}.PostInit(),
		tabs.BrowserTab{Name: "Browser"},
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
