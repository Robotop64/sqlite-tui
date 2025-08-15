package widgets

import (
	FWidget "fyne.io/fyne/v2/widget"
)

type Entry struct {
	FWidget.Entry
	OnFocusLost func()
}

func NewEntry() *Entry {
	entry := &Entry{}
	entry.ExtendBaseWidget(entry)

	return entry
}

func (e *Entry) FocusLost() {
	e.Entry.FocusLost()
	if e.OnFocusLost != nil {
		e.OnFocusLost()
	}
}
