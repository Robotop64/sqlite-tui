package tabs

import (
	"fmt"
	"path/filepath"
	"strings"

	style "github.com/Robotop64/sqlite-tui/internal/style"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	bubTxtIn "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
	lgTree "github.com/charmbracelet/lipgloss/tree"
	cfg "github.com/spf13/viper"
)

type ProfileTab struct {
	Name                 string
	Profiles             []utils.Profile
	IdxFocus             int
	IdxSelected          int
	AddProfile           bool
	AddProfile_textInput bubTxtIn.Model
}

var profile_popup_width = 40

func (b ProfileTab) PostInit() ProfileTab {
	ti := bubTxtIn.New()
	ti.Placeholder = "Enter profile path..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = profile_popup_width
	b.AddProfile_textInput = ti

	b.IdxSelected = cfg.GetInt("profiles.last_used")
	return b
}

func (b ProfileTab) Init() tea.Cmd {
	return bubTxtIn.Blink
}

func (b ProfileTab) GetName() string {
	return b.Name
}

func (b ProfileTab) View(width, height int) string {
	hint_height := 2
	sidecolumn_width := width / 5
	maincolumn_width := width - sidecolumn_width
	layout_height := height - hint_height

	sidecolumn := style.Box.
		Width(sidecolumn_width - 2*border).
		Height(layout_height - 2*border)
	maincolumn := style.Box.
		Width(maincolumn_width - 2*border).
		Height(layout_height - 2*border)

	tab := sidecolumn.
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(lipgloss.Width(sidecolumn.String())-6, lipgloss.Center, b.Name),
			"⏵",
		)

	var (
		sidebar,
		content,
		hints string
	)

	//=Side Column=============
	//-Explorer----------------
	explorer_height := layout_height - lipgloss.Height(tab.String()) - 2*border

	tree := lgTree.New()
	loadProfiles(tree, &b.Profiles, b.IdxSelected)
	tree.ItemStyleFunc(func(_ lgTree.Children, i int) lipgloss.Style {
		if len(b.Profiles) == 0 {
			return style.Normal
		}

		if i == b.IdxFocus {
			return style.Selected
		}
		return style.Normal
	})

	explorer := sidecolumn.
		Padding(0, 1).
		Height(explorer_height)

	explorer = explorer.SetString(tree.String())
	//-------------------------
	sidebar = lipgloss.JoinVertical(
		lipgloss.Top,
		tab.Render(),
		explorer.Render(),
	)

	//=Hints===================
	hints = lipgloss.JoinHorizontal(
		lipgloss.Top,
		"quit:\nctrl+c|q",
		"│\n│",
		"save:\nctrl+s",
		"│\n│",
		"tabs:\nalt+(</>)",
		"│\n│",
		"profiles:\n   ↑/↓",
		"│\n│",
		"add:\n +",
		"│\n│",
		"remove:\n   -",
		"│\n│",
	)
	//=========================
	//=Content Column==========
	temp_content := maincolumn

	if b.IdxFocus >= 0 && b.IdxFocus < len(b.Profiles) && b.Profiles[b.IdxFocus] != nil {
		title := style.Title.SetString("Profile Properties:").Render()
		list := lgList.New(
			fmt.Sprintf("Name: %s", b.Profiles[b.IdxFocus].GetString("profile.name")),
			fmt.Sprintf("Path: %s", b.Profiles[b.IdxFocus].GetString("profile.path")),
			"Database:",
			lgList.New(
				fmt.Sprintf("Path: %s", b.Profiles[b.IdxFocus].GetString("database.path")),
				fmt.Sprintf("Type: %s", b.Profiles[b.IdxFocus].GetString("database.type")),
			),
		).String()
		temp_content = temp_content.
			Padding(0, 1).
			SetString(
				lipgloss.JoinVertical(
					lipgloss.Top,
					title,
					list,
				),
			)
	}

	content = temp_content.
		Render()
	//=========================

	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			content,
		),
		hints,
	)

	if !b.AddProfile {
		return layout
	} else {
		overlay, err := utils.Overlay(
			layout,
			addProfilePrompt(b).Render(),
			utils.Center, utils.Center,
		)
		if err != nil {
			return fmt.Sprintf("Error overlaying popup: %v", err)
		}
		return overlay
	}
}

func (b ProfileTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch b.AddProfile {
		case true:
			switch msg.String() {
			case "esc":
				b.AddProfile = false
				b.AddProfile_textInput.Reset()
				return b, nil
			case "enter":
				raw_path := b.AddProfile_textInput.Value()
				var path string
				if !utils.EndsWith(raw_path, "Profile.yaml") {
					path = filepath.Join(raw_path, "Profile.yaml")
				} else {
					path = raw_path
				}

				fileExists := utils.CheckPath(path)
				var prof utils.Profile
				var err error
				if !fileExists {
					if prof, err = utils.GenProfile(b.AddProfile_textInput.Value()); err != nil {
						b.AddProfile = false
						b.AddProfile_textInput.Reset()
						return b, nil
					}
				} else {
					if prof, err = utils.LoadProfile(path); err != nil {
						b.AddProfile = false
						b.AddProfile_textInput.Reset()
						return b, nil
					}
				}
				prof.Set("profile.path", path)
				b.Profiles = append(b.Profiles, prof)
				paths := cfg.GetStringSlice("profiles.paths")
				paths = append(paths, path)
				cfg.Set("profiles.paths", paths)

				b.AddProfile = false
				b.AddProfile_textInput.Reset()
				return b, nil
			}
			b.AddProfile_textInput, cmd = b.AddProfile_textInput.Update(msg)
			return b, cmd
		case false:
			switch msg.String() {
			case "up":
				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "down":
				b.IdxFocus = min(b.IdxFocus+1, len(b.Profiles)-1)
				return b, nil
			case "+":
				b.AddProfile = true
				return b, nil
			case "-":
				if len(b.Profiles) > 0 && b.IdxFocus < len(b.Profiles) {
					b.Profiles = append(b.Profiles[:b.IdxFocus], b.Profiles[b.IdxFocus+1:]...)
					paths := cfg.GetStringSlice("profiles.paths")
					paths = append(paths[:b.IdxFocus], paths[b.IdxFocus+1:]...)
					cfg.Set("profiles.paths", paths)
				}
				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "enter":
				cfg.Set("profiles.last_used", b.IdxFocus)
				b.IdxSelected = b.IdxFocus
				return b, nil
			case "c":
				if len(b.Profiles) > 0 && b.IdxFocus < len(b.Profiles) {

				}
				return b, nil
			default:
				return b, nil
			}
		}
	}
	return b, nil
}

func loadProfiles(t *lgTree.Tree, profiles *[]utils.Profile, idxSel int) {
	for i, profile := range *profiles {
		if profile == nil {
			t.Child(fmt.Sprintf("Faulty Profile!\n@Position %d", i+1))
		} else {
			if i == idxSel {
				t.Child(">" + profile.GetString("profile.name") + "<")
			} else {
				t.Child(profile.GetString("profile.name"))
			}
		}
	}

	if len(*profiles) == 0 {
		t.Child("...")
	}
}

func addProfilePrompt(b ProfileTab) lipgloss.Style {
	title := style.Title.SetString("Add Profile:").Render()
	msg := style.Normal.SetString("Enter the path to the profile file:\n").Render()

	label_cancel, key_cancel := " Cancel ", "esc"
	label_confirm, key_confirm := " Confirm ", "enter"

	button_cancel := style.Button.SetString(label_cancel).Render() + "\n" +
		style.Normal.SetString(fmt.Sprintf(" (%s)", key_cancel)).Render()
	button_confirm := style.Button.SetString(label_confirm).Render() + "\n" +
		style.Normal.SetString(fmt.Sprintf(" (%s)", key_confirm)).Render()

	buffer_length := profile_popup_width - 2 - lipgloss.Width(button_cancel) - lipgloss.Width(button_confirm)

	hints := lipgloss.JoinHorizontal(
		lipgloss.Top,
		button_cancel,
		strings.Repeat(" ", buffer_length),
		button_confirm,
	)

	// hints := lipgloss.PlaceHorizontal(
	// 	profile_popup_width-2,
	// 	lipgloss.Right,
	// 	button_confirm,
	// )
	// hints = button_cancel //+ ansi.Cut(hints, ansi.StringWidth(button_cancel), ansi.StringWidth(hints))

	return style.Box.
		Padding(0, 1).
		Width(profile_popup_width).
		SetString(
			lipgloss.JoinVertical(
				lipgloss.Top,
				title,
				msg,
				b.AddProfile_textInput.View(),
				"",
				hints,
			),
		)
}
