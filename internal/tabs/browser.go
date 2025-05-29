package tabs

import (
	comp "github.com/Robotop64/sqlite-tui/internal/components"
	style "github.com/Robotop64/sqlite-tui/internal/style"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgTable "github.com/charmbracelet/lipgloss/table"
	lgTree "github.com/charmbracelet/lipgloss/tree"
)

type BrowserTab struct {
	Name     string
	Explorer string
	Tree     *lgTree.Tree
	Table    *comp.Table
}

const border = 1

func (b BrowserTab) Init() tea.Cmd {
	return nil
}

func (b BrowserTab) GetName() string {
	return b.Name
}

func (b BrowserTab) View(width, height int) string {
	hint_height := 2
	sidecolumn_width := width / 5
	maincolumn_width := width - sidecolumn_width
	layout_height := height - hint_height

	sidecolumn := style.Box.
		Width(sidecolumn_width - 2*border).
		Height(layout_height - 2*border)
	// maincolumn := box.
	// Width(maincolumn_width - 2*border).
	// Height(layout_height - 2*border)

	tab := sidecolumn.
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(lipgloss.Width(sidecolumn.String())-6, lipgloss.Center, b.Name),
			"⏵",
		)

	var (
		sidebar,
		content,
		hints string
	)

	//=Side Column========
	//-Explorer------------
	explorer_height := layout_height - lipgloss.Height(tab.String()) - 2*border

	tree := lgTree.New()
	loadSchema(tree)

	explorer := sidecolumn.
		Padding(0, 1).
		Height(explorer_height)

	explorer = explorer.SetString(tree.String())
	//---------------------
	sidebar = lipgloss.JoinVertical(
		lipgloss.Top,
		tab.Render(),
		explorer.Render(),
	)
	//=====================
	//region Main Column=========
	actions_height := 1
	//-Table---------------
	table_height := layout_height - actions_height - 2*border

	table := lgTable.New().
		Width(maincolumn_width).
		Height(table_height)

	table = table.
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
			switch {
			case row == lgTable.HeaderRow:
				return s.Foreground(lipgloss.Color("205")).Bold(true)
			default:
				return s.Foreground(lipgloss.Color("205"))
			}
		})

	table_dims := utils.Dimensions{}
	loadTable(table, &table_dims)
	fillTable(table,
		max(0, table_height-3*border-1-table_dims.Height),
		table_dims.Width,
	)

	//---------------------
	//-Actions-------------
	actions := lgTable.New().
		Width(maincolumn_width).
		Height(actions_height).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center).Foreground(lipgloss.Color("205"))
		})
	actions = actions.Row([]string{"Query", "Filter", "Add", "Sort", "Update"}...)
	//---------------------
	content = lipgloss.JoinVertical(
		lipgloss.Top,
		table.Render(),
		actions.Render(),
	)
	//=====================
	//=Hints===============
	hints = lipgloss.JoinHorizontal(
		lipgloss.Top,
		" tabs:\nalt+(</>)",
		"│\n│",
		" quit:\nctrl+c, q",
	)
	//=====================

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			content,
		),
		hints,
	)
}

func (b BrowserTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	return b, nil
}

func loadTable(t *lgTable.Table, dims *utils.Dimensions) {
	headers := []string{"ID", "Name", "Adress", "Status"}
	data := [][]string{
		{"1", "Alice", "123 Main St", "Completed"},
		{"2", "Bob", "456 Elm St", "Pending"},
		{"3", "Charlie", "789 Oak St", "Cancelled"},
		{"4", "David", "321 Pine St", "In Progress"},
		{"5", "Eve", "654 Maple St", "Completed"},
		{"6", "Frank", "987 Cedar St", "Pending"},
		{"7", "Grace", "159 Birch St", "Cancelled"},
		{"8", "Hank", "753 Spruce St", "In Progress"},
	}

	t.Headers(headers...)
	t.Rows(data...)

	dims.Width = len(headers)
	dims.Height = len(data)
}

func loadSchema(t *lgTree.Tree) {
	t.Root("DB Schema:").
		Child("Tables").
		Child("Users").
		Child("Products").
		Child("Orders")
}

func fillTable(t *lgTable.Table, rows int, cols int) {
	empty_row := make([]string, cols)
	for i := 0; i < rows; i++ {
		t.Row(empty_row...)
	}
}
