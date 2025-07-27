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

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	"SQLite-GUI/internal/persistent"
)

func main() {
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

	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
