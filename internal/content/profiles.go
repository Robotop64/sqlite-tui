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
}

func (t *ProfileTab) CreateContent() *FContainer.TabItem {
	return FContainer.NewTabItem("Profiles", FContainer.NewStack(FContainer.NewBorder(nil, nil, createProfilePanel(t), nil, createEditorForm(t))))
}

func createProfilePanel(t *ProfileTab) *fyne.Container {
	//# list buttons of tracked profiles
	t.components.list_btn_profiles = FContainer.NewVBox(createProfileButtons(t)...)
	buttonScroll := FContainer.NewVScroll(t.components.list_btn_profiles)
	//#

	//# list action / edit buttons
	addButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() { addPopup(t) })
	removeButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func() {})
	saveButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDocumentSave), func() {})
	removeButton.OnTapped = func() {}
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
	buttons := []fyne.CanvasObject{}
	buttons = make([]fyne.CanvasObject, len(persistent.Profiles))

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

			t.selected_profile = i

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

func updateProfileButtons(t *ProfileTab) {
	t.components.list_btn_profiles.Objects = createProfileButtons(t)
	t.components.list_btn_profiles.Refresh()
}

func createEditorForm(t *ProfileTab) *FWidget.Form {
	form := FWidget.NewForm()

	nonValidatedEntry := func(data FBind.String) *FWidget.Entry {
		entry := FWidget.NewEntryWithData(data)
		entry.Validator = nil
		entry.Refresh()
		return entry
	}

	form.Append("Name", nonValidatedEntry(t.bindings.profile_name))
	form.Append("File Location", FWidget.NewLabelWithData(t.bindings.profile_location))
	form.Append("Note", nonValidatedEntry(t.bindings.profile_note))
	t.components.list_form_targets = FContainer.NewVBox(createTargetForm(t))
	form.Append("Targets", t.components.list_form_targets)
	return form
}

func createTargetForm(t *ProfileTab) *fyne.Container {
	targets := persistent.Profiles[t.selected_profile].Targets
	btn_add := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() {
		persistent.Profiles[t.selected_profile].Targets = append(persistent.Profiles[t.selected_profile].Targets, persistent.Target{})
		t.components.list_form_targets.Objects = []fyne.CanvasObject{createTargetForm(t)}
		t.components.list_form_targets.Refresh()
	})
	if len(targets) == 0 {
		return FContainer.NewHBox(FWidget.NewLabel("No targets defined for this profile."), FLayout.NewSpacer(), btn_add)
	}

	targetForm := func(target persistent.Target) *fyne.Container {

		return FContainer.NewVBox()
	}

	forms := make([]fyne.CanvasObject, len(targets))
	for i, target := range targets {
		forms[i] = targetForm(target)
	}

	return FContainer.NewVBox(forms...)
}

// targetEntry := func(target persistent.) *fyne.Container {
// 		return FContainer.NewVBox()
// 	}

// 	targets := make([]*fyne.Container, len(t.bindings.profile.Targets))
// 	for i, target := range t.bindings.profile.Targets {
// 		targets[i] = targetEntry()
// 	}

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
