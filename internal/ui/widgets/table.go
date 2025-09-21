package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FBind "fyne.io/fyne/v2/data/binding"
	FLayout "fyne.io/fyne/v2/layout"
	FTheme "fyne.io/fyne/v2/theme"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"

	persistent "SQLite-GUI/internal/persistent"
	utils "SQLite-GUI/internal/utils"
)

func NewPrimitiveTable(data [][]string) *FWidget.Table {
	table := FWidget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return FWidget.NewLabel("Placeholder")
		},
		func(i FWidget.TableCellID, o fyne.CanvasObject) {
			o.(*FWidget.Label).SetText(data[i.Row][i.Col])
		},
	)
	return table
}

// TODO make generic
func NewPrimitiveEditableTable(data *[][]string, dirtyRows *[]int) *FWidget.Table {
	type CellFocus struct {
		Row int
		Col int
	}
	var cell_edit *CellFocus
	var table *FWidget.Table

	table = FWidget.NewTable(
		func() (int, int) {
			return len(*data), len((*data)[0])
		},
		func() fyne.CanvasObject {
			lbl := FWidget.NewLabel("Placeholder")
			entry := FWidget.NewEntry()
			entry.Hide()
			return FContainer.NewStack(lbl, entry)
		},
		func(i FWidget.TableCellID, o fyne.CanvasObject) {
			cell := o.(*fyne.Container)
			val := (*data)[i.Row][i.Col]

			lbl := cell.Objects[0].(*FWidget.Label)
			entry := cell.Objects[1].(*FWidget.Entry)

			if lbl.Text != val {
				lbl.SetText(val)
			}
			if entry.Text != val {
				entry.SetText(val)
			}

			updateData := func(val string) {
				if (*data)[i.Row][i.Col] != val && !utils.Contains(*dirtyRows, i.Row) {
					*dirtyRows = append(*dirtyRows, i.Row)
					fmt.Println("Row marked as dirty:", i.Row)
				}

				//TODO Check if type casting is valid
				(*data)[i.Row][i.Col] = val
			}

			entry.OnSubmitted = func(val string) {
				updateData(val)
				cell_edit = nil
				table.Refresh()
			}

			if cell_edit != nil && cell_edit.Row == i.Row && cell_edit.Col == i.Col {
				lbl.Hide()
				entry.Show()
			} else {
				lbl.Show()
				entry.Hide()
			}
		},
	)
	table.OnSelected = func(id FWidget.TableCellID) {
		cell_edit = &CellFocus{Row: id.Row, Col: id.Col}
		table.Refresh()
	}

	return table
}

type TableConfig struct {
	Editable bool
	Enabled  struct {
		Title struct {
			Enabled   bool
			Text      string
			Alignment string
			Style     struct {
				Bold   bool
				Italic bool
			}
		}
		Headers struct {
			Enabled bool
			Column  bool
			Row     bool
		}
		Filter FilterConfig
	}
}

type FilterConfig struct {
	Enabled bool
	Table   bool
	Columns bool
	SortBy  bool
	SortDir bool
	GroupBy bool
	Filter  bool
	Limit   bool
	Page    bool
}

func (cfg FilterConfig) num_components() int {
	num := 0
	items := []bool{cfg.Table, cfg.Columns, cfg.SortBy, cfg.SortDir, cfg.GroupBy, cfg.Filter, cfg.Limit, cfg.Page}
	for _, v := range items {
		if v {
			num++
		}
	}
	return num
}

func (cfg TableConfig) FromLuaTable(tbl *lua.LTable) (TableConfig, error) {
	if tbl == nil {
		return cfg, nil
	}

	if utils.CheckVal(tbl.RawGetString("editable"), true) {
		cfg.Editable = true
	}
	if enable_tbl, ok := tbl.RawGetString("enable").(*lua.LTable); ok {
		if title_tbl, ok := enable_tbl.RawGetString("title").(*lua.LTable); ok {
			cfg.Enabled.Title.Enabled = true
			if text, ok := title_tbl.RawGetString("text").(lua.LString); ok {
				cfg.Enabled.Title.Text = text.String()
			}
			if alignment, ok := title_tbl.RawGetString("alignment").(lua.LString); ok {
				switch alignment.String() {
				case "left", "center", "right":
					cfg.Enabled.Title.Alignment = alignment.String()
				default:
					return cfg, fmt.Errorf("invalid title.alignment value: %s", alignment.String())
				}
			}
			if style_tbl, ok := title_tbl.RawGetString("style").(*lua.LTable); ok {
				if bold, ok := style_tbl.RawGetString("bold").(lua.LBool); ok && bold == lua.LTrue {
					cfg.Enabled.Title.Style.Bold = true
				}
				if italic, ok := style_tbl.RawGetString("italic").(lua.LBool); ok && italic == lua.LTrue {
					cfg.Enabled.Title.Style.Italic = true
				}
			}
		}

		if header_tbl, ok := enable_tbl.RawGetString("headers").(*lua.LTable); ok {
			cfg.Enabled.Headers.Enabled = true
			if utils.CheckVal(header_tbl.RawGetString("column"), true) {
				cfg.Enabled.Headers.Column = true
			}
			if utils.CheckVal(header_tbl.RawGetString("row"), true) {
				cfg.Enabled.Headers.Row = true
			}
		}
		if filter_tbl, ok := enable_tbl.RawGetString("filter").(*lua.LTable); ok {
			cfg.Enabled.Filter.Enabled = true
			if utils.CheckVal(filter_tbl.RawGetString("table"), true) {
				cfg.Enabled.Filter.Table = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("columns"), true) {
				cfg.Enabled.Filter.Columns = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("sort_by"), true) {
				cfg.Enabled.Filter.SortBy = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("sort_dir"), true) {
				cfg.Enabled.Filter.SortDir = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("group_by"), true) {
				cfg.Enabled.Filter.GroupBy = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("filter"), true) {
				cfg.Enabled.Filter.Filter = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("limit"), true) {
				cfg.Enabled.Filter.Limit = true
			}
			if utils.CheckVal(filter_tbl.RawGetString("page"), true) {
				cfg.Enabled.Filter.Page = true
			}
		}
	}

	return cfg, nil
}

type TableModel struct {
	Source  *persistent.Source
	Columns []string
	Rows    []int
}

type FilterModel struct {
	Table string // which table to apply the fitler on

	ColsVisible []bool // which columns are visible
	GroupByCol  int    // which column to group by
	Sort        struct {
		ByCol int
		Asc   bool
	} // sorting options
	Values []struct {
		ByCol int
		Val   string
	} // which values to filter on which columns

	Limit int // how many rows should be displayed per page
	Page  int // which page to display
}

func NewFilterModel(cols int) FilterModel {
	return FilterModel{
		Table:       "",
		ColsVisible: make([]bool, cols),
		Values: make([]struct {
			ByCol int
			Val   string
		}, 0),
		Limit: 50,
		Page:  1,
	}
}

func GenQuery(table TableModel, cfg_filter FilterModel) string {
	var cols strings.Builder
	for i, col := range table.Columns {
		if cfg_filter.ColsVisible[i] {
			if i != 0 {
				cols.WriteString(",")
			}
			cols.WriteString(col)
		}
	}

	filter := ""
	sort := strings.Join([]string{table.Columns[cfg_filter.Sort.ByCol], func() string {
		if cfg_filter.Sort.Asc {
			return "ASC"
		}
		return "DESC"
	}()}, " ")

	page := (cfg_filter.Page - 1) * cfg_filter.Limit

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d", cols.String(), "%s", filter, sort, cfg_filter.Limit, page)
}

func NewFilter(cfg FilterConfig, model *FilterModel, table_model TableModel) *fyne.Container {
	num_elements := cfg.num_components()
	if num_elements == 0 {
		return FContainer.NewBorder(FWidget.NewLabel("Table Filter"), nil, nil, nil, FWidget.NewLabel("No filter options enabled."))
	}

	var content fyne.CanvasObject
	num_cols_visible := 0
	last_checked := -1
	cols_visible := make([]string, 0)
	bind_cols_visible := FBind.BindStringList(&cols_visible)
	components := make([]fyne.CanvasObject, 0, num_elements*2)
	canvas := fyne.CurrentApp().Driver().AllWindows()[0].Canvas()

	// set number of visible columns and create items
	create_items := func(arr *int) []string {
		*arr = 0
		for _, visible := range model.ColsVisible {
			if visible {
				*arr++
			}
		}
		avail_items := make([]string, 0, *arr)
		for i, col := range table_model.Columns {
			if model.ColsVisible[i] {
				avail_items = append(avail_items, col)
			}
		}
		return avail_items
	}

	if cfg.Table {

	}

	if cfg.Columns {
		longest_col_name := ""
		for _, col := range table_model.Columns {
			if len(col) > len(longest_col_name) {
				longest_col_name = col
			}
		}

		pop_visible_cols := FWidget.NewPopUp(nil, canvas)
		list_cols := FWidget.NewList(
			func() int {
				return len(table_model.Columns)
			},
			func() fyne.CanvasObject {
				return FWidget.NewCheck(longest_col_name, func(b bool) {})
			},
			func(i int, o fyne.CanvasObject) {
				o.(*FWidget.Check).SetText(table_model.Columns[i])
				o.(*FWidget.Check).SetChecked(model.ColsVisible[i])
				o.(*FWidget.Check).OnChanged = func(b bool) {
					model.ColsVisible[i] = b
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
					Height: float32(len(table_model.Columns))*list_cols.MinSize().Height + float32(len(table_model.Columns)-1)*6,
				},
				utils.FyneToDimensions(content.Size().Subtract(fyne.NewSize(0, btn_pop_visible_cols.Position().Y+btn_pop_visible_cols.Size().Height))),
			).ToFyneSize())
			pop_visible_cols.ShowAtPosition(btnPos.Add(fyne.NewPos(0, btn_pop_visible_cols.Size().Height)))
		})

		components = append(components, FWidget.NewLabel("Columns:"), btn_pop_visible_cols)
	}

	if cfg.SortBy {
		select_col := FWidget.NewSelect([]string{}, func(selected string) {
			for i, col := range table_model.Columns {
				if col == selected {
					model.Sort.ByCol = i
					break
				}
			}
		})
		select_col.PlaceHolder = "Select Column"
		select_col.Alignment = fyne.TextAlignCenter
		if len(select_col.Options) > 0 {
			select_col.Selected = select_col.Options[model.Sort.ByCol]
		}
		bind_cols_visible.AddListener(FBind.NewDataListener(func() {
			select_col.Options = cols_visible
			if last_checked >= 0 && select_col.Selected == table_model.Columns[last_checked] && !model.ColsVisible[last_checked] && num_cols_visible > 0 {
				select_col.SetSelectedIndex(0)
			}
			if num_cols_visible == 0 {
				select_col.ClearSelected()
			}
			select_col.Refresh()
		}))

		components = append(components, FWidget.NewLabel("Sort by Column:"), select_col)
	}

	if cfg.SortDir {
		var btn_sort_dir *FWidget.Button
		btn_sort_dir = FWidget.NewButton("Asc", func() {
			if model.Sort.Asc {
				model.Sort.Asc = false
				btn_sort_dir.SetText("Desc")
			} else {
				model.Sort.Asc = true
				btn_sort_dir.SetText("Asc")
			}
		})

		components = append(components, FWidget.NewLabel("Sort Direction:"), btn_sort_dir)
	}

	if cfg.GroupBy {
		select_col_group := FWidget.NewSelect([]string{}, func(selected string) {
			for i, col := range table_model.Columns {
				if col == selected {
					model.GroupByCol = i
					break
				}
			}
		})
		select_col_group.PlaceHolder = "Select Column"
		select_col_group.Alignment = fyne.TextAlignCenter
		if len(select_col_group.Options) > 0 {
			select_col_group.Selected = select_col_group.Options[model.GroupByCol]
		}
		bind_cols_visible.AddListener(FBind.NewDataListener(func() {
			select_col_group.Options = cols_visible
			if last_checked >= 0 && select_col_group.Selected == table_model.Columns[last_checked] && !model.ColsVisible[last_checked] && num_cols_visible > 0 {
				select_col_group.SetSelectedIndex(0)
			}
			if num_cols_visible == 0 {
				select_col_group.ClearSelected()
			}
			select_col_group.Refresh()
		}))

		components = append(components, FWidget.NewLabel("Group by Column:"), select_col_group)
	}

	if cfg.Filter {

	}

	if cfg.Limit {
		entry_num_rows := NewNumericalEntry(false)
		entry_num_rows.SetText(strconv.Itoa(model.Limit))
		entry_num_rows.OnChanged = func(val string) {
			if val == "" {
				// entry_num_rows.SetText(strconv.Itoa(table.Filter.Limit))
				return
			}
			if ival, err := strconv.Atoi(val); err == nil {
				model.Limit = ival
			}
		}
		// entry_num_rows.OnChanged = entry_num_rows.OnSubmitted
		entry_num_rows.OnFocusLost = func() {
			if entry_num_rows.Text != "" {
				model.Limit, _ = strconv.Atoi(entry_num_rows.Text)
			} else {
				entry_num_rows.SetText(strconv.Itoa(model.Limit))
			}
		}

		components = append(components, FWidget.NewLabel("Row Limit:"), entry_num_rows)
	}

	if cfg.Page {
		entry_page := NewNumericalEntry(false)
		entry_page.SetText(strconv.Itoa(model.Page))
		entry_page.OnChanged = func(val string) {
			if val == "" {
				// entry_page.SetText(strconv.Itoa(model.Page))
				return
			}
			if ival, err := strconv.Atoi(val); err == nil {
				if ival < 1 {
					entry_page.SetText("1")
					return
				}
				model.Page = ival
			}
		}
		entry_page.OnFocusLost = func() {
			if entry_page.Text != "" {
				model.Page, _ = strconv.Atoi(entry_page.Text)
			} else {
				entry_page.SetText(strconv.Itoa(model.Page))
			}
		}

		components = append(components, FWidget.NewLabel("Page:"), entry_page)
	}

	content = FContainer.New(FLayout.NewFormLayout(), components...)

	confirm := FWidget.NewButton("Done", func() {
		fmt.Println(GenQuery(table_model, *model))
	})
	confirm.Importance = FWidget.LowImportance
	confirm.Refresh()

	return FContainer.NewBorder(FWidget.NewLabel("Filter:"), confirm, nil, nil, content)
}

func NewTable(cfg TableConfig) *fyne.Container {
	data := [][]string{
		{"C1R1", "C2R1"},
		{"C1R2", "C2R2"},
		{"C1R3", "C2R3"},
		{"C1R4", "C2R4"},
		{"C1R5", "C2R5"},
		{"C1R6", "C2R6"},
		{"C1R7", "C2R7"},
		{"C1R8", "C2R8"},
		{"C1R9", "C2R9"},
		{"C1R10", "C2R10"},
		{"C1R1", "C2R1"},
		{"C1R2", "C2R2"},
		{"C1R3", "C2R3"},
		{"C1R4", "C2R4"},
		{"C1R5", "C2R5"},
		{"C1R6", "C2R6"},
		{"C1R7", "C2R7"},
		{"C1R8", "C2R8"},
		{"C1R9", "C2R9"},
		{"C1R10", "C2R10"},
	}
	table_model := TableModel{
		Columns: []string{"C1", "C2"},
		Rows:    make([]int, len(data)),
	}
	var table *FWidget.Table = nil
	var title *FWidget.Label = nil
	var filter *fyne.Container = nil

	if cfg.Editable {
		dirtyRows := make([]int, 0)
		table = NewPrimitiveEditableTable(&data, &dirtyRows)
	} else {
		table = NewPrimitiveTable(data)
	}

	if cfg.Enabled.Title.Enabled {
		if cfg.Enabled.Title.Text != "" {
			title = FWidget.NewLabel(cfg.Enabled.Title.Text)
		} else {
			title = FWidget.NewLabel("Undefined Title")
		}
		switch cfg.Enabled.Title.Alignment {
		case "left":
			title.Alignment = fyne.TextAlignLeading
		case "center":
			title.Alignment = fyne.TextAlignCenter
		case "right":
			title.Alignment = fyne.TextAlignTrailing
		}
		title.TextStyle = fyne.TextStyle{
			Bold:   cfg.Enabled.Title.Style.Bold,
			Italic: cfg.Enabled.Title.Style.Italic,
		}
		title.Refresh()
	}

	if cfg.Enabled.Headers.Enabled {
		if cfg.Enabled.Headers.Column {
			table.ShowHeaderRow = true
		}
		if cfg.Enabled.Headers.Row {
			table.ShowHeaderColumn = true
		}
	}

	if cfg.Enabled.Filter.Enabled {
		filter_model := NewFilterModel(len(table_model.Columns))
		filter = NewFilter(cfg.Enabled.Filter, &filter_model, table_model)
	}

	cellSize := table.MinSize()
	table.Resize(fyne.NewSize(
		cellSize.Width*float32(len(data[0])),
		cellSize.Height*float32(len(data)),
	))

	var sidePanel *fyne.Container
	var sidePanelContent *fyne.Container
	var filter_toggle *FWidget.Button
	current_visibility := true
	actions := make([]fyne.CanvasObject, 1)
	filter_toggle = FWidget.NewButtonWithIcon("", FTheme.Icon(utils.IconFromName("visibility")), func() {
		current_visibility = !current_visibility
		if current_visibility {
			filter_toggle.SetIcon(FTheme.Icon(utils.IconFromName("visibility_off")))
			sidePanelContent.Show()
		} else {
			filter_toggle.SetIcon(FTheme.Icon(utils.IconFromName("visibility")))
			sidePanelContent.Hide()
		}
		sidePanel.Refresh()
	})
	actions[0] = filter_toggle

	sidePanelContent = FContainer.NewHBox(FWidget.NewSeparator(), filter)
	sidePanel = FContainer.NewBorder(nil, nil, nil, FContainer.NewHBox(FWidget.NewSeparator(), FContainer.NewVBox(actions...)), sidePanelContent)

	return FContainer.NewBorder(nil, nil, nil, sidePanel, FContainer.NewBorder(title, nil, nil, nil, table))
}
