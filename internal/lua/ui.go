package lua

import (
	"fmt"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FLayout "fyne.io/fyne/v2/layout"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"

	CLayout "SQLite-GUI/internal/ui/layout"
	CWidget "SQLite-GUI/internal/ui/widgets"
)

func buildLayout(L *lua.LState, widgetTable *lua.LTable) fyne.CanvasObject {
	return FContainer.NewBorder(nil, nil, nil, nil, buildComponent(L, widgetTable))
}

func buildComponent(L *lua.LState, widgetTable *lua.LTable) fyne.CanvasObject {
	var component fyne.CanvasObject
	var err_msg string
	widgetType := widgetTable.RawGetString("WType").String()

WidgetSwitch:
	switch widgetType {
	case "LBox", "LBBox", "LFill", "LWBox":

		dir, ok := widgetTable.RawGetString("dir").(lua.LString)
		if !ok && widgetType != "LFill" {
			err_msg = fmt.Sprintf("Layout '%s' requires a 'dir' property", widgetType)
			break
		}
		switch widgetType {
		case "LBox":
			if dir == "vertical" {
				component = FContainer.NewVBox()
			} else {
				component = FContainer.NewHBox()
			}
		case "LBBox":
			if dir == "vertical" {
				component = FContainer.New(&CLayout.BVBox{})
			} else {
				component = FContainer.New(&CLayout.BHBox{})
			}
		case "LFill":
			component = FContainer.New(&CLayout.Fill{})
		case "LWBox":
			weights := make([]float32, 0)
			if weightsTable, ok := widgetTable.RawGetString("weights").(*lua.LTable); ok {
				weights = make([]float32, 0, weightsTable.Len())
				weightsTable.ForEach(func(k, v lua.LValue) {
					if weight, ok := v.(lua.LNumber); ok {
						weights = append(weights, float32(weight))
					}
				})
			} else {
				err_msg = fmt.Sprintf("Layout '%s' requires a 'weights' property", widgetType)
				break WidgetSwitch
			}
			component = FContainer.New(&CLayout.WBox{Weights: weights, Dir: CLayout.DirFromStr(dir.String())})
		}

		if err := fillContainer(L, component.(*fyne.Container), widgetTable); err != nil {
			fmt.Println("Error filling container:", err)
		}
		return component

	case "WTable":
		// idx_source, ok := widgetTable.RawGetString("idx_source").(lua.LNumber)
		// if !ok {
		// 	err_msg = "A table requires a source!\nIt can be defined with the 'idx_source' property.\nThis index is given by the index (starting with 1) of the sources registered in the current target."
		// 	break
		// }
		// curr_profile := persistent.Data.Profiles.LastProfileUsed
		// curr_target := persistent.Data.Profiles.LastTargetUsed
		// sourcepath := persistent.Profiles[curr_profile].Targets[curr_target].ScriptPaths[int(idx_source)-1]

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
		var table *FWidget.Table
		if cfg_tbl, ok := widgetTable.RawGetString("editable").(lua.LBool); ok && cfg_tbl == lua.LTrue {
			dirtyRows := make([]int, 0)
			table = CWidget.NewEditableTable(&data, &dirtyRows)
		} else {
			table = CWidget.NewTable(data)
		}
		if cfg_header, ok := widgetTable.RawGetString("header").(*lua.LTable); ok {
			if cfg_cols, ok := cfg_header.RawGetString("column").(lua.LBool); ok && cfg_cols == lua.LTrue {
				table.ShowHeaderRow = true
			}
			if cfg_rows, ok := cfg_header.RawGetString("row").(lua.LBool); ok && cfg_rows == lua.LTrue {
				table.ShowHeaderColumn = true
			}
		}

		cellSize := table.MinSize()
		table.Resize(fyne.NewSize(
			cellSize.Width*float32(len(data[0])),
			cellSize.Height*float32(len(data)),
		))

		if cfg_title, ok := widgetTable.RawGetString("title").(*lua.LTable); ok {
			if title_text, ok := cfg_title.RawGetString("text").(lua.LString); ok {
				lbl := FWidget.NewLabel(title_text.String())
				if title_alignment, ok := cfg_title.RawGetString("alignment").(lua.LString); ok {
					switch title_alignment.String() {
					case "left":
						lbl.Alignment = fyne.TextAlignLeading
					case "center":
						lbl.Alignment = fyne.TextAlignCenter
					case "right":
						lbl.Alignment = fyne.TextAlignTrailing
					}
				}
				if title_style, ok := cfg_title.RawGetString("style").(*lua.LTable); ok {

					var bold, italic bool
					if boldVal, ok := title_style.RawGetString("bold").(lua.LBool); ok {
						bold = boldVal == lua.LTrue
					}
					if italicVal, ok := title_style.RawGetString("italic").(lua.LBool); ok {
						italic = italicVal == lua.LTrue
					}

					lbl.TextStyle = fyne.TextStyle{
						Bold:   bold,
						Italic: italic,
					}
					lbl.Refresh()
				}

				component = FContainer.NewBorder(
					lbl, nil, nil, nil, table,
				)
			}
		} else {
			component = table
		}
	case "WFilter":
		if cfg_type, ok := widgetTable.RawGetString("type").(lua.LString); ok {
			switch cfg_type {
			case "table":
				table := CWidget.NewTableModel([]string{"C1", "C2", "C3", "C4", "C5", "C6"})
				component = CWidget.NewFilter_Table(&table)
			default:
				err_msg = fmt.Sprintf("Filter widget type '%s' is not supported.", cfg_type.String())
				break WidgetSwitch
			}
		} else {
			err_msg = "Filter widget requires a 'type' property to be defined."
			break
		}
	case "WSeparator":
		component = FWidget.NewSeparator()
	case "WSpacer":
		component = FLayout.NewSpacer()
	case "WButton":
		text := widgetTable.RawGetString("text").String()
		component = FWidget.NewButton(text, func() {
			if action, ok := widgetTable.RawGetString("action").(*lua.LFunction); ok {
				if err := L.CallByParam(lua.P{
					Fn:      action,
					NRet:    0,
					Protect: true,
				}); err != nil {
					fmt.Println("Error calling button action:", err)
				}
			} else {
				fmt.Println("The pressed button has no assigned action to be performed.")
			}
		})

	case "WLabel":
		component = FWidget.NewLabel(widgetTable.RawGetString("text").String())

	case "WCheckList":
		component = FWidget.NewLabel("Checklist placeholder")
	case "WView":
		component = FWidget.NewLabel(widgetTable.RawGetString("text").String())
	}

	if err_msg != "" {
		return FWidget.NewLabel(err_msg)
	}

	return component
}

func fillContainer(L *lua.LState, container *fyne.Container, widgetTable *lua.LTable) error {
	var childrenTbl *lua.LTable
	if children := widgetTable.RawGetString("children"); children.Type() == lua.LTTable {
		childrenTbl = children.(*lua.LTable)
	} else if children := widgetTable.RawGetInt(1); children.Type() == lua.LTTable {
		childrenTbl = children.(*lua.LTable)
	}
	if childrenTbl != nil {
		if childrenTbl.Len() == 0 {
			return fmt.Errorf("children table is empty")
		}
		items := make([]fyne.CanvasObject, 0, childrenTbl.Len())
		childrenTbl.ForEach(func(_, child lua.LValue) {
			if childTbl, ok := child.(*lua.LTable); ok {
				childWidget := buildComponent(L, childTbl)
				items = append(items, childWidget)
			}
		})
		container.Objects = items
		container.Refresh()
	} else {
		return fmt.Errorf("no children table found")
	}

	return nil
}
