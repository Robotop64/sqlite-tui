package widgets

import (
	"strconv"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FLayout "fyne.io/fyne/v2/layout"
	FWidget "fyne.io/fyne/v2/widget"

	utils "SQLite-GUI/internal/utils"
)

type TableModel struct {
	Columns []string
	Rows    int
	Filter  Filter
}
type Filter struct {
	ColsVisible []bool
	NumRows     int
	Page        int
	SortByCol   int
	SortAsc     bool
}

func NewTableModel(cols []string) TableModel {
	return TableModel{
		Columns: cols,
		Filter: Filter{
			ColsVisible: make([]bool, len(cols)),
			NumRows:     0,
			Page:        1,
			SortByCol:   0,
			SortAsc:     true,
		},
	}
}

func NewFilter_Table(table *TableModel) *fyne.Container {
	var content *fyne.Container
	var btn_pop_vis_col, btn_pop_sort_col, btn_sort_dir *FWidget.Button
	var list_cols *FWidget.List
	var radio_cols *FWidget.RadioGroup

	canvas := fyne.CurrentApp().Driver().AllWindows()[0].Canvas()
	pop_vis_col := FWidget.NewPopUp(nil, canvas)
	pop_sort_col := FWidget.NewPopUp(nil, canvas)

	num_cols_visible := 0

	create_radio_items := func(arr *int) []string {
		*arr = 0
		for _, visible := range table.Filter.ColsVisible {
			if visible {
				*arr++
			}
		}
		avail_radio_items := make([]string, 0, *arr)
		for i, col := range table.Columns {
			if table.Filter.ColsVisible[i] {
				avail_radio_items = append(avail_radio_items, col)
			}
		}
		return avail_radio_items
	}

	list_cols = FWidget.NewList(
		func() int {
			return len(table.Columns)
		},
		func() fyne.CanvasObject {
			return FWidget.NewCheck("Placeholder", func(b bool) {})
		},
		func(i int, o fyne.CanvasObject) {
			o.(*FWidget.Check).SetText(table.Columns[i])
			o.(*FWidget.Check).SetChecked(table.Filter.ColsVisible[i])
			o.(*FWidget.Check).OnChanged = func(b bool) {
				table.Filter.ColsVisible[i] = b
				radio_cols.Options = create_radio_items(&num_cols_visible)
				if radio_cols.Selected == "" && num_cols_visible > 0 {
					radio_cols.SetSelected(radio_cols.Options[0])
				}
				if radio_cols.Selected == table.Columns[i] && !b && num_cols_visible > 0 {
					radio_cols.SetSelected(radio_cols.Options[0])
				}
				radio_cols.Refresh()
			}
		},
	)
	pop_vis_col.Content = list_cols
	pop_vis_col.Refresh()

	radio_cols = FWidget.NewRadioGroup([]string{}, func(selected string) {
		for i, col := range table.Columns {
			if col == selected {
				table.Filter.SortByCol = i
				break
			}
		}
	})
	radio_cols.Options = create_radio_items(&num_cols_visible)
	if len(radio_cols.Options) > 0 {
		radio_cols.Selected = radio_cols.Options[table.Filter.SortByCol]
	}
	radio_cols.Required = true
	pop_sort_col.Content = FContainer.NewStack(radio_cols, FWidget.NewLabel("No columns selected"))
	pop_sort_col.Refresh()

	btn_pop_vis_col = FWidget.NewButton("Select", func() {
		btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(btn_pop_vis_col)
		pop_vis_col.Resize(utils.FitSize(
			utils.Dimensions[float32]{
				Width:  list_cols.MinSize().Width,
				Height: float32(len(table.Columns))*list_cols.MinSize().Height + float32(len(table.Columns)-1)*6,
			},
			utils.FyneToDimensions(content.Size().Subtract(fyne.NewSize(0, btn_pop_vis_col.Position().Y+btn_pop_vis_col.Size().Height))),
		).ToFyneSize())
		pop_vis_col.ShowAtPosition(btnPos.Add(fyne.NewPos(0, btn_pop_vis_col.Size().Height)))
	})
	btn_pop_sort_col = FWidget.NewButton("Select", func() {
		if num_cols_visible == 0 {
			pop_sort_col.Content.(*fyne.Container).Objects[0].Hide()
			pop_sort_col.Content.(*fyne.Container).Objects[1].Show()
			pop_sort_col.Refresh()
		} else {
			pop_sort_col.Content.(*fyne.Container).Objects[0].Show()
			pop_sort_col.Content.(*fyne.Container).Objects[1].Hide()
			pop_sort_col.Refresh()
		}
		btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(btn_pop_sort_col)
		pop_sort_col.ShowAtPosition(btnPos.Add(fyne.NewPos(0, btn_pop_sort_col.Size().Height)))
	})
	btn_sort_dir = FWidget.NewButton("Asc", func() {
		if table.Filter.SortAsc {
			table.Filter.SortAsc = false
			btn_sort_dir.SetText("Desc")
		} else {
			table.Filter.SortAsc = true
			btn_sort_dir.SetText("Asc")
		}
	})

	entry_num_rows := NewNumericalEntry()
	entry_num_rows.SetText(strconv.Itoa(table.Filter.NumRows))
	entry_num_rows.OnSubmitted = func(val string) {
		if val == "" {
			entry_num_rows.SetText(strconv.Itoa(table.Filter.NumRows))
			return
		}
		table.Filter.NumRows, _ = strconv.Atoi(val)
	}
	entry_num_rows.OnFocusLost = func() {
		if entry_num_rows.Text != "" {
			table.Filter.NumRows, _ = strconv.Atoi(entry_num_rows.Text)
		} else {
			entry_num_rows.SetText(strconv.Itoa(table.Filter.NumRows))
		}
	}

	content = FContainer.New(FLayout.NewFormLayout(),
		FWidget.NewLabel("Columns:"), btn_pop_vis_col,
		FWidget.NewLabel("Sort by Column:"), btn_pop_sort_col,
		FWidget.NewLabel("Number of Rows:"), entry_num_rows,
		FWidget.NewLabel("Sort Direction:"), btn_sort_dir,
	)

	return FContainer.NewBorder(FWidget.NewLabel("Table Filter"), nil, nil, nil, content)
}
