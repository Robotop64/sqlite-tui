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
	selectedFile_or_Location := FBind.NewString()
	selectedFile_or_Location.Set("No file or location selected!")

	cancelBtn := FWidget.NewButton("Cancel", func() {})
	confirmBtn := FWidget.NewButton("Confirm", func() {})
	confirmBtn.Disable()
	confirmBtn.Importance = FWidget.HighImportance
	confirmBtn.Refresh()

	selectionTrigger := FWidget.NewButton("Select Location / File", func() {
		var fSelector FDialog.Dialog

		switch actionMode {
		case OPEN_FOLDER:
			fSelector = FDialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if err == nil {
					if uri != nil {
						selectedFile_or_Location.Set(uri.Path())
						confirmBtn.Enable()
					}
				}
			}, *WindowHandle)

		case OPEN_FILE:
			fSelector = FDialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil {
					if reader != nil {
						selectedFile_or_Location.Set(reader.URI().Path())
						confirmBtn.Enable()
						reader.Close()
					}
				}
			}, *WindowHandle)
		}

		windowSize := (*WindowHandle).Canvas().Size()
		fSelector.Resize(fyne.NewSize(windowSize.Width*0.8, windowSize.Height*0.8))
		fSelector.Show()
	})
	selectionTrigger.Disable()

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
					selectedFile_or_Location.Set("No file or location selected!")
					confirmBtn.Disable()
				}
				selectionTrigger.Enable()
			}),
		),
		selectionTrigger,
		FContainer.New(FLayout.NewCustomPaddedVBoxLayout(-10),
			FWidget.NewLabel("Current selection:"),
			FWidget.NewLabelWithData(selectedFile_or_Location),
		),
		FContainer.NewGridWithColumns(
			2,
			cancelBtn,
			confirmBtn,
		),
	)

	form := FDialog.NewCustomWithoutButtons(
		"Manage a new Profile",
		formContent,
		*WindowHandle,
	)
	cancelBtn.OnTapped = func() {
		form.Hide()
	}

	form.Show()
}
