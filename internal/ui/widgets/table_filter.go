package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FLayout "fyne.io/fyne/v2/layout"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"

	utils "SQLite-GUI/internal/utils"
)

type TableModel struct {
	Columns []string
	Rows    int
	Filter  Filter
}
type Filter struct {
	Table       string
	ColsVisible []bool
	Limit       int
	Page        int
	SortByCol   int
	SortAsc     bool
}

func NewTableModel(cols []string) TableModel {
	return TableModel{
		Columns: cols,
		Filter: Filter{
			ColsVisible: make([]bool, len(cols)),
			Limit:       0,
			Page:        1,
			SortByCol:   0,
			SortAsc:     true,
		},
	}
}

func NewFilter_Table(table *TableModel, cfg *lua.LTable) *fyne.Container {
	var num_elements int
	if cfg != nil {
		cfg.ForEach(func(k, v lua.LValue) {
			if b, ok := v.(lua.LBool); ok && b == lua.LTrue {
				num_elements++
			}
		})
	}
	components := make([]fyne.CanvasObject, 0, num_elements*2)
	var content *fyne.Container

	var btn_pop_visible_cols, btn_sort_dir *FWidget.Button
	var list_cols *FWidget.List
	var select_col *FWidget.Select

	canvas := fyne.CurrentApp().Driver().AllWindows()[0].Canvas()

	if utils.CheckVal(cfg.RawGetString("table"), true) {

	}

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

	if utils.CheckVal(cfg.RawGetString("columns"), true) {
		longest_col_name := ""
		for _, col := range table.Columns {
			if len(col) > len(longest_col_name) {
				longest_col_name = col
			}
		}

		pop_visible_cols := FWidget.NewPopUp(nil, canvas)
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
		pop_visible_cols.Content = list_cols
		pop_visible_cols.Refresh()

		btn_pop_visible_cols = FWidget.NewButton("Select", func() {
			btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(btn_pop_visible_cols)
			pop_visible_cols.Resize(utils.FitSize(
				utils.Dimensions[float32]{
					Width:  list_cols.MinSize().Width,
					Height: float32(len(table.Columns))*list_cols.MinSize().Height + float32(len(table.Columns)-1)*6,
				},
				utils.FyneToDimensions(content.Size().Subtract(fyne.NewSize(0, btn_pop_visible_cols.Position().Y+btn_pop_visible_cols.Size().Height))),
			).ToFyneSize())
			pop_visible_cols.ShowAtPosition(btnPos.Add(fyne.NewPos(0, btn_pop_visible_cols.Size().Height)))
		})

		components = append(components, FWidget.NewLabel("Columns:"), btn_pop_visible_cols)
	}

	if utils.CheckVal(cfg.RawGetString("sort_by"), true) {
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

		components = append(components, FWidget.NewLabel("Sort by Column:"), select_col)
	}

	if utils.CheckVal(cfg.RawGetString("sort_dir"), true) {
		btn_sort_dir = FWidget.NewButton("Asc", func() {
			if table.Filter.SortAsc {
				table.Filter.SortAsc = false
				btn_sort_dir.SetText("Desc")
			} else {
				table.Filter.SortAsc = true
				btn_sort_dir.SetText("Asc")
			}
		})

		components = append(components, FWidget.NewLabel("Sort Direction:"), btn_sort_dir)
	}

	if utils.CheckVal(cfg.RawGetString("limit"), true) {
		entry_num_rows := NewNumericalEntry()
		entry_num_rows := NewNumericalEntry(false)
		entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
		entry_num_rows.OnSubmitted = func(val string) {
			if val == "" {
				entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
				return
			}
			table.Filter.Limit, _ = strconv.Atoi(val)
		}
		entry_num_rows.OnFocusLost = func() {
			if entry_num_rows.Text != "" {
				table.Filter.Limit, _ = strconv.Atoi(entry_num_rows.Text)
			} else {
				entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
			}
		}

		components = append(components, FWidget.NewLabel("Row Limit:"), entry_num_rows)
	}

	content = FContainer.New(FLayout.NewFormLayout(), components...)

	confirm := FWidget.NewButton("Done", func() {
		fmt.Println(table.GetQuery())
	})
	confirm.Importance = FWidget.LowImportance
	confirm.Refresh()

	return FContainer.NewBorder(FWidget.NewLabel("Table Filter"), confirm, nil, nil, content)
}

func (model *TableModel) GetQuery() string {
	var cols strings.Builder
	for i, col := range model.Columns {
		if model.Filter.ColsVisible[i] {
			if i != 0 {
				cols.WriteString(",")
			}
			cols.WriteString(col)
		}
	}

	filter := ""
	sort := strings.Join([]string{model.Columns[model.Filter.SortByCol], func() string {
		if model.Filter.SortAsc {
			return "ASC"
		}
		return "DESC"
	}()}, " ")

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s", cols.String(), "%s", filter, sort)
}
