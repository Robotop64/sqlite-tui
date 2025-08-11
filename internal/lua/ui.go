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
	widgetType := widgetTable.RawGetString("type").String()

	switch widgetType {
	case "LBox", "LBBox", "LFill":
		dir := widgetTable.RawGetString("dir").String()
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
		data := [][]string{
			[]string{"top left", "top right"},
			[]string{"bottom left", "bottom right"},
		}
		component = FWidget.NewTable(
			func() (int, int) {
				return 2, 2
			},
			func() fyne.CanvasObject {
				return FWidget.NewLabel("Placeholder")
			},
			func(i FWidget.TableCellID, o fyne.CanvasObject) {
				o.(*FWidget.Label).SetText(data[i.Row][i.Col])
			},
		)
		cellSize := component.MinSize()
		component.Resize(fyne.NewSize(
			cellSize.Width*float32(len(data[0])),
			cellSize.Height*float32(len(data)),
		))
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
