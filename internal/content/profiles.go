package content

import (
	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
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
	addButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameContentAdd), func() {})
	removeButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDelete), func() {})
	saveButton := FWidget.NewButtonWithIcon("", FTheme.Icon(FTheme.IconNameDocumentSave), func() {})
	addButton.OnTapped = func() {}
	removeButton.OnTapped = func() {}
	saveButton.OnTapped = func() {}
	editButtons := FContainer.NewVBox(FWidget.NewSeparator(), addButton, removeButton, saveButton)
	//====================================================

	profileList := FContainer.New(FLayout.NewCustomPaddedLayout(5, 0, 0, 0), FContainer.NewBorder(nil, editButtons, nil, nil, buttonScroll))

	profilesPane := FContainer.NewHBox(profileList, FWidget.NewSeparator())
	content := FContainer.NewBorder(nil, nil, profilesPane, nil, editor)
	return FContainer.NewTabItem("Profiles", content)
}
