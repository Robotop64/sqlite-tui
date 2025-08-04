package content

import (
	// "fmt"
	// "path/filepath"

	// "fyne.io/fyne/v2"
	// FBind "fyne.io/fyne/v2/data/binding"
	// FDialog "fyne.io/fyne/v2/dialog"
	// FTheme "fyne.io/fyne/v2/theme"

	"fmt"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FWidget "fyne.io/fyne/v2/widget"

	// "SQLite-GUI/internal/utils"
	"SQLite-GUI/internal/persistent"
	ui "SQLite-GUI/internal/ui"
	"SQLite-GUI/internal/utils"
)

type BrowserTab struct {
	selected_profile int
	selected_target  int
	scripts          []persistent.Script
	components       struct {
		tab_targets *fyne.Container
		tab_views   *fyne.Container
		content     *fyne.Container
	}
	bindings struct {
	}
}

func (t *BrowserTab) Init() {
	t.scripts = []persistent.Script{}

	t.components.tab_targets = FContainer.New(&ui.Fill{}, FWidget.NewLabel("No targets"))
	t.components.tab_views = FContainer.New(&ui.Fill{}, FWidget.NewLabel("No views"))
	t.components.content = FContainer.NewStack()

	t.Update()
}

func (t *BrowserTab) Update() {
	t.selected_profile = persistent.Data.Profiles.LastProfileUsed
	if len(persistent.Profiles[t.selected_profile].Targets) > persistent.Data.Profiles.LastTargetUsed {
		t.selected_target = persistent.Data.Profiles.LastTargetUsed
		t.components.tab_targets.Objects = []fyne.CanvasObject{createTargetButtons(t)}
		t.components.tab_targets.Refresh()
		t.components.tab_views.Objects = []fyne.CanvasObject{createViewButtons(t)}
		t.components.tab_views.Refresh()
	} else {
		t.selected_target = 0
	}
}

func (t *BrowserTab) CreateContent() *FContainer.TabItem {
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
	if len(targets) == 0 {
		return FWidget.NewLabel("No targets")
	}

	list := FWidget.NewList(
		func() int {
			return len(targets)
		},
		func() fyne.CanvasObject {
			return FWidget.NewLabel("Target")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*FWidget.Label).SetText(targets[i].Name)
		},
	)
	list.OnSelected = func(id int) {
		t.selected_target = id
		target := targets[t.selected_target]
		t.scripts = utils.Map(target.ScriptPaths, func(i int, path string) persistent.Script {
			script, err := persistent.LoadScript(path)
			if err != nil {
				fmt.Printf("Error loading script from %s: %v\n", path, err)
				return persistent.Script{}
			}
			return script
		})
		t.components.tab_views.Objects = []fyne.CanvasObject{createViewButtons(t)}
		t.components.tab_views.Refresh()
	}
	if len(targets) > 0 && t.selected_target < len(targets) {
		list.Select(t.selected_target)
	}
	list.Resize(fyne.NewSize(list.MinSize().Width, float32(len(targets))*list.MinSize().Height+float32(len(targets)-1)*4))
	return list
}

func createViewButtons(t *BrowserTab) fyne.CanvasObject {
	views := []*persistent.Script{}
	for _, script := range t.scripts {
		if script.MetaData.Type == persistent.SCRIPT_VIEW {
			views = append(views, &script)
		}
	}
	if len(views) == 0 {
		return FWidget.NewLabel("No views")
	}
	var longestName string
	for _, view := range views {
		if len(view.MetaData.Name) > len(longestName) {
			longestName = view.MetaData.Name
		}
	}
	list := FWidget.NewList(
		func() int {
			return len(views)
		},
		func() fyne.CanvasObject {
			return FWidget.NewLabel(longestName)
		},
		func(i int, o fyne.CanvasObject) {
			o.(*FWidget.Label).SetText(views[i].MetaData.Name)
		},
	)
	list.OnSelected = func(id int) {
	}
	list.Resize(fyne.NewSize(list.MinSize().Width, float32(len(views))*list.MinSize().Height+float32(len(views)-1)*4))
	return list
}
