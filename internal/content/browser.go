package content

import (
	// "fmt"
	// "path/filepath"

	// "fyne.io/fyne/v2"
	// FBind "fyne.io/fyne/v2/data/binding"
	// FDialog "fyne.io/fyne/v2/dialog"
	// FTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FWidget "fyne.io/fyne/v2/widget"

	// "SQLite-GUI/internal/utils"
	"SQLite-GUI/internal/persistent"
	ui "SQLite-GUI/internal/ui"
)

type BrowserTab struct {
	selected_profile int
	selected_target  int
	components       struct {
		tab_targets *fyne.Container
		tab_views   *fyne.Container
		content     *fyne.Container
	}
	bindings struct {
	}
}

func (t *BrowserTab) Init() {
	t.selected_profile = persistent.Data.Profiles.LastProfileUsed
	t.selected_target = persistent.Data.Profiles.LastTargetUsed
}

func (t *BrowserTab) CreateContent() *FContainer.TabItem {
	t.components.tab_targets = FContainer.New(&ui.Fill{}, createTargetButtons(t))
	t.components.tab_views = FContainer.NewVBox()
	t.components.content = FContainer.NewStack()

	sidebar := FContainer.NewAppTabs(
		FContainer.NewTabItem("Targets", t.components.tab_targets),
		FContainer.NewTabItem("Views", t.components.tab_views),
	)

	return FContainer.NewTabItem("Browser",
		FContainer.NewBorder(
			nil, nil, FContainer.NewHBox(FContainer.New(&ui.MinVBox{MinWidth: 130}, sidebar), FWidget.NewSeparator()), nil,
			t.components.content,
		),
	)
}

func createTargetButtons(t *BrowserTab) fyne.CanvasObject {
	targets := persistent.Profiles[t.selected_profile].Targets
	list := FWidget.NewList(
		func() int {
			return len(targets)
		},
		func() fyne.CanvasObject {
			return FWidget.NewLabel("Target     2")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*FWidget.Label).SetText(targets[i].Name)
		},
	)
	list.Resize(fyne.NewSize(130, 400))
	return list
}
