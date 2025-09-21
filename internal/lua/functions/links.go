package functions

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type LinkManager struct {
	LinkGroups [][]*any
}

func (m *LinkManager) newGroup() int {
	if m.LinkGroups == nil {
		m.LinkGroups = make([][]*any, 1)
	}
	m.LinkGroups = append(m.LinkGroups, []*any{})
	return len(m.LinkGroups) - 1
}

func (m *LinkManager) NewLink(group int, link *any) error {
	if group < 0 || group >= len(m.LinkGroups) {
		return fmt.Errorf("invalid link group: %d", group)
	}
	m.LinkGroups[group] = append(m.LinkGroups[group], link)
	return nil
}

var Link_Manager LinkManager = LinkManager{}

func NewLink(env *lua.LState) {
	env.SetGlobal("NewLink", env.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(Link_Manager.newGroup()))
		return 1
	}))
}
