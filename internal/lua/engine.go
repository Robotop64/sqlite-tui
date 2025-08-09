package lua

import (
	"SQLite-GUI/internal/persistent"
	"fmt"

	"fyne.io/fyne/v2"
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

	if err := Env.DoString(string(script.Script)); err != nil {
		return nil, fmt.Errorf("failed to load Lua script: %w", err)
	}

	lua_layout := Env.GetGlobal("layout").(*lua.LTable)

	return buildLayout(Env, lua_layout), nil
}
