package ui

import (
	FBind "fyne.io/fyne/v2/data/binding"
	FWidget "fyne.io/fyne/v2/widget"
)

func NonValidatedEntry(data FBind.String) *FWidget.Entry {
	entry := FWidget.NewEntryWithData(data)
	entry.Validator = nil
	entry.Refresh()
	return entry
}
