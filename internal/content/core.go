package content

import (
	FContainer "fyne.io/fyne/v2/container"
)

type Tab interface {
	Init()
	GenerateContent() *FContainer.TabItem
}

type TabCore struct {
	Tabs []Tab
}
