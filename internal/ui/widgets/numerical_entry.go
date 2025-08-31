package widgets

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

type NumericalEntry struct {
	Entry
	allow_signed bool
}

func NewNumericalEntry(signed bool) *NumericalEntry {
	entry := &NumericalEntry{allow_signed: signed}
	entry.ExtendBaseWidget(entry)

	return entry
}

func (e *NumericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' || (e.allow_signed && r == '-') {
		e.Entry.TypedRune(r)
	}
}

func (e *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
		if !e.allow_signed {
			e.SetText(strings.ReplaceAll(e.Text, "-", ""))
		}
	}
}
