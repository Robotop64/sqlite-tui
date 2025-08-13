package lua

import (
	"SQLite-GUI/internal/ui"
	"fmt"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"
)

func buildLayout(L *lua.LState, widgetTable *lua.LTable) fyne.CanvasObject {
	return FContainer.NewBorder(nil, nil, nil, nil, buildComponent(L, widgetTable))
}

func buildComponent(L *lua.LState, widgetTable *lua.LTable) fyne.CanvasObject {
	var component fyne.CanvasObject
	var err_msg string
	widgetType := widgetTable.RawGetString("type").String()

	switch widgetType {
	case "LBox", "LBBox", "LFill":

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
				component = FContainer.New(&ui.BVBox{})
			} else {
				component = FContainer.New(&ui.BHBox{})
			}
		case "LFill":
			component = FContainer.New(&ui.Fill{})
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
			[]string{"top left", "top right"},
			[]string{"bottom left", "bottom right"},
		}
		dirtyRows := make([]int, 0)
		table := ui.EditableTable(&data, &dirtyRows)

		cellSize := table.MinSize()
		table.Resize(fyne.NewSize(
			cellSize.Width*float32(len(data[0])),
			cellSize.Height*float32(len(data)),
		))
		if cfg_header, ok := widgetTable.RawGetString("header").(*lua.LTable); ok {
			if cfg_cols, ok := cfg_header.RawGetString("column").(lua.LBool); ok && cfg_cols == lua.LTrue {
				table.ShowHeaderRow = true
			}
			if cfg_rows, ok := cfg_header.RawGetString("row").(lua.LBool); ok && cfg_rows == lua.LTrue {
				table.ShowHeaderColumn = true
			}
		}

		component = table
	case "WFilter":
		component = FWidget.NewLabel("Filter placeholder")
	case "WCheckList":
		component = FWidget.NewLabel("Checklist placeholder")
	case "WView":
		component = FWidget.NewLabel(widgetTable.RawGetString("text").String())
	case "WButton":
		text := widgetTable.RawGetString("text").String()
		action := widgetTable.RawGetString("action")
		component = FWidget.NewButton(text, func() {
			if err := L.CallByParam(lua.P{
				Fn:      action,
				NRet:    0,
				Protect: true,
			}); err != nil {
				fmt.Println("Error calling button action:", err)
			}
		})
	case "WLabel":
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
		items := make([]fyne.CanvasObject, 0)
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
