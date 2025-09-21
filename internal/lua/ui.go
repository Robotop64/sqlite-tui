package lua

import (
	"fmt"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FLayout "fyne.io/fyne/v2/layout"
	FTheme "fyne.io/fyne/v2/theme"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"

	CLayout "SQLite-GUI/internal/ui/layout"
	CWidget "SQLite-GUI/internal/ui/widgets"
	"SQLite-GUI/internal/utils"
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

		cfg, err := CWidget.TableConfig{}.FromLuaTable(widgetTable)
		component = CWidget.NewTable(cfg)
		if err != nil {
			fmt.Println("Error creating table from config:", err)
		}
	case "WSeparator":
		component = FWidget.NewSeparator()
	case "WSpacer":
		component = FLayout.NewSpacer()
	case "WButton":
		var text, icon string
		var ok_text, ok_icon bool
		var Ltext, Licon lua.LValue
		if Ltext, ok_text = widgetTable.RawGetString("text").(lua.LString); ok_text {
			text = Ltext.String()
		}
		if Licon, ok_icon = widgetTable.RawGetString("icon").(lua.LString); ok_icon {
			icon = Licon.String()
		}

		if !ok_text && !ok_icon {
			err_msg = "A button requires either a 'text' or 'icon' or both property to be defined."
			break
		}

		lambda_action := func() {
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
		}

		if !ok_icon {
			component = FWidget.NewButton(text, lambda_action)
			break
		}

		component = FWidget.NewButtonWithIcon(text, FTheme.Icon(utils.IconFromName(icon)), lambda_action)
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
