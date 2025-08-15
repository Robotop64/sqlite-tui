package lua

import lua "github.com/yuin/gopher-lua"

type Widget interface {
	// SetData(data interface{})
	// SetActions(actions []interface{})
}


func registerWidgets() {
	registerWidget("WContainer")
	registerWidget("LBox")
	registerWidget("LFill")
	registerWidget("LBBox")
	registerWidget("LWBox")

	registerWidget("WTable")
	registerWidget("WFilter")
	registerWidget("WView")
	registerWidget("WButton")
	registerWidget("WCheckList")
	registerWidget("WLabel")
}

func registerWidget(name string) {
	Env.SetGlobal(name, Env.NewFunction(func(L *lua.LState) int {
		properties := L.CheckTable(1)

		out := L.NewTable()
		out.RawSetString("WType", lua.LString(name))

		properties.ForEach(func(k, v lua.LValue) {
			out.RawSet(k, v)
		})

		L.Push(out)
		return 1
	}))
}
