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

func LoadScript(script persistent.Script) error {
	Env.SetGlobal("load_sources", lua.LNil)
	Env.SetGlobal("layout", lua.LNil)

	if err := Env.DoString(string(script.Script)); err != nil {
		return fmt.Errorf("failed to load Lua script: %w", err)
	}
	return nil
}

func LoadView() (fyne.CanvasObject, error) {
	var lua_layout *lua.LTable

	if layout := Env.GetGlobal("layout"); layout.Type() == lua.LTTable {
		lua_layout = layout.(*lua.LTable)
	} else {
		return FWidget.NewLabel("The selected view does not contain a layout."), fmt.Errorf("layout not found in Lua script")
	}

	return buildLayout(Env, lua_layout), nil
}

func LoadSources() error {
	persistent.Sources = make([]*persistent.Source, 0)

	var lua_sources *lua.LTable
	if sources := Env.GetGlobal("load_sources"); sources.Type() == lua.LTTable {
		lua_sources = sources.(*lua.LTable)
	} else {
		return fmt.Errorf("load_sources not found in Lua script")
	}

	for i := 1; i <= lua_sources.Len(); i++ {
		target := persistent.ActiveTarget()
		idx_source := int(lua.LVAsNumber(lua_sources.RawGetInt(i)))
		if _, err := persistent.NewSource(target.SourcePaths[idx_source]); err != nil {
			return fmt.Errorf("failed to create source from path '%s': %w", target.SourcePaths[idx_source-1], err)
		}
	}

	fmt.Printf("Loaded %d source(s) from Lua script.\n", len(persistent.Sources))

	return nil
}
