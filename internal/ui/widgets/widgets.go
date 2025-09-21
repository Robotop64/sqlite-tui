package widgets

import (
	FBind "fyne.io/fyne/v2/data/binding"
	FWidget "fyne.io/fyne/v2/widget"
)

func NonValidatedEntry() *FWidget.Entry {
	entry := FWidget.NewEntry()
	entry.Validator = nil
	entry.Refresh()
	return entry
}

func NonValidatedEntryWithData(data FBind.String) *FWidget.Entry {
	entry := FWidget.NewEntryWithData(data)
	entry.Validator = nil
	entry.Refresh()
	return entry
}
