package content

import (
	"fyne.io/fyne/v2"
	FApp "fyne.io/fyne/v2/app"
)

var AppHandle *fyne.App
var WindowHandle *fyne.Window

func Init() {
	app := FApp.NewWithID("SQLite-GUI")
	AppHandle = &app
	window := app.NewWindow("SQLite-GUI")
	WindowHandle = &window

	window.SetContent(InitCore())
	window.Resize(fyne.NewSize(800, 600))
	window.SetMaster()
	window.ShowAndRun()
}
