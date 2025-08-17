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
	var btn_pop_vis_col, btn_sort_dir *FWidget.Button
	var list_cols *FWidget.List
	var select_col *FWidget.Select

	canvas := fyne.CurrentApp().Driver().AllWindows()[0].Canvas()
	pop_vis_col := FWidget.NewPopUp(nil, canvas)

	num_cols_visible := 0

	create_items := func(arr *int) []string {
		*arr = 0
		for _, visible := range table.Filter.ColsVisible {
			if visible {
				*arr++
			}
		}
		avail_items := make([]string, 0, *arr)
		for i, col := range table.Columns {
			if table.Filter.ColsVisible[i] {
				avail_items = append(avail_items, col)
			}
		}
		return avail_items
	}
	longest_col_name := ""
	for _, col := range table.Columns {
		if len(col) > len(longest_col_name) {
			longest_col_name = col
		}
	}
	list_cols = FWidget.NewList(
		func() int {
			return len(table.Columns)
		},
		func() fyne.CanvasObject {
			return FWidget.NewCheck(longest_col_name, func(b bool) {})
		},
		func(i int, o fyne.CanvasObject) {
			o.(*FWidget.Check).SetText(table.Columns[i])
			o.(*FWidget.Check).SetChecked(table.Filter.ColsVisible[i])
			o.(*FWidget.Check).OnChanged = func(b bool) {
				table.Filter.ColsVisible[i] = b
				select_col.Options = create_items(&num_cols_visible)
				if select_col.Selected == table.Columns[i] && !b && num_cols_visible > 0 {
					select_col.SetSelectedIndex(0)
				}
				if num_cols_visible == 0 {
					select_col.ClearSelected()
				}
				select_col.Refresh()
			}
		},
	)
	pop_vis_col.Content = list_cols
	pop_vis_col.Refresh()

	select_col = FWidget.NewSelect([]string{}, func(selected string) {
		for i, col := range table.Columns {
			if col == selected {
				table.Filter.SortByCol = i
				break
			}
		}
	})
	select_col.PlaceHolder = "Select Column"
	select_col.Alignment = fyne.TextAlignCenter
	select_col.Options = create_items(&num_cols_visible)
	if len(select_col.Options) > 0 {
		select_col.Selected = select_col.Options[table.Filter.SortByCol]
	}

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
		FWidget.NewLabel("Sort by Column:"), select_col,
		FWidget.NewLabel("Number of Rows:"), entry_num_rows,
		FWidget.NewLabel("Sort Direction:"), btn_sort_dir,
	)

	return FContainer.NewBorder(FWidget.NewLabel("Table Filter"), nil, nil, nil, content)
}
