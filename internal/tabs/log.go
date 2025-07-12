package tabs

import (
	"time"

	style "github.com/Robotop64/sqlite-tui/internal/style"
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgTable "github.com/charmbracelet/lipgloss/table"
)

type Log struct {
	Timestamps []time.Time
	Messages   []string
	Size       int
}

type LogTab struct {
	name        string
	SelectedTab string
}

var logMap map[string]Log = make(map[string]Log)

func (t *LogTab) GetName() string {
	return t.name
}

func (t *LogTab) Init() tea.Cmd {
	return nil
}

func (t *LogTab) Setup() Tab {
	t.name = "Logs"
	return t
}

func (t *LogTab) Activate() {
}

func (t *LogTab) View(width, height int) string {
	tab_Box := style.Box.
		Width(width/5-2).
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(width/5-6, lipgloss.Center, "Logs "+"("+t.SelectedTab+")"),
			"⏵",
		).
		Bold(true).
		Foreground(color.TextHighlight)

	table := lgTable.New()

	table = table.Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == lgTable.HeaderRow:
				return lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
			case col == 0:
				return lipgloss.NewStyle().Padding(0, 1).Width(10)
			default:
				return lipgloss.NewStyle().Padding(0, 1).Width(width - 13)
			}
		})

	table = table.Headers("Time", "Message")
	for i := 0; i < logMap[t.SelectedTab].Size; i++ {
		entry := logMap[t.SelectedTab]
		table = table.Row(entry.Timestamps[i].Format("15:04:05"), entry.Messages[i])
	}

	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		tab_Box.Render(),
		table.Render(),
	)

	return layout
	//=========================================================================
}

func (t *LogTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		}
	}
	return t, nil
}

func AddLog(tab string, message string) {
	timestamp := time.Now()
	log := logMap[tab]
	log.Timestamps = append(logMap[tab].Timestamps, timestamp)
	log.Messages = append(log.Messages, message)
	log.Size++
	logMap[tab] = log
}
