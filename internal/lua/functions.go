package lua

import funcs "SQLite-GUI/internal/lua/functions"

func registerFunctions() {
	funcs.NewLink(Env)
}
