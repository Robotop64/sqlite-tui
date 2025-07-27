package tabs

import (
	"fmt"
	"path/filepath"

	comp "github.com/Robotop64/sqlite-tui/internal/components"
	database "github.com/Robotop64/sqlite-tui/internal/database"
	persistent "github.com/Robotop64/sqlite-tui/internal/persistent"
	style "github.com/Robotop64/sqlite-tui/internal/style"
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	ui "github.com/Robotop64/sqlite-tui/internal/ui"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
)

type BrowserTab struct {
	name       string
	ElemFocus  ElemFocus
	ExplMode   ExplMode
	ActiveList *comp.ListModel[string]
	Lists      []comp.ListModel[string]
	Scripts    []persistent.Script
	Layout     ui.Layout
}

type ExplMode int

const (
	Target ExplMode = iota
	Schema
	View
)

var headers = [3]string{"Targets", "Schema", "Views"}

type ContentMode int

const (
	None ContentMode = iota
	Edit
	Display
)

type ElemFocus int

const (
	Explorer ElemFocus = iota
	Content
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
	b.ExplMode = Target

	b.Lists = make([]comp.ListModel[string], 3)
	b.ActiveList = &b.Lists[0]

	b.Layout = ui.Layout{}

	//this is later in the scripting stuff
	table := &ui.TableWidget{
		Title:   "Test Table",
		Headers: []string{"A", "B", "C"},
		Columns: []any{
			ui.Column[string]{Cells: []string{" ", " ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}},
			ui.Column[string]{Cells: []string{" ", " ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}},
			ui.Column[string]{Cells: []string{" ", " ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}},
		},
		// Headers: []string{"ID", "Name", "Address", "Status", "Actions", "Notes", "Tags", "Created At", "Updated At", "Deleted At", "Archived", "Priority", "Category", "Assigned To", "Due Date", "Completed", "Progress", "Rating", "Feedback", "Attachments", "Comments", "Links", "Related Items", "Custom Field 1", "Custom Field 2", "Custom Field 3"},
		// Columns: []any{
		// 	ui.Column[int]{Cells: []int{1, 2, 3, 4, 5}},
		// 	ui.Column[string]{Cells: []string{"Alice", "Bob", "Charlie", "David", "Eve"}},
		// 	ui.Column[string]{Cells: []string{"123 Main St", "456 Elm St", "789 Oak St", "321 Pine St", "654 Maple St"}},
		// 	ui.Column[string]{Cells: []string{"Completed", "Pending", "Cancelled", "In Progress", "Completed"}},
		// },
		HorizViewPort: ui.Viewport{
			Offset: 0,
		},
		Style: ui.TableStyle{
			Title:      false,
			Scrollbars: [2]bool{false, false},
		},
	}
	b.Layout.Widgets = append(b.Layout.Widgets, table)
	b.Layout.Positions = append(b.Layout.Positions, ui.Position{X: 0, Y: 0})
	b.Layout.Dimensions = append(b.Layout.Dimensions, ui.Dimensions{Width: 80, Height: 20})

	AddLog(b.name, "[STATUS] : Initialized")

	return b
}

func (b *BrowserTab) Activate() {
	targets := persistent.ActiveProfile().Targets

	b.Lists[Target].Items = utils.Map(targets, func(i int, t persistent.Target) string {
		return t.Name
	})

	target := targets[persistent.Data.Profiles.LastTargetUsed]

	AddLog(b.name, "[STATUS] : Activated")
	AddLog(b.name, fmt.Sprintf("[TARGET] : %s", target.Name))
	selectTarget(b, target)
}

func (b *BrowserTab) View(width, height int) string {
	//=Calculations============================================================
	explorer_size := utils.Dimensions{
		Width:  width / 5,
		Height: height,
	}

	content_size := utils.Dimensions{
		Width:  width - explorer_size.Width,
		Height: height,
	}
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
		Bold(true).
		Foreground(color.TextHighlight)

	selector_Box := gen_explorer(b, utils.Dimensions{Width: explorer_size.Width, Height: explorer_size.Height - 3})

	left_Column := lipgloss.JoinVertical(
		lipgloss.Top,
		tab_Box.Render(),
		selector_Box.Render(),
	)
	//=========================================================================

	//=Content Column==========================================================
	right_Column := gen_content(b, content_size).Render()
	//=========================================================================

	//=Hints===================================================================
	// hints := lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	"quit:\nctrl+c|q", "│\n│",
	// 	"save:\nctrl+s", "│\n│",
	// 	"tabs:\nalt+(</>)", "│\n│",
	// 	"profiles:\n   ↑/↓", "│\n│",
	// 	"add:\n +", "│\n│",
	// 	"remove:\n   -", "│\n│",
	// )
	//=========================================================================

	//=Layout==================================================================
	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left_Column,
			right_Column,
		),
		// hints,
	)

	return layout
	//=========================================================================
}

func (b *BrowserTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch b.ElemFocus {
		case Explorer:
			switch msg.String() {
			case "left":
				b.ExplMode = max(b.ExplMode-1, Target)
				AddLog(b.name, fmt.Sprintf("[MODE] : %s", headers[b.ExplMode]))
				b.ActiveList = &b.Lists[b.ExplMode]
				return b, nil
			case "right":

				if b.ExplMode == Target && len(persistent.ActiveProfile().Targets) == 0 {
					AddLog(b.name, "[ERROR] : No targets available, blocking switch")
					return b, nil
				}
				b.ExplMode = min(b.ExplMode+1, View)
				AddLog(b.name, fmt.Sprintf("[MODE] : %s", headers[b.ExplMode]))
				b.ActiveList = &b.Lists[b.ExplMode]
				return b, nil
			case "up":
				b.ActiveList.Focused = max(b.ActiveList.Focused-1, 0)
				return b, nil
			case "down":
				b.ActiveList.Focused = min(b.ActiveList.Focused+1, len(b.ActiveList.Items)-1)
				return b, nil
			case "enter":
				switch b.ExplMode {
				case Target:
					targetIdx := b.ActiveList.Focused
					target := persistent.ActiveProfile().Targets[targetIdx]
					AddLog(b.name, fmt.Sprintf("[TARGET] : %s", target.Name))

					if err := selectTarget(b, target); err == nil {
						b.ActiveList.Selected = b.ActiveList.Focused
					}
				default:
					b.ActiveList.Selected = b.ActiveList.Focused
				}
			}

		}
	}
	return b, nil
}

func gen_explorer(b *BrowserTab, dims utils.Dimensions) lipgloss.Style {
	view := style.Box.
		Padding(0, 1).
		Width(dims.Width - 2).
		Height(dims.Height - 2).
		BorderForeground(utils.Ifelse(b.ElemFocus == Explorer, color.BoxSelected, color.BoxUnselected).(lipgloss.Color))

	list := lgList.New()

	header := headers[b.ExplMode]
	selected := b.ActiveList.Selected
	focused := b.ActiveList.Focused
	list.Items(b.ActiveList.Items)

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

	view = view.SetString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			style.Title.SetString(header).Render(),
			utils.Ifelse(len(b.ActiveList.Items) > 0, list.String(), "...").(string),
		),
	)

	return view
}

func gen_content(b *BrowserTab, dims utils.Dimensions) lipgloss.Style {
	view := style.Box.
		// Padding(0, 1).
		Width(dims.Width - 2).
		Height(dims.Height - 2)

	switch b.ExplMode {
	case Target:
		return view
	case Schema:
		header := "Creation Statement:"
		body := database.ActiveSchema.CreationSQL[b.ActiveList.Selected]
		content := lipgloss.JoinVertical(
			lipgloss.Top,
			style.Title.SetString(header).Render(),
			utils.Ifelse(len(body) > 0, body, "No creation statement available").(string),
		)
		view = view.SetString(content)
		return view
	case View:

		table := b.Layout.Widgets[0].(*ui.TableWidget)

		table.Style.BaseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

		table.Style.MaxDimensions = ui.Dimensions{
			Width:  dims.Width - 2,
			Height: dims.Height - 4,
		}

		return view.SetString(b.Layout.Render())
	}

	// content = content.SetString(
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Top,
	// 		style.Title.SetString("Content").Render(),
	// 		utils.Ifelse(len(b.ActiveList.Items) > 0, b.ActiveList.String(), "...").(string),
	// 	),
	// )

	return style.Box
}

func selectTarget(b *BrowserTab, target persistent.Target) error {
	if dberr := database.SetTarget(persistent.ActiveProfilePath(), target); dberr != nil {
		AddLog(b.name, fmt.Sprintf("[ERROR] : Updating target: %v", dberr))
		return dberr
	}

	if err := database.SetSchema(); err != nil {
		AddLog(b.name, fmt.Sprintf("[ERROR] : Updating schema: %v", err))
	}
	b.Lists[Schema].Items = utils.Map(database.ActiveSchema.TableNames, func(i int, name string) string {
		return name
	})

	b.Scripts = make([]persistent.Script, len(target.ScriptPaths))
	for i, path := range target.ScriptPaths {
		path = utils.RelativeToAbsolutePath(filepath.Dir(persistent.ActiveProfilePath()), path)
		script, err := persistent.LoadScript(path)
		if err != nil {
			AddLog(b.name, fmt.Sprintf("[ERROR] :  Loading script [%s]: %v", path, err))
			continue
		}
		b.Scripts[i] = script
	}

	b.Lists[View].Items = utils.Map(b.Scripts, func(i int, script persistent.Script) string {
		return script.MetaData.Name
	})
	return nil
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
