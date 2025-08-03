package content

import (
	FContainer "fyne.io/fyne/v2/container"

	"SQLite-GUI/internal/utils"
)

type Tab interface {
	Init()
	CreateContent() *FContainer.TabItem
}

type TabCore struct {
	Tabs []Tab
}

func InitCore() *FContainer.AppTabs {
	tabCore := &TabCore{}
	tabCore.Tabs = append(tabCore.Tabs,
		&ProfileTab{},
		&BrowserTab{},
	)

	tabs := FContainer.NewAppTabs(
		utils.Map(tabCore.Tabs, func(i int, tab Tab) *FContainer.TabItem {
			tab.Init()
			return tab.CreateContent()
		})...,
	)

	tabs.SetTabLocation(FContainer.TabLocationTop)

	labels := utils.Map(tabs.Items, func(i int, item *FContainer.TabItem) string {
		return item.Text
	})
	tabs.OnSelected = func(item *FContainer.TabItem) {
		index := utils.IndexOf(labels, item.Text)
		tabs.Items[index].Content = tabCore.Tabs[index].CreateContent().Content
		tabs.Items[index].Content.Refresh()
	}

	return tabs
}
