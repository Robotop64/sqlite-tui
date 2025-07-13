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

	Rows       int
	ViewOffset int
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
	if logMap[t.SelectedTab].Size > t.Rows {
		t.ViewOffset = logMap[t.SelectedTab].Size - t.Rows
	}
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

	t.Rows = height - 3 - 4

	table := lgTable.New().
		Height(height - 3)

	table = table.Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(color.BoxUnselected)).
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
	// fill table to span the entire height
	if logMap[t.SelectedTab].Size < t.Rows {
		diff := t.Rows - logMap[t.SelectedTab].Size
		for i := 0; i < diff; i++ {
			table = table.Row("", "")
		}
	}

	table = table.Offset(t.ViewOffset)

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
		case "up":
			if t.ViewOffset > 0 {
				t.ViewOffset--
			}
		case "down":
			if t.ViewOffset+t.Rows < logMap[t.SelectedTab].Size {
				t.ViewOffset++
			}
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
