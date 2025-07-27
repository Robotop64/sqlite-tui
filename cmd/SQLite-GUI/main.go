// package main

// import (
// 	"fmt"
// 	"os"

// 	persistent "github.com/Robotop64/sqlite-tui/internal/persistent"
// 	tabs "github.com/Robotop64/sqlite-tui/internal/tabs"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// func main() {
// 	// load main config
// 	if err := persistent.LoadConfig(); err != nil {
// 		fmt.Println("Error loading config:", err)
// 		os.Exit(1)
// 	}
// 	if err := persistent.LoadData(); err != nil {
// 		fmt.Println("Error loading data:", err)
// 		os.Exit(1)
// 	}
// 	// load profiles
// 	persistent.LoadProfiles()

// 	// initialize the core model
// 	c := tabs.Init()

// 	p := tea.NewProgram(c, tea.WithAltScreen(), tea.WithMouseAllMotion())

// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }

package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	FApp "fyne.io/fyne/v2/app"
	FContainer "fyne.io/fyne/v2/container"

	"SQLite-GUI/internal/content"
	"SQLite-GUI/internal/persistent"
	"SQLite-GUI/internal/utils"
)

func main() {
	setup()

	app := FApp.New()
	window := app.NewWindow("SQLite-GUI")

	tabCore := &content.TabCore{}
	tabCore.Tabs = append(tabCore.Tabs, &content.ProfileTab{})

	for _, tab := range tabCore.Tabs {
		tab.Init()
	}

	tabs := FContainer.NewAppTabs(
		utils.Map(tabCore.Tabs, func(i int, tab content.Tab) *FContainer.TabItem {
			return tab.GenerateContent()
		})...,
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(FContainer.TabLocationTop)

	window.SetContent(tabs)
	window.Resize(fyne.NewSize(800, 600))
	window.SetMaster()
	window.ShowAndRun()

	shutdown()
}

func setup() {
	// load main config
	if err := persistent.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	if err := persistent.LoadData(); err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}
	// load profiles
	persistent.LoadProfiles()
}

func shutdown() {

}
