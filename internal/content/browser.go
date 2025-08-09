package content

import (
	// "fmt"
	// "path/filepath"

	// "fyne.io/fyne/v2"
	// FBind "fyne.io/fyne/v2/data/binding"
	// FDialog "fyne.io/fyne/v2/dialog"
	// FTheme "fyne.io/fyne/v2/theme"

	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FTheme "fyne.io/fyne/v2/theme"
	FWidget "fyne.io/fyne/v2/widget"

	// "SQLite-GUI/internal/utils"
	lua "SQLite-GUI/internal/lua"
	"SQLite-GUI/internal/persistent"
	ui "SQLite-GUI/internal/ui"
	utils "SQLite-GUI/internal/utils"
)

type BrowserTab struct {
	selection struct {
		profile int
		target  int
		view    int
	}
	scripts struct {
		paths []string
		items []persistent.Script
	}
	components struct {
		tab_targets *fyne.Container
		tab_views   *fyne.Container
		content     *fyne.Container
	}
	bindings struct {
	}
}

func (t *BrowserTab) Init() {
	t.scripts.items = []persistent.Script{}

	t.components.tab_targets = FContainer.New(&ui.Fill{}, FWidget.NewLabel("No targets"))
	t.components.tab_views = FContainer.New(&ui.Fill{}, FWidget.NewLabel("No views"))
	t.components.content = FContainer.NewStack()

	t.Update()

	lua.Init()
}

func (t *BrowserTab) Update() {
	t.selection.profile = persistent.Data.Profiles.LastProfileUsed
	if len(persistent.Profiles[t.selection.profile].Targets) > persistent.Data.Profiles.LastTargetUsed {
		t.selection.target = persistent.Data.Profiles.LastTargetUsed
		t.components.tab_targets.Objects = []fyne.CanvasObject{createTargetButtons(t)}
		t.components.tab_targets.Refresh()
		t.components.tab_views.Objects = []fyne.CanvasObject{createViewButtons(t)}
		t.components.tab_views.Refresh()
	} else {
		t.selection.target = 0
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
	targets := persistent.Profiles[t.selection.profile].Targets
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
		t.selection.target = id
		target := targets[t.selection.target]
		t.scripts.paths = make([]string, len(target.ScriptPaths))
		t.scripts.items = make([]persistent.Script, len(target.ScriptPaths))
		for i := 0; i < len(target.ScriptPaths); i++ {
			//load paths
			path := target.ScriptPaths[i]
			if !filepath.IsAbs(path) {
				path = utils.RelativeToAbsolutePath(filepath.Dir(persistent.Data.Profiles.Paths[t.selection.profile]), path)
			}
			t.scripts.paths[i] = path

			//find views
			if script, err := persistent.LoadScript(path); err != nil {
				fmt.Printf("Error loading script from %s: %v\n", path, err)
			} else {
				t.scripts.items = append(t.scripts.items, script)
			}
		}

		t.components.tab_views.Objects = []fyne.CanvasObject{createViewButtons(t)}
		t.components.tab_views.Refresh()
	}
	if len(targets) > 0 && t.selection.target < len(targets) {
		list.Select(t.selection.target)
	}
	list.Resize(fyne.NewSize(list.MinSize().Width, float32(len(targets))*list.MinSize().Height+float32(len(targets)-1)*4))
	return list
}

func createViewButtons(t *BrowserTab) fyne.CanvasObject {
	views := make([]*persistent.Script, 0, len(t.scripts.items))
	for _, script := range t.scripts.items {
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

	setContent := func(script persistent.Script) {
		if layout, err := lua.LoadView(script); err == nil {
			t.components.content.Objects = []fyne.CanvasObject{layout}
			t.components.content.Refresh()
		} else {
			fmt.Println("Error loading view \"", script.MetaData.Name, "\":\n  ", err)
		}
	}

	list := FWidget.NewList(
		func() int {
			return len(views)
		},
		func() fyne.CanvasObject {
			label := FWidget.NewLabel(longestName)
			btn := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameMediaReplay), func() {})
			btn.Importance = FWidget.LowImportance
			return FContainer.NewBorder(
				nil, nil,
				label, btn,
			)
		},
		func(i int, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			label := c.Objects[0].(*FWidget.Label)
			label.SetText(views[i].MetaData.Name)
			button := c.Objects[1].(*FWidget.Button)
			button.OnTapped = func() {
				path := t.scripts.paths[i]
				if script, err := persistent.LoadScript(path); err != nil {
					fmt.Println("Error re-loading script\"", script.MetaData.Name, "\":\n  ", err)
				} else {
					t.scripts.items[i] = script
					*views[i] = script
					fmt.Println("Re-loaded script:", script.MetaData.Name)
				}

				setContent(*views[i])
			}
		},
	)

	list.OnSelected = func(id int) {
		t.selection.view = id
		setContent(*views[id])
	}
	list.Resize(fyne.NewSize(list.MinSize().Width, float32(len(views))*list.MinSize().Height+float32(len(views)-1)*4))
	return list
}
