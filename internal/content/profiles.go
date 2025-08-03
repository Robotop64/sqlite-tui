package content

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FBind "fyne.io/fyne/v2/data/binding"
	FDialog "fyne.io/fyne/v2/dialog"
	FLayout "fyne.io/fyne/v2/layout"
	FTheme "fyne.io/fyne/v2/theme"
	FWidget "fyne.io/fyne/v2/widget"

	"SQLite-GUI/internal/persistent"
	"SQLite-GUI/internal/utils"
)

type ProfileTab struct {
	selected_profile int
	components       struct {
		list_btn_profiles *fyne.Container
		list_form_targets *fyne.Container
	}
	bindings struct {
		profile          *persistent.Profile
		profile_name     FBind.String
		profile_location FBind.String
		profile_note     FBind.String
	}
}

func (t *ProfileTab) Init() {
	t.selected_profile = persistent.Data.Profiles.LastProfileUsed

	t.bindings.profile_name = FBind.NewString()
	t.bindings.profile_location = FBind.NewString()
	t.bindings.profile_note = FBind.NewString()

	t.components.list_form_targets = FContainer.NewVBox()
}

func (t *ProfileTab) Update() {
}

func (t *ProfileTab) CreateContent() *FContainer.TabItem {
	return FContainer.NewTabItem("Profiles",
		FContainer.NewStack(
			FContainer.NewBorder(
				nil, nil, createProfilePanel(t), nil,
				FContainer.NewVScroll(createEditorForm(t)),
			),
		),
	)
}

func createProfilePanel(t *ProfileTab) *fyne.Container {
	//# list buttons of tracked profiles
	t.components.list_btn_profiles = FContainer.NewVBox(createProfileButtons(t)...)
	buttonScroll := FContainer.NewVScroll(t.components.list_btn_profiles)
	//#

	//# list action / edit buttons
	addButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() { addPopup(t) })
	removeButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func() { removePopup(t) })
	saveButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDocumentSave), func() {})
	saveButton.OnTapped = func() {
		persistent.SaveProfiles()
		updateProfileButtons(t)
	}
	editButtons := FContainer.NewVBox(FWidget.NewSeparator(), addButton, removeButton, saveButton)
	//#

	return FContainer.NewHBox(
		FContainer.New(FLayout.NewCustomPaddedLayout(5, 0, 0, 0), FContainer.NewBorder(nil, editButtons, nil, nil, buttonScroll)),
		FWidget.NewSeparator(),
	)
}

func createProfileButtons(t *ProfileTab) []fyne.CanvasObject {
	buttons := make([]fyne.CanvasObject, len(persistent.Profiles))

	for i, profile := range persistent.Profiles {
		button := FWidget.NewButton(profile.Name, func() {})
		buttons[i] = button
		button.OnTapped = func() {
			//untoggle all buttons
			for j := range buttons {
				buttons[j].(*FWidget.Button).Importance = FWidget.MediumImportance
				buttons[j].Refresh()
			}
			//toggle the tapped button
			button.Importance = FWidget.HighImportance
			button.Refresh()

			setSelectedProfile(t, i)

			updateEditorForm(t)
		}
	}

	// preselect selected button
	if len(persistent.Profiles) > 0 {
		buttons[t.selected_profile].(*FWidget.Button).Importance = FWidget.HighImportance
		buttons[t.selected_profile].Refresh()
		updateEditorForm(t)
	}

	return buttons
}

func setSelectedProfile(t *ProfileTab, i int) {
	t.selected_profile = i
	persistent.Data.Profiles.LastProfileUsed = i
	persistent.SaveData()
}

func updateProfileButtons(t *ProfileTab) {
	t.components.list_btn_profiles.Objects = createProfileButtons(t)
	t.components.list_btn_profiles.Refresh()
}

func nonValidatedEntry(data FBind.String) *FWidget.Entry {
	entry := FWidget.NewEntryWithData(data)
	entry.Validator = nil
	entry.Refresh()
	return entry
}

func createEditorForm(t *ProfileTab) *FWidget.Form {
	form := FWidget.NewForm()

	form.Append("Name", nonValidatedEntry(t.bindings.profile_name))
	form.Append("File Location", FWidget.NewLabelWithData(t.bindings.profile_location))
	form.Append("Note", nonValidatedEntry(t.bindings.profile_note))
	t.components.list_form_targets = FContainer.NewVBox(createTargetFormList(t)...)
	form.Append("Targets", t.components.list_form_targets)
	return form
}

func createTargetForm(t *ProfileTab, target *persistent.Target) *fyne.Container {

	nameBind := FBind.BindString(&target.Name)
	noteBind := FBind.BindString(&target.Note)

	nameBind.AddListener(FBind.NewDataListener(func() {
		target.Name, _ = nameBind.Get()
	}))
	noteBind.AddListener(FBind.NewDataListener(func() {
		target.Note, _ = noteBind.Get()
	}))

	entry_name := nonValidatedEntry(nameBind)
	entry_note := nonValidatedEntry(noteBind)

	list_files := func(list *[]string, t *ProfileTab) fyne.CanvasObject {
		items := make([]fyne.CanvasObject, len(*list))
		for i := 0; i < len(*list); i++ {
			bind := FBind.BindString(&(*list)[i])
			label := FWidget.NewLabelWithData(bind)
			btn_rmv := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func(idx int) func() {
				return func() {
					*list = append((*list)[:idx], (*list)[idx+1:]...)
					updateEditorForm(t)
				}
			}(i))
			items[i] = FContainer.NewBorder(nil, nil, nil, btn_rmv, label)
		}
		btn_add_file := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() {
			dlg_file_selector := FDialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader != nil {
					*list = append(*list, reader.URI().Path())
					reader.Close()
					updateEditorForm(t)
				}
			}, *WindowHandle)
			windowSize := (*WindowHandle).Canvas().Size()
			dlg_file_selector.Resize(fyne.NewSize(windowSize.Width*0.8, windowSize.Height*0.8))
			dlg_file_selector.Show()
		})
		return FContainer.NewBorder(nil, nil, nil, FContainer.NewHBox(FWidget.NewSeparator(), FContainer.NewVBox(btn_add_file, FLayout.NewSpacer())), FContainer.NewVBox(items...))
	}

	return FContainer.New(FLayout.NewFormLayout(),
		FWidget.NewLabel("Name"),
		entry_name,
		FWidget.NewLabel("Source Paths"),
		list_files(&target.SourcePaths, t),
		FWidget.NewLabel("Script Paths"),
		list_files(&target.ScriptPaths, t),
		FWidget.NewLabel("Note"),
		entry_note,
	)
}

func createTargetFormList(t *ProfileTab) []fyne.CanvasObject {
	profile := t.bindings.profile

	btn_add_target := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() {
		profile.Targets = append(profile.Targets, persistent.Target{Name: "New Target"})
		updateEditorForm(t)
	})

	if len(profile.Targets) == 0 {
		return []fyne.CanvasObject{FContainer.NewHBox(FWidget.NewLabel("No targets available."), FLayout.NewSpacer(), btn_add_target)}
	}

	var targetForms []fyne.CanvasObject
	for i := range profile.Targets {
		target := &profile.Targets[i]
		btn_remove_target := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func(idx int) func() {
			return func() {
				profile.Targets = append(profile.Targets[:idx], profile.Targets[idx+1:]...)
				updateEditorForm(t)
			}
		}(i))
		targetForms = append(targetForms,
			FContainer.NewBorder(nil, nil, nil, FContainer.NewVBox(btn_remove_target, FLayout.NewSpacer()),
				createTargetForm(t, target),
			),
		)
		targetForms = append(targetForms, FWidget.NewSeparator())
	}
	targetForms = append(targetForms, FContainer.NewHBox(FLayout.NewSpacer(), btn_add_target))

	return targetForms
}

func updateEditorForm(t *ProfileTab) {
	t.bindings.profile = persistent.Profiles[t.selected_profile]
	t.bindings.profile_name.Set(t.bindings.profile.Name)
	t.bindings.profile_location.Set(persistent.ProfilePath(t.bindings.profile))
	t.bindings.profile_note.Set(t.bindings.profile.Note)

	t.bindings.profile_name.AddListener(FBind.NewDataListener(func() {
		t.bindings.profile.Name, _ = t.bindings.profile_name.Get()
	}))
	t.bindings.profile_note.AddListener(FBind.NewDataListener(func() {
		t.bindings.profile.Note, _ = t.bindings.profile_note.Get()
	}))

	t.components.list_form_targets.Objects = createTargetFormList(t)
	t.components.list_form_targets.Refresh()
}

func addPopup(t *ProfileTab) {
	const (
		NO_SELECTION = 0
		OPEN_FOLDER  = 1
		OPEN_FILE    = 2
	)
	actionMode := NO_SELECTION
	selected_file_or_location := FBind.NewString()
	selected_file_or_location.Set("No file or location selected!")

	var dlg_form *FDialog.CustomDialog
	var btn_selection, btn_cancel, btn_confirm *FWidget.Button

	btn_selection = FWidget.NewButton("Select Location / File", func() {})
	btn_selection.Disable()

	btn_cancel = FWidget.NewButton("Cancel", func() { dlg_form.Hide() })

	btn_confirm = FWidget.NewButton("Confirm", func() {})
	btn_confirm.Disable()
	btn_confirm.Importance = FWidget.HighImportance
	btn_confirm.Refresh()
	btn_confirm.OnTapped = func() {
		path, err := selected_file_or_location.Get()
		if err == nil && path != "" {
			success := false
			var profile *persistent.Profile
			var err error
			if actionMode == OPEN_FOLDER {
				if profile, err = persistent.CreateProfile(path); err == nil {
					fmt.Println("Created new profile at:", path)
					path = filepath.Join(path, "Profile.yaml")
					success = true
				}
			}
			if actionMode == OPEN_FILE {
				if profile, err = persistent.LoadProfile(path); err == nil {
					fmt.Println("Tracking profile from file:", path)
					success = true
				}
			}
			if success {
				persistent.Profiles = append(persistent.Profiles, profile)
				persistent.Data.Profiles.Paths = append(persistent.Data.Profiles.Paths, path)
				persistent.SaveData()

				updateProfileButtons(t)
			}
		} else {
			FDialog.ShowError(err, *WindowHandle)
			return
		}

		dlg_form.Hide()
	}

	btn_selection.OnTapped = func() {
		var dlg_file_selector FDialog.Dialog

		switch actionMode {
		case OPEN_FOLDER:
			dlg_file_selector = FDialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if err == nil && uri != nil {
					selected_file_or_location.Set(uri.Path())
					btn_confirm.Enable()
				}
			}, *WindowHandle)

		case OPEN_FILE:
			dlg_file_selector = FDialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader != nil {
					selected_file_or_location.Set(reader.URI().Path())
					btn_confirm.Enable()
					reader.Close()
				}
			}, *WindowHandle)
		}

		windowSize := (*WindowHandle).Canvas().Size()
		dlg_file_selector.Resize(fyne.NewSize(windowSize.Width*0.8, windowSize.Height*0.8))
		dlg_file_selector.Show()
	}

	formContent := FContainer.New(FLayout.NewCustomPaddedVBoxLayout(20),
		FContainer.New(FLayout.NewCustomPaddedVBoxLayout(-5),
			FWidget.NewLabel("Select Action:"),
			FWidget.NewRadioGroup([]string{"Create a new Profile", "Track an existing Profile"}, func(selected string) {
				var newMode int
				switch selected {
				case "Create a new Profile":
					newMode = OPEN_FOLDER
				case "Track an existing Profile":
					newMode = OPEN_FILE
				}
				if actionMode != newMode {
					actionMode = newMode
					selected_file_or_location.Set("No file or location selected!")
					btn_confirm.Disable()
				}
				btn_selection.Enable()
			}),
		),
		btn_selection,
		FContainer.New(FLayout.NewCustomPaddedVBoxLayout(-10),
			FWidget.NewLabel("Current selection:"),
			FWidget.NewLabelWithData(selected_file_or_location),
		),
		FContainer.NewGridWithColumns(
			2,
			btn_cancel,
			btn_confirm,
		),
	)

	dlg_form = FDialog.NewCustomWithoutButtons(
		"Manage a new Profile",
		formContent,
		*WindowHandle,
	)

	dlg_form.Show()
}

func removePopup(t *ProfileTab) {
	var dlg *FDialog.CustomDialog
	selected_profiles := make([]int, 0)
	btn_cancel := FWidget.NewButton("Cancel", func() { dlg.Hide() })
	btn_confirm := FWidget.NewButton("Confirm", func() {
		for i := 0; i < len(selected_profiles); i++ {
			persistent.Profiles = utils.RemoveIdx(persistent.Profiles, selected_profiles[i])
			persistent.Data.Profiles.Paths = utils.RemoveIdx(persistent.Data.Profiles.Paths, selected_profiles[i])
			if t.selected_profile == selected_profiles[i] {
				setSelectedProfile(t, 0)
			}
		}
		persistent.SaveData()
		updateProfileButtons(t)
		updateEditorForm(t)
		dlg.Hide()
	})

	btns_profiles := func() []fyne.CanvasObject {
		buttons := make([]fyne.CanvasObject, len(persistent.Profiles))
		for i, profile := range persistent.Profiles {
			buttons[i] = FWidget.NewCheck(profile.Name, func(checked bool) {
				if checked {
					selected_profiles = append(selected_profiles, i)
				} else {
					selected_profiles = utils.RemoveItem(selected_profiles, i)
				}
			})
		}
		return buttons
	}

	content := FContainer.New(FLayout.NewCustomPaddedVBoxLayout(20),
		FWidget.NewLabel("Select profiles to remove / untrack:"),
		FContainer.NewVBox(btns_profiles()...),
		FWidget.NewLabelWithStyle("Note: This will not delete the profile files, only untrack them.", fyne.TextAlign(FWidget.ButtonAlignLeading), fyne.TextStyle{Italic: true}),
		FContainer.NewGridWithColumns(
			2,
			btn_cancel,
			btn_confirm,
		),
	)

	dlg = FDialog.NewCustomWithoutButtons(
		"Manage tracked Profiles",
		content,
		*WindowHandle,
	)
	dlg.Show()
}
