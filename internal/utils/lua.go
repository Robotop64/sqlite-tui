package utils

import lua "github.com/yuin/gopher-lua"

// check if the given lua value is equal to a go value
func CheckVal(luaVal lua.LValue, goVal interface{}) bool {
	switch v := goVal.(type) {
	case int:
		if l, ok := luaVal.(lua.LNumber); ok {
			return int(l) == v
		}
	case string:
		if l, ok := luaVal.(lua.LString); ok {
			return string(l) == v
		}
	case bool:
		if l, ok := luaVal.(lua.LBool); ok {
			return bool(l) == v
		}
	}
	return false
}
