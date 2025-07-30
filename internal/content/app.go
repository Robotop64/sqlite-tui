package content

import (
	"fyne.io/fyne/v2"
	FApp "fyne.io/fyne/v2/app"
	FContainer "fyne.io/fyne/v2/container"

	"SQLite-GUI/internal/utils"
)

var AppHandle *fyne.App
var WindowHandle *fyne.Window

func Init() {
	app := FApp.New()
	AppHandle = &app
	window := app.NewWindow("SQLite-GUI")
	WindowHandle = &window

	tabCore := &TabCore{}
	tabCore.Tabs = append(tabCore.Tabs, &ProfileTab{})

	for _, tab := range tabCore.Tabs {
		tab.Init()
	}

	tabs := FContainer.NewAppTabs(
		utils.Map(tabCore.Tabs, func(i int, tab Tab) *FContainer.TabItem {
			return tab.GenerateContent()
		})...,
	)

	tabs.SetTabLocation(FContainer.TabLocationTop)

	window.SetContent(tabs)
	window.Resize(fyne.NewSize(800, 600))
	window.SetMaster()
	window.ShowAndRun()
}
