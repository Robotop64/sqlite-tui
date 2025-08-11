package lua

import (
	"SQLite-GUI/internal/persistent"
	"fmt"

	"fyne.io/fyne/v2"
	FWidget "fyne.io/fyne/v2/widget"
	lua "github.com/yuin/gopher-lua"
)

var Env *lua.LState

func Init() {
	Env = lua.NewState()
	if Env == nil {
		panic("failed to create Lua state")
	}

	// Load standard libraries
	// if err := Env.DoString(`
	// 	require("os")
	// 	require("io")
	// 	require("math")
	// 	require("string")
	// 	require("table")
	// `); err != nil {
	// 	panic(err)
	// }

	registerWidgets()
}

func Clean() {
	if Env != nil {
		Env.Close()
	}
}

func LoadView(script persistent.Script) (fyne.CanvasObject, error) {

	Env.SetGlobal("layout", lua.LNil)

	if err := Env.DoString(string(script.Script)); err != nil {
		return FWidget.NewLabel("Failed to load the script of the selected view."), fmt.Errorf("failed to load Lua script: %w", err)
	}

	var lua_layout *lua.LTable

	if layout := Env.GetGlobal("layout"); layout.Type() == lua.LTTable {
		lua_layout = layout.(*lua.LTable)
	} else {
		return FWidget.NewLabel("The selected view does not contain a layout."), fmt.Errorf("layout not found in Lua script")
	}

	return buildLayout(Env, lua_layout), nil
}
