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
	CLayout "SQLite-GUI/internal/ui/layout"
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

	t.components.tab_targets = FContainer.New(&CLayout.Fill{}, FWidget.NewLabel("No targets"))
	t.components.tab_views = FContainer.New(&CLayout.Fill{}, FWidget.NewLabel("No views"))
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
		t.components.content.Objects = []fyne.CanvasObject{FWidget.NewLabel("Select a view to display its layout")}
		t.components.content.Refresh()
	} else {
		setTarget(t, 0)
	}
}

func (t *BrowserTab) CreateContent() *FContainer.TabItem {
	sidebar := FContainer.NewAppTabs(
		FContainer.NewTabItem("Targets", t.components.tab_targets),
		FContainer.NewTabItem("Views", t.components.tab_views),
	)

	return FContainer.NewTabItem("Browser",
		FContainer.NewBorder(
			nil, nil, FContainer.NewHBox(FContainer.New(&CLayout.MinVBox{MinWidth: 130}, sidebar), FWidget.NewSeparator()), nil,
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
		setTarget(t, id)
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

	setContent := func(script persistent.Script) error {
		if err := lua.LoadScript(script); err != nil {
			t.components.content.Objects = []fyne.CanvasObject{FWidget.NewLabel(fmt.Sprintf("Failed to load the script of the selected view.\nCaused by error:\n\t %v", err))}
			t.components.content.Refresh()
			return fmt.Errorf("failed to load Lua script \"%s\":\n  %w", script.MetaData.Name, err)
		}

		if err := lua.LoadSources(); err != nil {
			fmt.Println("Error loading sources for script", script.MetaData.Name, ":", err)
		}

		if layout, err := lua.LoadView(); err == nil {
			t.components.content.Objects = []fyne.CanvasObject{layout}
			t.components.content.Refresh()
		} else {
			if layout != nil {
				t.components.content.Objects = []fyne.CanvasObject{layout}
				t.components.content.Refresh()
			}
			return fmt.Errorf("failed to load view \"%s\":\n  %w", script.MetaData.Name, err)
		}
		return nil
	}

	ref_btns := make([]*FWidget.Button, len(views))
	var ref_btn_cur int
	list := FWidget.NewList(
		func() int {
			return len(views)
		},
		func() fyne.CanvasObject {
			label := FWidget.NewLabel(longestName)
			btn := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameMediaReplay), func() {})
			btn.Importance = FWidget.LowImportance
			btn.Hide()
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
			ref_btns[i] = button
			button.OnTapped = func() {
				path := t.scripts.paths[i]
				if script, err := persistent.LoadScript(path); err != nil {
					fmt.Println("Error re-loading script\"", script.MetaData.Name, "\":\n  ", err)
				} else {
					t.scripts.items[i] = script
					*views[i] = script
					if err := setContent(script); err != nil {
						fmt.Println("Error setting content for script\"", script.MetaData.Name, "\":\n  ", err)
					}
					fmt.Println("Re-loaded script:", script.MetaData.Name)
				}

			}
		},
	)

	list.OnSelected = func(id int) {
		t.selection.view = id
		setContent(*views[id])
		ref_btns[ref_btn_cur].Hide()
		ref_btn_cur = id
		ref_btns[id].Show()
	}
	list.Resize(fyne.NewSize(list.MinSize().Width, float32(len(views))*list.MinSize().Height+float32(len(views)-1)*4))
	return list
}

func setTarget(t *BrowserTab, target int) {
	t.selection.target = target
	persistent.Data.Profiles.LastTargetUsed = target
}
