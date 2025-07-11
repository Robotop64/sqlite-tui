package tabs

import (
	utils "github.com/Robotop64/sqlite-tui/internal/utils"
	"github.com/Robotop64/sqlite-tui/internal/utils/persistent"

	tea "github.com/charmbracelet/bubbletea"
)

type Core struct {
	Dim      utils.Dimensions
	Tabs     map[*utils.FocusElement]Tab
	TabFocus *utils.FocusElement
}

var (
	Profile,
	Browser,
	Hints,
	Logs utils.FocusElement
)

func (m Core) Init() tea.Cmd {
	return nil
}

func (m Core) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Dim = utils.Dimensions{Width: msg.Width, Height: msg.Height}
		updated_tab, cmd := m.Tabs[m.TabFocus].Update(msg)
		m.Tabs[m.TabFocus] = updated_tab
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			persistent.SaveConfig()
			persistent.SaveData()
			return m, tea.Quit
		case "ctrl+up", "ctrl+down", "ctrl+left", "ctrl+right":
			dir := msg.String()[5:]
			curr := m.TabFocus
			next := curr.Move(dir)
			if next == curr {
				return m, nil
			}

			if dir == "up" {
				if _, ok := m.Tabs[next].(*LogTab); ok {
					next.Down = curr
					m.Tabs[next].(*LogTab).SelectedTab = m.Tabs[curr].GetName()
				}
			}
			if dir == "down" {
				if _, ok := m.Tabs[next].(*HintsTab); ok {
					next.Up = curr
					m.Tabs[next].(*HintsTab).SelectedTab = m.Tabs[next]
				}
			}

			nextTab := m.Tabs[next]
			if nextTab.GetName() == "Browser" && len(persistent.Data.Profiles.Paths) == 0 {
				return m, nil
			}
			if next != curr {
				m.TabFocus = next
				nextTab.Activate()
			}
			return m, nil
		default:
			updated_tab, cmd := m.Tabs[m.TabFocus].Update(msg)
			m.Tabs[m.TabFocus] = updated_tab
			return m, cmd
		}
	}
	return m, nil
}

func (m Core) View() string {
	return m.Tabs[m.TabFocus].View(m.Dim.Width, m.Dim.Height)
}

func Init() Core {
	m := Core{
		Tabs: make(map[*utils.FocusElement]Tab),
	}

	m.Tabs[&Profile] = &ProfileTab{}
	m.Tabs[&Browser] = &BrowserTab{}
	m.Tabs[&Hints] = &HintsTab{}
	m.Tabs[&Logs] = &LogTab{}

	utils.FocusGrid([][]*utils.FocusElement{
		{&Logs, &Logs},
		{&Profile, &Browser},
		{&Hints, &Hints},
	})

	for _, t := range m.Tabs {
		t.Setup()
	}

	m.TabFocus = &Profile

	m.Tabs[m.TabFocus].Activate()
	return m
}
