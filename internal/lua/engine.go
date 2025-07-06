package lua

import (
	lua "github.com/yuin/gopher-lua"
)

var Env *lua.LState

func Init() {
	Env = lua.NewState()
	if Env == nil {
		panic("failed to create Lua state")
	}

	// Load standard libraries
	if err := Env.DoString(`
		require("os")
		require("io")
		require("math")
		require("string")
		require("table")
	`); err != nil {
		panic(err)
	}
}

func Clean() {
	if Env != nil {
		Env.Close()
	}
}
