package content

import (
	FContainer "fyne.io/fyne/v2/container"
)

type Tab interface {
	Init()
	CreateContent() *FContainer.TabItem
}

type TabCore struct {
	Tabs []Tab
}
