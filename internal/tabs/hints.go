package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Hint struct{}

type HintsTab struct {
	name        string
	Hints       map[Tab]Hint
	SelectedTab Tab
}

func (b *HintsTab) GetName() string {
	return b.name
}

func (b *HintsTab) Init() tea.Cmd {
	return nil
}

func (b *HintsTab) Setup() Tab {
	b.name = "Hints"

	return b
}

func (b *HintsTab) Activate() {

}

func (b *HintsTab) View(width, height int) string {
	return ""
}

func (b *HintsTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		}
	}
	return b, nil
}
