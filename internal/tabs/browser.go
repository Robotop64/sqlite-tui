package tabs

import (
	comp "github.com/Robotop64/sqlite-tui/internal/components"
	style "github.com/Robotop64/sqlite-tui/internal/style"
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
)

type BrowserTab struct {
	name       string
	ElemFocus  ElemFocus
	ExplMode   ExplMode
	SchemaList comp.ListModel[string]
	ViewsList  comp.ListModel[string]
}

type ExplMode int

const (
	Schema ExplMode = iota
	Views
)

type ElemFocus int

const (
	Explorer ElemFocus = iota
	Viewer
)

func (b *BrowserTab) GetName() string {
	return b.name
}

func (b *BrowserTab) Init() tea.Cmd {
	return nil
}

func (b *BrowserTab) Setup() Tab {
	b.name = "Browser"
	b.ElemFocus = Explorer
	b.ExplMode = Schema

	return b
}

func (b *BrowserTab) Activate() {
}

func (b *BrowserTab) View(width, height int) string {
	//=Calculations============================================================
	hint_height := 2

	explorer_size := utils.Dimensions{
		Width:  width / 5,
		Height: height - hint_height,
	}

	// content_size := utils.Dimensions{
	// 	Width:  width - explorer_size.Width,
	// 	Height: height - hint_height,
	// }
	//=========================================================================

	//=Left Column=============================================================
	tab_Box := style.Box.
		Width(explorer_size.Width-2).
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(explorer_size.Width-6, lipgloss.Center, b.name),
			"⏵",
		).
		Foreground(color.TextHighlight)

	selector_Box := gen_explorer(b, utils.Dimensions{Width: explorer_size.Width, Height: explorer_size.Height - 3})

	left_Column := lipgloss.JoinVertical(
		lipgloss.Top,
		tab_Box.Render(),
		selector_Box.Render(),
	)
	//=========================================================================

	//=Content Column==========================================================
	// right_Column := gen_editor(b, content_size).Render()
	//=========================================================================

	//=Hints===================================================================
	hints := lipgloss.JoinHorizontal(
		lipgloss.Top,
		"quit:\nctrl+c|q", "│\n│",
		"save:\nctrl+s", "│\n│",
		"tabs:\nalt+(</>)", "│\n│",
		"profiles:\n   ↑/↓", "│\n│",
		"add:\n +", "│\n│",
		"remove:\n   -", "│\n│",
	)
	//=========================================================================

	//=Layout==================================================================
	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left_Column,
			// right_Column,
		),
		hints,
	)

	return layout
	//=========================================================================
}

func (b *BrowserTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	return b, nil
}

func gen_explorer(b *BrowserTab, dims utils.Dimensions) lipgloss.Style {
	view := style.Box.
		Padding(0, 1).
		Width(dims.Width - 2).
		Height(dims.Height - 2).
		BorderForeground(utils.Ifelse(b.ElemFocus == Explorer, color.BoxSelected, color.BoxUnselected).(lipgloss.Color))

	list := lgList.New()

	var selected int
	var focused int
	var header string

	switch b.ExplMode {
	case Schema:
		header = "DB Schema:"
		selected = b.SchemaList.Selected
		focused = b.SchemaList.Focused
		list.Items(b.SchemaList.Items)
	case Views:
		header = "Views:"
		selected = b.ViewsList.Selected
		focused = b.ViewsList.Focused
		list.Items(b.ViewsList.Items)
	}

	list.Enumerator(func(l lgList.Items, i int) string {
		if i == selected {
			return ">"
		}
		return "•"
	})

	list.ItemStyleFunc(func(_ lgList.Items, i int) lipgloss.Style {
		if i == focused {
			return style.Selected
		}
		return style.Normal
	})

	view.SetString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			style.Title.SetString(header).Render(),
			list.String(),
		),
	)

	return view
}

// func loadTable(t *lgTable.Table, dims *utils.Dimensions) {
// 	headers := []string{"ID", "Name", "Adress", "Status"}
// 	data := [][]string{
// 		{"1", "Alice", "123 Main St", "Completed"},
// 		{"2", "Bob", "456 Elm St", "Pending"},
// 		{"3", "Charlie", "789 Oak St", "Cancelled"},
// 		{"4", "David", "321 Pine St", "In Progress"},
// 		{"5", "Eve", "654 Maple St", "Completed"},
// 		{"6", "Frank", "987 Cedar St", "Pending"},
// 		{"7", "Grace", "159 Birch St", "Cancelled"},
// 		{"8", "Hank", "753 Spruce St", "In Progress"},
// 	}

// 	t.Headers(headers...)
// 	t.Rows(data...)

// 	dims.Width = len(headers)
// 	dims.Height = len(data)
// }

// func loadSchema(t *lgTree.Tree) {
// 	t.Root("DB Schema:").
// 		Child("Tables").
// 		Child("Users").
// 		Child("Products").
// 		Child("Orders")
// }

// func fillTable(t *lgTable.Table, rows int, cols int) {
// 	empty_row := make([]string, cols)
// 	for i := 0; i < rows; i++ {
// 		t.Row(empty_row...)
// 	}
// }
