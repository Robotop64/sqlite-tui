package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FBind "fyne.io/fyne/v2/data/binding"
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
	Table string // which table to apply the fitler on

	ColsVisible []bool         // which columns are visible
	GroupByCol  int            // which column to group by
	Sort        SortFilter     // sorting options
	Values      []ValuesFilter // which values to filter on which columns

	Limit int // how many rows should be displayed per page
	Page  int // which page to display
}
type SortFilter struct {
	ByCol int
	Asc   bool
}
type ValuesFilter struct {
	ByCol int
	Val   string
}

func NewTableModel(cols []string) TableModel {
	return TableModel{
		Columns: cols,
		Filter: Filter{
			Table:       "",
			ColsVisible: make([]bool, len(cols)),
			Sort: SortFilter{
				ByCol: 0,
				Asc:   true,
			},
			GroupByCol: 0,
			Values:     make([]ValuesFilter, 0),
			Limit:      50,
			Page:       1,
		},
	}
}

func NewFilter_Table(table *TableModel, cfg *lua.LTable) *fyne.Container {
	var content *fyne.Container
	var num_elements int = 0
	var num_cols_visible int = 0
	var last_checked int = -1

	cols_visible := make([]string, 0)
	bind_cols_visible := FBind.BindStringList(&cols_visible)

	//calculate required number of editor elements
	if cfg != nil {
		cfg.ForEach(func(k, v lua.LValue) {
			if b, ok := v.(lua.LBool); ok && b == lua.LTrue {
				num_elements++
			}
		})
	}
	components := make([]fyne.CanvasObject, 0, num_elements*2)
	canvas := fyne.CurrentApp().Driver().AllWindows()[0].Canvas()

	// set number of visible columns and create items
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

	// add the elements to the widget container, if the lua flags are set to true

	if utils.CheckVal(cfg.RawGetString("table"), true) {

	}

	if utils.CheckVal(cfg.RawGetString("columns"), true) {
		longest_col_name := ""
		for _, col := range table.Columns {
			if len(col) > len(longest_col_name) {
				longest_col_name = col
			}
		}

		pop_visible_cols := FWidget.NewPopUp(nil, canvas)
		list_cols := FWidget.NewList(
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
					bind_cols_visible.Set(create_items(&num_cols_visible))
					last_checked = i
				}
			},
		)
		pop_visible_cols.Content = list_cols
		pop_visible_cols.Refresh()
		var btn_pop_visible_cols *FWidget.Button
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
		select_col := FWidget.NewSelect([]string{}, func(selected string) {
			for i, col := range table.Columns {
				if col == selected {
					table.Filter.Sort.ByCol = i
					break
				}
			}
		})
		select_col.PlaceHolder = "Select Column"
		select_col.Alignment = fyne.TextAlignCenter
		if len(select_col.Options) > 0 {
			select_col.Selected = select_col.Options[table.Filter.Sort.ByCol]
		}
		bind_cols_visible.AddListener(FBind.NewDataListener(func() {
			select_col.Options = cols_visible
			if last_checked >= 0 && select_col.Selected == table.Columns[last_checked] && !table.Filter.ColsVisible[last_checked] && num_cols_visible > 0 {
				select_col.SetSelectedIndex(0)
			}
			if num_cols_visible == 0 {
				select_col.ClearSelected()
			}
			select_col.Refresh()
		}))

		components = append(components, FWidget.NewLabel("Sort by Column:"), select_col)
	}

	if utils.CheckVal(cfg.RawGetString("sort_dir"), true) {
		var btn_sort_dir *FWidget.Button
		btn_sort_dir = FWidget.NewButton("Asc", func() {
			if table.Filter.Sort.Asc {
				table.Filter.Sort.Asc = false
				btn_sort_dir.SetText("Desc")
			} else {
				table.Filter.Sort.Asc = true
				btn_sort_dir.SetText("Asc")
			}
		})

		components = append(components, FWidget.NewLabel("Sort Direction:"), btn_sort_dir)
	}

	if utils.CheckVal(cfg.RawGetString("group_by"), true) {
		select_col_group := FWidget.NewSelect([]string{}, func(selected string) {
			for i, col := range table.Columns {
				if col == selected {
					table.Filter.GroupByCol = i
					break
				}
			}
		})
		select_col_group.PlaceHolder = "Select Column"
		select_col_group.Alignment = fyne.TextAlignCenter
		if len(select_col_group.Options) > 0 {
			select_col_group.Selected = select_col_group.Options[table.Filter.GroupByCol]
		}
		bind_cols_visible.AddListener(FBind.NewDataListener(func() {
			select_col_group.Options = cols_visible
			if last_checked >= 0 && select_col_group.Selected == table.Columns[last_checked] && !table.Filter.ColsVisible[last_checked] && num_cols_visible > 0 {
				select_col_group.SetSelectedIndex(0)
			}
			if num_cols_visible == 0 {
				select_col_group.ClearSelected()
			}
			select_col_group.Refresh()
		}))

		components = append(components, FWidget.NewLabel("Group by Column:"), select_col_group)
	}

	if utils.CheckVal(cfg.RawGetString("filter"), true) {

	}

	if utils.CheckVal(cfg.RawGetString("limit"), true) {
		entry_num_rows := NewNumericalEntry(false)
		entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
		entry_num_rows.OnChanged = func(val string) {
			if val == "" {
				// entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
				return
			}
			if ival, err := strconv.Atoi(val); err == nil {
				table.Filter.Limit = ival
			}
		}
		// entry_num_rows.OnChanged = entry_num_rows.OnSubmitted
		entry_num_rows.OnFocusLost = func() {
			if entry_num_rows.Text != "" {
				table.Filter.Limit, _ = strconv.Atoi(entry_num_rows.Text)
			} else {
				entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
			}
		}

		components = append(components, FWidget.NewLabel("Row Limit:"), entry_num_rows)
	}

	if utils.CheckVal(cfg.RawGetString("page"), true) {
		entry_page := NewNumericalEntry(false)
		entry_page.SetText(strconv.Itoa(table.Filter.Page))
		entry_page.OnChanged = func(val string) {
			if val == "" {
				// entry_page.SetText(strconv.Itoa(table.Filter.Page))
				return
			}
			if ival, err := strconv.Atoi(val); err == nil {
				if ival < 1 {
					entry_page.SetText("1")
					return
				}
				table.Filter.Page = ival
			}
		}
		entry_page.OnFocusLost = func() {
			if entry_page.Text != "" {
				table.Filter.Page, _ = strconv.Atoi(entry_page.Text)
			} else {
				entry_page.SetText(strconv.Itoa(table.Filter.Page))
			}
		}

		components = append(components, FWidget.NewLabel("Page:"), entry_page)
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
	sort := strings.Join([]string{model.Columns[model.Filter.Sort.ByCol], func() string {
		if model.Filter.Sort.Asc {
			return "ASC"
		}
		return "DESC"
	}()}, " ")

	page := (model.Filter.Page - 1) * model.Filter.Limit

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d", cols.String(), "%s", filter, sort, model.Filter.Limit, page)
}
