package content

import (
	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FBind "fyne.io/fyne/v2/data/binding"
	FDialog "fyne.io/fyne/v2/dialog"
	FLayout "fyne.io/fyne/v2/layout"
	FTheme "fyne.io/fyne/v2/theme"
	FWidget "fyne.io/fyne/v2/widget"

	Yaml "gopkg.in/yaml.v3"

	"SQLite-GUI/internal/components"
	"SQLite-GUI/internal/persistent"
)

type ProfileTab struct {
	profiles components.ListModel[*persistent.Profile]
}

func (t *ProfileTab) Init() {
	t.profiles = components.ListModel[*persistent.Profile]{}
	t.profiles.Items = persistent.Profiles
}

func (t *ProfileTab) GenerateContent() *FContainer.TabItem {
	editor := FWidget.NewMultiLineEntry()

	textdata, _ := Yaml.Marshal(t.profiles.Items[t.profiles.Selected])
	editor.SetText(string(textdata))

	//profile buttons =================================
	var buttons components.ListModel[fyne.CanvasObject]
	buttons.Items = make([]fyne.CanvasObject, len(t.profiles.Items))

	for i, profile := range t.profiles.Items {
		button := FWidget.NewButton(profile.Name, func() {})
		buttons.Items[i] = button
		button.OnTapped = func() {
			//untoggle all buttons
			for j := range buttons.Items {
				buttons.Items[j].(*FWidget.Button).Importance = FWidget.MediumImportance
				buttons.Items[j].Refresh()
			}
			//toggle the tapped button
			button.Importance = FWidget.HighImportance
			button.Refresh()

			t.profiles.Selected = i

			textdata, _ := Yaml.Marshal(t.profiles.Items[i])
			editor.SetText(string(textdata))
		}
	}

	//select initial profile/button
	buttons.Items[t.profiles.Selected].(*FWidget.Button).Importance = FWidget.HighImportance
	buttons.Items[t.profiles.Selected].Refresh()
	buttonScroll := FContainer.NewVScroll(FContainer.NewVBox(buttons.Items...))
	//=============================================
	//edit buttons ========================================
	addButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() { addPopup() })
	removeButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func() {})
	saveButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDocumentSave), func() {})
	removeButton.OnTapped = func() {}
	saveButton.OnTapped = func() {}
	editButtons := FContainer.NewVBox(FWidget.NewSeparator(), addButton, removeButton, saveButton)
	//====================================================

	profilesPane := FContainer.NewHBox(
		FContainer.New(FLayout.NewCustomPaddedLayout(5, 0, 0, 0), FContainer.NewBorder(nil, editButtons, nil, nil, buttonScroll)),
		FWidget.NewSeparator(),
	)
	content := FContainer.NewBorder(nil, nil, profilesPane, nil, editor)
	return FContainer.NewTabItem("Profiles", content)
}

func addPopup() {
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
